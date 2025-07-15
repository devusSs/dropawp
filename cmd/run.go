package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/devusSs/dropawp/internal/config"
	"github.com/devusSs/dropawp/internal/csfloat"
	"github.com/devusSs/dropawp/internal/lastrun"
	"github.com/devusSs/dropawp/internal/secret"
	"github.com/devusSs/dropawp/internal/steam"
	"github.com/devusSs/dropawp/internal/storage"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var lastRun *lastrun.LastRun

// TODO: add logging
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs the application with the given configuration",
	PreRun: func(_ *cobra.Command, _ []string) {
		var err error
		cfg, err = config.Read()
		cobra.CheckErr(err)

		lastRun, err = lastrun.Read()
		if err != nil {
			if !errors.Is(err, lastrun.ErrLastRunNotExist) {
				cobra.CheckErr(fmt.Sprintf("failed to read last run: %v", err))
			}

			lastRun = &lastrun.LastRun{ProjectName: cfg.ProjectName}
		}
	},
	Run: func(_ *cobra.Command, _ []string) {
		if lastRun.ProjectName != cfg.ProjectName {
			cobra.CheckErr(
				fmt.Errorf(
					"last run project name '%s' does not match current project name '%s'",
					lastRun.ProjectName,
					cfg.ProjectName,
				),
			)
		}

		if time.Since(lastRun.LastRun) < cfg.CooldownDuration {
			cobra.CheckErr(
				fmt.Errorf(
					"cooldown period not met, last run was at %s, cooldown duration is %s",
					lastRun.LastRun.Format(time.RFC3339),
					cfg.CooldownDuration,
				),
			)
		}

		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()

		if !cfg.SkipSteamServicesCheck {
			apiKey, err := getSecret(secret.SteamAPIKey)
			cobra.CheckErr(err)

			var status *steam.CSServerStatus
			status, err = steam.GetCSServerStatus(ctx, apiKey)
			cobra.CheckErr(err)

			if status.Community == "offline" || status.Sessions == "offline" {
				cobra.CheckErr(
					fmt.Sprintf(
						"required Steam services have issues: community: %s, sessions: %s",
						status.Community,
						status.Sessions,
					),
				)
			}
		}

		if !cfg.SkipSteamUserCheck {
			apiKey, err := getSecret(secret.SteamAPIKey)
			cobra.CheckErr(err)

			var user *steam.UserSummary
			user, err = steam.GetUserSummary(ctx, apiKey, cfg.SteamID64)
			cobra.CheckErr(err)

			if user.CommunityVisibilityState != steam.CommunityVisibilityPublic {
				cobra.CheckErr(
					fmt.Sprintf(
						"steam user profile is not public, visibility state: %s",
						user.CommunityVisibilityState.String(),
					),
				)
			}

			if user.ProfileState != steam.ProfileStateCreated {
				cobra.CheckErr(
					fmt.Sprintf(
						"steam user profile is not set up, profile state: %s",
						user.ProfileState.String(),
					),
				)
			}
		}

		inv, err := steam.GetCSInventory(ctx, cfg.SteamID64)
		cobra.CheckErr(err)

		items := inv.MarketableItems
		if !cfg.SkipFilterUntradableItems {
			items = inv.MarketableAndTradableItems
		}

		var additionalItems *additionalItems
		additionalItems, err = loadAdditionalItemsFile()
		cobra.CheckErr(err)

		for item := range additionalItems.Items {
			items = append(items, steam.CSInventoryItem{
				MarketHashName: item,
			})
		}

		var apiKey string
		apiKey, err = getSecret(secret.CSFloatAPIKey)
		cobra.CheckErr(err)

		itemsPriceMap := make(map[string]int)
		itemsAmountMap := make(map[string]int)
		mutex := &sync.Mutex{}
		itemsWg := &sync.WaitGroup{}
		itemsNoPrice := make(map[string]string)

		for _, item := range items {
			itemsWg.Add(1)
			go func() {
				defer itemsWg.Done()

				var price int
				price, err = csfloat.GetMedianItemPrice(ctx, apiKey, item.MarketHashName)
				if err != nil {
					mutex.Lock()
					defer mutex.Unlock()
					itemsNoPrice[item.MarketHashName] = err.Error()
					return
				}

				mutex.Lock()
				defer mutex.Unlock()

				itemsPriceMap[item.MarketHashName] = price
				itemsAmountMap[item.MarketHashName]++
			}()
		}

		itemsWg.Wait()

		if runPrintResults {
			err = printItemMap(itemsPriceMap, itemsAmountMap)
			cobra.CheckErr(err)

			fmt.Println()
			if len(itemsNoPrice) > 0 {
				fmt.Println("Items with no price:")
				for item, reason := range itemsNoPrice {
					fmt.Println("-", item, ":", reason)
				}
			}
		}

		if runExitOnNoPrice && len(itemsNoPrice) > 0 {
			cobra.CheckErr(
				fmt.Errorf(
					"some items have no price: %v",
					itemsNoPrice,
				),
			)
		}

		if len(itemsNoPrice) == len(items) {
			cobra.CheckErr(
				"all items have no price, check network conditions",
			)
		}

		storageItems := make([]storage.InventoryItem, 0, len(itemsPriceMap))
		for _, item := range items {
			price, ok := itemsPriceMap[item.MarketHashName]
			if !ok {
				continue
			}

			amount, ok := itemsAmountMap[item.MarketHashName]
			if !ok {
				continue
			}

			storageItems = append(
				storageItems,
				storage.InventoryItem{
					IconURL:           item.IconURL,
					ActionInspectLink: item.ActionInspectLink,
					Name:              item.Name,
					NameColor:         item.NameColor,
					MarketName:        item.MarketName,
					MarketHashName:    item.MarketHashName,
					MarketInspectLink: item.MarketInspectLink,
					Marketable:        item.Marketable,
					Tradable:          item.Tradable,
					Amount:            amount,
					Price:             price,
					Currency:          "USD",
				},
			)
		}

		err = storage.Write(cfg.ProjectName, storageItems)
		cobra.CheckErr(err)

		err = lastrun.Write(cfg.ProjectName)
		cobra.CheckErr(err)

		fmt.Println("Run completed successfully")
	},
}

var (
	runPrintResults  bool
	runExitOnNoPrice bool
)

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().BoolVar(&runPrintResults, "print-results", false, "print results to stdout")
	runCmd.Flags().
		BoolVar(&runExitOnNoPrice, "exit-on-no-price", false, "exit if any item has no price")
}

type additionalItems struct {
	Items map[string]int `json:"items"`
}

func (i *additionalItems) String() string {
	return fmt.Sprintf("%+v", *i)
}

func loadAdditionalItemsFile() (*additionalItems, error) {
	if cfg.AdditionalItemsFile == "" {
		return &additionalItems{}, nil
	}

	f, err := os.Open(cfg.AdditionalItemsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open additional items file: %w", err)
	}
	defer f.Close()

	var items additionalItems
	err = json.NewDecoder(f).Decode(&items)
	if err != nil {
		return nil, fmt.Errorf("failed to decode additional items file: %w", err)
	}

	if len(items.Items) == 0 {
		return nil, errors.New("additional items file contains no items")
	}

	return &items, nil
}

const priceConversionFactor = 100

func printItemMap(pricesMap map[string]int, amountsMap map[string]int) error {
	if len(pricesMap) == 0 {
		return errors.New("no items with prices to print")
	}

	if len(amountsMap) == 0 {
		return errors.New("no items with amounts to print")
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"Item", "Price (USD)", "Amount", "Total Price (USD)"})

	for item, price := range pricesMap {
		amount := amountsMap[item]
		totalPrice := float64(price) * float64(amount) / priceConversionFactor
		err := table.Append(
			[]string{
				item,
				fmt.Sprintf("%.2f", float64(price)/priceConversionFactor),
				strconv.Itoa(amount),
				fmt.Sprintf("%.2f", totalPrice),
			},
		)
		if err != nil {
			return fmt.Errorf("failed to append item to table: %w", err)
		}
	}

	return table.Render()
}

package steam

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type CS2Inventory struct {
	AllItems                   []CSInventoryItem `json:"all_items"`
	MarketableItems            []CSInventoryItem `json:"marketable_items"`
	MarketableAndTradableItems []CSInventoryItem `json:"marketable_and_tradable_items"`
}

func (i *CS2Inventory) String() string {
	return fmt.Sprintf("%+v", *i)
}

func (i *CS2Inventory) UniqueItemsCount() map[string]int {
	counts := make(map[string]int)
	for _, item := range i.AllItems {
		counts[item.MarketHashName]++
	}
	return counts
}

func (i *CS2Inventory) UniqueMarketableItemsCount() map[string]int {
	counts := make(map[string]int)
	for _, item := range i.MarketableItems {
		counts[item.MarketHashName]++
	}
	return counts
}

func (i *CS2Inventory) UniqueMarketableAndTradableItemsCount() map[string]int {
	counts := make(map[string]int)
	for _, item := range i.MarketableAndTradableItems {
		counts[item.MarketHashName]++
	}
	return counts
}

type CSInventoryItem struct {
	Name            string `json:"name"`
	Marketable      bool   `json:"marketable"`
	MarketName      string `json:"market_name"`
	MarketHashName  string `json:"market_hash_name"`
	Tradable        bool   `json:"tradable"`
	NameColor       string `json:"name_color"`
	BackgroundColor string `json:"background_color"`
	IconURL         string `json:"icon_url"`
	InspectLink     string `json:"inspect_link"`
}

func (i CSInventoryItem) String() string {
	return fmt.Sprintf(
		"Name: %s, Marketable: %t, MarketName: %s, MarketHashName: %s, Tradable: %t, NameColor: %s, BackgroundColor: %s, IconURL: %s, InspectLink: %s",
		i.Name,
		i.Marketable,
		i.MarketName,
		i.MarketHashName,
		i.Tradable,
		i.NameColor,
		i.BackgroundColor,
		i.IconURL,
		i.InspectLink,
	)
}

func LoadCS2Inventory(ctx context.Context, steamUserID64 string) (*CS2Inventory, error) {
	if ctx == nil {
		return nil, errors.New("context cannot be nil")
	}

	if steamUserID64 == "" {
		return nil, errors.New("steamUserID64 cannot be empty")
	}

	u, err := url.Parse(cs2InventoryBaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	u.Path, err = url.JoinPath(
		u.Path,
		steamUserID64,
		strconv.Itoa(cs2InventoryAppID),
		strconv.Itoa(cs2InventoryContextID),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to join path: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	res := &cs2InventoryResponse{}
	err = json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if res.Success != 1 {
		return nil, fmt.Errorf("failed to load CS2 inventory: success code is %d", res.Success)
	}

	return res.toCS2Inventory(), nil
}

const (
	cs2InventoryBaseURL   = "https://steamcommunity.com/inventory/"
	cs2InventoryAppID     = 730
	cs2InventoryContextID = 2
)

type cs2InventoryResponse struct {
	Assets []struct {
		Appid      int    `json:"appid"`
		Contextid  string `json:"contextid"`
		Assetid    string `json:"assetid"`
		Classid    string `json:"classid"`
		Instanceid string `json:"instanceid"`
		Amount     string `json:"amount"`
	} `json:"assets"`
	Descriptions []struct {
		Appid           int    `json:"appid"`
		Classid         string `json:"classid"`
		Instanceid      string `json:"instanceid"`
		Currency        int    `json:"currency"`
		BackgroundColor string `json:"background_color"`
		IconURL         string `json:"icon_url"`
		Descriptions    []struct {
			Type  string `json:"type"`
			Value string `json:"value"`
			Name  string `json:"name"`
			Color string `json:"color,omitempty"`
		} `json:"descriptions"`
		Tradable int `json:"tradable"`
		Actions  []struct {
			Link string `json:"link"`
			Name string `json:"name"`
		} `json:"actions,omitempty"`
		OwnerDescriptions []struct {
			Type  string `json:"type"`
			Value string `json:"value"`
			Color string `json:"color,omitempty"`
		} `json:"owner_descriptions,omitempty"`
		Name           string `json:"name"`
		NameColor      string `json:"name_color"`
		Type           string `json:"type"`
		MarketName     string `json:"market_name"`
		MarketHashName string `json:"market_hash_name"`
		MarketActions  []struct {
			Link string `json:"link"`
			Name string `json:"name"`
		} `json:"market_actions,omitempty"`
		Commodity                   int `json:"commodity"`
		MarketTradableRestriction   int `json:"market_tradable_restriction"`
		MarketMarketableRestriction int `json:"market_marketable_restriction"`
		Marketable                  int `json:"marketable"`
		Tags                        []struct {
			Category              string `json:"category"`
			InternalName          string `json:"internal_name"`
			LocalizedCategoryName string `json:"localized_category_name"`
			LocalizedTagName      string `json:"localized_tag_name"`
			Color                 string `json:"color,omitempty"`
		} `json:"tags"`
	} `json:"descriptions"`
	TotalInventoryCount int `json:"total_inventory_count"`
	Success             int `json:"success"`
	Rwgrsn              int `json:"rwgrsn"`
}

const cs2InventoryIconBaseURL = "https://community.fastly.steamstatic.com/economy/image/"

func (r *cs2InventoryResponse) toCS2Inventory() *CS2Inventory {
	i := &CS2Inventory{
		AllItems:                   make([]CSInventoryItem, 0, len(r.Descriptions)),
		MarketableItems:            make([]CSInventoryItem, 0),
		MarketableAndTradableItems: make([]CSInventoryItem, 0),
	}

	for _, desc := range r.Descriptions {
		item := CSInventoryItem{
			Name:            desc.Name,
			Marketable:      desc.Marketable == 1,
			MarketName:      desc.MarketName,
			MarketHashName:  desc.MarketHashName,
			Tradable:        desc.Tradable == 1,
			NameColor:       desc.NameColor,
			BackgroundColor: desc.BackgroundColor,
			IconURL:         cs2InventoryIconBaseURL + desc.IconURL,
			InspectLink:     "",
		}

		for _, action := range desc.Actions {
			if action.Name == "Inspect in Game..." {
				item.InspectLink = action.Link
				break
			}
		}

		i.AllItems = append(i.AllItems, item)

		if item.Marketable {
			i.MarketableItems = append(i.MarketableItems, item)
		}

		if item.Marketable && item.Tradable {
			i.MarketableAndTradableItems = append(i.MarketableAndTradableItems, item)
		}
	}

	return i
}

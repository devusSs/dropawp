package steam

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type CSInventory struct {
	AllItems                   []CSInventoryItem `json:"all_items"`
	MarketableItems            []CSInventoryItem `json:"marketable_items"`
	MarketableAndTradableItems []CSInventoryItem `json:"marketable_and_tradable_items"`
}

func (i *CSInventory) String() string {
	return fmt.Sprintf("%+v", *i)
}

type CSInventoryItem struct {
	IconURL           string `json:"icon_url"`
	ActionInspectLink string `json:"inspect_url"`
	Name              string `json:"name"`
	NameColor         string `json:"name_color"`
	MarketName        string `json:"market_name"`
	MarketHashName    string `json:"market_hash_name"`
	MarketInspectLink string `json:"market_inspect_link"`
	Marketable        bool   `json:"marketable"`
	Tradable          bool   `json:"tradable"`
}

func (i CSInventoryItem) String() string {
	return fmt.Sprintf(
		"CSInventoryItem{IconURL: %s, ActionInspectLink: %s, Name: %s, NameColor: %s, MarketName: %s, MarketHashName: %s, MarketInspectLink: %s, Marketable: %t, Tradable: %t}",
		i.IconURL,
		i.ActionInspectLink,
		i.Name,
		i.NameColor,
		i.MarketName,
		i.MarketHashName,
		i.MarketInspectLink,
		i.Marketable,
		i.Tradable,
	)
}

func GetCSInventory(ctx context.Context, steamID64 uint64) (*CSInventory, error) {
	if ctx == nil {
		return nil, ErrContextNil
	}

	err := validateSteamID64(steamID64)
	if err != nil {
		return nil, fmt.Errorf("invalid steamID64: %w", err)
	}

	inventoryURL := fmt.Sprintf(inventoryURLFormat, steamID64)

	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, inventoryURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	err = applyHeaders(req)
	if err != nil {
		return nil, fmt.Errorf("failed to apply headers: %w", err)
	}

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	res := &csInventoryResponse{}
	err = json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if res.Success != 1 {
		return nil, fmt.Errorf("failed to retrieve inventory: success code %d", res.Success)
	}

	return res.toCSInventory(), nil
}

const inventoryURLFormat = "https://steamcommunity.com/inventory/%d/730/2"

type csInventoryResponse struct {
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

const iconURLBase = "https://community.fastly.steamstatic.com/economy/image/"

func (r *csInventoryResponse) toCSInventory() *CSInventory {
	i := &CSInventory{
		AllItems:                   make([]CSInventoryItem, 0, len(r.Descriptions)),
		MarketableItems:            make([]CSInventoryItem, 0),
		MarketableAndTradableItems: make([]CSInventoryItem, 0),
	}

	for _, desc := range r.Descriptions {
		item := CSInventoryItem{
			IconURL:        iconURLBase + desc.IconURL,
			Name:           desc.Name,
			NameColor:      desc.NameColor,
			MarketName:     desc.MarketName,
			MarketHashName: desc.MarketHashName,
			Marketable:     desc.Marketable == 1,
			Tradable:       desc.Tradable == 1,
		}

		for _, action := range desc.Actions {
			if action.Name == "Inspect in Game..." {
				item.ActionInspectLink = action.Link
				break
			}
		}

		for _, action := range desc.MarketActions {
			if action.Name == "Inspect in Game..." {
				item.MarketInspectLink = action.Link
				break
			}
		}

		i.AllItems = append(i.AllItems, item)

		if item.Marketable && item.Tradable {
			i.MarketableAndTradableItems = append(i.MarketableAndTradableItems, item)
		}

		if item.Marketable {
			i.MarketableItems = append(i.MarketableItems, item)
		}
	}

	return i
}

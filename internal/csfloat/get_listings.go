package csfloat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"time"
)

// TODO: add additional parameters like float etc. to narrow down listings
func GetMedianItemPrice(ctx context.Context, apiKey string, marketHashName string) (int, error) {
	if ctx == nil {
		return 0, ErrContextNil
	}

	if apiKey == "" {
		return 0, errors.New("apiKey cannot be empty")
	}

	if marketHashName == "" {
		return 0, errors.New("marketHashName cannot be empty")
	}

	u, err := url.Parse(getAllListingsURL)
	if err != nil {
		return 0, fmt.Errorf("failed to parse URL: %w", err)
	}

	q := u.Query()
	q.Set("limit", "50")
	q.Set("page", "0")
	q.Set("sort_by", "lowest_price")
	q.Set("market_hash_name", marketHashName)
	u.RawQuery = q.Encode()

	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	err = applyHeaders(req, apiKey)
	if err != nil {
		return 0, fmt.Errorf("failed to apply headers: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API returned non-OK status: %s", resp.Status)
	}

	var listingsResponse getListingsResponse
	err = json.NewDecoder(resp.Body).Decode(&listingsResponse)
	if err != nil {
		return 0, fmt.Errorf("failed to decode response body: %w", err)
	}

	if len(listingsResponse.Data) == 0 {
		return 0, errors.New("no listings found for the given market hash name")
	}

	var prices []int
	for _, listing := range listingsResponse.Data {
		if listing.State == "listed" || listing.State == "buffered" {
			prices = append(prices, listing.Price)
		}
	}

	if len(prices) == 0 {
		return 0, errors.New("no active listings found to calculate median price")
	}

	sort.Ints(prices)

	const medianPriceDivisionFactor = 2

	mid := len(prices) / medianPriceDivisionFactor
	if len(prices)%2 == 0 {
		return (prices[mid-1] + prices[mid]) / medianPriceDivisionFactor, nil
	}

	return prices[mid], nil
}

const getAllListingsURL = "https://csfloat.com/api/v1/listings"

type getListingsResponse struct {
	Data []listing `json:"data"`
}

type listing struct {
	ID               string    `json:"id"`
	CreatedAt        time.Time `json:"created_at"`
	Type             string    `json:"type"`
	Price            int       `json:"price"`
	State            string    `json:"state"`
	Seller           seller    `json:"seller,omitempty"`
	Reference        reference `json:"reference"`
	Item             item      `json:"item,omitempty"`
	IsSeller         bool      `json:"is_seller"`
	IsWatchlisted    bool      `json:"is_watchlisted"`
	Watchers         int       `json:"watchers"`
	MinOfferPrice    int       `json:"min_offer_price,omitempty"`
	MaxOfferDiscount int       `json:"max_offer_discount,omitempty"`
}

type seller struct {
	Avatar              string `json:"avatar"`
	Away                bool   `json:"away"`
	Flags               int    `json:"flags"`
	HasValidSteamAPIKey bool   `json:"has_valid_steam_api_key"`
	ObfuscatedID        string `json:"obfuscated_id"`
	Online              bool   `json:"online"`
	StallPublic         bool   `json:"stall_public"`
	Statistics          struct {
		MedianTradeTime     int `json:"median_trade_time"`
		TotalAvoidedTrades  int `json:"total_avoided_trades"`
		TotalFailedTrades   int `json:"total_failed_trades"`
		TotalTrades         int `json:"total_trades"`
		TotalVerifiedTrades int `json:"total_verified_trades"`
	} `json:"statistics"`
	SteamID  string `json:"steam_id"`
	Username string `json:"username"`
}

type reference struct {
	BasePrice      int       `json:"base_price"`
	FloatFactor    float64   `json:"float_factor"`
	PredictedPrice int       `json:"predicted_price"`
	Quantity       int       `json:"quantity"`
	LastUpdated    time.Time `json:"last_updated"`
}

type item struct {
	AssetID           string    `json:"asset_id"`
	DefIndex          int       `json:"def_index"`
	PaintIndex        int       `json:"paint_index"`
	PaintSeed         int       `json:"paint_seed"`
	FloatValue        float64   `json:"float_value"`
	IconURL           string    `json:"icon_url"`
	DParam            string    `json:"d_param"`
	IsStattrak        bool      `json:"is_stattrak"`
	IsSouvenir        bool      `json:"is_souvenir"`
	Rarity            int       `json:"rarity"`
	Quality           int       `json:"quality"`
	MarketHashName    string    `json:"market_hash_name"`
	Stickers          []sticker `json:"stickers,omitempty"`
	HighRank          int       `json:"high_rank,omitempty"`
	LowRank           int       `json:"low_rank,omitempty"`
	Tradable          int       `json:"tradable"`
	InspectLink       string    `json:"inspect_link"`
	Cs2ScreenshotID   string    `json:"cs2_screenshot_id"`
	Cs2ScreenshotAt   time.Time `json:"cs2_screenshot_at"`
	IsCommodity       bool      `json:"is_commodity"`
	Type              string    `json:"type"`
	RarityName        string    `json:"rarity_name"`
	TypeName          string    `json:"type_name"`
	ItemName          string    `json:"item_name"`
	WearName          string    `json:"wear_name"`
	Description       string    `json:"description"`
	Collection        string    `json:"collection"`
	SerializedInspect string    `json:"serialized_inspect"`
	GsSig             string    `json:"gs_sig"`
}

type sticker struct {
	StickerID int    `json:"stickerId"`
	Slot      int    `json:"slot"`
	IconURL   string `json:"icon_url"`
	Name      string `json:"name"`
	Reference struct {
		Price     int       `json:"price"`
		Quantity  int       `json:"quantity"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"reference"`
	OffsetX  float64 `json:"offset_x,omitempty"`
	OffsetY  float64 `json:"offset_y,omitempty"`
	Rotation int     `json:"rotation,omitempty"`
}

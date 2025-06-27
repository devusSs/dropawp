package steam

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type CS2ServerStatus struct {
	AppVersion     int       `json:"app_version"`
	AppTimestamp   time.Time `json:"app_timestamp"`
	SessionsLogon  string    `json:"sessions_logon"`
	SteamCommunity string    `json:"steam_community"`
}

func (s *CS2ServerStatus) String() string {
	return fmt.Sprintf("%+v", *s)
}

func (s *CS2ServerStatus) IsAvailable() bool {
	return s.SessionsLogon != "offline" && s.SteamCommunity != "offline"
}

func GetCS2ServerStatus(ctx context.Context, apiKey string) (*CS2ServerStatus, error) {
	if ctx == nil {
		return nil, errors.New("context cannot be nil")
	}

	err := validateSteamAPIKey(apiKey)
	if err != nil {
		return nil, fmt.Errorf("invalid Steam API key: %w", err)
	}

	u, err := url.Parse(getCS2ServerStatusURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	query := u.Query()
	query.Set("key", apiKey)
	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	res := &cs2ServerStatusResponse{}
	err = json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return res.toCS2ServerStatus(), nil
}

const getCS2ServerStatusURL = "https://api.steampowered.com/ICSGOServers_730/GetGameServersStatus/v1/"

type cs2ServerStatusResponse struct {
	Result struct {
		App struct {
			Version   int    `json:"version"`
			Timestamp int    `json:"timestamp"`
			Time      string `json:"time"`
		} `json:"app"`
		Services struct {
			SessionsLogon  string `json:"SessionsLogon"`
			SteamCommunity string `json:"SteamCommunity"`
			IEconItems     string `json:"IEconItems"`
			Leaderboards   string `json:"Leaderboards"`
		} `json:"services"`
		Datacenters struct {
			Peru struct {
				Capacity string `json:"capacity"`
				Load     string `json:"load"`
			} `json:"Peru"`
			EUGermany struct {
				Capacity string `json:"capacity"`
				Load     string `json:"load"`
			} `json:"EU Germany"`
			EUAustria struct {
				Capacity string `json:"capacity"`
				Load     string `json:"load"`
			} `json:"EU Austria"`
			EUPoland struct {
				Capacity string `json:"capacity"`
				Load     string `json:"load"`
			} `json:"EU Poland"`
			HongKong struct {
				Capacity string `json:"capacity"`
				Load     string `json:"load"`
			} `json:"Hong Kong"`
			EUSpain struct {
				Capacity string `json:"capacity"`
				Load     string `json:"load"`
			} `json:"EU Spain"`
			Chile struct {
				Capacity string `json:"capacity"`
				Load     string `json:"load"`
			} `json:"Chile"`
			USCalifornia struct {
				Capacity string `json:"capacity"`
				Load     string `json:"load"`
			} `json:"US California"`
			USAtlanta struct {
				Capacity string `json:"capacity"`
				Load     string `json:"load"`
			} `json:"US Atlanta"`
			ChinaGuangdong struct {
				Capacity string `json:"capacity"`
				Load     string `json:"load"`
			} `json:"China Guangdong"`
			EUSweden struct {
				Capacity string `json:"capacity"`
				Load     string `json:"load"`
			} `json:"EU Sweden"`
			Emirates struct {
				Capacity string `json:"capacity"`
				Load     string `json:"load"`
			} `json:"Emirates"`
			USSeattle struct {
				Capacity string `json:"capacity"`
				Load     string `json:"load"`
			} `json:"US Seattle"`
			SouthAfrica struct {
				Capacity string `json:"capacity"`
				Load     string `json:"load"`
			} `json:"South Africa"`
			Brazil struct {
				Capacity string `json:"capacity"`
				Load     string `json:"load"`
			} `json:"Brazil"`
			USVirginia struct {
				Capacity string `json:"capacity"`
				Load     string `json:"load"`
			} `json:"US Virginia"`
			USChicago struct {
				Capacity string `json:"capacity"`
				Load     string `json:"load"`
			} `json:"US Chicago"`
			Japan struct {
				Capacity string `json:"capacity"`
				Load     string `json:"load"`
			} `json:"Japan"`
			ChinaPudong struct {
				Capacity string `json:"capacity"`
				Load     string `json:"load"`
			} `json:"China Pudong"`
			EUFinland struct {
				Capacity string `json:"capacity"`
				Load     string `json:"load"`
			} `json:"EU Finland"`
			IndiaMumbai struct {
				Capacity string `json:"capacity"`
				Load     string `json:"load"`
			} `json:"India Mumbai"`
			IndiaChennai struct {
				Capacity string `json:"capacity"`
				Load     string `json:"load"`
			} `json:"India Chennai"`
			USDallas struct {
				Capacity string `json:"capacity"`
				Load     string `json:"load"`
			} `json:"US Dallas"`
			Argentina struct {
				Capacity string `json:"capacity"`
				Load     string `json:"load"`
			} `json:"Argentina"`
			SouthKorea struct {
				Capacity string `json:"capacity"`
				Load     string `json:"load"`
			} `json:"South Korea"`
			UnitedKingdom struct {
				Capacity string `json:"capacity"`
				Load     string `json:"load"`
			} `json:"United Kingdom"`
			Singapore struct {
				Capacity string `json:"capacity"`
				Load     string `json:"load"`
			} `json:"Singapore"`
			Australia struct {
				Capacity string `json:"capacity"`
				Load     string `json:"load"`
			} `json:"Australia"`
			ChinaBeijing struct {
				Capacity string `json:"capacity"`
				Load     string `json:"load"`
			} `json:"China Beijing"`
			ChinaChengdu struct {
				Capacity string `json:"capacity"`
				Load     string `json:"load"`
			} `json:"China Chengdu"`
		} `json:"datacenters"`
		Matchmaking struct {
			Scheduler        string `json:"scheduler"`
			OnlineServers    int    `json:"online_servers"`
			OnlinePlayers    int    `json:"online_players"`
			SearchingPlayers int    `json:"searching_players"`
			SearchSecondsAvg int    `json:"search_seconds_avg"`
		} `json:"matchmaking"`
		Perfectworld struct {
			Logon struct {
				Availability string `json:"availability"`
				Latency      string `json:"latency"`
			} `json:"logon"`
			Purchase struct {
				Availability string `json:"availability"`
				Latency      string `json:"latency"`
			} `json:"purchase"`
		} `json:"perfectworld"`
	} `json:"result"`
}

func (s *cs2ServerStatusResponse) toCS2ServerStatus() *CS2ServerStatus {
	return &CS2ServerStatus{
		AppVersion:     s.Result.App.Version,
		AppTimestamp:   time.Unix(int64(s.Result.App.Timestamp), 0),
		SessionsLogon:  strings.ToLower(s.Result.Services.SessionsLogon),
		SteamCommunity: strings.ToLower(s.Result.Services.SteamCommunity),
	}
}

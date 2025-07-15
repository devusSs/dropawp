package steam

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type CSServerStatus struct {
	Version   int       `json:"version"`
	Timestamp time.Time `json:"timestamp"`
	Sessions  string    `json:"sessions"`
	Community string    `json:"community"`
}

func (s *CSServerStatus) String() string {
	return fmt.Sprintf("%+v", *s)
}

func GetCSServerStatus(ctx context.Context, apiKey string) (*CSServerStatus, error) {
	if ctx == nil {
		return nil, ErrContextNil
	}

	err := validateSteamAPIKey(apiKey)
	if err != nil {
		return nil, fmt.Errorf("invalid steam api key: %w", err)
	}

	var u *url.URL
	u, err = url.Parse(csServerStatusURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}

	q := u.Query()
	q.Set("key", apiKey)
	u.RawQuery = q.Encode()

	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
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
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	res := &csServerStatusResponse{}
	err = json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return res.toCSServerStatus(), nil
}

const csServerStatusURL = "https://api.steampowered.com/ICSGOServers_730/GetGameServersStatus/v1/"

type csServerStatusResponse struct {
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

func (r *csServerStatusResponse) toCSServerStatus() *CSServerStatus {
	return &CSServerStatus{
		Version:   r.Result.App.Version,
		Timestamp: time.Unix(int64(r.Result.App.Timestamp), 0),
		Sessions:  r.Result.Services.SessionsLogon,
		Community: r.Result.Services.SteamCommunity,
	}
}

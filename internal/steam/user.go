package steam

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const StateUnknown = "Unknown"

type CommunityVisibilityState int

const (
	CommunityVisibilityPrivate     CommunityVisibilityState = 1
	CommunityVisibilityFriendsOnly CommunityVisibilityState = 2
	CommunityVisibilityPublic      CommunityVisibilityState = 3
)

func (c CommunityVisibilityState) String() string {
	switch c {
	case CommunityVisibilityPrivate:
		return "Private"
	case CommunityVisibilityFriendsOnly:
		return "Friends Only"
	case CommunityVisibilityPublic:
		return "Public"
	default:
		return StateUnknown
	}
}

type ProfileState int

const (
	ProfileStateNone    ProfileState = 0
	ProfileStateCreated ProfileState = 1
)

func (p ProfileState) String() string {
	switch p {
	case ProfileStateNone:
		return "None"
	case ProfileStateCreated:
		return "Created"
	default:
		return StateUnknown
	}
}

type CommentPermission int

const (
	CommentPermissionPrivate     CommentPermission = 0
	CommentPermissionFriendsOnly CommentPermission = 1
	CommentPermissionPublic      CommentPermission = 2
)

func (c CommentPermission) String() string {
	switch c {
	case CommentPermissionPrivate:
		return "Private"
	case CommentPermissionFriendsOnly:
		return "Friends Only"
	case CommentPermissionPublic:
		return "Public"
	default:
		return StateUnknown
	}
}

type PersonaState int

const (
	PersonaStateOffline        PersonaState = 0
	PersonaStateOnline         PersonaState = 1
	PersonaStateBusy           PersonaState = 2
	PersonaStateAway           PersonaState = 3
	PersonaStateSnooze         PersonaState = 4
	PersonaStateLookingToTrade PersonaState = 5
	PersonaStateLookingToPlay  PersonaState = 6
)

func (p PersonaState) String() string {
	switch p {
	case PersonaStateOffline:
		return "Offline"
	case PersonaStateOnline:
		return "Online"
	case PersonaStateBusy:
		return "Busy"
	case PersonaStateAway:
		return "Away"
	case PersonaStateSnooze:
		return "Snooze"
	case PersonaStateLookingToTrade:
		return "Looking to Trade"
	case PersonaStateLookingToPlay:
		return "Looking to Play"
	default:
		return StateUnknown
	}
}

type UserSummary struct {
	SteamID64                uint64                   `json:"steam_id_64"`
	CommunityVisibilityState CommunityVisibilityState `json:"community_visibility_state"`
	ProfileState             ProfileState             `json:"profile_state"`
	PersonaName              string                   `json:"persona_name"`
	CommentPermission        CommentPermission        `json:"comment_permission"`
	ProfileURL               string                   `json:"profile_url"`
	Avatar                   string                   `json:"avatar"`
	AvatarMedium             string                   `json:"avatar_medium"`
	AvatarFull               string                   `json:"avatar_full"`
	AvatarHash               string                   `json:"avatar_hash"`
	LastLogoff               time.Time                `json:"last_logoff"`
	PersonaState             PersonaState             `json:"persona_state"`
	PrimaryClanID            string                   `json:"primary_clan_id"`
	TimeCreated              time.Time                `json:"time_created"`
	PersonaStateFlags        int                      `json:"persona_state_flags"`
	LocCountryCode           string                   `json:"loc_country_code"`
	LocStateCode             string                   `json:"loc_state_code"`
}

func (s *UserSummary) String() string {
	return fmt.Sprintf("%+v", *s)
}

func GetUserSummary(ctx context.Context, apiKey string, steamID64 uint64) (*UserSummary, error) {
	if ctx == nil {
		return nil, ErrContextNil
	}

	err := validateSteamAPIKey(apiKey)
	if err != nil {
		return nil, fmt.Errorf("invalid api key: %w", err)
	}

	err = validateSteamID64(steamID64)
	if err != nil {
		return nil, fmt.Errorf("invalid steamID64: %w", err)
	}

	var u *url.URL
	u, err = url.Parse(userSummaryURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user summary URL: %w", err)
	}

	q := u.Query()
	q.Set("key", apiKey)
	q.Set("steamids", strconv.FormatUint(steamID64, 10))
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

	res := &userSummaryResponse{}
	err = json.NewDecoder(resp.Body).Decode(res)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(res.Response.Players) == 0 {
		return nil, fmt.Errorf("no player data found for steamID64: %d", steamID64)
	}

	player := res.Response.Players[0]

	var id uint64
	id, err = strconv.ParseUint(player.Steamid, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse steamID: %w", err)
	}

	return &UserSummary{
		SteamID64:                id,
		CommunityVisibilityState: CommunityVisibilityState(player.Communityvisibilitystate),
		ProfileState:             ProfileState(player.Profilestate),
		PersonaName:              player.Personaname,
		CommentPermission:        CommentPermission(player.Commentpermission),
		ProfileURL:               player.Profileurl,
		Avatar:                   player.Avatar,
		AvatarMedium:             player.Avatarmedium,
		AvatarFull:               player.Avatarfull,
		AvatarHash:               player.Avatarhash,
		LastLogoff:               time.Unix(int64(player.Lastlogoff), 0),
		PersonaState:             PersonaState(player.Personastate),
		PrimaryClanID:            player.Primaryclanid,
		TimeCreated:              time.Unix(int64(player.Timecreated), 0),
		PersonaStateFlags:        player.Personastateflags,
		LocCountryCode:           player.Loccountrycode,
		LocStateCode:             player.Locstatecode,
	}, nil
}

const userSummaryURL = "https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v2/"

type userSummaryResponse struct {
	Response struct {
		Players []struct {
			Steamid                  string `json:"steamid"`
			Communityvisibilitystate int    `json:"communityvisibilitystate"`
			Profilestate             int    `json:"profilestate"`
			Personaname              string `json:"personaname"`
			Commentpermission        int    `json:"commentpermission"`
			Profileurl               string `json:"profileurl"`
			Avatar                   string `json:"avatar"`
			Avatarmedium             string `json:"avatarmedium"`
			Avatarfull               string `json:"avatarfull"`
			Avatarhash               string `json:"avatarhash"`
			Lastlogoff               int    `json:"lastlogoff"`
			Personastate             int    `json:"personastate"`
			Primaryclanid            string `json:"primaryclanid"`
			Timecreated              int    `json:"timecreated"`
			Personastateflags        int    `json:"personastateflags"`
			Loccountrycode           string `json:"loccountrycode"`
			Locstatecode             string `json:"locstatecode"`
		} `json:"players"`
	} `json:"response"`
}

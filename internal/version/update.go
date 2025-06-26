package version

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Masterminds/semver/v3"
)

func CheckForUpdate() (bool, error) {
	latestVersion, err := getLatestGithubReleaseVersion()
	if err != nil {
		if errors.Is(err, errNoLatestRelease) {
			return false, nil
		}

		return false, err
	}

	current, err := semver.NewVersion(Version)
	if err != nil {
		return false, err
	}

	latest, err := semver.NewVersion(latestVersion)
	if err != nil {
		return false, err
	}

	return latest.GreaterThan(current), nil
}

const githubRepoReleases = "https://api.github.com/repos/devusSs/dropawp/releases"

type githubReleases []struct {
	URL       string `json:"url"`
	AssetsURL string `json:"assets_url"`
	UploadURL string `json:"upload_url"`
	HTMLURL   string `json:"html_url"`
	ID        int    `json:"id"`
	Author    struct {
		Login             string `json:"login"`
		ID                int    `json:"id"`
		NodeID            string `json:"node_id"`
		AvatarURL         string `json:"avatar_url"`
		GravatarID        string `json:"gravatar_id"`
		URL               string `json:"url"`
		HTMLURL           string `json:"html_url"`
		FollowersURL      string `json:"followers_url"`
		FollowingURL      string `json:"following_url"`
		GistsURL          string `json:"gists_url"`
		StarredURL        string `json:"starred_url"`
		SubscriptionsURL  string `json:"subscriptions_url"`
		OrganizationsURL  string `json:"organizations_url"`
		ReposURL          string `json:"repos_url"`
		EventsURL         string `json:"events_url"`
		ReceivedEventsURL string `json:"received_events_url"`
		Type              string `json:"type"`
		UserViewType      string `json:"user_view_type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	NodeID          string    `json:"node_id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Draft           bool      `json:"draft"`
	Prerelease      bool      `json:"prerelease"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"`
	Assets          []any     `json:"assets"`
	TarballURL      string    `json:"tarball_url"`
	ZipballURL      string    `json:"zipball_url"`
	Body            string    `json:"body"`
	MentionsCount   int       `json:"mentions_count"`
	Reactions       struct {
		URL        string `json:"url"`
		TotalCount int    `json:"total_count"`
		Num1       int    `json:"+1"`
		Num10      int    `json:"-1"`
		Laugh      int    `json:"laugh"`
		Hooray     int    `json:"hooray"`
		Confused   int    `json:"confused"`
		Heart      int    `json:"heart"`
		Rocket     int    `json:"rocket"`
		Eyes       int    `json:"eyes"`
	} `json:"reactions,omitempty"`
}

var errNoLatestRelease = errors.New("no latest release found")

func getLatestGithubReleaseVersion() (string, error) {
	resp, err := http.Get(githubRepoReleases)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	model := githubReleases{}
	err = json.NewDecoder(resp.Body).Decode(&model)
	if err != nil {
		return "", err
	}

	if len(model) == 0 {
		return "", errNoLatestRelease
	}

	return strconv.Itoa(model[0].ID), nil
}

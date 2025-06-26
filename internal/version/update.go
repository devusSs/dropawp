package version

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/minio/selfupdate"
)

func CheckForUpdate(ctx context.Context) (bool, error) {
	latestVersion, err := getLatestGithubReleaseVersion(ctx)
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

func GetAndApplyUpdate(ctx context.Context) error {
	latestVersion, err := getLatestGithubReleaseVersion(ctx)
	if err != nil {
		return err
	}

	release, err := getRelease(ctx, latestVersion)
	if err != nil {
		return err
	}

	return downloadAndApplyUpdate(ctx, release)
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

func getLatestGithubReleaseVersion(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, githubRepoReleases, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to fetch releases: " + resp.Status)
	}

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

type githubRelease struct {
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
	Assets          []struct {
		URL      string `json:"url"`
		ID       int    `json:"id"`
		NodeID   string `json:"node_id"`
		Name     string `json:"name"`
		Label    string `json:"label"`
		Uploader struct {
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
		} `json:"uploader"`
		ContentType        string    `json:"content_type"`
		State              string    `json:"state"`
		Size               int       `json:"size"`
		Digest             string    `json:"digest"`
		DownloadCount      int       `json:"download_count"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		BrowserDownloadURL string    `json:"browser_download_url"`
	} `json:"assets"`
	TarballURL string `json:"tarball_url"`
	ZipballURL string `json:"zipball_url"`
	Body       string `json:"body"`
}

func getRelease(ctx context.Context, id string) (*githubRelease, error) {
	u, err := url.Parse(githubRepoReleases)
	if err != nil {
		return nil, err
	}

	u.Path, err = url.JoinPath(u.Path, id)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch release: " + resp.Status)
	}

	release := &githubRelease{}
	err = json.NewDecoder(resp.Body).Decode(release)
	if err != nil {
		return nil, err
	}

	return release, nil
}

func downloadAndApplyUpdate(ctx context.Context, release *githubRelease) error {
	assetURL, assetName, err := findSuitableAsset(release)
	if err != nil {
		return err
	}

	tmpFilePath, err := downloadAssetToTempFile(ctx, assetURL, assetName)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFilePath)

	switch filepath.Ext(assetName) {
	case ".gz":
		return extractTarGzAndApply(tmpFilePath)
	case ".zip":
		return extractZipAndApply(tmpFilePath)
	default:
		return errors.New("unsupported file extension: " + filepath.Ext(assetName))
	}
}

func findSuitableAsset(release *githubRelease) (string, string, error) {
	wantOS := strings.ToLower(runtime.GOOS)
	wantArch := strings.ToLower(runtime.GOARCH)

	for _, asset := range release.Assets {
		name := strings.ToLower(asset.Name)
		if strings.Contains(name, wantOS) && strings.Contains(name, wantArch) {
			return asset.BrowserDownloadURL, asset.Name, nil
		}
	}

	return "", "", errors.New("no suitable asset found for " + wantOS + "/" + wantArch)
}

func downloadAssetToTempFile(ctx context.Context, url, name string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to download update: " + resp.Status)
	}

	tmpFile, err := os.CreateTemp("", name)
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func extractTarGzAndApply(tarGzPath string) error {
	f, err := os.Open(tarGzPath)
	if err != nil {
		return err
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		var header *tar.Header
		header, err = tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}

		if header.Typeflag == tar.TypeReg {
			return selfupdate.Apply(tr, selfupdate.Options{})
		}
	}

	return errors.New("no regular file found in tar.gz")
}

func extractZipAndApply(zipPath string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, file := range r.File {
		if file.FileInfo().IsDir() || file.UncompressedSize64 == 0 {
			continue
		}

		var rc io.ReadCloser
		rc, err = file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		return selfupdate.Apply(rc, selfupdate.Options{})
	}

	return errors.New("no suitable file found in zip")
}

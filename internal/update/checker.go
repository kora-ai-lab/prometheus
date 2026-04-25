package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/blang/semver"
)

var localVersion = "1.0.0"

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Body   string `json:"body"`
	Assets []struct {
		Name        string `json:"name"`
		BrowserURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func CheckForUpdate(owner, repo string) (hasUpdate bool, latestVersion string, err error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, "", fmt.Errorf("cannot create request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, "", fmt.Errorf("cannot check for updates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return false, "", fmt.Errorf("parse release: %w", err)
	}

	latestVersion = strings.TrimPrefix(release.TagName, "v")

	cmp, err := CompareVersions(localVersion, latestVersion)
	if err != nil {
		return false, "", err
	}

	return cmp < 0, latestVersion, nil
}

func CompareVersions(local, remote string) (int, error) {
	localV, err := semver.Make(local)
	if err != nil {
		return 0, fmt.Errorf("invalid local version %q: %w", local, err)
	}
	remoteV, err := semver.Make(remote)
	if err != nil {
		return 0, fmt.Errorf("invalid remote version %q: %w", remote, err)
	}
	return localV.Compare(remoteV), nil
}
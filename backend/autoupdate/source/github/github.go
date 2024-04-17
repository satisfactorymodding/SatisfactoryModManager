package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Masterminds/semver/v3"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/autoupdate/updater"
)

type source struct {
	repo string
}

func MakeGithubSource(repo string) updater.Source {
	return &source{
		repo: repo,
	}
}

func (g *source) GetLatestVersion(includePrerelease bool) (string, error) {
	if !includePrerelease {
		release, err := g.getLatestReleaseData()
		if err != nil {
			return "", err
		}
		return release.TagName, nil
	}

	// GitHub does not return pre-releases on the /latest endpoint
	allReleases, err := g.getReleasesData()
	var latest *semver.Version
	var latestTagName string
	if err != nil {
		return "", err
	}
	for _, release := range allReleases {
		version, err := semver.NewVersion(release.TagName)
		if err != nil {
			continue
		}
		if latest == nil || version.GreaterThan(latest) {
			latest = version
			latestTagName = release.TagName
		}
	}
	if latest == nil {
		return "", fmt.Errorf("no releases found")
	}
	return latestTagName, nil
}

func (g *source) GetFile(version string, filename string) (io.ReadCloser, int64, error) {
	release, err := g.getReleaseData(version)
	if err != nil {
		return nil, 0, err
	}
	fileURL := getAssetURL(release, filename)
	if fileURL == "" {
		return nil, 0, fmt.Errorf("failed to find asset")
	}
	response, err := http.Get(fileURL)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to download asset: %w", err)
	}
	return response.Body, response.ContentLength, nil
}

func getAssetURL(release *Release, assetName string) string {
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			return asset.BrowserDownloadURL
		}
	}
	return ""
}

func (g *source) GetChangelogs() (map[string]string, error) {
	releases, err := g.getReleasesData()
	if err != nil {
		return nil, err
	}
	changelogs := make(map[string]string)
	for _, release := range releases {
		changelogs[release.TagName] = release.Body
	}
	return changelogs, nil
}

func (g *source) getLatestReleaseData() (*Release, error) {
	response, err := http.Get("https://api.github.com/repos/" + g.repo + "/releases/latest")
	if err != nil {
		return nil, fmt.Errorf("failed to get latest release: %w", err)
	}
	defer response.Body.Close()
	var release Release
	err = json.NewDecoder(response.Body).Decode(&release)
	if err != nil {
		return nil, fmt.Errorf("failed to decode latest release: %w", err)
	}
	return &release, nil
}

func (g *source) getReleasesData() ([]Release, error) {
	response, err := http.Get("https://api.github.com/repos/" + g.repo + "/releases")
	if err != nil {
		return nil, fmt.Errorf("failed to get releases: %w", err)
	}
	defer response.Body.Close()
	var releases []Release
	err = json.NewDecoder(response.Body).Decode(&releases)
	if err != nil {
		return nil, fmt.Errorf("failed to decode releases: %w", err)
	}
	return releases, nil
}

func (g *source) getReleaseData(tagName string) (*Release, error) {
	response, err := http.Get("https://api.github.com/repos/" + g.repo + "/releases/tags/" + tagName)
	if err != nil {
		return nil, fmt.Errorf("failed to get releases: %w", err)
	}
	defer response.Body.Close()
	var release Release
	err = json.NewDecoder(response.Body).Decode(&release)
	if err != nil {
		return nil, fmt.Errorf("failed to decode releases: %w", err)
	}
	return &release, nil
}

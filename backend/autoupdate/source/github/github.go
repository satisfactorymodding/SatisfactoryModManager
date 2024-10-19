package github

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"

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
	// This is called very often, so use the atom feed to avoid being rate limited
	allReleases, err := g.getReleasesAtom()
	var latest *semver.Version
	var latestTagName string
	if err != nil {
		return "", err
	}
	for _, releaseVersion := range allReleases {
		version, err := semver.NewVersion(releaseVersion)
		if err != nil {
			continue
		}
		if !includePrerelease && version.Prerelease() != "" {
			continue
		}
		if latest == nil || version.GreaterThan(latest) {
			latest = version
			latestTagName = releaseVersion
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

var releaseURLTagRegex = regexp.MustCompile(`/tag/([^/]+)$`)

func (g *source) getReleasesAtom() ([]string, error) {
	response, err := http.Get(fmt.Sprintf("https://github.com/%s/releases.atom", g.repo))
	if err != nil {
		return nil, fmt.Errorf("failed to get releases: %w", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	var feed Feed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal feed: %w", err)
	}

	releases := make([]string, 0, len(feed.Entry))

	for _, entry := range feed.Entry {
		tagMatches := releaseURLTagRegex.FindStringSubmatch(entry.Link.Href)
		if len(tagMatches) != 2 {
			slog.Warn("failed to parse tag from release URL", slog.String("url", entry.Link.Href))
			continue
		}
		releases = append(releases, tagMatches[1])
		slog.Debug("found release", slog.String("tag", tagMatches[1]))
	}

	return releases, nil
}

func (g *source) getReleasesData() ([]Release, error) {
	response, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/releases", g.repo))
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
	response, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/releases/tags/%s", g.repo, tagName))
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

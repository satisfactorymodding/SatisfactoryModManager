package github

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Masterminds/semver/v3"
)

type Provider struct {
	repo             string
	checksumArtifact string
}

func MakeGithubProvider(repo string, checksumArtifact string) *Provider {
	return &Provider{
		repo:             repo,
		checksumArtifact: checksumArtifact,
	}
}

func (g *Provider) GetLatestVersion(includePrerelease bool) (string, error) {
	if !includePrerelease {
		release, err := g.getLatestReleaseData()
		if err != nil {
			return "", fmt.Errorf("failed to get latest release: %w", err)
		}
		return release.TagName, nil
	}

	// GitHub does not return pre-releases on the /latest endpoint
	allReleases, err := g.getReleasesData()
	var latest *semver.Version
	var latestTagName string
	if err != nil {
		return "", fmt.Errorf("failed to get releases: %w", err)
	}
	for _, release := range allReleases {
		version, err := semver.NewVersion(release.TagName)
		if err != nil {
			continue
		}
		if !includePrerelease && version.Prerelease() != "" {
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

func (g *Provider) GetFile(version string, filename string) (io.ReadCloser, int64, []byte, error) {
	release, err := g.getReleaseData(version)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("failed to get latest release: %w", err)
	}
	fileURL := getAssetURL(release, filename)
	if fileURL == "" {
		return nil, 0, nil, fmt.Errorf("failed to find asset")
	}
	checksum, err := g.getFileChecksum(release, filename)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("failed to get checksum: %w", err)
	}
	response, err := http.Get(fileURL)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("failed to download asset: %w", err)
	}
	return response.Body, response.ContentLength, checksum, nil
}

func getAssetURL(release *Release, assetName string) string {
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			return asset.BrowserDownloadURL
		}
	}
	return ""
}

func (g *Provider) GetChangelogs() (map[string]string, error) {
	releases, err := g.getReleasesData()
	if err != nil {
		return nil, fmt.Errorf("failed to get latest release: %w", err)
	}
	changelogs := make(map[string]string)
	for _, release := range releases {
		changelogs[release.TagName] = release.Body
	}
	return changelogs, nil
}

func (g *Provider) getLatestReleaseData() (*Release, error) {
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

func (g *Provider) getReleasesData() ([]Release, error) {
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

func (g *Provider) getReleaseData(tagName string) (*Release, error) {
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

func (g *Provider) getFileChecksum(release *Release, filename string) ([]byte, error) {
	if g.checksumArtifact == "" {
		return nil, nil
	}
	url := getAssetURL(release, g.checksumArtifact)
	if url == "" {
		return nil, fmt.Errorf("failed to find checksum asset")
	}
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download checksum asset: %w", err)
	}
	defer response.Body.Close()
	checksum, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read checksum: %w", err)
	}
	checksums := parseChecksumFile(checksum)
	if sum, ok := checksums[filename]; ok {
		return sum, nil
	}
	return nil, fmt.Errorf("failed to find checksum for file")
}

func parseChecksumFile(checksumFile []byte) map[string][]byte {
	checksums := make(map[string][]byte)
	lines := strings.Split(string(checksumFile), "\n")
	for _, line := range lines {
		if len(line) == 0 {
			// Skip empty lines
			continue
		}
		parts := strings.Split(line, "  ")
		if len(parts) != 2 {
			slog.Debug("invalid checksum entry", slog.String("entry", line))
			continue
		}
		hexSum := parts[0]
		filename := parts[1]
		sum, err := hex.DecodeString(hexSum)
		if err != nil {
			slog.Debug("failed to decode checksum", slog.String("checksum", hexSum), slog.String("filename", filename), slog.Any("error", err))
			continue
		}
		checksums[parts[1]] = sum
	}
	return checksums
}

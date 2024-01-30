package github

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type Provider struct {
	repo string
}

func MakeGithubProvider(repo string) *Provider {
	return &Provider{
		repo: repo,
	}
}

func (g *Provider) GetLatestVersion() (string, error) {
	release, err := g.getLatestReleaseData()
	if err != nil {
		return "", errors.Wrap(err, "failed to get latest release")
	}
	return release.TagName, nil
}

func (g *Provider) GetFile(filename string) (io.ReadCloser, int64, error) {
	release, err := g.getLatestReleaseData()
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to get latest release")
	}
	var url string
	for _, asset := range release.Assets {
		if asset.Name == filename {
			url = asset.BrowserDownloadURL
			break
		}
	}
	if url == "" {
		return nil, 0, errors.Errorf("failed to find asset")
	}
	response, err := http.Get(url)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "failed to download asset")
	}
	return response.Body, response.ContentLength, nil
}

func (g *Provider) GetChangelogs() (map[string]string, error) {
	releases, err := g.getReleasesData()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get latest release")
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
		return nil, errors.Wrap(err, "failed to get latest release")
	}
	defer response.Body.Close()
	var release Release
	err = json.NewDecoder(response.Body).Decode(&release)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode latest release")
	}
	return &release, nil
}

func (g *Provider) getReleasesData() ([]Release, error) {
	response, err := http.Get("https://api.github.com/repos/" + g.repo + "/releases")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get releases")
	}
	defer response.Body.Close()
	var releases []Release
	err = json.NewDecoder(response.Body).Decode(&releases)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode releases")
	}
	return releases, nil
}

package app

import (
	"strings"

	"github.com/spf13/viper"
)

func (a *app) GetVersion() string {
	return viper.GetString("version")
}

func (a *app) GetCommit() string {
	return viper.GetString("commit")
}

func (a *app) GetDate() string {
	return viper.GetString("date")
}

func (a *app) GetAPIEndpoint() string {
	return viper.GetString("api-base") + viper.GetString("graphql-api")
}

func (a *app) GetSiteEndpoint() string {
	return strings.Replace(viper.GetString("api-base"), "api.", "", 1)
}

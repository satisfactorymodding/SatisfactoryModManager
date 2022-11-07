package ficsitcli_bindings

import (
	"github.com/satisfactorymodding/SatisfactoryModManager/install_finders"
	"github.com/satisfactorymodding/ficsit-cli/cli"
)

type InstallationInfo struct {
	Installation *cli.Installation             `json:"installation"`
	Info         *install_finders.Installation `json:"info"`
}

type Progress struct {
	Item     string  `json:"item"`
	Message  string  `json:"message"`
	Progress float64 `json:"progress"`
}

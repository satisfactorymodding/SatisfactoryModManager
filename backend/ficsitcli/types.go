package ficsitcli

import "github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"

type InstallState string

const (
	InstallStateUnknown InstallState = "unknown"
	InstallStateLoading InstallState = "loading"
	InstallStateInvalid InstallState = "invalid"
	InstallStateValid   InstallState = "valid"
)

type installationMetadata struct {
	State InstallState         `json:"state"`
	Info  *common.Installation `json:"info"`
}

type Progress struct {
	Item     string  `json:"item"`
	Message  string  `json:"message"`
	Progress float64 `json:"progress"`
}

var AllInstallationStates = []struct {
	Value  InstallState
	TSName string
}{
	{InstallStateUnknown, "UNKNOWN"},
	{InstallStateLoading, "LOADING"},
	{InstallStateInvalid, "INVALID"},
	{InstallStateValid, "VALID"},
}

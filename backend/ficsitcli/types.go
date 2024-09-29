package ficsitcli

import (
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/utils"
)

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

type Action string

const (
	ActionInstall       Action = "install"
	ActionUninstall     Action = "uninstall"
	ActionEnable        Action = "enable"
	ActionDisable       Action = "disable"
	ActionSelectInstall Action = "selectInstall"
	ActionToggleMods    Action = "toggleMods"
	ActionSelectProfile Action = "selectProfile"
	ActionImportProfile Action = "importProfile"
	ActionUpdate        Action = "update"
	ActionApply         Action = "apply"
)

type Progress struct {
	Action Action                    `json:"action"`
	Item   ProgressItem              `json:"item"`
	Tasks  map[string]utils.Progress `json:"tasks"`
}

type ProgressItem struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

var noItem = ProgressItem{}

func newSimpleItem(name string) ProgressItem {
	return ProgressItem{
		Name: name,
	}
}

func newItem(name, version string) ProgressItem {
	return ProgressItem{
		Name:    name,
		Version: version,
	}
}

func newProgress(action Action, item ProgressItem) *Progress {
	return &Progress{
		Action: action,
		Item:   item,
		Tasks:  make(map[string]utils.Progress),
	}
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

var AllActionTypes = []struct {
	Value  Action
	TSName string
}{
	{ActionInstall, "INSTALL"},
	{ActionUninstall, "UNINSTALL"},
	{ActionEnable, "ENABLE"},
	{ActionDisable, "DISABLE"},
	{ActionSelectInstall, "SELECT_INSTALL"},
	{ActionToggleMods, "TOGGLE_MODS"},
	{ActionSelectProfile, "SELECT_PROFILE"},
	{ActionImportProfile, "IMPORT_PROFILE"},
	{ActionUpdate, "UPDATE"},
	{ActionApply, "APPLY"},
}

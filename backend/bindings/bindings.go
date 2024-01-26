package bindings

import (
	"context"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/bindings/ficsitcli"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
)

type Bindings struct {
	App       *App
	Update    *Update
	FicsitCLI *ficsitcli.FicsitCLI
	settings  *Settings
	debugInfo *DebugInfo
}

var AllInstallTypes = []struct {
	Value  common.InstallType
	TSName string
}{
	{common.InstallTypeWindowsClient, "WINDOWS"},
	{common.InstallTypeWindowsServer, "WINDOWS_SERVER"},
	{common.InstallTypeLinuxServer, "LINUX_SERVER"},
}

var AllBranches = []struct {
	Value  common.GameBranch
	TSName string
}{
	{common.BranchEarlyAccess, "EARLY_ACCESS"},
	{common.BranchExperimental, "EXPERIMENTAL"},
}

var BindingsInstance *Bindings

func MakeBindings() (*Bindings, error) {
	if BindingsInstance != nil {
		return BindingsInstance, nil
	}

	ficsitCLI, err := ficsitcli.MakeFicsitCLI()
	if err != nil {
		return nil, err
	}

	BindingsInstance = &Bindings{
		App:       MakeApp(),
		Update:    MakeUpdate(),
		FicsitCLI: ficsitCLI,
		settings:  MakeSettings(),
		debugInfo: MakeDebugInfo(),
	}

	return BindingsInstance, nil
}

func (b *Bindings) Startup(ctx context.Context) {
	b.App.startup(ctx)
	b.Update.startup(ctx)
	b.FicsitCLI.Startup(ctx)
	b.settings.startup(ctx)
	b.debugInfo.startup(ctx)
}

func (b *Bindings) GetBindings() []interface{} {
	return []interface{}{
		b.App,
		b.Update,
		b.FicsitCLI,
		b.settings,
		b.debugInfo,
	}
}

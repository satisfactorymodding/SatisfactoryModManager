package bindings

import (
	"context"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/bindings/ficsitcli"
	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
)

type Bindings struct {
	FicsitCLI *ficsitcli.FicsitCLI
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

func MakeBindings() *Bindings {
	if BindingsInstance != nil {
		return BindingsInstance
	}

	BindingsInstance = &Bindings{
		FicsitCLI: ficsitcli.MakeFicsitCLI(),
		debugInfo: MakeDebugInfo(),
	}

	return BindingsInstance
}

func (b *Bindings) Startup(ctx context.Context) error {
	b.debugInfo.startup(ctx)
	return b.FicsitCLI.Startup(ctx)
}

func (b *Bindings) GetBindings() []interface{} {
	return []interface{}{
		b.FicsitCLI,
		b.debugInfo,
	}
}

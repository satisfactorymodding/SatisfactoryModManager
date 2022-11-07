package bindings

import (
	"context"

	"github.com/satisfactorymodding/SatisfactoryModManager/bindings/ficsitcli_bindings"
)

type Bindings struct {
	App       *App
	FicsitCLI *ficsitcli_bindings.FicsitCLI
	settings  *Settings
	debugInfo *DebugInfo
}

var BindingsInstance *Bindings

func MakeBindings() (*Bindings, error) {
	if BindingsInstance != nil {
		return BindingsInstance, nil
	}

	app := MakeApp()
	ficsitCLI, err := ficsitcli_bindings.MakeFicsitCLI()
	if err != nil {
		return nil, err
	}
	settings := MakeSettings()
	debugInfo := MakeDebugInfo()

	BindingsInstance = &Bindings{
		App:       app,
		FicsitCLI: ficsitCLI,
		settings:  settings,
		debugInfo: debugInfo,
	}

	return BindingsInstance, nil
}

func (b *Bindings) Startup(ctx context.Context) {
	b.App.startup(ctx)
	b.FicsitCLI.Startup(ctx)
	b.settings.startup(ctx)
	b.debugInfo.startup(ctx)
}

func (b *Bindings) GetBindings() []interface{} {
	return []interface{}{
		b.App,
		b.FicsitCLI,
		b.settings,
		b.debugInfo,
	}
}

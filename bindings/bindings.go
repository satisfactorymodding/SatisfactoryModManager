package bindings

import (
	"context"

	"github.com/pkg/errors"
)

type Bindings struct {
	App         *App
	FicsitCLI   *FicsitCLI
	Settings    *Settings
	ProjectFile *ProjectFile
}

var BindingsInstance *Bindings

func MakeBindings() (*Bindings, error) {
	if BindingsInstance != nil {
		return BindingsInstance, nil
	}

	app := MakeApp()
	ficsitCLI, err := MakeFicsitCLI()
	if err != nil {
		return nil, errors.Wrap(err, "failed to make ficsitCLI bindings")
	}
	settings, err := MakeSettings()
	if err != nil {
		return nil, errors.Wrap(err, "failed to make settings bindings")
	}
	projectFile, err := MakeProjectFile()
	if err != nil {
		return nil, errors.Wrap(err, "failed to make project file bindings")
	}

	BindingsInstance = &Bindings{
		App:         app,
		FicsitCLI:   ficsitCLI,
		Settings:    settings,
		ProjectFile: projectFile,
	}

	return BindingsInstance, nil
}

func (b *Bindings) Startup(ctx context.Context) {
	b.App.startup(ctx)
	b.FicsitCLI.startup(ctx)
	b.Settings.startup(ctx)
	b.ProjectFile.startup(ctx)
}

func (b *Bindings) GetBindings() []interface{} {
	return []interface{}{
		b.App,
		b.FicsitCLI,
		b.Settings,
		b.ProjectFile,
	}
}

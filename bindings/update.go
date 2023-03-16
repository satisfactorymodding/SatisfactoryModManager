package bindings

import (
	"context"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/satisfactorymodding/SatisfactoryModManager/autoupdate"
)

type Update struct {
	ctx     context.Context
	Restart bool
}

func MakeUpdate() *Update {
	return &Update{}
}

func (u *Update) startup(ctx context.Context) {
	u.ctx = ctx
}

func (u *Update) CheckForUpdates() {
	_ = autoupdate.Updater.CheckForUpdate()
}

func (u *Update) UpdateAndRestart() {
	u.Restart = true
	wailsRuntime.Quit(u.ctx)
}

func (u *Update) UpdateAvailable(latestVersion string, changelogs map[string]string) {
	wailsRuntime.EventsEmit(u.ctx, "updateAvailable", latestVersion, changelogs)
}

func (u *Update) UpdateDownloadProgress(downloaded, total int64) {
	wailsRuntime.EventsEmit(u.ctx, "updateDownloadProgress", downloaded, total)
}

func (u *Update) UpdateReady() {
	wailsRuntime.EventsEmit(u.ctx, "updateReady")
}

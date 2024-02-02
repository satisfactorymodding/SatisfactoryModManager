package backend

import (
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/app"
)

func ProcessArguments(args []string) {
	if len(args) < 1 {
		return
	}
	if strings.HasPrefix(args[0], "smmanager://") {
		uri := args[0]
		err := handleURI(uri)
		if err != nil {
			slog.Error("failed to handle smmanager:// URI", slog.Any("error", err), slog.String("uri", uri))
		}
	} else {
		err := handleFile(args[0])
		if err != nil {
			slog.Error("failed to handle file", slog.Any("error", err), slog.String("path", args[0]))
		}
	}
	app.App.Show()
}

func handleURI(uri string) error {
	u, err := url.Parse(uri)
	if err != nil {
		return fmt.Errorf("failed to parse URI: %w", err)
	}
	switch u.Host {
	case "install":
		modID := u.Query().Get("modID")
		version := u.Query().Get("version")
		app.App.ExternalInstallMod(modID, version)
		return nil
	default:
		return fmt.Errorf("unknown URI action " + u.Host)
	}
}

func handleFile(path string) error {
	if strings.HasSuffix(path, ".smmprofile") {
		println(path)
		app.App.ExternalImportProfile(path)
		return nil
	}
	return fmt.Errorf("unknown file type " + path)
}

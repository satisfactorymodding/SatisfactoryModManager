package main

import (
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/satisfactorymodding/SatisfactoryModManager/bindings"
)

func processArguments(args []string) {
	if len(args) < 2 {
		return
	}
	if strings.HasPrefix(args[1], "smmanager://") {
		uri := args[1]
		err := handleURI(uri)
		if err != nil {
			log.Error().Err(err).Str("uri", uri).Msg("Failed to handle smmanager:// URI")
		}
	} else {
		handleFile(args[1])
	}
	// TODO make window flash/bring to front
}

func handleURI(uri string) error {
	u, err := url.Parse(uri)
	if err != nil {
		return errors.Wrap(err, "failed to parse URI")
	}
	switch u.Host {
	case "install":
		modId := u.Query().Get("modID")
		version := u.Query().Get("version")
		bindings.BindingsInstance.App.ExternalInstallMod(modId, version)
		return nil
	}
	return errors.New("unknown URI action " + u.Host)
}

func handleFile(path string) error {
	if strings.HasSuffix(path, ".smmprofile") {
		println(path)
		bindings.BindingsInstance.App.ExternalImportProfile(path)
		return nil
	}
	return errors.New("unknown file type " + path)
}

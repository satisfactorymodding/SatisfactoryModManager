package file_scheme_association

import (
	"github.com/pkg/errors"
	"golang.org/x/sys/windows/registry"
)

func SetAsDefaultSchemeHandler(uriScheme string) error {
	uriKeyPath := "SOFTWARE\\Classes\\" + uriScheme
	permission := uint32(registry.QUERY_VALUE | registry.SET_VALUE)

	realExecutablePath, err := getRealExecutablePath()
	if err != nil {
		return errors.Wrap(err, "failed to get real executable path")
	}

	// create key
	k, _, err := registry.CreateKey(registry.CURRENT_USER, uriKeyPath, permission)
	if err != nil {
		return errors.Wrap(err, "failed to create base key")
	}

	// set description
	k.SetStringValue("", "URL:"+uriScheme)
	k.SetStringValue("URL Protocol", "")

	// set open command
	err = writeOpenCommand(k, realExecutablePath, permission)
	if err != nil {
		return errors.Wrap(err, "failed to write open command")
	}

	return nil
}

func SetAsDefaultFileHandler(extension string) error {
	uriKeyPath := "SOFTWARE\\Classes\\" + extension
	permission := uint32(registry.QUERY_VALUE | registry.SET_VALUE)

	realExecutablePath, err := getRealExecutablePath()
	if err != nil {
		return errors.Wrap(err, "failed to get real executable path")
	}

	// create key
	k, _, err := registry.CreateKey(registry.CURRENT_USER, uriKeyPath, permission)
	if err != nil {
		return errors.Wrap(err, "failed to create base key")
	}

	// set open command
	err = writeOpenCommand(k, realExecutablePath, permission)
	if err != nil {
		return errors.Wrap(err, "failed to write open command")
	}

	return nil
}

func writeOpenCommand(key registry.Key, executable string, permission uint32) error {
	// create tree
	shellKey, _, err := registry.CreateKey(key, "shell", permission)
	if err != nil {
		return errors.Wrap(err, "failed to create shell key")
	}
	openKey, _, err := registry.CreateKey(shellKey, "open", permission)
	if err != nil {
		return errors.Wrap(err, "failed to create open key")
	}
	commandKey, _, err := registry.CreateKey(openKey, "command", permission)
	if err != nil {
		return errors.Wrap(err, "failed to create command key")
	}

	// set open command
	err = commandKey.SetStringValue("", executable+" \"%1\"")
	if err != nil {
		return errors.Wrap(err, "failed to set command value")
	}

	return nil
}

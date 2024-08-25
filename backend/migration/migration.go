package migration

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type migration struct {
	smm2Dir                    string
	migrationSuccessMarkerPath string
}

var Migration *migration

func Init() {
	if Migration == nil {
		Migration = &migration{}
		Migration.smm2Dir = filepath.Join(viper.GetString("smm-local-dir"), "profiles")
		Migration.migrationSuccessMarkerPath = filepath.Join(Migration.smm2Dir, migrationSuccessMarkerFile)
	}
}

const migrationSuccessMarkerFile = ".smm3_migration_acknowledged"

// https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go
func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	slog.Warn("error when checking path exists, so assuming it does not exist", slog.String("path", path), slog.Any("error", err))
	return false
}

func (m *migration) NeedsSmm2Migration() bool {
	if pathExists(m.smm2Dir) {
		return !pathExists(Migration.migrationSuccessMarkerPath)
	}
	return false
}

func (m *migration) MarkSmm2MigrationSuccess() error {
	file, err := os.Create(Migration.migrationSuccessMarkerPath)
	if err != nil {
		return fmt.Errorf("failed to create migration success marker file: %w", err)
	}
	err = file.Close()
	if err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}
	return nil
}

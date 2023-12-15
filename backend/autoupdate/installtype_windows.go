package autoupdate

import (
	"path/filepath"

	"github.com/spf13/viper"
	"golang.org/x/sys/windows/registry"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/autoupdate/updater/apply"
)

func getInstallType() (string, apply.Apply) {
	nsis, elevation := isNsis()
	if nsis {
		return "Satisfactory-Mod-Manager-Setup.exe", apply.MakeNsisApply(apply.NsisApplyConfig{
			InstallerDownloadPath: filepath.Join(viper.GetString("cache-dir"), "Satisfactory-Mod-Manager-Setup.exe"),
			Elevation:             elevation,
		})
	}
	return "Satisfactory Mod Manager.exe", apply.MakeSingleFileApply()
}

func isNsis() (bool, bool) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\Satisfactory Mod Manager`, registry.READ)
	if err == nil {
		k.Close()
		return true, true
	}
	k, err = registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\Satisfactory Mod Manager`, registry.READ)
	if err == nil {
		k.Close()
		return true, false
	}
	return false, false
}

package autoupdate

import (
	"path/filepath"

	"github.com/spf13/viper"
	"golang.org/x/sys/windows/registry"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/autoupdate/apply"
)

func init() {
	// Currently there is no standalone release for Windows
	// registerUpdateType("standalone", func() UpdateType {
	//   return UpdateType{
	//	   ArtifactName: "SatisfactoryModManager.exe",
	//	   Apply:        apply.MakeSingleFileApply(),
	//   }
	// })
	registerUpdateType("nsis", func() UpdateType {
		return UpdateType{
			ArtifactName: "Satisfactory-Mod-Manager-Setup.exe",
			Apply: apply.MakeNsisApply(apply.NsisApplyConfig{
				InstallerDownloadPath: filepath.Join(viper.GetString("smm-cache-dir"), "Satisfactory-Mod-Manager-Setup.exe"),
				Elevation:             requireElevation(),
			}),
		}
	})
}

func requireElevation() bool {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\Satisfactory Mod Manager`, registry.READ)
	if err == nil {
		k.Close()
		return true
	}
	return false
}

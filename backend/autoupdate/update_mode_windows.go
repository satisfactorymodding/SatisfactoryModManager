package autoupdate

import (
	"os"
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
			ArtifactName: "SatisfactoryModManager-Setup.exe",
			Apply: apply.MakeNsisApply(apply.NsisApplyConfig{
				InstallerDownloadPath: filepath.Join(viper.GetString("smm-cache-dir"), "SatisfactoryModManager-Setup.exe"),
				IsAllUsers:            isAllUsers(),
			}),
		}
	})
}

func isAllUsers() bool {
	executable, _ := os.Executable()
	allUsersInstallPath := getInstallPath(registry.LOCAL_MACHINE)
	currentUserInstallPath := getInstallPath(registry.CURRENT_USER)
	if allUsersInstallPath != "" && currentUserInstallPath != "" {
		// Installed in both modes, so we need to check if the currently running executable is all-users or per-user
		return allUsersInstallPath == filepath.Dir(executable)
	}
	if allUsersInstallPath == "" && currentUserInstallPath == "" {
		// This should never happen, but since we don't know if the user has admin rights, we will default to per-user
		return false
	}
	return allUsersInstallPath != ""
}

func getInstallPath(key registry.Key) string {
	k, err := registry.OpenKey(key, `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall\Satisfactory Mod Manager`, registry.READ)
	if err != nil {
		return ""
	}
	defer k.Close()
	installPath, _, err := k.GetStringValue("InstallLocation")
	if err != nil {
		return ""
	}
	return installPath
}

package common

type Platform interface {
	ProcessPath(path string) string
	CacheDir() (string, error)
	Os() string
}

type LauncherPlatform struct {
	Platform
	launcherCommand func(arg string) []string
}

func MakeLauncherPlatform(platform Platform, launcherCommand func(arg string) []string) LauncherPlatform {
	return LauncherPlatform{Platform: platform, launcherCommand: launcherCommand}
}

func (p LauncherPlatform) LauncherCommand(arg string) []string {
	if p.launcherCommand != nil {
		return p.launcherCommand(arg)
	}
	return nil
}

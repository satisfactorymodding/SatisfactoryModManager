# NSIS Multi User Plugin
NSIS plugin that allows "per-user" (no admin required) and "per-machine" (asks elevation *only when necessary*) installations. This plugin was inspired by [MultiUser.nsh (by Joost Verburg)](http://nsis.sourceforge.net/Docs/MultiUser/Readme.html), but supports a lot of new features and is easier to use.

## How It Works

### Installer
The plugin creates a custom Install Options page based on the nsDisalogs library that is displayed before the Components page. The page is displayed always and has two options: install for all users (per-machine) and install for current user only (per-user). When the user starts the setup, he is not forced to elevate in the beginning. If the user selects per-user install, he can install only for himself without being asked for elevation (except when there is per-machine installation that needs to be removed first). If the user selects per-machine install, the Windows shield is displayed on the Next button and elevation is required. Limited users can also install per-machine as long as they know the administrator credentials. 

### Uninstaller
The plugin creates the same custom page and shows it in the beginning of the uninstaller if there are two installations. Elevation is required only when per-machine version is uninstalled. If there is only one installed version or if command-line parameters are passed specifying which version to uninstall, the page is not displayed. In this case, the uninstaller asks for elevation if per-machine version is to be uninstalled. When invoked from the Windows Uninstall dialog or from the Start menu, a parameter to the uninstaller is passed, so that it detects which verion to uninstall, and the page is not displayed.

### Both
An option (`MULTIUSER_INSTALLMODE_ALLOW_ELEVATION`) defines whether elevation if allowed. If elevation is disabled, the per-machine option becomes available only if the (un)installer is started elevated from Windows and is disabled otherwise.

An option (`MULTIUSER_INSTALLMODE_ALLOW_BOTH_INSTALLATIONS`) defines whether simultaneous per-user and per-machine installations on the same machine are allowed. If set to disallow, the installer alaways requires elevation when there's per-machine installation in order to remove it first.

## Features
- Not tied or dependant on any particular user interface. Supports Modern UI 1/2, ModernUIEx, Ultra Modern UI,  the native NSIS interface, as well as any other interface that supports nsDialogs pages.
- Includes fully functional demos for all supported interfaces that you can use as skeletons to start your scripts from.
- Support for 64-bit installations
- Correctly creates and removes full registry uninstall information like icon and estimated size (separate per-user and per-machine entries)
- Fully supports silent mode, command-line switches and error level handling
- Fully documented

## Screenshots

When `MULTIUSER_INSTALLMODE_ALLOW_ELEVATION` is `1`, there is no existing istallation and running as a regular user (Ultra Modern UI).
Installation for current user requires no elevation:

![Installation for current user requires no elevation](./Screenshots/01.png?raw=true "Installation for current user requires no elevation")
![Per-user installation folder](./Screenshots/02.png?raw=true "Per-user installation folder")

Installation for all users requires elevation:

![Installation for all users requires elevation](./Screenshots/03.png?raw=true "Installation for all users requires elevation")

When running as admin, no elevation is required:

![When running as admin, no elevation is required](./Screenshots/04.png?raw=true "When running as admin, no elevation is required")

When there is an existing installation, it is always selected (Modern UI 2):

![Existing instalation is always selected](./Screenshots/05.png?raw=true "Existing instalation is always selected")

When `MULTIUSER_INSTALLMODE_ALLOW_BOTH_INSTALLATIONS` is `0`, there is existing per-machine installation and running as a regular user, elevation to install per-user is required (Modern UI 2):

![When there is per-machine installation and elevation per-user is required](./Screenshots/06.png?raw=true "When there is per-machine installation and elevation per-user is required")

When `MULTIUSER_INSTALLMODE_ALLOW_ELEVATION` is `0` and running as a regular user, per-machine option is disabled (native NSIS interface):

![Per-machine option is disabled](./Screenshots/07.png?raw=true "Per-machine option is disabled")

When invoked with the `/allusers` parameter and `MULTIUSER_INSTALLMODE_ALLOW_ELEVATION` is `1` (native NSIS interface):

![/allusers parameter](./Screenshots/08.png?raw=true "/allusers parameter")

When invoked with the `/allusers` parameter and `MULTIUSER_INSTALLMODE_ALLOW_ELEVATION` is `0`:

![/allusers parameter and elevation is disabled](./Screenshots/09.png?raw=true "/allusers parameter and elevation is disabled")

When there are both per-user and per-machine installations and uninstaller is invoked without parameters, page is displayed (Ultra Modern UI):

![Uninstaller page](./Screenshots/10.png?raw=true "Uninstaller page")

The Windows Uninstall list of programs will show individual entries when there are both per-machine and per-user installations (one is stored in `HKLM` and other in `HKCU`):

![The Windows Uninstall list of programs](./Screenshots/11.png?raw=true "The Windows Uninstall list of programs")

The help dialog, invoked with the `/?` parameter:

![The help dialog](./Screenshots/12.png?raw=true "The help dialog")

## Usage

Please look at the fully functional demos in the `Demos` folder.

## Documentation

The full NsisMultiUser documentation is avaialable on the [Wiki](https://github.com/Drizin/NsisMultiUser/wiki).

You can also look at:
- [UAC plugin page](http://nsis.sourceforge.net/UAC_plug-in)
- [The original MultiUser.nsh plugin](http://nsis.sourceforge.net/Docs/MultiUser/Readme.html)
- [nsDialogs plugin](http://nsis.sourceforge.net/Docs/nsDialogs/Readme.html)
- [Modern UI](http://nsis.sourceforge.net/Docs/Modern%20UI/Readme.html)
- [Ultra Modern UI](http://ultramodernui.sourceforge.net/)
- [MSDN documentation](https://msdn.microsoft.com/en-us/library/windows/desktop/dd765197.aspx)


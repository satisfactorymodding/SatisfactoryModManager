# Satisfactory Mod Manager

Mod manager for [Satisfactory](https://www.satisfactorygame.com/).
Handles all the steps of installing mods for you.

## Installation and Usage

Check the [modding documentation](https://docs.ficsit.app/satisfactory-modding/latest/index.html#_for_users)
for installation and usage instructions.

## Features

* One-click install of any mod
* Automatically downloads the dependencies of any mod
* Mod update notifications
* Mod profiles and sharing of them
* View the mod description in the app

## Reporting issues

If you have any questions or run into issues, you can ask in the
[Satisfactory Modding discord](https://discord.gg/TShj39G)
for quicker responses than the GitHub issues.
If there is any error message, please include it along with the generated debug info zip.

## Troubleshooting

Check the [modding documentation](https://docs.ficsit.app/satisfactory-modding/latest/ForUsers/SatisfactoryModManager.html)
for common issues and their solutions.

* Profile and installation records are located in `%APPDATA%\SatisfactoryModManager\`
* Downloads are cached in `%LOCALAPPDATA%\SatisfactoryModManager\downloadCache\`

## Development

### Dependencies

* [Go 1.21](https://go.dev/doc/install)
* [pnpm](https://pnpm.io/installation)
* [nodejs](https://nodejs.org/en/download/)
* wails (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)
* IDE of Choice. Goland or VSCode suggested.

### Configuration

If you want to point to the SMR staging API instead of production, edit `main.go` to set `api-base` to the staging api (`https://api.ficsit.dev`)

### Building

```bash
wails build
```

### Linting

Install `golangci-lint` via the directions [here](https://golangci-lint.run/usage/install/#local-installation),
but make sure to install the version specified in `.github/workflows/push.yaml` instead of whatever it suggests.

Then, to run it, use:

```bash
golangci-lint run --fix
```

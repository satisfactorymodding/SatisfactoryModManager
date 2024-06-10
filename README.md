# Satisfactory Mod Manager

Mod manager for [Satisfactory](https://www.satisfactorygame.com/).
Handles all the steps of installing mods for you.

Implemented in [Wails](https://wails.io/).

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

* Profile and installation records are located in `%APPDATA%\ficsit\`
* Downloads are cached in `%LOCALAPPDATA%\ficsit\downloadCache\`
* Logs are stored in `%LOCALAPPDATA%\SatisfactoryModManager\logs`

## Development

### Dependencies

* [Go 1.22](https://go.dev/doc/install)
* [pnpm](https://pnpm.io/installation)
* [nodejs](https://nodejs.org/en/download/)
* Wails (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)
* IDE of Choice. Goland or VSCode suggested.

### Configuration

If you want to point to the SMR staging API instead of production, edit:

* `main.go` to set `api-base` to the staging api (`https://api.ficsit.dev`)
* `frontend\.graphqlrc` to set `schema` to the staging api (`https://api.ficsit.dev/v2/query`)
* `frontend\codegen.yml` to set `schema` to the staging api (`https://api.ficsit.dev/v2/query`)

### Development Server

The development server will hot reload whenever you make changes to the Go or Typescript code.

```bash
wails dev
```

Upon successful launch of the development server, the application will open automatically.
The command line output will also include a localhost URL you can visit in your browser if you wish to utilize browser developer tools for debugging.

Make sure you don't already have a copy of the application running or the command will silently fail.

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

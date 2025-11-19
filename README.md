# Satisfactory Mod Manager

Mod manager for [Satisfactory](https://www.satisfactorygame.com/).
Handles all the steps of installing mods for you.

Implemented in [Wails](https://wails.io/) using [Svelte](https://svelte.dev/) and [Skeleton](https://www.skeleton.dev/).

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

* [Go](https://go.dev/doc/install) (the version in [the Setup CI action `go-version`](./.github/actions/setup/action.yml))
* [pnpm](https://pnpm.io/installation)
* [nodejs](https://nodejs.org/en/download/) (the version in [the Setup CI action `node-version`](./.github/actions/setup/action.yml))
* Wails (to install, execute the `run` command from [the Setup CI action's `Install Wails` step](./.github/actions/setup/action.yml))
* IDE of Choice. Goland or VSCode suggested.

### Configuration

If you want to point to the SMR staging API instead of production, edit:

* `main.go` to set `api-base` to the staging api (`https://api.ficsit.dev`)
* `frontend\.graphqlrc.yml` to set `schema` to the staging api (`https://api.ficsit.dev/v2/query`)

### Development Server

The development server will hot reload whenever you make changes to the Go or Typescript code.

```bash
wails dev
```

Upon successful launch of the development server, the application will open automatically.
The command line output will also include a localhost URL you can visit in your browser if you wish to utilize browser developer tools for debugging.

Make sure you don't already have a copy of the application running or the command will silently fail.

Although `wails dev` should run these commands for you under normal circumstances,
you may need to run `pnpm graphql-codegen` in the `frontend` directory to update the code generated from the SMR API,
or run `pnpm translations` to update the translation data.

### IDE Configuration

Make sure that your IDE is connecting with the frontend's installation of ESLint to get the best experience.

VSCode users, a preconfigured workspace is provided in `.vscode/`
that allows editing both Go and Svelte files
while maintaining correct ESLint functionality.

### Building

```bash
wails build
```

To build a version that includes browser devtools (extending production debugging capabilities),:

```bash
wails build -devtools
```

### Linting

Install `golangci-lint` via the directions [in the Golangci-lint documentation](https://golangci-lint.run/docs/welcome/install/#local-installation),
but make sure to install the version specified in [`.github/workflows/push.yaml`](./.github/workflows/push.yml) instead of whatever it suggests.

Then, to run it, use:

```bash
golangci-lint run --fix
```

You may also need to manually run the frontend linter. First, navigate to the `frontend` directory, then run:

```bash
pnpm run format
pnpm run check
```

### Assorted Development Helpers

* Browser dev tools should automatically open when running `wails dev`.
* You can open the URL the frontend is being served on in a browser to use browser extensions like the Svelte Devtools.
  Check the logs for `Using DevServer URL:` to find the correct URL.
* You can manually trigger an error popup for testing purposes by running `debugCauseErrorMessage("your message here")` in the devtools console.

### Localization

If you'd like to help translate and localize SMM to different languages, join our [discord server](https://discord.ficsit.app/).

SMM handles localization through the Tolgee Svelte integration.
This allows for [in-context translation](https://tolgee.io/js-sdk/) - simply alt-click on a translatable element to open the Tolgee interface.

In order to edit translations in-context, you will need to provide a tolgee API key with edit permissions.
You can create an API key for yourself [in our Tolgee instance](https://translate.ficsit.app/projects/4/integrate) once you're added to the project.
To supply this API key at development time, create or edit `/frontend/.env.local` and supply the key in a similar format as `.env`.

The in-context translation screenshot feature requires installing the _Tolgee Tools_ browser extension.
After running `wails dev`, open `http://localhost:34115/` in your browser of choice to access the application.

## Code Signing Policy

<!-- markdownlint-disable -->
<table>
  <tr>
    <td colspan=3>
      Free code signing provided by <a href="https://about.signpath.io/">SignPath.io</a>, certificate by <a href="https://signpath.org/">SignPath Foundation</a>
    </td>
  </tr>
  <tr>
    <th>
      Committers and reviewers
    </th>
    <th>
      Approvers
    </th>
    <th>
      Privacy Policy
    </th>
  </tr>
  <tr>
    <td>
      <img src="https://github.com/satisfactorymodding.png?size=24" alt="Satisfactory Modding Logo" align="center" /><a href="https://github.com/orgs/satisfactorymodding/teams/members">Organization members</a>
    </td>
    <td>
      <img src="https://github.com/satisfactorymodding.png?size=24" alt="Satisfactory Modding Logo" align="center" /><a href="https://github.com/orgs/satisfactorymodding/people?query=role%3Aowner">Owners</a>
    </td>
    <td>
      <a href="https://ficsit.app/privacy-policy">https://ficsit.app/privacy-policy</a>
    </td>
  </tr>
</table>
<!-- markdownlint-enable -->

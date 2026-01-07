# Deployment

A [GitHub Action](https://github.com/satisfactorymodding/SatisfactoryModManager/blob/master/.github/workflows/push.yml)
automatically creates testing builds from commits on any branch and uploads them to the run as build artifacts.

## Development Builds

Since the build action is triggered on every push to any branch,
non developers can obtain testing builds easily from [GitHub Actions](https://github.com/satisfactorymodding/SatisfactoryModManager/actions/workflows/push.yml) artifact output.
You need to be signed into GitHub to download them, and they expire after some time.

TODO how to point dev builds to ficsit.dev (staging) instead?

## Releases

To create a new published release:

1. Ensure all intended changes are merged into the `master` branch.
2. Ensure someone with Signpath access is around to approve the release in the next few minutes.
   At the time of writing, that's Vilsol and Mircea.
3. Make a `Bump version` commit on the master branch to increase the `productVersion` field in `wails.json`. Version numbers should follow [Semantic Versioning](https://semver.org/).
4. Manually create and push a tag of the format `v*` to that commit, which triggers a
   [GitHub Action to make a release](https://github.com/satisfactorymodding/SatisfactoryModManager/blob/master/.github/workflows/release.yml).
5. The action will automatically request signing approval, sending an email to approvers and including the link in the action output.
   The action will busy-wait until approval is granted.
6. Review the draft release automatically created by Goreleaser and edit its description as needed.
7. Publish the release.

Users will automatically be prompted to update next time they launch SMM.

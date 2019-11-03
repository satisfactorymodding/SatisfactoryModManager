import semver from 'semver'

const isGood = function (version) {
  return semver.valid(version) || semver.valid(version + '.0')
}

export default function (version) {
  if (version.startsWith('v')) version = version.slice(1)
  let trimmedVersion = version
  while (!isGood(trimmedVersion)) {
    trimmedVersion = trimmedVersion.slice(0, -1)
  }
  if (!semver.valid(trimmedVersion)) trimmedVersion = trimmedVersion + '.0'
  return trimmedVersion
}

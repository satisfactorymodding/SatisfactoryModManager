import cliInterface from './cliInterface'
import { Model } from 'vue-mc'
import semver from 'semver'
const request = require('request')

const smlGitHubReleasesAPIurl = 'https://api.github.com/repos/satisfactorymodding/SatisfactoryModLoader/releases'

class SMLRelease extends Model {
  defaults () {
    return {
      tag_name: '',
      body: ''
    }
  }

  get changelog () {
    return this.body.replace(/#\s*Installation(.+\r?\n)*\r?\n/g, '', '')
  }
}

const cache = {
  smlVersions: []
}

const getInstalledSMLVersion = function (satisfactory) {
  return new Promise((resolve, reject) => {
    cliInterface.RunCommand('sml_version', `-p "${satisfactory.binariesDir}"`).then((output) => {
      resolve(output)
    })
  })
}

const installSML = function (satisfactory, version) {
  if (version) {
    return cliInterface.RunCommand('install_sml', `-p "${satisfactory.binariesDir}"`, `-v "${version}"`)
  } else {
    return cliInterface.RunCommand('install_sml', `-p "${satisfactory.binariesDir}"`)
  }
}

const uninstallSML = function (satisfactory) {
  return cliInterface.RunCommand('uninstall_sml', `-p "${satisfactory.binariesDir}"`)
}

const updateSML = function (satisfactory) {
  return cliInterface.RunCommand('update_sml', `-p "${satisfactory.binariesDir}"`)
}

const getAvailableSMLVersions = function () {
  return new Promise((resolve, reject) => {
    if (cache.smlVersions.length === 0) {
      rebuildCache().then(() => {
        resolve(cache.smlVersions)
      })
    } else {
      resolve(cache.smlVersions)
    }
  })
}

const rebuildCache = function () {
  return new Promise((resolve, reject) => {
    request({
      url: smlGitHubReleasesAPIurl, headers: { 'User-Agent': 'SatisfactoryModLauncher' }
    }, function (error, response, body) {
      if (error) { reject(error) }
      let smlGitHubVersions = JSON.parse(body)
      let finalVersions = []
      smlGitHubVersions.forEach(version => {
        finalVersions.push(new SMLRelease(version))
      })
      finalVersions.sort((a, b) => {
        return semver.gt(a.tag_name, b.tag_name)
      })
      cache.smlVersions = finalVersions
      resolve(cache)
    })
  })
}

const checkForUpdates = function () {
  rebuildCache()
}

const getLatestSMLVersion = function () {
  return new Promise((resolve, reject) => {
    getAvailableSMLVersions().then((versions) => resolve(versions[versions.length - 1]))
  })
}

export default {
  SMLRelease,
  getInstalledSMLVersion,
  installSML,
  uninstallSML,
  updateSML,
  getAvailableSMLVersions,
  getLatestSMLVersion,
  checkForUpdates,
  cache
}

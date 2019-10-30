import { spawn } from 'child_process'

const path = require('path')
const fs = require('fs')
const { getDataFolders } = require('platform-folders')
const EpicManifestsFolder = path.join(getDataFolders()[0], '\\Epic\\EpicGamesLauncher\\Data\\Manifests')
const { Model, Collection } = require('vue-mc')

class SatisfactoryInstall extends Model {
  // Default attributes that define the "empty" state.
  defaults () {
    return {
      id: null,
      name: '',
      version: '',
      installLocation: '',
      launchExecutable: ''
    }
  }

  // Attribute mutations.
  mutations () {
    return {
      id: (id) => Number(id) || null,
      name: String,
      version: String,
      installLocation: String,
      launchExecutable: String
    }
  }

  get launchPath () {
    return path.join(this.installLocation, this.launchExecutable)
  }

  get binariesDir () {
    return path.join(this.installLocation, 'FactoryGame', 'Binaries', 'Win64')
  }

  get displayName () {
    return `${this.name} (${this.version})`
  }

  get modsDir () {
    return path.join(this.binariesDir, 'mods')
  }

  launch () {
    spawn(this.launchPath, { detached: true }).unref()
  }
}

const getInstalls = function () {
  return new Promise((resolve, reject) => {
    let foundInstalls = new Collection()
    fs.readdirSync(EpicManifestsFolder).forEach(file => {
      if (file.endsWith('.item')) {
        file = path.join(EpicManifestsFolder, file)
        let jsonString = fs.readFileSync(file, 'utf8')
        let manifest = JSON.parse(jsonString)
        if (manifest.CatalogNamespace === 'crab') {
          foundInstalls.add(new SatisfactoryInstall({ name: manifest.DisplayName, version: manifest.AppVersionString, installLocation: manifest.InstallLocation, launchExecutable: manifest.LaunchExecutable }))
        }
      }
    })
    foundInstalls.sort(['version', 'name'])
    resolve(foundInstalls)
  })
}

export default { SatisfactoryInstall, getInstalls }

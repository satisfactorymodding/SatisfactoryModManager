import cliInterface from './cliInterface'
const path = require('path')
const fs = require('fs')
const { Model, Collection } = require('vue-mc')
const StreamZip = require('node-stream-zip')

let cachedModsDir = ''

const getModsDir = function () {
  return new Promise((resolve, reject) => {
    if (cachedModsDir) { return resolve(cachedModsDir) }
    cliInterface.RunCommand('mods_dir').then((output) => { cachedModsDir = output; resolve(cachedModsDir) })
  })
}

class ModData extends Model {
  // Default attributes that define the "empty" state.
  defaults () {
    return {
      mod_id: null,
      name: '',
      version: '',
      description: '',
      authors: new Collection(),
      objects: new Collection()
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
}

const getModData = function (modZipPath) {
  return new Promise((resolve, reject) => {
    const zip = new StreamZip({
      file: modZipPath,
      storeEntries: true
    })
    zip.on('error', () => {
      zip.close()
      resolve(null)
    })
    zip.on('ready', () => {
      const data = zip.entryDataSync('data.json')
      const dataJSON = JSON.parse(data.toString())
      zip.close()
      resolve(new ModData(dataJSON))
    })
  })
}

const isDirectory = source => fs.existsSync(source) && fs.lstatSync(source).isDirectory()
const isFile = source => fs.existsSync(source) && fs.lstatSync(source).isFile()
const getDirectories = source => fs.readdirSync(source).map(name => path.join(source, name)).filter(isDirectory)
const getFiles = source => fs.readdirSync(source).map(name => path.join(source, name)).filter(isFile)

const getDownloadedMods = function () {
  return new Promise((resolve, reject) => {
    let mods = new Collection()
    let modPaths = []
    getModsDir().then((modsDir) => {
      getDirectories(modsDir).forEach(dir => {
        getFiles(dir).forEach((file) => {
          if (file.endsWith('.zip')) { modPaths.push(file) }
        })
      })
      Promise.all(modPaths.map(mod => getModData(mod))).then((modsData) => {
        modsData.forEach(modData => { if (modData) mods.add(modData) })
        mods.sort(['name'])
        resolve(mods)
      })
    })
  })
}

const getInstalledMods = function (satisfactory) {
  return new Promise((resolve, reject) => {
    let mods = new Collection()
    let modPaths = []
    if (!fs.existsSync(satisfactory.modsDir)) { return resolve(mods) }
    getFiles(satisfactory.modsDir).forEach(file => {
      if (file.endsWith('.zip')) {
        modPaths.push(file)
      }
    })
    Promise.all(modPaths.map(mod => getModData(mod))).then((modsData) => {
      modsData.forEach(modData => mods.add(modData))
      mods.sort(['name'])
      resolve(mods)
    })
  })
}

const installModVersion = function (mod, satisfactory) {
  return cliInterface.RunCommand('install', `-m "${mod.mod_id}"`, `-v "${mod.version}"`, `-p "${satisfactory.binariesDir}"`)
}

const uninstallModVersion = function (mod, satisfactory) {
  return cliInterface.RunCommand('uninstall', `-m "${mod.mod_id}"`, `-v "${mod.version}"`, `-p "${satisfactory.binariesDir}"`)
}

const downloadModVersion = function (mod) {
  return cliInterface.RunCommand('download', `-m "${mod.mod_id}"`, `-v "${mod.version}"`)
}

const removeModVersion = function (mod) {
  return cliInterface.RunCommand('remove', `-m "${mod.mod_id}"`, `-v "${mod.version}"`)
}

const isModVersionDownloaded = function (mod, version) {
  return new Promise((resolve, reject) => {
    getDownloadedMods().then((mods) => {
      resolve(mods.has({ mod_id: mod, version: version }))
    })
  })
}

export default { ModData, getDownloadedMods, getInstalledMods, installModVersion, uninstallModVersion, downloadModVersion, removeModVersion, isModVersionDownloaded }

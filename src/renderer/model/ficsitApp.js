import { Model, Collection } from 'vue-mc'
import semver from 'semver'
import tryFixSemver from './tryFixSemver'

const API = require('graphql-client')({
  url: 'https://api.ficsit.app/v2/query'
})

const getAllModsQuery = `
query
{
  getMods(filter: {limit: 100})
  {
    mods
    {
      id,
      name,
      short_description,
      full_description,
      logo,
      authors
      {
        user
        {
          username
        }
      },
      versions
      {
        mod_id,
        version,
        sml_version,
        stability,
        link
      }
    }
  }
}
`
const getModVersionsQuery = `
query($modID: ModID!){
  getMod(modId: $modID)
  {
    versions
    {
      mod_id,
      version,
      sml_version,
      stability,
      link
    }
  }
}
`

class FicsitAppModVersion extends Model {
  defaults () {
    return {
      mod_id: '',
      version: '',
      sml_version: '',
      stability: '',
      link: '',
      isDownloaded: false,
      inProgress: false
    }
  }
}

class FicsitAppMod extends Model {
  defaults () {
    return {
      id: '',
      name: '',
      short_description: '',
      full_description: '',
      logo: '',
      authors: [],
      versions: new Collection([], { model: FicsitAppModVersion })
    }
  }
}

const getMods = function () {
  return new Promise((resolve, reject) => {
    API.query(getAllModsQuery, {}).then((body) => {
      const mods = new Collection()
      body.data.getMods.mods.forEach((mod) => {
        var newMod = new FicsitAppMod(mod)
        newMod.authors = newMod.defaults().authors
        newMod.versions = newMod.defaults().versions
        mod.authors.forEach((author) => newMod.authors.push(author.user.username))
        mod.versions.forEach((version) => newMod.versions.add(new FicsitAppModVersion(version)))
        mods.add(newMod)
      })
      resolve(mods)
    })
  })
}

const getModVersions = function (modID) {
  return new Promise((resolve, reject) => {
    API.query(getModVersionsQuery, { modID: modID }).then((body) => {
      const modVersions = []
      body.data.getMod.versions.forEach((version) => {
        modVersions.push(new FicsitAppModVersion(version))
      })
      modVersions.sort((a, b) => {
        return semver.gt(tryFixSemver(a.version), tryFixSemver(b.version))
      })
      resolve(modVersions)
    })
  })
}

const getLatestModVersion = function (modID) {
  return new Promise((resolve, reject) => {
    getModVersions(modID).then((versions) => {
      resolve(versions[versions.length - 1])
    })
  })
}

const noImageURL = 'https://ficsit.app/static/assets/images/no_image.png'

export default {
  FicsitAppMod,
  FicsitAppModVersion,
  getMods,
  getModVersions,
  getLatestModVersion,
  noImageURL
}

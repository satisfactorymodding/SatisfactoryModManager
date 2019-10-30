import { Model, Collection } from 'vue-mc'

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
      let mods = new Collection()
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

const noImageURL = 'https://ficsit.app/static/assets/images/no_image.png'

export default {
  FicsitAppMod,
  FicsitAppModVersion,
  getMods,
  noImageURL
}

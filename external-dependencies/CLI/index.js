const request = require('request')
const fs = require('fs')
const cliDependency = require('./dependency-config.json')

module.exports = {
  download: function (context) {
    return new Promise((resolve, reject) => {
      const downloadURL = cliDependency.url.replace('<version>', cliDependency.version).replace('<extension>', '.exe')
      if (fs.existsSync(cliDependency.path.replace('<extension>', '.exe'))) { return resolve() }
      console.log(`Downloading CLI ${cliDependency.version} from ${downloadURL}`)
      request(downloadURL)
        .pipe(fs.createWriteStream(cliDependency.path.replace('<extension>', '.exe'))).on('close', () => {
          console.log('Finished downloading CLI')
          resolve()
        }).on('error', error => {
          reject(error)
        })
    })
  },
  cleanup: function (context) {
    return new Promise((resolve, reject) => {
      fs.unlinkSync(cliDependency.path.replace('<extension>', '.exe'))
      console.log('Removed CLI')
      resolve()
    })
  }
}

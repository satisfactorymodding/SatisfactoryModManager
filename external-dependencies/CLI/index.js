const request = require('request')
const fs = require('fs')
const cliDependency = require('./dependency-config.json')

module.exports = function (context) {
  return new Promise((resolve, reject) => {
    if (!fs.existsSync(cliDependency.cacheFile)) {
      fs.writeFileSync(cliDependency.cacheFile, JSON.stringify({
        version: '0.0.0'
      }))
    }
    var cached = JSON.parse(fs.readFileSync(cliDependency.cacheFile, 'UTF8'))
    if (cached.version === cliDependency.version) { return resolve() }
    const downloadURL = cliDependency.url.replace('<version>', cliDependency.version).replace('<extension>', '.exe')
    if (fs.existsSync(cliDependency.path.replace('<extension>', '.exe'))) { return resolve() }
    console.log(`Downloading CLI ${cliDependency.version} from ${downloadURL}`)
    request(downloadURL)
      .pipe(fs.createWriteStream(cliDependency.path.replace('<extension>', '.exe'))).on('close', () => {
        console.log('Finished downloading CLI')
        cached.version = cliDependency.version
        fs.writeFileSync(cliDependency.cacheFile, JSON.stringify(cached))
        resolve()
      }).on('error', error => {
        reject(error)
      })
  })
}

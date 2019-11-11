module.exports = function (context) {
  const externalDependencies = require('./external-dependencies.json')
  console.log('Downloading external dependencies')
  return Promise.all(externalDependencies.map(dependency => require(dependency).download(context)))
}

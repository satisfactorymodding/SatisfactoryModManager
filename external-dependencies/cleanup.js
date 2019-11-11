module.exports = async function (context) {
  const externalDependencies = require('./external-dependencies.json')
  console.log('Cleaning up external dependencies')
  return Promise.all(externalDependencies.map(dependency => require(dependency).cleanup(context)))
}

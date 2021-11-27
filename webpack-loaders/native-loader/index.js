/*
  MIT License http://www.opensource.org/licenses/mit-license.php
  Author Tobias Koppers @sokra
*/

// Modified version of https://github.com/webpack-contrib/node-loader for electron, removing __webpack_public_path__ (https://github.com/webpack-contrib/node-loader/issues/37#issuecomment-796963091)

const { interpolateName } = require('loader-utils');
const schema = require('./options.json');

module.exports = function loader(content) {
  const options = this.getOptions(schema);

  const name = interpolateName(
    this,
    typeof options.name !== 'undefined' ? options.name : '[contenthash].[ext]',
    {
      context: this.rootContext,
      content,
    },
  );

  if (options.emit) {
    this.emitFile(name, content);
  }

  return `
try {
  process.dlopen(module, __dirname + require("path").sep + ${JSON.stringify(
    name,
  )}${
  typeof options.flags !== 'undefined'
    ? `, ${JSON.stringify(options.flags)}`
    : ''
});
} catch (error) {
  throw new Error('node-loader:\\n' + error);
}
`;
};

module.exports.raw = true;

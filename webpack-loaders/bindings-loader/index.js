// Modified version of https://github.com/alessiopcc/node-bindings-loader to support typescript transpiled bindings imports (`bindings_1.default('module')`)
// @ts-nocheck
const {OriginalSource, SourceMapSource, ReplaceSource} = require("webpack-sources");
const {runInNewContext} = require('vm');
const path = require('path');

const SUPPORTED_PACKAGES = {
    'bindings': {handler: _bindings},
    'node-gyp-build': {handler: _node_gyp_build},
    'node-bindings-loader': {method: 'glob', handler: _custom_glob},
}

function _bindings(loader, match, code)
{
    return new Promise((resolve, reject) =>
    {
        loader.resolve(loader.context, 'bindings', (error, module_path) =>
        {
            if(error)
                return reject(error);

            try
            {
                const root = path.dirname(loader.resourcePath);

                const node_module = require(module_path);

                const args = {
                    bindings: runInNewContext(match[1], {
                        __dirname: root,
                        __filename: loader.resourcePath,
                    }),
                    path: true,
                    module_root: node_module.getRoot(loader.resourcePath),
                };
                
                const resolve_path = path.relative(root, node_module(args)).replace(/\\/g, '/');

                code.replace(match.index, match.index + match[0].length - 1, `require('./${resolve_path}')`);
            }
            catch(module_error)
            {
                return reject(module_error);
            }

            return resolve();
        });
    });
}

function _node_gyp_build(loader, match, code)
{
    return new Promise((resolve, reject) =>
    {
        loader.resolve(loader.context, 'node-gyp-build', (error, module_path) =>
        {
            if(error)
                return reject(error);

            try
            {
                const root = path.dirname(loader.resourcePath);

                const node_module = require(module_path);

                const args = runInNewContext(match[1], {
                    __dirname: root,
                    __filename: loader.resourcePath,
                });

                const resolve_path = path.relative(root, node_module.path(args)).replace(/\\/g, '/');
                code.replace(match.index, match.index + match[0].length - 1, `require('./${resolve_path}')`);
            }
            catch(module_error)
            {
                return reject(module_error);
            }

            return resolve();
        });
    });
}

function _custom_glob(loader, match, code)
{
    return new Promise((resolve, reject) =>
    {
        loader.resolve(loader.context, 'glob', (error, module_path) =>
        {
            if(error)
                return reject(error);

            try
            {
                const root = path.dirname(loader.resourcePath);

                const node_module = require(module_path);

                const args = runInNewContext(match[1], {
                    __dirname: root,
                    __filename: loader.resourcePath,
                });

                const binding = node_module.sync(args, {cwd: root})[0];

                if(!binding)
                    throw new Error('Glob cannot find module');

                const resolve_path = binding.replace(/\\/g, '/');
                code.replace(match.index, match.index + match[0].length - 1, `require('./${resolve_path}')`);
            }
            catch(module_error)
            {
                return reject(module_error);
            }

            return resolve();
        });
    });
}

async function run(source, map)
{

    const callback = this.async();

    const balanced_parenthesis_regex = '\\(((?:[^)(]+|\\((?:[^)(]+|\\([^)(]*\\))*\\))*)\\)'

    const code = new ReplaceSource(map ? new SourceMapSource(source, this.resourcePath, map) : new OriginalSource(source, this.resourcePath));

    try
    {
        for(const package of Object.keys(SUPPORTED_PACKAGES))
        {
            const method = SUPPORTED_PACKAGES[package].method ? `\\.${SUPPORTED_PACKAGES[package].method}` : '';

            const regex = new RegExp(`\\b(?:require\\((?:'|")${package}(?:'|")\\)${method}|${package}_1.default)\\s*${balanced_parenthesis_regex}`, 'g');

            while(match = regex.exec(source))
                await SUPPORTED_PACKAGES[package].handler(this, match, code);
        }
    }
    catch(error)
    {
        return callback(error);
    }

    const loader_code = code.sourceAndMap();
    return callback(null, loader_code.source, loader_code.map);
};

function glob(pattern)
{
    const glob_module = require('glob');
    const root = path.dirname(require('parent-module')())

    const binding = glob_module.sync(pattern, {cwd: root})[0];

    if(!binding)
        throw new Error('Glob cannot find module');

    return require(path.resolve(root, binding));
}

module.exports = run;
module.exports.glob = glob;
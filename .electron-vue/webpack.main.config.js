'use strict'

process.env.BABEL_ENV = 'main'

const path = require('path');
const { dependencies } = require('../package.json');
const webpack = require('webpack');
const ESLintPlugin = require('eslint-webpack-plugin');
const TerserPlugin = require('terser-webpack-plugin');

let mainConfig = {
  entry: {
    main: path.join(__dirname, '../src/main/index.js')
  },
  externals: [
    'bindings'
  ],
  module: {
    rules: [
      {
        test: /\.js$/,
        use: 'babel-loader',
        exclude: /node_modules/
      },
      {
        test: /\.js$/,
        use: path.resolve('webpack-loaders/bindings-loader'),
      },
      {
        test: /\.node$/,
        use: {
          loader: path.resolve('webpack-loaders/native-loader'),
          options: {
            name: process.env.NODE_ENV !== 'production' ? '[name].[ext]' : '[name]-[hash].[ext]', // Use original in dev
            emit: process.env.NODE_ENV === 'production' // Do not emit in dev, needed when using linked packages
          }
        }
      },
    ]
  },
  node: {
    __dirname: process.env.NODE_ENV !== 'production',
    __filename: process.env.NODE_ENV !== 'production'
  },
  output: {
    filename: '[name].js',
    libraryTarget: 'commonjs2',
    path: path.join(__dirname, '../dist/electron')
  },
  plugins: [
    new webpack.NoEmitOnErrorsPlugin(),
    new ESLintPlugin({
      formatter: require('eslint-friendly-formatter')
    }),
  ],
  resolve: {
    extensions: ['.js', '.json', '.node']
  },
  target: 'electron-main',
  optimization: {},
}

/**
 * Adjust mainConfig for development settings
 */
if (process.env.NODE_ENV !== 'production') {
  mainConfig.plugins.push(
    new webpack.DefinePlugin({
      '__static': `"${path.join(__dirname, '../static').replace(/\\/g, '\\\\')}"`
    })
  )
}

/**
 * Adjust mainConfig for production settings
 */
if (process.env.NODE_ENV === 'production') {
  mainConfig.plugins.push(
    new webpack.DefinePlugin({
      'process.env.NODE_ENV': '"production"'
    })
  )
  mainConfig.optimization.minimize = true;
  mainConfig.optimization.minimizer = [new TerserPlugin()];
}

module.exports = mainConfig

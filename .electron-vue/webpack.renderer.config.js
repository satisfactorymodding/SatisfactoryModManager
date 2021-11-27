'use strict'

process.env.BABEL_ENV = 'renderer'

const path = require('path')
const webpack = require('webpack')

const TerserPlugin = require('terser-webpack-plugin');
const CopyWebpackPlugin = require('copy-webpack-plugin')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const { VueLoaderPlugin } = require('vue-loader')
const ESLintPlugin = require('eslint-webpack-plugin');

let rendererConfig = {
  devtool: 'source-map',
  entry: {
    renderer: path.join(__dirname, '../src/renderer/main.js')
  },
  externals: [
    'bindings',
    'ps-list',
  ],
  module: {
    rules: [
      {
        test: /\.vue$/,
        use: {
          loader: 'vue-loader',
          options: {
            extractCSS: process.env.NODE_ENV === 'production',
            loaders: {
              sass: 'vue-style-loader!css-loader!sass-loader?indentedSyntax=1',
              scss: 'vue-style-loader!css-loader!sass-loader',
              less: 'vue-style-loader!css-loader!less-loader'
            }
          }
        }
      },
      {
        test: /\.s(c|a)ss$/,
        use: [
          'vue-style-loader',
          'css-loader',
          {
            loader: 'sass-loader',
            options: {
              implementation: require('sass'),
              sassOptions: {
                indentedSyntax: true // optional
              },
            },
          },
        ],
      },
      {
        test: /\.less$/,
        use: ['vue-style-loader', 'css-loader', 'less-loader']
      },
      {
        test: /\.css$/,
        use: ['vue-style-loader', 'css-loader']
      },
      {
        test: /\.html$/,
        use: 'vue-html-loader'
      },
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
      {
        test: /\.(png|jpe?g|gif|svg)(\?.*)?$/,
        type: 'asset',
        generator: {
          filename: 'imgs/[name]--[folder].[ext]'
        }
      },
      {
        test: /\.(mp4|webm|ogg|mp3|wav|flac|aac)(\?.*)?$/,
        type: 'asset',
        generator: {
          filename: 'media/[name]--[folder].[ext]'
        }
      },
      {
        test: /\.(woff2?|eot|ttf|otf)(\?.*)?$/,
        type: 'asset',
        generator: {
          filename: 'fonts/[name]--[folder].[ext]'
        }
      }
    ]
  },
  node: {
    __dirname: process.env.NODE_ENV !== 'production',
    __filename: process.env.NODE_ENV !== 'production'
  },
  plugins: [
    new VueLoaderPlugin(),
    new MiniCssExtractPlugin({ filename: 'styles.css' }),
    new HtmlWebpackPlugin({
      filename: 'index.html',
      template: path.resolve(__dirname, '../src/index.ejs'),
      templateParameters(compilation, assets, options) {
        return {
          compilation: compilation,
          webpack: compilation.getStats().toJson(),
          webpackConfig: compilation.options,
          htmlWebpackPlugin: {
            files: assets,
            options: options
          },
          process,
        };
      },
      minify: {
        collapseWhitespace: true,
        removeAttributeQuotes: true,
        removeComments: true
      },
      nodeModules: process.env.NODE_ENV !== 'production'
        ? path.resolve(__dirname, '../node_modules')
        : false
    }),
    new webpack.NoEmitOnErrorsPlugin(),
    new ESLintPlugin({
      formatter: require('eslint-friendly-formatter')
    }),
  ],
  output: {
    filename: '[name].js',
    libraryTarget: 'commonjs2',
    path: path.join(__dirname, '../dist/electron')
  },
  resolve: {
    alias: {
      '@': path.join(__dirname, '../src/renderer'),
      '~': path.join(__dirname, '../src'),
      'vue$': 'vue/dist/vue.esm.js',
      'jszip': require.resolve('jszip/lib/index.js'), // https://github.com/Stuk/jszip/issues/663#issuecomment-675737399
    },
    extensions: ['.js', '.vue', '.json', '.css', '.node'],
    symlinks: false,
  },
  target: 'electron-renderer',
  optimization: {},
}

/**
 * Adjust rendererConfig for development settings
 */
if (process.env.NODE_ENV !== 'production') {
  rendererConfig.plugins.push(
    new webpack.DefinePlugin({
      '__static': `"${path.join(__dirname, '../static').replace(/\\/g, '\\\\')}"`
    })
  )
}

/**
 * Adjust rendererConfig for production settings
 */
if (process.env.NODE_ENV === 'production') {
  rendererConfig.plugins.push(
    new CopyWebpackPlugin({ patterns: [
      {
        from: path.join(__dirname, '../static'),
        to: path.join(__dirname, '../dist/electron/static'),
        globOptions: {
          ignore: ['.*'],
        }
      },
      {
        from: path.join(__dirname, '../node_modules/win-ca/lib/roots.exe'),
        to: path.join(__dirname, '../dist/electron/roots.exe')
      }
    ]}),
    new webpack.DefinePlugin({
      'process.env.NODE_ENV': '"production"'
    }),
    new webpack.LoaderOptionsPlugin({
      minimize: true
    })
  );
  rendererConfig.optimization.minimize = true;
  rendererConfig.optimization.minimizer = [
    new TerserPlugin({
      terserOptions: {
        sourceMap: true
      }
    })
  ];
}

module.exports = rendererConfig

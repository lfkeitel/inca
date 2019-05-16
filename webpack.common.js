const path = require('path');
const CleanWebpackPlugin = require('clean-webpack-plugin');

module.exports = {
  entry: {
    main: './frontend/source/scripts/main.js',
    deviceListPage: './frontend/source/scripts/deviceListPage.js',
    archive: './frontend/source/scripts/archive.js'
  },

  plugins: [
    new CleanWebpackPlugin()
  ],

  output: {
    filename: '[name].min.js',
    path: path.resolve(__dirname, 'frontend/dist/js')
  },

  module: {
    rules: [
      {
        test: /\.tsx?$/,
        use: 'ts-loader',
        exclude: /node_modules/
      }
    ]
  }
};

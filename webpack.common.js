const path = require('path');
const CleanWebpackPlugin = require('clean-webpack-plugin');

module.exports = {
  entry: {
    main: './frontend/src/scripts/main.ts',
    deviceListPage: './frontend/src/scripts/deviceListPage.ts',
    archive: './frontend/src/scripts/archive.ts'
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
  },

  resolve: {
    extensions: ['.tsx', '.ts', '.js']
  },
};

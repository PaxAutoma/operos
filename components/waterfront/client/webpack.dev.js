/*
Copyright 2018 Pax Automa Systems, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

const webpack = require('webpack');
const merge = require('webpack-merge');
const common = require('./webpack.common.js');

module.exports = merge(common.config, {
  devtool: 'inline-source-map',
  devServer: {
    contentBase: common.srcDir,
    historyApiFallback: true,
    port: 10000
  },
  plugins: [
    new webpack.EnvironmentPlugin({
      API_BASE_URL: 'http://localhost:2780/api/v1',
      KUBE_DASHBOARD_URL: '/kube-dashboard/'
    })
  ]
});

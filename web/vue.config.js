const path = require('path')

module.exports = {
  productionSourceMap: process.env.VUE_APP_ENV !== 'prod',
  "transpileDependencies": [
    "vuetify"
  ],
  pluginOptions: {
    webpackBundleAnalyzer: {
      openAnalyzer: false,
      analyzerMode: process.env.VUE_APP_ENV !== 'prod' ? 'server' : ''
    }
  },
  configureWebpack: {
    optimization: {
      splitChunks: {
        minSize: 10000,
        maxSize: 250000,
      }
    }
  },
}
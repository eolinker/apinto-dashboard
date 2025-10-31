const mf = require('@angular-architects/module-federation/webpack')
const path = require('path')
const sharedMappings = new mf.SharedMappings()
sharedMappings.register(path.join(__dirname, '../../tsconfig.json'), [
  /* mapped paths to share */
])

module.exports = mf.withModuleFederationPlugin({

  name: 'main',
  filename: 'apinto.js',
  shared: {
    ...mf.shareAll({ singleton: true, strictVersion: true, requiredVersion: 'auto' })
  }

})

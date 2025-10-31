export const PLUGIN_CONFIG = [
  {
    name: 'guide',
    driver: 'apinto.builtIn.component',
    router: [
      { path: 'guide', type: 'normal' }
    ]
  },
  {
    name: 'cluster',
    driver: 'apinto.builtIn.module',
    router: [
      { path: 'deploy/cluster', type: 'normal' }
    ]
  },
  {
    name: 'global-env',
    driver: 'apinto.builtIn.module',
    router: [
      { path: 'deploy/variable', type: 'normal' }
    ]
  },
  {
    name: 'node-plugin',
    driver: 'apinto.builtIn.module',
    router: [
      { path: 'deploy/plugin', type: 'normal' }
    ]
  }, {
    name: 'application',
    driver: 'apinto.builtIn.module',
    router: [
      { path: 'application', type: 'normal' }
    ]
  }, {
    name: 'api',
    driver: 'apinto.builtIn.module',
    router: [
      { path: 'router/api', type: 'normal' }
    ]
  }, {
    name: 'plugin-template',
    driver: 'apinto.builtIn.module',
    router: [
      { path: 'router/plugin-template', type: 'normal' }
    ]
  }, {
    name: 'traffic-strategy',
    driver: 'apinto.builtIn.module',
    router: [
      { path: 'serv-governance/traffic', type: 'normal' }
    ]
  }, {
    name: 'fuse-strategy',
    driver: 'apinto.builtIn.module',
    router: [
      { path: 'serv-governance/fuse', type: 'normal' }
    ]
  }, {
    name: 'visit-strategy',
    driver: 'apinto.builtIn.module',
    router: [
      { path: 'serv-governance/visit', type: 'normal' }
    ]
  }, {
    name: 'cache-strategy',
    driver: 'apinto.builtIn.module',
    router: [
      { path: 'serv-governance/cache', type: 'normal' }
    ]
  }, {
    name: 'grey-strategy',
    driver: 'apinto.builtIn.module',
    router: [
      { path: 'serv-governance/grey', type: 'normal' }
    ]
  },
  {
    name: 'open-api',
    driver: 'apinto.builtIn.module',
    router: [
      { path: 'system/ext-app', type: 'normal' }
    ]
  }, {
    name: 'email',
    driver: 'apinto.builtIn.component',
    router: [
      { path: 'system/email', type: 'normal' }
    ]
  }, {
    name: 'webhook',
    driver: 'apinto.builtIn.component',
    router: [
      { path: 'system/webhook', type: 'normal' }
    ]
  }, {
    name: 'audit-log',
    driver: 'apinto.builtIn.module',
    router: [
      { path: 'audit-log', type: 'normal' }
    ]
  }, {
    name: 'module-plugin',
    driver: 'apinto.builtIn.module',
    router: [
      { path: 'module-plugin', type: 'normal' }
    ]
  }, {
    name: 'log',
    driver: 'apinto.builtIn.module',
    router: [
      { path: 'log', type: 'normal' }
    ]
  }, {
    name: 'apispace',
    driver: 'apinto.remote.normal',
    router: [
      { path: 'remote/apispace', type: 'normal' }
    ]
  }, {
    name: 'discovery',
    driver: 'apinto.intelligent.normal',
    router: [
      { path: 'template/discovery', type: 'normal' }
    ]
  },
  {
    name: 'user',
    driver: 'apinto.local.preload',
    router: [
      { path: 'login', expose: 'LoginModule', type: 'root' },
      { path: 'user', expose: 'AppModule', type: 'normal' }
    ]
  },
  {
    name: 'auth',
    driver: 'apinto.local.preload',
    router: [
      { path: 'auth', expose: 'AppModule', type: 'root' },
      { path: 'auth-info', expose: 'AuthInfoModule', type: 'normal' }
    ]
  },
  {
    name: 'monitor',
    driver: 'apinto.local.router',
    router: [
      { path: 'monitor', expose: 'AppModule', type: 'normal' }
    ]
  }
]

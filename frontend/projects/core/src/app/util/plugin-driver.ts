import { Route } from '@angular/router'
import { AppComponent } from '../app.component'
import { routerMap } from '../constant/app.config'
import { PluginLoaderService } from '../service/plugin-loader.service'
import { PluginProviderService } from '../service/plugin-provider.service'
import { PluginLifecycleGuard } from './plugin-lifecycle.guard'

export const DEFAULT_LOCAL_PLUGIN_PATH = '/plugin-frontend/'
interface PluginRouterConfig {
  path:string;
  type:string;
  expose?:string
}
interface PluginConfig {
  name: string;
  router: Array<PluginRouterConfig>;
  path?: string;
  driver:string
}
interface CoreObj {
  routerConfig: Route[];
  executeList: any[];
  pluginLoader: PluginLoaderService;
  pluginLifecycleGuard: PluginLifecycleGuard;
  pluginProvider: PluginProviderService;
  builtInPluginLoader: (name: string) => any;
}

const defaultBuiltInPlugin:Array<{path:string, pathMatch?:'full' | 'prefix', componentName?:string, type?:string}> = [
  { path: '', pathMatch: 'full', componentName: 'redirectPage' },
  { path: '', componentName: 'basicLayout' },
  { path: 'dynamic-demo', componentName: 'dynamic-demo' },
  { path: '**', componentName: 'redirectPage' }

]

export const apinto:{[key:string]:any} = {
  builtIn: {
    // apinto主项目驱动，在core中自动调，不根据插件配置表
    default: (coreObj:CoreObj) => {
      const url = new URL(window.location.href)
      const navHidden = url.searchParams.get('nav_hidden') || sessionStorage.getItem('nav_hidden')
      if (navHidden === 'true') sessionStorage.setItem('nav_hidden', navHidden)
      // TODO Test
      const routes = defaultBuiltInPlugin.map(plugin => (
        {
          path: plugin.path,
          component: navHidden === 'true' && plugin.componentName === 'basicLayout' ? routerMap.get('nav-hidden').component : routerMap.get(plugin.componentName || plugin.path)?.component,
          children: [],
          data: {
            type: plugin.type || plugin.componentName || plugin.path
          },
          pathMatch: plugin.pathMatch || 'prefix'
        })
      )
      coreObj.routerConfig.push(...routes)
      return coreObj
    },
    component: (coreObj:CoreObj, pluginConfig:PluginConfig) => {
      for (const pluginRouter of pluginConfig.router) {
        coreObj.pluginProvider.setRouterConfig(pluginRouter.type === 'root', {
          path: pluginRouter.path,
          component: routerMap.get(pluginConfig.name).component
        }, coreObj.routerConfig)
      }
      return coreObj
    },
    module: (coreObj:CoreObj, pluginConfig:PluginConfig) => {
      for (const pluginRouter of pluginConfig.router) {
        coreObj.pluginProvider.setRouterConfig(pluginRouter.type === 'root', {
          path: pluginRouter.path,
          component: AppComponent,
          loadChildren: coreObj.builtInPluginLoader(pluginConfig.name)
        }, coreObj.routerConfig)
      }
      return coreObj
    },
    httpApi: () => {}
  },
  remote: {
    normal: (coreObj:CoreObj, pluginConfig:PluginConfig) => {
      const remoteRouter = coreObj.routerConfig.find((item:Route) => item?.data?.['type'] === 'remotePlugin')
      if (!remoteRouter) {
        coreObj.pluginProvider.setRouterConfig(false, {
          path: 'remote',
          children: [
            {
              path: ':moduleName',
              component: routerMap.get('remote').component
            }
          ],
          data: {
            type: 'remotePlugin'
          }
        }, coreObj.routerConfig)
      }
      return coreObj
    }
  },
  intelligent: {
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    normal: (coreObj:CoreObj, pluginConfig:PluginConfig) => {
      const remoteRouter = coreObj.routerConfig.find((item:Route) => item?.data?.['type'] === 'intelligentPlugin')
      if (!remoteRouter) {
        coreObj.pluginProvider.setRouterConfig(false, {
          path: 'template',
          loadChildren: coreObj.builtInPluginLoader('intelligent'),
          data: {
            type: 'intelligentPlugin'
          }
        }, coreObj.routerConfig)
      }
      return coreObj
    }
  },
  local: {
    router: (coreObj:CoreObj, pluginConfig:PluginConfig) => {
      for (const pluginRouter of pluginConfig.router) {
        if (pluginRouter.type === 'sub') {
          continue
        }
        updateRouterConfigWithPlugin(coreObj, pluginRouter, pluginConfig)
      }
      return coreObj
    },
    preload: (coreObj:CoreObj, pluginConfig:PluginConfig) => {
      coreObj.executeList.push({ ...pluginConfig, expose: 'Bootstrap', bootstrap: 'BootstrapModule.bootstrap' })
      for (const pluginRouter of pluginConfig.router) {
        updateRouterConfigWithPlugin(coreObj, pluginRouter, pluginConfig)
      }
      return coreObj
    }
    // extender: (coreObj:CoreObj, pluginConfig:PluginConfig) => {}
  }
}
function updateRouterConfigWithPlugin (coreObj: CoreObj, pluginRouter: PluginRouterConfig, pluginConfig: PluginConfig) {
  if (!pluginRouter.expose) {
    throw new Error('pluginRouter.expose is required')
  } else {
    coreObj.pluginProvider.setRouterConfig(pluginRouter.type === 'root', {
      path: pluginRouter.path,
      loadChildren: () => coreObj.pluginLoader.loadModule(
        pluginRouter.path,
        pluginConfig.name,
        pluginRouter.expose!,
        pluginConfig.path || `${DEFAULT_LOCAL_PLUGIN_PATH}${pluginConfig.name}/apinto.js`
      ),
      canActivate: [coreObj.pluginLifecycleGuard],
      canActivateChild: [coreObj.pluginLifecycleGuard],
      canDeactivate: [coreObj.pluginLifecycleGuard],
      canLoad: [coreObj.pluginLifecycleGuard]
    }, coreObj.routerConfig)
  }
}

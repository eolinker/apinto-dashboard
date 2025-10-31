import { loadRemoteModule } from '@angular-architects/module-federation'
import { Injectable, Injector } from '@angular/core'
import { Router } from '@angular/router'
import { generateRemoteModuleTemplate, validateExportLifecycle } from '../util/module-federation'
import { PluginProviderService } from './plugin-provider.service'
import { EoNgNavigationService } from './eo-ng-navigation.service'
import { PluginLifecycleGuard } from '../util/plugin-lifecycle.guard'
import { DEFAULT_LOCAL_PLUGIN_PATH, apinto } from '../util/plugin-driver'
import { PluginEventHubService } from './plugin-event-hub.service'
import { PluginSlotHubService } from './plugin-slot-hub.service'
import { EoNgMessageService } from './eo-ng-message.service'
import { EoNgFeedbackModalService } from 'eo-ng-feedback'
import { NzModalService } from 'ng-zorro-antd/modal'
import { PlatformProviderService } from './platform-provider.service'
import { routerMap } from '../constant/app.config'
import { ApiService } from './api.service'
import { BaseInfoService } from './base-info.service'

@Injectable({
  providedIn: 'root'
})
export class PluginLoaderService {
  private modules:Map<string, Object> = new Map<string, Object>()
  private baseHref: string = ''; // 传递给模块的 baseHref
  private executeList: any[] = []
  redirectUrl:string = ''

  constructor (
    private api:ApiService,
    private router: Router,
    private baseInfo:BaseInfoService,
    private pluginProvider:PluginProviderService,
    private navigation:EoNgNavigationService,
    private pluginEventHub:PluginEventHubService,
    private pluginSlotHub:PluginSlotHubService,
    private platformProvider:PlatformProviderService,
    private injector:Injector,
    private message:EoNgMessageService,
    private modalService:NzModalService,
    private eoModalService:EoNgFeedbackModalService) {
  }

  getExecuteList () {
    return this.executeList
  }

  // 涉及路由修改，所以在在项目运行之初通过APP_INITIALIZER执行
  loadPlugins:()=>Promise<void> = () => {
    // eslint-disable-next-line no-async-promise-executor
    return new Promise((resolve) => {
      const routerConfig: any[] = []
      // 防止access接口报错时无路由
      apinto['builtIn'].default({ routerConfig })
      this.router.resetConfig(routerConfig)
      // 需要先获取用户权限后安装插件，否则传递数据为空?
      // const subscription = this.navigation.getMenuList().subscribe(async () => {
      //   subscription.unsubscribe()
      // })
      this.installPlugin(routerConfig, this.executeList).then(
        async () => {
          this.router.resetConfig(routerConfig)
          await this.loadExecutedPlugin()
          resolve() // 当一切完成后，解析 Promise
        }
      )
    })
  }

  // 根据接口给的插件配置列表安装插件
  installPlugin = (routerConfig:any, executeList:any) : Promise<void> => {
    return new Promise((resolve, reject) => {
      this.api.get('system/plugins').subscribe((resp:any) => {
        if (resp.code === 0) {
          // 更新版本号与构建时间
          this.baseInfo.updateDate = resp.data.buildAt
          this.baseInfo.version = resp.data.version
          this.baseInfo.powered = resp.data.powered
          this.baseInfo.product = resp.data.product
          this.baseInfo.showGuide = resp.data.guide
          const driverMethod:{[key:string]:any} = { apinto: apinto }
          const pluginConfigList = resp.data.plugins
          const pluginLoader = { loadModule: this.loadModule }
          const pluginLifecycleGuard = PluginLifecycleGuard
          const builtInPluginLoader = this.loadBuiltInModule
          const pluginProvider = this.pluginProvider
          this.pluginSlotHub.addSlot('renewMenu', () => {
            this.navigation.dataUpdated = true
            const subscription = this.navigation.getMenuList().subscribe(async () => {
              subscription.unsubscribe()
            })
          })

          for (const plugin of pluginConfigList) {
            try {
              const driverName:string = plugin.driver!
              if (!driverName) {
                console.error(' no driver name')
                continue
              }
              const driver:any = driverName.split('.').reduce((driverMethod, driverName) => driverMethod[driverName], driverMethod)
              driver({ routerConfig, executeList, pluginLoader, pluginProvider, pluginLifecycleGuard, builtInPluginLoader }, plugin)
            } catch (err:any) {
              console.warn('安装插件出错：', err)
            }
          }
          this.router.resetConfig(routerConfig)
          resolve()
        } else {
          this.message.error(resp.msg || '获取插件配置列表失败，请重试!')
          reject(new Error(resp.msg || '获取插件配置列表失败'))
        }
      })
    })
  }

  // 加载内置插件(模块级)
  loadBuiltInModule = (pluginName:string) => {
    try {
      const { module } = routerMap.get(pluginName)
      return module
    } catch (err:any) {
      console.warn(`安装内置插件[${pluginName}]出错：`, err)
    }
  }

   loadModule =async (routerPrefix:string, pluginName:string, exposedModule:string, pluginPath:string) => {
     if (!this.modules.get(routerPrefix)) {
       try {
         const Module = await loadRemoteModule(generateRemoteModuleTemplate(pluginName, exposedModule, pluginPath))
         this.modules.set(routerPrefix, Module)
         if (!validateExportLifecycle(Module)) {
           console.error(' 需要导出插件生命周期函数')
           return
         }
         await Module.bootstrap?.({
           pluginProvider: this.pluginProvider,
           pluginEventHub: this.pluginEventHub.initHub(),
           pluginSlotHub: this.pluginSlotHub
         })
         return Module[exposedModule]
       } catch (error) {
         console.error(' 导入插件失败：', error)
       }
     }
     return this.getModule(routerPrefix, true)[exposedModule]
   }

   getModule (routerPrefix: string, specific?:boolean):any {
     if (routerPrefix.startsWith('/')) {
       routerPrefix = routerPrefix.substring(1)
     }
     if (specific) {
       return this.modules.get(routerPrefix)
     }
     let matchedModule = null
     let matchedLength = 0

     this.modules.forEach((value, key) => {
       // 检查是否为前缀且为更长的匹配
       if (routerPrefix.startsWith(key) && key.length > matchedLength) {
         matchedModule = value
         matchedLength = key.length
       }
     })
     return matchedModule
   }

   // 执行 立即执行队列里的插件导出的方法
   // 默认执行Bootstrap模块的bootstrap方法，如plugin配置有expose和bootstrap则执行对应方法
   loadExecutedPlugin = async () => {
     for (const plugin of this.executeList) {
       try {
         const Module = await loadRemoteModule(generateRemoteModuleTemplate(plugin.name, plugin?.expose || 'Bootstrap', plugin.path || `${DEFAULT_LOCAL_PLUGIN_PATH}${plugin.name}/apinto.js`))
         const bootstrap = Module.bootstrap
         if (!bootstrap) {
           console.warn(' 立即执行插件未导出Bootstrap模块或bootstrap函数')
         } else {
           await bootstrap(
             {
               pluginEventHub: this.pluginEventHub.initHub(),
               pluginSlotHub: this.pluginSlotHub,
               pluginProvider: this.pluginProvider,
               platformProvider: this.platformProvider,
               closeModal: this.closeModal,
               router: this.router,
               messageService: this.message,
               modalService: this.modalService,
               injector: this.injector,
               apiService: this.api
             }
           )
         }
       } catch (error) {
         console.error(' 执行插件失败：', error)
       }
     }
   }

   closeModal = () => {
     this.modalService.closeAll()
     this.eoModalService.closeAll()
   }

   setBaseHref (baseHref: string) {
     this.baseHref = baseHref
   }

   getBaseHref () {
     return this.baseHref
   }
}

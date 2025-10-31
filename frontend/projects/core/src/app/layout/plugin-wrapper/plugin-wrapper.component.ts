/* eslint-disable no-useless-catch */
import { loadRemoteModule } from '@angular-architects/module-federation'
import { Compiler, Component, ElementRef, Injector, NgModuleFactory, NgModuleRef, Type, ViewChild, ViewContainerRef, createNgModule } from '@angular/core'
import { ActivatedRoute } from '@angular/router'
import { generateRemoteModuleTemplate } from '../../util/module-federation'
import { PlatformProviderService } from '../../service/platform-provider.service'

@Component({
  selector: 'eo-ng-plugin-wrapper',
  template: `
  `,
  styles: [
  ]
})
export class PluginWrapperComponent {
  @ViewChild('modulePlaceHolder', { read: ViewContainerRef, static: true }) modulePlaceHolder:any
  private userApp: any;
  private moduleRef: any;

  constructor (private platformProvider: PlatformProviderService, private el: ElementRef, private activatedRoute:ActivatedRoute, private compiler: Compiler, private injector:Injector) {}
  async ngOnInit () {
    const { pluginPath, pluginExpose, pluginName } = this.activatedRoute.snapshot.data
    if (!pluginPath || !pluginExpose) {
      console.error(' 插件配置有误，请检查')
      return
    }
    // loadRemoteModule(generateRemoteModuleTemplate(pluginName, pluginExpose, pluginPath))
    //   .then(moduleElement => {
    //     // console.log(moduleElement)
    //     if (!validateExportLifecycle(moduleElement)) {
    //       console.error(' 需要导出插件生命周期函数')
    //       return
    //     }
    //     // this.instantiateModule(moduleElement)
    //     moduleElement.mount({ containerRef: this.modulePlaceHolder, baseRouter: `/${pluginName}` })

    this.userApp = await loadRemoteModule(generateRemoteModuleTemplate(pluginName, pluginExpose, pluginPath))
    const platformRef = this.platformProvider.getPlatformRef()
    const ngZone = this.platformProvider.getNgZone()
    if (this.moduleRef) {
      this.moduleRef.destroy()
    }

    // 这里的 'UserAppModule' 是 User 项目中 Angular 模块的名称
    this.moduleRef = await this.userApp.mount(platformRef, ngZone, this.el.nativeElement, this.injector, this.compiler, 'UserAppModule')

    // this.instantiateModule(remoteModule)
    // this.loadModule(remoteModule)
  }

  instantiateModule (remoteModule:Promise<any>) {
    // console.log(remoteModule)
    remoteModule
      .then(moduleElement => {
        // console.log(moduleElement)
        const module = moduleElement['AppModule']
        const moduleRef:NgModuleRef<any> = createNgModule(module, this.injector)
        moduleElement.mount({ containerRef: this.modulePlaceHolder, module: moduleRef })
        // console.log(this.modulePlaceHolder)
      })
  }

  loadModule (path: any) {
    (path as Promise<NgModuleFactory<any> | Type<any>>)
      .then((res:any) => {
        const elementModuleOrFactory = res['AppModule']
        // console.log(elementModuleOrFactory)
        if (elementModuleOrFactory instanceof NgModuleFactory) {
          // if ViewEngine
          return elementModuleOrFactory
        } else {
          try {
            // if Ivy
            return this.compiler.compileModuleAsync(elementModuleOrFactory)
          } catch (err) {
            throw err
          }
        }
      })
      .then((moduleFactory:any) => {
        // console.log(moduleFactory)
        try {
          const elementModuleRef = moduleFactory.create(this.injector)
          const moduleInstance = elementModuleRef.instance
          // do something with the module...
        } catch (err) {
          throw err
        }
      })
  }
}

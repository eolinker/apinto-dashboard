import { Component, OnInit } from '@angular/core'
import { EoIframeComponent } from '../../component/iframe/iframe.component'
import { HttpClient } from '@angular/common/http'
import { NavigationEnd, Router } from '@angular/router'
import { ApiService } from '../../service/api.service'
import { BaseInfoService } from '../../service/base-info.service'
import { EoNgNavigationService } from '../../service/eo-ng-navigation.service'
import { IframeHttpService } from '../../service/iframe-http.service'
import { ModuleOpenConfigData } from '../../constant/type'

@Component({
  selector: 'eo-ng-remote-plugin',
  templateUrl: '../../component/iframe/iframe.component.html',
  styles: [
    `
    :host{
      display:block;
      height:100%;
      overflow-y:hidden;
    }
    :host ::ng-deep{
      nz-spin.iframe-spin,
      nz-spin.iframe-spin >.ant-spin-container,
      #iframePanel,
      #iframePanel > iframe{
        width:100%;
        height:100%;
        border:none;
      }
    }`
  ]
})
export class RemotePluginComponent extends EoIframeComponent implements OnInit {
  constructor (
    iframeService:IframeHttpService,
    api:ApiService,
    router:Router,
    baseInfo:BaseInfoService,
    navigation:EoNgNavigationService,
   private http:HttpClient) {
    super(iframeService, api, router, baseInfo, navigation)
  }

  override ngOnInit (): void {
    this.moduleName = this.baseInfo.allParamsInfo.moduleName
    this.iframeService.moduleName = this.moduleName
    // 此处监听的是切换module事件，需要判断moduleName是否变化
    this.subscription = this.router.events.subscribe((event) => {
      if (event instanceof NavigationEnd) {
        if (this.moduleName !== this.baseInfo.allParamsInfo.moduleName) {
          this.moduleName = this.baseInfo.allParamsInfo.moduleName
          this.iframeService.moduleName = this.moduleName
          this.iframeService.subscription.unsubscribe()
          this.showIframe()
        }
      }
    })

    this.subscription2 = this.iframeService.repFlashIframe().subscribe((event) => {
      this.showIframe(true, `${event ? `/${event}` : ''}`)
    })
  }

  override createIframe = (id: string, url: string) => {
    this.navigation.reqFlashBreadcrumb([])
    this.start = true // TODO 暂时关闭loading，等apispace对接完成后改回false
    const iframe = document.createElement('iframe')
    iframe.id = id
    iframe.width = '100%'
    iframe.height = '100%'
    iframe.src = url
    iframe.onload = () => {
      this.start = true
    }
    return iframe
  }

  override showIframe = (noChangeUrl?:boolean, innerUrl?:string) => {
    this.api.get(`remote/${this.moduleName}`).subscribe((resp:{code:number, data:{module:ModuleOpenConfigData}, msg:string}) => {
      if (resp.code === 0) {
        const url:string = resp.data.module.url
        const initData:{[k:string]:any} = {}
        for (const init of resp.data.module.initialize) {
          switch (init.type) {
            case 'number':
              initData[init.name] = Number(init.value)
              break
            case 'boolean':
              initData[init.name] = Boolean(init.value)
              break
            default:
              initData[init.name] = init.value
          }
        }
        this.getRemoteUrl(`${noChangeUrl && innerUrl ? url + innerUrl : url}`, resp.data.module.header, resp.data.module.query, initData)
      }
    })
  }

  getRemoteUrl (url:string, header?:any, query?:any, initData?:any) {
    let newUrl = url
    let newQuery:string = ''
    for (const que of query) {
      newQuery = `${que.name}=${que.value}`
    }
    if (newUrl.indexOf('?') !== -1) {
      newUrl = `${newUrl.split('?')[0]}?${newQuery}${newUrl.split('?')[1] || ''}`
    } else {
      newUrl = `${newUrl.split('?')[0]}?${newQuery}`
    }
    setTimeout(() => {
      this.iframe = this.createIframe('iframe', newUrl)
      this.loadIframe(initData)
    }, 0)
  }
}

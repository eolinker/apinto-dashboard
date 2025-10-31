import { Component, OnInit } from '@angular/core'
import { Router, NavigationEnd } from '@angular/router'
import { environment } from 'projects/core/src/environments/environment'
import { Subscription } from 'rxjs'
import { ApiService } from '../../service/api.service'
import { BaseInfoService } from '../../service/base-info.service'
import { EoNgNavigationService } from '../../service/eo-ng-navigation.service'
import { IframeHttpService } from '../../service/iframe-http.service'
import { underline } from '../../service/http-interceptors/error.interceptor'

@Component({
  selector: 'eo-ng-iframe',
  templateUrl: './iframe.component.html',
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
export class EoIframeComponent implements OnInit {
  iframe:any = null
  start:boolean = false
  moduleName:string = ''
  initMessage:object|null = null
  subscription: Subscription = new Subscription()
  subscription2: Subscription = new Subscription()

  constructor (
    public iframeService:IframeHttpService,
    public api:ApiService,
    public router:Router,
    public baseInfo:BaseInfoService,
    public navigation:EoNgNavigationService) {}

  ngOnInit (): void {
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

  ngAfterViewInit () {
    this.showIframe()
    this.moduleName = this.baseInfo.allParamsInfo.moduleName
    this.iframeService.moduleName = this.moduleName
  }

  ngOnDestroy () {
    this.stopIframe()
    this.subscription.unsubscribe()
    this.iframeService.subscription.unsubscribe()
    this.subscription2.unsubscribe()
  }

  proxyHandler:{[k:string]:any} ={
    ...this.iframeService.apinto2PluginApi,
    test: function (params:any) {
      const response = params
      return new Promise(resolve => {
        setTimeout(function () {
          resolve('this is response for call test("' + response + '")')
        }, 1)
      })
    }
  }

  listenMessage =async (event:any) => {
    if (event && event.data.apinto && event.data.type === 'request') {
      this.start = true
      const handler = this.proxyHandler[event.data.path as any]
      if (typeof handler === 'function') {
        const args = event.data.data
        const result = await handler(...args)
        try {
          result.data = underline(result.data)
        } catch {
          console.warn(' 转化接口数据命名法出现问题')
        }
        ;(this.iframe as any).contentWindow.postMessage({
          requestId: event.data.requestId,
          type: 'response',
          data: JSON.parse(JSON.stringify(result)),
          apinto: true,
          magic: 'apinto'
        }, '*')
      } else {
        ;(this.iframe as any).contentWindow.postMessage({
          requestId: event.data.requestId,
          apinto: true,
          type: 'error',
          magic: 'apinto',
          data: 'unknown function for:' + event.data.path
        }, '*')
      }
    }
  }

   createIframe = (id: string, url: string) => {
     this.navigation.reqFlashBreadcrumb([])
     const iframe = document.createElement('iframe')
     iframe.id = id
     iframe.width = '100%'
     iframe.height = '100%'
     if (environment.production) {
       iframe.src = url
     } else { // 调试用
       iframe.src = this.moduleName === 'monitor' ? 'http://localhost:4444' : 'http://localhost:5555'
     }
     iframe.onload = () => {
       this.start = true
     }

     return iframe
   }

    onLoadCallback = () => {
      this.start = true
      ;(this.iframe as any).contentWindow.postMessage({
        apinto: true,
        type: 'initialize',
        data: {
          userId: this.navigation.getUserId(),
          userRoleId: this.navigation.getUserRoleId(),
          userModuleAccess: this.navigation.getUserModuleAccess(this.moduleName)
        }
      }, '*')
      window.addEventListener('message', this.listenMessage)
    }

    // changeUrl=true时，表示传入的url是已经处理好的，不需要再根据router.url拼接。暂时用在面包屑场景
 showIframe = (noChangeUrl?:boolean, innerUrl?:string) => {
   const url:string = `agent/${this.moduleName}`

   if (noChangeUrl) {
     if (environment.production) {
       this.iframe.src = `${url}${innerUrl}`
     } else { // 调试用
       this.iframe.src = `http://localhost:4444${innerUrl}`
     }
     return
   }

   setTimeout(() => {
     this.iframe = this.createIframe('iframe', `${url}${this.router.url.includes('#') ? this.router.url.split('#')[1] : ''}`)
     this.loadIframe()
   })
 }

 loadIframe (initData?:any) {
   const onLoadCallback = () => {
     this.start = true
     ;(this.iframe as any).contentWindow.postMessage({
       apinto: true,
       type: 'initialize',
       data: {
         userId: this.navigation.getUserId(),
         userRoleId: this.navigation.getUserRoleId(),
         userModuleAccess: this.navigation.getUserModuleAccess(this.moduleName),
         ...initData
       }
     }, '*')
     window.addEventListener('message', this.listenMessage)
   }
   if ((this.iframe as any).attachEvent) {
     (this.iframe as any).attachEvent('onload', onLoadCallback)
   } else {
     (this.iframe as any).addEventListener('load', onLoadCallback)
   }

   const panel = document.getElementById('iframePanel')
   while (panel?.hasChildNodes()) {
     panel?.firstChild && panel.removeChild(panel?.firstChild)
   }
   panel?.appendChild(this.iframe)
 }

 // 当组件销毁时需要通知iframe注销
 stopIframe () {
   window.removeEventListener('message', this.listenMessage)
   ;(this.iframe as any).contentWindow?.postMessage({
     type: 'stopConnection',
     apinto: true
   }, '*')
 }

 count:number = 0
 testIframe () {
   this.iframeService.batchPublishResModal('api', 'online', {
     uuids: ['apispace-119762', 'apispace-119763', 'apispace-119767',

       '0e0548d2-2ce3-6a10-8643-2d2d74fd5af6',
       'd322e493-dd68-2e3f-cd17-9e0ed8cc9f47'
     ],
     clusters: ['cluster_84bf7525_2d1c_4bb9_81d7_fdcc56de5ccf', 'remote']
   }, true)
 }

 testIframe1 () {
   this.iframeService.apinto2PluginApi.chooseEnvVar().then((resp) => {
     // console.log(resp)
   })
 }
}

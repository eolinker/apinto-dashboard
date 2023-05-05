import { Component, OnInit } from '@angular/core'
import { IframeHttpService } from '../../service/iframe-http.service'
import { ApiService } from '../../service/api.service'
import { EoNgNavigationService } from '../../service/eo-ng-navigation.service'
import { NavigationEnd, Router } from '@angular/router'
import { BaseInfoService } from '../../service/base-info.service'
import { Subscription } from 'rxjs'

@Component({
  selector: 'eo-ng-iframe-page',
  template: `
  <nz-spin class="iframe-spin" [nzSpinning]="!start">
  <div  #iframe id="iframePanel" style="height: 100%;display: block"></div>
  </nz-spin>
  `,
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
export class IframePageComponent implements OnInit {
  iframe:any = null
   proxyHandler:{[k:string]:any} ={
     ...this.iframeService.apinto2PluginApi,
     test: function (params:any) {
       const response = params
       console.log('-----测试是否调用到父窗口')
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
           result.data = this.api.underline(result.data)
         } catch {
           console.warn('转化接口数据命名法出现问题')
         }
         console.log('=----', result)

         ;(this.iframe as any).contentWindow.postMessage({
           requestId: event.data.requestId,
           magic: 'apinto',
           type: 'response',
           data: JSON.parse(JSON.stringify(result)),
           apinto: true
         }, '*')
       } else {
         ;(this.iframe as any).contentWindow.postMessage({
           requestId: event.data.requestId,
           magic: 'apinto',
           apinto: true,
           type: 'error',
           data: 'unknown function for:' + event.data.path
         }, '*')
       }
     }
   }

   // changeUrl=true时，表示传入的url是已经处理好的，不需要再根据router.url拼接。暂时用在面包屑场景
  showIframe = (id: any, url: any, initData: any, noChangeUrl?:boolean) => {
    const createIframe = (id: string, url: string) => {
      console.log(url)
      const iframe = document.createElement('iframe')
      console.log(iframe)
      iframe.id = id
      iframe.width = '100%'
      iframe.height = '100%'
      iframe.src = 'http://localhost:5555' || url
      iframe.onload = () => {
        this.start = true
      }

      return iframe
    }

    if (noChangeUrl) {
      this.iframe.src = url
      return
    }
    this.iframe = createIframe('iframe', `${url}${this.router.url.includes('#') ? this.router.url.split('#')[1] : ''}`)
    const onLoadCallback = () => {
      console.log('load-iframe')
      this.start = true
      ;(this.iframe as any).contentWindow.postMessage({ apinto: true, type: 'initialize', data: initData }, '*')
      window.addEventListener('message', this.listenMessage)
    }
    console.log(this.iframe)
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

  path:string =''

  start:boolean = false
  iframeSrc:string = ''
  iframeDom:Window|null = null
  moduleName:string = ''
  initMessage:object|null = null
  private subscription: Subscription = new Subscription()
  private subscription2: Subscription = new Subscription()

  constructor (private iframeService:IframeHttpService, private api:ApiService,
    private navigation:EoNgNavigationService, private router:Router,
    private baseInfo:BaseInfoService) {}

  ngOnInit (): void {
    console.log('------------20230504')
    this.moduleName = this.baseInfo.allParamsInfo.moduleName
    console.log(this.router.url)
    // 此处监听的是切换module事件，需要判断moduleName是否变化
    this.subscription = this.router.events.subscribe((event) => {
      if (event instanceof NavigationEnd) {
        console.log(this.router.url)
        if (this.moduleName !== this.baseInfo.allParamsInfo.moduleName) {
          this.moduleName = this.baseInfo.allParamsInfo.moduleName
          this.iframeService.moduleName = this.moduleName
          this.subscription.unsubscribe()
          this.iframeService.subscription.unsubscribe()
          this.showIframe('test', `agent/${this.moduleName}`, {})
        }
        // this.getPath()
      }
    })

    this.subscription2 = this.iframeService.repFlashIframe().subscribe((event) => {
      this.showIframe('test', `agent/${this.moduleName}${event ? `/${event}` : ''}`, {}, true)
    })
    // this.getPath()
  }

  ngAfterViewInit () {
    this.showIframe('test', `agent/${this.moduleName}`, {})
  }

  ngOnDestroy () {
    this.stopIframe()
    this.subscription.unsubscribe()
    this.iframeService.subscription.unsubscribe()
    this.subscription2.unsubscribe()
  }

  getPath () {
    this.initMessage = null
    this.api.get('system/module', { name: this.moduleName }).subscribe((resp:any) => {
      if (resp.code === 0) {
        if (resp.data.module.query) {
          this.path = `${resp.data.module.path}${this.router.url.includes('#') ? '/' + this.router.url.split('#')[1] : ''}?`
        }
        for (const queryParams of resp.data.module.query) {
          this.path = `${this.path}${queryParams.name}=${queryParams.value}&`
        }
        this.getIframeSrc({ headers: resp.data.module.header })
        if (resp.data.module.initialize) {
          this.initMessage = this.getInitMessage(resp.data.module.initialize)
        }
      }
    })
  }

  getInitMessage (initializeData:Array<{key:string, value:any, type:string}>) {
    const res:{[k:string]:any} = {}
    for (const data of initializeData) {
      switch (data.type) {
        case 'string': {
          res[data.key] = res[data.value]
          break
        }
        case 'number': {
          res[data.key] = Number(res[data.value])
          break
        }
        case 'boolean': {
          res[data.key] = Boolean(res[data.value])
          break
        }
        case 'array': {
          res[data.key] = [this.getInitMessage(data.value)]
          break
        }
        case 'object': {
          res[data.key] = { ...this.getInitMessage(data.value) as any }
          break
        }
      }
    }
    return res
  }

  // 打开的iframe可能需要传入header
  getIframeSrc (options?:any) {
    this.start = true
    this.iframeService.openIframe(this.path, options).subscribe((blob) => {
      this.iframeSrc = blob
      if (this.initMessage) {
        window.top?.postMessage({ apinto: true, type: 'init', initialize: this.initMessage }, '*')
      }
    })
  }

  onLoad () {
    if (this.start) {
      this.iframeDom = (<HTMLIFrameElement>document.getElementById('iframe')).contentWindow
      this.iframeDom?.postMessage({ data: 'test', apinto: true }, '*')
      this.initIframeChange()
    }
  }

  initIframeChange () {
    const elemIfram = document.getElementById('iframe')!
    if (window.MutationObserver) {
      // chrome
      const callback = (mutations: any[]) => {
        mutations.forEach((mutation: { oldValue: any; target: { src: any; }; }) => {
          this.iframeSrcChanged(mutation.oldValue, mutation.target.src, mutation.target)
        })
      }
      const observer = new MutationObserver(callback)

      observer.observe(elemIfram, {
        attributes: true,
        attributeOldValue: true
      })
    }
  }

  iframeSrcChanged (oldValue:string, newValue:string, iframeObj:any) {
    console.log('旧地址：' + oldValue)
    console.log('新地址：' + newValue)
    // if (newValue.indexOf('aaaa') > -1) {
    //   console.log('有危险，请马上离开……')
    //   iframeObj.src = oldValue// 钓鱼地址，恢复原url
    // } else {
    //   console.log('安全地址，允许跳转……')
    // }
  }
}

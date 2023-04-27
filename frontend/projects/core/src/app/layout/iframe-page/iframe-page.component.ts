import { Component, ElementRef, HostListener, OnInit, ViewChild } from '@angular/core'
import { IframeHttpService } from '../../service/iframe-http.service'
import { ApiService } from '../../service/api.service'
import { EoNgNavigationService } from '../../service/eo-ng-navigation.service'
import { NavigationEnd, Router } from '@angular/router'
import { BaseInfoService } from '../../service/base-info.service'
import { Subscription, take } from 'rxjs'

@Component({
  selector: 'eo-ng-iframe-page',
  template: `
  <nz-spin class="iframe-spin" [nzSpinning]="!start">
  <div  *ngIf="start" #iframe id="iframePanel" style="height: 100%;display: block"></div>
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

  showIframe = (id: any, url: any, initData: any) => {
    const iframe = createIframe('iframe', url)
    const onLoadCallback = () => {
      console.log(this.iframe)
      ;(iframe as any).contentWindow.postMessage({ apinto: true, type: 'initialize', data: initData }, '*')
      window.addEventListener('message', async (event) => {
        if (event && event.data.apinto && event.data.type === 'request') {
          // const msg = {

          // }
          const handler = this.proxyHandler[event.data.path as any]
          console.log(event, handler)
          if (typeof handler === 'function') {
            const args = event.data.data
            const result = await handler(...args)
            console.log(iframe, result, this.proxyHandler)
            ;(iframe as any).contentWindow.postMessage({
              requestId: event.data.requestId,
              magic: 'apinto',
              type: 'response',
              data: JSON.parse(JSON.stringify(result)),
              apinto: true
            }, '*')
          } else {
            ;(iframe as any).contentWindow.postMessage({
              requestId: event.data.requestId,
              magic: 'apinto',
              apinto: true,
              type: 'error',
              data: 'unknown function for:' + event.data.path
            }, '*')
          }
        }
      })
    }
    if ((iframe as any).attachEvent) {
      (iframe as any).attachEvent('onload', onLoadCallback)
    } else {
      (iframe as any).addEventListener('load', onLoadCallback)
    }
    const panel = document.getElementById('iframePanel')
    while (panel?.hasChildNodes()) {
      panel?.firstChild && panel.removeChild(panel?.firstChild)
    }
    panel?.appendChild(iframe)

    function createIframe (id: string, url: string) {
      const iframe = document.createElement('iframe')
      iframe.id = id
      iframe.width = '100%'
      iframe.height = '100%'
      iframe.src = url

      return iframe
    }
  }

  @ViewChild('iframe') iframe: ElementRef|undefined;
  path:string =''

  @HostListener('window:message', ['$event'])
  onMessage (e:any) {
    if (e.data.apinto) {
      switch (e.data.type) {
        case 'request': {
          this.getRequestFromIframe(e.data)
          break
        }
        case 'breadcrumb': {
          // TODO要对每个面包屑导航做特殊处理 ${modulesName}/...
          for (const breadcrumb of e.data.breadcrumbOption) {
            if (breadcrumb.routerLink) {
              breadcrumb.routerLink = `${this.navigation.iframePrefix}/${this.moduleName}/${breadcrumb.routerLink}`
            }
          }
          this.navigation.reqFlashBreadcrumb(e.data.breadcrumbOption)
          break
        }
        case 'router': {
          let newRouterArr:Array<string> = this.router.url.split('#')
          if (this.router.url.includes('#')) {
            newRouterArr.pop()
          }
          newRouterArr = newRouterArr.join('').split('/')
          newRouterArr[newRouterArr.length - 1] = `${newRouterArr[newRouterArr.length - 1]}#${e.data.url}`
          window.location.href = newRouterArr.join('/')
          break
        }
        case 'menu': {
          this.navigation.reqFlashMenu()
          break
        }
        case 'userInfo': {
          this.getUserRequestFromIframe(e.data)
          break
        }
      }
    }
  }

  start:boolean = true
  iframeSrc:string = ''
  iframeDom:Window|null = null
  moduleName:string = ''
  initMessage:object|null = null
  private subscription: Subscription = new Subscription()

  constructor (private iframeService:IframeHttpService, private api:ApiService,
    private navigation:EoNgNavigationService, private router:Router,
    private baseInfo:BaseInfoService) {}

  getRequestFromIframe ({ func, body, params, url, callback, callbackPrefixParam, callbackSuffixParam, others }:any) {
    switch (func) {
      case 'apinto.get': {
        this.api.get(url, params).subscribe((resp:any) => {
          resp.data = this.api.underline(resp.data)
          this.iframeDom?.postMessage({
            callback: callback,
            callbackPrefixParam: callbackPrefixParam,
            callbackSuffixParam: callbackSuffixParam,
            others: others,
            type: 'response',
            response: resp,
            apinto: true
          }, '*')
        })
        break
      }
      case 'apinto.post': {
        this.api.post(url, body, params).subscribe((resp:any) => {
          resp.data = this.api.underline(resp.data)
          this.iframeDom?.postMessage({
            callback: callback,
            callbackPrefixParam: callbackPrefixParam,
            callbackSuffixParam: callbackSuffixParam,
            others: others,
            type: 'response',
            response: resp,
            apinto: true
          }, '*')
        })
        break
      }
      case 'apinto.put': {
        this.api.put(url, body, params).subscribe((resp:any) => {
          resp.data = this.api.underline(resp.data)
          this.iframeDom?.postMessage({
            callback: callback,
            callbackPrefixParam: callbackPrefixParam,
            callbackSuffixParam: callbackSuffixParam,
            others: others,
            type: 'response',
            response: resp,
            apinto: true
          }, '*')
        })
        break
      }
      case 'apinto.delete': {
        this.api.delete(url, params).subscribe((resp:any) => {
          resp.data = this.api.underline(resp.data)
          this.iframeDom?.postMessage({
            callback: callback,
            callbackPrefixParam: callbackPrefixParam,
            callbackSuffixParam: callbackSuffixParam,
            others: others,
            type: 'response',
            response: resp,
            apinto: true
          }, '*')
        })
        break
      }
      case 'apinto.patch': {
        this.api.patch(url, body, params).subscribe((resp:any) => {
          resp.data = this.api.underline(resp.data)
          this.iframeDom?.postMessage({
            callback: callback,
            callbackPrefixParam: callbackPrefixParam,
            callbackSuffixParam: callbackSuffixParam,
            others: others,
            type: 'response',
            response: resp,
            apinto: true
          }, '*')
        })
      }
    }
  }

  getUserRequestFromIframe ({ func }:any) {
    switch (func) {
      case 'apinto.renewUserInfo': {
        // 控制台重新获取用户信息，以便头像中使用，等待接口
        this.navigation.reqUpdateRightList()
        this.navigation.repUpdateRightList().pipe(take(1)).subscribe(() => {
          this.iframeDom?.postMessage({
            type: 'userInfo',
            func: 'apinto.renewUserInfo',
            userId: this.navigation.getUserId(),
            userRoleId: this.navigation.getUserRoleId(),
            userAccess: this.navigation.getRightsList(),
            apinto: true
          }, '*')
        })
        break
      }
      case 'apinto.getCurrentUserInfo': {
        this.iframeDom?.postMessage({
          type: 'userInfo',
          func: 'apinto.getCurrentUserInfo',
          userId: this.navigation.getUserId(),
          userRoleId: this.navigation.getUserRoleId(),
          userModuleAccess: this.navigation.originAccessData[this.moduleName],
          apinto: true
        }, '*')
        break
      }
    }
  }

  ngOnInit (): void {
    this.moduleName = this.baseInfo.allParamsInfo.moduleName
    this.subscription = this.router.events.subscribe((event) => {
      if (event instanceof NavigationEnd) {
        console.log(this.router.url)
        this.moduleName = this.baseInfo.allParamsInfo.moduleName
        // this.getPath()
      }
    })
    // this.getPath()
  }

  ngAfterViewInit () {
    window.onload = () => {
      this.showIframe('test', `agent/${this.moduleName}`, {})
    }
  }

  ngOnDestroy () {
    this.subscription.unsubscribe()
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

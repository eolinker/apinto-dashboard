import { Component, ElementRef, HostListener, Input, OnInit, ViewChild } from '@angular/core'
import { IframeHttpService } from '../../service/iframe-http.service'
import { ApiService } from '../../service/api.service'
import { EoNgNavigationService } from '../../service/eo-ng-navigation.service'
import { Router } from '@angular/router'
import { BaseInfoService } from '../../service/base-info.service'

@Component({
  selector: 'eo-ng-iframe-page',
  template: `
  <nz-spin class="iframe-spin" [nzSpinning]="!start">
      <iframe *ngIf="start "[src] ="path | safe:'resourceUrl'" #iframe id="iframe" (load)="onLoad()" scrolling="yes" >
      <p>Your browser does not support iframe.</p>
      </iframe>
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
    iframe{
      width:100%;
      height:100%;
      border:none;
    }
    }`
  ]
})
export class IframePageComponent implements OnInit {
  @ViewChild('iframe') iframe: ElementRef|undefined;
  @Input() path:string =''

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
        }
      }
    }
  }

  start:boolean = false
  iframeSrc:string = ''
  iframeDom:Window|null = null
  moduleName:string = ''
  constructor (private iframeService:IframeHttpService, private api:ApiService,
    private navigation:EoNgNavigationService, private router:Router,
    private baseInfo:BaseInfoService) {}

  getRequestFromIframe ({ func, body, params, url, callback, callbackPrefixParam, callbackSuffixParam, others }:any) {
    switch (func) {
      case 'apinto.get': {
        this.api.get(url, params).subscribe((resp:any) => {
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

  ngOnInit (): void {
    this.moduleName = this.baseInfo.allParamsInfo.moduleName
    if (!this.path && this.router.url.includes('#')) {
      this.path = this.router.url.split('#')[1]
      this.getIframeSrc()
    }
    if (!this.path && !this.router.url.includes('#')) {
      this.getPath()
    }
    if (this.path) {
      this.start = true
    }
  }

  getPath () {
    this.api.get('system/module', { name: this.moduleName }).subscribe((resp:any) => {
      if (resp.code === 0) {
        if (resp.data.module.query) {
          this.path = `${resp.data.module.path}?`
        }
        for (const queryParams of resp.data.module.query) {
          this.path = `${this.path}${queryParams.name}=${queryParams.value}&`
        }
        this.getIframeSrc({ headers: resp.data.module.header, initialize: resp.data.module.initialize })
      }
    })
  }

  testIframe () {
    return 'test2'
  }

  // 打开的iframe可能需要传入header
  getIframeSrc (options?:any) {
    this.start = true
    this.iframeService.openIframe(this.path, options).subscribe((blob) => {
      this.iframeSrc = blob
    })
  }

  // ngAfterViewInit () {
  //   const doc = this.iframe!.nativeElement.contentDocument || this.iframe!.nativeElement.contentWindow
  // }

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
    if (newValue.indexOf('aaaa') > -1) {
      console.log('有危险，请马上离开……')
      iframeObj.src = oldValue// 钓鱼地址，恢复原url
    } else {
      console.log('安全地址，允许跳转……')
    }
  }
}

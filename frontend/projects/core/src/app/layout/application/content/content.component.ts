/* eslint-disable dot-notation */
/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2022-08-17 23:42:52
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2022-08-24 00:45:28
 * @FilePath: /apinto/src/app/layout/application/application-content/application-content.component.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { ChangeDetectorRef, Component, ElementRef, OnInit, Renderer2, TemplateRef, ViewChild } from '@angular/core'
import { Router } from '@angular/router'
import { TabTemplateContext } from 'ng-zorro-antd/tabs'
import { TabsOptions } from 'eo-ng-tabs'
import { BaseInfoService } from '../../../service/base-info.service'
import { EoNgApplicationService } from '../application.service'

@Component({
  selector: 'eo-ng-application-content',
  templateUrl: './content.component.html',
  styles: [`

  :host{
      height:100%;
      display:block;

      ::ng-deep{
        nz-spin>.ant-spin-container{
          height:100%;
          display:block;
        }
      }
    }`
  ]
})
export class ApplicationContentComponent implements OnInit {
  @ViewChild('tab1Tpl', { read: TemplateRef, static: true }) tab1Tpl: TemplateRef<TabTemplateContext> | string = '额外参数'
  @ViewChild('tab2Tpl', { read: TemplateRef, static: true }) tab2Tpl: TemplateRef<TabTemplateContext> | string = '应用设置'
  @ViewChild('tab3Tpl', { read: TemplateRef, static: true }) tab3Tpl: TemplateRef<TabTemplateContext> | string = '访问鉴权'
  appId:string=''
  options:Array<any>=[]

  tabOptions:TabsOptions[]=[]
  showTopBlank:boolean = false // 是否显示表单上方空隙

  constructor (
    private baseInfo:BaseInfoService,
    private router:Router,
    private cdRef: ChangeDetectorRef,
    private elem: ElementRef,
    private renderer: Renderer2,
    public service:EoNgApplicationService) {
  }

  ngOnInit (): void {
    this.appId = this.baseInfo.allParamsInfo.appId
    if (!this.appId) {
      this.router.navigate(['/', 'application'])
    }
    if (!this.router.url.includes('/message')) {
      this.showTopBlank = false
      this.service.getApplicationData(this.appId)
    } else {
      this.showTopBlank = true
    }
  }

  ngAfterViewInit () {
    this.tabOptions = [
      {
        title: this.tab3Tpl,
        routerLink: 'authentication',
        queryParamsHandling: 'merge'
      },
      {
        title: this.tab1Tpl,
        routerLink: 'extra',
        queryParamsHandling: 'merge'
      },
      {
        title: this.tab2Tpl,
        routerLink: 'message',
        queryParamsHandling: 'merge'
      }
    ]
    if (this.appId === 'anonymous') {
      this.tabOptions.shift()
    }

    this.cdRef.detectChanges()
  }

  ngAfterViewChecked () {
    const element = this.elem.nativeElement.querySelector('[nz-tabs-ink-bar]')
    this.renderer.removeAttribute(element, 'hidden')
  }

  ngOnDestroy () {
    this.service.clearData()
  }

  backToList () {
    this.router.navigate(['/', 'application'])
  }
}

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
import { Router, ActivatedRoute } from '@angular/router'
import { TabTemplateContext } from 'ng-zorro-antd/tabs'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { TabsOptions } from 'eo-ng-tabs'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { BaseInfoService } from '../../../service/base-info.service'

@Component({
  selector: 'eo-ng-application-content',
  templateUrl: './content.component.html',
  styles: [
  ]
})
export class ApplicationContentComponent implements OnInit {
  @ViewChild('tab1Tpl', { read: TemplateRef, static: true }) tab1Tpl: TemplateRef<TabTemplateContext> | string = '上线管理'
  @ViewChild('tab2Tpl', { read: TemplateRef, static: true }) tab2Tpl: TemplateRef<TabTemplateContext> | string = '应用信息'
  @ViewChild('tab3Tpl', { read: TemplateRef, static: true }) tab3Tpl: TemplateRef<TabTemplateContext> | string = '鉴权管理'
  appId:string=''
  options:Array<any>=[]

  tabOptions:TabsOptions[]=[]
  selectedIndex:number = 0
  readonly nowUrl:string = this.router.routerState.snapshot.url

  constructor (private message: EoNgFeedbackMessageService,
    private baseInfo:BaseInfoService,
    private api:ApiService, private router:Router, private activateInfo:ActivatedRoute, private cdRef: ChangeDetectorRef, private elem: ElementRef, private renderer: Renderer2) {
  }

  ngOnInit (): void {
    this.appId = this.baseInfo.allParamsInfo.appId
    if (!this.appId) {
      this.router.navigate(['/', 'application'])
    }
  }

  ngAfterViewInit () {
    this.tabOptions = [
      {
        title: this.tab1Tpl,
        routerLink: '.',
        queryParamsHandling: 'merge'
      },
      {
        title: this.tab2Tpl,
        routerLink: 'message',
        queryParamsHandling: 'merge'
      },
      {
        title: this.tab3Tpl,
        routerLink: 'authentication',
        queryParamsHandling: 'merge'
      }
    ]
    if (this.appId === 'anonymous') {
      this.tabOptions.pop()
    }

    this.cdRef.detectChanges()
  }

  ngAfterViewChecked () {
    const element = this.elem.nativeElement.querySelector('[nz-tabs-ink-bar]')
    this.renderer.removeAttribute(element, 'hidden')
  }
}

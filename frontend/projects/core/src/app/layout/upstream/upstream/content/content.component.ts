/*
 * @Author:
 * @Date: 2022-08-17 23:42:52
 * @LastEditors:
 * @LastEditTime: 2022-08-23 17:47:24
 * @FilePath: /apinto/src/app/layout/upstream/upstream/upstream-content/upstream-content.component.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
/* eslint-disable dot-notation */
import { ChangeDetectorRef, Component, ElementRef, OnInit, Renderer2, TemplateRef, ViewChild } from '@angular/core'
import { Router, ActivatedRoute } from '@angular/router'
import { TabTemplateContext } from 'ng-zorro-antd/tabs'
import { TabsOptions } from 'eo-ng-tabs'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'

@Component({
  selector: 'eo-ng-upstream-content',
  templateUrl: './content.component.html',
  styles: [
  ]
})
export class UpstreamContentComponent implements OnInit {
  @ViewChild('tab1Tpl', { read: TemplateRef, static: true }) tab1Tpl: TemplateRef<TabTemplateContext> | string = '环境变量'
  @ViewChild('tab2Tpl', { read: TemplateRef, static: true }) tab2Tpl: TemplateRef<TabTemplateContext> | string = '证书管理'
  @ViewChild('tab3Tpl', { read: TemplateRef, static: true }) tab3Tpl: TemplateRef<TabTemplateContext> | string = '网关节点'
  serviceName:string=''
  options:Array<any>=[]

  tabOptions:TabsOptions[]=[]
  selectedIndex:number = 0
  readonly nowUrl:string = this.router.routerState.snapshot.url

  constructor (private router:Router, private baseInfo:BaseInfoService, private activateInfo:ActivatedRoute, private cdRef: ChangeDetectorRef, private elem: ElementRef, private renderer: Renderer2) {
  }

  ngOnInit (): void {
    this.serviceName = this.baseInfo.allParamsInfo.serviceName
    if (!this.serviceName) {
      this.router.navigate(['/'])
    }
  }

  ngAfterViewInit () {
    this.tabOptions = [
      // {
      //   title: this.tab1Tpl,
      //   routerLink: '.',
      //   queryParams: {
      //     serviceName: this.serviceName
      //   },
      //   queryParamsHandling: 'merge'
      // },
      {
        title: this.tab2Tpl,
        routerLink: 'publish',
        queryParamsHandling: 'merge'
      },
      {
        title: this.tab3Tpl,
        routerLink: 'message',
        queryParamsHandling: 'merge'
      }
    ]
    this.cdRef.detectChanges()
  }

  ngAfterViewChecked () {
    const element = this.elem.nativeElement.querySelector('[nz-tabs-ink-bar]')
    this.renderer.removeAttribute(element, 'hidden')
  }
}

/* eslint-disable dot-notation */
/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2022-08-17 23:42:52
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2022-09-20 23:16:05
 * @FilePath: /apinto/src/app/layout/api/api-content/api-content.component.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { ChangeDetectorRef, Component, ElementRef, OnInit, Renderer2, TemplateRef, ViewChild } from '@angular/core'
import { Router } from '@angular/router'
import { TabTemplateContext } from 'ng-zorro-antd/tabs'
import { TabsOptions } from 'eo-ng-tabs'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'

@Component({
  selector: 'eo-ng-api-content',
  templateUrl: './content.component.html',
  styles: [
  ]
})
export class ApiContentComponent implements OnInit {
  @ViewChild('tab1Tpl', { read: TemplateRef, static: true }) tab1: TemplateRef<TabTemplateContext> | string = '环境变量'
  @ViewChild('tab2Tpl', { read: TemplateRef, static: true }) tab2: TemplateRef<TabTemplateContext> | string = '证书管理'
  apiUuid:string=''
  options:Array<any>=[]

  tabOptions:TabsOptions[]=[]
  selectedIndex:number = 0
  readonly nowUrl:string = this.router.routerState.snapshot.url

  constructor (
     private baseInfo:BaseInfoService,
     private router:Router,
     private cdRef: ChangeDetectorRef,
      private elem: ElementRef,
      private renderer: Renderer2) {
  }

  ngOnInit (): void {
    this.apiUuid = this.baseInfo.allParamsInfo.apiId
    if (!this.apiUuid) {
      this.router.navigate(['/', 'router', 'api', 'group'])
    }
  }

  ngAfterViewInit () {
    // this.tabOptions = [
    //   {
    //     title: this.tab1,
    //     routerLink: '.',
    //     queryParamsHandling: 'merge'
    //   },
    //   {
    //     title: this.tab2,
    //     routerLink: 'publish',
    //     queryParamsHandling: 'merge'
    //   }
    // ]
    // this.cdRef.detectChanges()
  }

  ngAfterViewChecked () {
    // const element = this.elem.nativeElement.querySelector('[nz-tabs-ink-bar]')
    // this.renderer.removeAttribute(element, 'hidden')
  }

  ngOnDestroy () {
  }
}

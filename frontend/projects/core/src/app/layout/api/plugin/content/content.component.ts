import { ChangeDetectorRef, Component, ElementRef, OnInit, Renderer2, TemplateRef, ViewChild } from '@angular/core'
import { Router } from '@angular/router'
import { TabTemplateContext, TabsOptions } from 'eo-ng-tabs'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'

@Component({
  selector: 'eo-ng-api-plugin-template-content',
  template: `
  <div class="content-tabs-fixed">
    <eo-ng-tabs
      [nzLinkRouter]="true"
      [(nzSelectedIndex)]="selectedIndex"
      [nzOptions]="tabOptions"
    ></eo-ng-tabs>
    <ng-template #tab1Tpl class="test">
      <span>上线管理</span>
    </ng-template>

    <ng-template #tab2Tpl>
      <span>模板信息</span>
    </ng-template>
  </div>
  <div class="inside-content-fixed-top-51">
    <router-outlet></router-outlet>
  </div>
  `,
  styles: [
  ]
})
export class ApiPluginTemplateContentComponent implements OnInit {
  @ViewChild('tab1Tpl', { read: TemplateRef, static: true }) tab1: TemplateRef<TabTemplateContext> | string = '上线管理'
  @ViewChild('tab2Tpl', { read: TemplateRef, static: true }) tab2: TemplateRef<TabTemplateContext> | string = '模板信息'
  uuid:string=''
  options:Array<any>=[]
  tabOptions:TabsOptions[]=[]
  selectedIndex:number = 0

  constructor (
     private baseInfo:BaseInfoService,
     private router:Router,
     private cdRef: ChangeDetectorRef,
      private elem: ElementRef,
      private renderer: Renderer2) {
  }

  ngOnInit (): void {
    this.uuid = this.baseInfo.allParamsInfo.pluginTemplateId
    if (!this.uuid) {
      this.router.navigate(['/', 'router', 'plugin'])
    }
  }

  ngAfterViewInit () {
    this.tabOptions = [
      {
        title: this.tab1,
        routerLink: '.',
        queryParamsHandling: 'merge'
      },
      {
        title: this.tab2,
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

/* eslint-disable dot-notation */
import { Component, Input, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { ActivatedRoute, NavigationEnd, Router } from '@angular/router'
import { TabsOptions } from 'eo-ng-tabs'
import { Subscription } from 'rxjs'
import { AppConfigService } from '../../../service/app-config.service'
import { BaseInfoService } from '../../../service/base-info.service'
import { EoNgMonitorTabsService } from '../../../service/eo-ng-monitor-tabs.service'

@Component({
  selector: 'eo-ng-monitor-alarm-area',
  template: `
    <ng-container *ngIf="!initFlag">
      <eo-ng-tabs
              [ngStyle]="{'margin-top':showTabs? '0px' : '4px'}"
              class="eo-ng-monitor-area-inside-tabs"
              [(nzSelectedIndex)]="selectedIndex"
              [nzTabBarExtraContent]="extraTpl"
              [nzOptions]="tabOptions"
              [ngClass]="{ 'eo-ng-monitor-tabs-hidden': showTabs }"
      ></eo-ng-tabs>
      <ng-template #extraTpl>
      <button eoNgFeedbackTooltip  style="margin-right:var(--LAYOUT_MARGIN)" class="ant-btn ant-btn-text" nzTooltipPlacement="bottom" nzTooltipTrigger='hover' nzTooltipTitle="编辑分区配置" >
        <svg class="iconpark-icon" (click)="goToConfig()"><use href="#setting"></use></svg>
      </button>
      </ng-template>
      <ng-template #totalTpl>
        <eo-ng-monitor-alarm-area-total></eo-ng-monitor-alarm-area-total>
      </ng-template>
      <ng-template #apiTpl>
        <eo-ng-monitor-alarm-area-api></eo-ng-monitor-alarm-area-api>
      </ng-template>
      <ng-template #appTpl>
        <eo-ng-monitor-alarm-area-app></eo-ng-monitor-alarm-area-app>
      </ng-template>
      <ng-template #upstreamTpl>
        <eo-ng-monitor-alarm-area-service></eo-ng-monitor-alarm-area-service>
      </ng-template>
      <ng-template #strategyTpl>
        <eo-ng-monitor-alarm-strategy></eo-ng-monitor-alarm-strategy>
      </ng-template>
      <ng-template #historyTpl>
        <eo-ng-monitor-alarm-history></eo-ng-monitor-alarm-history>
      </ng-template>
    </ng-container>
    <div class="" *ngIf="initFlag">
      <router-outlet></router-outlet>
    </div>
  `,
  styles: [
  ]
})
export class MonitorAlarmAreaComponent implements OnInit {
  @ViewChild('totalTpl') totalTpl: TemplateRef<any> | undefined;
  @ViewChild('apiTpl') apiTpl: TemplateRef<any> | undefined;
  @ViewChild('appTpl') appTpl: TemplateRef<any> | undefined;
  @ViewChild('upstreamTpl') upstreamTpl: TemplateRef<any> | undefined;
  @ViewChild('strategyTpl') strategyTpl: TemplateRef<any> | undefined;
  @ViewChild('historyTpl') historyTpl: TemplateRef<any> | undefined;
  @Input() initFlag:boolean = false // 是否显示配置页
  partitionId:string = ''
  tabOptions:TabsOptions[]= []

  tabIndex:string[]= ['total', 'api', 'app', 'service', 'strategy', 'history']
  showTabs:boolean = false

  selectedIndex:number = 0
  private subscription: Subscription = new Subscription()
  private subscription2: Subscription = new Subscription()
  constructor (private activateInfo:ActivatedRoute,
    private appConfigService:AppConfigService,
    private baseInfo:BaseInfoService,
    private tabs:EoNgMonitorTabsService,
    private router:Router) {
  }

  ngOnInit (): void {
    this.partitionId = this.baseInfo.allParamsInfo.partitionId
    this.tabOptions = [
      {
        title: '监控总览',
        routerLink: 'total/' + this.partitionId,
        queryParamsHandling: '',
        lazyLoad: true
      },
      {
        title: 'API调用统计',
        routerLink: 'api/' + this.partitionId,
        queryParamsHandling: '',
        lazyLoad: true
      },
      {
        title: '应用调用统计',
        routerLink: 'app/' + this.partitionId,
        queryParamsHandling: '',
        lazyLoad: true
      },
      {
        title: '上游调用统计',
        routerLink: 'service/' + this.partitionId,
        queryParamsHandling: '',
        lazyLoad: true
      },
      {
        title: '告警策略',
        routerLink: 'strategy/' + this.partitionId,
        queryParamsHandling: '',
        lazyLoad: true
      },
      {
        title: '告警历史',
        routerLink: 'history/' + this.partitionId,
        queryParamsHandling: '',
        lazyLoad: true
      }
    ]
    this.initTabs()
    this.subscription = this.router.events.subscribe((event) => {
      if (event instanceof NavigationEnd) {
        this.initTabs()
      }
    })
  }

  ngDoCheck () {
  }

  ngAfterViewInit () {
    this.renderTemplate()
  }

  renderTemplate () {
    if (this.partitionId) {
      this.tabOptions[0].content = this.totalTpl
      this.tabOptions[1].content = this.apiTpl
      this.tabOptions[2].content = this.appTpl
      this.tabOptions[3].content = this.upstreamTpl
      this.tabOptions[4].content = this.strategyTpl
      this.tabOptions[5].content = this.historyTpl
    }
  }

  initTabs () {
    this.partitionId = this.baseInfo.allParamsInfo.partitionId
    this.tabOptions[0].routerLink = 'total/' + this.partitionId
    this.tabOptions[1].routerLink = 'api/' + this.partitionId
    this.tabOptions[2].routerLink = 'app/' + this.partitionId
    this.tabOptions[3].routerLink = 'service/' + this.partitionId
    this.tabOptions[4].routerLink = 'strategy/' + this.partitionId
    this.tabOptions[5].routerLink = 'history/' + this.partitionId

    this.initFlag = this.router.url.includes('area/message') || this.router.url.includes('area/config')
    if (this.initFlag) {
      this.appConfigService.reqFlashBreadcrumb([{ title: '监控告警' }])
    }
    this.selectedIndex = this.tabIndex.indexOf(this.router.url.split('?')[0].split('/')[3] || 'total')
    if (this.router.url.split('?')[0].split('/').length === 5 && !this.initFlag) {
      this.appConfigService.reqFlashBreadcrumb([{ title: '监控告警', routerLink: 'monitor-alarm' }, { title: this.tabs.getTabName(this.partitionId), routerLink: 'monitor-alarm/area/total/' + this.partitionId }, { title: this.tabOptions[this.selectedIndex].title }])
    }
    setTimeout(() => {
      if (!this.tabOptions[this.selectedIndex]?.content && this.totalTpl) {
        this.renderTemplate()
      }
    }, 0)

    if (this.router.url.includes('/detail') || (this.router.url.includes('strategy') && this.router.url.includes('config')) || (this.router.url.includes('strategy') && this.router.url.includes('/message'))) {
      this.showTabs = true
    } else {
      this.showTabs = false
    }
  }

  ngOnDestroy () {
    this.subscription.unsubscribe()
    this.subscription2.unsubscribe()
  }

  goToConfig () {
    this.initFlag = true
    this.router.navigate(['/', 'monitor-alarm', 'area', 'message', this.partitionId])
  }
}

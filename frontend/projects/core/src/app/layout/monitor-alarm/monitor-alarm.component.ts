/* eslint-disable dot-notation */
import { Component, OnInit, ViewChildren } from '@angular/core'
import { NavigationEnd, Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { Subscription } from 'rxjs/internal/Subscription'
import { TabHostDirective } from '../../directive/tab-host.directive'
import { ApiService } from '../../service/api.service'
import { AppConfigService } from '../../service/app-config.service'
import { BaseInfoService } from '../../service/base-info.service'
import { EoNgMonitorTabsService } from '../../service/eo-ng-monitor-tabs.service'
import { MonitorPartition } from './types/types'

@Component({
  selector: 'eo-ng-monitor-alarm',
  templateUrl: './monitor-alarm.component.html',
  styles: [
    `
    `
  ]
})
export class MonitorAlarmComponent implements OnInit {
  @ViewChildren(TabHostDirective) tabHost?: TabHostDirective
  private subscription: Subscription = new Subscription()
  partitionId:string = ''
  isLoading:boolean = true
  createConfig:boolean = true
  showTabs:boolean = false
  // eslint-disable-next-line no-useless-constructor
  constructor (private api:ApiService, private message: EoNgFeedbackMessageService, private appConfigService:AppConfigService, public tabs:EoNgMonitorTabsService,
    private baseInfo:BaseInfoService, private router:Router) {
  }

  ngOnInit (): void {
    this.getTabs(true) // 查询分区列表
    this.partitionId = this.baseInfo.allParamsInfo.partitionId
    this.appConfigService.reqFlashBreadcrumb([{ title: '监控告警' }])
    this.subscription = this.router.events.subscribe((event) => {
      if (event instanceof NavigationEnd) {
        this.partitionId = this.baseInfo.allParamsInfo.partitionId
        if (this.partitionId) { // 路由变化时，如果有分区id，需要将tab切换至相应分区，同时记录该分区当前url
          this.tabs.changeTab(this.partitionId)
          const rawUrl = this.router.url.split('/')
          rawUrl.shift()
          const newUrl = rawUrl.join('/').split('?')[0]
          this.tabs.changeRouter(this.partitionId, '/' + newUrl)
        } else if (!this.router.url.includes('area/config') && this.router.url.includes('monitor-alarm')) {
          this.getTabs(true)
        }
      }
    })
  }

  // 当进入分区配置页（area/config）、分区调用统计详情页（detail）及告警策略配置页（strategy/config、strategy/message)时
  ngDoCheck () {
    if (this.router.url.includes('area/config')) {
      this.createConfig = true
    } else if (this.router.url.includes('area')) {
      this.createConfig = false
    }
    if (this.router.url.includes('/detail') || (this.router.url.includes('strategy') && this.router.url.includes('config')) || (this.router.url.includes('strategy') && this.router.url.includes('/message'))) {
      this.showTabs = true
    } else {
      this.showTabs = false
    }
  }

  ngOnDestroy () {
    this.subscription.unsubscribe()
  }

  // 从接口获取分区列表
  getTabs (changeRouter?:boolean) {
    this.api.get('monitor/partitions').subscribe((resp:{code:number, data:{partitions:MonitorPartition[]}, msg:string}) => {
      if (resp.code === 0) {
        if (resp.data.partitions.length > 0) {
          this.tabs.getTabs(resp.data.partitions, this.partitionId)
          this.isLoading = false
          if (changeRouter && !this.partitionId) {
            this.router.navigate(['/', 'monitor-alarm', 'area', 'total', resp.data.partitions[0].uuid])
          }
        }
      } else {
        this.message.error(resp.msg || '获取分区列表失败，请重试！')
      }
    })
  }

  // 添加分区
  addArea () {
    this.tabs.showList = []
    this.tabs.newTab()
    this.tabs.addTabFlag = true
    this.router.navigate(['/', 'monitor-alarm', 'area', 'config'])
  }

  changeTab () {
    this.partitionId = this.tabs.list[this.tabs.index].uuid || this.partitionId
  }
}

/* eslint-disable camelcase */
/* eslint-disable dot-notation */
/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2022-07-20 22:34:58
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2022-08-08 00:28:56
 * @FilePath: /apinto/src/app/layout/deploy/deploy-cluster-content/deploy-cluster-content.component.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { Component, OnInit, TemplateRef, ViewChild, ChangeDetectorRef } from '@angular/core'
import { Router } from '@angular/router'
import { TabsOptions, TabTemplateContext } from 'eo-ng-tabs'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { DeployService } from '../../deploy.service'

@Component({
  selector: 'eo-ng-deploy-cluster-content',
  templateUrl: './content.component.html',
  styleUrls: [
    './content.component.scss'
  ]
})
export class DeployClusterContentComponent implements OnInit {
  @ViewChild('tab1', { read: TemplateRef, static: true }) tab1: TemplateRef<TabTemplateContext> | string = '全局变量'
  @ViewChild('tab2', { read: TemplateRef, static: true }) tab2: TemplateRef<TabTemplateContext> | string = 'SSL证书'
  @ViewChild('tab3', { read: TemplateRef, static: true }) tab3: TemplateRef<TabTemplateContext> | string = '节点列表'
  @ViewChild('tab5', { read: TemplateRef, static: true }) tab5: TemplateRef<TabTemplateContext> | string = '节点插件'
  @ViewChild('tab6', { read: TemplateRef, static: true }) tab6: TemplateRef<TabTemplateContext> | string = '集群管理'
  clusterName:string=''
  clusterDesc:string=''
  _clusterDesc:string=''
  options:Array<any>=[]
  disabled:boolean = true

  tabOptions:TabsOptions[]=[]
  readonly nowUrl:string = this.router.routerState.snapshot.url
  constructor (
    public service:DeployService,
    private baseInfo:BaseInfoService,
     private api:ApiService, private router:Router, private cdRef: ChangeDetectorRef) {
  }

  ngOnInit (): void {
    this.clusterName = this.baseInfo.allParamsInfo.clusterName
    if (!this.clusterName) {
      this.router.navigate(['/'])
    }
    if (!this.router.url.includes('/message')) {
      this.getClustersData()
    }
  }

  ngAfterViewInit () {
    this.tabOptions = [
      {
        title: this.tab3,
        routerLink: '.',
        queryParamsHandling: 'merge'
      },
      {
        title: this.tab5,
        routerLink: 'plugin',
        queryParamsHandling: 'merge'
      },
      {
        title: this.tab1,
        routerLink: 'env',
        queryParamsHandling: 'merge'
      },
      {
        title: this.tab2,
        routerLink: 'cert',
        queryParamsHandling: 'merge'
      },
      {
        title: this.tab6,
        routerLink: 'message',
        queryParamsHandling: 'merge'
      }
    ]
    this.cdRef.detectChanges()
  }

  ngAfterViewChecked () {
    document.getElementsByClassName('ant-tabs-ink-bar')[0]?.removeAttribute('hidden')
  }

  getClustersData () {
    this.api.get('cluster', { clusterName: this.clusterName }).subscribe((resp:{code:number, data:{cluster:{desc:string, title:string, [key:string]:any}}, msg:string}) => {
      if (resp.code === 0) {
        this.service.clusterName = resp.data.cluster.title
        this.service.clusterDesc = resp.data.cluster.desc
        this.clusterDesc = resp.data.cluster.desc
        this._clusterDesc = resp.data.cluster.desc
      }
    })
  }

  save () {
    this.api.put('cluster/' + this.clusterName + '/desc', { desc: this._clusterDesc }).subscribe((resp:{code:number, data:any, msg:string}) => {
      if (resp.code === 0) {
        this.clusterDesc = this._clusterDesc
      } else {
        this._clusterDesc = this.clusterDesc
      }
    })
    this.disabled = true
  }

  backToList () {
    this.router.navigate(['/', 'deploy', 'cluster'])
  }
}

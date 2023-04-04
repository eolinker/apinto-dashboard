/* eslint-disable camelcase */
/* eslint-disable dot-notation */
/*
 * @Author:
 * @Date: 2022-07-20 22:34:58
 * @LastEditors:
 * @LastEditTime: 2022-08-08 00:28:56
 * @FilePath: /apinto/src/app/layout/deploy/deploy-cluster-content/deploy-cluster-content.component.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { Component, OnInit, TemplateRef, ViewChild, ChangeDetectorRef } from '@angular/core'
import { Router, ActivatedRoute } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { TabsOptions, TabTemplateContext } from 'eo-ng-tabs'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'

@Component({
  selector: 'eo-ng-deploy-cluster-content',
  templateUrl: './content.component.html',
  styleUrls: [
    './content.component.scss'
  ]
})
export class DeployClusterContentComponent implements OnInit {
  @ViewChild('descInput', { read: TemplateRef, static: false }) descInput: TemplateRef<any> | undefined
  @ViewChild('tab1', { read: TemplateRef, static: true }) tab1: TemplateRef<TabTemplateContext> | string = '环境变量'
  @ViewChild('tab2', { read: TemplateRef, static: true }) tab2: TemplateRef<TabTemplateContext> | string = '证书管理'
  @ViewChild('tab3', { read: TemplateRef, static: true }) tab3: TemplateRef<TabTemplateContext> | string = '网关节点'
  @ViewChild('tab5', { read: TemplateRef, static: true }) tab5: TemplateRef<TabTemplateContext> | string = '插件管理'
  clusterName:string=''
  clusterDesc:string=''
  _clusterDesc:string=''
  options:Array<any>=[]
  disabled:boolean = true

  tabOptions:TabsOptions[]=[]
  readonly nowUrl:string = this.router.routerState.snapshot.url
  constructor (private baseInfo:BaseInfoService,
     private message: EoNgFeedbackMessageService, private api:ApiService, private router:Router, private activateInfo:ActivatedRoute, private cdRef: ChangeDetectorRef) {
  }

  ngOnInit (): void {
    this.clusterName = this.baseInfo.allParamsInfo.clusterName
    if (!this.clusterName) {
      this.router.navigate(['/'])
    }
    this.getClustersData()
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
        routerLink: 'cert',
        queryParamsHandling: 'merge'
      },
      {
        title: this.tab3,
        routerLink: 'nodes',
        queryParamsHandling: 'merge'
      },
      {
        title: this.tab5,
        routerLink: 'plugin',
        queryParamsHandling: 'merge'
      }
    ]
    this.cdRef.detectChanges()
  }

  ngAfterViewChecked () {
    document.getElementsByClassName('ant-tabs-ink-bar')[0]?.removeAttribute('hidden')
  }

  getClustersData () {
    this.api.get('cluster', { clusterName: this.clusterName }).subscribe((resp:{code:number, data:{cluster:{desc:string, [key:string]:any}}, msg:string}) => {
      if (resp.code === 0) {
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
}

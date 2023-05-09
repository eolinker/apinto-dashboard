/* eslint-disable dot-notation */
import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { EoNgNavigationService } from 'projects/core/src/app/service/eo-ng-navigation.service'
import { CommonPublishTableBody, CommonPublishTableHeadName } from '../../../constant/conf'
import { EmptyHttpResponse, PublishManagementData } from '../../../constant/type'
import { BaseInfoService } from '../../../service/base-info.service'
import { PublishFailService } from '../../../service/publish-fail.service'

@Component({
  selector: 'eo-ng-application-publish',
  templateUrl: './publish.component.html',
  styles: [
  ]
})
export class ApplicationPublishComponent implements OnInit {
  @ViewChild('clusterStatusTpl', { read: TemplateRef, static: true }) clusterStatusTpl: TemplateRef<any> | undefined
  @ViewChild('disabledStatusTpl', { read: TemplateRef, static: true }) disabledStatusTpl: TemplateRef<any> | undefined
  appId:string = ''
  clustersList : Array<object> = []
  clustersTableHeadName:THEAD_TYPE[]= [...CommonPublishTableHeadName]
  clustersTableBody:TBODY_TYPE[] = [...CommonPublishTableBody]
  type:string = ''
  solutionRouter:string = ''
  solutionParam:any = {}
  nzDisabled:boolean = false

  constructor (
                private message: EoNgFeedbackMessageService,
                public api:ApiService,
                private baseInfo:BaseInfoService,
                private router:Router,
                private navigationService:EoNgNavigationService,
                private publishFailModal:PublishFailService) {
    this.navigationService.reqFlashBreadcrumb([{ title: '应用管理', routerLink: 'application' }, { title: '上线管理' }])
  }

  ngOnInit (): void {
    this.appId = this.baseInfo.allParamsInfo.appId
    this.initTable()
    if (!this.appId) {
      this.router.navigate(['/', 'application'])
    }
    this.getClustersData()
  }

  initTable () {
    // 状态:待更新&启用 - 可操作:更新,下线,禁用
    this.clustersTableBody[6].btns[0].disabledFn = () => { return this.nzDisabled }
    this.clustersTableBody[6].btns[0].click = (item:any) => { this.updateOrOnline(item.data, '更新') }
    this.clustersTableBody[6].btns[1].disabledFn = () => { return this.nzDisabled }
    this.clustersTableBody[6].btns[1].click = (item:any) => { this.offline(item.data) }
    this.clustersTableBody[6].btns[2].disabledFn = () => { return this.nzDisabled }
    this.clustersTableBody[6].btns[2].click = (item:any) => { this.disable(item.data) }

    // 状态:待更新&禁用 - 可操作:更新,下线,启用
    this.clustersTableBody[7].btns[0].disabledFn = () => { return this.nzDisabled }
    this.clustersTableBody[7].btns[0].click = (item:any) => { this.updateOrOnline(item.data, '更新') }
    this.clustersTableBody[7].btns[1].disabledFn = () => { return this.nzDisabled }
    this.clustersTableBody[7].btns[1].click = (item:any) => { this.offline(item.data) }
    this.clustersTableBody[7].btns[2].disabledFn = () => { return this.nzDisabled }
    this.clustersTableBody[7].btns[2].click = (item:any) => { this.enable(item.data) }

    // 状态:已上线&启用 - 可操作:下线,禁用
    this.clustersTableBody[8].btns[0].disabledFn = () => { return this.nzDisabled }
    this.clustersTableBody[8].btns[0].click = (item:any) => { this.offline(item.data) }
    this.clustersTableBody[8].btns[1].disabledFn = () => { return this.nzDisabled }
    this.clustersTableBody[8].btns[1].click = (item:any) => { this.disable(item.data) }

    // 状态:已上线&禁用 - 可操作:下线,启用
    this.clustersTableBody[9].btns[0].disabledFn = () => { return this.nzDisabled }
    this.clustersTableBody[9].btns[0].click = (item:any) => { this.offline(item.data) }
    this.clustersTableBody[9].btns[1].disabledFn = () => { return this.nzDisabled }
    this.clustersTableBody[9].btns[1].click = (item:any) => { this.enable(item.data) }

    // 状态:已下线&启用 - 可操作:上线,禁用
    this.clustersTableBody[10].btns[0].disabledFn = () => { return this.nzDisabled }
    this.clustersTableBody[10].btns[0].click = (item:any) => { this.updateOrOnline(item.data, '上线') }
    this.clustersTableBody[10].btns[1].disabledFn = () => { return this.nzDisabled }
    this.clustersTableBody[10].btns[1].click = (item:any) => { this.disable(item.data) }

    // 状态:已下线&禁用 - 可操作:上线,启用
    this.clustersTableBody[11].btns[0].disabledFn = () => { return this.nzDisabled }
    this.clustersTableBody[11].btns[0].click = (item:any) => { this.updateOrOnline(item.data, '上线') }
    this.clustersTableBody[11].btns[1].disabledFn = () => { return this.nzDisabled }
    this.clustersTableBody[11].btns[1].click = (item:any) => { this.enable(item.data) }
  }

  ngAfterViewInit () {
    this.clustersTableBody[2].title = this.clusterStatusTpl
    this.clustersTableBody[3].title = this.disabledStatusTpl
  }

  getClustersData () {
    this.api.get('application/onlines', { appId: this.appId })
      .subscribe((resp:{code:number, data:{clusters:PublishManagementData[]}, msg:string}) => {
        if (resp.code === 0) {
          this.clustersList = resp.data.clusters
        }
      })
  }

  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  updateOrOnline (item:PublishManagementData, type:string) {
    this.api.put('application/online', { clusterName: item.name || '' }, { appId: this.appId })
      .subscribe((resp:{code:number, data :{router:{params:{[k:string]:string}, name:string}}, msg:string}) => {
        if (resp.code === 0) {
          this.message.success(resp.msg || (type + '成功'), { nzDuration: 1000 })
          this.getClustersData()
        } else {
          this.solutionRouter = resp.data?.router?.name ? resp.data.router.name : ''
          this.solutionParam = resp.data?.router?.params ? resp.data.router.params : {}
          if (this.solutionRouter) {
            this.publishFailModal.openModal(resp.msg, '应用', this.solutionRouter, this.solutionParam)
          } else {
            this.message.error(resp.msg || '操作失败')
          }
        }
      })
  }

  offline (item:PublishManagementData) {
    this.api.put('application/offline', { clusterName: (item.name || '') }, { appId: this.appId })
      .subscribe((resp:EmptyHttpResponse) => {
        if (resp.code === 0) {
          this.message.success(resp.msg || '下线成功', { nzDuration: 1000 })
          this.getClustersData()
        }
      })
  }

  enable (item:PublishManagementData) {
    this.api.put('application/enable', { clusterName: (item.name || '') }, { appId: this.appId })
      .subscribe((resp:EmptyHttpResponse) => {
        if (resp.code === 0) {
          this.message.success(resp.msg || '启用成功', { nzDuration: 1000 })
          this.getClustersData()
        }
      })
  }

  disable (item:PublishManagementData) {
    this.api.put('application/disable', { clusterName: item.name || '' }, { appId: this.appId })
      .subscribe((resp:EmptyHttpResponse) => {
        if (resp.code === 0) {
          this.message.success(resp.msg || '停用成功', { nzDuration: 1000 })
          this.getClustersData()
        }
      })
  }
}

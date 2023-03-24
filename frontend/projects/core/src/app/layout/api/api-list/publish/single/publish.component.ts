/* eslint-disable dot-notation */
/*
 * @Author:
 * @Date: 2022-08-17 23:42:52
 * @LastEditors:
 * @LastEditTime: 2022-09-20 23:15:51
 * @FilePath: /apinto/src/app/layout/api/api-publish/api-publish.component.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { CommonPublishTableHeadName, CommonPublishTableBody } from 'projects/core/src/app/constant/conf'
import { EmptyHttpResponse, PublishManagementData } from 'projects/core/src/app/constant/type'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { PublishFailService } from 'projects/core/src/app/service/publish-fail.service'

@Component({
  selector: 'eo-ng-api-publish',
  templateUrl: './publish.component.html',
  styles: [
  ]
})
export class ApiPublishComponent implements OnInit {
  @ViewChild('clusterStatusTpl', { read: TemplateRef, static: true }) clusterStatusTpl: TemplateRef<any> | undefined
  @ViewChild('disabledStatusTpl', { read: TemplateRef, static: true }) disabledStatusTpl: TemplateRef<any> | undefined
  apiUuid:string = ''
  nzDisabled:boolean = false
  clustersList : PublishManagementData[] = []
  clustersTableHeadName: THEAD_TYPE[]= [...CommonPublishTableHeadName]
  clustersTableBody: TBODY_TYPE[] = [...CommonPublishTableBody]
  type:string = ''
  solutionRouter:string = ''
  solutionParam:{[k:string]:any} = {}

  constructor (
    private message: EoNgFeedbackMessageService,
    public api:ApiService,
    private router:Router,
    private publishFailModal:PublishFailService,
    private baseInfo:BaseInfoService,
    private appConfigService:AppConfigService
  ) {
    this.apiUuid = this.baseInfo.allParamsInfo.apiId
    if (!this.apiUuid) {
      this.router.navigate(['/', 'router', 'api', 'group'])
    }
    this.appConfigService.reqFlashBreadcrumb([{ title: 'API管理', routerLink: 'router/api/group/list' }, { title: '上线管理' }])
  }

  ngOnInit (): void {
    this.getClustersData()
    this.initTable()
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
    this.clustersTableBody[8].btns[0].disabledFn = () => { return this.nzDisabled }
    this.clustersTableBody[8].btns[0].click = (item:any) => { this.offline(item.data) }
    this.clustersTableBody[8].btns[1].disabledFn = () => { return this.nzDisabled }
    this.clustersTableBody[8].btns[1].click = (item:any) => { this.enable(item.data) }

    // 状态:已下线&启用 - 可操作:上线,禁用
    this.clustersTableBody[9].btns[0].disabledFn = () => { return this.nzDisabled }
    this.clustersTableBody[9].btns[0].click = (item:any) => { this.updateOrOnline(item.data, '上线') }
    this.clustersTableBody[9].btns[1].disabledFn = () => { return this.nzDisabled }
    this.clustersTableBody[9].btns[1].click = (item:any) => { this.disable(item.data) }

    // 状态:已下线&禁用 - 可操作:上线,启用
    this.clustersTableBody[10].btns[0].disabledFn = () => { return this.nzDisabled }
    this.clustersTableBody[10].btns[0].click = (item:any) => { this.updateOrOnline(item.data, '上线') }
    this.clustersTableBody[10].btns[1].disabledFn = () => { return this.nzDisabled }
    this.clustersTableBody[10].btns[1].click = (item:any) => { this.enable(item.data) }
  }

  ngAfterViewInit () {
    this.clustersTableBody[2].title = this.clusterStatusTpl
    this.clustersTableBody[3].title = this.disabledStatusTpl
  }

  getClustersData () {
    this.api.get('router/onlines', { uuid: this.apiUuid })
      .subscribe((resp:{code:number, data:{clusters:PublishManagementData[]}, msg:string}) => {
        if (resp.code === 0) {
          this.clustersList = resp.data.clusters
        } else {
          this.message.error(resp.msg || '获取列表数据失败!')
        }
      })
  }

  modalRef:NzModalRef|undefined

  updateOrOnline (item:PublishManagementData, type:string) {
    this.api.put('router/online', { clusterName: item.name || '' }, { uuid: this.apiUuid })
      .subscribe((resp:{code:number, data :{router:{params:{serviceName:string, templateUuid:string}, name:string}}, msg:string}) => {
        if (resp.code === 0) {
          this.message.success(resp.msg || (type + '成功'), { nzDuration: 1000 })
          this.getClustersData()
        } else {
          this.type = type
          this.solutionRouter = resp.data?.router?.name ? resp.data.router.name : ''
          this.solutionParam = resp.data?.router?.params ? resp.data.router.params : {}
          if (this.solutionRouter) {
            this.modalRef = this.publishFailModal.openFooterModal(resp.msg, 'API', [
              {
                label: '取消',
                type: 'default',
                onClick: () => {
                  this.modalRef?.close()
                  return true
                }
              },
              {
                label: '跳转至模板上线管理',
                type: 'primary',
                onClick: () => {
                  window.open(`router/plugin/content/${resp.data.router.params.templateUuid}`)
                  return true
                },
                disabled: !resp.data.router.params.templateUuid
              },
              {
                label: '跳转至服务上线管理',
                type: 'primary',
                onClick: () => {
                  window.open(`upstream/upstream/content/${resp.data.router.params.serviceName}/publish`)
                  return true
                },
                disabled: !resp.data.router.params.serviceName
              }
            ])
          } else {
            this.message.error(resp.msg || (type + '失败'))
          }
        }
      })
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  offline (item:PublishManagementData) {
    this.api.put('router/offline', { clusterName: (item.name || '') }, { uuid: this.apiUuid })
      .subscribe((resp:EmptyHttpResponse) => {
        if (resp.code === 0) {
          this.message.success(resp.msg || '下线成功', { nzDuration: 1000 })
          this.getClustersData()
        } else {
          this.type = '下线'
          this.message.error(resp.msg || '下线失败')
        }
      })
  }

  enable (item:PublishManagementData) {
    this.api.put('router/enable', { clusterName: (item.name || '') }, { uuid: this.apiUuid })
      .subscribe((resp:EmptyHttpResponse) => {
        if (resp.code === 0) {
          this.message.success(resp.msg || '启用成功', { nzDuration: 1000 })
          this.getClustersData()
        } else {
          this.message.error(resp.msg || '启用失败')
        }
      })
  }

  disable (item:PublishManagementData) {
    this.api.put('router/disable', { clusterName: (item.name || '') }, { uuid: this.apiUuid })
      .subscribe((resp:EmptyHttpResponse) => {
        if (resp.code === 0) {
          this.message.success(resp.msg || '禁用成功', { nzDuration: 1000 })
          this.getClustersData()
        } else {
          this.message.error(resp.msg || '禁用失败')
        }
      })
  }
}

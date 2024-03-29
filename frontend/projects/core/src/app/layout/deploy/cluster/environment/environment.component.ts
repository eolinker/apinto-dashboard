/* eslint-disable dot-notation */
/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2022-07-20 22:34:58
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2022-08-30 23:21:12
 * @FilePath: /apinto/src/app/layout/deploy/deploy-cluster-environment/deploy-cluster-environment.component.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgFeedbackModalService, EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { MODAL_NORMAL_SIZE, MODAL_SMALL_SIZE } from 'projects/core/src/app/constant/app.config'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { EoNgNavigationService } from 'projects/core/src/app/service/eo-ng-navigation.service'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { DeployService } from '../../deploy.service'
import { DeployClusterEnvConfigThead } from '../types/conf'
import { DeployClusterEnvironmentConfigFormComponent } from './config/form/form.component'
import { DeployClusterEnvironmentConfigUpdateComponent } from './config/update/update.component'
import { DeployClusterEnvironmentHistoryChangeComponent } from './history/change/change.component'
import { DeployClusterEnvironmentHistoryPublishComponent } from './history/publish/publish.component'
import { DeployClusterEnvironmentPublishComponent } from './publish/publish.component'

@Component({
  selector: 'eo-ng-deploy-cluster-environment',
  templateUrl: './environment.component.html',
  styles: [
    `label{
      line-height:32px;
    }

    .ant-col-4{
      text-align:right;
    }

    :host{
      overflow-y:auto;
      height:100%;
      display:block;
    }
`
  ]
})
export class DeployClusterEnvironmentComponent implements OnInit {
  @ViewChild('operatorStatusTpl', { read: TemplateRef, static: true }) operatorStatusTpl: TemplateRef<any> | undefined
  @ViewChild('publishStatusTpl', { read: TemplateRef, static: true }) publishStatusTpl: TemplateRef<any> | undefined
  @ViewChild('publishTypeTpl', { read: TemplateRef, static: true }) publishTypeTpl: TemplateRef<any> | undefined

  clusterName:string=''
  readonly nowUrl:string = this.router.routerState.snapshot.url
   configsList: Array<{ key: string, value: string, variableId: number, publish:string, status:string, desc:string, operator:string, updateTime:string, createTime:string, id: number, checked:boolean}> = []

  modalRef:NzModalRef | undefined
  nzDisabled:boolean = false

  configsTableHeadName: THEAD_TYPE[]= [...DeployClusterEnvConfigThead]
  configsTableBody: TBODY_TYPE[]=[...this.service.createClusterEnvConfigTbody(this)]

  autoTips: Record<string, Record<string, string>> = defaultAutoTips

  constructor (
          private baseInfo:BaseInfoService,
          private message: EoNgFeedbackMessageService,
          private modalService:EoNgFeedbackModalService,
          private api:ApiService,
          private router:Router,
          private navigationService:EoNgNavigationService,
          private service:DeployService) {
    this.navigationService.reqFlashBreadcrumb([{ title: '网关集群', routerLink: 'deploy/cluster' }, { title: '全局变量' }])
  }

  ngOnInit (): void {
    this.clusterName = this.baseInfo.allParamsInfo.clusterName
    if (!this.clusterName) {
      this.router.navigate(['/'])
    }

    this.getConfigsList()
  }

  ngAfterViewInit () {
    this.configsTableBody[0].title = this.operatorStatusTpl
    this.configsTableBody[3].title = this.publishStatusTpl
  }

  configTableClick = (item:any) => {
    this.openDrawer('editConfig', item.data)
  }

  getConfigsList () {
    this.api.get('cluster/' + this.clusterName + '/variables').subscribe(resp => {
      if (resp.code === 0) {
        this.configsList = resp.data.variables
      }
    })
  }

  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  delete (item:any, e?:Event) {
    e?.stopPropagation()
    this.modalService.create({
      nzTitle: '删除',
      nzContent: '该数据删除后将无法找回，请确认是否删除？',
      nzClosable: true,
      nzClassName: 'delete-modal',
      nzWidth: MODAL_SMALL_SIZE,
      nzOkDanger: true,
      nzOnOk: () => {
        this.deleteConfig(item)
      }
    })
  }

  deleteConfig (item:any) {
    this.api.delete('cluster/' + this.clusterName + '/variable', { key: item.key }).subscribe(resp => {
      if (resp.code === 0) {
        this.message.success(resp.msg || '删除成功', { nzDuration: 1000 })
        this.getConfigsList()
      }
    })
  }

  openDrawer (usage:string, data?:any, e?:Event) {
    e?.stopPropagation()
    switch (usage) {
      case 'addConfig':
        this.modalRef = this.modalService.create({
          nzTitle: '新建配置',
          nzWidth: MODAL_SMALL_SIZE,
          nzContent: DeployClusterEnvironmentConfigFormComponent,
          nzComponentParams: { clusterName: this.clusterName, closeModal: this.closeModal },
          nzOkText: '保存',
          nzOnOk: (component:DeployClusterEnvironmentConfigFormComponent) => {
            component.save()
            return false
          }
        })
        break
      case 'editConfig':
        this.modalRef = this.modalService.create({
          nzTitle: '编辑配置',
          nzWidth: MODAL_SMALL_SIZE,
          nzContent: DeployClusterEnvironmentConfigFormComponent,
          nzComponentParams: { clusterName: this.clusterName, editData: data, closeModal: this.closeModal },
          nzOkDisabled: this.nzDisabled,
          nzOkText: '提交',
          nzOnOk: (component:DeployClusterEnvironmentConfigFormComponent) => {
            component.save()
            return false
          }
        })
        break
      case 'updateConfig':
        this.modalRef = this.modalService.create({
          nzTitle: '同步配置',
          nzWidth: MODAL_NORMAL_SIZE,
          nzContent: DeployClusterEnvironmentConfigUpdateComponent,
          nzComponentParams: { clusterName: this.clusterName, closeModal: this.closeModal },
          nzOkDisabled: this.nzDisabled,
          nzOkText: '提交',
          nzOnOk: (component:DeployClusterEnvironmentConfigUpdateComponent) => {
            component.save()
            return false
          }
        })
        break
      case 'operateRecords':
        this.modalRef = this.modalService.create({
          nzTitle: '更改历史',
          nzWidth: MODAL_NORMAL_SIZE,
          nzContent: DeployClusterEnvironmentHistoryChangeComponent,
          nzComponentParams: {
            publishTypeTpl: this.publishTypeTpl,
            clusterName: this.clusterName
          },
          nzFooter: null
        })
        break
      case 'publishRecords':
        this.modalRef = this.modalService.create({
          nzTitle: '发布历史',
          nzWidth: MODAL_NORMAL_SIZE,
          nzContent: DeployClusterEnvironmentHistoryPublishComponent,
          nzComponentParams: {
            clusterName: this.clusterName,
            publishTypeTpl: this.publishTypeTpl
          },
          nzFooter: null
        })
        break
      case 'publish':
        this.modalRef = this.modalService.create({
          nzTitle: '发布',
          nzWidth: MODAL_NORMAL_SIZE,
          nzContent: DeployClusterEnvironmentPublishComponent,
          nzComponentParams: {
            clusterName: this.clusterName,
            publishTypeTpl: this.publishTypeTpl,
            closeModal: this.closeModal
          },
          nzOkDisabled: this.nzDisabled,
          nzOkText: '提交',
          nzOnOk: (component:DeployClusterEnvironmentPublishComponent) => {
            component.save(usage)
            return false
          }
        })
        break
    }
  }

  closeModal = (fresh?:boolean) => {
    fresh && this.getConfigsList()
    this.modalRef?.close()
  }

  copyCallback () {
    this.message.success('复制成功', {
      nzDuration: 1000
    })
  }
}

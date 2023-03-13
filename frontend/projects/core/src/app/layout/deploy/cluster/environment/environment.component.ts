/* eslint-disable camelcase */
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
import { NzModalRef } from 'ng-zorro-antd/modal'
import { MODAL_SMALL_SIZE } from 'projects/core/src/app/constant/app.config'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { MODAL_NORMAL_SIZE } from 'projects/eo-ng-apinto-user/src/public-api'
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
`
  ]
})
export class DeployClusterEnvironmentComponent implements OnInit {
  @ViewChild('operatorStatusTpl', { read: TemplateRef, static: true }) operatorStatusTpl: TemplateRef<any> | undefined
  @ViewChild('publishStatusTpl', { read: TemplateRef, static: true }) publishStatusTpl: TemplateRef<any> | undefined
  @ViewChild('publishTypeTpl', { read: TemplateRef, static: true }) publishTypeTpl: TemplateRef<any> | undefined

  clusterName:string=''
  readonly nowUrl:string = this.router.routerState.snapshot.url
   configsList: Array<{ key: string, value: string, variable_id: number, publish:string, status:string, desc:string, operator:string, update_time:string, create_time:string, id: number, checked:boolean}> = []

  drawerRef:NzModalRef | undefined
  nzDisabled:boolean = false

  configsTableHeadName: Array<object> = [
    { title: 'KEY', resizeable: true },
    { title: 'VALUE', resizeable: true },
    { title: '描述', resizeable: true },
    { title: '发布状态', resizeable: true },
    { title: '更新者', resizeable: true },
    { title: '更新时间' },
    {
      title: '操作',
      right: true
    }
  ]

  configsTableBody: Array<any> =[
    {
      key: 'key'
    },
    {
      key: 'value'
    },
    {
      key: 'desc'
    },
    {
      key: 'publish'
    },
    {
      key: 'operator'
    },
    {
      key: 'update_time'
    },
    {
      type: 'btn',
      right: true,
      showFn: (item:any) => {
        return item.publish !== 'DEFECT'
      },
      btns: [
        {
          title: '编辑',
          disabledFn: () => {
            return this.nzDisabled
          },
          click: (item:any) => {
            this.openDrawer('editConfig', item.data)
          }
        },
        {
          title: '删除',
          click: (item:any) => {
            this.delete(item.data)
          },
          disabledFn: () => {
            return this.nzDisabled
          }
        }
      ]
    },
    {
      type: 'btn',
      right: true,
      showFn: (item:any) => {
        return item.publish === 'DEFECT'
      },
      btns: [
        {
          title: '编辑',
          click: (item:any) => {
            this.openDrawer('editConfig', item.data)
          },
          disabledFn: () => {
            return this.nzDisabled
          }
        }
      ]
    }

  ]

  autoTips: Record<string, Record<string, string>> = defaultAutoTips

  // eslint-disable-next-line no-useless-constructor
  constructor (
          private baseInfo:BaseInfoService,
          private message: EoNgFeedbackMessageService,
          private modalService:EoNgFeedbackModalService,
          private api:ApiService,
          private router:Router,
          private appConfigService:AppConfigService) {
    this.appConfigService.reqFlashBreadcrumb([{ title: '网关集群', routerLink: 'deploy/cluster' }, { title: '环境变量' }])
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
      } else {
        this.message.error(resp.msg || '获取列表数据失败!')
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
      } else {
        this.message.error(resp.msg || '删除失败!')
      }
    })
  }

  openDrawer (usage:string, data?:any, e?:Event) {
    e?.stopPropagation()
    switch (usage) {
      case 'addConfig':
        this.drawerRef = this.modalService.create({
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
        this.drawerRef = this.modalService.create({
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
        this.drawerRef = this.modalService.create({
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
        this.drawerRef = this.modalService.create({
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
        this.drawerRef = this.modalService.create({
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
        this.drawerRef = this.modalService.create({
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
    this.drawerRef?.close()
  }
}

import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { THEAD_TYPE, TBODY_TYPE } from 'eo-ng-table'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { DeployService } from '../../deploy.service'
import { DeployClusterPluginPublishComponent } from './publish/publish.component'
import { DeployClusterPluginThead } from '../types/conf'
import { ClusterPluginItem } from '../types/types'
import { DeployClusterPluginConfigFormComponent } from './config/config.component'
import { DeployClusterPluginHistoryChangeComponent } from './history/change/change.component'
import { DeployClusterPluginHistoryPublishComponent } from './history/publish/publish.component'
import { MODAL_NORMAL_SIZE, MODAL_SMALL_SIZE } from 'projects/core/src/app/constant/app.config'

@Component({
  selector: 'eo-ng-deploy-cluster-plugin',
  templateUrl: './plugin.component.html',
  styles: [
  ]
})
export class DeployClusterPluginComponent implements OnInit {
  @ViewChild('operatorStatusTpl', { read: TemplateRef, static: true }) operatorStatusTpl: TemplateRef<any> | undefined
  @ViewChild('publishStatusTpl', { read: TemplateRef, static: true }) publishStatusTpl: TemplateRef<any> | undefined
  @ViewChild('publishTypeTpl', { read: TemplateRef, static: true }) publishTypeTpl: TemplateRef<any> | undefined
  @ViewChild('pluginStatusTpl', { read: TemplateRef, static: true }) pluginStatusTpl: TemplateRef<any> | undefined

  clusterName:string=''
  readonly nowUrl:string = this.router.routerState.snapshot.url
   pluginsList: ClusterPluginItem[] = []

  drawerRef:NzModalRef | undefined
  nzDisabled:boolean = false

  pluginsTableHeadName: THEAD_TYPE[]= [...DeployClusterPluginThead]
  pluginsTableBody: TBODY_TYPE[]=[...this.service.createPluginTbody(this)]

  autoTips: Record<string, Record<string, string>> = defaultAutoTips

  // eslint-disable-next-line no-useless-constructor
  constructor (
          private baseInfo:BaseInfoService,
          private message: EoNgFeedbackMessageService,
          private modalService:EoNgFeedbackModalService,
          private api:ApiService,
          private router:Router,
          private appConfigService:AppConfigService,
          private service:DeployService) {
    this.appConfigService.reqFlashBreadcrumb([{ title: '网关集群', routerLink: 'deploy/cluster' }, { title: '插件管理' }])
  }

  ngOnInit (): void {
    this.clusterName = this.baseInfo.allParamsInfo.clusterName
    if (!this.clusterName) {
      this.router.navigate(['/'])
    }
    this.getPluginsList()
  }

  ngAfterViewInit () {
    this.pluginsTableBody[0].title = this.operatorStatusTpl
    this.pluginsTableBody[1].title = this.publishStatusTpl
    this.pluginsTableBody[2].title = this.pluginStatusTpl
  }

  pluginTableClick = (item:any) => {
    this.openDrawer('editConfig', item.data)
  }

  getPluginsList () {
    this.api.get('cluster/' + this.clusterName + '/plugins').subscribe((resp:{code:number, data:{plugins:ClusterPluginItem[]}, msg:string}) => {
      if (resp.code === 0) {
        this.pluginsList = resp.data.plugins
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
        this.getPluginsList()
      } else {
        this.message.error(resp.msg || '删除失败!')
      }
    })
  }

  openDrawer (usage:string, data?:any, e?:Event) {
    e?.stopPropagation()
    switch (usage) {
      case 'editConfig':
        this.drawerRef = this.modalService.create({
          nzTitle: '编辑配置',
          nzWidth: MODAL_NORMAL_SIZE,
          nzContent: DeployClusterPluginConfigFormComponent,
          nzComponentParams: { clusterName: this.clusterName, editData: data, closeModal: this.closeModal },
          nzOkDisabled: this.nzDisabled,
          nzOkText: '提交',
          nzOnOk: (component:DeployClusterPluginConfigFormComponent) => {
            component.save()
            return false
          }
        })
        break
      case 'operateRecords':
        this.drawerRef = this.modalService.create({
          nzTitle: '更改历史',
          nzWidth: MODAL_NORMAL_SIZE,
          nzContent: DeployClusterPluginHistoryChangeComponent,
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
          nzContent: DeployClusterPluginHistoryPublishComponent,
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
          nzContent: DeployClusterPluginPublishComponent,
          nzComponentParams: {
            clusterName: this.clusterName,
            publishTypeTpl: this.publishTypeTpl,
            closeModal: this.closeModal
          },
          nzOkDisabled: this.nzDisabled,
          nzOkText: '提交',
          nzOnOk: (component:DeployClusterPluginPublishComponent) => {
            component.save()
            return false
          }
        })
        break
    }
  }

  closeModal = (fresh?:boolean) => {
    fresh && this.getPluginsList()
    this.drawerRef?.close()
  }
}

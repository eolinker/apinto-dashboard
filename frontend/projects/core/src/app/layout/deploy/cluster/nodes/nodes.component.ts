/* eslint-disable camelcase */
/* eslint-disable dot-notation */
import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { FormGroup } from '@angular/forms'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { MODAL_NORMAL_SIZE } from 'projects/core/src/app/constant/app.config'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { DeployClusterNodeTbody, DeployClusterNodeThead } from '../types/conf'
import { DeployClusterNodesFormComponent } from './form/form.component'

@Component({
  selector: 'eo-ng-deploy-cluster-nodes',
  templateUrl: './nodes.component.html',
  styles: [
    `
`
  ]
})
export class DeployClusterNodesComponent implements OnInit {
  @ViewChild('nodeStatusTpl', { read: TemplateRef, static: true }) nodeStatusTpl: TemplateRef<any> | undefined
  @ViewChild('adminAddrTpl', { read: TemplateRef, static: true }) adminAddrTpl: TemplateRef<any> | undefined
  @ViewChild('serviceAddrTpl', { read: TemplateRef, static: true }) serviceAddrTpl: TemplateRef<any> | undefined
  clusterName:string=''
  readonly nowUrl:string = this.router.routerState.snapshot.url

  nodesForms:{nodes:Array<{id:number, name:string, serviceAddr:string, adminAddr:string, status:string}>, isUpdate:boolean}=
    {
      nodes: [],
      isUpdate: false
    }

  nodesTableHeadName:THEAD_TYPE[] = [...DeployClusterNodeThead]
  nodesTableBody: TBODY_TYPE[]=[...DeployClusterNodeTbody]

  autoTips: Record<string, Record<string, string>> = defaultAutoTips

  nzDisabled:boolean = false
  validateResetNodeForm: FormGroup = new FormGroup({})
  constructor (private modalService:EoNgFeedbackModalService,
                private baseInfo:BaseInfoService,
                private message: EoNgFeedbackMessageService,
                private api:ApiService, private router:Router,
                private appConfigService:AppConfigService) {
    this.appConfigService.reqFlashBreadcrumb([{ title: '网关集群', routerLink: 'deploy/cluster' }, { title: '网关节点' }])
  }

  ngOnInit (): void {
    this.clusterName = this.baseInfo.allParamsInfo.clusterName
    if (!this.clusterName) {
      this.router.navigate(['/'])
    }
    this.getNodeslist()
  }

  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  getNodeslist () {
    this.api.get('cluster/' + this.clusterName + '/nodes').subscribe(resp => {
      if (resp.code === 0) {
        this.nodesForms = resp.data
      } else {
        this.message.error(resp.msg || '获取列表数据失败！', { nzDuration: 1000 })
      }
    })
  }

  ngAfterViewInit () {
    this.nodesTableBody[3].title = this.nodeStatusTpl
    this.nodesTableBody[2].title = this.serviceAddrTpl
    this.nodesTableBody[1].title = this.adminAddrTpl
  }

  updateNodes () {
    this.api.put('cluster/' + this.clusterName + '/node').subscribe(resp => {
      if (resp.code === 0) {
        this.message.success(resp.msg || '操作成功！', { nzDuration: 1000 })
        this.getNodeslist()
      } else {
        this.message.error(resp.msg || '操作失败！')
      }
    })
  }

  nodesDrawerRef:NzModalRef | undefined

  openDrawer (usage:string) {
    switch (usage) {
      case 'resetNodes':
        this.nodesDrawerRef = this.modalService.create({
          nzTitle: '重置配置',
          nzWidth: MODAL_NORMAL_SIZE,
          nzContent: DeployClusterNodesFormComponent,
          nzComponentParams: {
            nodeStatusTpl: this.nodeStatusTpl,
            serviceAddrTpl: this.serviceAddrTpl,
            adminAddrTpl: this.adminAddrTpl,
            nodesTableBody: this.nodesTableBody,
            clusterName: this.clusterName,
            closeModal: this.closeModal
          },
          nzOkDisabled: this.nzDisabled,
          nzOkText: '提交',
          nzOnOk: (component:DeployClusterNodesFormComponent) => {
            component.save()
            return false
          }
        })
    }
  }

  closeModal = () => {
    this.getNodeslist()
    this.nodesDrawerRef?.close()
  }
}

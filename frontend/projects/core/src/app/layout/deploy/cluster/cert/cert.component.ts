/* eslint-disable camelcase */
/* eslint-disable dot-notation */
import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgFeedbackModalService, EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { MODAL_SMALL_SIZE } from 'projects/core/src/app/constant/app.config'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { DeployClusterCertFormComponent } from './form/form.component'

@Component({
  selector: 'eo-ng-deploy-cluster-cert',
  templateUrl: './cert.component.html',
  styles: [
    `
    `
  ]
})
export class DeployClusterCertComponent implements OnInit {
  @ViewChild('certContentTpl', { read: TemplateRef, static: true }) certContentTpl: TemplateRef<any> | undefined

  readonly nowUrl:string = this.router.routerState.snapshot.url

  clusterName:string=''

  certsList:Array<{id:string, name:string, valid_time:string, operator:string, create_time:string, update_time:string}>=[]

  nzDisabled:boolean = false

  certsTableHeadName: Array<object> = [
    { title: '证书' },
    { title: '证书有效期' },
    { title: '更新者' },
    { title: '更新时间' },
    {
      title: '操作',
      right: true
    }
  ]

  certsTableBody: Array<any> =[
    { key: 'name' },
    { key: 'valid_time' },
    { key: 'operator' },
    { key: 'update_time' },
    {
      type: 'btn',
      right: true,
      btns: [
        {
          title: '修改',
          disabledFn: () => { return this.nzDisabled },
          click: (item:any) => {
            this.openDrawer('editCert', item.data)
          }
        },
        {
          title: '删除',
          disabledFn: () => { return this.nzDisabled },
          click: (item:any) => {
            this.delete(item.data)
          }
        }
      ]
    }
  ]

  modalRef:NzModalRef | undefined

  constructor (
                private message: EoNgFeedbackMessageService,
                private modalService:EoNgFeedbackModalService,
                public api:ApiService,
                private baseInfo:BaseInfoService,
                private router:Router,
                private appConfigService:AppConfigService) {
    this.appConfigService.reqFlashBreadcrumb([{ title: '网关集群', routerLink: 'deploy/cluster' }, { title: '证书管理' }])
  }

  ngOnInit (): void {
    this.clusterName = this.baseInfo.allParamsInfo.clusterName
    if (!this.clusterName) {
      this.router.navigate(['/'])
    }
    this.getCertsList()
  }

  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  certTableClick = (item:any) => {
    this.openDrawer('editCert', item.data)
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
        this.deleteCert(item)
      }
    })
  }

  getCertsList () {
    this.api.get('cluster/' + this.clusterName + '/certificates').subscribe(resp => {
      if (resp.code === 0) {
        this.certsList = resp.data.certificates
      } else {
        this.message.error(resp.msg || '获取列表数据失败!')
      }
    })
  }

  openDrawer (usage:string, data?:any, e?:Event):void {
    e?.stopPropagation()
    this.modalRef = this.modalService.create({
      nzTitle: usage === 'addCert' ? '新建证书' : '修改证书',
      nzWidth: MODAL_SMALL_SIZE,
      nzContent: DeployClusterCertFormComponent,
      nzComponentParams: { certId: data?.id, clusterName: this.clusterName, closeModal: this.closeModal },
      nzOkText: usage === 'addCert' ? '保存' : '提交',
      nzOnOk: (component:DeployClusterCertFormComponent) => {
        component.save(usage)
        return false
      }
    })
  }

  closeModal = () => {
    this.getCertsList()
    this.modalRef?.close()
  }

  deleteCert (row:{id:number, name:string, operator:string, update_time:string, valid_time:string, create_time:string}):void {
    this.api.delete('cluster/' + (this.clusterName || '') + '/certificate/' + (row.id || '')).subscribe(resp => {
      if (resp.code === 0) {
        this.getCertsList()
        this.message.success(resp.msg || '删除成功!', { nzDuration: 1000 })
      } else {
        this.message.error(resp.msg || '删除证书失败!', { nzDuration: 1000 })
      }
    })
  }
}

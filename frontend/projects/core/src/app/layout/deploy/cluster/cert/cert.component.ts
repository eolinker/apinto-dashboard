/* eslint-disable camelcase */
import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgFeedbackModalService, EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { EoNgNavigationService } from 'projects/core/src/app/service/eo-ng-navigation.service'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { MODAL_SMALL_SIZE } from 'projects/core/src/app/constant/app.config'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { DeployClusterCertFormComponent } from './form/form.component'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { deployCertsTableBody, deployCertsTableHeadName } from '../types/conf'
import { DeployCertListData } from '../types/types'
import { EmptyHttpResponse } from 'projects/core/src/app/constant/type'

@Component({
  selector: 'eo-ng-deploy-cluster-cert',
  templateUrl: './cert.component.html',
  styles: [
    `
    `
  ]
})
export class DeployClusterCertComponent implements OnInit {
  @ViewChild('dnsNameTpl', { read: TemplateRef, static: true }) dnsNameTpl: TemplateRef<any> | undefined
  clusterName:string=''
  nzDisabled:boolean = false
  modalRef:NzModalRef | undefined
  certsList:DeployCertListData[]=[]
  certsTableHeadName: THEAD_TYPE[] = [...deployCertsTableHeadName]
  certsTableBody: TBODY_TYPE[] =[...deployCertsTableBody]

  constructor (
                private message: EoNgFeedbackMessageService,
                private modalService:EoNgFeedbackModalService,
                private api:ApiService,
                private baseInfo:BaseInfoService,
                private router:Router,
                private navigationService:EoNgNavigationService) {
    this.navigationService.reqFlashBreadcrumb([{ title: '网关集群', routerLink: 'deploy/cluster' }, { title: 'SSL证书' }])
  }

  ngOnInit (): void {
    this.clusterName = this.baseInfo.allParamsInfo.clusterName
    if (!this.clusterName) {
      this.router.navigate(['/'])
    }
    this.certsTableBody[1].title = this.dnsNameTpl
    this.certsTableBody[5].btns[0].disabledFn = () => { return this.nzDisabled }
    this.certsTableBody[5].btns[0].click = (item:any) => { this.openDrawer('editCert', item.data) }
    this.certsTableBody[5].btns[1].disabledFn = () => { return this.nzDisabled }
    this.certsTableBody[5].btns[1].click = (item:any) => { this.delete(item.data) }

    this.getCertsList()
  }

  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  certTableClick = (item:{data:DeployCertListData}) => {
    this.openDrawer('editCert', item.data)
  }

  delete (item:DeployCertListData, e?:Event) {
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
    this.api.get('cluster/' + this.clusterName + '/certificates')
      .subscribe((resp:{code:number, data:{certificates:DeployCertListData[]}, msg:string}) => {
        if (resp.code === 0) {
          this.certsList = resp.data.certificates.map((cert:DeployCertListData) => {
            cert.dnsNameStr = cert.dnsName.join('/')
            return cert
          })
        }
      })
  }

  openDrawer (usage:string, data?:DeployCertListData, e?:Event):void {
    e?.stopPropagation()
    this.modalRef = this.modalService.create({
      nzTitle: usage === 'addCert' ? '新建证书' : '修改证书',
      nzWidth: MODAL_SMALL_SIZE,
      nzContent: DeployClusterCertFormComponent,
      nzComponentParams: { certId: data?.id, clusterName: this.clusterName, editPage: usage === 'editCert', closeModal: this.closeModal },
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

  deleteCert (row:DeployCertListData):void {
    this.api.delete('cluster/' + (this.clusterName || '') + '/certificate/' + (row.id || ''))
      .subscribe((resp:EmptyHttpResponse) => {
        if (resp.code === 0) {
          this.getCertsList()
          this.message.success(resp.msg || '删除成功!', { nzDuration: 1000 })
        }
      })
  }

  copyCallback (event?:Event) {
    event?.stopPropagation()
    this.message.success('复制成功', {
      nzDuration: 1000
    })
  }
}

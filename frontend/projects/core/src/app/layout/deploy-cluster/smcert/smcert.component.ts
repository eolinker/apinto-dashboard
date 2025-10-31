/* eslint-disable camelcase */
import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgFeedbackModalService, EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { EoNgNavigationService } from 'projects/core/src/app/service/eo-ng-navigation.service'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { MODAL_SMALL_SIZE } from 'projects/core/src/app/constant/app.config'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { DeployClusterSmcertFormComponent } from './form/form.component'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { deploySmcertsTableBody, deploySmcertsTableHeadName } from '../types/conf'
import { DeploySmcertListData } from '../types/types'
import { EmptyHttpResponse } from 'projects/core/src/app/constant/type'

@Component({
  selector: 'eo-ng-deploy-cluster-smcert',
  templateUrl: './smcert.component.html',
  styles: [
    `
    :host{
      overflow-y:auto;
      height:100%;
      display:block;
    }
    `
  ]
})
export class DeployClusterSmcertComponent implements OnInit {
  @ViewChild('signDnsNameTpl', { read: TemplateRef, static: true }) signDnsNameTpl: TemplateRef<any> | undefined
  @ViewChild('encDnsNameTpl', { read: TemplateRef, static: true }) encDnsNameTpl: TemplateRef<any> | undefined
  clusterName:string=''
  nzDisabled:boolean = false
  modalRef:NzModalRef | undefined
  certsList:DeploySmcertListData[]=[]
  certsTableHeadName: THEAD_TYPE[] = [...deploySmcertsTableHeadName]
  certsTableBody: TBODY_TYPE[] =[...deploySmcertsTableBody]

  constructor (
                private message: EoNgFeedbackMessageService,
                private modalService:EoNgFeedbackModalService,
                private api:ApiService,
                private baseInfo:BaseInfoService,
                private router:Router,
                private navigationService:EoNgNavigationService) {
    this.navigationService.reqFlashBreadcrumb([{ title: '网关集群', routerLink: 'deploy/cluster' }, { title: '国密证书' }])
  }

  ngOnInit (): void {
    this.clusterName = this.baseInfo.allParamsInfo.clusterName
    if (!this.clusterName) {
      this.router.navigate(['/'])
    }
    this.certsTableBody[1].title = this.signDnsNameTpl
    this.certsTableBody[3].title = this.encDnsNameTpl
    this.certsTableBody[7].btns[0].disabledFn = () => { return this.nzDisabled }
    this.certsTableBody[7].btns[0].click = (item:any) => { this.openDrawer('editSmcert', item.data) }
    this.certsTableBody[7].btns[1].disabledFn = () => { return this.nzDisabled }
    this.certsTableBody[7].btns[1].click = (item:any) => { this.delete(item.data) }

    this.getCertsList()
  }

  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  certTableClick = (item:{data:DeploySmcertListData}) => {
    this.openDrawer('editSmcert', item.data)
  }

  delete (item:DeploySmcertListData, e?:Event) {
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
    this.api.get('cluster/' + this.clusterName + '/gm_certificates')
      .subscribe((resp:{code:number, data:{certificates:DeploySmcertListData[]}, msg:string}) => {
        if (resp.code === 0) {
          this.certsList = resp.data.certificates.map((cert:DeploySmcertListData) => {
            cert.signDnsName = cert.signCert.dnsName
            cert.signValidTime = cert.signCert.validTime
            cert.encDnsName = cert.encCert.dnsName
            cert.encValidTime = cert.encCert.validTime
            return cert
          })
        }
      })
  }

  openDrawer (usage:string, data?:DeploySmcertListData, e?:Event):void {
    e?.stopPropagation()
    this.modalRef = this.modalService.create({
      nzTitle: usage === 'addSmcert' ? '新建证书' : '修改证书',
      nzWidth: MODAL_SMALL_SIZE,
      nzContent: DeployClusterSmcertFormComponent,
      nzComponentParams: { smcertId: data?.id, clusterName: this.clusterName, editPage: usage === 'editSmcert', closeModal: this.closeModal },
      nzOkText: usage === 'addSmcert' ? '保存' : '提交',
      nzOnOk: (component:DeployClusterSmcertFormComponent) => {
        component.save(usage)
        return false
      }
    })
  }

  closeModal = () => {
    this.getCertsList()
    this.modalRef?.close()
  }

  deleteCert (row:DeploySmcertListData):void {
    this.api.delete('cluster/' + (this.clusterName || '') + '/gm_certificate/' + (row.id || ''))
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

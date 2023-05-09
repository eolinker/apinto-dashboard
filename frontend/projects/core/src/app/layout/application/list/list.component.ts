/* eslint-disable camelcase */
/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2022-08-17 23:42:52
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2022-08-24 00:35:02
 * @FilePath: /apinto/src/app/layout/application/application-management-list/application-management-list.component.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { Component, OnInit } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { EoNgNavigationService } from 'projects/core/src/app/service/eo-ng-navigation.service'
import { MODAL_SMALL_SIZE } from '../../../constant/app.config'
import { EmptyHttpResponse } from '../../../constant/type'
import { applicationsTableBody, applicationsTableHeadName } from '../types/conf'
import { ApplicationListData } from '../types/types'

@Component({
  selector: 'eo-ng-application-management-list',
  templateUrl: 'list.component.html',
  styles: [
  ]
})
export class ApplicationManagementListComponent implements OnInit {
  applicationName:string = ''
  applicationNameForSear:string = ''
  nzDisabled:boolean = false
  applicationsForm: {applications:ApplicationListData[], total:number, pageNum:number, pageSize:number}={
    applications: [],
    total: 0,
    pageNum: 1,
    pageSize: 20
  }

  applicationsTableHeadName: THEAD_TYPE[]= [...applicationsTableHeadName]
  applicationsTableBody: TBODY_TYPE[] =[...applicationsTableBody]

  constructor (
    private message: EoNgFeedbackMessageService,
    private modalService:EoNgFeedbackModalService,
    private api:ApiService,
    private router:Router,
    private navigationService:EoNgNavigationService
  ) {
    this.navigationService.reqFlashBreadcrumb([{ title: '应用管理', routerLink: 'application' }])
  }

  ngOnInit (): void {
    this.applicationsTableBody[5].btns[0].click = (item:any) => { this.router.navigate(['/', 'application', 'content', item.data.id]) }
    this.applicationsTableBody[5].btns[1].disabledFn = (data:any, item:any) => { return !item.data.isDelete || this.nzDisabled }
    this.applicationsTableBody[5].btns[1].click = (item:any) => { this.delete(item.data) }
    this.getApplicationsList()
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  applicationsTableClick = (item:{data:ApplicationListData}) => {
    this.router.navigate(['/', 'application', 'content', item.data.id])
  }

  getApplicationsList () {
    this.api.get('applications', { name: this.applicationNameForSear, pageNum: this.applicationsForm.pageNum, pageSize: this.applicationsForm.pageSize })
      .subscribe((resp:{code:number, data:{applications:ApplicationListData[], total:number}, msg:string}) => {
        if (resp.code === 0) {
          this.applicationsForm.applications = resp.data.applications
          this.applicationsForm.total = resp.data.total
          this.applicationName = this.applicationNameForSear
        }
      })
  }

  addApplication () {
    this.router.navigate(['/', 'application', 'create'])
  }

  delete (item:ApplicationListData, e?:Event) {
    e?.stopPropagation()
    this.modalService.create({
      nzTitle: '删除',
      nzContent: '该数据删除后将无法找回，请确认是否删除？',
      nzClosable: true,
      nzWidth: MODAL_SMALL_SIZE,
      nzClassName: 'delete-modal',
      nzOkDanger: true,
      nzOnOk: () => {
        this.deleteDiscovery(item.id)
      }
    })
  }

  deleteDiscovery (id:string) {
    this.api.delete('application', { appId: id }).subscribe((resp:EmptyHttpResponse) => {
      if (resp.code === 0) {
        this.getApplicationsList()
        this.message.success(resp.msg || '删除成功', { nzDuration: 1000 })
      }
    })
  }
}

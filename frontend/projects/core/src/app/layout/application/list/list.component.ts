/* eslint-disable no-dupe-class-members */
/* eslint-disable camelcase */
/* eslint-disable no-useless-constructor */
/*
 * @Author:
 * @Date: 2022-08-17 23:42:52
 * @LastEditors:
 * @LastEditTime: 2022-08-24 00:35:02
 * @FilePath: /apinto/src/app/layout/application/application-management-list/application-management-list.component.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { Component, OnInit } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'
import { MODAL_SMALL_SIZE } from '../../../constant/app.config'

@Component({
  selector: 'eo-ng-application-management-list',
  templateUrl: 'list.component.html',
  styles: [
  ]
})
export class ApplicationManagementListComponent implements OnInit {
  constructor (private message: EoNgFeedbackMessageService,
    private modalService:EoNgFeedbackModalService,
     private api:ApiService,
     private router:Router,
     private appConfigService:AppConfigService) {
    this.appConfigService.reqFlashBreadcrumb([{ title: '应用管理', routerLink: 'application' }])
  }

  applicationName:string = ''
  applicationNameForSear:string = ''
  applicationsForm: {applications:Array<{name:string, id:string, desc:string, operator:string, update_time:string, is_delete:boolean}>, total:number, page_num:number, page_size:number}={
    applications: [],
    total: 0,
    page_num: 1,
    page_size: 20
  }

  nzDisabled:boolean = false

  applicationsTableHeadName: Array<object> = [
    { title: '应用名称' },
    { title: '应用ID' },
    { title: '描述' },
    { title: '更新者' },
    { title: '更新时间' },
    {
      title: '操作',
      right: true
    }
  ]

  applicationsTableBody: Array<any> =[
    { key: 'name' },
    { key: 'id' },
    { key: 'desc' },
    { key: 'operator' },
    { key: 'update_time' },
    {
      type: 'btn',
      right: true,
      btns: [
        {
          title: '查看',
          click: (item:any) => {
            this.router.navigate(['/', 'application', 'content', item.data.id])
          }
        },
        {
          title: '删除',
          disabledFn: (data:any, item:any) => {
            return !item.data.is_delete || this.nzDisabled
          },
          click: (item:any) => {
            this.delete(item.data)
          }
        }
      ]
    }

  ]

  ngOnInit (): void {
    this.getApplicationsList()
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  applicationsTableClick = (item:any) => {
    this.router.navigate(['/', 'application', 'content', item.data.id])
  }

  getApplicationsList () {
    this.api.get('applications', { name: this.applicationNameForSear, page_num: this.applicationsForm.page_num, page_size: this.applicationsForm.page_size }).subscribe(resp => {
      if (resp.code === 0) {
        this.applicationsForm.applications = resp.data.applications
        this.applicationsForm.total = resp.data.total
        this.applicationName = this.applicationNameForSear
      } else {
        this.message.error(resp.msg || '刷新列表失败!')
      }
    })
  }

  addApplication () {
    this.router.navigate(['/', 'application', 'create'])
  }

  delete (item:any, e?:Event) {
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
    this.api.delete('application', { app_id: id }).subscribe(resp => {
      if (resp.code === 0) {
        this.getApplicationsList()
        this.message.success(resp.msg || '删除成功', { nzDuration: 1000 })
      } else {
        this.message.error(resp.msg || '删除失败!')
      }
    })
  }
}

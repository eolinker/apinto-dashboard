/* eslint-disable camelcase */
/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2022-08-17 23:42:52
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2022-08-24 00:35:02
 * @FilePath: /apinto/src/app/layout/application/application-management-list/application-management-list.component.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { Router } from '@angular/router'
import { MODAL_NORMAL_SIZE, MODAL_SMALL_SIZE } from '../../../constant/app.config'
import { ClusterSimpleOption } from '../../../constant/type'
import { Subscription } from 'rxjs'
import { ApplicationCreateComponent } from '../create/create.component'
import { ApplicationListData } from '../types/types'
import { EoNgApplicationService } from '../application.service'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { SelectOption } from 'eo-ng-select'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { ApiService } from '../../../service/api.service'
import { EoNgNavigationService } from '../../../service/eo-ng-navigation.service'

@Component({
  selector: 'eo-ng-application-management-list',
  templateUrl: '../../../component/intelligent-plugin/list/list.component.html',
  styles: [
  ]
})
export class ApplicationManagementListComponent implements OnInit {
  @ViewChild('clusterStatusTpl', { read: TemplateRef, static: true }) clusterStatusTpl: TemplateRef<any> | undefined
  @ViewChild('loadingTpl', { read: TemplateRef, static: true }) loadingTpl: TemplateRef<any> | undefined
  moduleName:string = ''
  pluginName:string = '应用'
  keyword:string = ''
  nzDisabled:boolean = false
  cluster:any = []
  clusterOptions:SelectOption[] = []
  tableBody:TBODY_TYPE[] = []
  tableHeadName:THEAD_TYPE[] = [...this.service.createAppListThead(this)]
  tableData:{data:any[], pagination:boolean, total:number, pageNum:number, pageSize:number}
  = { data: [], pagination: true, total: 1, pageSize: 20, pageNum: 1 }

  driverOptions:SelectOption[] = []
  renderSchema:any = {} // 动态渲染数据，是json schema
  modalRef:NzModalRef|undefined
  statusMap:{[k:string]:any} = {}
  tableLoading:boolean = true
  subscription: Subscription = new Subscription()

  constructor (
    public message: EoNgFeedbackMessageService,
    public service:EoNgApplicationService,
    public modalService:EoNgFeedbackModalService,
    public api:ApiService,
    public router:Router,
    public navigationService: EoNgNavigationService) {

  }

  ngOnInit (): void {
    this.navigationService.reqFlashBreadcrumb([
      { title: '应用管理' }
    ])
    this.getClusters()
    this.getTableData()
  }

  getClusters () {
    this.api.get('clusters/simple').subscribe((resp:{code:number, msg:string, data:{clusters:ClusterSimpleOption[]}}) => {
      if (resp.code === 0) {
        this.clusterOptions = resp.data.clusters.map((cluster:ClusterSimpleOption) => {
          return { label: cluster.title, value: cluster.name }
        })
        this.cluster = this.clusterOptions.map((cluster:SelectOption) => {
          return cluster.value
        })
      }
    })
  }

  getTableData () {
    this.tableLoading = true
    this.api.get('applications', {
      name: this.keyword,
      pageNum: this.tableData.pageNum,
      pageSize: this.tableData.pageSize,
      clusters: JSON.stringify(this.cluster)
    }).subscribe((resp:{code:number, data:{applications:ApplicationListData[], total:number, pageNum:number, pageSize:number}}) => {
      if (resp.code === 0) {
        this.tableData.data = resp.data.applications.map((item:ApplicationListData) => {
          if (item.publish?.length > 0) {
            for (const p of item.publish) {
              item[`cluster_${p.name}`] = p.status
            }
          }
          return item
        })
        if (resp.data.applications.length > 0) {
          this.tableBody = this.service.createAppListTbody(this, resp.data.applications[0].publish)
          this.tableHeadName = this.service.createAppListThead(this, resp.data.applications[0].publish)
        }

        this.tableData.total = resp.data.total || this.tableData.total
        this.tableData.pageNum = resp.data.pageNum || this.tableData.pageNum
        this.tableData.pageSize = resp.data.pageSize || this.tableData.pageSize
      }
      this.tableLoading = false
    })
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  publish (value:any) {
    this.service.publishAppModal(value.data, this)
  }

  addData () {
    this.modalRef = this.modalService.create({
      nzTitle: `新建${this.pluginName}`,
      nzWidth: MODAL_NORMAL_SIZE,
      nzContent: ApplicationCreateComponent,
      nzComponentParams: {
        editPage: false
      },
      nzOnOk: (component:ApplicationCreateComponent) => {
        return new Promise((resolve, reject) => {
          component.saveApplication()?.subscribe((resp) => {
            if (resp) {
              resolve()
              this.getTableData()
            } else {
              reject(new Error())
            }
          })
        })
      }
    })
  }

  editData (value:any) {
    this.router.navigate(['/', 'application', 'content', value.data.id, 'message'])
  }

  deleteDataModal (items:{id:string, [k:string]:any}) {
    this.modalService.create({
      nzTitle: '删除',
      nzContent: '该数据删除后将无法找回，请确认是否删除？',
      nzClosable: true,
      nzClassName: 'delete-modal',
      nzWidth: MODAL_SMALL_SIZE,
      nzOkDanger: true,
      nzOnOk: () => {
        this.deleteData(items)
      }
    })
  }

  // 删除单条数据
   deleteData = (items:{id:string, [k:string]:any}) => {
     this.api.delete('application', { appId: items.id }).subscribe((resp:any) => {
       if (resp.code === 0) {
         this.message.success(resp.msg || '删除成功!')
         this.getTableData()
       }
     })
   }
}

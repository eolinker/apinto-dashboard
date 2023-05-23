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
import { NavigationEnd } from '@angular/router'
import { MODAL_NORMAL_SIZE } from '../../../constant/app.config'
import { EmptyHttpResponse } from '../../../constant/type'
import { forkJoin, map } from 'rxjs'
import { IntelligentPluginDefaultThead } from '../../../component/intelligent-plugin/types/conf'
import { DynamicConfig, DynamicDriverData } from '../../../component/intelligent-plugin/types/types'
import { v4 as uuidv4 } from 'uuid'
import { IntelligentPluginListComponent } from '../../../component/intelligent-plugin/list/list.component'
import { ApplicationPublishComponent } from '../publish/publish.component'
import { ApplicationCreateComponent } from '../create/create.component'

@Component({
  selector: 'eo-ng-application-management-list',
  templateUrl: '../../../component/intelligent-plugin/list/list.component.html',
  styles: [
  ]
})
export class ApplicationManagementListComponent extends IntelligentPluginListComponent implements OnInit {
  override pluginName:string = ''

  override ngOnInit (): void {
    this.navigationService.reqFlashBreadcrumb([
      { title: '应用管理' }
    ])
    this.subscription = this.router.events.subscribe((event) => {
      if (event instanceof NavigationEnd) {
        this.pluginName = ''
        this.keyword = ''
        this.cluster = []
        this.clusterOptions = []
        this.tableBody = [...this.service.createTbody(this)]
        this.tableHeadName = [...IntelligentPluginDefaultThead]
        this.tableData = { data: [], pagination: true, total: 1, pageSize: 20, pageNum: 1 }

        this.driverOptions = []
        this.renderSchema = {} // 动态渲染数据，是json schema
        this.modalRef = undefined
        this.statusMap = {}
        this.tableLoading = true
        this.getClusters()
        this.getTableData()
      }
    })
    this.getClusters()
    this.getTableData()
  }

  override getTableData () {
    this.tableLoading = true
    // 表格内的其他数据与状态数据是分别获取的，如果list先返回，需要先展示除了状态数据以外的其他数据
    // TODO 修改接口
    forkJoin([this.api.get('dynamic/application/list', {
      page: this.tableData.pageNum,
      pageSize: this.tableData.pageSize,
      keyword: this.keyword,
      cluster: JSON.stringify(this.cluster)
    }).pipe(
      map(res => {
        if (res.code === 0) {
          this.getConfig(res.data)
        }
        return res
      })),
    // TODO 修改接口
    this.api.get('dynamic/application/status', {
      page: this.tableData.pageNum,
      pageSize: this.tableData.pageSize,
      keyword: this.keyword,
      cluster: JSON.stringify(this.cluster)
    })]).subscribe((val:Array<any>) => {
      this.refreshTableData(this.tableData.data, val[1].data)
    })
  }

  // 获取列表渲染配置、表单渲染配置
  override getConfig (data:DynamicConfig) {
    this.pluginName = data.title
    this.getTableConfig(data.fields) // 获取列表配置
    this.tableData.data = data.list // 获取列表数据
    this.driverOptions = data.drivers?.map((driver:DynamicDriverData) => {
      return { label: driver.title, value: driver.name }
    }) || []
  }

  override publish (value:any) {
    this.modalRef = this.modalService.create({
      nzTitle: `${value.data.title}上线管理`,
      nzWidth: MODAL_NORMAL_SIZE,
      nzContent: ApplicationPublishComponent,
      nzComponentParams: {
        name: value.data.title,
        id: value.data.id,
        desc: value.data.description,
        moduleName: this.moduleName,
        closeModal: this.closeModal,
        nzDisabled: this.nzDisabled
      },
      nzFooter: [{
        label: '取消',
        type: 'default',
        onClick: () => {
          this.modalRef?.close()
        }
      },
      {
        label: '下线',
        danger: true,
        onClick: (context:ApplicationPublishComponent) => {
          return new Promise((resolve, reject) => {
            context.offline() ? resolve(true) : reject(new Error())
          })
        },
        disabled: () => {
          return this.nzDisabled
        }
      },
      {
        label: '上线',
        type: 'primary',
        onClick: (context:ApplicationPublishComponent) => {
          return new Promise((resolve, reject) => {
            context.online() ? resolve(true) : reject(new Error())
          })
        },
        disabled: () => {
          return this.nzDisabled
        }
      }]
    })
  }

  override addData () {
    this.modalRef = this.modalService.create({
      nzTitle: `新建${this.pluginName}`,
      nzWidth: MODAL_NORMAL_SIZE,
      nzContent: ApplicationCreateComponent,
      nzComponentParams: {
        editPage: false
      },
      nzOnOk: (component:ApplicationCreateComponent) => {
        return new Promise((resolve, reject) => {
          component.saveApplication() ? resolve() : reject(new Error())
        })
      }
    })
  }

  override editData (value:any) {
    console.log(value)
    this.router.navigate(['/', 'application', 'content', value.data.id, 'message'])
    // this.modalRef = this.modalService.create({
    //   nzTitle: `编辑${this.pluginName}`,
    //   nzWidth: MODAL_NORMAL_SIZE,
    //   nzContent: ApplicationCreateComponent,
    //   nzComponentParams: {
    //     editPage: true,
    //     appId: value.data.id
    //   },
    //   nzOnOk: (component:ApplicationCreateComponent) => {
    //     return new Promise((resolve, reject) => {
    //       component.saveApplication() ? resolve() : reject(new Error())
    //     })
    //   }
    // })
  }

  override saveData (form:{[k:string]:any}, id:string = uuidv4(), editPage?:boolean) {
    if (editPage) {
      // TODO
      this.api.put(`dynamic/application/config/${id}`, { ...form }).subscribe((resp:EmptyHttpResponse) => {
        if (resp.code === 0) {
          this.message.success(resp.msg || '操作成功')
          this.getTableData()
          this.modalRef?.close()
        }
      })
    } else {
      // TODO
      this.api.post('dynamic/application', { ...form }).subscribe((resp:EmptyHttpResponse) => {
        if (resp.code === 0) {
          this.message.success(resp.msg || '操作成功')
          this.getTableData()
          this.modalRef?.close()
        }
      })
    }
  }

  // 删除单条数据
  override deleteData = (items:{id:string, [k:string]:any}) => {
    // TODO
    this.api.delete('dynamic/batch', { uuids: JSON.stringify([items.id]) }).subscribe((resp:any) => {
      if (resp.code === 0) {
        this.message.success(resp.msg || '删除成功!')
        this.getTableData()
      }
    })
  }
}

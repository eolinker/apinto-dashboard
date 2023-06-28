import { Injectable } from '@angular/core'
import { ApplicationManagementListComponent } from './list/list.component'
import { EO_TBODY_TYPE } from 'projects/eo-ng-apinto-table/src/public-api'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { ApplicationData } from './types/types'
import { ApiService } from '../../service/api.service'
import { ApplicationPublishComponent } from './publish/publish.component'
import { FilterOpts } from '../../constant/conf'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { EoNgFeedbackModalService } from 'eo-ng-feedback'
import { MODAL_NORMAL_SIZE } from '../../constant/app.config'

@Injectable({
  providedIn: 'root'
})
export class EoNgApplicationService {
  appName:string = ''
  appDesc:string = ''

  appData:ApplicationData|null = null
  loading:boolean = true
  modalRef:NzModalRef|undefined
  constructor (private api:ApiService, private modalService:EoNgFeedbackModalService) {}

  getApplicationData (appId:string) {
    this.loading = true
    this.appData = null
    this.api.get('application', { appId }).subscribe((resp:{code:number, data:{application:ApplicationData}, msg:string}) => {
      if (resp.code === 0) {
        this.appData = resp.data.application
        this.appName = this.appData.name
        this.appDesc = this.appData.desc
      }
      this.loading = false
    })
  }

  clearData () {
    this.appName = ''
    this.appDesc = ''
    this.appData = null
  }

  createAppListThead (context:ApplicationManagementListComponent, publishList?:Array<any>):THEAD_TYPE[] {
    return [
      {
        title: '名称'
      },
      {
        title: 'ID'
      },
      {
        title: '描述'
      },
      ...(publishList?.length
        ? publishList.map((p) => {
          return {
            title: `状态：${p.title}`,
            tooltip: `状态：${p.title}`,
            titleString: `状态：${p.title}`,
            filterMultiple: true,
            filterOpts: [...FilterOpts],
            filterFn: (list: string[], item: any) => {
              return list.some((name) => item.data[`cluster_${p.name}`] === name)
            }
          }
        })
        : []),
      {
        title: '更新时间'
      },
      {
        title: '操作',
        right: true
      }
    ]
  }

  createAppListTbody (context:ApplicationManagementListComponent, publishList?:Array<any>):EO_TBODY_TYPE[] {
    return [
      {
        key: 'name',
        copy: true
      },
      {
        key: 'id'
      },
      {
        key: 'desc'
      },
      ...(publishList?.length
        ? publishList.map((p) => {
          return { key: `cluster_${p.name}`, title: context.clusterStatusTpl }
        })
        : []),
      {
        key: 'updateTime'
      },
      {
        type: 'btn',
        right: true,
        btns: [{
          title: '发布管理',
          click: (item:any) => {
            context.publish(item)
          }
        },
        {
          title: '查看',
          click: (item:any) => {
            context.router.navigate(['/', 'application', 'content', item.data.id, 'message'])
          }
        },
        {
          title: '删除',
          click: (item:any) => {
            context.deleteDataModal(item.data)
          },
          disabledFn: (data:any, item:any) => {
            return !item.data.isDelete || context.nzDisabled
          }
        }
        ]
      }
    ]
  }

  createApplicationPublicTbody (component:ApplicationPublishComponent):TBODY_TYPE[] {
    const tbody:EO_TBODY_TYPE[] = [
      {
        type: 'checkbox',
        click: () => {
          component.clickData()
        },
        disabledFn: () => {
          return component.nzDisabled
        }
      },
      {
        key: 'title'
      },
      { key: 'status', title: component.clusterStatusTpl },
      { key: 'operator' },
      { key: 'updateTime' }
    ]
    return tbody
  }

  publishAppModal (data:{name:string, id:string, desc?:string}, component?:ApplicationManagementListComponent, returnToSdk?:Function) {
    this.modalRef = this.modalService.create({
      nzTitle: `${data.name}发布管理`,
      nzWidth: MODAL_NORMAL_SIZE,
      nzContent: ApplicationPublishComponent,
      nzComponentParams: {
        name: data.name,
        id: data.id,
        desc: data?.desc,
        closeModal: () => { this.modalRef?.close() },
        nzDisabled: component?.nzDisabled,
        returnToSdk
      },
      nzOnCancel: () => {
        returnToSdk && returnToSdk({ data: { closeModal: true } })
      },
      nzFooter: [{
        label: '取消',
        type: 'default',
        onClick: () => {
          this.modalRef?.close()
          returnToSdk && returnToSdk({ data: { closeModal: true } })
        }
      },
      {
        label: '下线',
        danger: true,
        onClick: (context:ApplicationPublishComponent) => {
          return new Promise((resolve, reject) => {
            context.offline().subscribe((resp) => {
              if (resp) {
                this.modalRef?.close()
                resolve(true)
                component?.getTableData()
              } else {
                reject(new Error())
              }
            })
          })
        },
        disabled: () => {
          return !!component?.nzDisabled
        }
      },
      {
        label: '上线',
        type: 'primary',
        onClick: (context:ApplicationPublishComponent) => {
          return new Promise((resolve, reject) => {
            context.online().subscribe((resp) => {
              if (resp) {
                resolve(true)
                this.modalRef?.close()
                component?.getTableData()
              } else {
                reject(new Error())
              }
            })
          })
        },
        disabled: () => {
          return !!component?.nzDisabled
        }
      }]
    })
  }
}

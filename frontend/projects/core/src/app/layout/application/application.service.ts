import { Injectable } from '@angular/core'
import { ApplicationManagementListComponent } from './list/list.component'
import { EO_TBODY_TYPE } from 'projects/eo-ng-apinto-table/src/public-api'
import { THEAD_TYPE } from 'eo-ng-table'
import { ApplicationData } from './types/types'
import { ApiService } from '../../service/api.service'

@Injectable({
  providedIn: 'root'
})
export class EoNgApplicationService {
  appName:string = ''
  appDesc:string = ''

  appData:ApplicationData|null = null

  constructor (private api:ApiService) {}

  getApplicationData (appId:string) {
    this.appData = null
    this.api.get('application', { appId }).subscribe((resp:{code:number, data:{application:ApplicationData}, msg:string}) => {
      if (resp.code === 0) {
        this.appData = resp.data.application
        this.appName = this.appData.name
        this.appDesc = this.appData.desc
      }
    })
  }

  clearData () {
    this.appName = ''
    this.appDesc = ''
    this.appData = null
  }

  createApiListThead (context:ApplicationManagementListComponent, publishList?:Array<any>):THEAD_TYPE[] {
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
          return { title: `状态：${p.title}`, tooltip: `状态：${p.title}`, titleString: `状态：${p.title}` }
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

  createApiListTbody (context:ApplicationManagementListComponent, publishList?:Array<any>):EO_TBODY_TYPE[] {
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
          title: '上线管理',
          click: (item:any) => {
            context.publish(item.data.uuid)
          }
        },
        {
          title: '查看',
          click: (item:any) => {
            context.router.navigate(['/', 'router', 'api', item.data.scheme === 'websocket' ? 'message-ws' : 'message', item.data.uuid])
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
}

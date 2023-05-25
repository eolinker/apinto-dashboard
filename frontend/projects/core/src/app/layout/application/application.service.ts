import { Injectable } from '@angular/core'
import { ApplicationManagementListComponent } from './list/list.component'
import { EO_TBODY_TYPE } from 'projects/eo-ng-apinto-table/src/public-api'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { ApplicationData } from './types/types'
import { ApiService } from '../../service/api.service'
import { ApplicationPublishComponent } from './publish/publish.component'
import { FilterOpts } from '../../constant/conf'

@Injectable({
  providedIn: 'root'
})
export class EoNgApplicationService {
  appName:string = ''
  appDesc:string = ''

  appData:ApplicationData|null = null
  loading:boolean = true
  constructor (private api:ApiService) {}

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
          title: '上线管理',
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
}

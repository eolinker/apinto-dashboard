import { Injectable } from '@angular/core'
import { IntelligentPluginListComponent } from './list/list.component'
import { EO_TBODY_TYPE } from 'projects/eo-ng-apinto-table/src/public-api'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { IntelligentPluginPublishComponent } from './publish/publish.component'

@Injectable({
  providedIn: 'root'
})
export class IntelligentPluginService {
  createTbody (component:IntelligentPluginListComponent, key?:string):TBODY_TYPE[] {
    const btnConfig:EO_TBODY_TYPE[] = [{

      type: 'btn',
      right: true,
      showFn: () => {
        return !component.tableLoading
      },
      btns: [{
        title: '上线管理',
        click: (item:any) => {
          component.publish(item)
        }
      },
      {
        title: '查看',
        click: (item:any) => {
          component.editData(item)
        }
      }, {
        title: '删除',
        click: (item:any) => {
          component.deleteDataModal(item.data)
        },
        disabledFn: () => {
          return component.nzDisabled
        }
      }
      ]
    },
    {

      type: 'btn',
      right: true,
      showFn: () => {
        return component.tableLoading
      },
      btns: [{
        title: '上线管理',
        click: (item:any) => {
          component.publish(item)
        }
      },
      {
        title: '查看',
        click: (item:any) => {
          component.editData(item)
        }
      }
      ]
    }
    ]
    const tbody:EO_TBODY_TYPE[] = [
      {
        key: 'name',
        left: true
      },
      { key: 'id' },
      { key: 'desc' },
      { key: 'status' },
      { key: 'operator' },
      { key: 'update_time' },
      ...btnConfig
    ]
    return key === 'btn' ? btnConfig : tbody
  }

  createPluginThead (component:IntelligentPluginPublishComponent):THEAD_TYPE[] {
    console.log(component.nzDisabled)
    const thead:THEAD_TYPE[] =
    [{
      type: 'checkbox',
      click: () => {
        component.checkAll()
      },
      disabled: component.nzDisabled
    },
    { title: '集群' },
    { title: '状态' },
    { title: '更新人' },
    { title: '更新时间' }]
    return thead
  }

  createPluginTbody (component:IntelligentPluginPublishComponent):TBODY_TYPE[] {
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
        key: 'name'
      },
      { key: 'status' },
      { key: 'updater' },
      { key: 'update_time' }
    ]
    return tbody
  }
}

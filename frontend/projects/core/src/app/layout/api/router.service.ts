import { Injectable } from '@angular/core'
import { EO_TBODY_TYPE } from 'projects/eo-ng-apinto-table/src/public-api'
import { ApiPluginTemplateListComponent } from './plugin/list/list.component'
import { THEAD_TYPE } from 'eo-ng-table'
import { ApiManagementListComponent } from './api-list/list/list.component'
import { ApiPublishComponent } from './api-list/publish/single/publish.component'

@Injectable({
  providedIn: 'root'
})
export class RouterService {
  createApiListThead (context:ApiManagementListComponent, publishList?:Array<any>):THEAD_TYPE[] {
    return [
      {
        type: 'checkbox',
        resizeable: false,
        click: (item:any) => {
          context.changeApisSet(item, 'all')
        },
        showFn: () => {
          return !context.nzDisabled
        }
      },
      {
        title: 'API名称'
      },
      {
        title: '协议'
      },
      {
        title: '方法',
        width: 140,
        resizeable: false
      },
      {
        title: '请求路径'
      },
      ...(publishList?.length
        ? publishList.map((p) => {
          return { title: `状态：${p.name}`, tooltip: `状态：${p.name}`, titleString: `状态：${p.name}` }
        })
        : []),
      {
        title: '来源',
        filterMultiple: true,
        filterOpts: [{
          text: '自建',
          value: 'build'
        },
        {
          text: '导入',
          value: 'import'
        }
        ],
        filterFn: () => {
          return true
        }
      },
      {
        title: '更新时间'
      },
      {
        title: '操作',
        right: true
      }
    ]
  }

  createApiListTbody (context:ApiManagementListComponent, publishList?:Array<any>):EO_TBODY_TYPE[] {
    return [
      {
        key: 'checked',
        type: 'checkbox',
        click: (item:any) => {
          context.changeApisSet(item)
        },
        showFn: () => {
          return !context.nzDisabled
        }
      },
      {
        key: 'name',
        copy: true
      },
      {
        key: 'scheme'
      },
      {
        key: 'method',
        title: context.methodTpl
      },
      {
        key: 'requestPath',
        copy: true
      },
      ...(publishList?.length
        ? publishList.map((p) => {
          return { key: `cluster_${p.name}`, title: context.clusterStatusTpl }
        })
        : []),
      {
        key: 'source'
      },
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
            context.deleteApiModal(item.data)
          },
          disabledFn: (data:any, item:any) => {
            return !item.data.isDelete || context.nzDisabled
          }
        }
        ]
      }
    ]
  }

  createPluginTemplateTbody (context:ApiPluginTemplateListComponent):EO_TBODY_TYPE[] {
    return [
      {
        key: 'name',
        copy: true
      },
      {
        key: 'desc'
      },
      {
        key: 'createTime'
      },
      {
        key: 'updateTime'
      },
      {
        type: 'btn',
        right: true,
        btns: [
          {
            title: '查看',
            click: (item:any) => {
              context.router.navigate(['/', 'router', 'plugin-template', 'content', item.data.uuid])
            }
          },
          {
            title: '删除',
            disabledFn: (item:any) => {
              return context.nzDisabled || !item.isDelete
            },
            click: (item:any) => {
              context.delete(item)
            }
          }
        ]
      }
    ]
  }

  createApiPublishThead (component:ApiPublishComponent):THEAD_TYPE[] {
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

  createApiPublishTbody (component:ApiPublishComponent):EO_TBODY_TYPE[] {
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
      { key: 'status', title: component.clusterStatusTpl },
      { key: 'operator' },
      { key: 'updateTime' }
    ]
    return tbody
  }
}

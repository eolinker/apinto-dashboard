/*
 * @Date: 2023-12-12 18:57:19
 * @LastEditors: maggieyyy
 * @LastEditTime: 2023-12-15 16:41:04
 * @FilePath: \apinto\projects\core\src\app\layout\node-plugin\nodePlugin.service.ts
 */
import { Injectable } from '@angular/core'
import { EO_TBODY_TYPE } from 'projects/eo-ng-apinto-table/src/public-api'
import { DeployPluginListComponent } from './list/list.component'

@Injectable({
  providedIn: 'root'
})
export class NodePluginService {
  createPluginsTbody = (context:DeployPluginListComponent):EO_TBODY_TYPE[] => {
    return [
      {
        type: 'sort'
      },
      { title: context.pluginName },
      {
        key: 'title',
        copy: true
      },
      {
        key: 'desc'
      },
      { key: 'updateTime' },
      {
        type: 'btn',
        right: true,
        btns: [
          {
            title: '编辑',
            click: (item: any) => {
              context.router.navigate(['/', 'deploy', 'plugin', 'message', item.data.name])
            }
          },
          {
            title: '删除',
            disabledFn: (item:any) => {
              return !item.isDelete || context.nzDisabled
            },
            click: (item:any) => {
              context.delete(item.data)
            }
          }
        ]
      }
    ]
  }
}

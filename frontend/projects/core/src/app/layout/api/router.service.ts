import { Injectable } from '@angular/core'
import { TBODY_TYPE } from 'eo-ng-table'
import { ApiPluginTemplateListComponent } from './plugin/list/list.component'

@Injectable({
  providedIn: 'root'
})
export class RouterService {
  constructor () { }

  createPluginTemplateTbody (context:ApiPluginTemplateListComponent):TBODY_TYPE[] {
    return [
      {
        key: 'name'
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
              context.router.navigate(['/', 'router', 'plugin', 'content', item.data.uuid])
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
}

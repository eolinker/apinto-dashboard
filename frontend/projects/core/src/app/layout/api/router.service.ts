import { Injectable } from '@angular/core'
import { EO_TBODY_TYPE } from 'projects/eo-ng-apinto-table/src/public-api'
import { ApiPluginTemplateListComponent } from './plugin/list/list.component'

@Injectable({
  providedIn: 'root'
})
export class RouterService {
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
}

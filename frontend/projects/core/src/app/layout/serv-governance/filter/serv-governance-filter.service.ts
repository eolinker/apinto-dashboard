import { Injectable } from '@angular/core'
import { TBODY_TYPE } from 'eo-ng-table'
import { MODAL_SMALL_SIZE } from '../../../constant/app.config'
import { FilterTableComponent } from './table/table.component'

@Injectable({
  providedIn: 'root'
})
export class ServGovernanceFilterService {
  constructor () { }

  createFilterTbody (context:FilterTableComponent):TBODY_TYPE[] {
    return [
      { key: 'title' },
      { key: 'label' },
      {
        type: 'btn',
        right: true,
        btns: [
          {
            title: '配置',
            click: (item: any) => {
              context.openDrawer('editFilter', item.data)
            },
            disabledFn: () => {
              return context.nzDisabled
            }
          },
          {
            title: '删除',
            click: (item: any) => {
              context.modalService.create({
                nzTitle: '删除',
                nzContent: '该数据删除后将无法找回，请确认是否删除？',
                nzClosable: true,
                nzClassName: 'delete-modal',
                nzOkDanger: true,
                nzWidth: MODAL_SMALL_SIZE,
                nzOnOk: () => {
                  context.filterDelete(context, item.data)
                }
              })
            },
            disabledFn: () => {
              return context.nzDisabled
            }
          }
        ]
      }
    ]
  }
}

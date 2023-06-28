import { Injectable } from '@angular/core'
import { EoIntelligentPluginListComponent } from './list/list.component'
import { EO_TBODY_TYPE } from 'projects/eo-ng-apinto-table/src/public-api'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { EoIntelligentPluginPublishComponent } from './publish/publish.component'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { EoNgFeedbackModalService } from 'eo-ng-feedback'
import { MODAL_NORMAL_SIZE } from '../../constant/app.config'
import { IntelligentPluginListComponent } from '../../layout/intelligent-plugin/list/list.component'

@Injectable({
  providedIn: 'root'
})
export class EoIntelligentPluginService {
  constructor (private modalService:EoNgFeedbackModalService) {}
  createTbody (component:EoIntelligentPluginListComponent, key?:string):TBODY_TYPE[] {
    const btnConfig:EO_TBODY_TYPE[] = [{
      type: 'btn',
      right: true,
      showFn: () => {
        return !component.tableLoading
      },
      btns: [{
        title: '发布管理',
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
        title: '发布管理',
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

  createPluginThead (component:EoIntelligentPluginPublishComponent):THEAD_TYPE[] {
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

  createPluginTbody (component:EoIntelligentPluginPublishComponent):TBODY_TYPE[] {
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
      { key: 'updater' },
      { key: 'update_time' }
    ]
    return tbody
  }

  modalRef:NzModalRef|undefined

  publishPluginModal (moduleName:string, data:{name:string, id:string, desc?:string}, component?:IntelligentPluginListComponent, returnToSdk?:Function) {
    this.modalRef = this.modalService.create({
      nzTitle: `${data.name}发布管理`,
      nzWidth: MODAL_NORMAL_SIZE,
      nzContent: EoIntelligentPluginPublishComponent,
      nzComponentParams: {
        name: data.name,
        id: data.id,
        desc: data.desc,
        moduleName: moduleName,
        closeModal: () => {
          this.modalRef?.close()
          component?.getTableData()
        },
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
        onClick: (context:EoIntelligentPluginPublishComponent) => {
          context.offline()
        },
        disabled: () => {
          return !!component?.nzDisabled
        }
      },
      {
        label: '上线',
        type: 'primary',
        onClick: (context:EoIntelligentPluginPublishComponent) => {
          context.online()
        },
        disabled: () => {
          return !!component?.nzDisabled
        }
      }]
    })
  }
}

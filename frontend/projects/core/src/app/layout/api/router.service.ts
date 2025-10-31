import { Injectable } from '@angular/core'
import { EO_TBODY_TYPE } from 'projects/eo-ng-apinto-table/src/public-api'
import { THEAD_TYPE } from 'eo-ng-table'
import { ApiManagementListComponent } from './api-list/list/list.component'
import { ApiPublishComponent } from './api-list/publish/single/publish.component'
import { FilterOpts } from '../../constant/conf'
import { MODAL_NORMAL_SIZE, MODAL_SMALL_SIZE } from '../../constant/app.config'
import { EoNgFeedbackModalService } from 'eo-ng-feedback'
import { ModalOptions, NzModalRef } from 'ng-zorro-antd/modal'
import { ApiBatchPublishResultComponent } from './api-list/publish/batch/result.component'
import { apiBatchOnlineVerifyTableBody, apiBatchOnlineVerifyTableHeadName, apiBatchPublishResultTableBody, apiBatchPublishResultTableHeadName } from './types/conf'
import { ApiManagementEditGroupComponent } from './api-list/group/edit-group/edit-group.component'
import { ApiManagementGroupComponent } from './api-list/group/group.component'
import { EmptyHttpResponse } from '../../constant/type'
import { ApiBatchPublishChooseClusterComponent } from './api-list/publish/batch/choose-cluster.component'

@Injectable({
  providedIn: 'root'
})
export class RouterService {
  groupModal: NzModalRef |undefined
  constructor (
    private modalService:EoNgFeedbackModalService) {}

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
      {
        title: '拦截请求',
        width: 80
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
        title: '来源',
        filterMultiple: true,
        filterOpts: context.sourcesList.length > 0
          ? [...context.sourcesList]
          : [{
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
      {
        key: 'isDisable',
        title: context.clusterStatusTpl
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
          title: '发布管理',
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

  createApiPublishThead (component:ApiPublishComponent):THEAD_TYPE[] {
    const thead:THEAD_TYPE[] =
    [{
      type: 'checkbox',
      click: () => {
        component.checkSelectedCluster()
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
          component.checkSelectedCluster()
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

  modalRef:NzModalRef|undefined

  publishApiModal (uuid:string, component?:ApiManagementListComponent, returnToSdk?:Function) {
    this.modalRef = this.modalService.create({
      nzTitle: '发布管理',
      nzWidth: MODAL_NORMAL_SIZE,
      nzContent: ApiPublishComponent,
      nzComponentParams: {
        apiUuid: uuid,
        closeModal: () => { this.modalRef?.close() },
        getApisData: () => { component?.getApisData() },
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
        onClick: (context:ApiPublishComponent) => {
          context.offline()
        },
        disabled: () => {
          return component?.nzDisabled || false
        }
      },
      {
        label: '上线',
        type: 'primary',
        onClick: (context:ApiPublishComponent) => {
          context.online()
        },
        disabled: () => {
          return component?.nzDisabled || false
        }
      }]
    })
  }

  batchPublishApiModal (type:'online'|'offline', data:{uuids:Array<string>, clusters?:Array<string>}, returnToSdk?:Function, component?:ApiManagementListComponent) {
    const batchPublishApiConfig:ModalOptions =
    {
      nzTitle: `批量${type === 'online' ? '上' : '下'}线`,
      nzWidth: MODAL_SMALL_SIZE,
      nzContent: ApiBatchPublishChooseClusterComponent,
      nzComponentParams: {
        type,
        apisSet: new Set(data.uuids),
        clustersSet: new Set(data.clusters || []),
        returnToSdk
      },
      nzOnCancel: () => {
        returnToSdk && returnToSdk({ data: { closeModal: true } })
      },
      nzFooter: [
        {
          label: '取消',
          type: 'default',
          onClick: () => {
            this.modalRef?.close()
            returnToSdk && returnToSdk({ data: { closeModal: true } })
          }
        },
        {
          label: type === 'online' ? '下一步' : '提交',
          type: 'primary',
          loading: (context:ApiBatchPublishChooseClusterComponent) => {
            return context.loading
          },
          onClick: (context:ApiBatchPublishChooseClusterComponent) => {
            this.modalRef?.close()
            this.batchPublishApiResModal(type,
              { uuids: [...context.apisSet], clusters: [...context.clustersSet] },
              () => { this.batchPublishApiModal(type, { ...data, clusters: [...context.clustersSet] }, returnToSdk, component) },
              (resp:any) => { returnToSdk && returnToSdk(resp) },
              component
            )
          },
          disabled: (context:ApiBatchPublishChooseClusterComponent) => {
            return context.clustersSet.size === 0
          }
        }
      ]
    }
    this.modalRef = this.modalService.create(batchPublishApiConfig)
  }

  batchPublishApiResModal (type:'online'|'offline', data:{uuids:Array<string>, clusters:Array<string>}, chooseCluster?: Function, returnToSdk?: Function, component?:ApiManagementListComponent) {
    this.modalRef?.close()
    const checkModalConfig:ModalOptions = {
      nzTitle: '检测结果',
      nzWidth: MODAL_NORMAL_SIZE,
      nzContent: ApiBatchPublishResultComponent,
      nzComponentParams: {
        publishType: 'online',
        stepType: 'check',
        apisSet: new Set(data.uuids),
        clustersSet: new Set(data.clusters),
        chooseCluster,
        renewApiList: () => {
          component?.getApisData()
          component?.apisSet.clear()
        },
        closeModal: () => { this.modalRef?.close() },
        batchPublishTableBody: [...apiBatchOnlineVerifyTableBody],
        batchPublishTableHeadName: [...apiBatchOnlineVerifyTableHeadName],
        onlineApisModal: (context?:ApiBatchPublishResultComponent) => {
          this.modalRef?.close()
          this.modalRef = this.modalService.create(this.getResModalConfig(type, data, returnToSdk, component, context))
          returnToSdk && returnToSdk({ data: { finishPublish: true } })
        }
      },
      nzOnCancel: () => {
        returnToSdk && returnToSdk({ data: { closeModal: true } })
      },
      nzFooter: [{
        label: '重新检测',
        type: 'primary',
        loading: (context:ApiBatchPublishResultComponent) => {
          return context.loading
        },
        show: (context?:ApiBatchPublishResultComponent) => (
          context?.stepType === 'check'
        ),
        onClick: (context:ApiBatchPublishResultComponent) => {
          context.onlineApisCheck()
        }
      },
      {
        label: '上一步',
        loading: (context:ApiBatchPublishResultComponent) => {
          return context.loading
        },
        show: (context?:ApiBatchPublishResultComponent) => (
          context?.stepType === 'check' && !!context?.chooseCluster
        ),
        onClick: (context:ApiBatchPublishResultComponent) => {
          this.modalRef?.close()
          context.chooseCluster && context.chooseCluster()
        }
      },
      {
        label: '批量上线',
        type: 'primary',
        loading: (context:ApiBatchPublishResultComponent) => {
          return context.loading
        },
        disabled: (context:ApiBatchPublishResultComponent) => {
          return !context.onlineToken
        },
        show: (context?:ApiBatchPublishResultComponent) => (
          context?.stepType === 'check'
        ),
        onClick: (context:ApiBatchPublishResultComponent) => {
          this.modalRef?.close()
          this.modalRef = this.modalService.create(this.getResModalConfig(type, data, returnToSdk, component, context))
          returnToSdk && returnToSdk({ data: { finishPublish: true } })
        }
      }]

    }

    this.modalRef = this.modalService.create(type === 'online' ? checkModalConfig : this.getResModalConfig(type, data, returnToSdk, component))
    if (type === 'offline') {
      returnToSdk && returnToSdk({ data: { finishPublish: true } })
    }
  }

  getResModalConfig:(...args:any)=>ModalOptions = (type:'online'|'offline', data:{uuids:Array<string>, clusters:Array<string>}, returnToSdk?:Function, component?:ApiManagementListComponent, context?:ApiBatchPublishResultComponent) => {
    return {
      nzTitle: type === 'online' ? '批量上线结果' : '批量下线结果',
      nzWidth: MODAL_NORMAL_SIZE,
      nzContent: ApiBatchPublishResultComponent,
      nzComponentParams: {
        publishType: type,
        stepType: 'result',
        apisSet: new Set(data.uuids),
        clustersSet: new Set(data.clusters),
        closeModal: () => { this.modalRef?.close() },
        batchPublishTableBody: [...apiBatchPublishResultTableBody],
        batchPublishTableHeadName: [...apiBatchPublishResultTableHeadName],
        onlineToken: context?.onlineToken,
        renewApiList: () => {
          component?.getApisData()
          component?.apisSet.clear()
        },
        returnToSdk: returnToSdk
      },
      nzFooter: [
        {
          label: '关闭',
          loading: (context:ApiBatchPublishResultComponent) => {
            return context.loading
          },
          onClick: () => {
            this.modalRef?.close()
          }
        }
      ]
    }
  }

  addOrEditGroupModal (type:'add'|'edit', uuid?:string, name?:string, component?:ApiManagementGroupComponent) {
    return new Promise((resolve) => {
      const title:string = type === 'add' ? (uuid !== 'root' ? '添加子分组' : '添加分组') : '编辑分组'

      this.groupModal = this.modalService.create({
        nzTitle: title,
        nzContent: ApiManagementEditGroupComponent,
        nzWidth: MODAL_SMALL_SIZE,
        nzComponentParams: {
          ...(type === 'edit' ? { groupName: name } : {}),
          uuid: uuid,
          type: type,
          closeModal: (resp:EmptyHttpResponse) => { resolve(resp); this.groupModal?.close(); component?.getMenuList() },
          showUuid: type === 'add' && uuid !== 'root'
        },
        nzClosable: true,
        nzCancelText: '取消',
        nzOkText: '确定',
        nzOnOk: (context:ApiManagementEditGroupComponent) => {
          if (type === 'add') {
            component && (component.editParentUuid = uuid === 'root' ? '' : uuid || '')
            context.addGroup(uuid!)
          } else {
            component && (component.editUuid = uuid!)
            context.editGroup(uuid!)
          }
          return false
        },
        nzOnCancel: () => {
          resolve({ closeModal: true })
        }
      })
    })
  }
}

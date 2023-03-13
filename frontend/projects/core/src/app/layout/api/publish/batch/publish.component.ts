/* eslint-disable no-useless-constructor */
/* eslint-disable no-undef */
import { Component, EventEmitter, Input, OnInit, Output, TemplateRef, ViewChild } from '@angular/core'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { MODAL_NORMAL_SIZE, MODAL_SMALL_SIZE } from 'projects/core/src/app/constant/app.config'
import { ApiService } from 'projects/core/src/app/service/api.service'

@Component({
  selector: 'eo-ng-api-batch-publish',
  templateUrl: './publish.component.html',
  styles: [
    `label{
      height: 26px;
      line-height: 26px;
      margin-bottom: 16px;
      display: block;
      margin-top: 2px;
    }`
  ]
})
export class ApiBatchPublishComponent implements OnInit {
  @ViewChild('disabledCheckboxTpl', { read: TemplateRef, static: true }) disabledCheckboxTpl: TemplateRef<any> | undefined
  @ViewChild('startContentTpl', { read: TemplateRef, static: true }) startContentTpl: TemplateRef<any> | undefined
  @ViewChild('startFooterTpl', { read: TemplateRef, static: true }) startFooterTpl: TemplateRef<any> | undefined
  @ViewChild('resultContentTpl', { read: TemplateRef, static: true }) resultContentTpl: TemplateRef<any> | undefined
  @ViewChild('resultFooterTpl', { read: TemplateRef, static: true }) resultFooterTpl: TemplateRef<any> | undefined
  @Input()
  get apisSet () {
    return this._apisSet
  }

  set apisSet (val) {
    this._apisSet = val
    this.apisSetChange.emit(this._apisSet)
  }

  @Output() apisSetChange = new EventEmitter()
  @Output() flashList:EventEmitter<any> = new EventEmitter()

  _apisSet:Set<string> = new Set()
  clustersList:Array<{label:string, value:any, disabled:boolean, checked:boolean, template?:any}>=[]
  clustersSet:Set<string> = new Set()
  onlineResultList:Array<{service:string, cluster:string, status:boolean, statusString:string, result:string, solution:{params:any, name:string}, name:string}>=[]
  onlineToken:string = ''
  onlineModalRef:NzModalRef | undefined
  offlineModalRef:NzModalRef | undefined
  resultList:Array<{api:string, cluster:string, status:boolean, statusString:string, result:string}>=[]

  onlineResultTableHeadName:Array<object> = [
    {
      title: '上游服务名称'
    },
    { title: '集群名称' },
    { title: '状态' },
    { title: '失败原因' },
    {
      title: '操作',
      width: 94,
      right: true
    }
  ]

  onlineResultTableBody:Array<any> = [
    {
      key: 'service',
      styleFn: (item:any) => {
        if (!item.status) {
          return 'color:#ff3b30'
        }
        return ''
      }
    },
    {
      key: 'cluster',
      styleFn: (item:any) => {
        if (!item.status) {
          return 'color:#ff3b30'
        }
        return ''
      }
    },
    {
      key: 'statusString',
      styleFn: (item:any) => {
        if (!item.status) {
          return 'color:#ff3b30'
        }
        return ''
      }
    },
    {
      key: 'result',
      styleFn: (item:any) => {
        if (!item.status) {
          return 'color:#ff3b30'
        }
        return ''
      }
    },
    {
      type: 'btn',
      right: true,
      showFn: (item:any) => {
        return item.solution.name
      },
      btns: [{
        title: '解决方案',
        click: (item:any) => {
          let routerS:string = '/' + item.data.solution.name + '?'
          if (Object.keys(item.data.solution.params).length > 0) {
            for (const index in Object.keys(item.data.solution.params)) {
              routerS = routerS + Object.keys(item.data.solution.params)[index] + '=' + item.data.solution.params[Object.keys(item.data.solution.params)[index]] + '&'
            }
          }
          window.open(routerS, '')
        },
        type: 'text'
      }
      ]
    },
    {
      type: 'btn',
      right: true,
      showFn: (item:any) => {
        return !item.solution.name
      }
    }
  ]

  resultTableHeadName:Array<object> = [
    {
      title: 'API名称'
    },
    { title: '集群名称' },
    { title: '状态' },
    { title: '失败原因' }
  ]

  resultTableBody:Array<any> = [
    {
      key: 'api',
      styleFn: (item:any) => {
        if (!item.status) {
          return 'color:#ff3b30'
        }
        return ''
      }
    },
    {
      key: 'cluster',
      styleFn: (item:any) => {
        if (!item.status) {
          return 'color:#ff3b30'
        }
        return ''
      }
    },
    {
      key: 'statusString',
      styleFn: (item:any) => {
        if (!item.status) {
          return 'color:#ff3b30'
        }
        return ''
      }
    },
    {
      key: 'result',
      styleFn: (item:any) => {
        if (!item.status) {
          return 'color:#ff3b30'
        }
        return ''
      }
    }
  ]

  constructor (
  private modalService:EoNgFeedbackModalService, private api:ApiService, private message: EoNgFeedbackMessageService) { }

  ngOnInit (): void {
  }

  // 获取集群列表
  getClusterList (type:string):void {
    this.api.get('clusters').subscribe(resp => {
      if (resp.code === 0) {
        this.clustersList = []
        for (const index in resp.data.clusters) {
          this.clustersList.push({ label: `${resp.data.clusters[index].name}_${resp.data.clusters[index].env}`, value: resp.data.clusters[index].name, checked: (resp.data.clusters.length === 1 && resp.data.clusters[index].status !== 'ABNORMAL'), disabled: resp.data.clusters[index].status === 'ABNORMAL', template: resp.data.clusters[index].status === 'ABNORMAL' ? this.disabledCheckboxTpl : '' })
          if (resp.data.clusters.length === 1 && resp.data.clusters[index].status !== 'ABNORMAL') {
            this.clustersSet.add(resp.data.clusters[index].name)
          }
        }
        this.openDrawer(type)
      } else {
        this.message.error(resp.msg || '获取列表数据失败!')
      }
    })
  }

  // 打开弹窗-批量上线\下线 选择集群页
  openDrawer (type:string) {
    switch (type) {
      case 'online': {
        this.onlineModalRef?.close()
        this.onlineModalRef = this.modalService.create({
          nzTitle: '批量上线',
          nzWidth: MODAL_SMALL_SIZE,
          nzContent: this.startContentTpl,
          nzComponentParams: { type: 'online' },
          nzFooter: this.startFooterTpl,
          nzOnCancel: () => {
            this.clustersSet.clear()
            this.clustersList = []
          }
        })
        break
      }
      case 'offline': {
        this.offlineModalRef?.close()
        this.offlineModalRef = this.modalService.create({
          nzTitle: '批量下线',
          nzWidth: MODAL_SMALL_SIZE,
          nzContent: this.startContentTpl,
          nzComponentParams: { type: 'offline' },
          nzFooter: this.startFooterTpl,
          nzOnCancel: () => {
            this.clustersSet.clear()
            this.clustersList = []
          }
        })
      }
    }
  }

  // 批量上线/下线 结果页
  apisOperatorResult (type:string) {
    switch (type) {
      case 'online': {
        this.onlineApisCheck()
        this.onlineModalRef?.close()
        this.onlineModalRef = this.modalService.create({
          nzTitle: '检测结果',
          nzWidth: MODAL_NORMAL_SIZE,
          nzContent: this.resultContentTpl,
          nzComponentParams: { type: 'online', online: true },
          nzFooter: this.resultFooterTpl,
          nzOnCancel: () => {
            this.clustersSet.clear()
            this.clustersList = []
          }
        })
        break
      }
      case 'online-res': {
        this.onlineModalRef?.close()
        this.onlineModalRef = this.modalService.create({
          nzTitle: '批量上线结果',
          nzWidth: MODAL_NORMAL_SIZE,
          nzContent: this.resultContentTpl,
          nzComponentParams: { type: 'res', online: false },
          nzFooter: this.resultFooterTpl,
          nzOnCancel: () => {
            this.clustersSet.clear()
            this.clustersList = []
          }
        })
        break
      }
      case 'offline': {
        this.offlineApis()
        this.offlineModalRef?.close()
        this.offlineModalRef = this.modalService.create({
          nzTitle: '批量下线结果',
          nzWidth: MODAL_NORMAL_SIZE,
          nzContent: this.resultContentTpl,
          nzComponentParams: { type: 'res' },
          nzFooter: this.resultFooterTpl,
          nzOnCancel: () => {
            this.clustersSet.clear()
            this.clustersList = []
          }
        })
        break
      }
    }
  }

  // 在批量上\下线检测页和结果页中，上\下线成功的表格行字体为绿色，失败的为红色
  trStyleFn (item:any) {
    if (item.status) {
      return 'color:green'
    } else {
      return 'color:red'
    }
  }

  // 关闭弹窗
  closeDrawer (type:string) {
    switch (type) {
      case 'online':
        this.onlineModalRef?.close()
        break
      case 'offline':
        this.offlineModalRef?.close()
        break
      case 'res':
        this.onlineModalRef?.close()
        this.offlineModalRef?.close()
        break
    }
  }

  // 检测批量上线的api
  onlineApisCheck () {
    this.onlineToken = ''
    this.api.post('routers/batch-online/check', { api_uuids: [...this.apisSet], cluster_names: [...this.clustersSet] }).subscribe((resp:any) => {
      if (resp.code === 0) {
        this.onlineResultList = resp.data.list
        for (const index in this.onlineResultList) {
          this.onlineResultList[index].statusString = this.onlineResultList[index].status ? '成功' : '失败'
        }
        this.onlineToken = resp.data.online_token
      } else {
        this.message.error(resp.msg || '获取列表数据失败!')
      }
    })
  }

  // 批量上线api
  onlineApis () {
    this.api.post('routers/batch-online', { online_token: this.onlineToken }).subscribe((resp:any) => {
      if (resp.code === 0) {
        this.apisOperatorResult('online-res')
        this.resultList = resp.data.list
        for (const index in this.resultList) {
          this.resultList[index].statusString = this.resultList[index].status ? '成功' : '失败'
        }
        this.apisSet = new Set()
        this.flashList.emit(true)
      } else {
        this.message.error(resp.msg || '批量上线失败!')
        this.apisOperatorResult('online-res')
        this.resultList = resp.data.list
      }
    })
  }

  // 批量下线api
  offlineApis () {
    this.api.post('routers/batch-offline', { api_uuids: [...this.apisSet], cluster_names: [...this.clustersSet] }).subscribe((resp:any) => {
      if (resp.code === 0) {
        this.flashList.emit(true)
        this.apisSet = new Set()
      } else {
        this.message.error(resp.msg || '批量下线失败!')
      }
      this.resultList = resp.data.list
      for (const index in this.resultList) {
        this.resultList[index].statusString = this.resultList[index].status ? '成功' : '失败'
      }
    })
  }

  // 批量上下线的集群被选中或取消时，clustersSet随之变化
  changeClustersSet (value:any) {
    for (const index in value) {
      if (value[index].checked && !value[index].disabled) {
        this.clustersSet.add(value[index].value)
      } else {
        this.clustersSet.delete(value[index].value)
      }
    }
  }
}

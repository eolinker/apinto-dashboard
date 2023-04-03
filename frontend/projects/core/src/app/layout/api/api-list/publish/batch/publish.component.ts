import { Component, EventEmitter, Input, OnInit, Output, TemplateRef, ViewChild } from '@angular/core'
import { CheckBoxOptionInterface } from 'eo-ng-checkbox'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { MODAL_NORMAL_SIZE, MODAL_SMALL_SIZE } from 'projects/core/src/app/constant/app.config'
import { ClustersData } from 'projects/core/src/app/constant/type'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { APIBatchOnlineVerifyData, APIBatchPublishData } from '../../../types/types'
import { apiBatchOnlineVerifyTableBody, apiBatchOnlineVerifyTableHeadName, apiBatchPublishResultTableBody, apiBatchPublishResultTableHeadName } from '../../../types/conf'

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
  @Output() flashList:EventEmitter<boolean> = new EventEmitter()

  _apisSet:Set<string> = new Set()
  clustersList:CheckBoxOptionInterface[]=[]
  clustersSet:Set<string> = new Set()
  onlineResultList:APIBatchOnlineVerifyData[]=[]
  onlineToken:string = ''
  onlineModalRef:NzModalRef | undefined
  offlineModalRef:NzModalRef | undefined
  resultList:APIBatchPublishData[]=[]
  onlineResultTableHeadName:THEAD_TYPE[]= [...apiBatchOnlineVerifyTableHeadName]
  onlineResultTableBody:TBODY_TYPE[] = [...apiBatchOnlineVerifyTableBody]
  resultTableHeadName:THEAD_TYPE[] = [...apiBatchPublishResultTableHeadName]
  resultTableBody:TBODY_TYPE[] = [...apiBatchPublishResultTableBody]

  constructor (private modalService:EoNgFeedbackModalService, private api:ApiService, private message: EoNgFeedbackMessageService) { }

  ngOnInit (): void {
  }

  // 获取集群列表
  getClusterList (type:'online'|'offline'):void {
    this.api.get('clusters').subscribe((resp:{code:number, data:{clusters:ClustersData[]}, msg:string}) => {
      if (resp.code === 0) {
        this.clustersList = []
        for (const index in resp.data.clusters) {
          this.clustersList.push({ label: `${resp.data.clusters[index].name}_${resp.data.clusters[index].env}`, value: resp.data.clusters[index].name, checked: (resp.data.clusters.length === 1 && resp.data.clusters[index].status !== 'ABNORMAL'), disabled: resp.data.clusters[index].status === 'ABNORMAL', template: resp.data.clusters[index].status === 'ABNORMAL' ? this.disabledCheckboxTpl : undefined })
          if (resp.data.clusters.length === 1 && resp.data.clusters[index].status !== 'ABNORMAL') {
            this.clustersSet.add(resp.data.clusters[index].name)
          }
        }
        this.openDrawer(type)
      }
    })
  }

  // 打开弹窗-批量上线\下线 选择集群页
  openDrawer (type:'online'|'offline') {
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
  trStyleFn (item:APIBatchOnlineVerifyData|APIBatchPublishData) {
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
    this.api.post(
      'routers/batch-online/check',
      { apiUuids: [...this.apisSet], clusterNames: [...this.clustersSet] }
    )
      .subscribe((resp:{code:number, data:{list:APIBatchOnlineVerifyData[], onlineToken:string}, msg:string}) => {
        if (resp.code === 0) {
          this.onlineResultList = resp.data.list
          for (const index in this.onlineResultList) {
            this.onlineResultList[index].statusString = this.onlineResultList[index].status ? '成功' : '失败'
          }
          this.onlineToken = resp.data.onlineToken
        }
      })
  }

  // 批量上线api
  onlineApis () {
    this.api.post(
      'routers/batch-online',
      { onlineToken: this.onlineToken }
    )
      .subscribe((resp:{code:number, data:{list:APIBatchPublishData[]}, msg:string}) => {
        if (resp.code === 0) {
          this.apisOperatorResult('online-res')
          this.resultList = resp.data.list
          for (const index in this.resultList) {
            this.resultList[index].statusString = this.resultList[index].status ? '成功' : '失败'
          }
          this.apisSet = new Set()
          this.flashList.emit(true)
        } else {
          this.apisOperatorResult('online-res')
          this.resultList = resp.data.list
        }
      })
  }

  // 批量下线api
  offlineApis () {
    this.api.post(
      'routers/batch-offline',
      { apiUuids: [...this.apisSet], clusterNames: [...this.clustersSet] }
    )
      .subscribe((resp:{code:number, data:{ list:APIBatchPublishData[] }, msg:string}) => {
        if (resp.code === 0) {
          this.flashList.emit(true)
          this.apisSet = new Set()
        } 
        this.resultList = resp.data.list
        for (const index in this.resultList) {
          this.resultList[index].statusString = this.resultList[index].status ? '成功' : '失败'
        }
      })
  }

  // 批量上下线的集群被选中或取消时，clustersSet随之变化
  changeClustersSet (value:CheckBoxOptionInterface[]) {
    for (const index in value) {
      if (value[index].checked && !value[index].disabled) {
        this.clustersSet.add(value[index].value)
      } else {
        this.clustersSet.delete(value[index].value)
      }
    }
  }
}

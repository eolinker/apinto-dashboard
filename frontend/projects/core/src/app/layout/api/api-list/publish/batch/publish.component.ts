import { Component, EventEmitter, Input, OnInit, Output, TemplateRef, ViewChild } from '@angular/core'
import { CheckBoxOptionInterface } from 'eo-ng-checkbox'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { MODAL_SMALL_SIZE } from 'projects/core/src/app/constant/app.config'
import { ClusterSimpleOption } from 'projects/core/src/app/constant/type'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { APIBatchOnlineVerifyData, APIBatchPublishData } from '../../../types/types'
import { apiBatchOnlineVerifyTableBody, apiBatchOnlineVerifyTableHeadName, apiBatchPublishResultTableBody, apiBatchPublishResultTableHeadName } from '../../../types/conf'
import { RouterService } from '../../../router.service'

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
    }
  `
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

  constructor (private modalService:EoNgFeedbackModalService, private service:RouterService, private api:ApiService, private message: EoNgFeedbackMessageService) { }

  ngOnInit (): void {
  }

  // 获取集群列表
  getClusterList (type:'online'|'offline'):void {
    this.clustersSet.clear()
    this.clustersList = []
    this.api.get('clusters/simple').subscribe((resp:{code:number, data:{clusters:ClusterSimpleOption[]}, msg:string}) => {
      if (resp.code === 0) {
        this.clustersList = []
        for (const index in resp.data.clusters) {
          this.clustersList = resp.data.clusters.map((cluster:ClusterSimpleOption) => {
            return { label: cluster.title, value: cluster.name }
          })
          if (this.clustersList.length === 1) {
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
  apisOperatorResult (type:'online'|'offline') {
    this.service.batchPublishApiResModal(
      type,
      { uuids: [...this.apisSet], clusters: [...this.clustersSet] },
      () => { this.openDrawer(type) },
      () => {}, // 给sdk预留的参数
      this)
    this.onlineModalRef?.close()
    this.offlineModalRef?.close()
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

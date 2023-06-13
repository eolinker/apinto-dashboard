import { Component } from '@angular/core'
import { APIBatchOnlineVerifyData, APIBatchPublishData } from '../../../types/types'
import { THEAD_TYPE, TBODY_TYPE } from 'eo-ng-table'
import { ApiService } from 'projects/core/src/app/service/api.service'

@Component({
  selector: 'eo-ng-batch-publish-result',
  templateUrl: './result.component.html',
  styles: [
  ]
})
export class ApiBatchPublishResultComponent {
  publishType:'online'|'offline' = 'online'
  stepType:'check'|'result' = 'check'
  batchPublishTableList:APIBatchPublishData[] | APIBatchOnlineVerifyData[]=[]
  batchPublishTableBody:TBODY_TYPE[] = []
  batchPublishTableHeadName:THEAD_TYPE[] = []
  onlineToken:string = ''
  apisSet:Set<string> = new Set()
  clustersSet:Set<string> = new Set()
  closeModal:Function|undefined
  chooseCluster:Function|undefined
  renewApiList:Function|undefined
  onlineApisModal:Function|undefined // 当检测结果成功时，直接进入批量上线
  loading:boolean = false
  returnToSdk:Function | undefined // 从sdk调用上下线弹窗时，需要的回调
  constructor (private api:ApiService) { }

  // 在批量上\下线检测页和结果页中，上\下线成功的表格行字体为绿色，失败的为红色
  trStyleFn (item:APIBatchOnlineVerifyData|APIBatchPublishData) {
    if (item.status) {
      return 'color:green'
    } else {
      return 'color:red'
    }
  }

  ngOnInit () {
    if (this.publishType === 'online' && this.stepType === 'check') {
      this.onlineApisCheck()
    } else if (this.publishType === 'online' && this.stepType === 'result') {
      this.onlineApis()
    } else {
      this.offlineApis()
    }
  }

  // 检测批量上线的api
  onlineApisCheck () {
    this.loading = true
    this.onlineToken = ''
    this.api.post(
      'routers/batch-online/check',
      { apiUuids: [...this.apisSet], clusterNames: [...this.clustersSet] }
    )
      .subscribe((resp:{code:number, data:{list:APIBatchOnlineVerifyData[], onlineToken:string}, msg:string}) => {
        this.loading = false
        if (resp.code === 0) {
          this.batchPublishTableList = resp.data.list
          for (const index in this.batchPublishTableList) {
            this.batchPublishTableList[index].statusString = this.batchPublishTableList[index].status ? '成功' : '失败'
          }
          this.onlineToken = resp.data.onlineToken
          if (this.onlineToken) {
            this.onlineApisModal && this.onlineApisModal(this)
          }
        }
        if (this.onlineToken) {
          this.returnToSdk && this.returnToSdk(resp)
        }
      })
  }

  // 批量上线api
  onlineApis () {
    this.loading = true
    this.api.post(
      'routers/batch-online',
      { onlineToken: this.onlineToken }
    )
      .subscribe((resp:{code:number, data:{list:APIBatchPublishData[]}, msg:string}) => {
        this.loading = false
        if (resp.code === 0) {
          // this.apisOperatorResult('online-res')
          this.batchPublishTableList = resp.data.list
          for (const index in this.batchPublishTableList) {
            this.batchPublishTableList[index].statusString = this.batchPublishTableList[index].status ? '成功' : '失败'
          }
          this.apisSet = new Set()
          this.renewApiList && this.renewApiList()
        } else {
          // this.apisOperatorResult('online-res')
          this.batchPublishTableList = resp.data.list
        }
        this.returnToSdk && this.returnToSdk(resp)
      })
  }

  // 批量下线api
  offlineApis () {
    this.loading = true
    this.api.post(
      'routers/batch-offline',
      { apiUuids: [...this.apisSet], clusterNames: [...this.clustersSet] }
    )
      .subscribe((resp:{code:number, data:{ list:APIBatchPublishData[] }, msg:string}) => {
        this.loading = false
        if (resp.code === 0) {
          this.renewApiList && this.renewApiList()
          this.apisSet = new Set()
        }
        this.batchPublishTableList = resp.data.list
        for (const index in this.batchPublishTableList) {
          this.batchPublishTableList[index].statusString = this.batchPublishTableList[index].status ? '成功' : '失败'
        }
        this.returnToSdk && this.returnToSdk(resp)
      })
  }

  getTip () {
    if (this.stepType === 'check') {
      return `正在检测选中 API 是否可以${(this.publishType === 'online' ? '上线' : '下线')} ...`
    } else {
      return `正在${(this.publishType === 'online' ? '上线' : '下线')}选中 API...`
    }
  }
}

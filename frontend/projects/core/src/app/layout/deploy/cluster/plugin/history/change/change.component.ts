import { Component, Input, OnInit, TemplateRef } from '@angular/core'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { SelectOption } from 'eo-ng-select'
import { THEAD_TYPE, TBODY_TYPE } from 'eo-ng-table'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { DeployClusterPluginStatusOptions, DeployClusterPluginChangeHistoryTbody, DeployClusterPluginChangeHistoryThead } from '../../../types/conf'
import { ClusterPluginChangeHistoryItem } from '../../../types/types'

@Component({
  selector: 'eo-ng-deploy-cluster-plugin-history-change',
  templateUrl: './change.component.html',
  styles: [
  ]
})
export class DeployClusterPluginHistoryChangeComponent implements OnInit {
  @Input() publishTypeTpl: TemplateRef<any> | undefined
  // eslint-disable-next-line camelcase
  operateRecordsData:{historys:ClusterPluginChangeHistoryItem[], total:number}=
      {
        historys: [],
        total: 0
      }

  operateRecordTabelHeadName: THEAD_TYPE[] = [...DeployClusterPluginChangeHistoryThead]
  operateRecordTableBody:TBODY_TYPE[]=[...DeployClusterPluginChangeHistoryTbody]

  // 更改历史分页
  // eslint-disable-next-line camelcase
  operateRecordsPage:{pageNum:number, pageSize:number, total:number}={
    pageNum: 1,
    pageSize: 15,
    total: 0
  }

  pageSizeOptions:Array<number>=[15, 20, 50, 100]
  clusterName:string = ''
  statusList:SelectOption[] = [...DeployClusterPluginStatusOptions]
  constructor (
    private message: EoNgFeedbackMessageService,
    private api:ApiService) { }

  ngOnInit (): void {
    this.getOperateRecords()
  }

  ngAfterViewInit () {
    this.operateRecordTableBody[3].title = this.publishTypeTpl
  }

  getOperateRecords () {
    this.api.get('cluster/' + this.clusterName + '/plugin/update-history', { pageNum: this.operateRecordsPage.pageNum, pageSize: this.operateRecordsPage.pageSize })
      .subscribe((resp:{code:number, data:{ histories:ClusterPluginChangeHistoryItem[], total:number}, msg:string}) => {
        if (resp.code === 0) {
          this.operateRecordsData.total = resp.data.total
          this.operateRecordsData.historys = resp.data.histories.map((history:ClusterPluginChangeHistoryItem) => {
            history.oldValue = `状态：${this.getStatusString(history.oldConfig.status)}，配置信息：${history.oldConfig.config}`
            history.newValue = `状态：${this.getStatusString(history.newConfig.status)}，配置信息：${history.newConfig.config}`
            return history
          })
        }
      })
  }

  getStatusString (status:'GLOBAL'|'DISABLE'|'ENABLE') {
    if (!status) {
      return '无'
    }
    return this.statusList.filter((item:SelectOption) => { return item.value === status })[0]?.label
  }
}

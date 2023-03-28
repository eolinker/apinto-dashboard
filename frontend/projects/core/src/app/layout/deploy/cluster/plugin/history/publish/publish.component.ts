import { Component, Input, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { SelectOption } from 'eo-ng-select'
import { THEAD_TYPE, TBODY_TYPE } from 'eo-ng-table'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { DeployClusterPublishRecordThead, DeployClusterPublishRecordTbody, DeployClusterPluginChangeHistoryTbody, DeployClusterPluginChangeHistoryThead, DeployClusterPluginStatusOptions } from '../../../types/conf'
import { ClusterPluginPublishHistoryItem } from '../../../types/types'

@Component({
  selector: 'eo-ng-deploy-cluster-plugin-history-publish',
  templateUrl: './publish.component.html',
  styles: [
  ]
})
export class DeployClusterPluginHistoryPublishComponent implements OnInit {
  @ViewChild('showInnerTableBtnTpl', { read: TemplateRef, static: true }) showInnerTableBtnTpl: TemplateRef<any> | undefined
  @ViewChild('innerTableOldValueTpl', { read: TemplateRef, static: true }) innerTableOldValueTpl: TemplateRef<any> | undefined
  @ViewChild('innerTableNewValueTpl', { read: TemplateRef, static: true }) innerTableNewValueTpl: TemplateRef<any> | undefined
  @Input() publishTypeTpl: TemplateRef<any> | undefined

  publishRecordsData:{historys:ClusterPluginPublishHistoryItem[], total:number}=
      {
        historys: [],
        total: 0
      }

  publishRecordTabelHeadName: THEAD_TYPE[]= [...DeployClusterPublishRecordThead]
  publishRecordTableBody: TBODY_TYPE[]=[...DeployClusterPublishRecordTbody]

  publishRecordDetailsTabelHeadName: THEAD_TYPE[] = [...DeployClusterPluginChangeHistoryThead]
  publishRecordDetailsTableBody: TBODY_TYPE[] = [...DeployClusterPluginChangeHistoryTbody]
  pageSizeOptions:Array<number>=[15, 20, 50, 100]
  statusList:SelectOption[] = [...DeployClusterPluginStatusOptions]

  // 发布历史分页
  publishRecordsPage:{pageNum:number, pageSize:number, total:number}={
    pageNum: 1,
    pageSize: 15,
    total: 0
  }

  clusterName:string = ''
  constructor (private message: EoNgFeedbackMessageService,
    private api:ApiService) {
  }

  ngOnInit (): void {
    this.getPublishRecords()
  }

  ngAfterViewInit () {
    this.publishRecordTableBody[0].title = this.showInnerTableBtnTpl

    this.publishRecordDetailsTableBody[1].title = this.innerTableOldValueTpl
    this.publishRecordDetailsTableBody[2].title = this.innerTableNewValueTpl
    this.publishRecordDetailsTableBody[3].title = this.publishTypeTpl
  }

  getPublishRecords () {
    this.api.get('cluster/' + this.clusterName + '/plugin/publish-history', { pageNum: this.publishRecordsPage.pageNum, pageSize: this.publishRecordsPage.pageSize })
      .subscribe((resp:{code:number, data:{histories:ClusterPluginPublishHistoryItem[], total:number}, msg:string}) => {
        if (resp.code === 0) {
          for (const history of resp.data.histories) {
            history.isExpand = false
          }
          this.publishRecordsData.historys = resp.data.histories
          this.publishRecordsData.total = resp.data.total
        }
      })
  }

  expandChange () {
    this.publishRecordsData.historys = [...this.publishRecordsData.historys]
  }

  publishRTableClick = (item:any) => {
    item.data.isExpand = !item.data.isExpand
  }

  getStatusString (status:'GLOBAL'|'DISABLE'|'ENABLE') {
    return this.statusList.filter((item:SelectOption) => { return item.value === status })[0]?.label || '无'
  }

  transferToJson (str:string) {
    return str.replace(/(,)/g, ',\n').replace(/(，)/g, '，\n')
  }
}

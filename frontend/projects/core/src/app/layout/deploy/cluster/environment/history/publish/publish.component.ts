import { Component, Input, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { ApiService } from 'projects/core/src/app/service/api.service'

@Component({
  selector: 'eo-ng-deploy-cluster-environment-history-publish',
  template: `
  <div class="drawer-list-content drawer-table" style="margin-top: 0px">
    <eo-ng-apinto-table
      class="drawer-table"
      [nzTbody]="publishRecordTableBody"
      [nzThead]="publishRecordTabelHeadName"
      [(nzData)]="publishRecordsData.historys"
      [nzTrBottomTmp]="nzTrBottomTmp"
      [nzTrClick]="publishRTableClick"
      [nzNoScroll]="true"
      [nzShowPagination]="true"
      [nzFrontPagination]="false"
      [nzTotal]="publishRecordsPage.total"
      [(nzPageIndex)]="publishRecordsPage.page_num"
      [(nzPageSize)]="publishRecordsPage.page_size"
      [nzPageSizeOptions]="pageSizeOptions"
      (nzPageIndexChange)="getPublishRecords()"
      (nzPageSizeChange)="getPublishRecords()"
    >
    </eo-ng-apinto-table>
  </div>
  <ng-template #nzTrBottomTmp let-item="item" let-apis="apis" let-index="index">
    <tr *ngIf="item.data.isExpand">
      <td style="border-right: none; background-color: #e8e8e8"></td>
      <td
        colspan="3"
        style="
          background-color: pink;
          border: none;
          padding: 0;
          margin-top: -2px;
        "
      >
        <div style="width: 100%">
          <eo-ng-apinto-table
            class="floatR innerTable"
            [nzTbody]="publishRecordDetailsTableBody"
            [nzThead]="publishRecordDetailsTabelHeadName"
            [(nzData)]="item.data.details"
            nzTableLayout="fixed"
            [nzScrollY]="9999"
          >
          </eo-ng-apinto-table>
        </div>
      </td>
    </tr>
  </ng-template>


<ng-template #showInnerTableBtnTpl let-item="item">
<button eo-ng-row-expand-button [expand]="item.isExpand"></button>
</ng-template>
  `,
  styles: [
  ]
})
export class DeployClusterEnvironmentHistoryPublishComponent implements OnInit {
  @ViewChild('showInnerTableBtnTpl', { read: TemplateRef, static: true }) showInnerTableBtnTpl: TemplateRef<any> | undefined
  @Input() publishTypeTpl: TemplateRef<any> | undefined

  // eslint-disable-next-line camelcase
  publishRecordsData:{historys:Array<{id:number, name:string, create_time:string, opt_type:string, operator:string, detailShow:boolean, detail:Array<{key:string, old_value:string, new_value:string, opt_type:string, create_time:string}>}>, total:number}=
      {
        historys: [],
        total: 0
      }

  publishRecordTabelHeadName: Array<object> = [
    { width: 45 },
    { title: '版本名称', resizeable: true },
    { title: '发布者', resizeable: true },
    { title: '发布时间' }
  ]

  publishRecordTableBody: Array<any> =[
    {
      key: ''
    },
    { key: 'name' },
    { key: 'operator' },
    { key: 'create_time' }
  ]

  publishRecordDetailsTabelHeadName: Array<object> = [
    { title: 'KEY', resizeable: true },
    { title: 'OLD VALUE', resizeable: true },
    { title: 'NEW VALUE', resizeable: true },
    { title: '类型', resizeable: true },
    { title: '操作时间' }
  ]

  publishRecordDetailsTableBody: Array<any> =[
    { key: 'key' },
    { key: 'old_value' },
    { key: 'new_value' },
    { key: 'opt_type' },
    { key: 'create_time' }
  ]

  pageSizeOptions:Array<number>=[15, 20, 50, 100]

  // 发布历史分页
  // eslint-disable-next-line camelcase
  publishRecordsPage:{page_num:number, page_size:number, total:number}={
    page_num: 1,
    page_size: 15,
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
    this.publishRecordDetailsTableBody[3].title = this.publishTypeTpl
  }

  getPublishRecords () {
    this.api.get('cluster/' + this.clusterName + '/variable/publish-history', { page_num: 1, page_size: 20 }).subscribe(resp => {
      if (resp.code === 0) {
        this.publishRecordsData = resp.data
        for (const history of resp.data.historys) {
          history.isExpand = false
        }
        this.publishRecordsPage.total = resp.data.total
      } else {
        this.message.error(resp.msg || '获取列表数据失败!')
      }
    })
  }

  expandChange () {
    this.publishRecordsData.historys = [...this.publishRecordsData.historys]
  }

  publishRTableClick = (item:any) => {
    item.data.isExpand = !item.data.isExpand
  }
}

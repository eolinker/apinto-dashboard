import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { ApiService } from 'projects/core/src/app/service/api.service'

@Component({
  selector: 'eo-ng-deploy-environment-detail',
  template: `
  <div class="drawer-table">
    <eo-ng-apinto-table
      class="mr10 mt10"
      [nzTbody]="globalEnvDetailTableBody"
      [nzThead]="globalEnvDetailTableHeadName"
      [nzData]="globalEnvDetailList"
      [nzNoScroll]="true"
    >
    </eo-ng-apinto-table>
  </div>


  <ng-template #variableDetailStatusTpl let-item="item">
  <ng-container [ngSwitch]="item.publish_status">
    <span *ngSwitchCase="'UNPUBLISHED'" class="red-bold">未发布</span>
    <span *ngSwitchCase="'PUBLISHED'" class="green-bold">已发布</span>
    <span *ngSwitchCase="'DEFECT'" class="grey-bold">缺失</span>
  </ng-container>
  </ng-template>
  `,
  styles: [
  ]
})
export class DeployEnvironmentDetailComponent implements OnInit {
  @ViewChild('variableDetailStatusTpl', { read: TemplateRef, static: true })
  variableDetailStatusTpl: TemplateRef<any> | undefined

  envKey:string = ''
  globalEnvDetailList: Array<{
    // eslint-disable-next-line camelcase
    clusterName: string
    environment: string
    value: string
    // eslint-disable-next-line camelcase
    publish_status: string
  }> = []

  globalEnvDetailTableHeadName: Array<any> = [
    { title: '集群' },
    { title: '环境' },
    { title: 'VALUE' },
    { title: '状态' }
  ]

  globalEnvDetailTableBody: Array<any> = [
    { key: 'clusterName' },
    { key: 'environment' },
    { key: 'value' },
    { key: 'publish_status' }
  ]

  constructor (
    private message: EoNgFeedbackMessageService,
    private api: ApiService
  ) {
  }

  ngOnInit (): void {
    this.getEnvDetail()
  }

  ngAfterViewInit () {
    this.globalEnvDetailTableBody[3].title = this.variableDetailStatusTpl
  }

  getEnvDetail () {
    this.api.get('variable', { key: this.envKey || '' }).subscribe((resp) => {
      if (resp.code === 0) {
        this.globalEnvDetailList = resp.data.variables
      } else {
        this.message.error(resp.msg || '获取列表数据失败！')
      }
    })
  }
}

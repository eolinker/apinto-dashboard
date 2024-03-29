import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { THEAD_TYPE } from 'eo-ng-table'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { EO_TBODY_TYPE } from 'projects/eo-ng-apinto-table/src/public-api'
import { DeployGlobalEnvDetailTableHeadName, DeployGlobalEnvDetailTableBody } from '../types/conf'

@Component({
  selector: 'eo-ng-deploy-environment-detail',
  template: `
  <div class="drawer-table">
    <eo-ng-apinto-table
    class="drawer-table"
      class="mr-btnbase mt-btnbase"
      [nzTbody]="globalEnvDetailTableBody"
      [nzThead]="globalEnvDetailTableHeadName"
      [nzData]="globalEnvDetailList"
      [nzNoScroll]="true"
    >
    </eo-ng-apinto-table>
  </div>


  <ng-template #variableDetailStatusTpl let-item="item">
  <eo-ng-table-publish-status
    [publish]="item.publishStatus"
  ></eo-ng-table-publish-status>
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
    clusterName: string
    environment: string
    value: string
    publishStatus: string
  }> = []

  globalEnvDetailTableHeadName: THEAD_TYPE[] = [...DeployGlobalEnvDetailTableHeadName]
  globalEnvDetailTableBody: EO_TBODY_TYPE[] = [...DeployGlobalEnvDetailTableBody]

  constructor (
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
      }
    })
  }
}

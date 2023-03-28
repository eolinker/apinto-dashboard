import { Component, Input, OnInit } from '@angular/core'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { DeployService } from '../../../../deploy.service'

@Component({
  selector: 'eo-ng-deploy-cluster-environment-config-update',
  template: `<div class="update-config-drawer">
  <div class="drawer-list-header">
    <div nz-row nzJustify="start" nzAligh="top" class="mb-[20px]">
      <div nz-col>
        <label class="label table-label w-[82px]"
          ><span class="required-symbol">*</span>同步集群：</label
        >
      </div>
      <div nz-col class="" style="width: 100%">
        <div class="drawer-list-content">
          <eo-ng-apinto-table
            class="drawer-table"
            [nzTbody]="clusterTableBody"
            [nzThead]="clusterTableHeadName"
            [(nzData)]="clustersList"
            [nzNoScroll]="true"
            [nzScrollY]="120"
          >
          </eo-ng-apinto-table>

          <div *ngIf="startValidate && (!updateConfigForm.clusters || updateConfigForm.clusters.length === 0 )" class="ant-form-item-with-help">
          <div class="ant-form-item-explain">
            <div
              role="alert"
              class="ant-form-item-explain-error"
              style="margin-left: var(--LAYOUT_PADDING)"
            >
              必填项！
            </div>
          </div>
        </div>
        </div>
      </div>
    </div>
  </div>
  <div nz-row nzJustify="start" nzAligh="top">
    <div nz-col>
      <label class="label table-label  w-[82px]" ><span class="required-symbol">*</span>同步配置：</label>
    </div>
    <div nz-col class="" style="width: 100%">
      <div class="drawer-list-content" style="padding-bottom: 0px">
        <eo-ng-apinto-table
          class="drawer-table"
          [nzTbody]="configsTable2Body"
          [nzThead]="configsTable2HeadName"
          [(nzData)]="updateConfigsList"
          [nzNoScroll]="true"
        >
        </eo-ng-apinto-table>

        <div *ngIf="startValidate && (!updateConfigForm.variables ||  updateConfigForm.variables.length === 0)" class="ant-form-item-with-help">
        <div class="ant-form-item-explain">
          <div
            role="alert"
            class="ant-form-item-explain-error"
            style="margin-left: var(--LAYOUT_PADDING)"
          >
          必填项！
          </div>
        </div>
      </div>
      </div>
    </div>
  </div>
</div>
  `,
  styles: [
  ]
})
export class DeployClusterEnvironmentConfigUpdateComponent implements OnInit {
  @Input() closeModal?:(value?:any)=>void
  clustersList:Array<{env:string, status:string, name:string, checked:boolean, id:number}>=[]
  clusterTableHeadName:THEAD_TYPE[] = [...this.service.createClusterEnvUpdateThead(this)]
  clusterTableBody:TBODY_TYPE[] = [...this.service.createClusterEnvUpdateTbody(this)]
  configsTable2HeadName: THEAD_TYPE[] = [...this.service.createClusterEnvUpdate2Thead(this)]
  configsTable2Body: TBODY_TYPE[]=[...this.service.createClusterEnvUpdate2Tbody(this)]

  // eslint-disable-next-line camelcase
  updateConfigsList: Array<{ key: string, value: string, variableId: number, publish:string, status:string, desc:string, operator:string, updateTime:string, createTime:string, id: number, checked:boolean}> = []
  clusterName:string = ''

  // eslint-disable-next-line camelcase
  updateConfigForm:{clusters:Array<{name:string, env:string, id:number}>, variables:Array<{key:string, value:string, variableId:number, id:number}>}=
      {
        clusters: [],
        variables: []
      }

  startValidate:boolean = false // 开始校验数据，当用户点击过提交按钮才触发
  constructor (
    private message: EoNgFeedbackMessageService,
    private api:ApiService,
    private service:DeployService
  ) { }

  ngOnInit (): void {
    this.getUpdateData()
  }

  getUpdateData () {
    this.api.get('cluster/' + this.clusterName + '/variable/sync-conf').subscribe(resp => {
      if (resp.code === 0) {
        this.clustersList = resp.data.info.clusters
        this.updateConfigsList = resp.data.info.variables
      }
    })
  }

  getVarCheckedList () {
    setTimeout(() => {
      this.updateConfigForm.variables = this.updateConfigsList?.filter(config => {
        return config.checked
      })
    }, 0)
  }

  getClusterCheckedList () {
    setTimeout(() => {
      this.updateConfigForm.clusters = this.clustersList.filter(cluster => {
        return cluster.checked
      })
    }, 0)
  }

  save () {
    this.getClusterCheckedList()
    this.getVarCheckedList()

    this.startValidate = true

    if (this.updateConfigForm.clusters?.length > 0 && this.updateConfigForm.variables?.length > 0) {
      this.api.post('cluster/' + this.clusterName + '/variable/sync-conf', this.updateConfigForm).subscribe(resp => {
        if (resp.code === 0) {
          this.message.success(resp.msg || '同步成功', { nzDuration: 1000 })
          this.closeModal && this.closeModal()
        }
      })
    }
  }
}

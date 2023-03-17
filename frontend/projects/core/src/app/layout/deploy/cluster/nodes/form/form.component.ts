/* eslint-disable dot-notation */
import { Component, Input, OnInit, TemplateRef } from '@angular/core'
import { FormGroup, UntypedFormBuilder } from '@angular/forms'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { ApiService } from 'projects/core/src/app/service/api.service'

@Component({
  selector: 'eo-ng-deploy-cluster-nodes-form',
  template: `
  <form
  nz-form
  [nzNoColon]="true"
  [nzAutoTips]="autoTips"
  [formGroup]="validateResetNodeForm"
>
  <nz-form-item class="mb-0">
    <label class="label w-[70px]">集群地址：</label>

    <nz-form-control [nzErrorTip]="nodeErrorTpl">
      <div class="flex">
        <input
          eo-ng-input
          class=""
          formControlName="clusterAddr"
          placeholder="请输入"
        />

        <ng-template #nodeErrorTpl let-control>
          <ng-container *ngIf="control.hasError('required')"
            >请输入集群地址</ng-container
          >
          <ng-container *ngIf="control.hasError('pattern')"
            >集群地址输入错误，请重新输入</ng-container
          >

          <ng-container *ngIf="control.hasError('source')"
            >集群地址需要通过测试</ng-container
          >
        </ng-template>
        <button
          eo-ng-button
          type="button"
          class="ant-btn-primary ml-btnybase"
          (click)="testCluster()"
        >
          测试
        </button>
      </div>
    </nz-form-control>
  </nz-form-item>

  <nz-form-item class="mt-btnybase mb-0" *ngIf="nodesTestTableShow">
    <label class="label" style="width: 64px"></label>
    <nz-form-control>
      <div style="width: 100%">
        <eo-ng-apinto-table
          [nzTbody]="nodesTableBody"
          [nzThead]="nodesTableHeadName2"
          [nzData]="nodesTestList"
          [nzNoScroll]="true"
          nzTableLayout="fixed"
        >
        </eo-ng-apinto-table>
      </div>
    </nz-form-control>
  </nz-form-item>
</form>


  `,
  styles: [
    `
    input.ant-input[eo-ng-input]{
      width:494px !important;
    }`
  ]
})
export class DeployClusterNodesFormComponent implements OnInit {
  @Input() nodeStatusTpl: TemplateRef<any> | undefined
  @Input() serviceAddrTpl: TemplateRef<any> | undefined
  @Input() adminAddrTpl: TemplateRef<any> | undefined
  @Input() nodesTableBody: Array<any> =[]

  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  validateResetNodeForm: FormGroup = new FormGroup({})
  clusterName:string = ''
  testFlag:boolean = false
  testPassAddr:string = ''
  nodesTestTableShow = false
  clusterCanBeCreated:boolean = false

  nodesTestList:Array<any> = []
  source:string = ''

  nodesTableHeadName2:Array<any> = [
    {
      title: '名称',
      width: 150
    },
    { title: '管理地址' },
    { title: '服务地址' },
    {
      title: '状态'
    }
  ]

  constructor (
    private message: EoNgFeedbackMessageService,
    private api:ApiService,
    private fb: UntypedFormBuilder) {
    this.validateResetNodeForm = this.fb.group({
      clusterAddr: ['', [this.clusterAddrVadidator]]
    })
  }

  ngOnInit (): void {
  }

  // 集群地址为必填、有格式要求、且测试后需通过测试才有效
  clusterAddrVadidator = (control: any): { [s: string]: boolean } => {
    if (!control.value) {
      return { error: true, required: true }
    } else if (this.testFlag && !/(\w+):\/\/([a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+\.?)(:\d*)/.test(control.value) && !/(\w+):\/\/(((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3})(:\d*)/.test(control.value)) {
      return { error: true, pattern: true }
    } else if (this.testFlag && control.value !== this.testPassAddr) {
      return { source: true, error: true }
    }
    return {}
  }

  testCluster ():void {
    if (this.validateResetNodeForm.controls['clusterAddr'].value && (/(\w+):\/\/([a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+\.?)(:\d+)/.test(this.validateResetNodeForm.controls['clusterAddr'].value) || /(\w+):\/\/(((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3})(:\d*)/.test(this.validateResetNodeForm.controls['clusterAddr'].value))) {
      this.testFlag = true
      this.api.get('cluster-test', { cluster_addr: this.validateResetNodeForm.controls['clusterAddr'].value }).subscribe(resp => {
        if (resp.code === 0) {
          this.nodesTestList = resp.data.nodes
          this.source = resp.data.source
          this.clusterCanBeCreated = resp.data.is_update
          if (this.nodesTestList.length > 0) {
            this.nodesTestTableShow = true
          }
          if (this.clusterCanBeCreated) {
            this.testPassAddr = this.validateResetNodeForm.controls['clusterAddr'].value
          }
          this.validateResetNodeForm.controls['clusterAddr'].updateValueAndValidity({
            onlySelf: true
          })
        } else {
          this.validateResetNodeForm.controls['clusterAddr'].markAsDirty()
          this.validateResetNodeForm.controls['clusterAddr'].updateValueAndValidity({ onlySelf: true })
          this.message.error(resp.msg || '获取列表数据失败！', { nzDuration: 1000 })
        }
      })
    } else {
      this.validateResetNodeForm.controls['clusterAddr'].markAsDirty()
      this.validateResetNodeForm.controls['clusterAddr'].updateValueAndValidity()
    }
  }

  save () {
    this.testFlag = true
    this.validateResetNodeForm.controls['clusterAddr'].updateValueAndValidity()
    if (this.validateResetNodeForm.controls['clusterAddr'].valid && this.source) {
      this.api.post('cluster/' + this.clusterName + '/node/reset', { source: this.source || '', cluster_addr: this.validateResetNodeForm.controls['clusterAddr'].value || '' }).subscribe(resp => {
        if (resp.code === 0) {
          this.closeModal()
          this.message.success(resp.msg || '重置成功!', { nzDuration: 1000 })
          return true
        } else {
          this.message.error(resp.msg || '重置失败!')
          return false
        }
      })
    } else {
      Object.values(this.validateResetNodeForm.controls).forEach((control) => {
        if (control.invalid) {
          control.markAsDirty()
          control.updateValueAndValidity({ onlySelf: true })
        }
      })
      return false
    }
    return false
  }

  closeModal = () => {
  }
}

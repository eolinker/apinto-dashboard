/* eslint-disable dot-notation */
import { Component, Input, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { defaultAutoTips } from '../../../constant/conf'
import { EmptyHttpResponse } from '../../../constant/type'
import { ApiService } from '../../../service/api.service'
import { publishTableBody, publishTableHeadName } from '../types/conf'
import { StrategyPublishListData } from '../types/types'

@Component({
  selector: 'eo-ng-service-governance-publish',
  template: `
  <form
    nz-form
    [nzNoColon]="true"
    [nzAutoTips]="autoTips"
    [formGroup]="validateForm"
    autocomplete="off"
  >
    <nz-form-item>
      <label class="label" style="width: 80px"
        ><span class="required-symbol">*</span>发布名称：</label
      >
      <nz-form-control>
        <input
          eo-ng-input
          class="w-INPUT_NORMAL"
          formControlName="versionName"
          eoNgUserAccess="{{ 'serv-governance/' + strategyType }}"
        />
      </nz-form-control>
    </nz-form-item>
    <nz-form-item>
      <label class="label" style="width: 80px">描述：</label>
      <nz-form-control>
        <textarea
          class="w-INPUT_NORMAL"
          name="desc"
          rows="6"
          eo-ng-input
          formControlName="desc"
          placeholder="请输入"
          eoNgUserAccess="{{ 'serv-governance/' + strategyType }}"
          (disabledEdit)="disabledEdit($event)"
        ></textarea>
      </nz-form-control>
    </nz-form-item>
    <nz-form-item class="mb-0">
      <label class="label table-label" style="width: 80px"
        ><span class="required-symbol">*</span>策略列表：</label
      >
      <nz-form-control
        [nzValidateStatus]="
          !strategyIsPublish && strategyUnpublishMsg ? 'error' : ''
        "
        [nzErrorTip]="strategyUnpublishMsg"
      >
        <div>
          <eo-ng-apinto-table
            [nzTbody]="publishTableBody"
            [nzThead]="publishTableHeadName"
            [(nzData)]="publishList"
            [nzNoScroll]="true"
            nzTableLayout="fixed"
          >
          </eo-ng-apinto-table>
        </div>
      </nz-form-control>
    </nz-form-item>
  </form>


<ng-template #strategiesStatusTpl let-item="item">
  <eo-ng-table-status [status]="item.status"></eo-ng-table-status>
</ng-template>
  `,
  styles: [
  ]
})
export class ServiceGovernancePublishComponent implements OnInit {
  @ViewChild('strategiesStatusTpl', { read: TemplateRef, static: true }) strategiesStatusTpl: TemplateRef<any> | undefined
  @Input() closeModal?:(value?:any)=>void
  validateForm: FormGroup = new FormGroup({})
  nzDisabled:boolean = false
  strategyType:string = ''

  publishTableHeadName:THEAD_TYPE[] = [...publishTableHeadName]
  publishTableBody:TBODY_TYPE[]= [...publishTableBody]
  publishList: Array<any> = []

  strategySource:string = ''
  strategyUnpublishMsg:string = ''

  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  strategyIsPublish:boolean = false
  clusterName:string = ''
  returnToSdk:Function|undefined
  constructor (
                private message: EoNgFeedbackMessageService,
                private api:ApiService,
                private fb: UntypedFormBuilder) {
    this.validateForm = this.fb.group({
      versionName: ['', [Validators.required]],
      desc: ['']
    })
  }

  ngOnInit (): void {
    this.getPublishList()
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  ngAfterContentInit () {
    this.publishTableBody[2].title = this.strategiesStatusTpl
  }

  // 获取待发布的策略列表
  getPublishList () {
    this.api.get('strategy/' + this.strategyType + '/to-publishs', { clusterName: this.clusterName })
      .subscribe((resp:{ code:number, data:StrategyPublishListData, msg:string}) => {
        if (resp.code === 0) {
          this.publishList = resp.data.strategies
          this.strategyIsPublish = resp.data.isPublish
          this.strategySource = resp.data.source
          this.validateForm.controls['versionName'].setValue(resp.data.versionName)
          this.strategyUnpublishMsg = resp.data.unpublishMsg
          if (!resp.data.isPublish && !this.strategyUnpublishMsg) {
            if (this.publishList?.length > 0) {
              this.strategyUnpublishMsg = '当前策略不可发布'
            } else {
              this.strategyUnpublishMsg = '当前无可发布策略'
            }
          }
        } else {
          this.returnToSdk && this.returnToSdk(resp)
        }
      })
  }

  // 发布策略，仅当表单校验通过且策略可发布时才可发布，否则需显示提示语
  publish () {
    if (this.validateForm.valid && this.strategyIsPublish) {
      this.api.post('strategy/' + this.strategyType + '/publish', { versionName: (this.validateForm.controls['versionName'].value || ''), desc: (this.validateForm.controls['desc'].value || ''), source: (this.strategySource || '') }, { clusterName: this.clusterName })
        .subscribe((resp:EmptyHttpResponse) => {
          if (resp.code === 0) {
            this.message.success(resp.msg || '发布策略成功!', { nzDuration: 1000 })
            this.closeModal && this.closeModal()
          }
          this.returnToSdk && this.returnToSdk(resp)
        })
    } else {
      Object.values(this.validateForm.controls).forEach((control) => {
        if (control.invalid) {
          control.markAsDirty()
          control.updateValueAndValidity({ onlySelf: true })
        }
      })
    }
  }
}

/* eslint-disable dot-notation */
import { ViewportScroller } from '@angular/common'
import { Component, Input, OnInit, TemplateRef } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { ActivatedRoute, Router } from '@angular/router'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { defaultAutoTips } from '../../../constant/conf'
import { ApiService } from '../../../service/api.service'
import { BaseInfoService } from '../../../service/base-info.service'

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
          formControlName="version_name"
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
          !strategyIsPublish && strategyUnpulishMsg ? 'error' : ''
        "
        [nzErrorTip]="strategyUnpulishMsg"
      >
        <div>
          <eo-ng-apinto-table
            [nzTbody]="publishTableBody"
            [nzThead]="publishTabelHeadName"
            [(nzData)]="publishList"
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
  ]
})
export class ServiceGovernancePublishComponent implements OnInit {
  @Input() closeModal?:(value?:any)=>void
  validateForm: FormGroup = new FormGroup({})
  nzDisabled:boolean = false
  strategyType:string = ''

  publishTabelHeadName:Array<any> = [
    {
      title: '策略名称'
    },
    {
      title: '优先级'
    },
    {
      title: '状态'
    },
    {
      title: '操作时间'
    }
  ]

  publishTableBody:Array<any> = [
    { key: 'name' },
    { key: 'priority' },
    { key: 'status' },
    { key: 'opt_time' }
  ]

  publishList: Array<any> = []

  strategySource:string = ''
  strategyUnpulishMsg:string = ''

  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  strategiesStatusTpl:TemplateRef<any>|undefined
  strategyIsPublish:string = ''
  clusterName:string = ''
  constructor (private baseInfo:BaseInfoService,
    private viewportScroller: ViewportScroller,
                private message: EoNgFeedbackMessageService,
                private modalService:EoNgFeedbackModalService,
                private api:ApiService,
                private activateInfo:ActivatedRoute,
                private router:Router,
                private fb: UntypedFormBuilder) {
    this.validateForm = this.fb.group({
      version_name: ['', [Validators.required]],
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
    this.api.get('strategy/' + this.strategyType + '/to-publishs', { cluster_name: this.clusterName }).subscribe((resp:any) => {
      if (resp.code === 0) {
        this.publishList = resp.data.strategies
        this.strategyIsPublish = resp.data.is_publish
        this.strategySource = resp.data.source
        this.validateForm.controls['version_name'].setValue(resp.data.version_name)
        this.strategyUnpulishMsg = resp.data.unpublish_msg
        if (!resp.data.is_publish && !this.strategyUnpulishMsg) {
          if (this.publishList?.length > 0) {
            this.strategyUnpulishMsg = '当前策略不可发布'
          } else {
            this.strategyUnpulishMsg = '当前无可发布策略'
          }
        }
      } else {
        this.message.error(resp.msg || '获取数据失败!')
      }
    })
  }

  // 发布策略，仅当表单校验通过且策略可发布时才可发布，否则需显示提示语
  publish () {
    if (this.validateForm.valid && this.strategyIsPublish) {
      this.api.post('strategy/' + this.strategyType + '/publish', { version_name: (this.validateForm.controls['version_name'].value || ''), desc: (this.validateForm.controls['desc'].value || ''), source: (this.strategySource || '') }, { cluster_name: this.clusterName }).subscribe((resp:any) => {
        if (resp.code === 0) {
          this.message.success(resp.msg || '发布策略成功!', { nzDuration: 1000 })
          this.closeModal && this.closeModal()
        } else {
          this.message.error(resp.msg || '发布策略失败!')
        }
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

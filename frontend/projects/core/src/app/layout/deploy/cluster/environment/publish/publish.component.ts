import { Component, Input, OnInit, TemplateRef } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { DeployClusterPublishTbody, DeployClusterPublishThead } from '../../types/conf'

@Component({
  selector: 'eo-ng-deploy-cluster-environment-publish',
  template: `
  <div>
  <form
    nz-form
    [nzNoColon]="true"
    [nzAutoTips]="autoTips"
    [formGroup]="validatePublishForm"
  >
  <div class='drawer-list-header'>
    <nz-form-item>
      <label class="label" style="width: 100px"
        ><span class="required-symbol">*</span>发布名称：</label
      >
      <nz-form-control>
        <input
          class="w-INPUT_NORMAL"
          eo-ng-input
          eoNgUserAccess="deploy/cluster"
          formControlName="versionName"
          placeholder="请输入"
        />
      </nz-form-control>
    </nz-form-item>
    <nz-form-item>
      <label class="label" style="width: 100px">描述：</label>
      <nz-form-control>
        <textarea
          class="w-INPUT_NORMAL"
          rows="8"
          eo-ng-input
          eoNgUserAccess="deploy/cluster"
          formControlName="desc"
          placeholder="请输入"
        ></textarea>
      </nz-form-control>
    </nz-form-item>
</div>
    <nz-form-item class="mb-0">
      <label class="label table-label" style="width: 100px"
        >环境变量列表：</label
      >
      <nz-form-control
        [nzValidateStatus]="publishData.isPublish ? '' : 'error'"
        [nzErrorTip]="unpublishMsgTpl"
      >
        <div style="width: 100%">
          <eo-ng-apinto-table
          class="drawer-table"
            [nzTbody]="publishTableBody"
            [nzThead]="publishTabelHeadName"
            [(nzData)]="publishData.variables"
            [nzNoScroll]="true"
          >
          </eo-ng-apinto-table>
        </div>
        <ng-template #unpublishMsgTpl>
        <div class="drawer-list-footer">
        {{unpublishMsg}}
        </div>
        </ng-template>
      </nz-form-control>
    </nz-form-item>
  </form>
</div>
  `,
  styles: [
  ]
})
export class DeployClusterEnvironmentPublishComponent implements OnInit {
  @Input() publishTypeTpl: TemplateRef<any> | undefined
  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  validatePublishForm:FormGroup = new FormGroup({})
  clusterName:string=''
  // eslint-disable-next-line camelcase
  public unpublishMsg:string = ''
  nzDisabled:boolean = false
  publishSource:string = ''
  // eslint-disable-next-line camelcase
  publishData:{source:string, variables:Array<{key:string, finishValue:string, noReleasedValue:string, createTime:string, optType:string}>, isPublish:boolean, versionName:string}=
      {
        source: '',
        variables: [],
        isPublish: false,
        versionName: ''
      }

  publishTabelHeadName: THEAD_TYPE[] = [...DeployClusterPublishThead]
  publishTableBody: TBODY_TYPE[] =[...DeployClusterPublishTbody]

  constructor (
    private message: EoNgFeedbackMessageService,
    private api:ApiService,
    private fb: UntypedFormBuilder) {
    this.validatePublishForm = this.fb.group({
      versionName: ['', [Validators.required]],
      desc: ['']
    })
  }

  ngOnInit (): void {
    this.getPublishData()
  }

  ngAfterViewInit () {
    this.publishTableBody[3].title = this.publishTypeTpl
  }

  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  getPublishData () {
    this.api.get('cluster/' + this.clusterName + '/variable/to-publishs').subscribe(resp => {
      if (resp.code === 0) {
        this.publishData = resp.data
        this.publishData.variables = this.publishData.variables || []
        // eslint-disable-next-line dot-notation
        this.validatePublishForm.controls['versionName'].setValue(resp.data.versionName)
        this.publishSource = resp.data.source
        this.unpublishMsg = resp.data.unpublishMsg
        if (!this.publishData.isPublish && !this.unpublishMsg) {
          if (this.publishData.variables.length === 0) {
            this.unpublishMsg = '无环境变量可发布'
          } else {
            this.unpublishMsg = '当前环境变量不可发布'
          }
        }
      } else {
        this.message.error(resp.msg || '获取列表数据失败!')
      }
    })
  }

  save (type:string):boolean {
    switch (type) {
      case 'publish':
        if (this.validatePublishForm.valid && this.publishData.isPublish) {
          this.api.post('cluster/' + this.clusterName + '/variable/publish', {
            versionName: this.validatePublishForm.value.versionName,
            desc: this.validatePublishForm.value.desc,
            source: this.publishSource
          }).subscribe((resp: { code: number; msg: any }) => {
            if (resp.code === 0) {
              this.message.success(resp.msg || '发布成功', { nzDuration: 1000 })
              this.closeModal(true)
              return true
            } else {
              this.message.error(resp.msg || '发布失败!')
              return false
            }
          })
        } else {
          Object.values(this.validatePublishForm.controls).forEach(control => {
            if (control.invalid) {
              control.markAsDirty()
              control.updateValueAndValidity({ onlySelf: true })
            }
          })
          return false
        }
        return false
      default:
        return false
    }
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  closeModal = (fresh?:boolean) => {}
}

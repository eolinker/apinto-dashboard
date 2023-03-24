import { Component, Input, OnInit, TemplateRef } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { SelectOption } from 'eo-ng-select'
import { THEAD_TYPE, TBODY_TYPE } from 'eo-ng-table'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { DeployClusterPluginPublishTbody, DeployClusterPluginPublishThead, DeployClusterPluginStatusOptions } from '../../types/conf'
import { ClusterPluginPublishData } from '../../types/types'

@Component({
  selector: 'eo-ng-deploy-cluster-plugin-publish',
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
        >插件列表：</label
      >
      <nz-form-control
        [nzValidateStatus]="publishData.isPublish ? '' : 'error'"
        [nzErrorTip]="unpublishMsgTpl"
      >
        <div style="width: 100%">
          <eo-ng-apinto-table
          class="drawer-table"
            [nzTbody]="publishTableBody"
            [nzThead]="publishTableHeadName"
            [(nzData)]="publishData.plugins"
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
export class DeployClusterPluginPublishComponent implements OnInit {
  @Input() publishTypeTpl: TemplateRef<any> | undefined
  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  validatePublishForm:FormGroup = new FormGroup({})
  clusterName:string=''
  // eslint-disable-next-line camelcase
  public unpublishMsg:string = ''
  nzDisabled:boolean = false
  publishSource:string = ''
  // eslint-disable-next-line camelcase
  publishData:{source:string, plugins:ClusterPluginPublishData[], isPublish:boolean, versionName:string}=
      {
        source: '',
        plugins: [],
        isPublish: false,
        versionName: ''
      }

  publishTableHeadName: THEAD_TYPE[] = [...DeployClusterPluginPublishThead]
  publishTableBody: TBODY_TYPE[] =[...DeployClusterPluginPublishTbody]
  statusList:SelectOption[] = [...DeployClusterPluginStatusOptions]
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
    this.api.get('cluster/' + this.clusterName + '/plugin/to-publish').subscribe((resp:{code:number, data:{source:string, plugins:ClusterPluginPublishData[], isPublish:boolean, versionName:string, unpublishedMsg:string}, msg:string}) => {
      if (resp.code === 0) {
        this.publishData = resp.data
        this.publishData.plugins = resp.data.plugins?.map((item) => {
          item.finishValue = `插件顺序：${item.releasedSort}，状态：${this.getStatusString(item.releasedConfig.status)}，配置：${item.releasedConfig.config}`
          item.noReleasedValue = `插件顺序：${item.nowSort}，状态：${this.getStatusString(item.noReleasedConfig.status)}，配置：${item.noReleasedConfig.config}`
          return item
        }) || []
        // eslint-disable-next-line dot-notation
        this.validatePublishForm.controls['versionName'].setValue(resp.data.versionName)
        this.publishSource = resp.data.source
        this.unpublishMsg = resp.data.unpublishedMsg
        if (!this.publishData.isPublish && !this.unpublishMsg) {
          if (this.publishData.plugins.length === 0) {
            this.unpublishMsg = '无插件可发布'
          } else {
            this.unpublishMsg = '当前插件不可发布'
          }
        }
      } else {
        this.message.error(resp.msg || '获取列表数据失败!')
      }
    })
  }

  getStatusString (status:'GLOBAL'|'DISABLE'|'ENABLE') {
    return this.statusList.filter((item:SelectOption) => { return item.value === status })[0].label
  }

  save ():boolean {
    if (this.validatePublishForm.valid && this.publishData.isPublish) {
      this.api.post('cluster/' + this.clusterName + '/plugin/publish', {
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
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  closeModal = (fresh?:boolean) => {}
}

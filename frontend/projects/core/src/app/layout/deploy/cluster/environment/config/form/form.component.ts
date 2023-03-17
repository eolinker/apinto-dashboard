/* eslint-disable dot-notation */
import { Component, Input, OnInit } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { ApiService } from 'projects/core/src/app/service/api.service'

@Component({
  selector: 'eo-ng-deploy-cluster-environment-config-form',
  template: `
  <form
  class="deploy-drawer"
  nz-form
  [nzNoColon]="true"
  [nzAutoTips]="autoTips"
  [formGroup]="validateAddConfigForm"
>
  <nz-form-item>
    <nz-form-label [nzSpan]="6" nzRequired>KEY：</nz-form-label>
    <nz-form-control [nzSpan]="13" [nzErrorTip]="matchKeyErrorTpl">
      <input
        class="w-INPUT_NORMAL"
        eo-ng-input
        formControlName="key"
        eoNgUserAccess="deploy/cluster"
        placeholder="英文、下划线及圆点组合"
      />
      <ng-template #matchKeyErrorTpl let-control>
        <ng-container *ngIf="control.hasError('pattern')"
          >英文、下划线及圆点组合</ng-container
        >
        <ng-container *ngIf="control.hasError('required')"
          >必填项</ng-container
        ></ng-template
      >
    </nz-form-control>
  </nz-form-item>

  <nz-form-item>
    <nz-form-label [nzSpan]="6">VALUE：</nz-form-label>
    <nz-form-control [nzSpan]="13" [nzExtra]="valueExtraTpl">
      <textarea
        class="w-INPUT_NORMAL"
        rows="8"
        eoNgUserAccess="deploy/cluster"
        eo-ng-input
        formControlName="value"
        placeholder="请输入"
      ></textarea>
    </nz-form-control>
  </nz-form-item>

  <ng-template #valueExtraTpl>
    <p class="w-INPUT_NORMAL">
      注意：隐藏字符(空格、换行符、制表符Tab)容易导致配置出错，如果需要检测Value中隐藏字符，请点击<a
        (click)="detect()"
        >检测隐藏字符</a
      >
    </p>
    <p
      class="detect_value_bg w-INPUT_NORMAL"
      *ngIf="showValueAfterTest"
      [innerHtml]="valueAfterTest"
    ></p>
  </ng-template>

  <nz-form-item class="mb-0">
    <nz-form-label [nzSpan]="6">描述：</nz-form-label>

    <nz-form-control [nzSpan]="13">
      <textarea
        class="w-INPUT_NORMAL"
        eo-ng-input
        formControlName="desc"
        placeholder="请输入"
        eoNgUserAccess="deploy/cluster"
      ></textarea>
    </nz-form-control>
  </nz-form-item>
</form>
  `,
  styles: [
    `
    .detect_value_bg {
      background-color: var(--bar-background-color);;
      padding-left: var(--LAYOUT_PADDING);
      word-break:break-all;
    }`
  ]
})
export class DeployClusterEnvironmentConfigFormComponent implements OnInit {
  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  validateAddConfigForm:FormGroup = new FormGroup({})
  valueAfterTest:string = '' // 通过检查隐藏字符后的值
  showValueAfterTest:boolean = false // 是否显示检测隐藏字符的值
  clusterName:string = ''
  @Input() editData?:{key:string, value:string, desc:string}
  constructor (
          private message: EoNgFeedbackMessageService,
          private api:ApiService,
          private fb: UntypedFormBuilder) {
  }

  ngOnInit (): void {
    this.validateAddConfigForm = this.fb.group({
      key: [this.editData?.key || '', [Validators.required, Validators.pattern(/^[A-Za-z_\\.]*$/)]],
      value: [this.editData?.value || ''],
      desc: [this.editData?.desc || '']
    })
    if (this.editData) {
      this.validateAddConfigForm.controls['key'].disable()
      this.validateAddConfigForm.controls['desc'].disable()
    }
  }

  detect ():void {
    if (this.validateAddConfigForm.value.value !== '') {
      this.valueAfterTest = this.validateAddConfigForm.value.value
        .replace(/(\n|\r|\r\n|↵)/g, '#换行符#')
        .replace(/\t/g, '#制表符#')
        .replace(/\s/g, '#空格#')
        .replace(/#空格#/g, '<span class="detected-symbol">#空格#</span>')
        .replace(/#换行符#/g, '<span class="detected-symbol">#换行符#</span>')
        .replace(/#制表符#/g, '<span class="detected-symbol">#制表符#</span>')
      this.showValueAfterTest = true
    }
  }

  save () {
    if (this.validateAddConfigForm.valid) {
      if (this.editData) {
        this.api.put('cluster/' + this.clusterName + '/variable', { value: this.validateAddConfigForm.value.value || '' }, { key: this.validateAddConfigForm.controls['key'].value || '' }).subscribe(resp => {
          if (resp.code === 0) {
            this.message.success(resp.msg || '编辑成功', { nzDuration: 1000 })
            this.closeModal(true)
            return true
          } else {
            this.message.error(resp.msg || '操作失败!')
            return false
          }
        })
      } else {
        this.api.post('cluster/' + this.clusterName + '/variable', {
          key: this.validateAddConfigForm.value.key || '',
          value: this.validateAddConfigForm.value.value || '',
          desc: this.validateAddConfigForm.value.desc || ''
        },
        { key: this.validateAddConfigForm.value.key || '' }).subscribe(resp => {
          if (resp.code === 0) {
            this.message.success(resp.msg || '添加成功', { nzDuration: 1000 })
            this.closeModal(true)
          } else {
            this.message.error(resp.msg || '添加失败!')
          }
        })
      }
    } else {
      Object.values(this.validateAddConfigForm.controls).forEach(control => {
        if (control.invalid) {
          control.markAsDirty()
          control.updateValueAndValidity({ onlySelf: true })
        }
      })
    }
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  closeModal = (fresh?:boolean) => {}
}

/* eslint-disable dot-notation */
import { Component } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { AutoCompleteOption } from 'eo-ng-auto-complete'
import { NzUploadFile } from 'ng-zorro-antd/upload'
import { defaultAutoTips } from '../../../constant/conf'
import { EmptyHttpResponse } from '../../../constant/type'
import { ApiService } from '../../../service/api.service'
import { EoNgMessageService } from '../../../service/eo-ng-message.service'

@Component({
  selector: 'eo-ng-plugin-create',
  template: `
    <form
      nz-form
      [nzNoColon]="true"
      [nzAutoTips]="autoTips"
      [formGroup]="validateForm"
      autocomplete="off"
    >
      <nz-form-item class="form-row">
        <nz-form-label [nzSpan]="6" nzRequired>上传文件：</nz-form-label>
        <nz-form-control [nzSpan]="14">
          <nz-upload
            [(nzFileList)]="fileList"
            [nzBeforeUpload]="beforeUpload"
            [nzLimit]="1"
            [nzRemove]="removeFile"
            nzAccept=".zip"
          >
            <button id="uploadBtn" [nzDanger]="fileError" eo-ng-button>
              选择文件
            </button>
          </nz-upload>
            <div
              *ngIf="fileError"
              class="ant-form-item-explain-error"
              (click)="$event.stopPropagation()"
            >
              必填项
            </div>
            <div
              class="ant-form-item-extra activation-extra"
              style="padding-left: 0"
              (click)="$event.stopPropagation()"
            >
              仅支持官方提供插件配置模板文件
            </div>
        </nz-form-control>
      </nz-form-item>
  <nz-form-item class="form-row">
    <nz-form-label [nzSpan]="6" nzRequired>分组：</nz-form-label>
    <nz-form-control [nzSpan]="14">
      <eo-ng-auto-complete
            formControlName="name"
            [nzOptions]="groupOptions"
          ></eo-ng-auto-complete>
    </nz-form-control>
  </nz-form-item>
</form>
  `,
  styles: [
  ]
})
export class PluginCreateComponent {
  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  validateForm: FormGroup = new FormGroup({})
  fileError:boolean = false
  fileList: NzUploadFile[] = [];
  groupOptions:AutoCompleteOption[] = []
  name= ''
  constructor (
    public api:ApiService,
    private fb: UntypedFormBuilder,
    private message:EoNgMessageService) {
    this.validateForm = this.fb.group({
      file: [null, [Validators.required]],
      name: ['', [Validators.required]]
    })
    this.getGroupList()
  }

  getGroupList () {
    this.api.get('system/plugin/groups/enum').subscribe((resp:{code:number, data:{groups:Array<{uuid:string, name:string}>}, msg:string}) => {
      if (resp.code === 0) {
        this.groupOptions = resp.data.groups.map((group:{uuid:string, name:string}) => {
          return { label: group.name, value: group.name }
        })
      }
    })
  }

  // 手动上传文件
  beforeUpload = (file: NzUploadFile): boolean => {
    this.fileList = []
    this.fileList = this.fileList.concat(file)
    this.fileError = this.fileList.length === 0
    return false
  }

  // 移除文件
  removeFile () {
    this.fileList = []
    this.fileError = true
    return true
  }

  checkValid () {
    if (this.fileList.length === 0 || !this.validateForm.controls['name'].valid) {
      this.fileError = this.fileList.length === 0
      if (this.validateForm.controls['name'].invalid) {
        this.validateForm.controls['name'].markAsDirty()
        this.validateForm.controls['name'].updateValueAndValidity({ onlySelf: true })
      }
      return false
    }
    return true
  }

  submit () {
    if (!this.checkValid()) { return false }
    const formData = new FormData()
    formData.append('plugin', this.fileList[0] as any)
    // eslint-disable-next-line dot-notation
    formData.append('group_name', this.validateForm.controls['name'].value as any)
    return this.api.post('system/plugin/install', formData).subscribe((resp:EmptyHttpResponse) => {
      if (resp.code === 0) {
        this.message.success(resp.msg || '安装插件成功')
        return true
      } else {
        return false
      }
    })
  }
}

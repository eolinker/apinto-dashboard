/* eslint-disable dot-notation */
import { Component, Input, OnInit } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { v4 as uuidv4 } from 'uuid'

@Component({
  selector: 'eo-ng-api-management-edit-group',
  template: `
  <form
    class="api-group"
    nz-form
    [formGroup]="validateApiGroupForm"
    [nzNoColon]="true"
    [nzAutoTips]="autoTips"
    autocomplete="off"
  >

  <nz-form-item *ngIf="showUuid">
      <nz-form-label [nzSpan]="6" nzFor="groupName">分组ID：</nz-form-label>
      <nz-form-control [nzSpan]="15">
      <span
      class="overflow-ellipsis inline-block overflow-hidden align-middle"
      >{{ uuid }}
    </span>
    <button
      eo-copy
      eo-ng-button
      nzType="primary"
      nzGhost
      class="deploy-node-copy-btn ant-btn-text border-transparent h-[22px]"
      [copyText]="uuid"
      (copyCallback)="copyCallback()"
    >
      <svg class="iconpark-icon"><use href="#copy"></use></svg>
    </button>
      </nz-form-control>
      </nz-form-item>

    <nz-form-item class="mb-0">
      <nz-form-label [nzSpan]="6" nzFor="groupName">分组名称：</nz-form-label>
      <nz-form-control [nzSpan]="15">
        <input
          eo-ng-input
          id="groupName"
          eoNgAutoFocus
          placeholder="分组名称"
          formControlName="groupName"
          (keyup.enter)="
            type === 'edit' ? editGroup(uuid) : addGroup(uuid)
          "
        />
      </nz-form-control>
    </nz-form-item>
  </form>
  `,
  styles: [
  ]
})
export class ApiManagementEditGroupComponent implements OnInit {
  @Input() type:string = ''
  @Input() uuid:string = ''
  @Input() groupName:string = ''
  @Input() closeModal?:(value?:any)=>void
  @Input() showUuid:boolean = false
  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  validateApiGroupForm:FormGroup = new FormGroup({})
  constructor (private message: EoNgFeedbackMessageService,
    private api:ApiService,
    private fb: UntypedFormBuilder,
    private router:Router) {
    this.validateApiGroupForm = this.fb.group({
      groupName: [this.groupName, [Validators.required]]
    })
  }

  ngOnInit (): void {
    this.validateApiGroupForm = this.fb.group({
      groupName: [this.groupName, [Validators.required]]
    })
  }

  // 添加分组的请求
  addGroup (parentUuid:string) {
    if (this.validateApiGroupForm.valid) {
      const uuid = uuidv4()
      this.api.post('group/api', { name: this.validateApiGroupForm.controls['groupName'].value, uuid: uuid, parent_uuid: parentUuid }).subscribe((resp:any) => {
        if (resp.code === 0) {
          this.message.success(resp.msg || '添加成功', { nzDuration: 1000 })
          this.router.navigate(['/', 'router', 'group', 'list', uuid])
          this.closeModal && this.closeModal()
        } else {
          this.message.error(resp.msg || '添加失败!')
        }
      })
    } else {
      Object.values(this.validateApiGroupForm.controls).forEach(control => {
        if (control.invalid) {
          control.markAsDirty()
          control.updateValueAndValidity({ onlySelf: true })
        }
      })
    }
  }

  // 编辑分组
  editGroup (groupUuid:string) {
    if (this.validateApiGroupForm.valid) {
      this.api.put('group/api/' + groupUuid, { name: this.validateApiGroupForm.controls['groupName'].value }).subscribe((resp:any) => {
        if (resp.code === 0) {
          this.message.success(resp.msg || '修改成功', { nzDuration: 1000 })
          this.closeModal && this.closeModal()
        } else {
          this.message.error(resp.msg || '修改失败!')
        }
      })
    } else {
      Object.values(this.validateApiGroupForm.controls).forEach(control => {
        if (control.invalid) {
          control.markAsDirty()
          control.updateValueAndValidity({ onlySelf: true })
        }
      })
    }
  }

  copyCallback = () => {
    this.message.success('复制成功', { nzDuration: 1000 })
  };
}

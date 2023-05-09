import { Component, OnInit } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { SelectOption } from 'eo-ng-select'
import { setFormValue } from 'projects/core/src/app/constant/form'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { EoNgMessageService } from 'projects/core/src/app/service/eo-ng-message.service'
import { protocolsList } from '../../types/conf'
import { EmailData } from '../../types/type'
import { EoNgNavigationService } from 'projects/core/src/app/service/eo-ng-navigation.service'

@Component({
  selector: 'eo-ng-system-email-config',
  templateUrl: './config.component.html',
  styles: [
  ]
})
export class SystemEmailConfigComponent implements OnInit {
  editPage: boolean = false
  validateForm: FormGroup = new FormGroup({})
  nzDisabled: boolean = false
  emailId:string = ''
  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  listOfProtocols:SelectOption[] = [...protocolsList]
  constructor (private fb: UntypedFormBuilder,
    private navigationService: EoNgNavigationService,
    private api: ApiService,
    private message: EoNgMessageService) {
    this.validateForm = this.fb.group({
      smtpUrl: ['', [Validators.required]],
      smtpPort: [null, [Validators.required]],
      protocol: ['ssl', [Validators.required]],
      email: ['', [Validators.email]],
      account: [''],
      password: ['']
    })

    this.navigationService.reqFlashBreadcrumb([
      { title: '邮箱设置' }
    ])
  }

  ngOnInit (): void {
    this.getEmailMessage()
  }

  disabledEdit (editAccess: boolean) {
    this.nzDisabled = editAccess
  }

  getEmailMessage () {
    this.api.get('email').subscribe((resp:{code:number, data:{emailInfo?:EmailData}, msg:string}) => {
      if (resp.code === 0 && resp.data.emailInfo) {
        resp.data.emailInfo.protocol = resp.data.emailInfo.protocol === '' ? 'none' : resp.data.emailInfo.protocol
        setFormValue(this.validateForm, resp.data.emailInfo)
        this.emailId = resp.data.emailInfo.uuid
        this.editPage = true
      }
    })
  }

  save () {
    if (this.validateForm.valid) {
      const data:EmailData = { ...this.validateForm.value }
      if (data.protocol === 'none') { data.protocol = '' }
      if (this.editPage) {
        this.api.put('email', { ...data, uuid: this.emailId }).subscribe((resp:{code:number, data:{}, msg:string}) => {
          if (resp.code === 0) {
            this.message.success(resp.msg || '修改通知邮箱成功！')
          }
        })
      } else {
        this.api.post('email', { ...data }).subscribe((resp:{code:number, data:{}, msg:string}) => {
          if (resp.code === 0) {
            this.message.success(resp.msg || '创建通知邮箱成功！')
          }
        })
      }
    } else {
      Object.values(this.validateForm.controls).forEach(control => {
        if (control.invalid) {
          control.markAsDirty()
          control.updateValueAndValidity({ onlySelf: true })
        }
      })
    }
  }
}

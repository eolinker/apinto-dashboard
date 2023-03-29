/* eslint-disable dot-notation */
import { Component, Input, OnInit } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'

@Component({
  selector: 'eo-ng-external-app-create',
  templateUrl: './create.component.html',
  styles: [
  ]
})
export class ExternalAppCreateComponent implements OnInit {
  @Input() editPage: boolean = false
  @Input() appId: string = ''
  validateForm: FormGroup = new FormGroup({})
  submitButtonLoading:boolean = false

  constructor (
    private message: EoNgFeedbackMessageService,
    private api: ApiService,
    private router: Router,
    private fb: UntypedFormBuilder,
    private appConfigService: AppConfigService
  ) {
    this.appConfigService.reqFlashBreadcrumb([
      { title: '外部应用' },
      { title: '新建外部应用' }

    ])

    this.validateForm = this.fb.group({
      name: ['', [Validators.required]],
      id: [''],
      desc: ['']
    })
  }

  autoTips: Record<string, Record<string, string>> = defaultAutoTips

  nzDisabled: boolean = false
  ngOnInit (): void {
    if (this.editPage) {
      this.getApplicationMessage()
      this.appConfigService.reqFlashBreadcrumb([
        { title: '外部应用' },
        { title: '外部应用详情' }
      ])
    } else {
      this.getApplicationId()
    }
  }

  getApplicationMessage () {
    this.api
      .get('external-app', { id: this.appId })
      .subscribe((resp: {code:number, data:{info:{name:string, id:string, desc:string}}, msg:string}) => {
        if (resp.code === 0) {
          this.validateForm.controls['name'].setValue(
            resp.data.info.name
          )
          this.validateForm.controls['id'].setValue(resp.data.info.id)
          this.validateForm.controls['desc'].setValue(
            resp.data.info.desc
          )
          this.validateForm.controls['id'].disable()
        }
      })
  }

  getApplicationId () {
    this.api.get('random/external-app/id').subscribe((resp: any) => {
      if (resp.code === 0) {
        this.validateForm.controls['id'].setValue(resp.data.id)
      }
    })
  }

  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  // 保存鉴权，editPage = true时，表示页面为编辑页，false为新建页
  saveApplication () {
    if (this.validateForm.valid) {
      this.submitButtonLoading = true
      if (!this.editPage) {
        this.api
          .post('external-app', {
            ...this.validateForm.value
          })
          .subscribe((resp: any) => {
            this.submitButtonLoading = false
            if (resp.code === 0) {
              this.message.success(resp.msg || '添加成功')
              this.backToList()
            }
          })
      } else {
        this.api
          .put('external-app', {
            ...this.validateForm.value
          }, { id: this.validateForm.controls['id'].value })
          .subscribe((resp: any) => {
            this.submitButtonLoading = false
            if (resp.code === 0) {
              this.message.success(resp.msg || '修改成功')
              this.backToList()
            }
          })
      }
    } else {
      Object.values(this.validateForm.controls).forEach((control) => {
        if (control.invalid) {
          control.markAsDirty()
          control.updateValueAndValidity({ onlySelf: true })
        }
      })
    }
  }

  backToList () {
    this.router.navigate(['/', 'system', 'ext-app'])
  }
}

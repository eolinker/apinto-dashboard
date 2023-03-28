/* eslint-disable camelcase */
/* eslint-disable dot-notation */
import { Component, Input, OnInit } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { Router, ActivatedRoute } from '@angular/router'
import {
  EoNgFeedbackMessageService
} from 'eo-ng-feedback'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'

@Component({
  selector: 'eo-ng-deploy-environment-create',
  templateUrl: './create.component.html',
  styles: []
})
export class DeployEnvironmentCreateComponent implements OnInit {
  @Input() editPage: boolean = false
  validateForm: FormGroup = new FormGroup({})
  VariableName: string = ''
  autoTips: Record<string, Record<string, string>> = defaultAutoTips

  globalEnvDetailList: Array<{
    clusterName: string
    environment: string
    value: string
    publish_status: string
  }> = []

  constructor (
    private message: EoNgFeedbackMessageService,
    private api: ApiService,
    private fb: UntypedFormBuilder,
    private router: Router,
    private activateInfo: ActivatedRoute,
    private appConfigService: AppConfigService
  ) {
    this.validateForm = this.fb.group({
      key: [
        '',
        [Validators.required, Validators.pattern('^[a-zA-Z][a-zA-Z0-9/_]*')]
      ],
      desc: ['']
    })
    this.appConfigService.reqFlashBreadcrumb([
      { title: '环境变量', routerLink: 'deploy/env' },
      { title: '新建配置' }
    ])
  }

  ngOnInit (): void {
  }

  save () {
    if (this.validateForm.valid) {
      this.api
        .post('variable', {
          key: this.validateForm.controls['key'].value,
          desc: this.validateForm.controls['desc'].value || ''
        })
        .subscribe((resp) => {
          if (resp.code === 0) {
            this.message.success(resp.msg || '新增环境变量成功！', { nzDuration: 1000 })
            this.router.navigate(['/', 'deploy', 'env'])
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

  backToList () {
    this.router.navigate(['/', 'deploy', 'env'])
  }
}

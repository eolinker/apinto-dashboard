/* eslint-disable dot-notation */
import { Component, Input, OnInit } from '@angular/core'
import { FormGroup, UntypedFormBuilder } from '@angular/forms'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { SelectOption } from 'eo-ng-select'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { DeployClusterPluginStatusOptions } from '../../types/conf'
import { ClusterPluginConfig } from '../../types/types'

@Component({
  selector: 'eo-ng-deploy-cluster-plugin-config-form',
  templateUrl: './config.component.html',
  styles: [
  ]
})
export class DeployClusterPluginConfigFormComponent implements OnInit {
  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  validateConfigForm:FormGroup = new FormGroup({})
  clusterName:string = ''
  statusList:SelectOption[] = [...DeployClusterPluginStatusOptions]
  nzDisabled:boolean = false
  @Input() editData?:ClusterPluginConfig
  constructor (
          private message: EoNgFeedbackMessageService,
          private api:ApiService,
          private fb: UntypedFormBuilder) {
  }

  ngOnInit (): void {
    this.validateConfigForm = this.fb.group({
      name: [this.editData?.name || ''],
      status: [this.editData?.status || ''],
      config: [this.editData?.config || {}]
    })
  }

  save () {
    if (this.validateConfigForm.valid) {
      this.api.post('cluster/' + this.clusterName + '/plugin', {
        name: this.validateConfigForm.value.name || '',
        status: this.validateConfigForm.value.status || '',
        config: this.validateConfigForm.value.config || ''
      }).subscribe(resp => {
        if (resp.code === 0) {
          this.message.success(resp.msg || '添加成功', { nzDuration: 1000 })
          this.closeModal(true)
        } else {
          this.message.error(resp.msg || '添加失败!')
        }
      })
    } else {
      Object.values(this.validateConfigForm.controls).forEach(control => {
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

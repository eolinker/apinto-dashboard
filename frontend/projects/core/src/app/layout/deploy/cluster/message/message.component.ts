/* eslint-disable dot-notation */
import { Component } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { setFormValue } from 'projects/core/src/app/constant/form'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { ClusterData } from '../types/types'
import { EoNgNavigationService } from 'projects/core/src/app/service/eo-ng-navigation.service'
import { DeployService } from '../../deploy.service'

@Component({
  selector: 'eo-ng-cluster-message',
  templateUrl: './message.component.html',
  styles: [
  ]
})
export class DeployClusterMessageComponent {
  clusterName:string = ''
  validateForm: FormGroup = new FormGroup({})
  environmentList: Array<{ label: string; value: any }> = []

  submitButtonLoading:boolean = false

  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  constructor (
    private message: EoNgFeedbackMessageService,
    private api: ApiService,
    private router: Router,
    private fb: UntypedFormBuilder,
    private service:DeployService,
    private baseInfo:BaseInfoService,
    private navigationService: EoNgNavigationService) {
    this.validateForm = this.fb.group({
      title: ['', [Validators.required]],
      env: ['', [Validators.required]],
      desc: ['']
    })
  }

  ngOnInit (): void {
    this.navigationService.reqFlashBreadcrumb([{ title: '网关集群', routerLink: 'deploy/cluster' }, { title: '集群管理' }])
    this.clusterName = this.baseInfo.allParamsInfo.clusterName
    this.getClusterMessage()
    this.getEnvList()
  }

  getClusterMessage () {
    this.api.get('cluster', { clusterName: this.clusterName }).subscribe((resp:{code:number, data:{cluster:ClusterData}, msg:string}) => {
      if (resp.code === 0) {
        this.service.clusterDesc = resp.data.cluster.desc
        this.service.clusterName = resp.data.cluster.title
        setFormValue(this.validateForm, resp.data.cluster)
      }
    })
  }

  getEnvList () {
    this.api.get('enum/envs').subscribe((resp:{code:number, data:{envs:Array<{name:string, value:string}>}, msg:string}) => {
      if (resp.code === 0) {
        this.environmentList = resp.data.envs.map(
          (env: { name: string; value: string }) => ({
            label: env.name,
            value: env.value
          })
        )
        this.validateForm
          .controls['env']
          .setValue(this.environmentList[0].value)
        this.validateForm.controls['env'].updateValueAndValidity({
          onlySelf: true
        })
      }
    })
  }

  // 修改集群信息
  saveCluster () {
    if (this.validateForm.valid) {
      this.submitButtonLoading = true
      this.api.put(`cluster/${this.clusterName}`, {
        env: this.validateForm.controls['env'].value,
        title: this.validateForm.controls['title'].value,
        desc: this.validateForm.controls['desc'].value
      }).subscribe((resp) => {
        this.submitButtonLoading = false
        if (resp.code === 0) {
          this.service.clusterDesc = this.validateForm.controls['desc'].value
          this.service.clusterName = this.validateForm.controls['title'].value
          this.message.success(resp.msg || '修改集群信息成功')
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

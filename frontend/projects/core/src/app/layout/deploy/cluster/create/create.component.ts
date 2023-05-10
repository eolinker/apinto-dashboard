/* eslint-disable dot-notation */
import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { EoNgNavigationService } from 'projects/core/src/app/service/eo-ng-navigation.service'
import { DeployClusterNodeTbody, DeployClusterNodeThead } from '../types/conf'
import { v4 as uuidv4 } from 'uuid'
@Component({
  selector: 'eo-ng-cluster-create',
  templateUrl: './create.component.html',
  styles: [
    `
    `
  ]
})
export class DeployClusterCreateComponent implements OnInit {
  @ViewChild('nodeStatusTpl', { read: TemplateRef, static: true }) nodeStatusTpl:
    | TemplateRef<any>
    | undefined

  validateForm: FormGroup = new FormGroup({})
  source: string = '' // 集群地址通过测试后得到的source, 有source的情况才能新建集群成功
  environmentList: Array<{ label: string; value: any }> = []
  nodesList: Array<object> = []

  nodesTableHeadName: THEAD_TYPE[] = [...DeployClusterNodeThead]
  nodesTableBody: TBODY_TYPE[] = [...DeployClusterNodeTbody]
  submitButtonLoading:boolean = false
  testButtonLoading:boolean = false
  nodesTableShow = false
  clusterCanBeCreated: boolean = false
  testFlag:boolean = false
  testPassAddr:string = '' // 通过测试的集群地址

  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  clusterProtocol:string = 'http'
  startCheckCluster:boolean = false
  checkClusterError:boolean = false // 接口检查该网关节点无法连接
  checkClusterErrorAddr:string = ''// 接口检查该网关节点无法连接
  closeModal:Function|undefined
  newAddr:string = '' // 拼接后的地址
  constructor (
    private message: EoNgFeedbackMessageService,
    private api: ApiService,
    private router: Router,
    private fb: UntypedFormBuilder,
    private navigationService: EoNgNavigationService) {
    this.validateForm = this.fb.group({
      clusterName: ['', [Validators.required]],
      envValue: ['', [Validators.required]],
      clusterDesc: [''],
      clusterAddr: ['', [this.clusterAddrVadidator]]
    })
  }

  ngOnInit (): void {
    this.getEnvList()
  }

  ngAfterViewInit ():void {
    this.nodesTableBody[3].title = this.nodeStatusTpl
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
          .controls['envValue']
          .setValue(this.environmentList[0].value)
        this.validateForm.controls['envValue'].updateValueAndValidity({
          onlySelf: true
        })
      }
    })
  }

  testCluster (): void {
    this.startCheckCluster = true
    if (this.validateForm.valid) {
      this.testFlag = true
      this.testButtonLoading = true
      this.newAddr = this.validateForm.controls['clusterAddr'].value
      if ((/^(\/\/)/).test(this.validateForm.controls['clusterAddr'].value)) {
        this.newAddr = `${this.clusterProtocol}:${this.newAddr}`
      } else if ((/^(\/)/).test(this.validateForm.controls['clusterAddr'].value)) {
        this.newAddr = `${this.clusterProtocol}:/${this.newAddr}`
      } else {
        this.newAddr = `${this.clusterProtocol}://${this.newAddr}`
      }

      this.api
        .get('cluster-test', {
          clusterAddr: this.newAddr
        })
        .subscribe((resp) => {
          this.testButtonLoading = false
          if (resp.code === 0) {
            this.nodesList = resp.data.nodes
            this.clusterCanBeCreated = resp.data.isUpdate
            this.source = resp.data.source
            if (this.nodesList.length > 0) {
              this.nodesTableShow = true
            }
            if (this.source) {
              this.testPassAddr = this.validateForm.controls['clusterAddr'].value
            }
            this.validateForm.controls['clusterAddr'].updateValueAndValidity({
              onlySelf: true
            })
          } else {
            this.checkClusterError = true
            this.checkClusterErrorAddr = this.validateForm.controls['clusterAddr'].value
            this.validateForm.controls['clusterAddr'].markAsDirty()
            this.validateForm.controls['clusterAddr'].updateValueAndValidity({ onlySelf: true })
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

  // 集群地址为必填、有格式要求、且测试后需通过测试才有效
  clusterAddrVadidator = (control: any): { [s: string]: boolean } => {
    if (!control.value) {
      return { error: true, required: true }
    } else if (!/([a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})*\.?)*(:\d+)/.test(control.value) && !/(((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3})(:\d*)/.test(control.value)) {
      return { error: true, pattern: true }
    } else if (this.checkClusterError && control.value === this.checkClusterErrorAddr) {
      return { source: true, error: true }
    }
    return {}
  }

  // 新建集群
  saveCluster () {
    this.validateForm.markAllAsTouched()
    if (this.validateForm.valid || this.checkValidForm()) {
      const name:string = `cluster_${uuidv4().replace(/-/g, '_')}`
      const params = {
        title: this.validateForm.controls['clusterName'].value,
        name: name,
        desc: this.validateForm.controls['clusterDesc'].value || '',
        addr: this.newAddr,
        source: this.source || '',
        env: this.validateForm.controls['envValue'].value
      }
      this.submitButtonLoading = true
      this.api.post('cluster', params).subscribe((resp) => {
        this.submitButtonLoading = false
        if (resp.code === 0) {
          this.closeModal && this.closeModal()
          this.router.navigate(['/', 'deploy', 'cluster', 'content', name])
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

  checkValidForm () {
    for (const index in this.validateForm.controls) {
      if (this.validateForm.controls[index].invalid) {
        return false
      }
    }
    return true
  }

  // 取消新建集群
  cancel () {
    if (this.nodesTableShow) {
      this.nodesTableShow = false
    } else {
      this.closeModal && this.closeModal()
    }
  }
}

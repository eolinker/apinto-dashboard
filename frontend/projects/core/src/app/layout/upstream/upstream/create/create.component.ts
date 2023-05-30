/* eslint-disable dot-notation */
/* eslint-disable camelcase */
/* eslint-disable no-useless-constructor */
import { Component, Input, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { SelectOption } from 'eo-ng-select'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { EoNgNavigationService } from 'projects/core/src/app/service/eo-ng-navigation.service'
import { UpstreamBalanceList, UpstreamSchemeList } from '../types/conf'

@Component({
  selector: 'eo-ng-upstream-create',
  templateUrl: './create.component.html',
  styles: [
    `
    `
  ]
})
export class UpstreamCreateComponent implements OnInit {
  @ViewChild('environmentRef', { read: TemplateRef, static: true })
  environmentRef: TemplateRef<any> | undefined

  @Input() editPage: boolean = false
  @Input() serviceName: string = ''
  baseData: any = null
  createUpstreamForm: {
    name: string
    desc: string
    scheme: string
    balance: string
    discoveryName: string
    timeout: number
    config: {
      serviceName: string
      useVariable: boolean
      addrsVariable: string
      staticConf: Array<{ weight: number | null; addr: string }>
    }
  } = {
    name: '',
    desc: '',
    scheme: 'HTTP',
    balance: 'round-robin',
    discoveryName: 'static',
    timeout: 100,
    config: {
      addrsVariable: '',
      useVariable: false,
      serviceName: '',
      staticConf: [
        {
          weight: null,
          addr: ''
        }
      ]
    }
  }

  schemeList: SelectOption[] = [...UpstreamSchemeList]
  balanceList: SelectOption[] = [...UpstreamBalanceList]
  discoveryList: Array<{ label: string; value: string; render: any }> = []
  useEnvVar: boolean = false
  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  submitButtonLoading:boolean = false

  startValidateDynamic:boolean = false
  validateForm: FormGroup = new FormGroup({})
  nzDisabled:boolean = false // 权限控制
  constructor (
    private message: EoNgFeedbackMessageService,
    private api: ApiService,
    private router: Router,
    private fb: UntypedFormBuilder,
    private navigationService: EoNgNavigationService
  ) {
    this.validateForm = this.fb.group({
      name: ['', [Validators.required, Validators.pattern('^[a-zA-Z][a-zA-Z0-9/_]*')]],
      desc: [''],
      scheme: ['HTTP', [Validators.required]],
      balance: ['round-robin', [Validators.required]],
      discoveryName: ['static', [Validators.required]],
      timeout: [100, [Validators.required]]
    })
    this.navigationService.reqFlashBreadcrumb([{ title: '上游管理', routerLink: 'upstream/upstream' }, { title: '创建上游' }])
  }

  public _baseData: any = null
  showDynamicTips: boolean = false
  ngOnInit (): void {
    if (this.editPage) {
      this.getUpstreamMessage()
      this.navigationService.reqFlashBreadcrumb([{ title: '上游管理', routerLink: 'upstream/upstream' }, { title: '上游信息' }])
    } else {
      this.getDiscovery()
    }
  }

  getDiscovery () {
    this.api.get('discovery/enum').subscribe((resp) => {
      if (resp.code === 0) {
        for (const index in resp.data.discoveries) {
          if (resp.data.discoveries[index].driver !== 'static') {
            this.discoveryList.push({
              label: `${resp.data.discoveries[index].name}[${resp.data.discoveries[index].driver}]`,
              value: resp.data.discoveries[index].name,
              render: resp.data.discoveries[index].render
            })
          } else {
            this.discoveryList.push({
              label: '静态节点',
              value: resp.data.discoveries[index].name,
              render: resp.data.discoveries[index].render
            })
          }
          if (
            resp.data.discoveries[index].name ===
            this.validateForm.controls['discoveryName'].value
          ) {
            this.baseData = resp.data.discoveries[index].render
          }
        }
        this.discoveryList = [...this.discoveryList]
      }
    })
  }

  changeBasedata () {
    this.discoveryList.forEach((ele: any) => {
      if (this.validateForm.controls['discoveryName'].value === ele.value) {
        this.baseData = ele.render
      }
    })
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  getUpstreamMessage () {
    this.api.get('service', { name: this.serviceName }).subscribe((resp) => {
      if (resp.code === 0) {
        this.validateForm.controls['name']!.setValue(resp.data.service.name)
        this.validateForm.controls['desc']!.setValue(resp.data.service.desc)
        this.validateForm.controls['scheme']!.setValue(resp.data.service.scheme)
        this.validateForm.controls['balance']!.setValue(
          resp.data.service.balance
        )
        this.validateForm.controls['discoveryName']!.setValue(
          resp.data.service.discoveryName
        )
        this.validateForm.controls['timeout']!.setValue(
          resp.data.service.timeout
        )
        this.validateForm.controls['name']!.disable()
        this.createUpstreamForm = resp.data.service
        this.createUpstreamForm.config.serviceName = this.createUpstreamForm
          .config.serviceName
          ? this.createUpstreamForm.config.serviceName
          : ''
        this.createUpstreamForm.config.useVariable = this.createUpstreamForm
          .config.useVariable
          ? this.createUpstreamForm.config.useVariable
          : false
        this.createUpstreamForm.config.addrsVariable = this.createUpstreamForm
          .config.addrsVariable
          ? this.createUpstreamForm.config.addrsVariable
          : ''
        if (
          !this.createUpstreamForm.config.staticConf ||
          this.createUpstreamForm.config.staticConf.length === 0
        ) {
          this.createUpstreamForm.config.staticConf = [
            { addr: '', weight: null }
          ]
        }
        this.getDiscovery()
      }
    })
  }

  canBeSave: boolean = false

  // 动态渲染的组件内每次改变值时，都会调用该方法，以供判断表单是否可以提交
  getDataFromDynamicComponent (value: any) {
    // 静态节点
    if (value && this.validateForm.controls['discoveryName'].value === 'static') {
      // 地址选用环境变量
      if (value.config.useVariable && !value.config.addrsVariable) {
        this.canBeSave = false
        return
      }
      // 地址不选用环境变量
      if (!value.config.useVariable) {
        let flag = false
        for (const index in value.config.staticConf) {
          if (
            value.config.staticConf[index].addr &&
            value.config.staticConf[index].addr !== '' &&
            value.config.staticConf[index].weight > 0 &&
            value.config.staticConf[index].weight < 1000
          ) {
            flag = true
          }
        }
        if (!flag) {
          this.canBeSave = false
          return
        }
      }
    } else if (!value.config.serviceName) {
      // 动态节点
      this.canBeSave = false
      return
    }

    this.canBeSave = true
  }

  checkTimeout () {
    if (this.validateForm.controls['timeout'].value !== null && this.validateForm.controls['timeout'].value < 1) {
      this.validateForm.controls['timeout'].setValue(1)
    }
  }

  saveUpstream () {
    this.startValidateDynamic = true
    this.showDynamicTips = !this.canBeSave
    if ((this.validateForm.valid || this.checkValidForm()) && this.canBeSave) {
      this.createUpstreamForm.config.useVariable =
        this.createUpstreamForm.config.useVariable !== null
          ? this.createUpstreamForm.config.useVariable
          : false
      this.createUpstreamForm.config.addrsVariable =
        this.createUpstreamForm.config.addrsVariable !== null
          ? this.createUpstreamForm.config.addrsVariable
          : ''
      this.createUpstreamForm.config.serviceName =
        this.createUpstreamForm.config.serviceName !== null
          ? this.createUpstreamForm.config.serviceName
          : ''
      this.createUpstreamForm.config.useVariable =
        this.createUpstreamForm.config.useVariable !== null
          ? this.createUpstreamForm.config.useVariable
          : false
      if (this.createUpstreamForm.config.staticConf.length > 0) {
        for (const index in this.createUpstreamForm.config.staticConf) {
          this.createUpstreamForm.config.staticConf[index].weight = Number(
            this.createUpstreamForm.config.staticConf[index].weight
          )
          this.createUpstreamForm.config.staticConf[index].addr =
            this.createUpstreamForm.config.staticConf[index].addr === null
              ? ''
              : this.createUpstreamForm.config.staticConf[index].addr
        }
      }
      this.createUpstreamForm.timeout = Number(this.createUpstreamForm.timeout)
      this.submitButtonLoading = true
      if (!this.editPage) {
        this.api
          .post('service', {
            ...this.validateForm.value,
            config: this.createUpstreamForm.config
          })
          .subscribe((resp) => {
            this.submitButtonLoading = false
            if (resp.code === 0) {
              this.message.success(resp.msg || '新建服务成功!', { nzDuration: 1000 })
              this.backToList()
            }
          })
      } else {
        this.api
          .put(
            'service',
            {
              ...this.validateForm.value,
              name: this.validateForm.controls['name']!.value,
              config: this.createUpstreamForm.config
            },
            {
              name: this.validateForm.controls['name']!.value
            }
          )
          .subscribe((resp) => {
            this.submitButtonLoading = false
            if (resp.code === 0) {
              this.message.success(resp.msg || '修改服务成功!', { nzDuration: 1000 })
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

  checkValidForm () {
    for (const index in this.validateForm.controls) {
      if (this.validateForm.controls[index].invalid) {
        return false
      }
    }
    return true
  }

  backToList () {
    this.router.navigate(['/', 'upstream', 'upstream'])
  }
}

/* eslint-disable dot-notation */
/* eslint-disable camelcase */
/* eslint-disable no-useless-constructor */
import { Component, Input, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'

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
    discovery_name: string
    timeout: number
    config: {
      service_name: string
      use_variable: boolean
      addrs_variable: string
      static_conf: Array<{ weight: number | null; addr: string }>
    }
  } = {
    name: '',
    desc: '',
    scheme: 'HTTP',
    balance: 'round-robin',
    discovery_name: 'static',
    timeout: 100,
    config: {
      addrs_variable: '',
      use_variable: false,
      service_name: '',
      static_conf: [
        {
          weight: null,
          addr: ''
        }
      ]
    }
  }

  schemeList: Array<{ label: string; value: string }> = [
    { label: 'HTTP', value: 'HTTP' },
    { label: 'HTTPS', value: 'HTTPS' }
  ]

  balanceList: Array<{ label: string; value: string }> = [
    { label: 'round-robin', value: 'round-robin' }
  ]

  discoveryList: Array<{ label: string; value: string; render: any }> = []

  useEnvVar: boolean = false

  autoTips: Record<string, Record<string, string>> = defaultAutoTips


  startValidateDynamic:boolean = false
  validateForm: FormGroup = new FormGroup({})
  nzDisabled:boolean = false // 权限控制
  constructor (
    private message: EoNgFeedbackMessageService,
    private api: ApiService,
    private router: Router,
    private fb: UntypedFormBuilder,
    private appConfigService: AppConfigService
  ) {
    this.validateForm = this.fb.group({
      name: ['', [Validators.required, Validators.pattern('^[a-zA-Z][a-zA-Z0-9/_]*')]],
      desc: [''],
      scheme: ['HTTP', [Validators.required]],
      balance: ['round-robin', [Validators.required]],
      discovery_name: ['static', [Validators.required]],
      timeout: [100, [Validators.required]]
    })
    this.appConfigService.reqFlashBreadcrumb([{ title: '上游管理', routerLink: 'upstream/upstream' }, { title: '创建上游' }])
  }

  public _baseData: any = null
  showDynamicTips: boolean = false
  ngOnInit (): void {
    if (this.editPage) {
      this.getUpstreamMessage()
      this.appConfigService.reqFlashBreadcrumb([{ title: '上游管理', routerLink: 'upstream/upstream' }, { title: '上游信息' }])
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
            this.validateForm.controls['discovery_name'].value
          ) {
            this.baseData = resp.data.discoveries[index].render
          }
        }
        this.discoveryList = [...this.discoveryList]
      } else {
        this.message.error(resp.msg || '获取数据失败!')
      }
    })
  }

  changeBasedata () {
    this.discoveryList.forEach((ele: any) => {
      if (this.validateForm.controls['discovery_name'].value === ele.value) {
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
        this.validateForm.controls['discovery_name']!.setValue(
          resp.data.service.discovery_name
        )
        this.validateForm.controls['timeout']!.setValue(
          resp.data.service.timeout
        )
        this.validateForm.controls['name']!.disable()
        this.createUpstreamForm = resp.data.service
        this.createUpstreamForm.config.service_name = this.createUpstreamForm
          .config.service_name
          ? this.createUpstreamForm.config.service_name
          : ''
        this.createUpstreamForm.config.use_variable = this.createUpstreamForm
          .config.use_variable
          ? this.createUpstreamForm.config.use_variable
          : false
        this.createUpstreamForm.config.addrs_variable = this.createUpstreamForm
          .config.addrs_variable
          ? this.createUpstreamForm.config.addrs_variable
          : ''
        if (
          !this.createUpstreamForm.config.static_conf ||
          this.createUpstreamForm.config.static_conf.length === 0
        ) {
          this.createUpstreamForm.config.static_conf = [
            { addr: '', weight: null }
          ]
        }
        this.getDiscovery()
      } else {
        this.message.error(resp.msg || '获取数据失败!')
      }
    })
  }

  canBeSave: boolean = false

  // 动态渲染的组件内每次改变值时，都会调用该方法，以供判断表单是否可以提交
  getDataFromDynamicComponent (value: any) {
    // 静态节点
    if (value && this.validateForm.controls['discovery_name'].value === 'static') {
      // 地址选用环境变量
      if (value.config.use_variable && !value.config.addrs_variable) {
        this.canBeSave = false
        return
      }
      // 地址不选用环境变量
      if (!value.config.use_variable) {
        let flag = false
        for (const index in value.config.static_conf) {
          if (
            value.config.static_conf[index].addr &&
            value.config.static_conf[index].addr !== '' &&
            value.config.static_conf[index].weight > 0 &&
            value.config.static_conf[index].weight < 1000
          ) {
            flag = true
          }
        }
        if (!flag) {
          this.canBeSave = false
          return
        }
      }
    } else if (!value.config.service_name) {
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
      this.createUpstreamForm.config.use_variable =
        this.createUpstreamForm.config.use_variable !== null
          ? this.createUpstreamForm.config.use_variable
          : false
      this.createUpstreamForm.config.addrs_variable =
        this.createUpstreamForm.config.addrs_variable !== null
          ? this.createUpstreamForm.config.addrs_variable
          : ''
      this.createUpstreamForm.config.service_name =
        this.createUpstreamForm.config.service_name !== null
          ? this.createUpstreamForm.config.service_name
          : ''
      this.createUpstreamForm.config.use_variable =
        this.createUpstreamForm.config.use_variable !== null
          ? this.createUpstreamForm.config.use_variable
          : false
      if (this.createUpstreamForm.config.static_conf.length > 0) {
        for (const index in this.createUpstreamForm.config.static_conf) {
          this.createUpstreamForm.config.static_conf[index].weight = Number(
            this.createUpstreamForm.config.static_conf[index].weight
          )
          this.createUpstreamForm.config.static_conf[index].addr =
            this.createUpstreamForm.config.static_conf[index].addr === null
              ? ''
              : this.createUpstreamForm.config.static_conf[index].addr
        }
      }
      this.createUpstreamForm.timeout = Number(this.createUpstreamForm.timeout)
      if (!this.editPage) {
        this.api
          .post('service', {
            ...this.validateForm.value,
            config: this.createUpstreamForm.config
          })
          .subscribe((resp) => {
            if (resp.code === 0) {
              this.message.success(resp.msg || '新建服务成功!', { nzDuration: 1000 })
              this.backToList()
            } else {
              this.message.error(resp.msg || '新建服务失败!')
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
            if (resp.code === 0) {
              this.message.success(resp.msg || '修改服务成功!', { nzDuration: 1000 })
              this.backToList()
            } else {
              this.message.error(resp.msg || '修改服务失败!')
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

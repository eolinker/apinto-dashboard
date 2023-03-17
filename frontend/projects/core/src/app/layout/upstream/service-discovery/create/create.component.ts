/* eslint-disable camelcase */
/* eslint-disable dot-notation */
/* eslint-disable no-useless-constructor */
import { Component, Input, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'

@Component({
  selector: 'eo-ng-service-discovery-create',
  templateUrl: './create.component.html',
  styles: []
})
export class ServiceDiscoveryCreateComponent implements OnInit {
  @ViewChild('environmentRef', { read: TemplateRef, static: true })
  environmentRef: TemplateRef<any> | undefined

  @Input() editPage: boolean = false
  @Input() serviceName: string = ''

  startValidateDynamic:boolean = false
  // baseData:any = FormJson.serviceData
  baseData: any = ''
  public _baseData: any = null
  tempForDynamic: string = 'MockDynamic'
  canBeSave: boolean = false
  useEnvVar: boolean = false
  nzDisabled:boolean = false
  createServiceForm: {
    name: string
    desc: string
    driver: string
    config: {
      use_variable: boolean
      addrs_variable: string
      addrs: Array<any>
      params: Array<{ key: string; value: string }>
    }
  } = {
    name: '',
    desc: '',
    driver: 'nacos',
    config: {
      use_variable: false,
      addrs_variable: '',
      addrs: [],
      params: [{ key: '', value: '' }]
    }
  }

  driverList: Array<{ label: string; value: string; render: any }> = []

  environmentTableHeadName: Array<object> = [
    { title: 'KEY' },
    { title: '描述' }
  ]

  environmentTableBody: Array<any> = [
    {
      key: 'key'
    },
    {
      key: 'description'
    }
  ]

  environmentList: Array<any> = []

  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  
  validateForm: FormGroup = new FormGroup({})
  showDynamicTips: boolean = false

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
      driver: ['nacos', [Validators.required]]
    })
    this.appConfigService.reqFlashBreadcrumb([{ title: '服务发现', routerLink: 'upstream/serv-discovery' }, { title: '新建服务' }])
  }

  ngOnInit (): void {
    if (this.editPage) {
      this.appConfigService.reqFlashBreadcrumb([{ title: '服务发现', routerLink: 'upstream/serv-discovery' }, { title: '服务信息' }])
      this.validateForm.controls['driver'].disable()
      this.validateForm.controls['name'].disable()
      this.getServiceMessage()
    } else {
      this.getDriverList()
    }
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  getDriverList () {
    this.api.get('discovery/drivers').subscribe((resp) => {
      if (resp.code === 0) {
        resp.data.drivers.forEach((element: any) => {
          this.driverList.push({
            label: element.name,
            value: element.name,
            render: element.render
          })
          if (element.name === this.validateForm.controls['driver'].value) {
            this.baseData = element.render
          }
        })
        this.driverList = [...this.driverList]
      } else {
        this.message.error(resp.msg || '获取列表数据失败!')
      }
    })
  }

  changeBasedata () {
    this.driverList.forEach((ele: any) => {
      if (this.validateForm.controls['driver'].value === ele.label) {
        this.baseData = ele.render
      }
    })
  }

  getServiceMessage () {
    this.api.get('discovery', { name: this.serviceName }).subscribe((resp) => {
      if (resp.code === 0) {
        this.validateForm.controls['name'].setValue(resp.data.discovery.name)
        this.validateForm.controls['desc'].setValue(resp.data.discovery.desc)
        this.validateForm.controls['driver'].setValue(
          resp.data.discovery.driver
        )
        this.createServiceForm = resp.data.discovery
        if (
          !this.createServiceForm.config.params ||
          this.createServiceForm.config.params.length === 0
        ) {
          this.createServiceForm.config.params = [{ key: '', value: '' }]
        }
        this.createServiceForm.config.addrs = this.createServiceForm.config.addrs ? this.createServiceForm.config.addrs : []
        this.createServiceForm.config.addrs_variable = this.createServiceForm.config.addrs_variable ? this.createServiceForm.config.addrs_variable : ''
        this.getDriverList()
      } else {
        this.message.error(resp.msg || '获取数据失败!')
      }
    })
  }

  getDataFromDynamicComponent (value: any) {
    if (value) {
      // 地址选用环境变量
      if (!value.config.use_variable) {
        if (value.config.addrs.length === 0) {
          this.canBeSave = false
          return
        } else {
          let flag = false
          for (const index in value.config.addrs) {
            if (value.config.addrs[index]) {
              flag = true
            }
          }
          if (!flag) {
            this.canBeSave = false
            return
          }
        }
      }

      // 地址不选用环境变量
      if (value.config.use_variable && !value.config.addrs_variable) {
        this.canBeSave = false
        return
      }

      this.canBeSave = true
    }
  }

  saveService () {
    this.startValidateDynamic = true
    this.showDynamicTips = !this.canBeSave

    if ((this.validateForm.valid || this.checkValidForm()) && this.canBeSave) {
      if (this.createServiceForm.config.params.length > 0) {
        const tempRes: Array<{ key: string; value: string }> = []
        for (const index in this.createServiceForm.config.params) {
          if (
            this.createServiceForm.config.params[index].key &&
            this.createServiceForm.config.params[index].value
          ) {
            tempRes.push({
              key: this.createServiceForm.config.params[index].key,
              value: this.createServiceForm.config.params[index].value
            })
          }
        }
        this.createServiceForm.config.params = tempRes
      } else {
        this.createServiceForm.config.params = []
      }
      this.createServiceForm.config.addrs =
      (this.createServiceForm.config.addrs === null || this.createServiceForm.config.use_variable)
        ? []
        : Array.from(new Set(this.createServiceForm.config.addrs))
      this.createServiceForm.config.addrs_variable =
      (this.createServiceForm.config.addrs_variable === null || !this.createServiceForm.config.use_variable)
        ? ''
        : this.createServiceForm.config.addrs_variable
      this.createServiceForm.config.use_variable =
        this.createServiceForm.config.use_variable === null
          ? false
          : this.createServiceForm.config.use_variable
      if (!this.editPage) {
        this.api
          .post('discovery', {
            ...this.validateForm.value,
            config: this.createServiceForm.config
          })
          .subscribe((resp) => {
            if (resp.code === 0) {
              this.message.success(resp.msg || '新建服务成功', { nzDuration: 1000 })
              this.backToList()
            } else {
              this.message.error(resp.msg || '新建服务失败!')
            }
          })
      } else {
        this.api
          .put(
            'discovery',
            {
              ...this.validateForm.value,
              config: this.createServiceForm.config
            },
            {
              name: this.validateForm.controls['name'].value
            }
          )
          .subscribe((resp) => {
            if (resp.code === 0) {
              this.message.success(resp.msg || '修改服务成功', { nzDuration: 1000 })
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
    this.router.navigate(['/', 'upstream', 'serv-discovery'])
  }
}

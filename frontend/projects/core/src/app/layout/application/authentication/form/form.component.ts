/* eslint-disable dot-notation */
import { Component, OnInit } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { differenceInCalendarDays } from 'date-fns'
import { SelectOption } from 'eo-ng-select'
import { AuthData } from '../../types/types'
import { algorithmList, authLabelTableBody, positionList, verifyList } from '../../types/conf'
import { TBODY_TYPE } from 'eo-ng-table'

@Component({
  selector: 'eo-ng-application-authentication-form',
  templateUrl: './form.component.html',
  styles: [
  ]
})
export class ApplicationAuthenticationFormComponent implements OnInit {
  appId:string = ''
  authId:string = ''
  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  driverList:Array<{label:string, value:string, render:any}> = []
  validateForm: FormGroup = new FormGroup({})
  showDynamicTips:boolean = false
  baseData:any = null
  startValidateDynamic:boolean = false
  nzDisabled:boolean = false
  listOfAlgorithm:SelectOption[] = [...algorithmList]
  positionList:SelectOption[]=[...positionList]
  listOfVerify:SelectOption[] = [...verifyList]

  labelHeaderTableBody:TBODY_TYPE[]= [...authLabelTableBody]

  labelHeaderList: Array<{ key: string; value: string }> = [
    { key: '', value: '' },
    { key: '', value: '' }
  ]

  // 判断动态组件内的必填数据是否填写
  canBeSave:boolean = false

  // eslint-disable-next-line camelcase
  createAuthForm: {position:string, uuid?:string, token_name:string, is_transparent:boolean, expire_time:number, expire_time_date:Date|null, driver:string, config:any} = {
    position: 'header',
    uuid: '',
    token_name: 'Authorization',
    is_transparent: false,
    expire_time: 0,
    expire_time_date: null,
    driver: 'basic',
    config: {
      user_name: '',
      password: '',
      apikey: '',
      ak: '',
      sk: '',
      iss: '',
      algorithm: '',
      secret: '',
      public_key: '',
      label: [{ key: '', value: '' }]
    }
  }

  constructor (
    private message: EoNgFeedbackMessageService,
    public api:ApiService,
    private fb: UntypedFormBuilder) {
    this.validateForm = this.fb.group({
      driver: ['basic', [Validators.required]],
      is_transparent: [false],
      position: ['header'],
      token_name: ['Authorization'],
      expire_time_date: [''],
      iss: [''],
      algorithm: ['HS256'],
      secret: [''],
      public_key: [''],
      user: [''],
      user_path: [''],
      claims_to_verify: [[]],
      signature_is_base64: [false]
    })
  }

  ngOnInit (): void {
    this.getDriversList()
    if (this.authId) {
      this.getAuthMessage()
    }
    this.initTable()
  }

  initTable () {
    this.labelHeaderTableBody[0].disabledFn = () => {
      return this.nzDisabled
    }
    this.labelHeaderTableBody[1].disabledFn = () => {
      return this.nzDisabled
    }

    this.labelHeaderTableBody[2].btns[0].disabledFn = () => {
      return this.nzDisabled
    }

    this.labelHeaderTableBody[2].showFn = (item: any) => {
      return item === this.labelHeaderList[0]
    }

    this.labelHeaderTableBody[3].btns[0].disabledFn = () => {
      return this.nzDisabled
    }

    this.labelHeaderTableBody[3].btns[1].disabledFn = () => {
      return this.nzDisabled
    }

    this.labelHeaderTableBody[3].showFn = (item: any) => {
      return item !== this.labelHeaderList[0]
    }
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  getAuthMessage () {
    this.api.get('application/auth', { app_id: this.appId, uuid: this.authId }).subscribe((resp:any) => {
      if (resp.code === 0) {
        this.createAuthForm = resp.data.auth
        this.validateForm.controls['driver'].setValue(resp.data.auth.driver)
        this.validateForm.controls['is_transparent'].setValue(resp.data.auth.is_transparent)
        this.validateForm.controls['position'].setValue(resp.data.auth.position)
        this.validateForm.controls['token_name'].setValue(resp.data.auth.token_name)
        const data = resp.data.auth
        if (this.createAuthForm.expire_time) {
          this.validateForm.controls['expire_time_date'].setValue(new Date(Number(data.expire_time) * 1000))
        }

        if (data?.config?.label && Object.keys(data.config.label)) {
          const tempLabel:Array<{key:string, value:any}> = []
          for (const index in Object.keys(data.config.label)) {
            tempLabel.push({ key: Object.keys(data.config.label)[index], value: data.config.label[Object.keys(data.config.label)[index]] })
          }
          if (tempLabel.length < 2) {
            tempLabel.push({ key: '', value: '' })
          }
          this.labelHeaderList = [...tempLabel]
          this.createAuthForm.config.label = tempLabel
        } else {
          this.createAuthForm.config.label = [{ key: '', value: '' }]
          this.labelHeaderList = [{ key: '', value: '' }, { key: '', value: '' }]
        }
        this.createAuthForm.config.user_name = this.createAuthForm.config?.user_name ? this.createAuthForm.config.user_name : ''
        this.createAuthForm.config.password = this.createAuthForm.config?.password ? this.createAuthForm.config.password : ''
        this.createAuthForm.config.apikey = this.createAuthForm.config?.apikey ? this.createAuthForm.config.apikey : ''
        this.createAuthForm.config.ak = this.createAuthForm.config?.ak ? this.createAuthForm.config.ak : ''
        this.createAuthForm.config.sk = this.createAuthForm.config?.sk ? this.createAuthForm.config.sk : ''
        this.createAuthForm.config.iss = this.createAuthForm.config?.iss ? this.createAuthForm.config.iss : ''
        this.createAuthForm.config.algorithm = this.createAuthForm.config?.algorithm ? this.createAuthForm.config.algorithm : 'HS256'
        this.createAuthForm.config.secret = this.createAuthForm.config?.secret ? this.createAuthForm.config.secret : ''
        this.createAuthForm.config.public_key = this.createAuthForm.config?.public_key ? this.createAuthForm.config.public_key : ''
        this.validateForm.controls['iss'].setValue(resp.data.auth.config?.iss || '')
        this.validateForm.controls['algorithm'].setValue(resp.data.auth.config?.algorithm || 'HS256')
        this.validateForm.controls['secret'].setValue(resp.data.auth.config?.secret || '')
        this.validateForm.controls['public_key'].setValue(resp.data.auth.config?.public_key || '')
        this.validateForm.controls['user'].setValue(resp.data.auth.config?.user || '')
        this.validateForm.controls['user_path'].setValue(resp.data.auth.config?.user_path || '')
        this.validateForm.controls['claims_to_verify'].setValue(resp.data.auth.config?.claims_to_verify || [])
        this.validateForm.controls['signature_is_base64'].setValue(resp.data.auth.config?.signature_is_base64 || false)
      } else {
        this.message.error(resp.msg || '获取数据失败!')
      }
    })
  }

  getDriversList () {
    this.driverList = []
    this.api.get('application/drivers').subscribe((resp:any) => {
      if (resp.code === 0) {
        for (const index in resp.data.drivers) {
          this.driverList = [...this.driverList, { label: this.getAuthDriver(resp.data.drivers[index].name), value: resp.data.drivers[index].name, render: resp.data.drivers[index].render }]
          if (this.validateForm.controls['driver'].value === resp.data.drivers[index].name) {
            this.baseData = resp.data.drivers[index].render
          }
        }
      } else {
        this.message.error(resp.msg || '获取列表数据失败!')
      }
    })
  }

  changeBasedata () {
    this.showDynamicTips = false
    for (const index in this.driverList) {
      if (this.driverList[index].value === this.validateForm.controls['driver'].value) {
        this.baseData = this.driverList[index].render
      }
    }
    this.getDataFromDynamicComponent(this.createAuthForm)
  }

  // 判断动态组件内的必填数据是否填写
  getDataFromDynamicComponent (value:any) {
    this.canBeSave = false
    switch (this.validateForm.controls['driver'].value) {
      case 'basic': {
        if (!value.config.user_name || !value.config.password) {
          this.canBeSave = false
          return
        }
        this.canBeSave = true
        this.showDynamicTips = false
        break
      }
      case 'apikey': {
        if (!value.config.apikey) {
          this.canBeSave = false
          return
        }
        this.canBeSave = true
        this.showDynamicTips = false
        break
      }

      case 'aksk': {
        if (!value.config.ak || !value.config.sk) {
          this.canBeSave = false
          return
        }
        this.canBeSave = true
        this.showDynamicTips = false
        break
      }
      case 'jwt':
        if (value.config.algorithm?.includes('HS') && (!value.config.iss || !value.config.secret)) {
          this.canBeSave = false

          break
        } else if ((value.config.algorithm?.includes('RS') || value.config.algorithm?.includes('ES')) && (!value.config.iss || !value.config.public_key)) {
          this.canBeSave = false

          break
        } else if (!value.config.algorithm || (!value.config.algorithm.includes('HS') && !value.config.algorithm.includes('ES') && !value.config.algorithm.includes('RS'))) {
          this.canBeSave = false
          this.showDynamicTips = true

          break
        }
        this.canBeSave = true
        break
    }
  }

  getAuthDriver (driver:string):string {
    switch (driver) {
      case 'basic':
        return 'Basic'
      case 'apikey':
        return 'ApiKey'
      case 'aksk':
        return 'AkSk'
      case 'jwt':
        return 'Jwt'
      default:
        return driver
    }
  }

  saveAuth () {
    this.startValidateDynamic = true
    this.showDynamicTips = !this.canBeSave
    if (this.validateForm.valid && (this.canBeSave || this.validateForm.controls['driver'].value === 'jwt')) {
      const tempLabel: Map<string, any> = new Map()
      let body:AuthData|undefined
      if (this.validateForm.controls['driver'].value !== 'jwt') {
        this.createAuthForm.expire_time = this.validateForm.controls['expire_time_date'].value ? Math.floor(new Date(this.validateForm.controls['expire_time_date'].value.setHours(23, 59, 59)).getTime() / 1000) : 0
        this.createAuthForm.config.public_key = this.createAuthForm.config.public_key === null ? '' : this.createAuthForm.config.public_key
        this.createAuthForm.config.secret = this.createAuthForm.config.secret === null ? '' : this.createAuthForm.config.secret
        if (this.createAuthForm.config.label?.length > 0) {
          for (const index in this.createAuthForm.config.label) {
            if (this.createAuthForm.config.label[index].key !== null && this.createAuthForm.config.label[index].key !== '' && this.createAuthForm.config.label[index].value !== null && this.createAuthForm.config.label[index].value !== '') {
              tempLabel.set(this.createAuthForm.config.label[index].key, this.createAuthForm.config.label[index].value)
            }
          }
        }
        const obj = Object.create(null)
        for (const [k, v] of tempLabel) {
          obj[k] = v
        }
        this.createAuthForm.config.label = obj
        body = {
          ...this.createAuthForm,
          driver: this.validateForm.controls['driver'].value,
          is_transparent: this.validateForm.controls['is_transparent'].value,
          expire_time: this.validateForm.controls['expire_time_date'].value ? Math.floor(new Date(this.validateForm.controls['expire_time_date'].value.setHours(23, 59, 59)).getTime() / 1000) : 0,
          position: this.validateForm.controls['position'].value,
          token_name: this.validateForm.controls['token_name'].value
        }
      } else {
        for (const label of this.labelHeaderList) {
          if (label.key && label.value) {
            tempLabel.set(label.key, label.value)
          }
        }

        const obj = Object.create(null)
        for (const [k, v] of tempLabel) {
          obj[k] = v
        }

        body = {
          driver: this.validateForm.controls['driver'].value,
          is_transparent: this.validateForm.controls['is_transparent'].value,
          expire_time: this.validateForm.controls['expire_time_date'].value ? Math.floor(new Date(this.validateForm.controls['expire_time_date'].value.setHours(23, 59, 59)).getTime() / 1000) : 0,
          position: this.validateForm.controls['position'].value,
          token_name: this.validateForm.controls['token_name'].value,
          config: {
            iss: this.validateForm.controls['iss'].value,
            algorithm: this.validateForm.controls['algorithm'].value,
            secret: this.validateForm.controls['secret'].value,
            public_key: this.validateForm.controls['public_key'].value,
            user: this.validateForm.controls['user'].value,
            user_path: this.validateForm.controls['user_path'].value,
            claims_to_verify: this.validateForm.controls['claims_to_verify'].value,
            signature_is_base64: this.validateForm.controls['signature_is_base64'].value,
            label: obj
          }
        }

        if (body.config.algorithm!.includes('HS')) {
          delete body.config.public_key
        } else {
          delete body.config.secret
          delete body.config.signature_is_base64
        }
      }
      if (this.authId) {
        this.api.put('application/auth', { ...body }, { app_id: this.appId, uuid: this.createAuthForm.uuid }).subscribe((resp:any) => {
          if (resp.code === 0) {
            this.message.success(resp.msg || '修改鉴权成功', { nzDuration: 1000 })
            this.closeModal(true)
            return true
          } else {
            this.message.error(resp.msg || '修改鉴权失败!')
            this.closeModal()

            return true
          }
        })
      } else {
        this.api.post('application/auth', { ...body }, { app_id: this.appId }).subscribe((resp:any) => {
          if (resp.code === 0) {
            this.message.success(resp.msg || '新增鉴权成功', { nzDuration: 1000 })
            this.closeModal(true)
            return true
          } else {
            this.message.error(resp.msg || '新增鉴权失败!')
            this.closeModal()
            return true
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
    return false
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  closeModal = (freshList?:boolean) => {}

  // Can not select days after today and today
  disabledDate = (current: Date): boolean =>
    differenceInCalendarDays(current, new Date()) < 0;
}

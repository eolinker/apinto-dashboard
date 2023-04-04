/* eslint-disable dot-notation */
import { Component, Input, OnInit } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { differenceInCalendarDays } from 'date-fns'
import { SelectOption } from 'eo-ng-select'
import { ApplicationAuthForm, AuthData } from '../../types/types'
import { algorithmList, authLabelTableBody, positionList, verifyList } from '../../types/conf'
import { TBODY_TYPE } from 'eo-ng-table'
import { ArrayItemData, EmptyHttpResponse } from 'projects/core/src/app/constant/type'
import { setFormValue } from 'projects/core/src/app/constant/form'

@Component({
  selector: 'eo-ng-application-authentication-form',
  templateUrl: './form.component.html',
  styles: [
  ]
})
export class ApplicationAuthenticationFormComponent implements OnInit {
  @Input() closeModal? :(value?:boolean)=>void
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
  canBeSave:boolean = false// 判断动态组件内的必填数据是否填写
  labelHeaderTableBody:TBODY_TYPE[]= [...authLabelTableBody]
  labelHeaderList:ArrayItemData[] = [
    { key: '', value: '' },
    { key: '', value: '' }
  ]

  // eslint-disable-next-line camelcase
  createAuthForm:ApplicationAuthForm = {
    position: 'header',
    uuid: '',
    tokenName: 'Authorization',
    isTransparent: false,
    expireTime: 0,
    expireTimeDate: null,
    driver: 'basic',
    config: {
      userName: '',
      password: '',
      apikey: '',
      ak: '',
      sk: '',
      iss: '',
      algorithm: '',
      secret: '',
      publicKey: '',
      label: [{ key: '', value: '' }]
    }
  }

  constructor (
    private message: EoNgFeedbackMessageService,
    public api:ApiService,
    private fb: UntypedFormBuilder) {
    this.validateForm = this.fb.group({
      driver: ['basic', [Validators.required]],
      isTransparent: [false],
      position: ['header'],
      tokenName: ['Authorization'],
      expireTimeDate: [''],
      iss: [''],
      algorithm: ['HS256'],
      secret: [''],
      publicKey: [''],
      user: [''],
      userPath: [''],
      claimsToVerify: [[]],
      signatureIsBase64: [false]
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
    this.labelHeaderTableBody[0].disabledFn = () => { return this.nzDisabled }
    this.labelHeaderTableBody[1].disabledFn = () => { return this.nzDisabled }
    this.labelHeaderTableBody[2].btns[0].disabledFn = () => { return this.nzDisabled }
    this.labelHeaderTableBody[2].showFn = (item: any) => { return item === this.labelHeaderList[0] }
    this.labelHeaderTableBody[3].btns[0].disabledFn = () => { return this.nzDisabled }
    this.labelHeaderTableBody[3].btns[1].disabledFn = () => { return this.nzDisabled }
    this.labelHeaderTableBody[3].showFn = (item: any) => { return item !== this.labelHeaderList[0] }
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  getAuthMessage () {
    this.api.get('application/auth', { appId: this.appId, uuid: this.authId }).subscribe((resp:any) => {
      if (resp.code === 0) {
        this.createAuthForm = resp.data.auth
        setFormValue(this.validateForm, {
          ...resp.data.auth,
          algorithm: 'HS256',
          claimsToVerify: [],
          ...resp.data.auth.config
        })
        const data = resp.data.auth
        if (this.createAuthForm.expireTime) {
          this.validateForm.controls['expireTimeDate'].setValue(new Date(Number(data.expireTime) * 1000))
        }

        if (data?.config?.label && Object.keys(data.config.label)) {
          const tempLabel:ArrayItemData[] = []
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
        this.createAuthForm.config.userName = this.createAuthForm.config?.userName ? this.createAuthForm.config.userName : ''
        this.createAuthForm.config.password = this.createAuthForm.config?.password ? this.createAuthForm.config.password : ''
        this.createAuthForm.config.apikey = this.createAuthForm.config?.apikey ? this.createAuthForm.config.apikey : ''
        this.createAuthForm.config.ak = this.createAuthForm.config?.ak ? this.createAuthForm.config.ak : ''
        this.createAuthForm.config.sk = this.createAuthForm.config?.sk ? this.createAuthForm.config.sk : ''
        this.createAuthForm.config.iss = this.createAuthForm.config?.iss ? this.createAuthForm.config.iss : ''
        this.createAuthForm.config.algorithm = this.createAuthForm.config?.algorithm ? this.createAuthForm.config.algorithm : 'HS256'
        this.createAuthForm.config.secret = this.createAuthForm.config?.secret ? this.createAuthForm.config.secret : ''
        this.createAuthForm.config.publicKey = this.createAuthForm.config?.publicKey ? this.createAuthForm.config.publicKey : ''
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
        if (!value.config.userName || !value.config.password) {
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
        } else if ((value.config.algorithm?.includes('RS') || value.config.algorithm?.includes('ES')) && (!value.config.iss || !value.config.publicKey)) {
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
        this.createAuthForm.expireTime = this.validateForm.controls['expireTimeDate'].value ? Math.floor(new Date(this.validateForm.controls['expireTimeDate'].value.setHours(23, 59, 59)).getTime() / 1000) : 0
        this.createAuthForm.config.publicKey = this.createAuthForm.config.publicKey === null ? '' : this.createAuthForm.config.publicKey
        this.createAuthForm.config.secret = this.createAuthForm.config.secret === null ? '' : this.createAuthForm.config.secret
        if (this.createAuthForm.config?.label && this.createAuthForm.config?.label?.length > 0) {
          for (const index in this.createAuthForm.config.label) {
            if (this.createAuthForm.config.label[index].key !== null && this.createAuthForm.config.label[index].key !== '' && this.createAuthForm.config.label[index].value !== null && this.createAuthForm.config.label[index].value !== '') {
              tempLabel.set(this.createAuthForm.config.label[index].key, this.createAuthForm.config.label[index].value)
            }
          }
        }
        const obj:{[k:string]:any} = Object.create(null)
        for (const [k, v] of tempLabel) {
          obj[k] = v
        }
        this.createAuthForm.config.label = obj
        body = {
          ...this.createAuthForm as AuthData,
          driver: this.validateForm.controls['driver'].value,
          isTransparent: this.validateForm.controls['isTransparent'].value,
          expireTime: this.validateForm.controls['expireTimeDate'].value ? Math.floor(new Date(this.validateForm.controls['expireTimeDate'].value.setHours(23, 59, 59)).getTime() / 1000) : 0,
          position: this.validateForm.controls['position'].value,
          tokenName: this.validateForm.controls['tokenName'].value
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
          isTransparent: this.validateForm.controls['isTransparent'].value,
          expireTime: this.validateForm.controls['expireTimeDate'].value ? Math.floor(new Date(this.validateForm.controls['expireTimeDate'].value.setHours(23, 59, 59)).getTime() / 1000) : 0,
          position: this.validateForm.controls['position'].value,
          tokenName: this.validateForm.controls['tokenName'].value,
          config: {
            iss: this.validateForm.controls['iss'].value,
            algorithm: this.validateForm.controls['algorithm'].value,
            secret: this.validateForm.controls['secret'].value,
            publicKey: this.validateForm.controls['publicKey'].value,
            user: this.validateForm.controls['user'].value,
            userPath: this.validateForm.controls['userPath'].value,
            claimsToVerify: this.validateForm.controls['claimsToVerify'].value,
            signatureIsBase64: this.validateForm.controls['signatureIsBase64'].value,
            label: obj
          }
        }

        if (body.config.algorithm!.includes('HS')) {
          delete body.config.publicKey
        } else {
          delete body.config.secret
          delete body.config.signatureIsBase64
        }
      }
      if (this.authId) {
        this.api.put('application/auth', { ...body }, { appId: this.appId, uuid: this.createAuthForm.uuid })
          .subscribe((resp:EmptyHttpResponse) => {
            if (resp.code === 0) {
              this.message.success(resp.msg || '修改鉴权成功', { nzDuration: 1000 })
              this.closeModal && this.closeModal(true)
            } else {
              this.closeModal && this.closeModal()
            }
          })
      } else {
        this.api.post('application/auth', { ...body }, { appId: this.appId })
          .subscribe((resp:EmptyHttpResponse) => {
            if (resp.code === 0) {
              this.message.success(resp.msg || '新增鉴权成功', { nzDuration: 1000 })
              this.closeModal && this.closeModal(true)
            } else {
              this.closeModal && this.closeModal()
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

  // Can not select days after today and today
  disabledDate = (current: Date): boolean =>
    differenceInCalendarDays(current, new Date()) < 0;
}

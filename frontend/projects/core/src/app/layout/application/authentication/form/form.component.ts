/* eslint-disable dot-notation */
import { Component, Input, OnInit } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { differenceInCalendarDays } from 'date-fns'
import { SelectOption } from 'eo-ng-select'
import { ApplicationAuthForm, AuthData } from '../../types/types'
import { algorithmList, positionList, verifyList } from '../../types/conf'
import { EmptyHttpResponse } from 'projects/core/src/app/constant/type'
import { setFormValue } from 'projects/core/src/app/constant/form'

@Component({
  selector: 'eo-ng-application-authentication-form',
  templateUrl: './form.component.html',
  styles: []
})
export class ApplicationAuthenticationFormComponent implements OnInit {
  @Input() closeModal?: (value?: boolean) => void
  appId: string = ''
  authId: string = ''
  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  driverList: Array<{ label: string; value: string; render: any }> = []
  validateForm: FormGroup = new FormGroup({})
  showDynamicTips: boolean = false
  baseData: any = null
  startValidateDynamic: boolean = false
  nzDisabled: boolean = false
  listOfAlgorithm: SelectOption[] = [...algorithmList]
  positionList: SelectOption[] = [...positionList]
  listOfVerify: SelectOption[] = [...verifyList]
  canBeSave: boolean = false // 判断动态组件内的必填数据是否填写
  oauth2OriginClientSecret: string = ''
  // eslint-disable-next-line camelcase
  createAuthForm: ApplicationAuthForm = {
    title: '',
    position: 'header',
    uuid: '',
    tokenName: 'Authorization',
    hideCredential: false,
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
      clientId: '',
      clientSecret: '',
      clientType: '',
      hashSecret: true,
      hashed: false,
      appId: '',
      appKey: '',
      // acceptHttpIfAlreadyTerminated: true,
      // enableClientCredentials: true,
      // enableAuthorizationCode: false,
      // enableImplicitGrant: false,
      // enablePasswordGrant: false,
      // tokenExpiration: 7200,
      // refreshTokenTtl: 7200,
      redirectUrls: [{ url: '' }],
      // scopes: [{ scope: '' }],
      // mandatoryScope: true,
      // provisionKey: '',
      // reuseRefreshToken: true,
      // persistentRefreshToken: true,
      // pkce: '',
      issuer: '',
      authenticatedGroupsClaim: [{ value: '' }]
      // enableMode: [
      //   {
      //     label: '授权码模式',
      //     value: 'enable_authorization_code',
      //     checked: true
      //   },
      //   {
      //     label: '客户端凭证模式',
      //     value: 'enable_client_credentials'
      //   },
      //   {
      //     label: '隐式授权模式',
      //     value: 'enable_implicit_grant'
      //   }
      // ]
    }
  }

  constructor (
    private message: EoNgFeedbackMessageService,
    public api: ApiService,
    private fb: UntypedFormBuilder
  ) {
    this.validateForm = this.fb.group({
      title: ['', [Validators.required]],
      driver: ['basic', [Validators.required]],
      hideCredential: [false],
      position: ['header'],
      tokenName: ['Authorization'],
      expireTimeDate: '',
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
  }

  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  getAuthMessage () {
    this.api
      .get('application/auth', { appId: this.appId, uuid: this.authId })
      .subscribe((resp: any) => {
        if (resp.code === 0) {
          this.createAuthForm = {
            ...resp.data.auth,
            config: {
              ...resp.data.auth.config,
              // enableMode: [
              //   {
              //     label: '授权码模式',
              //     value: 'enable_authorization_code',
              //     checked: resp.data.auth.config.enableAuthorizationCode
              //   },
              //   {
              //     label: '客户端凭证模式',
              //     value: 'enable_client_credentials',
              //     checked: resp.data.auth.config.enableClientCredentials
              //   },
              //   {
              //     label: '隐式授权模式',
              //     value: 'enable_implicit_grant',
              //     checked: resp.data.auth.config.enableImplicitGrant
              //   }
              // ],
              redirectUrls:
                resp.data.auth.config.redirectUrls?.map((x: string) => ({
                  url: x
                })) || [],
              // scopes: resp.data.auth.config.scopes?.map((x: string) => ({
              //   scope: x
              // })) || [],
              authenticatedGroupsClaim: resp.data.auth.config.authenticatedGroupsClaim?.map((x: string) => ({
                value: x
              })) || []
            }
          }

          setFormValue(this.validateForm, {
            ...resp.data.auth,
            algorithm: 'HS256',
            claimsToVerify: [],
            ...resp.data.auth.config,
            expireTimeDate: ''
          })

          const data = resp.data.auth
          if (this.createAuthForm.expireTime) {
            this.validateForm.controls['expireTimeDate'].setValue(
              new Date(Number(data.expireTime) * 1000)
            )
          }

          this.oauth2OriginClientSecret = resp.data.auth.config.clientSecret

          this.createAuthForm.config.userName = this.createAuthForm.config
            ?.userName
            ? this.createAuthForm.config.userName
            : ''
          this.createAuthForm.config.password = this.createAuthForm.config
            ?.password
            ? this.createAuthForm.config.password
            : ''
          this.createAuthForm.config.apikey = this.createAuthForm.config?.apikey
            ? this.createAuthForm.config.apikey
            : ''
          this.createAuthForm.config.ak = this.createAuthForm.config?.ak
            ? this.createAuthForm.config.ak
            : ''
          this.createAuthForm.config.sk = this.createAuthForm.config?.sk
            ? this.createAuthForm.config.sk
            : ''
          this.createAuthForm.config.iss = this.createAuthForm.config?.iss
            ? this.createAuthForm.config.iss
            : ''
          this.createAuthForm.config.algorithm = this.createAuthForm.config
            ?.algorithm
            ? this.createAuthForm.config.algorithm
            : 'HS256'
          this.createAuthForm.config.secret = this.createAuthForm.config?.secret
            ? this.createAuthForm.config.secret
            : ''
          this.createAuthForm.config.publicKey = this.createAuthForm.config
            ?.publicKey
            ? this.createAuthForm.config.publicKey
            : ''
        }
      })
  }

  getDriversList () {
    this.driverList = []
    this.api.get('application/drivers').subscribe((resp: any) => {
      if (resp.code === 0) {
        for (const index in resp.data.drivers) {
          this.driverList = [
            ...this.driverList,
            {
              label: this.getAuthDriver(resp.data.drivers[index].name),
              value: resp.data.drivers[index].name,
              render: resp.data.drivers[index].render
            }
          ]
          if (
            this.validateForm.controls['driver'].value ===
            resp.data.drivers[index].name
          ) {
            this.baseData = resp.data.drivers[index].render
          }
        }
        if (this.authId) {
          this.getAuthMessage()
        }
      }
    })
  }

  changeBasedata () {
    this.showDynamicTips = false
    for (const index in this.driverList) {
      if (
        this.driverList[index].value ===
        this.validateForm.controls['driver'].value
      ) {
        this.baseData = this.driverList[index].render
      }
    }
    this.getDataFromDynamicComponent(this.createAuthForm)
  }

  // 判断动态组件内的必填数据是否填写
  getDataFromDynamicComponent (value: any) {
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
        if (
          value.config.algorithm?.includes('HS') &&
          (!value.config.iss || !value.config.secret)
        ) {
          this.canBeSave = false

          break
        } else if (
          (value.config.algorithm?.includes('RS') ||
            value.config.algorithm?.includes('ES')) &&
          (!value.config.iss || !value.config.publicKey)
        ) {
          this.canBeSave = false

          break
        } else if (
          !value.config.algorithm ||
          (!value.config.algorithm.includes('HS') &&
            !value.config.algorithm.includes('ES') &&
            !value.config.algorithm.includes('RS'))
        ) {
          this.canBeSave = false
          this.showDynamicTips = true

          break
        }
        this.canBeSave = true
        break

      case 'oauth2': {
        if (
          !value.config.clientId ||
          !value.config.clientSecret ||
          !value.config.clientType
        ) {
          this.canBeSave = false
          return
        }
        this.canBeSave = true
        this.showDynamicTips = false
        break
      }

      case 'openid-connect-jwt': {
        if (
          !value.config.issuer ||
          !value.config.authenticatedGroupsClaim ||
          value.config.authenticatedGroupsClaim.length === 0
        ) {
          this.canBeSave = false
          return
        }
        this.canBeSave = true
        this.showDynamicTips = false
        break
      }

      case 'para-hmac': {
        if (
          !value.config.appId ||
          !value.config.appKey
        ) {
          this.canBeSave = false
          return
        }
        this.canBeSave = true
        this.showDynamicTips = false
        break
      }
    }
  }

  getAuthDriver (driver: string): string {
    switch (driver) {
      case 'basic':
        return 'Basic'
      case 'apikey':
        return 'ApiKey'
      case 'aksk':
        return 'AkSk'
      case 'jwt':
        return 'Jwt'
      case 'oauth2':
        return 'Oauth2'
      case 'openid-connect-jwt':
        return 'OpenID-Connect(JWT)'
      default:
        return driver
    }
  }

  buildRequestBody () {
    const formData = this.validateForm.controls
    this.createAuthForm.config.publicKey =
        this.createAuthForm.config.publicKey === null
          ? ''
          : this.createAuthForm.config.publicKey
    this.createAuthForm.config.secret =
        this.createAuthForm.config.secret === null
          ? ''
          : this.createAuthForm.config.secret
    const finalConfig = this.createAuthForm.config

    const body: AuthData = {
      title: formData['title'].value,
      driver: formData['driver'].value,
      hideCredential: formData['hideCredential'].value,
      expireTime: this.calculateExpireTime(formData['expireTimeDate'].value),
      position: formData['position'].value,
      tokenName: formData['tokenName'].value,
      config: {
        ...(this.createAuthForm.config as any),
        redirectUrls: Array.from(new Set(finalConfig.redirectUrls?.map(x => x.url))),
        // scopes: Array.from(new Set(finalConfig.scopes?.map(x => x.scope))),
        authenticatedGroupsClaim: Array.from(new Set(finalConfig.authenticatedGroupsClaim?.map(x => x.value)))
      }
    }

    if (formData['driver'].value === 'jwt') {
      const formValue = this.validateForm.value
      delete formValue.title
      delete formValue.driver
      delete formValue.tokenName
      delete formValue.position
      body.config = {
        ...body.config,
        ...formValue
      }
    }

    if (formData['driver'].value === 'jwt' && body.config.algorithm?.includes('HS')) {
      delete body.config.publicKey
    } else {
      delete body.config.secret
      delete body.config.signatureIsBase64
    }
    return body
  }

  calculateExpireTime (expireTimeDate: Date | null): number {
    return expireTimeDate ? Math.floor(new Date(expireTimeDate.setHours(23, 59, 59)).getTime() / 1000) : 0
  }

  saveAuth () {
    this.startValidateDynamic = true
    this.showDynamicTips = !this.canBeSave
    // 解决出现未输入值但二次切换算法时，表单校验不通过的bug
    Object.values(this.validateForm.controls).forEach((control) => {
      if (control.invalid) {
        control.markAsDirty()
        control.updateValueAndValidity({ onlySelf: true })
      }
    })
    this.validateForm.updateValueAndValidity()
    if (
      this.validateForm.valid &&
        (this.canBeSave || this.validateForm.controls['driver'].value === 'jwt')
    ) {
      const body = this.buildRequestBody()
      const params = this.authId ? { appId: this.appId, uuid: this.createAuthForm.uuid } : { appId: this.appId }

      this.api[this.authId ? 'put' : 'post']('application/auth', body, params).subscribe((resp: EmptyHttpResponse) => {
        if (resp.code === 0) {
          const successMessage = this.authId ? '修改鉴权成功' : '新增鉴权成功'
          this.message.success(resp.msg || successMessage, { nzDuration: 1000 })
          this.closeModal && this.closeModal(true)
        } else {
          this.closeModal && this.closeModal()
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

  // Can not select days after today and today
  disabledDate = (current: Date): boolean =>
    differenceInCalendarDays(current, new Date()) < 0
}

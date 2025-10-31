import { CheckboxOptionType } from 'antd'
import { ArrayItemData, PublishStatus } from '../../../constant/type'

/* eslint-disable camelcase */
export interface ApplicationEnum {
  name: string
  id: string
}

export interface ApplicationListData {
  name: string
  id: string
  desc: string
  operator: string
  updateTime: string
  isDelete: boolean
  publish: Array<{ name: string; title: string; status: PublishStatus }>
  [k: string]: any
}

export type ApplicationParamData = {
  key: string
  value: string
  conflict: string
  position: string
  eoKey?: string
  conflictString?: string
}
export interface ApplicationData {
  name: string
  id: string
  desc: string
  customAttrList: ArrayItemData[]
  params?: ApplicationParamData[]
}

export interface ApplicationAuthForm {
  title: string
  position: string
  uuid?: string
  tokenName: string
  hideCredential: boolean
  expireTime: number
  expireTimeDate: Date | null
  driver: string
  config: {
    appId?:string
    appKey?:string
    userName?: string
    password?: string
    apikey?: string
    ak?: string
    sk?: string
    iss?: string
    algorithm?: string
    secret?: string
    publicKey?: string
    user?: string
    userPath?: string
    claimsToVerify?: Array<string>
    signatureIsBase64?: boolean
    hideCredential?: boolean
    label?:
      | Array<{ key: string; value: string | number }>
      | { [k: string]: any }
    redirectUrls?: Array<{ url: string }>
    scopes?:Array<{ scope: string }>
    mandatoryScope?:boolean,
    provisionKey?:string,
    clientId?: string
    clientSecret?: string
    clientType?: string
    hashSecret?: boolean
    hashed?: boolean
    acceptHttpIfAlreadyTerminated?: boolean
    // enableClientCredentials?: boolean
    // enableAuthorizationCode?: boolean
    // enableImplicitGrant?: boolean
    // enablePasswordGrant?: boolean
    tokenExpiration?: number
    refreshTokenTtl?: number
    reuseRefreshToken?: boolean
    persistentRefreshToken?: boolean
    pkce?: string
    // enableMode?: Array<{ label: string; value: string; checked?: boolean }>
    issuer?:string
    authenticatedGroupsClaim?:Array<{ value: string }>
  }
}

export interface AuthData {
  title: string
  driver: 'basic' | 'apikey' | 'aksk' | 'jwt'
  hideCredential: boolean
  expireTime: number
  position: string
  tokenName: string
  config: {
    userName?: string
    password?: string
    apikey?: string
    ak?: string
    sk?: string
    iss?: string
    algorithm?: string
    secret?: string
    publicKey?: string
    user?: string
    userPath?: string
    claimsToVerify?: Array<string>
    signatureIsBase64?: boolean
    hideCredential?: boolean
    // label?: { [key: string]: string }
    redirectUrls?: string[]
    scopes?:string[],
    mandatoryScope?:boolean,
    provisionKey?:string,
    clientId?: string
    clientSecret?: string
    clientType?: string
    hashSecret?: boolean
    hashed?: boolean
    acceptHttpIfAlreadyTerminated?: boolean
    // enableClientCredentials?: boolean
    // enableAuthorizationCode?: boolean
    // enableImplicitGrant?: boolean
    // enablePasswordGrant?: boolean
    tokenExpiration?: number
    refreshTokenTtl?: number
    reuseRefreshToken?: boolean
    persistentRefreshToken?: boolean
    pkce?: string
    // enableMode?: any
    issuer?:string
    authenticatedGroupsClaim?:string[]
  }
}

export interface AuthListData {
  uuid: string
  driver: string
  hideCredential: boolean | string
  expireTime: number
  expireTimeString?: string
  paramPosition: string
  paramName: string
  paramInfo: string
  operator: string
  updateTime: string
  ruleInfo: string
}

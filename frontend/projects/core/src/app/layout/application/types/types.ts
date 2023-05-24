import { ArrayItemData, PublishStatus } from '../../../constant/type'

/* eslint-disable camelcase */
export interface ApplicationEnum{
    name:string
    id:string
}

export interface ApplicationListData{
    name:string
    id:string
    desc:string
    operator:string
    updateTime:string
    isDelete:boolean
    publish:Array<{name:string, title:string, status:PublishStatus}>
    [k:string]:any
}

export type ApplicationParamData = {
    key:string
    value:string
    conflict:string
    position:string
}
export interface ApplicationData{
    name:string
    id:string
    desc:string
    customAttrList:ArrayItemData[]
    params?:ApplicationParamData[]
}

export interface ApplicationAuthForm{
    name:string
    position:string
    uuid?:string
    tokenName:string
    isTransparent:boolean
    expireTime:number
    expireTimeDate:Date|null
    driver:string
    config:{
        userName?:string
        password?:string
        apikey?:string
        ak?:string
        sk?:string
        iss?:string
        algorithm?:string
        secret?:string
        publicKey?:string
        user?:string
        userPath?:string
        claimsToVerify?:Array<string>
        signatureIsBase64?:boolean
        hideCredential?:boolean
        label?:Array<{key:string, value:string|number}>|{[k:string]:any}
    }
}

export interface AuthData{
    name:string
    driver:'basic' | 'apikey' | 'aksk' | 'jwt'
    isTransparent:boolean
    expireTime:number
    position:string
    tokenName:string
    config:{
        userName?:string
        password?:string
        apikey?:string
        ak?:string
        sk?:string
        iss?:string
        algorithm?:string
        secret?:string
        publicKey?:string
        user?:string
        userPath?:string
        claimsToVerify?:Array<string>
        signatureIsBase64?:boolean
        hideCredential?:boolean
        label?:{[key:string]:string}
    }
}

export interface AuthListData{
    uuid:string
    driver:string
    isTransparent:boolean|string
    expireTime:number
    expireTimeString?:string,
    paramPosition:string
    paramName:string
    paramInfo:string
    operator:string
    updateTime:string
    ruleInfo:string
}

export interface ExtraListData{
    key:string
    value:string
    conflict:string
    position:string
    conflictString?:string
}

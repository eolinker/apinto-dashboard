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
    eoKey?:string
    conflictString?:string
}
export interface ApplicationData{
    name:string
    id:string
    desc:string
    customAttrList:ArrayItemData[]
    params?:ApplicationParamData[]
}

export interface ApplicationAuthForm{
    title:string
    position:string
    uuid?:string
    tokenName:string
    hideCredential:boolean
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
    title:string
    driver:'basic' | 'apikey' | 'aksk' | 'jwt'
    hideCredential:boolean
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
    hideCredential:boolean|string
    expireTime:number
    expireTimeString?:string,
    paramPosition:string
    paramName:string
    paramInfo:string
    operator:string
    updateTime:string
    ruleInfo:string
}

/* eslint-disable camelcase */
export interface ApplicationEnum{
    name:string
    id:string
}

export interface AuthData{
    driver:'basic' | 'apikey' | 'aksk' | 'jwt'
    is_transparent:boolean
    expire_time:number
    position:string
    token_name:string
    config:{
        user_name?:string
        password?:string
        apikey?:string
        ak?:string
        sk?:string
        iss?:string
        algorithm?:string
        secret?:string
        public_key?:string
        user?:string
        user_path?:string
        claims_to_verify?:Array<string>
        signature_is_base64?:boolean
        hide_credential?:boolean
        label?:{[key:string]:string}
    }
}

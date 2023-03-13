/* eslint-disable camelcase */
export interface UserData{
    sex?:number
    avatar?:string
    email:string
    phone:string
    user_name:string
    nick_name:string
    role_ids:Array<string>
    desc:string
    notice_user_id:string
}

export interface UserListData extends UserData{
    id:number
    status:number
    last_login:string
    create_time:string
    update_time:string
    operate_disable:boolean
    operator:string
}

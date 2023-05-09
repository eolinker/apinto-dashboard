// email 相关参数
export interface EmailData{
  uuid:string
  smtpUrl:string
  smtpPort:string
  protocol:string
  email:string
  account:string
  password:string
}

// webhook 相关参数
export interface WebhookListData{
    uuid:string
    title:string
    url:string
    method:string
    contentType:string
    isDelete:boolean
    operator:string
    updateTime:string
    createTime:string
    [key:string]:any
}

export interface WebhookData{
    uuid?:string
    title:string
    desc:string
    url:string
    method:string
    contentType:string
    noticeType:string
    userSeparator?:string
    header:Map<string, string> | {[key:string]:string}
    template:string
    [key:string]:any
}

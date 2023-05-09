import { Operator } from '../../../constant/type'
export interface AuditLogDetail{
    attr:string
     value:string
}

export interface AuditLogsData{
    id:number,
    operator:Operator,
    operateType: string,
    kind:string,
    time:string
    ip:string,
    [k:string]:any
  }

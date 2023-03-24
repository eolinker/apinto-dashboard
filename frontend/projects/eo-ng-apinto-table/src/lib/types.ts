import { TemplateRef } from '@angular/core'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'

export interface EO_TBODY_TYPE extends TBODY_TYPE{
    check?:Function
    tooltip?:string | TemplateRef<any>
    json?:boolean
}

export interface EO_THEAD_TYPE extends THEAD_TYPE{
    required?:boolean
    titleString?:string // require用
}

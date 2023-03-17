import { TemplateRef } from '@angular/core'
import { TBODY_TYPE } from 'eo-ng-table'

export interface EO_TBODY_TYPE extends TBODY_TYPE{
    tooltip?:string | TemplateRef<any>
}

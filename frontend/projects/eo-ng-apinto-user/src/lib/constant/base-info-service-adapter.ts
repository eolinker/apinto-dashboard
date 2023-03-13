/* eslint-disable camelcase */

import { InjectionToken } from '@angular/core'
export interface RouteParams {
    // 用户角色ID
    roleId:string
    operate_disable:string
  }

export const BASEINFO_SERVICE_ADAPTER:

// eslint-disable-next-line no-use-before-define
InjectionToken<BaseInfoServiceAdapter> = new InjectionToken<BaseInfoServiceAdapter>('BaseInfoServiceAdapter')
export interface BaseInfoServiceAdapter {
    allParamsInfo:RouteParams
}

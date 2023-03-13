
import { InjectionToken } from '@angular/core'
export const APP_CONFIG_ADAPTER:
// eslint-disable-next-line no-use-before-define
InjectionToken<AppConfigAdapter> = new InjectionToken<AppConfigAdapter>('AppConfigAdapter')

export const MODAL_NORMAL_SIZE: number = 900
export const MODAL_SMALL_SIZE: number = 600
export const MODAL_LARGE_SIZE: number = 1200

export interface AppConfigAdapter {
    menuList:Array<any>
}

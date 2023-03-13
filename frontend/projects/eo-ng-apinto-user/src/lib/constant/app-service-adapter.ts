
import { InjectionToken } from '@angular/core'
import { Observable } from 'rxjs'
export const APP_SERVICE_ADAPTER:
// eslint-disable-next-line no-use-before-define
InjectionToken<AppServiceAdapter> = new InjectionToken<AppServiceAdapter>('AppServicegAdapter')

export interface AppServiceAdapter {
    dataUpdated:boolean;
    checkUpdateRight(router:string): Observable<any>;
    repUpdateRightList():Observable<Array<string>>;
    getRightsList(): Observable<any>;
    getUpdateRightsRouter():Array<string>;
    setUserRoleId(id:string):void
    getUserRoleId() :string
    setUserId(id:string):void
    getUserId() :string
    getMenuList():Observable<Array<string>>;
}

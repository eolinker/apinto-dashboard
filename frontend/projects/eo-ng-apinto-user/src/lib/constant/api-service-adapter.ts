
import { InjectionToken } from '@angular/core'
import { Observable } from 'rxjs'

export const API_SERVICE_ADAPTER:
// eslint-disable-next-line no-use-before-define
InjectionToken<ApiServiceAdapter> = new InjectionToken<ApiServiceAdapter>('apiServiceAdapter')

export interface ApiServiceAdapter {
    get (url: string, params?: {[key:string]:any}): Observable<any>;
    post (url: string, body?: any, params?: {[key:string]:any}): Observable<any>;
    put (url:string, body?:any, params?: {[key:string]:any}): Observable<any>;
    delete (url:string, body?: {[key:string]:any}):Observable<any>;
    patch (url:string, body?:any, params?: {[key:string]:any}): Observable<any>;
    login(body?:any, params?: {[key:string]:any}): Observable<any>;
    logout(body?:any, params?: {[key:string]:any}): Observable<any>;
}

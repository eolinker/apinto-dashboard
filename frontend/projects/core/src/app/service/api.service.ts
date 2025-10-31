/* eslint-disable dot-notation */
/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2022-07-30 00:40:51
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2024-01-16 15:42:06
 * @FilePath: /apinto/src/app/service/api.service.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import {
  HttpClient,
  HttpErrorResponse,
  HttpHeaders,
  HttpParams
} from '@angular/common/http'
import { Inject, Injectable, InjectionToken } from '@angular/core'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { catchError, Observable, throwError } from 'rxjs'

export const API_URL = new InjectionToken<string>('apiUrl')
@Injectable({
  providedIn: 'root'
})
export class ApiService {
  constructor(
    private message: EoNgFeedbackMessageService,
    private http: HttpClient,
    @Inject(API_URL) public urlPrefix: string
  ) {}

  // 登录接口
  login(body?: any, params?: { [key: string]: any }) {
    if (params) {
      params['namespace'] = 'default'
    } else {
      params = { namespace: 'default' }
    }
    if (params && params['query']) {
      params['query'] = JSON.stringify(params['query'])
    }
    const p = new HttpParams({
      fromObject: params
    })
    return this.http
      .post(this.urlPrefix + 'sso/login', body, {
        params: p,
        withCredentials: true
      })
      .pipe(catchError(this.handleError))
  }

  // 检查cookie是否合法,不合法则需要登录
  checkAuth(body?: any, params?: { [key: string]: any }) {
    if (params) {
      params['namespace'] = 'default'
    } else {
      params = { namespace: 'default' }
    }
    if (params && params['query']) {
      params['query'] = JSON.stringify(params['query'])
    }
    const p = new HttpParams({
      fromObject: params
    })
    return this.http
      .post(this.urlPrefix + 'sso/login/check', body, {
        params: p,
        withCredentials: true
      })
      .pipe(catchError(this.handleError))
  }

  // 退出登录
  logout(body?: any, params?: { [key: string]: any }) {
    if (params) {
      params['namespace'] = 'default'
    } else {
      params = { namespace: 'default' }
    }
    if (params && params['query']) {
      params['query'] = JSON.stringify(params['query'])
    }
    const p = new HttpParams({
      fromObject: params
    })
    return this.http
      .post(this.urlPrefix + 'sso/logout', body, {
        params: p,
        withCredentials: true
      })
      .pipe(catchError(this.handleError))
  }

  // 商业授权中激活时需要上传文件
  authPostWithFile(
    url: string,
    body?: any,
    params?: { [key: string]: any }
  ): Observable<any> {
    for (const index in params) {
      if (typeof params[index] === 'string') {
        params[index] = params[index].trim()
      }
    }

    for (const index in body) {
      if (typeof body[index] === 'string') {
        body[index] = body[index].trim()
      }
    }

    if (params) {
      params['namespace'] = 'default'
    } else {
      params = { namespace: 'default' }
    }
    const headers = new HttpHeaders()
    return this.http
      .post(this.urlPrefix + '_system/' + url, body, {
        headers,
        params: params,
        withCredentials: true
      })
      .pipe(catchError(this.handleError))
  }

  // 商业授权相关的get接口
  authGet(url: string, params?: { [key: string]: any }): Observable<any> {
    if (params) {
      params['namespace'] = 'default'
    } else {
      params = { namespace: 'default' }
    }
    if (params && params['query']) {
      params['query'] = JSON.stringify(params['query'])
    }

    for (const index in params) {
      if (typeof params[index] === 'string') {
        params[index] = params[index].trim()
      }
    }

    const p = new HttpParams({
      fromObject: params
    })
    return this.http
      .get(this.urlPrefix + '_system/' + url, {
        params: p,
        withCredentials: true
      })
      .pipe(
        // retry(3),

        catchError(this.handleError)
      )
  }

  get(url: string, params?: { [key: string]: any }): Observable<any> {
    if (params && params['query']) {
      params['query'] = JSON.stringify(params['query'])
    }

    for (const index in params) {
      if (typeof params[index] === 'string') {
        params[index] = params[index].trim()
      }
    }

    const p = new HttpParams({
      fromObject: params
    })
    return this.http
      .get(this.urlPrefix + 'api/' + url, {
        params: p,
        withCredentials: true
      })
      .pipe(
        // retry(3),
        catchError(this.handleError)
      )
  }

  post(
    url: string,
    body?: any,
    params?: { [key: string]: any }
  ): Observable<any> {
    for (const index in params) {
      if (typeof params[index] === 'string') {
        params[index] = params[index].trim()
      }
    }

    if (body && !(body instanceof FormData)) {
      for (const index in body) {
        if (typeof body[index] === 'string') {
          body[index] = body[index].trim()
        }
      }
    }
    return this.http
      .post(this.urlPrefix + 'api/' + url, body, {
        params: params,
        withCredentials: true
      })
      .pipe(catchError(this.handleError))
  }

  put(
    url: string,
    body?: any,
    params?: { [key: string]: any }
  ): Observable<any> {
    for (const index in params) {
      if (typeof params[index] === 'string') {
        params[index] = params[index].trim()
      }
    }

    for (const index in body) {
      if (typeof body[index] === 'string') {
        body[index] = body[index].trim()
      }
    }
    return this.http
      .put(this.urlPrefix + 'api/' + url, body, {
        params: params,
        withCredentials: true
      })
      .pipe(catchError(this.handleError))
  }

  delete(url: string, params?: { [key: string]: any }): Observable<any> {
    for (const index in params) {
      if (typeof params[index] === 'string') {
        params[index] = params[index].trim()
      }
    }

    return this.http
      .delete(this.urlPrefix + 'api/' + url, { params: params })
      .pipe(catchError(this.handleError))
  }

  patch(
    url: string,
    body?: any,
    params?: { [key: string]: any }
  ): Observable<any> {
    for (const index in params) {
      if (typeof params[index] === 'string') {
        params[index] = params[index].trim()
      }
    }

    for (const index in body) {
      if (typeof body[index] === 'string') {
        body[index] = body[index].trim()
      }
    }
    return this.http
      .patch(this.urlPrefix + 'api/' + url, body, {
        params: params,
        withCredentials: true
      })
      .pipe(catchError(this.handleError))
  }

  handleError = (error: HttpErrorResponse) => {
    if (error.status === 0) {
      // A client-side or network error occurred. Handle it accordingly.
      console.error('An error occurred:', error.error)
    } else {
      // The backend returned an unsuccessful response code.
      // The response body may contain clues as to what went wrong.
      console.error(
        `Backend returned code ${error.status}, body was: `,
        error.error
      )
    }
    console.log(error)
    if (error?.error?.msg) {
      this.message.error(error?.error?.msg)
    }
    // Return an observable with a user-facing error message.
    return throwError(
      () => new Error('Something bad happened; please try again later.')
    )
  }
}

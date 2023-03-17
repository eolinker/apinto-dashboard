/* eslint-disable no-useless-constructor */
import { Injectable } from '@angular/core'
import {
  HttpRequest,
  HttpHandler,
  HttpEvent,
  HttpInterceptor,
  HttpResponse
} from '@angular/common/http'
import { tap, Observable } from 'rxjs'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'

@Injectable()
export class ErrorInterceptor implements HttpInterceptor {
  private loadingMessageId:string = ''
  constructor (
    private message: EoNgFeedbackMessageService) {}

  intercept (request: HttpRequest<unknown>, next: HttpHandler): Observable<HttpEvent<unknown>> {
    return next.handle(request).pipe(
      tap((event:any) => {
        // this.hideLoader()
        if (event instanceof HttpResponse) {
          this.checkAccess(event.body.code, event)
          if (request.url.includes('monitor') || request.url.includes('warn') || request.url.includes('user/enum')) {
            event.body.data = this.camel(event.body.data)
          }
        }
      }
      )
    )
  }

  // 下划线转驼峰
  camel (data:any):any {
    if (typeof data !== 'object' || !data) return data
    if (Array.isArray(data)) {
      return (data as Array<any>).map((item:any) => { return this.camel(item) })
    }
    const newData:any = {}
    for (const key in data) {
      const newKey = key.replace(/_([a-z0-9])/g, (p, m) => m.toUpperCase())
      newData[newKey] = this.camel(data[key])
    }
    return newData
  }

  // 根据后端返回的code判断是否要提示无权限弹窗或跳转路由
  checkAccess (code:number, responseBody:any) {
    if (responseBody.url.includes('warn/') && code !== 0) {
      this.message.error(responseBody.body.msg || '操作失败！')
    }
  }
}

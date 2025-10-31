import { Injectable } from '@angular/core'
import {
  HttpRequest,
  HttpHandler,
  HttpEvent,
  HttpInterceptor,
  HttpResponse,
  HttpParams
} from '@angular/common/http'
import { Observable, from, mergeMap, map, of } from 'rxjs'
import { EoNgNavigationService } from '../eo-ng-navigation.service'
import { PluginEventHubService } from '../plugin-event-hub.service'
import { EoNgMessageService } from '../eo-ng-message.service'

// 判断对象是否需要做变量名转化，如果是表单对象、文件流等就不做处理
const isJsonObject:(obj: any)=> boolean = (obj:any) => {
  // 首先确保它是一个对象
  if (typeof obj !== 'object' || obj === null) {
    return false
  }

  // 检查对象不是特殊的 HTTP 相关类型
  return !(obj instanceof FormData ||
           obj instanceof Blob ||
           obj instanceof File ||
           obj instanceof ArrayBuffer ||
           obj instanceof URLSearchParams ||
           // 检查其他可能的类型
           obj instanceof ReadableStream ||
           obj instanceof FileList)
}

// 驼峰转下划线,其中监控的status_4xx和status_5xx需要特殊处理
export const underline:(data:any)=>any = (data:any) => {
  if (typeof data !== 'object' || !data) return data
  if (!isJsonObject(data)) { return data };
  if (Array.isArray(data)) {
    return data.map(item => underline(item))
  }
  const newData:any = {}
  for (const key in data) {
    // 首字母不参与转换
    let newKey = key[0] + key.substring(1).replace(/([A-Z])/g, (p, m) => `_${m.toLowerCase()}`
    )
    newKey = key === 'status4xx' ? 'status_4xx' : (key === 'status5xx' ? 'status_5xx' : newKey)
    newData[newKey] = underline(data[key])
  }
  return newData
}

export const underlinedStr:(str:string)=>string = (str:string) => {
  return str[0] + str.substring(1).replace(/([A-Z])/g, (p, m) => `_${m.toLowerCase()}`)
}

@Injectable()
export class ErrorInterceptor implements HttpInterceptor {
  authStatus:'normal'|'waring'|'freeze' = 'normal'
  constructor (
    private navigationService: EoNgNavigationService,
    private pluginEventHub:PluginEventHubService,
    private message:EoNgMessageService) {}

  intercept (request: HttpRequest<unknown>, next: HttpHandler): Observable<HttpEvent<unknown>> {
    const copiedReq = this.handleRequest(request)
    return next.handle(copiedReq).pipe(
      mergeMap((event:any) => {
        if (event instanceof HttpResponse) {
          return from(this.pluginEventHub.initHub()!.emit('httpResponse', { data: { req: request, res: event } })).pipe(
            map((res:any) => {
              event = res.res

              if (!request.url.includes('api/dynamic/')) {
                try {
                  event.body.data = this.camel(event.body.data)
                } catch {
                  console.warn(' 转化接口数据命名法出现问题')
                }
              }

              if (!request.url.startsWith('/remote') && event.body.code !== undefined && !(request.url.includes('sso/login/check')) && event.body.code !== 0 && event.body.code !== 30001) {
                let msg = event.body.msg
                if (event.url.includes('router/online') && request.method === 'PUT') {
                  msg = event.body.data.router.map((data:any) => {
                    return data.msg
                  }).join('  ')
                }
                this.message.error(msg || '操作失败！')
              }
              return event
            }))
        }

        return of(event)
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
      const newKey = (key === 'status_4xx' || key === 'status_5xx') ? key : key.replace(/_([a-z0-9])/g, (p, m) => m.toUpperCase())
      newData[newKey] = this.camel(data[key])
    }
    return newData
  }

  underlineParams (params: HttpParams) {
    let newParams = params.set('namespace', 'default')
    params.keys().forEach(key => {
      const newKey = underlinedStr(key)
      newParams = newParams.set(newKey, underline(params.get(key))!)
    })
    return newParams
  }

  handleRequest (request: HttpRequest<unknown>) {
    if (request.url.startsWith('dynamic') || request.url.startsWith('/api/dynamic')) {
      return request
    }
    const modifiedParams = this.underlineParams(request.params)
    const modifiedBody = underline(request.body)
    return request.clone({ params: modifiedParams, body: modifiedBody })
  }
}

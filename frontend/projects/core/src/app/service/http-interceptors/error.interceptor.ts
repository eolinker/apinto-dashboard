import { Component, Injectable } from '@angular/core'
import {
  HttpRequest,
  HttpHandler,
  HttpEvent,
  HttpInterceptor,
  HttpResponse
} from '@angular/common/http'
import { tap, Observable } from 'rxjs'
import { Router } from '@angular/router'
import { NzModalService } from 'ng-zorro-antd/modal'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { EoNgNavigationService } from '../eo-ng-navigation.service'
import { environment } from 'projects/core/src/environments/environment'

@Injectable()
export class ErrorInterceptor implements HttpInterceptor {
  authStatus:'normal'|'waring'|'freeze' = 'normal'
  constructor (
    private router: Router,
    private navigationService: EoNgNavigationService,
    private modalService:NzModalService,
    private eoModalService:EoNgFeedbackModalService,
    private message: EoNgFeedbackMessageService) {}

  intercept (request: HttpRequest<unknown>, next: HttpHandler): Observable<HttpEvent<unknown>> {
    return next.handle(request).pipe(
      tap((event:any) => {
        // this.hideLoader()
        if (event instanceof HttpResponse) {
          this.checkAuthStatus(event)
          this.checkAccess(event.body.code, event, request.method)
          if (!request.url.includes('api/dynamic/')) {
            try {
              event.body.data = this.camel(event.body.data)
            } catch {
              console.warn('转化接口数据命名法出现问题')
            }
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

  checkAuthStatus (event:HttpResponse<any>) {
    if (event.headers && event.headers.get('X-Apinto-Auth-Status') && this.authStatus !== event.headers.get('X-Apinto-Auth-Status')) {
      this.authStatus = event.headers.get('X-Apinto-Auth-Status') as 'normal' | 'waring' | 'freeze'
      this.navigationService.reqCheckAuthStatus()
    }
  }

  // 根据后端返回的code判断是否要提示无权限弹窗或跳转路由
  checkAccess (code:number, responseBody:any, requestMethod:string) {
    switch (code) {
      case -2:
        this.modalService.closeAll()
        this.eoModalService.closeAll()
        this.openAccessModal()
        break
      case -3:
        setTimeout(() => {
          this.modalService.closeAll()
          this.eoModalService.closeAll()
          if (!this.router.url.includes('/login')) {
            this.router.navigate(['/', 'login'], { queryParams: { callback: this.router.url }, queryParamsHandling: 'merge' })
          }
        }, 1000)
        break
      case -7:
        if (environment.isBusiness) {
          this.eoModalService.closeAll()
          this.modalService.closeAll()
          setTimeout(() => {
            if (!responseBody.url.includes('create_check')) {
              this.router.navigate(['/', 'auth'])
            }
          }, 1000)
        }
        break
      default:
        if (!this.router.url.startsWith('/remote') && code !== undefined && !(responseBody.url.includes('sso/login/check')) && code !== 0 && code !== 30001) {
          let msg = responseBody.body.msg
          if (responseBody.url.includes('router/online') && requestMethod === 'PUT') {
            msg = responseBody.body.data.router.map((data:any) => {
              return data.msg
            }).join('  ')
          }
          this.message.error(msg || '操作失败！')
        }
    }
  }

  openAccessModal () {
    this.modalService.confirm({
      nzWrapClassName: 'modal-header',
      nzTitle: '权限提示',
      nzIconType: 'exclamation-circle',
      nzContent: ModalContentComponent,
      nzClosable: true,
      nzOkText: '确定',
      nzCancelText: '取消',
      nzOnOk: () => {
        const mainPageUrl = this.navigationService.getPageRoute()
        if (mainPageUrl) {
          this.router.navigate([this.navigationService.getPageRoute()])
        }
      }
    })
  }
}

@Component({
  selector: 'modal-content',
  template: `
    <div class="modal-header">
    <p>无法获取您当前账号的相关权限信息，请确认是否赋予权限。</p>
    <p>具体信息，请咨询<span class='blue'>管理员</span></p>
    </div>
  `,
  styles: [
    `
    .blue{
      color:blue;
    }
    .modal-header{
      margin-top:24px;
    }
  `
  ]
})
export class ModalContentComponent {
}

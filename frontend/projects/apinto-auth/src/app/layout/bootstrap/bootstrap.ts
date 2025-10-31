/*
 * @Date: 2023-12-12 18:57:19
 * @LastEditors: maggieyyy
 * @LastEditTime: 2023-12-15 17:34:48
 * @FilePath: \apinto\projects\apinto-auth\src\app\layout\bootstrap\bootstrap.ts
 */
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { Router } from '@angular/router'
import { Subject } from 'rxjs'
import { Injector } from '@angular/core'
import { AuthButtonComponent } from '../../component/auth-button/auth-button.component'
import { AuthButtonModule } from '../../component/auth-button/auth-button.module'

export const checkAuthStatus: Subject<boolean> = new Subject<boolean>()

// 当本项目作为插件导入apinto-dashboard时，本模块里的bootstrap方法作为立即执行插件的立即执行函数被执行
// 采用module的方式，user项目被引入时，每个模块创建的实例只有一个，模块之间的通信通过Angular服务实现，保持了Angular的一致性和可维护性
export async function bootstrap (props: any): Promise<void> {
  const checkAccess = (code:number, responseBody:any, requestMethod:string, closeModal:Function, router:Router, messageService:EoNgFeedbackMessageService) => {
    switch (code) {
      case -7:
        closeModal()
        setTimeout(() => {
          if (!responseBody.url.includes('create_check')) {
            router.navigate(['/', 'auth'])
          } else {
            router.navigate(['/', 'auth-info'])
          }
        }, 1000)
        return false
      default:
        return true
    }
  }

  // 使用这些服务执行你需要的操作
  const { pluginEventHub, pluginSlotHub, platformProvider, closeModal, router, messageService } = props

  // pluginProvider.setRouterConfig(false, {
  //   path: 'auth-info',
  //   component: AuthInfoComponent
  // })

  pluginEventHub.on('httpResponse', (eventData:any) => {
    if (eventData.headers && eventData.headers.get('X-Apinto-Auth-Status') && eventData.authStatus !== eventData.headers.get('X-Apinto-Auth-Status')) {
      checkAuthStatus.next(true)
    }
    const continueRes = checkAccess(eventData.res.body.code, eventData.res, eventData.req.method, closeModal, router, messageService)
    return {
      data: eventData,
      continue: continueRes
    }
  })

  const ModuleRef = await platformProvider.getPlatformRef().bootstrapModule(AuthButtonModule, { ngZone: platformProvider.getNgZone() })
    .catch((error: any) => {
      console.error('Bootstrap error:', error)
      if (error instanceof Error) {
        console.error('Error message:', error.message)
        console.error('Error stack:', error.stack)
      } else {
        console.error('Error details:', error)
      }
      throw error
    })
  const ModuleInject: Injector = ModuleRef.injector
  pluginSlotHub.addSlot('renderAuthButton', [AuthButtonComponent, { ModuleInject }])
}

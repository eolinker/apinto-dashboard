/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2023-04-13 23:14:11
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2023-04-14 00:00:43
 * @FilePath: /apinto/projects/core/src/app/service/redirect-page.service.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { Injectable } from '@angular/core'
import { CanActivate, Router } from '@angular/router'
import { Observable, of } from 'rxjs'
import { EoNgNavigationService } from './app-config.service'

@Injectable({
  providedIn: 'root'
})
export class RedirectPageService implements CanActivate, CanActivate {
  // eslint-disable-next-line  no-useless-constructor
  constructor (
    private router: Router,
    private navigationService: EoNgNavigationService
  ) {}

  canActivate (): Observable<boolean> {
    if (!this.router.routerState.snapshot.url) {
      return new Observable((observer) => {
        this.navigationService.getMenuList().subscribe(() => {
          const pageRouter = this.navigationService.getPageRoute()
          if (pageRouter) {
            this.router.navigate([this.navigationService.getPageRoute()])
            observer.next(true)
          }
        })
        observer.next(true)
      })
    } else {
      return of(true)
    }
  }
}

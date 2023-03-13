import { Injectable } from '@angular/core'
import { CanActivate, Router } from '@angular/router'
import { Observable, of } from 'rxjs'
import { AppConfigService } from './app-config.service'

@Injectable({
  providedIn: 'root'
})
export class RedirectPageService implements CanActivate, CanActivate {
  // eslint-disable-next-line  no-useless-constructor
  constructor (private router: Router, private appConfigService: AppConfigService) {}

  canActivate (): Observable<boolean> {
    if (!this.router.routerState.snapshot.url) {
      return new Observable(observer => {
        this.appConfigService.getMenuList().subscribe(() => {
          const pageRouter = this.appConfigService.getPageRoute()
          if (pageRouter) {
            this.router.navigate([this.appConfigService.getPageRoute()])
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

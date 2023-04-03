import { Injectable } from '@angular/core'
import { CanActivate, Router } from '@angular/router'
import { Observable } from 'rxjs'
import { ApiService } from './api.service'
import { AppConfigService } from './app-config.service'

@Injectable({
  providedIn: 'root'
})
export class AuthGuardService implements CanActivate {
  // eslint-disable-next-line  no-useless-constructor
  constructor (private router: Router,
              private api: ApiService,
              private appConfigService: AppConfigService) {}

  canActivate (): Observable<boolean> {
    return new Observable(observer => {
      this.api.checkAuth().subscribe((resp:any) => {
        if (resp.code === 0) {
        // this.getInitMenuList(AppConfig.menuList)
          this.appConfigService.getMenuList().subscribe(() => {
            this.router.navigate([this.appConfigService.getPageRoute()])
            observer.next(true)
          })
        } else {
          observer.next(true)
        }
      })
    })
  }
  // return of(false)
}

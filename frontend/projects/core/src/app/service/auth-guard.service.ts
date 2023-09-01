import { Injectable } from '@angular/core'
import { CanActivate, Router } from '@angular/router'
import { Observable } from 'rxjs'
import { ApiService } from './api.service'
import { EoNgNavigationService } from './eo-ng-navigation.service'

@Injectable({
  providedIn: 'root'
})
export class AuthGuardService implements CanActivate {
  constructor (private router: Router,
              private api: ApiService,
              private navigationService: EoNgNavigationService) {}

  canActivate (): Observable<boolean> {
    return new Observable(observer => {
      this.api.checkAuth().subscribe((resp:any) => {
        if (resp.code === 0) {
        // this.getInitMenuList(AppConfig.menuList)
          this.navigationService.getMenuList().subscribe(() => {
            this.router.navigate([this.navigationService.getPageRoute()])
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

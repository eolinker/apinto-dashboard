/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2023-04-13 23:14:10
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2023-04-13 23:30:16
 * @FilePath: /apinto/projects/core/src/app/layout/login/login.component.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
/* eslint-disable no-useless-constructor */
import { Component, OnInit } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { Subscription } from 'rxjs'
import { ApiService } from '../../service/api.service'
import { EoNgNavigationService } from '../../service/eo-ng-navigation.service'
import { environment } from 'projects/core/src/environments/environment'

@Component({
  selector: 'eo-ng-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit {
  isBusiness:boolean = environment.isBusiness
  private subscription: Subscription = new Subscription()
  constructor (
    private navigationService: EoNgNavigationService,
    private api: ApiService,
    private router: Router
  ) {}

  ngOnInit () {
    this.api.checkAuth().subscribe((resp: any) => {
      if (resp.code === 0) {
        this.subscription = this.navigationService
          .getMenuList()
          .subscribe(() => {
            this.router.navigate([this.navigationService.getPageRoute()], { queryParamsHandling: 'merge' })
          })
      }
    })
  }

  ngOnDestroy () {
    this.subscription.unsubscribe()
  }
}

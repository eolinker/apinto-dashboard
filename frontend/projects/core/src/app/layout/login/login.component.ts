/* eslint-disable no-useless-constructor */
import { Component, OnInit } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { Subscription } from 'rxjs'
import { ApiService } from '../../service/api.service'
import { AppConfigService } from '../../service/app-config.service'

@Component({
  selector: 'eo-ng-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit {
  private subscription: Subscription = new Subscription()
  constructor (
    private appConfigService: AppConfigService,
    private api: ApiService,
    private router:Router,
    private message: EoNgFeedbackMessageService
  ) { }

  ngOnInit () {
    this.api.checkAuth().subscribe((resp:any) => {
      if (resp.code === 0) {
        this.subscription = this.appConfigService.getMenuList().subscribe(() => {
          this.router.navigate([this.appConfigService.getPageRoute()])
        })
      }
    })
  }

  ngOnDestroy () {
    this.subscription.unsubscribe()
  }
}

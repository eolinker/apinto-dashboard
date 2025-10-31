import { Component } from '@angular/core'
import { Router } from '@angular/router'
import { ApiService } from '../../service/api.service'
import { checkAuthStatus } from '../../layout/bootstrap/bootstrap'
import { Subscription } from 'rxjs'

@Component({
  selector: 'auth-auth-button',
  template: `
      <button
          *ngIf="showAuthBtn()"
          nzType="primary"
          [nzDanger]="true"
          eo-ng-button
          class="mr-btnbase"
          (click)="goToAuth()"
          eoNgFeedbackTooltip
          [nzTooltipTitle]="btnTooltip"
          nzTooltipPlacement="bottom"
          [nzTooltipVisible]="false"
          nzTooltipTrigger="hover"
        >
          {{ btnLabel }}
        </button>
  `,
  styles: [
    ''
  ]
})
export class AuthButtonComponent {
  authStatus:'normal' | 'waring' | 'freeze' = 'normal'
  btnLabel:string = ''
  btnTooltip:string = ''

  constructor (private router:Router, private api:ApiService) {}
  private subAuthCheck: Subscription = new Subscription()

  ngOnInit () {
    this.checkAuthStatus()
    this.subAuthCheck = checkAuthStatus.asObservable().subscribe(() => {
      this.checkAuthStatus()
    })
  }

  ngOnDestroy () {
    this.subAuthCheck.unsubscribe()
  }

  showAuthBtn () {
    return !this.router.url.includes('auth/info') && this.authStatus && this.authStatus !== 'normal'
  }

  goToAuth () {
    this.router.navigate(['/', 'auth', 'info'])
  }

  checkAuthStatus () {
    this.api.authGet('activation/check').subscribe((resp:{code:number, msg:string, data:{status:'normal'|'waring'|'freeze', prompt:string, label:string}}) => {
      if (resp.code === 0) {
        this.authStatus = resp.data.status
        this.btnLabel = resp.data.label
        this.btnTooltip = resp.data.prompt
        if (resp.data.status === 'freeze') {
          this.router.navigate(['/', 'auth', 'info'])
        }
      }
    })
  }
}

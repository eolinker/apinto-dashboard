/* eslint-disable no-useless-constructor */
import { Component, OnInit } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { ApiService } from '../../../service/api.service'
import { AppConfigService } from '../../../service/app-config.service'

@Component({
  selector: 'eo-ng-auth-info',
  templateUrl: './info.component.html',
  styles: [
  ]
})
export class AuthInfoComponent implements OnInit {
  authInfo:{title:string, infos:Array<{key:string, value:string}>}
  = { title: '标准版授权', infos: [] }

  constructor (
    private message: EoNgFeedbackMessageService,
    private api:ApiService,
    private router:Router,
    private appConfigService:AppConfigService) {
  }

  ngOnInit (): void {
    this.getInfo()
  }

  getInfo () {
    this.api.authGet('activation/info').subscribe((resp:{code:number, data:{infos:Array<{key:string, value:string}>, title:string}, msg:string}) => {
      if (resp.code === 0) {
        this.authInfo = resp.data
        this.appConfigService.reqFlashBreadcrumb([{ title: resp.data.title }])
      } else {
        this.message.error(resp.msg || '获取信息失败，请重试！')
      }
    })
  }

  updateAuth () {
    this.router.navigate(['/', 'auth-update'])
  }
}

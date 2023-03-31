/*
 * @Author:
 * @Date: 2022-08-17 23:42:52
 * @LastEditors:
 * @LastEditTime: 2022-08-24 00:31:06
 * @FilePath: /apinto/src/app/layout/application/application-message/application-message.component.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { Component, OnInit } from '@angular/core'
import { Router } from '@angular/router'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { BaseInfoService } from '../../../service/base-info.service'

@Component({
  selector: 'eo-ng-application-message',
  templateUrl: './message.component.html',
  styles: [
  ]
})
export class ApplicationMessageComponent implements OnInit {
  appId:string = ''

  constructor (public api:ApiService,
    private baseInfo:BaseInfoService,
     private router:Router) {
  }

  ngOnInit (): void {
    this.appId = this.baseInfo.allParamsInfo.appId
    if (!this.appId) {
      this.router.navigate(['/', 'application'])
    }
  }
}

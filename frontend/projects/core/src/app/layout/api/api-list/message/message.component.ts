/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2022-08-25 22:41:39
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2022-09-20 22:09:52
 * @FilePath: /apinto/src/app/layout/api/api-message/api-message.component.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
/* eslint-disable dot-notation */
import { Component } from '@angular/core'
import { Router } from '@angular/router'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { Subscription } from 'rxjs'

@Component({
  selector: 'eo-ng-api-message',
  templateUrl: './message.component.html',
  styles: [
  ]
})
export class ApiMessageComponent {
  readonly nowUrl:string = this.router.routerState.snapshot.url
  apiUuid:string = ''

  constructor (
    private baseInfo:BaseInfoService,
    public api:ApiService,
    private router:Router
  ) {
    this.apiUuid = this.baseInfo.allParamsInfo.apiId
    if (!this.apiUuid) {
      this.router.navigate(['/'])
    }
  }
}

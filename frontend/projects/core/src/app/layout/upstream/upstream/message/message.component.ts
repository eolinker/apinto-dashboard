/*
 * @Author:
 * @Date: 2022-08-17 23:42:52
 * @LastEditors:
 * @LastEditTime: 2022-08-23 17:46:34
 * @FilePath: /apinto/src/app/layout/upstream/upstream/upstream-message/upstream-message.component.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
/* eslint-disable dot-notation */
import { Component, OnInit } from '@angular/core'
import { Router, ActivatedRoute } from '@angular/router'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'

@Component({
  selector: 'eo-ng-upstream-message',
  templateUrl: './message.component.html',
  styles: [
  ]
})
export class UpstreamMessageComponent implements OnInit {
  readonly nowUrl:string = this.router.routerState.snapshot.url
  serviceName:string = ''

  constructor (private baseInfo:BaseInfoService, public api:ApiService, private appConfigService:AppConfigService, private router:Router, private activateInfo:ActivatedRoute) {

  }

  ngOnInit (): void {
    this.serviceName = this.baseInfo.allParamsInfo.serviceName
    if (!this.serviceName) {
      this.router.navigate(['/'])
    }
  }
}

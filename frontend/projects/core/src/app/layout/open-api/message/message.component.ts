/* eslint-disable dot-notation */
import { Component } from '@angular/core'
import { Router, ActivatedRoute } from '@angular/router'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'

@Component({
  selector: 'eo-ng-external-app-message',
  template: `
    <eo-ng-external-app-create [editPage]="true" [appId]="appId"></eo-ng-external-app-create>

  `,
  styles: [
  ]
})
export class ExternalAppMessageComponent {
  readonly nowUrl:string = this.router.routerState.snapshot.url
  appId:string = ''

  constructor (private baseInfo:BaseInfoService, public api:ApiService, private router:Router, private activateInfo:ActivatedRoute) {
    if (!this.baseInfo.allParamsInfo.extAppId) {
      this.router.navigate(['/'])
    } else {
      this.appId = this.baseInfo.allParamsInfo.extAppId
    }
  }

  ngOnDestroy () {
  }
}

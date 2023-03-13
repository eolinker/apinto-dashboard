/* eslint-disable dot-notation */
import { Component, Input } from '@angular/core'
import { Router, ActivatedRoute } from '@angular/router'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'

@Component({
  selector: 'eo-ng-grey-message',
  template: `
  <eo-ng-grey-create
  [editPage]="true"
  [clusterName]="clusterName"
  [strategyUuid]="strategyUuid"
></eo-ng-grey-create>
  `,
  styles: [
  ]
})
export class GreyMessageComponent {
  readonly nowUrl:string = this.router.routerState.snapshot.url
  @Input() strategyUuid:string = ''
  @Input() clusterName:string = ''

  constructor (private baseInfo:BaseInfoService,
    public api:ApiService, private router:Router, private activateInfo:ActivatedRoute) {
  }

  ngOnInit (): void {
    this.strategyUuid = this.baseInfo.allParamsInfo.strategyId
    if (!this.strategyUuid) {
      this.router.navigate(['/'])
    }
  }

}

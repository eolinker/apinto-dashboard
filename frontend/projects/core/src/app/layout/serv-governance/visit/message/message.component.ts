/* eslint-disable dot-notation */
import { Component, Input, OnInit } from '@angular/core'
import { Router } from '@angular/router'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'

@Component({
  selector: 'eo-ng-visit-message',
  template: `
  <eo-ng-visit-create
  [editPage]="true"
  [clusterName]="clusterName"
  [strategyUuid]="strategyUuid"
></eo-ng-visit-create>
  `,
  styles: [
  ]
})
export class VisitMessageComponent implements OnInit {
  @Input() strategyUuid:string = ''
  @Input() clusterName:string = ''

  constructor (
    private baseInfo:BaseInfoService, public api:ApiService, private router:Router) {
  }

  ngOnInit (): void {
    this.strategyUuid = this.baseInfo.allParamsInfo.strategyId
    if (!this.strategyUuid) {
      this.router.navigate(['/'])
    }
  }
}

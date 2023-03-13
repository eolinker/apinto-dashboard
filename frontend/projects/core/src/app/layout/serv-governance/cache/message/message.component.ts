/* eslint-disable dot-notation */
import { Component, Input, OnInit } from '@angular/core'
import { Router, ActivatedRoute } from '@angular/router'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'

@Component({
  selector: 'eo-ng-cache-message',
  template: `
  <eo-ng-cache-create
  [editPage]="true"
  [clusterName]="clusterName"
  [strategyUuid]="strategyUuid"
></eo-ng-cache-create>
  `,
  styles: [
  ]
})
export class CacheMessageComponent implements OnInit {
  readonly nowUrl:string = this.router.routerState.snapshot.url
  @Input() strategyUuid:string = ''
  @Input() clusterName:string = ''

  constructor (public api:ApiService,
    private baseInfo:BaseInfoService,
     private router:Router, private activateInfo:ActivatedRoute) {
  }

  ngOnInit (): void {
    this.strategyUuid = this.baseInfo.allParamsInfo.strategyId
    if (!this.strategyUuid) {
      this.router.navigate(['/'])
    }
  }

}

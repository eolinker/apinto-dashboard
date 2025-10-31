import { Component, Input, OnInit } from '@angular/core'
import { Router } from '@angular/router'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'

@Component({
  selector: 'eo-ng-data-mask-message',
  template: `
  <eo-ng-data-mask-create
    [editPage]="true"
    [clusterName]="clusterName"
    [strategyUuid]="strategyUuid" >
</eo-ng-data-mask-create>
  `,
  styles: [
  ]
})
export class DataMaskMessageComponent implements OnInit {
  @Input() strategyUuid:string = ''
  @Input() clusterName:string = ''

  constructor (
    private baseInfo:BaseInfoService,
     private router:Router) {
  }

  ngOnInit (): void {
    this.strategyUuid = this.baseInfo.allParamsInfo.strategyId
    if (!this.strategyUuid) {
      this.router.navigate(['/'])
    }
  }
}

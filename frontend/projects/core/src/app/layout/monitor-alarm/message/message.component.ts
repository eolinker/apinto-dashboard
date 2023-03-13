/* eslint-disable dot-notation */
import { Component, OnInit } from '@angular/core'
import { BaseInfoService } from '../../../service/base-info.service'

@Component({
  selector: 'eo-ng-monitor-alarm-message',
  template: `
    <eo-ng-monitor-alarm-config [editPage]="true" [partitionId]="partitionId"></eo-ng-monitor-alarm-config>
  `,
  styles: [
  ]
})
export class MonitorAlarmMessageComponent implements OnInit {
  partitionId:string = ''
  // eslint-disable-next-line no-useless-constructor
  constructor (private baseInfo:BaseInfoService) {
  }

  ngOnInit (): void {
    this.partitionId = this.baseInfo.allParamsInfo.partitionId
  }
}

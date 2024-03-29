import { Component, Input, OnInit } from '@angular/core'
import { Router, ActivatedRoute } from '@angular/router'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'

@Component({
  selector: 'eo-ng-traffic-message',
  templateUrl: './message.component.html',
  styles: [
  ]
})
export class TrafficMessageComponent implements OnInit {
  @Input() strategyUuid:string = ''
  @Input() clusterName:string = ''

  constructor (
    private baseInfo:BaseInfoService, public api:ApiService, private router:Router, private activateInfo:ActivatedRoute) {
  }

  ngOnInit (): void {
    this.strategyUuid = this.baseInfo.allParamsInfo.strategyId
    if (!this.strategyUuid) {
      this.router.navigate(['/'])
    }
  }
}

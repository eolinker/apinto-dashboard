/* eslint-disable dot-notation */
import { Component, OnInit } from '@angular/core'
import { Router, ActivatedRoute } from '@angular/router'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'

@Component({
  selector: 'eo-ng-service-discovery-message',
  templateUrl: './message.component.html',
  styles: [
  ]
})
export class ServiceDiscoveryMessageComponent implements OnInit {
  serviceName:string = ''

  constructor (private router:Router, private baseInfo:BaseInfoService, private activateInfo:ActivatedRoute) {
  }

  ngOnInit (): void {
    this.serviceName = this.baseInfo.allParamsInfo.discoveryName
    if (!this.serviceName) {
      this.router.navigate(['/'])
    }
  }

}

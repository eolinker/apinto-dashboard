
import { Component, OnInit } from '@angular/core'
import { Router, ActivatedRoute } from '@angular/router'

@Component({
  selector: 'eo-ng-monitor-alarm-area-service',
  templateUrl: './service.component.html',
  styles: [`
  .group-search-large{
    margin-bottom:16px;
  }
  .label{
    width:57px;
    display:inline-block;
    white-space:nowrap;
  }`
  ]
})
export class MonitorAlarmAreaUpstreamComponent implements OnInit {
  showTotal:boolean = true
  constructor (private router:Router, private activeInfo: ActivatedRoute) {
  }

  ngOnInit (): void {
  }

  ngDoCheck () {
    this.showTotal = this.router.url.split('?')[0].split('/').length !== 7 && this.router.url.split('?')[0].split('/')[5] !== 'detail'
  }
}

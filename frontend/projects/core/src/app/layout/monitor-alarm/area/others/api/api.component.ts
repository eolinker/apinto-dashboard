/* eslint-disable no-useless-constructor */
import { Component, OnInit } from '@angular/core'
import { Router } from '@angular/router'

@Component({
  selector: 'eo-ng-monitor-alarm-area-api',
  templateUrl: './api.component.html',
  styles: [
    `
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
export class MonitorAlarmAreaApiComponent implements OnInit {
  showTotal:boolean = true
  constructor (private router:Router) {
  }

  ngOnInit (): void {

  }

  ngDoCheck () {
    this.showTotal = this.router.url.split('?')[0].split('/').length !== 7 && this.router.url.split('?')[0].split('/')[5] !== 'detail'
  }
}

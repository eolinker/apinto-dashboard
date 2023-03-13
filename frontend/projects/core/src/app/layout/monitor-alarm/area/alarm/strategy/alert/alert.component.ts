import { Component, Input, OnInit } from '@angular/core'

@Component({
  selector: 'eo-ng-monitor-alarm-strategy-alert',
  template: `
    <div>
    请选择<a
      [routerLink]="['/', 'system', 'email']"
      (click)="closeModal && closeModal()"
      >配置邮箱告警</a
    >或<a
      [routerLink]="['/', 'system', 'webhook']"
      (click)="closeModal && closeModal()"
      >配置Webhook</a
    >告警渠道</div>
  `,
  styles: [
  ]
})
export class MonitorAlarmStrategyAlertComponent implements OnInit {
  @Input() closeModal?:(value?:any)=>void
  ngOnInit (): void {
  }
}

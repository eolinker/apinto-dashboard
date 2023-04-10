import { Component, EventEmitter, Input, Output } from '@angular/core'

@Component({
  selector: 'eo-ng-auth-info-detail',
  template: `
  <ol class="activation-info">
    <li *ngFor="let info of eoInfos">
      {{ info.key }} ：{{ info.value }}
    </li>
  </ol>
  <a class="auth-a" eoNgUserAccess="auth-info" (click)="updateAuth()"
    >更新授权</a
  >
  `,
  styles: [
    `
    ol {
      color: #666666;
      font-size: 14px;
      font-weight: 500;
      line-height: 26px;
      list-style-type: none;
      padding-inline-start: 0px;
      margin-bottom: 12px !important;
    }
    a {
      font-size: 14px;
      font-weight: 400;
      line-height: 22px;
    }
    `
  ]
})
export class AuthInfoDetailComponent {
  @Input() eoInfos:Array<{key:string, value:string}> = []
  @Output() eoUpdateAuth:EventEmitter<any> = new EventEmitter()

  updateAuth () {
    this.eoUpdateAuth.emit()
  }
}

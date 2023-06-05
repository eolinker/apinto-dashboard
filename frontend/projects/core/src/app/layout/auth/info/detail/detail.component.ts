import { Component, EventEmitter, Input, Output } from '@angular/core'

@Component({
  selector: 'eo-ng-auth-info-detail',
  template: `
  <ol
        class="text-DESC_TEXT block m-auto mt-[40px] text-MAIN_TEXT text-[14px] font-medium leading-[26px] list-none"
        style="padding-inline-start: 0px">
    <li *ngFor="let info of eoInfos">
      {{ info.key }} ：{{ info.value }}
    </li>
  </ol>
  <a class="m-auto mt-[24px] text-center" eoNgUserAccess="auth-info" (click)="updateAuth()"
    >更新授权</a
  >
  `,
  styles: [
    `
    ol {
      color: var(--TITLE_TEXT);
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

    :host{
        display: inline-block;
        width: auto;
        margin: auto;
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

import { Component, Input } from '@angular/core'

@Component({
  selector: 'eo-ng-api-method-tag',
  template: `
  <ng-container [ngSwitch]="methodItem">
    <span *ngSwitchCase="'PUT'" class="method-put-bg" [ngClass]="{'method':!inner, 'method-inner':inner}"
      >PUT</span
    >
    <span *ngSwitchCase="'POST'" class="method-post-bg"  [ngClass]="{'method':!inner, 'method-inner':inner}"
      >POST</span
    >
    <span *ngSwitchCase="'GET'" class="method-get-bg" [ngClass]="{'method':!inner, 'method-inner':inner}"
      >GET</span
    >
    <span *ngSwitchCase="'DELETE'" class="method-del-bg" [ngClass]="{'method':!inner, 'method-inner':inner}"
      >DELETE</span
    >
    <span *ngSwitchCase="'PATCH'" class="method-patch-bg" [ngClass]="{'method':!inner, 'method-inner':inner}"
      >PATCH</span
    >
    <span *ngSwitchCase="'HEAD'" class="method-head-bg" [ngClass]="{'method':!inner, 'method-inner':inner}"
      >HEAD</span
    >
    <span *ngSwitchCase="'OPTIONS'" class="method-opts-bg" [ngClass]="{'method':!inner, 'method-inner':inner}"
      >OPTIONS</span
    >
    <span
     *ngSwitchCase="'ALL'"
      class="method-all-bg" [ngClass]="{'method':!inner, 'method-inner':inner}"
      >ALL</span
    >
</ng-container>
  `,
  styles: [
    `
    .method {
      font-size: 12px !important;
      font-family: 'Helvetica Neue', 'Helvetica', 'PingFang SC', 'Hiragino Sans GB',
        'Microsoft YaHei', 'Noto Sans CJK SC', 'WenQuanYi Micro Hei', 'Arial',
        sans-serif;
      display: inline-block;
      border-radius: var(--DEFAULT_BORDER_RADIUS);
      height: 20px !important;
      line-height: 12px !important;
      text-align: center;
      padding: 4px 6px !important;
    }
    .method-inner{
      margin:2px
    }
    .method-put-bg {
      color: #d8830c !important;
      background-color: rgba(216, 131, 12, 0.15);
      border: 1px solid rgba(216, 131, 12, 0.15);
    }

    .method-post-bg {
      color: #10a54b !important;
      background-color: rgba(16, 165, 75, 0.15);
      border: 1px solid rgba(16, 165, 75, 0.15);
    }

    .method-get-bg {
      color: #067ddb !important;
      background-color: rgba(6, 125, 219, 0.15);
      border: 1px solid rgba(6, 125, 219, 0.15);
    }

    .method-del-bg {
      color: #c2161b !important;
      background-color: rgba(194, 22, 27, 0.15);
      border: 1px solid rgba(194, 22, 27, 0.15);
    }

    .method-all-bg {
      color: #7728f5 !important;
      background-color: rgba(119, 40, 245, 0.15);
      border: 1px solid rgba(119, 40, 245, 0.15);
    }

    .method-opts-bg {
      color: #0e5ab3 !important;
      background-color: rgba(14, 90, 179, 0.15);
      border: 1px solid rgba(14, 90, 179, 0.15);
    }

    .method-head-bg {
      color: #eec40c !important;
      background-color: rgba(238, 196, 12, 0.15);
      border: 1px solid rgba(238, 196, 12, 0.15);
    }

    .method-patch-bg {
      color: #ed863a !important;
      background-color: rgba(237, 134, 58, 0.15);
      border: 1px solid rgba(237, 134, 58, 0.15);
    }
`
  ]
})
export class ApiMethodTagComponent {
  @Input() methodItem:'PUT'|'POST'|'GET'|'DELETE'|'PATCH'|'HEAD'|'OPTIONS'|'ALL'|''=''
  @Input() inner:boolean = false
}

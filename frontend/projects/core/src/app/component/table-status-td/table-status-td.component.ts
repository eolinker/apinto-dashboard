import { Component, Input } from '@angular/core'

@Component({
  selector: 'eo-ng-table-status',
  template: `
    <ng-container [ngSwitch]="status">
    <span *ngSwitchCase="'TOUPDATE'" class="blue-bold">待更新</span>
    <span *ngSwitchCase="'待更新'" class="blue-bold">待更新</span>
    <span *ngSwitchCase="'待发布'" class="blue-bold">待发布</span>
    <span *ngSwitchCase="'GOONLINE'" class="green-bold">已上线</span>
    <span *ngSwitchCase="'已上线'" class="green-bold">已上线</span>
    <span *ngSwitchCase="'已发布'" class="green-bold">已发布</span>
    <span *ngSwitchCase="'OFFLINE'" class="grey-bold">已下线</span>
    <span *ngSwitchCase="'已下线'" class="grey-bold">已下线</span>
    <span *ngSwitchCase="'NOTGOONLINE'" class="grey-bold">未上线</span>
    <span *ngSwitchCase="'未上线'" class="grey-bold">未上线</span>
    <span *ngSwitchCase="'未发布'" class="grey-bold">未发布</span>
    <span *ngSwitchCase="'TODELETE'" class="orange-bold">待删除</span>
    </ng-container>
  `,
  styles: [
  ]
})
export class TableStatusTdComponent {
    @Input() status:string | undefined
}

@Component({
  selector: 'eo-ng-table-disabled',
  template: `
    <ng-container [ngSwitch]="disable">
      <span *ngSwitchCase="true" class="red-bold">已禁用</span>
      <span *ngSwitchCase="false" class="green-bold">未禁用</span>
    </ng-container>
    `,
  styles: [
  ]
})
export class TableDisabledStatusTdComponent {
      @Input() disable:string | undefined
}

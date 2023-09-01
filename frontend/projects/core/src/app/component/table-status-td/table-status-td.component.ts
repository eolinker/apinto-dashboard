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
    <span *ngSwitchCase="'放行'" class="green-bold">放行</span>
    <span *ngSwitchCase="'OFFLINE'" class="grey-bold">已下线</span>
    <span *ngSwitchCase="'已下线'" class="grey-bold">已下线</span>
    <span *ngSwitchCase="'NOTGOONLINE'" class="grey-bold">未上线</span>
    <span *ngSwitchCase="'未上线'" class="grey-bold">未上线</span>
    <span *ngSwitchCase="'未发布'" class="grey-bold">未发布</span>
    <span *ngSwitchCase="'TODELETE'" class="orange-bold">待删除</span>
    <span *ngSwitchCase="'拦截'" class="orange-bold">拦截</span>
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

@Component({
  selector: 'eo-ng-table-running',
  template: `
    <ng-container [ngSwitch]="running">
      <span *ngSwitchCase="'NOTRUNNING'" class="red-bold">未运行</span>
      <span *ngSwitchCase="'RUNNING'" class="green-bold">运行中</span>
    </ng-container>
    `,
  styles: [
  ]
})
export class TableRunningStatusTdComponent {
      @Input() running:string | undefined
}

@Component({
  selector: 'eo-ng-table-cluster-status',
  template: `
    <ng-container [ngSwitch]="status">
      <span *ngSwitchCase="'ABNORMAL'" class="red-bold">异常</span>
      <span *ngSwitchCase="'NORMAL'" class="green-bold">正常</span>
      <span *ngSwitchCase="'PARTIALLY_NORMAL'" class="orange-bold">部分正常</span>
    </ng-container>
    `,
  styles: [
  ]
})
export class TableClusterStatusTdComponent {
      @Input() status:string | undefined
}

@Component({
  selector: 'eo-ng-table-publish-status',
  template: `
    <ng-container [ngSwitch]="publish">
      <span *ngSwitchCase="'UNPUBLISHED'" class="red-bold">未发布</span>
      <span *ngSwitchCase="'PUBLISHED'" class="green-bold">已发布</span>
      <span *ngSwitchCase="'DEFECT'" class="font-bold">缺失</span>
    </ng-container>
    `,
  styles: [
  ]
})
export class TablePublishStatusTdComponent {
      @Input() publish:string | undefined
}

@Component({
  selector: 'eo-ng-table-publish-change-status',
  template: `
    <ng-container [ngSwitch]="status">
      <span class="bg-[#37a9fd1a] text-[#37a9fd]" *ngSwitchCase="'NEW'">新</span>
      <span class="bg-[#ff66001a] text-[#ff6600]" *ngSwitchCase="'MODIFY'">改</span>
      <span class="bg-[#7a7a7a1a] text-[#7a7a7a]" *ngSwitchCase="'DELETE'">删</span>
    </ng-container>
    `,
  styles: [
    `
    span{
        width: 16px;
        height: 16px;
        font-size: 12px;
        border-radius: 4px;
        display: flex;
        justify-content: center;
        flex-direction: column;
        align-items: flex-start;
        padding: 0px 2px;
        margin-left: 6px;
        margin-right: 6px;
    }`
  ]
})
export class TablePublishChangeStatusTdComponent {
      @Input() status:string | undefined
}

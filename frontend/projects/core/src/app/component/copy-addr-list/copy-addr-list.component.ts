import { Component, Input, OnInit, Output, EventEmitter } from '@angular/core'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'

// 网关集群-网关节点中,表格一行中有多个地址需要支持单行地址复制
@Component({
  selector: 'eo-ng-copy-addr-list',
  template: `
  <span
  *ngIf="!addrItem.expand || !addrItem[keyName].includes(','); else expandTpl"
  class="cluster-node-addr-btn overflow-ellipsis w-[100%] inline-block overflow-hidden align-middle"
  eoNgFeedbackTooltip
  [nzTooltipTitle]="tooltipTpl"
  [nzTooltipTitleContext]="{ addrLong: addrItem[keyName] }"
  [nzTooltipVisible]="false"
  [nzTooltipTrigger]="'hover'"
>
  <span>
    <span
      class="overflow-ellipsis inline-block overflow-hidden align-middle"
      [ngClass]="{ 'w-[86%]': addrItem[keyName].includes(',') }"
      >{{ addrItem[keyName] }}
    </span>
    <button
      *ngIf="!addrItem[keyName].includes(',')"
      eo-copy
      eo-ng-button
      nzType="primary"
      nzGhost
      class="deploy-node-copy-btn ant-btn-text border-transparent h-[22px]"
      [copyText]="addrItem[keyName]"
      (copyCallback)="copyCallback()"
    >
      <svg class="iconpark-icon"><use href="#copy"></use></svg>
    </button>
  </span>
  <svg
    *ngIf="addrItem[keyName].includes(',')"
    class="iconpark-icon ml-[8px]"
    (click)="addrItem.expand = true ;addrItemChange.emit(_addrItem)"
  >
    <use href="#down"></use>
  </svg>
</span>
<ng-template #expandTpl>
  <div class="flex flex-nowrap items-center justify-between">
    <div>
      <ng-container *ngFor="let addr of addrItem[keyName].split(',')">
        <div class="block w-[100%]">
          <span class="leading-[22px]">{{ addr }}</span
          ><button
            eo-copy
            eo-ng-button
            nzType="primary"
            nzGhost
            class="ml-[8px] border-transparent h-[22px]"
            [copyText]="addr"
            (copyCallback)="copyCallback()"
          >
            <svg class="iconpark-icon"><use href="#copy"></use></svg>
          </button>
        </div>
      </ng-container>
    </div>
    <svg class="iconpark-icon" (click)="addrItem.expand = false;addrItemChange.emit(_addrItem)">
      <use href="#up"></use>
    </svg>
  </div>
</ng-template>


<ng-template #tooltipTpl let-addrLong="addrLong">
  <ng-container *ngFor="let addr of addrLong.split(',')">
    <div class="flex justify-between">
      <span class="leading-[22px]">{{ addr }}</span>
    </div>
  </ng-container>
</ng-template>

  `,
  styles: [
  ]
})
export class CopyAddrListComponent implements OnInit {
  @Input()
  get addrItem () {
    return this._addrItem
  }

  set addrItem (val:{expand:boolean, [key:string]:any}) {
    this._addrItem = val
    this.addrItemChange.emit(val)
  }

   @Output() addrItemChange:EventEmitter<{expand:boolean, [key:string]:any}> = new EventEmitter()
  _addrItem:{expand:boolean, [key:string]:any} = { expand: false }
  @Input() keyName:string = ''
  constructor (private message: EoNgFeedbackMessageService) { }

  ngOnInit (): void {
  }

  copyCallback () {
    this.message.success('复制成功', {
      nzDuration: 1000
    })
  }
}

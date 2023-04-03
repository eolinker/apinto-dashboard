import { Component, EventEmitter, Input, Output } from '@angular/core'

@Component({
  selector: 'eo-ng-search-input-group',
  template: `
      <eo-ng-input-group
      [ngClass]="{'rounded-SEARCH_RADIUS':eoSingle, 'w-SEARCH':eoSingle && !eoNoWidth,'w-INPUT_NORMAL':!eoSingle&& !eoNoWidth}"
      [nzPrefix]="prefixTpl"
      [nzSuffix]="inputClearTpl"
    >
      <ng-content></ng-content>
    </eo-ng-input-group>
    <ng-template #inputClearTpl>
      <span
        class="ant-input-clear-icon"
        *ngIf="eoInputVal"
        (click)="eoClick.emit($event)"
      ><svg class="iconpark-icon"><use href="#close-small"></use></svg></span>
    </ng-template>
    <ng-template #prefixTpl
      ><svg class="iconpark-icon"><use href="#search"></use></svg>
    </ng-template>
  `,
  styles: [
  ]
})
export class EoNgSearchInputGroupComponent {
  @Input() eoSingle:boolean = true
  @Input() eoNoWidth:boolean = false
  @Input() eoInputVal:string|number|undefined // input组件绑定的value
  @Output() eoClick:EventEmitter<any> = new EventEmitter()
}

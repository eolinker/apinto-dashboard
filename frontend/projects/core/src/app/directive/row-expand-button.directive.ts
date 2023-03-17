/* eslint-disable quotes */
import { Directive, EventEmitter, Input, Output } from '@angular/core'

@Directive({
  selector: '[eo-ng-row-expand-button]',
  host: {
    class: 'ant-table-row-expand-icon',
    '[type]': `'button'`,
    '[class.ant-table-row-expand-icon-expanded]': `!spaceMode && expand`,
    '[class.ant-table-row-expand-icon-collapsed]': `!spaceMode && !expand`,
    '[class.ant-table-row-expand-icon-spaced]': 'spaceMode',
    '(click)': 'onHostClick()'
  }
})
export class RowExpandButtonDirective {
  @Input() expand = false;
  @Input() spaceMode = false;
  @Output() readonly expandChange = new EventEmitter();

  onHostClick (): void {
    if (!this.spaceMode) {
      this.expand = !this.expand
      this.expandChange.emit(this.expand)
    }
  }
}

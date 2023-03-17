import { ChangeDetectionStrategy, Component, Input, TemplateRef, ViewEncapsulation } from '@angular/core'
import { NzTransferListComponent } from 'ng-zorro-antd/transfer'

@Component({
  selector: 'eo-ng-transfer-list',
  templateUrl: './transfer-list.component.html',
  encapsulation: ViewEncapsulation.None,
  preserveWhitespaces: false,
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class TransferListComponent extends NzTransferListComponent {
  @Input() renderTop: TemplateRef<void> | null = null;
  @Input() renderTopParams: any = {}
  @Input() renderTitle:string|null = ''

  ngDoCheck () {
    // this.stat.shownCount = this.validData.length
  }
}

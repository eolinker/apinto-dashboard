import { ChangeDetectionStrategy, Component, Input, QueryList, TemplateRef, ViewChildren, ViewEncapsulation } from '@angular/core'
import { NzTransferComponent } from 'ng-zorro-antd/transfer'
import { TransferListComponent } from '../transfer-list/transfer-list.component'

@Component({
  selector: 'eo-ng-transfer',
  templateUrl: './transfer.component.html',
  encapsulation: ViewEncapsulation.None,
  changeDetection: ChangeDetectionStrategy.OnPush,
  preserveWhitespaces: false
})
export class TransferComponent extends NzTransferComponent {
  @ViewChildren(TransferListComponent) override lists!: QueryList<TransferListComponent>;
  @Input() nzRenderTopList: Array<TemplateRef<any> | null> | null = null;
  @Input() nzRenderTopParamsList: Array<any> | null = null;
  @Input() nzRenderTitleList: Array<string> | null = null;
  @Input() nzNoButton:boolean = false
}

import { NgModule } from '@angular/core'
import { NzTransferModule } from 'ng-zorro-antd/transfer'
import { TransferComponent } from './transfer/transfer.component'
import { TransferListComponent } from './transfer-list/transfer-list.component'
import { TransferSearchComponent } from './transfer-search/transfer-search.component'

import { BidiModule } from '@angular/cdk/bidi'
import { CommonModule } from '@angular/common'
import { FormsModule } from '@angular/forms'

import { NzButtonModule } from 'ng-zorro-antd/button'
import { NzCheckboxModule } from 'ng-zorro-antd/checkbox'
import { NzEmptyModule } from 'ng-zorro-antd/empty'
import { NzI18nModule } from 'ng-zorro-antd/i18n'
import { NzInputModule } from 'ng-zorro-antd/input'
import { ComponentModule } from '../component.module'

@NgModule({
  declarations: [
    TransferComponent,
    TransferListComponent,
    TransferSearchComponent
  ],
  imports: [
    NzTransferModule,
    BidiModule,
    CommonModule,
    FormsModule,
    NzCheckboxModule,
    NzButtonModule,
    NzInputModule,
    NzI18nModule,
    NzEmptyModule,
    ComponentModule
  ],
  exports: [
    TransferComponent
  ]
})
export class EoNgTransferModule { }

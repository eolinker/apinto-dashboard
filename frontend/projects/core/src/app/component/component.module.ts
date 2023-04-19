/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2022-08-09 23:06:58
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2022-08-23 16:45:46
 * @FilePath: /apinto/src/app/component/component.module.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { FormsModule, ReactiveFormsModule } from '@angular/forms'
import { EoNgTableModule } from 'eo-ng-table'
import { EditableEnvTableComponent } from './editable-env-table/editable-env-table.component'
import { DynamicComponentComponent } from './dynamic-component/dynamic-component.component'
import { NzSpaceModule } from 'ng-zorro-antd/space'
import { EoNgDatePickerModule } from 'eo-ng-date-picker'
import { EoNgSelectModule } from 'eo-ng-select'
import { EoNgCheckboxModule } from 'eo-ng-checkbox'
import { NzPopconfirmModule } from 'ng-zorro-antd/popconfirm'
import { DragDropModule } from '@angular/cdk/drag-drop'
import { EoNgInputModule } from 'eo-ng-input'
import { NzFormModule } from 'ng-zorro-antd/form'
import { EoNgButtonModule } from 'eo-ng-button'
import { EoNgFeedbackTooltipModule } from 'eo-ng-feedback'
import { EoNgSwitchModule } from 'eo-ng-switch'
import { NzTableModule } from 'ng-zorro-antd/table'
import { EoNgDropdownModule } from 'eo-ng-dropdown'
import { NzCheckboxModule } from 'ng-zorro-antd/checkbox'
import { NzOutletModule } from 'ng-zorro-antd/core/outlet'
import { NzDividerModule } from 'ng-zorro-antd/divider'
import { ScrollingModule } from '@angular/cdk/scrolling'
import { NzResizableModule } from 'ng-zorro-antd/resizable'
import { DirectiveModule } from '../directive/directive.module'
import { EoNgApintoTableModule } from 'projects/eo-ng-apinto-table/src/public-api'
import { TableDisabledStatusTdComponent, TableStatusTdComponent } from './table-status-td/table-status-td.component'
import { EoNgSearchInputGroupComponent } from './eo-ng-search-input-group/eo-ng-search-input-group.component'
import { ApiMethodTagComponent } from './api-method-tag/api-method-tag.component'
import { CopyAddrListComponent } from './copy-addr-list/copy-addr-list.component'
import { EoNgCopyModule } from 'eo-ng-copy'
import { CardListComponent } from './card-list/card-list.component'
import { NzCardModule } from 'ng-zorro-antd/card'
import { NzListModule } from 'ng-zorro-antd/list'
@NgModule({
  declarations: [
    EditableEnvTableComponent,
    DynamicComponentComponent,
    TableStatusTdComponent,
    TableDisabledStatusTdComponent,
    EoNgSearchInputGroupComponent,
    ApiMethodTagComponent,
    CopyAddrListComponent,
    CardListComponent
  ],
  imports: [
    CommonModule,
    FormsModule,
    EoNgTableModule,
    NzSpaceModule,
    EoNgDatePickerModule,
    EoNgSelectModule,
    EoNgCheckboxModule,
    NzPopconfirmModule,
    ReactiveFormsModule,
    DragDropModule,
    EoNgInputModule,
    NzFormModule,
    EoNgButtonModule,
    EoNgFeedbackTooltipModule,
    EoNgSwitchModule,
    NzTableModule,
    ScrollingModule,
    NzCheckboxModule,
    NzResizableModule,
    NzOutletModule,
    EoNgDropdownModule,
    EoNgSelectModule,
    NzPopconfirmModule,
    NzDividerModule,
    DirectiveModule,
    EoNgApintoTableModule,
    EoNgCopyModule,
    NzCardModule,
    NzListModule
  ],
  exports: [
    EditableEnvTableComponent,
    DynamicComponentComponent,
    TableStatusTdComponent,
    TableDisabledStatusTdComponent,
    EoNgSearchInputGroupComponent,
    ApiMethodTagComponent,
    CopyAddrListComponent,
    CardListComponent
  ]
})
export class ComponentModule { }

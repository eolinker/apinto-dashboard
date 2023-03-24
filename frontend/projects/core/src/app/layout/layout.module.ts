/*
 * @Author: maggieyyy im.ymj@hotmail.com
 * @Date: 2022-07-11 23:20:14
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2022-09-20 22:01:32
 * @FilePath: /apinto/src/app/app.module.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { CommonModule } from '@angular/common'
import { NgModule } from '@angular/core'
import { FormsModule, ReactiveFormsModule } from '@angular/forms'
import { NoopAnimationsModule } from '@angular/platform-browser/animations'
import { AppRoutingModule } from '../app-routing.module'
import { ComponentModule } from '../component/component.module'
import { BasicLayoutComponent } from './basic-layout/basic-layout.component'
import { NzTableModule } from 'ng-zorro-antd/table'
import { NzAvatarModule } from 'ng-zorro-antd/avatar'
import { EoNgBreadcrumbModule } from 'eo-ng-breadcrumb'
import { EoNgLayoutModule } from 'eo-ng-layout'
import { EoNgMenuModule } from 'eo-ng-menu'
import { EoNgSelectModule } from 'eo-ng-select'
import { EoNgFeedbackModalModule, EoNgFeedbackTooltipModule } from 'eo-ng-feedback'
import { EoNgTabsModule } from 'eo-ng-tabs'
import { EoNgCheckboxModule } from 'eo-ng-checkbox'
import { EoNgInputModule } from 'eo-ng-input'
import { EoNgCascaderModule } from 'eo-ng-cascader'
import { EoNgDropdownModule } from 'eo-ng-dropdown'
import { EoNgTreeModule } from 'eo-ng-tree'
import { EoNgTableModule } from 'eo-ng-table'
import { EoNgDatePickerModule } from 'eo-ng-date-picker'
import { EoNgSwitchModule } from 'eo-ng-switch'
import { EoNgButtonModule } from 'eo-ng-button'
import { EoNgTransferModule } from '../component/transfer/transfer.module'
import { EoNgCopyModule } from 'eo-ng-copy'
import { EoNgUploadModule } from 'eo-ng-upload'
import { LoginComponent } from './login/login.component'
import { NzFormModule } from 'ng-zorro-antd/form'
import { NzLayoutModule } from 'ng-zorro-antd/layout'
import { NzToolTipModule } from 'ng-zorro-antd/tooltip'
import { NzCheckboxModule } from 'ng-zorro-antd/checkbox'
import { NzGridModule } from 'ng-zorro-antd/grid'
import { NzModalModule } from 'ng-zorro-antd/modal'
import { NzHighlightModule } from 'ng-zorro-antd/core/highlight'
import { DirectiveModule } from '../directive/directive.module'
import { NzUploadModule } from 'ng-zorro-antd/upload'
import { EoNgEmptyModule } from 'eo-ng-empty'

@NgModule({
  declarations: [
    BasicLayoutComponent,
    LoginComponent
  ],
  imports: [
    AppRoutingModule,
    EoNgBreadcrumbModule,
    EoNgLayoutModule,
    EoNgMenuModule,
    CommonModule,
    EoNgSelectModule,
    NoopAnimationsModule,
    FormsModule,
    ComponentModule,
    EoNgTableModule,
    EoNgFeedbackModalModule,
    EoNgFeedbackTooltipModule,
    EoNgTabsModule,
    EoNgCheckboxModule,
    NzTableModule,
    EoNgInputModule,
    EoNgCascaderModule,
    EoNgTreeModule,
    EoNgDropdownModule,
    EoNgDatePickerModule,
    EoNgSwitchModule,
    NzAvatarModule,
    EoNgButtonModule,
    EoNgTransferModule,
    NzFormModule,
    ReactiveFormsModule,
    NzLayoutModule,
    NzToolTipModule,
    NzCheckboxModule,
    NzGridModule,
    NzModalModule,
    NzHighlightModule,
    EoNgCopyModule,
    DirectiveModule,
    EoNgUploadModule,
    NzUploadModule,
    EoNgEmptyModule
  ],
  exports: [
    BasicLayoutComponent
  ],
  providers: []
})
export class LayoutModule { }

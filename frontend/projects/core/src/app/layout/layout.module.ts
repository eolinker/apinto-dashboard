/*
 * @Author: maggieyyy im.ymj@hotmail.com
 * @Date: 2022-07-11 23:20:14
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-07-02 10:24:35
 * @FilePath: \apinto\projects\core\src\app\layout\layout.module.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { CommonModule } from '@angular/common'
import { NgModule } from '@angular/core'
import { FormsModule, ReactiveFormsModule } from '@angular/forms'
import { AppRoutingModule } from '../app-routing.module'
import { ComponentModule } from '../component/component.module'
import { BasicLayoutComponent } from './basic-layout/basic-layout.component'
import { NzAvatarModule } from 'ng-zorro-antd/avatar'
import { EoNgBreadcrumbModule } from 'eo-ng-breadcrumb'
import { EoNgLayoutModule } from 'eo-ng-layout'
import { EoNgSelectModule } from 'eo-ng-select'
import { EoNgFeedbackAlertModule, EoNgFeedbackModalModule, EoNgFeedbackTooltipModule } from 'eo-ng-feedback'
import { EoNgTabsModule } from 'eo-ng-tabs'
import { EoNgCheckboxModule } from 'eo-ng-checkbox'
import { EoNgInputModule } from 'eo-ng-input'
import { EoNgCascaderModule } from 'eo-ng-cascader'
import { EoNgDropdownModule } from 'eo-ng-dropdown'
import { EoNgTreeModule } from 'eo-ng-tree'
import { EoNgDatePickerModule } from 'eo-ng-date-picker'
import { EoNgSwitchModule } from 'eo-ng-switch'
import { EoNgButtonModule } from 'eo-ng-button'
import { EoNgCopyModule } from 'eo-ng-copy'
import { EoNgUploadModule } from 'eo-ng-upload'
// import { LoginComponent } from './login/login.component'
import { LogoComponent } from './login/logo/logo.component'
// import { PasswordComponent } from './login/password/password.component'
import { NzFormModule } from 'ng-zorro-antd/form'
import { NzGridModule } from 'ng-zorro-antd/grid'
import { NzHighlightModule } from 'ng-zorro-antd/core/highlight'
import { DirectiveModule } from '../directive/directive.module'
import { NzUploadModule } from 'ng-zorro-antd/upload'
import { EoNgEmptyModule } from 'eo-ng-empty'
import { EoNgMenuModule } from 'eo-ng-menu'
import { LocalPluginComponent } from './local-plugin/local-plugin.component'
import { GuideComponent } from './guide/guide.component'
import { NzDividerModule } from 'ng-zorro-antd/divider'
import { DynamicDemoComponent } from './dynamic-demo/dynamic-demo.component'
import { SafePipe } from '../pipe/safe.pipe'
import { NzSpinModule } from 'ng-zorro-antd/spin'
import { NotFoundPageComponent } from './not-found-page/not-found-page.component'
import { IntelligentPluginComponent } from './intelligent-plugin/intelligent-plugin.component'
import { RemotePluginComponent } from './remote-plugin/remote-plugin.component'
import { EoNgApintoTableModule } from 'projects/eo-ng-apinto-table/src/public-api'
import { PluginWrapperComponent } from './plugin-wrapper/plugin-wrapper.component'
import { SystemEmailConfigComponent } from './system/email/config/config.component'
import { SystemWebhookConfigComponent } from './system/webhook/config/config.component'
import { SystemWebhookListComponent } from './system/webhook/list/list.component'
import { RedirectPage } from './redirect-page/redirect-page.component'
import { RouterLayoutComponent } from './router-layout/router-layout.component'

const sharedEoLibraryModules = [
  EoNgButtonModule,
  EoNgBreadcrumbModule,
  EoNgLayoutModule,
  EoNgMenuModule,
  EoNgSelectModule,
  EoNgFeedbackAlertModule,
  EoNgFeedbackModalModule,
  EoNgFeedbackTooltipModule,
  EoNgTabsModule,
  EoNgCheckboxModule,
  EoNgInputModule,
  EoNgCascaderModule,
  EoNgTreeModule,
  EoNgDropdownModule,
  EoNgDatePickerModule,
  EoNgSwitchModule,
  EoNgCopyModule,
  EoNgApintoTableModule,
  EoNgUploadModule,
  EoNgEmptyModule
]
@NgModule({
  declarations: [
    BasicLayoutComponent,
    RedirectPage,
    LogoComponent,
    LocalPluginComponent,
    GuideComponent,
    DynamicDemoComponent,
    SafePipe,
    NotFoundPageComponent,
    IntelligentPluginComponent,
    RemotePluginComponent,
    PluginWrapperComponent,
    SystemWebhookListComponent,
    SystemEmailConfigComponent,
    SystemWebhookConfigComponent,
    RouterLayoutComponent
  ],
  imports: [
    AppRoutingModule,
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    ComponentModule,
    NzAvatarModule,
    NzFormModule,
    NzGridModule,
    NzHighlightModule,
    DirectiveModule,
    NzUploadModule,
    NzDividerModule,
    NzSpinModule,
    ...sharedEoLibraryModules
  ],
  exports: [
    BasicLayoutComponent,
    RedirectPage,
    SystemEmailConfigComponent,
    RouterLayoutComponent
  ],
  providers: []
})
export class LayoutModule { }

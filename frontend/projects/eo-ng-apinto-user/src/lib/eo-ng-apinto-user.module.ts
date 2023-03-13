import { Inject, NgModule } from '@angular/core'
import { ResetPswComponent } from './component/reset-psw/reset-psw.component'
import { UserAvatarComponent } from './component/user-avatar/user-avatar.component'
import { UserProfileComponent } from './component/user-profile/user-profile.component'
import { NzAvatarModule } from 'ng-zorro-antd/avatar'
import { CommonModule } from '@angular/common'
import { EoNgDropdownModule } from 'eo-ng-dropdown'
import { EoNgTreeModule } from 'eo-ng-tree'
import { EoNgSwitchModule } from 'eo-ng-switch'
import { RolesGroupComponent } from './layout/roles/group/group.component'
import { RouterModule } from '@angular/router'
import { API_SERVICE_ADAPTER, ApiServiceAdapter } from './constant/api-service-adapter'
import { EoNgInputModule } from 'eo-ng-input'
import { EoNgTableModule } from 'eo-ng-table'
import { EoNgButtonModule } from 'eo-ng-button'
import { EoNgFeedbackDrawerModule, EoNgFeedbackMessageModule, EoNgFeedbackTooltipModule } from 'eo-ng-feedback'
import { RolesListComponent } from './layout/roles/list/list.component'
import { FormsModule, ReactiveFormsModule } from '@angular/forms'
import { NzFormModule } from 'ng-zorro-antd/form'
import { EoNgSelectModule } from 'eo-ng-select'
import { RoleProfileComponent } from './component/role-profile/role-profile.component'
import { RoleAccessComponent } from './component/role-access/role-access.component'
import { AppConfigAdapter, APP_CONFIG_ADAPTER } from './constant/app-config-adapter'
import { APP_SERVICE_ADAPTER, AppServiceAdapter } from './constant/app-service-adapter'
import { EoNgCheckboxModule } from 'eo-ng-checkbox'
import { UserAccessDirective } from './directive/user-access.directive'
import { TreeDragDirective } from './directive/tree-drag.directive'
import { EoNgAutoCompleteModule } from 'eo-ng-auto-complete'
import { NzTableModule } from 'ng-zorro-antd/table'
import { NzCheckboxModule } from 'ng-zorro-antd/checkbox'
import { NzOutletModule } from 'ng-zorro-antd/core/outlet'
import { NzDividerModule } from 'ng-zorro-antd/divider'
import { NzPopconfirmModule } from 'ng-zorro-antd/popconfirm'
import { NzResizableModule } from 'ng-zorro-antd/resizable'
import { DragDropModule } from '@angular/cdk/drag-drop'
import { ScrollingModule } from '@angular/cdk/scrolling'
import { EoNgApintoTableModule } from '../../../eo-ng-apinto-table/src/public-api'
import { EoNgSearchInputGroupComponent } from './component/eo-ng-search-input-group/eo-ng-search-input-group.component'
import { BaseInfoServiceAdapter, BASEINFO_SERVICE_ADAPTER } from './constant/base-info-service-adapter';
import { RefreshPswComponent } from './layout/roles/refresh-psw/refresh-psw.component';
import { ChangeUserRoleComponent } from './layout/roles/change-user-role/change-user-role.component'

@NgModule({
  declarations: [
    UserAvatarComponent,
    UserProfileComponent,
    ResetPswComponent,
    RolesGroupComponent,
    RolesListComponent,
    UserProfileComponent,
    RoleProfileComponent,
    RoleAccessComponent,
    UserAccessDirective,
    TreeDragDirective,
    EoNgSearchInputGroupComponent,
    RefreshPswComponent,
    ChangeUserRoleComponent
  ],
  imports: [
    NzAvatarModule,
    CommonModule,
    EoNgDropdownModule,
    RouterModule,
    EoNgTreeModule,
    EoNgInputModule,
    EoNgTableModule,
    EoNgButtonModule,
    EoNgFeedbackDrawerModule,
    EoNgFeedbackMessageModule,
    FormsModule,
    NzFormModule,
    ReactiveFormsModule,
    EoNgSelectModule,
    EoNgCheckboxModule,
    EoNgSwitchModule,
    EoNgAutoCompleteModule,
    ScrollingModule,
    NzCheckboxModule,
    NzResizableModule,
    NzTableModule,
    NzOutletModule,
    DragDropModule,
    EoNgDropdownModule,
    NzPopconfirmModule,
    NzDividerModule,
    EoNgFeedbackTooltipModule,
    EoNgApintoTableModule
  ],
  exports: [
    UserAvatarComponent,
    RolesGroupComponent
  ]
})
export class EoNgApintoUserModule {
  constructor (
    @Inject(API_SERVICE_ADAPTER) private apiService: ApiServiceAdapter,
    @Inject(APP_CONFIG_ADAPTER) private appConfig:AppConfigAdapter,
    @Inject(APP_SERVICE_ADAPTER) private appService:AppServiceAdapter,
    @Inject(BASEINFO_SERVICE_ADAPTER) private baseInfo:BaseInfoServiceAdapter
  ) {
    if (!apiService) {
      throw new Error('You must provate a api adapter')
    }
    if (!appConfig) {
      throw new Error('You must provate a appConfig adapter')
    }
  }
}

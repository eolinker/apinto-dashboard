import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'

import { MonitorAlarmRoutingModule } from './monitor-alarm-routing.module'
import { MonitorAlarmComponent } from './monitor-alarm.component'
import { MonitorAlarmAreaComponent } from './area/area.component'
import { MonitorAlarmAreaApiComponent } from './area/others/api/api.component'
import { MonitorAlarmAreaAppComponent } from './area/others/app/app.component'
import { MonitorAlarmAreaOthersDetailComponent } from './area/others/detail/detail.component'
import { MonitorAlarmAreaUpstreamComponent } from './area/others/service/service.component'
import { MonitorAlarmAreaTotalComponent } from './area/total/total.component'
import { MonitorAlarmConfigComponent } from './config/config.component'
import { FormsModule, ReactiveFormsModule } from '@angular/forms'
import { EoNgFeedbackAlertModule, EoNgFeedbackModalModule, EoNgFeedbackTooltipModule } from 'eo-ng-feedback'
import { EoNgInputModule } from 'eo-ng-input'
import { EoNgButtonModule } from 'eo-ng-button'
import { EoNgCheckboxModule } from 'eo-ng-checkbox'
import { EoNgSelectModule } from 'eo-ng-select'
import { EoNgTabsModule } from 'eo-ng-tabs'
import { NzFormModule } from 'ng-zorro-antd/form'
import { EoNgApintoTableModule } from 'projects/eo-ng-apinto-table/src/public-api'
import { ComponentModule } from '../../component/component.module'
import { EoNgRadioModule } from 'eo-ng-radio'
import { EoNgDatePickerModule } from 'eo-ng-date-picker'
import { EoNgDropdownModule } from 'eo-ng-dropdown'
import { EoNgEmptyModule } from 'eo-ng-empty'
import { MonitorAlarmMessageComponent } from './message/message.component'
import { MonitorAlarmAreaOthersApiTotalComponent } from './area/others/api/total/total.component'
import { MonitorAlarmAreaOthersAppTotalComponent } from './area/others/app/total/total.component'
import { MonitorAlarmAreaOthersUpstreamTotalComponent } from './area/others/service/total/total.component'
import { RouterModule } from '@angular/router'
import { NzSpinModule } from 'ng-zorro-antd/spin'
import { DirectiveModule } from '../../directive/directive.module'
import { MonitorAlarmStrategyComponent } from './area/alarm/strategy/strategy.component'
import { MonitorAlarmHistoryComponent } from './area/alarm/history/history.component'
import { MonitorAlarmStrategyListComponent } from './area/alarm/strategy/list/list.component'
import { MonitorAlarmStrategyConfigComponent } from './area/alarm/strategy/config/config.component'
import { MonitorAlarmStrategyMessageComponent } from './area/alarm/strategy/message/message.component'
import { EoNgSwitchModule } from 'eo-ng-switch'
import { MonitorAlarmStrategyTransferComponent } from './area/alarm/strategy/transfer/transfer.component'
import { EoNgTransferModule } from '../../component/transfer/transfer.module'
import { EoNgTableModule } from 'eo-ng-table'
import { MonitorAlarmStrategyRuleComponent } from './area/alarm/strategy/rule/rule.component'
import { NzDividerModule } from 'ng-zorro-antd/divider'
import { NzTableModule } from 'ng-zorro-antd/table'
import { EoNgCascaderModule } from 'eo-ng-cascader'
import { MonitorAlarmStrategyAlertComponent } from './area/alarm/strategy/alert/alert.component'
import { NgxEchartsModule } from 'ngx-echarts'
import { MonitorPieGraphComponent } from './area/graph/pie/pie.component'
import { MonitorLineGraphComponent } from './area/graph/line/line.component'
import { MonitorAlarmTableComponent } from './table/table.component'

@NgModule({
  declarations: [
    MonitorAlarmComponent,
    MonitorAlarmConfigComponent,
    MonitorAlarmAreaComponent,
    MonitorAlarmAreaTotalComponent,
    MonitorAlarmAreaApiComponent,
    MonitorAlarmAreaAppComponent,
    MonitorAlarmAreaUpstreamComponent,
    MonitorAlarmAreaOthersApiTotalComponent,
    MonitorAlarmAreaOthersDetailComponent,
    MonitorAlarmMessageComponent,
    MonitorAlarmAreaOthersAppTotalComponent,
    MonitorAlarmAreaOthersUpstreamTotalComponent,
    MonitorAlarmStrategyComponent,
    MonitorAlarmHistoryComponent,
    MonitorAlarmStrategyListComponent,
    MonitorAlarmStrategyConfigComponent,
    MonitorAlarmStrategyMessageComponent,
    MonitorAlarmStrategyTransferComponent,
    MonitorAlarmStrategyRuleComponent,
    MonitorAlarmStrategyAlertComponent,
    MonitorPieGraphComponent,
    MonitorLineGraphComponent,
    MonitorAlarmTableComponent
  ],
  imports: [
    CommonModule,
    MonitorAlarmRoutingModule,
    ReactiveFormsModule,
    FormsModule,
    EoNgCheckboxModule,
    EoNgTabsModule,
    EoNgFeedbackModalModule,
    EoNgSelectModule,
    EoNgApintoTableModule,
    EoNgFeedbackAlertModule,
    EoNgRadioModule,
    EoNgDatePickerModule,
    EoNgFeedbackTooltipModule,
    NzFormModule,
    EoNgInputModule,
    EoNgButtonModule,
    ComponentModule,
    EoNgDropdownModule,
    EoNgEmptyModule,
    RouterModule,
    NzSpinModule,
    DirectiveModule,
    EoNgSwitchModule,
    EoNgTransferModule,
    EoNgTableModule,
    NzDividerModule,
    NzTableModule,
    EoNgCascaderModule,
    NgxEchartsModule
  ]
})
export class MonitorAlarmModule { }

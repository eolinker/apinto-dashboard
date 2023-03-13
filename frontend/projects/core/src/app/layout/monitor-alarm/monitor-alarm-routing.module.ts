import { NgModule } from '@angular/core'
import { RouterModule, Routes } from '@angular/router'
import { MonitorAlarmAreaApiComponent } from './area/others/api/api.component'
import { MonitorAlarmAreaComponent } from './area/area.component'
import { MonitorAlarmAreaTotalComponent } from './area/total/total.component'
import { MonitorAlarmConfigComponent } from './config/config.component'
import { MonitorAlarmComponent } from './monitor-alarm.component'
import { MonitorAlarmAreaOthersDetailComponent } from './area/others/detail/detail.component'
import { MonitorAlarmAreaAppComponent } from './area/others/app/app.component'
import { MonitorAlarmAreaUpstreamComponent } from './area/others/service/service.component'
import { MonitorAlarmAreaOthersApiTotalComponent } from './area/others/api/total/total.component'
import { MonitorAlarmMessageComponent } from './message/message.component'
import { MonitorAlarmAreaOthersAppTotalComponent } from './area/others/app/total/total.component'
import { MonitorAlarmAreaOthersUpstreamTotalComponent } from './area/others/service/total/total.component'
import { MonitorAlarmStrategyComponent } from './area/alarm/strategy/strategy.component'
import { MonitorAlarmHistoryComponent } from './area/alarm/history/history.component'
import { MonitorAlarmStrategyListComponent } from './area/alarm/strategy/list/list.component'
import { MonitorAlarmStrategyConfigComponent } from './area/alarm/strategy/config/config.component'
import { MonitorAlarmStrategyMessageComponent } from './area/alarm/strategy/message/message.component'

const routes:Routes = [{
  path: '',
  component: MonitorAlarmComponent,
  data: { id: '9' },
  children: [
    {
      path: 'area',
      component: MonitorAlarmAreaComponent,
      children: [
        { path: 'config', component: MonitorAlarmConfigComponent },
        { path: 'message/:partitionId', component: MonitorAlarmMessageComponent },
        {
          path: 'api/:partitionId',
          component: MonitorAlarmAreaApiComponent,
          children: [
            { path: '', component: MonitorAlarmAreaOthersApiTotalComponent },
            { path: 'detail/:monitorDataId', component: MonitorAlarmAreaOthersDetailComponent }
          ]
        },
        {
          path: 'app/:partitionId',
          component: MonitorAlarmAreaAppComponent,
          children: [
            { path: '', component: MonitorAlarmAreaOthersAppTotalComponent },
            { path: 'detail/:monitorDataId', component: MonitorAlarmAreaOthersDetailComponent }
          ]
        },
        {
          path: 'service/:partitionId',
          component: MonitorAlarmAreaUpstreamComponent,
          children: [
            { path: '', component: MonitorAlarmAreaOthersUpstreamTotalComponent },
            { path: 'detail/:monitorDataId', component: MonitorAlarmAreaOthersDetailComponent }
          ]
        },
        {
          path: 'strategy/:partitionId',
          component: MonitorAlarmStrategyComponent,
          children: [
            { path: '', component: MonitorAlarmStrategyListComponent },
            { path: 'config', component: MonitorAlarmStrategyConfigComponent },
            { path: 'message/:strategyUuid', component: MonitorAlarmStrategyMessageComponent }
          ]
        },
        { path: 'history/:partitionId', component: MonitorAlarmHistoryComponent },
        {
          path: 'total/:partitionId',
          component: MonitorAlarmAreaTotalComponent,
          pathMatch: 'full'
        }

      ]
    }]
}]

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class MonitorAlarmRoutingModule { }

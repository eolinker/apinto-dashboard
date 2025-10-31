import { NgModule } from '@angular/core'
import { RouterModule, Routes } from '@angular/router'
import { ListComponent } from '../list/list.component'
import { TrafficCreateComponent } from './create/create.component'
import { TrafficMessageComponent } from './message/message.component'
import { TrafficComponent } from './traffic.component'
import { GroupComponent } from '../group/group.component'

const routes: Routes = [
  {
    path: '',
    component: TrafficComponent,
    data: {
      id: '501'
    },
    children: [
      {
        path: 'group',
        component: GroupComponent,
        children: [
          {
            path: 'list/:clusterName',
            component: ListComponent
          }
        ]
      },
      {
        path: 'create/:clusterName',
        component: TrafficCreateComponent

      },
      { path: 'message/:clusterName/:strategyId', component: TrafficMessageComponent }
    ]
  }
]

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class TrafficStrategyRoutingModule { }

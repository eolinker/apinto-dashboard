import { NgModule } from '@angular/core'
import { RouterModule, Routes } from '@angular/router'
import { GroupComponent } from '../group/group.component'
import { ListComponent } from '../list/list.component'
import { VisitCreateComponent } from './create/create.component'
import { VisitMessageComponent } from './message/message.component'
import { VisitComponent } from './visit.component'

const routes: Routes = [
  {
    path: '',
    component: VisitComponent,
    data: {
      id: '503'
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
        component: VisitCreateComponent

      },
      { path: 'message/:clusterName/:strategyId', component: VisitMessageComponent }
    ]
  }
]

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class VisitStrategyRoutingModule { }

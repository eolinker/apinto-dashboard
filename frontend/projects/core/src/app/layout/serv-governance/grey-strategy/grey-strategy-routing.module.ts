import { NgModule } from '@angular/core'
import { RouterModule, Routes } from '@angular/router'
import { GroupComponent } from '../group/group.component'
import { ListComponent } from '../list/list.component'
import { GreyCreateComponent } from './create/create.component'
import { GreyComponent } from './grey.component'
import { GreyMessageComponent } from './message/message.component'

const routes: Routes = [{
  path: '',
  component: GreyComponent,
  data: {
    id: '505'
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
      component: GreyCreateComponent

    },
    { path: 'message/:clusterName/:strategyId', component: GreyMessageComponent }
  ]
}]

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class GreyStrategyRoutingModule { }

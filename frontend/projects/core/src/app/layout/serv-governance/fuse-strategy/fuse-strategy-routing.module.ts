import { NgModule } from '@angular/core'
import { RouterModule, Routes } from '@angular/router'
import { GroupComponent } from '../group/group.component'
import { ListComponent } from '../list/list.component'
import { FuseCreateComponent } from './create/create.component'
import { FuseComponent } from './fuse.component'
import { FuseMessageComponent } from './message/message.component'

const routes: Routes = [{
  path: '',
  component: FuseComponent,
  data: {
    id: '502'
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
      component: FuseCreateComponent

    },
    { path: 'message/:clusterName/:strategyId', component: FuseMessageComponent }
  ]
}]

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class FuseStrategyRoutingModule { }

import { NgModule } from '@angular/core'
import { RouterModule, Routes } from '@angular/router'
import { GroupComponent } from '../group/group.component'
import { ListComponent } from '../list/list.component'
import { CacheComponent } from './cache.component'
import { CacheCreateComponent } from './create/create.component'
import { CacheMessageComponent } from './message/message.component'

const routes: Routes = [{
  path: '',
  component: CacheComponent,
  data: {
    id: '504'
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
      component: CacheCreateComponent

    },
    { path: 'message/:clusterName/:strategyId', component: CacheMessageComponent }
  ]
}]

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class CacheStrategyRoutingModule { }

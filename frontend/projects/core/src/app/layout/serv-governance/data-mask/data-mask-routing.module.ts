import { NgModule } from '@angular/core'
import { RouterModule, Routes } from '@angular/router'
import { GroupComponent } from '../group/group.component'
import { ListComponent } from '../list/list.component'
import { DataMaskCreateComponent } from './create/create.component'
import { DataMaskMessageComponent } from './message/message.component'
import { DataMaskComponent } from './data-mask.component'

const routes: Routes = [{
  path: '',
  component: DataMaskComponent,
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
      component: DataMaskCreateComponent

    },
    {
      path: 'message/:clusterName/:strategyId',
      component: DataMaskMessageComponent
    }
  ]
}]

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class DataMaskRoutingModule { }

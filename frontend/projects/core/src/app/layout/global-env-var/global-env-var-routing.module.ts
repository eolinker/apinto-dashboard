import { NgModule } from '@angular/core'
import { RouterModule, Routes } from '@angular/router'
import { DeployEnvironmentCreateComponent } from './create/create.component'
import { DeployEnvironmentComponent } from './environment.component'
import { DeployEnvironmentListComponent } from './list/list.component'

const routes: Routes = [{
  path: '',
  component: DeployEnvironmentComponent,
  data: {
    id: '102',
    parentId: '1'
  },
  children: [
    {
      path: '',
      component: DeployEnvironmentListComponent
    },
    {
      path: 'create',
      component: DeployEnvironmentCreateComponent
    }
  ]
}
]

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class GlobalEnvVarRoutingModule { }

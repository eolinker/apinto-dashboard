import { NgModule } from '@angular/core'
import { RouterModule, Routes } from '@angular/router'
import { DeployClusterCertComponent } from './cluster/cert/cert.component'
import { DeployClusterComponent } from './cluster/cluster.component'
import { DeployClusterContentComponent } from './cluster/content/content.component'
import { DeployClusterCreateComponent } from './cluster/create/create.component'
import { DeployClusterConfComponent } from './cluster/conf/conf.component'
import { DeployClusterEnvironmentComponent } from './cluster/environment/environment.component'
import { DeployClusterListComponent } from './cluster/list/list.component'
import { DeployClusterNodesComponent } from './cluster/nodes/nodes.component'
import { DeployEnvironmentCreateComponent } from './environment/create/create.component'
import { DeployEnvironmentComponent } from './environment/environment.component'
import { DeployEnvironmentListComponent } from './environment/list/list.component'

const routes: Routes = [{
  path: '',
  component: DeployClusterComponent,
  data: {
    id: '1'
  },
  children: [
    {
      path: 'cluster',
      component: DeployClusterComponent,
      data: {
        id: '101',
        parentId: '1'
      },
      children: [
        {
          path: '',
          component: DeployClusterListComponent
        },
        {
          path: 'create',
          component: DeployClusterCreateComponent
        },
        {
          path: 'content/:clusterName',
          component: DeployClusterContentComponent,
          children: [
            {
              path: '',
              component: DeployClusterEnvironmentComponent
            },
            {
              path: 'cert',
              component: DeployClusterCertComponent
            },
            {
              path: 'nodes',
              component: DeployClusterNodesComponent
            },
            {
              path: 'conf',
              component: DeployClusterConfComponent
            }
          ]
        }]
    },
    {
      path: 'env',
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
}]

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class DeployRoutingModule { }

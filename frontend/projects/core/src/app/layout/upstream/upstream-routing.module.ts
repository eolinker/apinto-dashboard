import { NgModule } from '@angular/core'
import { RouterModule, Routes } from '@angular/router'
import { ServiceDiscoveryContentComponent } from './service-discovery/content/content.component'
import { ServiceDiscoveryCreateComponent } from './service-discovery/create/create.component'
import { ServiceDiscoveryListComponent } from './service-discovery/list/list.component'
import { ServiceDiscoveryMessageComponent } from './service-discovery/message/message.component'
import { ServiceDiscoveryPublishComponent } from './service-discovery/publish/publish.component'
import { ServiceDiscoveryComponent } from './service-discovery/service-discovery.component'
import { UpstreamContentComponent } from './upstream/content/content.component'
import { UpstreamCreateComponent } from './upstream/create/create.component'
import { UpstreamListComponent } from './upstream/list/list.component'
import { UpstreamMessageComponent } from './upstream/message/message.component'
import { UpstreamPublishComponent } from './upstream/publish/publish.component'
import { UpstreamComponent } from './upstream/upstream.component'

const routes: Routes = [{
  path: '',
  component: UpstreamComponent,
  data: {
    id: '2'
  },
  children: [
    {
      path: 'upstream',
      component: UpstreamComponent,
      data: {
        id: '201',
        parentId: '2'
      },
      children: [
        {
          path: '',
          component: UpstreamListComponent
        },
        {
          path: 'create',
          component: UpstreamCreateComponent
        },
        {
          path: 'content/:serviceName',
          component: UpstreamContentComponent,
          children: [
            {
              path: '',
              component: UpstreamPublishComponent
            },
            {
              path: 'publish',
              component: UpstreamPublishComponent
            },
            {
              path: 'message',
              component: UpstreamMessageComponent
            }
          ]
        }
      ]
    },
    {
      path: 'discovery',
      component: ServiceDiscoveryComponent,
      data: {
        id: '202',
        parentId: '2'
      },
      children: [
        {
          path: '',
          component: ServiceDiscoveryListComponent
        },
        {
          path: 'create',
          component: ServiceDiscoveryCreateComponent
        },
        {
          path: 'content/:discoveryName',
          component: ServiceDiscoveryContentComponent,
          children: [
            {
              path: '',
              component: ServiceDiscoveryPublishComponent
            },
            {
              path: 'message',
              component: ServiceDiscoveryMessageComponent
            }
          ]
        }
      ]
    }
  ]
}]

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class UpstreamRoutingModule { }

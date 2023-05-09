import { NgModule } from '@angular/core'
import { RouterModule, Routes } from '@angular/router'
import { TrafficCreateComponent } from './traffic/create/create.component'
import { TrafficComponent } from './traffic/traffic.component'
import { TrafficMessageComponent } from './traffic/message/message.component'
import { GroupComponent } from './group/group.component'
import { ServiceGovernanceComponent } from './service-governance.component'
import { ListComponent } from './list/list.component'
import { FuseCreateComponent } from './fuse/create/create.component'
import { FuseMessageComponent } from './fuse/message/message.component'
import { VisitComponent } from './visit/visit.component'
import { VisitCreateComponent } from './visit/create/create.component'
import { VisitMessageComponent } from './visit/message/message.component'
import { CacheComponent } from './cache/cache.component'
import { CacheCreateComponent } from './cache/create/create.component'
import { CacheMessageComponent } from './cache/message/message.component'
import { GreyComponent } from './grey/grey.component'
import { GreyCreateComponent } from './grey/create/create.component'
import { GreyMessageComponent } from './grey/message/message.component'
import { FuseComponent } from './fuse/fuse.component'

const routes: Routes = [{
  path: '',
  component: ServiceGovernanceComponent,
  data: {
    id: '5'
  },
  children: [
    {
      path: 'traffic',
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
    },
    {
      path: 'fuse',
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
    },
    {
      path: 'visit',
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
    },
    {
      path: 'cache',
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
    },
    {
      path: 'grey',
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
    }
  ]
}]

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class ServGovernanceRoutingModule { }

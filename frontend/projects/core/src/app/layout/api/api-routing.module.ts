import { NgModule } from '@angular/core'
import { RouterModule, Routes } from '@angular/router'
import { ApiListComponent } from './api-list/api-list.component'
import { ApiManagementGroupComponent } from './api-list/group/group.component'
import { ApiManagementListComponent } from './api-list/list/list.component'
import { ApiWebsocketCreateComponent } from './api-list/create/websocket-create/websocket-create.component'
import { ApiHttpCreateComponent } from './api-list/create/http-create/http-create.component'
import { ApiHttpMessageComponent } from './api-list/message/http-message/http-message.component'
import { ApiWebsocketMessageComponent } from './api-list/message/websocket-message/websocket-message.component'

export const routes: Routes = [{
  path: '',
  component: ApiListComponent,
  children: [
    {
      path: 'group',
      component: ApiManagementGroupComponent,
      children: [
        {
          path: 'list',
          component: ApiManagementListComponent,
          children: [{
            path: ':apiGroupId',
            component: ApiManagementListComponent
          }]
        }
      ]
    },
    {
      path: 'create',
      component: ApiHttpCreateComponent,
      children: [{
        path: ':apiGroupId',
        component: ApiHttpCreateComponent
      }]
    },
    {
      path: 'create-ws',
      component: ApiWebsocketCreateComponent,
      children: [{
        path: ':apiGroupId',
        component: ApiWebsocketCreateComponent
      }]
    },
    {
      path: 'message/:apiId',
      component: ApiHttpMessageComponent
    },
    {
      path: 'message-ws/:apiId',
      component: ApiWebsocketMessageComponent
    }
  ]
}
]

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class ApiRoutingModule { }

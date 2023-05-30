import { NgModule } from '@angular/core'
import { RouterModule, Routes } from '@angular/router'
import { ApiListComponent } from './api-list/api-list.component'
import { ApiManagementComponent } from './api-list/group/group.component'
import { ApiManagementListComponent } from './api-list/list/list.component'
import { ApiPluginTemplateContentComponent } from './plugin/content/content.component'
import { ApiPluginTemplateCreateComponent } from './plugin/create/create.component'
import { ApiPluginTemplateListComponent } from './plugin/list/list.component'
import { ApiPluginTemplateMessageComponent } from './plugin/message/message.component'
import { ApiPluginTemplateComponent } from './plugin/plugin.component'
import { ApiPluginTemplatePublishComponent } from './plugin/publish/publish.component'
import { RouterComponent } from './router/router.component'
import { ApiWebsocketCreateComponent } from './api-list/create/websocket-create/websocket-create.component'
import { ApiHttpCreateComponent } from './api-list/create/http-create/http-create.component'
import { ApiHttpMessageComponent } from './api-list/message/http-message/http-message.component'
import { ApiWebsocketMessageComponent } from './api-list/message/websocket-message/websocket-message.component'

const routes: Routes = [{
  path: '',
  component: RouterComponent,
  data: { id: '4' },
  children: [
    {
      path: 'api',
      component: ApiListComponent,
      children: [
        {
          path: 'group',
          component: ApiManagementComponent,
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
    },
    {
      path: 'plugin-template',
      component: ApiPluginTemplateComponent,
      children: [
        {
          path: '',
          component: ApiPluginTemplateListComponent
        },
        {
          path: 'create',
          component: ApiPluginTemplateCreateComponent
        },

        {
          path: 'content/:pluginTemplateId',
          component: ApiPluginTemplateContentComponent,
          children: [
            {
              path: '',
              component: ApiPluginTemplatePublishComponent
            },
            {
              path: 'message',
              component: ApiPluginTemplateMessageComponent
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
export class ApiRoutingModule { }

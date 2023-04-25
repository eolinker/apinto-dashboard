import { NgModule } from '@angular/core'
import { RouterModule, Routes } from '@angular/router'
import { ApiListComponent } from './api-list/api-list.component'
import { ApiContentComponent } from './api-list/content/content.component'
import { ApiCreateComponent } from './api-list/create/create.component'
import { ApiManagementComponent } from './api-list/group/group.component'
import { ApiManagementListComponent } from './api-list/list/list.component'
import { ApiMessageComponent } from './api-list/message/message.component'
import { ApiPublishComponent } from './api-list/publish/single/publish.component'
import { ApiPluginTemplateContentComponent } from './plugin/content/content.component'
import { ApiPluginTemplateCreateComponent } from './plugin/create/create.component'
import { ApiPluginTemplateListComponent } from './plugin/list/list.component'
import { ApiPluginTemplateMessageComponent } from './plugin/message/message.component'
import { ApiPluginTemplateComponent } from './plugin/plugin.component'
import { ApiPluginTemplatePublishComponent } from './plugin/publish/publish.component'
import { RouterComponent } from './router/router.component'

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
            },
            {
              path: 'create',
              component: ApiCreateComponent,
              children: [{
                path: ':apiGroupId',
                component: ApiCreateComponent
              }]
            },
            {
              path: 'message/:apiId',
              component: ApiMessageComponent
            }
          ]
        },
        {
          path: 'create',
          component: ApiCreateComponent,
          children: [{
            path: ':apiGroupId',
            component: ApiCreateComponent
          }]
        },
        {
          path: 'content/:apiId',
          component: ApiContentComponent,
          children: [
            {
              path: '',
              component: ApiMessageComponent
            },
            {
              path: 'publish',
              component: ApiPublishComponent
            }
          ]
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

import { NgModule } from '@angular/core'
import { RouterModule, Routes } from '@angular/router'
import { ApiPluginTemplateContentComponent } from './content/content.component'
import { ApiPluginTemplateCreateComponent } from './create/create.component'
import { ApiPluginTemplateListComponent } from './list/list.component'
import { ApiPluginTemplateMessageComponent } from './message/message.component'
import { ApiPluginTemplateComponent } from './plugin.component'
import { ApiPluginTemplatePublishComponent } from './publish/publish.component'
const routes: Routes = [
  {
    path: '',
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

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class PluginTemplateRoutingModule { }

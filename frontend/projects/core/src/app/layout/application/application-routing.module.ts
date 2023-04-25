import { NgModule } from '@angular/core'
import { RouterModule, Routes } from '@angular/router'
import { ApplicationManagementComponent } from './application.component'
import { ApplicationAuthenticationComponent } from './authentication/authentication.component'
import { ApplicationContentComponent } from './content/content.component'
import { ApplicationCreateComponent } from './create/create.component'
import { ApplicationManagementListComponent } from './list/list.component'
import { ApplicationMessageComponent } from './message/message.component'
import { ApplicationPublishComponent } from './publish/publish.component'

const routes: Routes = [{
  path: '',
  component: ApplicationManagementComponent,
  data: {
    id: '3'
  },
  children: [
    {
      path: '',
      component: ApplicationManagementListComponent
    },
    { path: 'create', component: ApplicationCreateComponent },
    {
      path: 'content/:appId',
      component: ApplicationContentComponent,
      children: [
        {
          path: '',
          component: ApplicationPublishComponent
        },
        {
          path: 'message',
          component: ApplicationMessageComponent
        },
        {
          path: 'authentication',
          component: ApplicationAuthenticationComponent
        }
      ]
    },
    {
      path: '**',
      component: ApplicationManagementListComponent
    }
  ]
}]

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class ApplicationRoutingModule { }

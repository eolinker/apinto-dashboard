import { NgModule } from '@angular/core'
import { RouterModule, Routes } from '@angular/router'
import { ApiContentComponent } from './content/content.component'
import { ApiCreateComponent } from './create/create.component'
import { ApiManagementComponent } from './group/group.component'
import { ApiManagementListComponent } from './list/list.component'
import { ApiMessageComponent } from './message/message.component'
import { ApiPublishComponent } from './publish/single/publish.component'
import { RouterComponent } from './router/router.component'

const routes: Routes = [{
  path: '',
  component: RouterComponent,
  data: { id: '4' },
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
}]

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class ApiRoutingModule { }

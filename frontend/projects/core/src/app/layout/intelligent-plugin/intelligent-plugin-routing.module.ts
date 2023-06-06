import { NgModule } from '@angular/core'
import { RouterModule, Routes } from '@angular/router'
import { IntelligentPluginComponent } from './intelligent-plugin.component'
import { IntelligentPluginLayoutComponent } from './layout/layout.component'

const routes: Routes = [{
  path: '',
  component: IntelligentPluginComponent,
  children: [
    {
      path: ':moduleName',
      children: [{
        path: '**',
        component: IntelligentPluginLayoutComponent
      }
      ]
    }
  ]
}]

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class IntelligentPluginRoutingModule { }

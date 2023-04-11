import { NgModule } from '@angular/core'
import { RouterModule, Routes } from '@angular/router'
import { InterceptorComponent } from './interceptor.component'

const routes: Routes = [{
  path: '',
  component: InterceptorComponent
}]

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class InterceptorRoutingModule { }

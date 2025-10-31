
import { NgModule } from '@angular/core'
import { RouterModule, Routes } from '@angular/router'
import { AuthInfoComponent } from './auth-info.component'

const routes: Routes = [
  {
    path: '',
    component: AuthInfoComponent
  }]

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class AuthInfoRoutingModule { }

/*
 * @Date: 2023-12-12 18:57:19
 * @LastEditors: maggieyyy
 * @LastEditTime: 2023-12-14 17:24:00
 * @FilePath: \apinto\projects\apinto-auth\src\app\app-routing.module.ts
 */
import { NgModule } from '@angular/core'
import { RouterModule, Routes } from '@angular/router'
import { AuthActivationComponent } from './layout/auth/activation/activation.component'
import { AuthUpdateComponent } from './layout/auth/update/update.component'
import { AuthInfoComponent } from './layout/auth-info/auth-info.component'

const routes: Routes = [
  {
    path: '',
    component: AuthActivationComponent,
    children: [
      {
        path: 'info',
        component: AuthInfoComponent
      },
      {
        path: 'update',
        component: AuthUpdateComponent
      }
    ]
  }]

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }

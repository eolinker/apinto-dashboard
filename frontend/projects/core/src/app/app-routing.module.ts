/*
 * @Author:
 * @Date: 2022-07-11 23:20:14
 * @LastEditors:
 * @LastEditTime: 2022-09-20 23:14:26
 * @FilePath: /apinto/src/app/app-routing.module.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { NgModule } from '@angular/core'
import { RouterModule, Routes } from '@angular/router'
import { RedirectPageService } from './service/redirect-page.service'
import { AuthGuardService } from './service/auth-guard.service'
import { LoginComponent } from './layout/login/login.component'
import { BasicLayoutComponent } from './layout/basic-layout/basic-layout.component'
import { CustomPreloadingStrategy } from './custom-preloading-strategy'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
const routes: Routes = [
  {
    path: 'login',
    component: LoginComponent
  },
  {
    path: '',
    component: BasicLayoutComponent,
    children: [
      {
        path: 'deploy',
        data: {
          id: '1'
        },
        loadChildren: () => import('./layout/deploy/deploy.module').then(m => m.DeployModule)
      },
      {
        path: 'upstream',
        data: {
          id: '2'
        },
        loadChildren: () => import('./layout/upstream/upstream.module').then(m => m.UpstreamModule)
      },

      {
        path: 'application',
        data: {
          id: '3'
        },
        loadChildren: () => import('./layout/application/application.module').then(m => m.ApplicationModule)
      },

      {
        path: 'router',
        data: {
          id: '4'
        },
        loadChildren: () => import('./layout/api/api.module').then(m => m.ApiModule)
      },
      {
        path: 'serv-governance',
        data: {
          id: '5'
        },
        loadChildren: () => import('./layout/serv-governance/serv-governance.module').then(m => m.ServGovernanceModule)
      },
      {
        path: 'system',
        data: {
          id: '6'
        },
        loadChildren: () => import('./layout/system/system.module').then(m => m.SystemModule)
      },
      {
        path: 'audit-log',
        data: {
          id: '7'
        },
        loadChildren: () => import('./layout/audit-log/audit-log.module').then(m => m.AuditLogModule)
      }
    ]
  }

]

@NgModule({
  imports: [RouterModule.forRoot(routes, { preloadingStrategy: CustomPreloadingStrategy })],
  exports: [RouterModule],
  providers: [AuthGuardService, RedirectPageService, CustomPreloadingStrategy, EoNgFeedbackMessageService]
})
export class AppRoutingModule { }

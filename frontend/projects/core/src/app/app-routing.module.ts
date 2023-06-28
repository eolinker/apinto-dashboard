/*
 * @Author: maggieyyy im.ymj@hotmail.com
 * @Date: 2022-07-11 23:20:14
 * @LastEditors: MengjieYang yangmengjie@eolink.com
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
import { AuthActivationComponent } from './layout/auth/activation/activation.component'
import { LocalPluginComponent } from './layout/local-plugin/local-plugin.component'
import { GuideComponent } from './layout/guide/guide.component'
import { DynamicDemoComponent } from './layout/dynamic-demo/dynamic-demo.component'
import { environment } from '../environments/environment'
import { NotFoundPageComponent } from './layout/not-found-page/not-found-page.component'
import { RemotePluginComponent } from './layout/remote-plugin/remote-plugin.component'
import { AuthInfoComponent } from './layout/auth/info/info.component'
const routes: Routes = [
  {
    path: 'auth',
    component: AuthActivationComponent
  },
  {
    path: 'login',
    component: LoginComponent
  },
  {
    path: '',
    pathMatch: 'full',
    component: LoginComponent
  },
  {
    path: '',
    component: BasicLayoutComponent,
    children: [
      {
        path: 'guide',
        component: GuideComponent
      },
      ...(environment.isBusiness
        ? [{
            path: 'auth-info',
            component: AuthInfoComponent
          }]
        : []),
      {
        path: 'deploy',
        data: {
          id: '1'
        },
        loadChildren: () => import('./layout/deploy/deploy.module').then(m => m.DeployModule)
      },
      {
        path: 'application',
        data: {
          id: '3'
        },
        pathMatch: 'prefix',
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
      },
      {
        path: 'module-plugin',
        data: {
          id: '10'
        },
        loadChildren: () => import('./layout/plugin/plugin-management.module').then(m => m.PluginManagementModule)
      },
      {
        path: 'log',
        data: {
          id: '11'
        },
        loadChildren: () => import('./layout/log-retrieval/log-retrieval.module').then(m => m.LogRetrievalModule)
      },
      {
        path: 'dynamic-demo',
        data: {

        },
        children: [{
          path: '**',
          component: DynamicDemoComponent
        }
        ],
        component: DynamicDemoComponent
      },
      {
        path: 'module',
        children: [
          {
            path: ':moduleName',
            children: [{ path: ':subPath', component: LocalPluginComponent }, {
              path: '**',
              component: LocalPluginComponent
            }
            ]
          }
        ]
      },
      {
        path: 'remote',
        children: [
          {
            path: ':moduleName',
            children: [{ path: ':subPath', component: RemotePluginComponent }, {
              path: '**',
              component: RemotePluginComponent
            }
            ]
          }
        ]
      },
      {
        path: 'template',
        loadChildren: () => import('./layout/intelligent-plugin/intelligent-plugin.module').then(m => m.IntelligentPluginModule)
      },
      {
        path: '**',
        data: {
        },
        component: NotFoundPageComponent
      }
    ]
  }

]

@NgModule({
  imports: [RouterModule.forRoot(routes, { preloadingStrategy: CustomPreloadingStrategy })],
  exports: [RouterModule],
  providers: [AuthGuardService, RedirectPageService, CustomPreloadingStrategy]
})
export class AppRoutingModule { }

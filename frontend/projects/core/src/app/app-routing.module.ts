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
import { AuthInfoComponent } from './layout/auth/info/info.component'
import { AuthUpdateComponent } from './layout/auth/update/update.component'
import { IframePageComponent } from './layout/iframe-page/iframe-page.component'
import { GuideComponent } from './layout/guide/guide.component'
import { DynamicDemoComponent } from './layout/dynamic-demo/dynamic-demo.component'
import { OuterComponent } from './layout/outer/outer.component'
import { IntelligentPluginLayoutComponent } from './component/intelligent-plugin/layout/layout.component'
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
    component: BasicLayoutComponent,
    children: [

      {
        path: 'guide',
        component: GuideComponent
      },
      {
        path: 'auth-info',
        component: AuthInfoComponent,
        data: {
          id: '8'
        }
      },

      {
        path: 'auth-update',
        component: AuthUpdateComponent
      },
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
      },
      {
        path: 'plugin',
        data: {
          id: '10'
        },
        loadChildren: () => import('./layout/plugin/plugin-management.module').then(m => m.PluginManagementModule)
      },
      {
        path: 'navigation',
        data: {
          id: '11'
        },
        loadChildren: () => import('./layout/navigation/navigation.module').then(m => m.NavigationModule)
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
        path: 'template',
        data: {
        },
        component: OuterComponent,
        children: [
          {
            path: 'iframe',
            component: IframePageComponent
          },
          {
            path: '**',
            component: IntelligentPluginLayoutComponent
          }
        ]
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

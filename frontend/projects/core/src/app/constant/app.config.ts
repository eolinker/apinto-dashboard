/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2022-07-28 22:12:29
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-07-02 10:26:00
 * @FilePath: \apinto\projects\core\src\app\constant\app.config.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

import { GuideComponent } from '../layout/guide/guide.component'
import { DynamicDemoComponent } from '../layout/dynamic-demo/dynamic-demo.component'
import { BasicLayoutComponent } from '../layout/basic-layout/basic-layout.component'
import { RemotePluginComponent } from '../layout/remote-plugin/remote-plugin.component'
import { PluginWrapperComponent } from '../layout/plugin-wrapper/plugin-wrapper.component'
import { SystemEmailConfigComponent } from '../layout/system/email/config/config.component'
import { SystemWebhookListComponent } from '../layout/system/webhook/list/list.component'
import { RedirectPage } from '../layout/redirect-page/redirect-page.component'
import { RouterLayoutComponent } from '../layout/router-layout/router-layout.component'

export const MODAL_NORMAL_SIZE: number = 900
export const MODAL_SMALL_SIZE: number = 600
export const MODAL_LARGE_SIZE: number = 1200

// 内置插件与对应组件/模块
export const routerMap:Map<string, any> = new Map([
  ['basicLayout', { type: 'component', component: BasicLayoutComponent }],
  ['redirectPage', { type: 'component', component: RedirectPage }],
  ['guide', { type: 'component', component: GuideComponent }],
  ['cluster', { type: 'module', module: () => import('projects/core/src/app/layout/deploy-cluster/deploy-cluster.module').then(m => m.DeployClusterModule) }],
  ['global-env', { type: 'module', module: () => import('projects/core/src/app/layout/global-env-var/global-env-var.module').then(m => m.GlobalEnvVarModule) }],
  ['node-plugin', { type: 'module', module: () => import('projects/core/src/app/layout/node-plugin/node-plugin.module').then(m => m.NodePluginModule) }],
  ['application', { type: 'module', module: () => import('projects/core/src/app/layout/application/application.module').then(m => m.ApplicationModule) }],
  ['api', { type: 'module', module: () => import('projects/core/src/app/layout/api/api.module').then(m => m.ApiModule) }],
  ['plugin-template', { type: 'module', module: () => import('projects/core/src/app/layout/plugin-template/plugin-template.module').then(m => m.PluginTemplateModule) }],
  ['traffic-strategy', { type: 'module', module: () => import('projects/core/src/app/layout/serv-governance/traffic-strategy/traffic-strategy.module').then(m => m.TrafficStrategyModule) }],
  ['fuse-strategy', { type: 'module', module: () => import('projects/core/src/app/layout/serv-governance/fuse-strategy/fuse-strategy.module').then(m => m.FuseStrategyModule) }],
  ['visit-strategy', { type: 'module', module: () => import('projects/core/src/app/layout/serv-governance/visit-strategy/visit-strategy.module').then(m => m.VisitStrategyModule) }],
  ['cache-strategy', { type: 'module', module: () => import('projects/core/src/app/layout/serv-governance/cache-strategy/cache-strategy.module').then(m => m.CacheStrategyModule) }],
  ['grey-strategy', { type: 'module', module: () => import('projects/core/src/app/layout/serv-governance/grey-strategy/grey-strategy.module').then(m => m.GreyStrategyModule) }],
  ['data-mask-strategy', { type: 'module', module: () => import('projects/core/src/app/layout/serv-governance/data-mask/data-mask.module').then(m => m.DataMaskModule) }],
  ['open-api', { type: 'module', module: () => import('projects/core/src/app/layout/open-api/open-api.module').then(m => m.OpenApiModule) }],
  ['email', { type: 'component', component: SystemEmailConfigComponent }],
  ['webhook', { type: 'component', component: SystemWebhookListComponent }],
  ['audit-log', { type: 'module', module: () => import('projects/core/src/app/layout/audit-log/audit-log.module').then(m => m.AuditLogModule) }],
  ['module-plugin', { type: 'module', module: () => import('projects/core/src/app/layout/plugin/plugin-management.module').then(m => m.PluginManagementModule) }],
  ['log', { type: 'module', module: () => import('projects/core/src/app/layout/log-retrieval/log-retrieval.module').then(m => m.LogRetrievalModule) }],
  ['intelligent', { type: 'module', module: () => import('projects/core/src/app/layout/intelligent-plugin/intelligent-plugin.module').then(m => m.IntelligentPluginModule) }],
  ['dynamic-demo', { type: 'component', component: DynamicDemoComponent }],
  ['remote', { type: 'component', component: RemotePluginComponent }],
  ['plugin-wrapper', { type: 'component', component: PluginWrapperComponent }],
  ['nav-hidden', { type: 'component', component: RouterLayoutComponent }]
])

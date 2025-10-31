/*
 * @Author: maggieyyy im.ymj@hotmail.com
 * @Date: 2022-07-11 23:20:14
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-07-02 10:27:27
 * @FilePath: \apinto\projects\core\src\app\app.module.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { HttpClient, HttpClientModule } from '@angular/common/http'
import { APP_INITIALIZER, NgModule, PlatformRef } from '@angular/core'
import { environment } from 'projects/core/src/environments/environment'

import { AppRoutingModule } from './app-routing.module'
import { AppComponent } from './app.component'
import { ComponentModule } from './component/component.module'
import { LayoutModule } from './layout/layout.module'
import { ApiService, API_URL } from './service/api.service'
import { NzTransferModule } from 'ng-zorro-antd/transfer'
import { httpInterceptorProviders } from './service/http-interceptors'
import { EoNgBreadcrumbModule } from 'eo-ng-breadcrumb'
import { EoNgFeedbackDrawerModule, EoNgFeedbackMessageService, EoNgFeedbackModalModule } from 'eo-ng-feedback'
import { EoNgLayoutModule } from 'eo-ng-layout'
import { EoNgMenuModule } from 'eo-ng-menu'
import { NzConfig, NZ_CONFIG } from 'ng-zorro-antd/core/config'
import { CommonModule, registerLocaleData } from '@angular/common'
import zh from '@angular/common/locales/en-GB'
import { NzSpinModule } from 'ng-zorro-antd/spin'
import { ChangeWordColorPipe } from './pipe/change-word-color.pipe'
import { MarkdownModule, MarkedOptions, MarkedRenderer } from 'ngx-markdown'
import { EoNgApintoTableModule } from 'projects/eo-ng-apinto-table/src/public-api'
import { NzFormModule } from 'ng-zorro-antd/form'
import { platformBrowserDynamic } from '@angular/platform-browser-dynamic'
import { PlatformProviderService } from './service/platform-provider.service'
import { PluginLoaderService } from './service/plugin-loader.service'
import { ScrollingModule } from '@angular/cdk/scrolling'
import { OverlayModule, ScrollStrategyOptions } from '@angular/cdk/overlay'
import { NzOverlayModule } from 'ng-zorro-antd/core/overlay'
import { FormsModule, ReactiveFormsModule } from '@angular/forms'
import { EoNgTreeModule } from 'eo-ng-tree'
import { NoopAnimationsModule } from '@angular/platform-browser/animations'
import { cloneDeep } from 'lodash-es'
import { AsIframeService } from './service/as-iframe.service'

if (typeof window.structuredClone !== 'function') {
  console.warn(' 浏览器版本过低，部分功能可能无法正常使用，请及时升级浏览器版本')
  window.structuredClone = function (obj:any) {
    // 使用 Lodash 的 _.cloneDeep 或者其他自定义的深拷贝方法
    return cloneDeep(obj)
  }
}
registerLocaleData(zh)
const ngZorroConfig: NzConfig = {
  // 注意组件名称没有 nz 前缀
  message: { nzMaxStack: 1, nzDuration: 2000 },
  notification: { nzMaxStack: 1, nzDuration: 2000 }
}

export function markedOptionsFactory (): MarkedOptions {
  const renderer = new MarkedRenderer()
  const linkRenderer = renderer.link
  renderer.link = (href, title, text) => {
    let html = linkRenderer.call(renderer, href, title, text)
    html = html.replace(/^<a /, '<a role="link"  tabindex="0" target="_blank" rel="nofollow noopener noreferrer" ')
    return html
  }

  return {
    renderer: renderer,
    gfm: true,
    breaks: false,
    pedantic: false,
    smartLists: true,
    smartypants: false
  }
}

export function initializeApp (pluginService: PluginLoaderService, asIframeService:AsIframeService) {
  (window as any)._apinto_mf = true
  const isIframe = window.self !== window.top
  if (isIframe) {
    asIframeService.startReceiveMessage()
  }
  return (): Promise<void> => {
    // 确保 loadPlugins 返回一个 Promise
    return pluginService.loadPlugins()
  }
}

@NgModule({
  declarations: [
    AppComponent,
    ChangeWordColorPipe
  ],
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    EoNgBreadcrumbModule,
    EoNgFeedbackDrawerModule,
    EoNgLayoutModule,
    EoNgMenuModule,
    EoNgTreeModule,
    EoNgFeedbackModalModule,
    HttpClientModule,
    EoNgApintoTableModule,
    NzFormModule,
    LayoutModule,
    ComponentModule,
    NzTransferModule,
    AppRoutingModule,
    NzSpinModule,
    ScrollingModule,
    OverlayModule,
    MarkdownModule.forRoot({
      loader: HttpClient,
      markedOptions: {
        provide: MarkedOptions,
        useFactory: markedOptionsFactory
      }
    }),
    NoopAnimationsModule,
    NzOverlayModule
  ],
  providers: [
    ApiService,
    EoNgFeedbackMessageService,
    { provide: API_URL, useValue: environment.urlPrefix },
    { provide: NZ_CONFIG, useValue: ngZorroConfig },
    httpInterceptorProviders,
    {
      provide: PlatformRef,
      useFactory: () => platformBrowserDynamic().bootstrapModule(AppModule)
    },
    PlatformProviderService,
    {
      provide: APP_INITIALIZER,
      useFactory: initializeApp,
      deps: [PluginLoaderService, AsIframeService],
      multi: true
    },
    ScrollStrategyOptions],
  bootstrap: [AppComponent]
})
export class AppModule { }

/*
 * @Author:  
 * @Date: 2022-07-11 23:20:14
 * @LastEditors:  
 * @LastEditTime: 2022-07-12 00:15:53
 * @FilePath: /apinto/src/app/app.module.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { HttpClientModule } from '@angular/common/http'
import { NgModule } from '@angular/core'
import { BrowserModule } from '@angular/platform-browser'
import { environment } from 'projects/core/src/environments/environment'

import { AppRoutingModule } from './app-routing.module'
import { AppComponent } from './app.component'
import { ComponentModule } from './component/component.module'
import { LayoutModule } from './layout/layout.module'
import { ApiService, API_URL } from './service/api.service'
import { NzTransferModule } from 'ng-zorro-antd/transfer'
import { httpInterceptorProviders } from './service/http-interceptors'
import { EoNgBreadcrumbModule } from 'eo-ng-breadcrumb'
import { EoNgFeedbackDrawerModule } from 'eo-ng-feedback'
import { EoNgLayoutModule } from 'eo-ng-layout'
import { EoNgMenuModule } from 'eo-ng-menu'
import { EoNgTableModule } from 'eo-ng-table'
import { NzConfig, NZ_CONFIG } from 'ng-zorro-antd/core/config'
import { registerLocaleData } from '@angular/common'
import zh from '@angular/common/locales/en-GB'
import { NgxEchartsModule } from 'ngx-echarts'
import { NzSpinModule } from 'ng-zorro-antd/spin'
import { ChangeWordColorPipe } from './pipe/change-word-color.pipe'

registerLocaleData(zh)
const ngZorroConfig: NzConfig = {
  // 注意组件名称没有 nz 前缀
  message: { nzMaxStack: 1, nzDuration: 2000 },
  notification: { nzMaxStack: 1, nzDuration: 2000 }
}

@NgModule({
  declarations: [
    AppComponent,
    ChangeWordColorPipe
  ],
  imports: [
    BrowserModule,
    EoNgBreadcrumbModule,
    EoNgFeedbackDrawerModule,
    EoNgLayoutModule,
    EoNgMenuModule,
    EoNgTableModule,
    HttpClientModule,
    LayoutModule,
    ComponentModule,
    NzTransferModule,
    AppRoutingModule,
    NgxEchartsModule.forRoot({
      echarts: () => import('echarts')
    }),
    NzSpinModule
  ],
  providers: [ApiService,
    { provide: API_URL, useValue: environment.urlPrefix },
    { provide: NZ_CONFIG, useValue: ngZorroConfig },
    httpInterceptorProviders],
  bootstrap: [AppComponent]
})
export class AppModule { }

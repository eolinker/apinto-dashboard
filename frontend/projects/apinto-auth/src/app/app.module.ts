/*
 * @Date: 2023-12-12 18:57:19
 * @LastEditors: maggieyyy
 * @LastEditTime: 2023-12-14 17:29:37
 * @FilePath: \apinto\projects\apinto-auth\src\app\app.module.ts
 */
import { NgModule } from '@angular/core'
import { BrowserModule } from '@angular/platform-browser'

import { AppRoutingModule } from './app-routing.module'
import { AppComponent } from './app.component'
import { ApiService } from './service/api.service'
import { AuthModule } from './layout/auth/auth.module'
import { ModuleFederationService } from './service/module-federation.service'
import { AuthInfoModule } from './layout/auth-info/auth-info.module'

let coreUserService:any

@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    AuthModule,
    AuthInfoModule

  ],
  providers: [ApiService],
  bootstrap: [AppComponent]
})
export class AppModule {
  constructor (private moduleFederationService:ModuleFederationService) {
    if (!this.moduleFederationService.providerFromCore) {
      this.moduleFederationService.providerFromCore = coreUserService
      this.moduleFederationService.initialized = true
    }
  }
}

export function bootstrap (props:any) {
  coreUserService = props.pluginProvider
}
export async function beforeMount (props:any) {
  coreUserService.redirectUrl = props.redirectUrl
};
export async function mount (props:any) {
};
export function beforeUnmount (props:any) { }
export function unmount (props:any) { }

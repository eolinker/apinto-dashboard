import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { EoNgButtonModule } from 'eo-ng-button'
import { EoNgCopyModule } from 'eo-ng-copy'
import { EoNgFeedbackAlertModule, EoNgFeedbackTooltipModule } from 'eo-ng-feedback'
import { EoNgInputModule } from 'eo-ng-input'
import { EoNgLayoutModule } from 'eo-ng-layout'
import { NzUploadModule } from 'ng-zorro-antd/upload'
import { AuthInfoComponent } from './auth-info.component'
import { FormsModule, ReactiveFormsModule } from '@angular/forms'
import { ModuleFederationService } from '../../service/module-federation.service'
import { AuthInfoRoutingModule } from './auth-info-routing.module'
let coreUserService:any

@NgModule({
  declarations: [AuthInfoComponent],
  imports: [
    CommonModule,
    AuthInfoRoutingModule,
    FormsModule,
    ReactiveFormsModule,
    EoNgLayoutModule,
    EoNgInputModule,
    EoNgButtonModule,
    EoNgCopyModule,
    NzUploadModule,
    EoNgFeedbackAlertModule,
    EoNgCopyModule,
    EoNgFeedbackTooltipModule
  ],
  exports: [AuthInfoComponent]
})
export class AuthInfoModule {
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

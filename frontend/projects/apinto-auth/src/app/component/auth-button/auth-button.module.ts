import { NgModule } from '@angular/core'
import { AuthButtonComponent } from './auth-button.component'
import { BrowserModule } from '@angular/platform-browser'
import { EoNgButtonModule } from 'eo-ng-button'
import { EoNgCopyModule } from 'eo-ng-copy'
import { EoNgFeedbackAlertModule, EoNgFeedbackTooltipModule } from 'eo-ng-feedback'
import { ApiService } from '../../service/api.service'
import { RouterModule } from '@angular/router'

@NgModule({
  declarations: [AuthButtonComponent],
  imports: [
    BrowserModule,
    EoNgButtonModule,
    EoNgFeedbackAlertModule,
    EoNgCopyModule,
    EoNgFeedbackTooltipModule,
    RouterModule
  ],
  providers: [
    ApiService
  ]
})
export class AuthButtonModule {
  ngDoBootstrap () {}
}

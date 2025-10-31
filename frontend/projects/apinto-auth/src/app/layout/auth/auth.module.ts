import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { AuthActivationComponent } from './activation/activation.component'
import { AuthUpdateComponent } from './update/update.component'
import { EoNgLayoutModule } from 'eo-ng-layout'
import { FormsModule, ReactiveFormsModule } from '@angular/forms'
import { EoNgInputModule } from 'eo-ng-input'
import { EoNgButtonModule } from 'eo-ng-button'
import { EoNgCopyModule } from 'eo-ng-copy'
import { NzUploadModule } from 'ng-zorro-antd/upload'
import { AuthInfoDetailComponent } from './info/detail/detail.component'
import { EoNgFeedbackAlertModule, EoNgFeedbackTooltipModule } from 'eo-ng-feedback'

@NgModule({
  declarations: [
    AuthInfoDetailComponent,
    AuthUpdateComponent,
    AuthActivationComponent
  ],
  imports: [
    CommonModule,
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
  ]
})
export class AuthModule { }

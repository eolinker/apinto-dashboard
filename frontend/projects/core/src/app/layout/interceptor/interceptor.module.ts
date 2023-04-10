import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { DragDropModule } from '@angular/cdk/drag-drop'

import { InterceptorRoutingModule } from './interceptor-routing.module'
import { InterceptorComponent } from './interceptor.component'
import { ComponentModule } from '../../component/component.module'
import { EoNgButtonModule } from 'eo-ng-button'
import { EoNgInputModule } from 'eo-ng-input'
import { FormsModule } from '@angular/forms'
import { NzListModule } from 'ng-zorro-antd/list'
import { EoNgEmptyModule } from 'eo-ng-empty'

@NgModule({
  declarations: [InterceptorComponent],
  imports: [
    CommonModule,
    FormsModule,
    InterceptorRoutingModule,
    ComponentModule,
    DragDropModule,
    EoNgButtonModule,
    EoNgInputModule,
    NzListModule,
    EoNgEmptyModule
  ],
  exports: [
    InterceptorComponent
  ]
})
export class InterceptorModule { }

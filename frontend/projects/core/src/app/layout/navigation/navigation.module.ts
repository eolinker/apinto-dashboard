import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { NzListModule } from 'ng-zorro-antd/list'
import { NavigationRoutingModule } from './navigation-routing.module'
import { NavigationComponent } from './navigation.component'
import { EoNgButtonModule } from 'eo-ng-button'
import { NavigationCreateComponent } from './create/create.component'
import { FormsModule, ReactiveFormsModule } from '@angular/forms'
import { NzUploadModule } from 'ng-zorro-antd/upload'
import { EoNgApintoTableModule } from 'projects/eo-ng-apinto-table/src/public-api'
import { EoNgInputModule } from 'eo-ng-input'
import { NzFormModule } from 'ng-zorro-antd/form'
import { DragDropModule } from '@angular/cdk/drag-drop'
import { EoNgEmptyModule } from 'eo-ng-empty'
import { NzIconModule } from 'ng-zorro-antd/icon'

@NgModule({
  declarations: [
    NavigationComponent,
    NavigationCreateComponent
  ],
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    NzUploadModule,
    NavigationRoutingModule,
    NzListModule,
    NzFormModule,
    EoNgButtonModule,
    EoNgApintoTableModule,
    EoNgInputModule,
    DragDropModule,
    EoNgEmptyModule,
    NzIconModule
  ],
  exports: [
    NavigationComponent
  ]
})
export class NavigationModule { }

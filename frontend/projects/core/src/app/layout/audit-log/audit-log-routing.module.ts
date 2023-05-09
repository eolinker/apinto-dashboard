import { NgModule } from '@angular/core'
import { RouterModule, Routes } from '@angular/router'
import { AuditLogComponent } from './audit-log.component'
import { AuditLogListComponent } from './list/list.component'

const routes: Routes = [{
  path: '',
  component: AuditLogComponent,
  data: { id: '7' },
  children: [
    {
      path: '',
      component: AuditLogListComponent
    }
  ]
}]

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class AuditLogRoutingModule { }

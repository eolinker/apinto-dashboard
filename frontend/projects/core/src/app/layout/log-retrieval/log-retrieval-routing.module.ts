import { NgModule } from '@angular/core'
import { RouterModule, Routes } from '@angular/router'
import { LogRetrievalComponent } from './log-retrieval.component'

const routes: Routes = [
  {
    path: '',
    component: LogRetrievalComponent
  }
]

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class LogRetrievalRoutingModule { }

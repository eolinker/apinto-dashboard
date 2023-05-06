import { Component } from '@angular/core'
import { ApiWebsocketCreateComponent } from '../websocket-create/websocket-create.component'

@Component({
  selector: 'eo-ng-api-http-create',
  templateUrl: './http-create.component.html',
  styles: [
  ]
})
export class ApiHttpCreateComponent extends ApiWebsocketCreateComponent {
  // 提交api数据
  override saveApi () {
    if (this.createApiForm.method.length === 0 && !this.allChecked) {
      this.showCheckboxGroupValid = true
      return
    } else {
      this.showCheckboxGroupValid = false
    }
    super.saveApi()
  }
}

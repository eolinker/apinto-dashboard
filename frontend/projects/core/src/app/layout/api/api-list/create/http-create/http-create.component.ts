import { Component } from '@angular/core'
import { ApiWebsocketCreateComponent } from '../websocket-create/websocket-create.component'
import { setFormValue } from 'projects/core/src/app/constant/form'

@Component({
  selector: 'eo-ng-api-http-create',
  templateUrl: './http-create.component.html',
  styles: [
    `
      eo-ng-table.ant-table {
        min-width: 508px !important;
      }

      .ant-form-item-control:first-child:not([class^='ant-col-']):not(
          [class*=' ant-col-']
        ) {
        width: auto !important;
      }

      nz-form-item.ant-row.checkbox-group-api.ant-form-item.ant-form-item-has-error {
        margin-bottom: 0 !important;
      }

      :host ::ng-deep{
        .arrayItem.hosts input{
          width:508px;
        }
      }
    `
  ]
})
export class ApiHttpCreateComponent extends ApiWebsocketCreateComponent {
  // 当编辑api时，需要获取api信息
  override getApiMessage () {
    this.api.get('router', { uuid: this.apiUuid }).subscribe((resp) => {
      if (resp.code === 0) {
        setFormValue(this.validateForm, resp.data.api)
        // eslint-disable-next-line dot-notation
        this.validateForm.controls['requestPath'].setValue(resp.data.api.requestPath.slice(1))
        this.createApiForm = resp.data.api
        if (
          !this.createApiForm.method ||
          this.createApiForm.method.length === 0
        ) {
          this.createApiForm.method = [
            'POST',
            'PUT',
            'GET',
            'DELETE',
            'PATCH',
            'HEAD',
            'OPTIONS'
          ]
          this.allChecked = true
          this.updateAllChecked()
        } else {
          this.initCheckbox()
        }
        this.getHeaderList()
        this.hostsList = [...resp.data.api.hosts?.map((x:string) => ({ key: x })) || [], { key: '' }]
      }
    })
  }

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

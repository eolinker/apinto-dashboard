import { Component } from '@angular/core'
import { SelectOption } from 'eo-ng-select'
// 该组件是为了方便开发react组件使用的demo
@Component({
  selector: 'eo-ng-formily2-react',
  template: `
  <div id="myReactComponentContainer">
    <div>
  schema：<eo-ng-codebox
        style="display: block"
        [(code)]="renderSchema"
        mode="json"
      ></eo-ng-codebox>
  <button (click)="changeSchema()">生成表单</button>
</div>
  form：<eo-ng-codebox
        style="display: block"
        [(code)]="initFormValue"
        mode="json"
      ></eo-ng-codebox>
   <formily2-react-wrapper
        #formily
        [editPage]="true"
        [initFormValue]="initFormValue"
        [driverSelectOptions]="driverSelectOptions"
        [demoSchema]="schema"
        [demo]="true"
       ></formily2-react-wrapper>
 </div>
  `,
  styles: [
    `
    #myReactComponentContainer{
      width:700px;
      margin:auto;
      margin-top:10px;
    }
    `
  ]
})
export class Formily2ReactComponent {
  renderSchema:any
  initFormValue:any
  driverSelectOptions:SelectOption[] = []
  schema:any
  changeSchema () {
    this.schema = JSON.parse(this.renderSchema)
  }
}

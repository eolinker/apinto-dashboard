import { Component } from '@angular/core'
import { SelectOption } from 'eo-ng-select'
// 该组件是为了方便开发react组件使用的demo
@Component({
  selector: 'eo-ng-formily2-react',
  template: `
  <div id="myReactComponentContainer" class="flex">
    <div class="w-[50%] mr-[50px]">

  schema：<eo-ng-codebox
        style="display: block;width:550px"
        [(code)]="renderSchema"
        mode="json"
      ></eo-ng-codebox>

  form：<eo-ng-codebox
        style="display: block"
        [(code)]="initFormValue"
        mode="json"
      ></eo-ng-codebox>
  <button (click)="changeSchema()">生成表单</button>
</div>
   <formily2-react-wrapper
        #formily
        [editPage]="true"
        [initFormValue]="form"
        [driverSelectOptions]="driverSelectOptions"
        [demoSchema]="schema"
        [demo]="true"
       ></formily2-react-wrapper>

 </div>

  `,
  styles: [
    `
    #myReactComponentContainer{
      width:100%;
      height:100%;
    }
    formily2-react-wrapper{
      width:100%;
      border-left:1px solid;
      padding-left:50px;
      padding-top:10px;
    }
    `
  ]
})
export class Formily2ReactComponent {
  renderSchema:any
  initFormValue:any
  driverSelectOptions:SelectOption[] = []
  schema:any
  form:any
  changeSchema () {
    console.log(this.renderSchema)
    if (!this.renderSchema) {
      return
    }
    this.schema = JSON.parse(this.renderSchema)
    try {
      this.form = JSON.parse(this.initFormValue)
    } catch {
      console.warn('dynamic-demo生成表单失败，请检查form值格式')
    }
  }
}

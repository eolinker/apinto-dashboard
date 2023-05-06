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
<div id="render-block">
   <formily2-react-wrapper
        #formily
        [editPage]="true"
        [initFormValue]="form"
        [driverSelectOptions]="driverSelectOptions"
        [demoSchema]="schema"
        [demo]="true"
        (onSubmit)="handlerSubmit($event)"
       ></formily2-react-wrapper>

       <div  class="mt-[50px] " *ngIf="value">
          <p class="font-bold">提交后的数据如下：</p>
          <p class="border-[1px] border-solid border-BORDER">{{value}}
</p>
    </div>
<div>
 </div>

  `,
  styles: [
    `
    #myReactComponentContainer{
      width:100%;
      height:100%;
    }
    #render-block{
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
  value:any
  handlerSubmit ($event:any) {
    console.log($event)
    this.value = JSON.stringify($event)
  }

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

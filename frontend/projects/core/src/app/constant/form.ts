// 表单的常用方法
import { FormGroup } from '@angular/forms'
// 表单初始化

// 表单赋值
export function setFormValue (form:FormGroup, data:{[key:string]:any}):void {
  Object.keys(form.controls).forEach(key => {
    if (form.controls[key] instanceof FormGroup) {
      setFormValue(form.controls[key] as FormGroup, data[key])
    } else {
      form.controls[key].setValue(data[key])
    }
  })
}

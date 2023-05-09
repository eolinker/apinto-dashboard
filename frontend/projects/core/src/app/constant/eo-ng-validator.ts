import { Validators, ValidatorFn, AbstractControl } from '@angular/forms'
import { NzSafeAny } from 'ng-zorro-antd/core/types'

// current locale is key of the MyErrorsOptions
export type EoNgMyErrorsOptions = { 'zh-cn': string; en: string } & Record<string, NzSafeAny>;
export type EoNgMyValidationErrors = Record<string, EoNgMyErrorsOptions>;

export class EoNgMyValidators extends Validators {
  static override minLength (minLength: number): ValidatorFn {
    return (control: AbstractControl): EoNgMyValidationErrors | null => {
      if (Validators.minLength(minLength)(control) === null) {
        return null
      }
      return { minlength: { 'zh-cn': `最小长度为 ${minLength}`, en: `MinLength is ${minLength}` } }
    }
  }

  static override maxLength (maxLength: number): ValidatorFn {
    return (control: AbstractControl): EoNgMyValidationErrors | null => {
      if (Validators.maxLength(maxLength)(control) === null) {
        return null
      }
      return { maxlength: { 'zh-cn': `最大长度为 ${maxLength}`, en: `MaxLength is ${maxLength}` } }
    }
  }

  static maxByteLength (maxByteLength:number): ValidatorFn {
    return (control: AbstractControl): EoNgMyValidationErrors | null => {
      let len = 0
      for (let i = 0; i < control.value.length; i++) {
        const length = control.value.charCodeAt(i)
        if (length >= 0 && length <= 128) {
          len += 1
        } else {
          len += 2
        }
      }
      if (len <= maxByteLength) {
        return null
      }
      return { maxByteLength: { 'zh-cn': `最大字符长度为${maxByteLength}`, en: `MaxByteLength is ${maxByteLength}` } }
    }
  }

  static roleAccess (control:AbstractControl): EoNgMyValidationErrors | null {
    const value = control.value
    if (value.size > 0) {
      return null
    } else {
      return { roleAccess: { 'zh-cn': '角色权限不能为空', en: 'Not Empty' } }
    }
  }

  static priority (control:AbstractControl): EoNgMyValidationErrors | null {
    const value = control.value
    if (
      value &&
      (value < 1 ||
        value > 999 ||
        !/^\d+$/.test(value + ''))
    ) {
      return { priority: { 'zh-cn': '优先级范围需在1-999之间，且为整数', en: '优先级范围需在1-999之间，且为整数' } }
    }
    return null
  }
}

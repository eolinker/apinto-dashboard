import { Component, Input, OnInit } from '@angular/core'
import { AbstractControl, FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { SelectOption } from 'eo-ng-select'
import { defaultAutoTips } from '../../../constant/conf'
import { setFormValue } from '../../../constant/form'
import { optTypeList } from '../types/conf'

@Component({
  selector: 'eo-ng-proxy',
  templateUrl: './proxy.component.html',
  styles: [
  ]
})
export class ApiManagementProxyComponent implements OnInit {
  @Input() data:any = {}
  @Input() editPage:boolean = false
  nzDisabled:boolean = false
  listOfOptTypes:SelectOption[] = optTypeList
  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  validateProxyHeaderForm:FormGroup = new FormGroup({})
  constructor (private fb: UntypedFormBuilder) {
    this.validateProxyHeaderForm = this.fb.group({
      key: ['', [Validators.required]],
      value: [''],
      opt_type: ['', [Validators.required]]
    })
  }

  ngOnInit (): void {
    if (this.editPage) {
      setFormValue(this.validateProxyHeaderForm, this.data)
    }
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  valueValidator = (control:AbstractControl) => {
    // eslint-disable-next-line dot-notation
    if (this.validateProxyHeaderForm.controls['opt_type']?.value === 'DELETE') {
      return null
    } else {
      if (!control.value) {
        return { error: true, required: true }
      }
    }

    return null
  }

  changeValidator () {
    this.validateProxyHeaderForm.patchValue({
      key: '', value: ''
    })
    // eslint-disable-next-line dot-notation
    if (this.validateProxyHeaderForm.controls['opt_type'].value !== 'DELETE') {
      // eslint-disable-next-line dot-notation
      this.validateProxyHeaderForm.controls['value'].setValidators([Validators.required])
    } else {
    // eslint-disable-next-line dot-notation
      this.validateProxyHeaderForm.controls['value'].setValidators([])
    }
  }
}

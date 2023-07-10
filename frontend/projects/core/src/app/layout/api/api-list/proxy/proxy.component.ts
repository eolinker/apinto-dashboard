import { Component, Input, OnInit } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { SelectOption } from 'eo-ng-select'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { OptTypeList } from '../../types/conf'

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
  listOfOptTypes:SelectOption[] = [...OptTypeList]
  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  validateProxyHeaderForm:FormGroup = new FormGroup({})
  constructor (private fb: UntypedFormBuilder) {
    this.validateProxyHeaderForm = this.fb.group({
      key: ['', [Validators.required]],
      value: [''],
      optType: ['', [Validators.required]]
    })
  }

  ngOnInit (): void {
    if (this.editPage) {
      this.validateProxyHeaderForm.patchValue(this.data)
    }
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  changeValidator () {
    this.validateProxyHeaderForm.patchValue({
      key: '', value: ''
    })
    if (this.validateProxyHeaderForm.controls['optType'].value !== 'DELETE') {
      this.validateProxyHeaderForm.controls['value'].setValidators([Validators.required])
    } else {
      this.validateProxyHeaderForm.controls['value'].setValidators([])
    }
  }
}

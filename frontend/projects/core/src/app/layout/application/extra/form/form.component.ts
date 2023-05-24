/* eslint-disable dot-notation */
import { Component, EventEmitter, Input, Output } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { Router } from '@angular/router'
import { SelectOption } from 'eo-ng-select'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { extraConflictList, positionList } from '../../types/conf'
import { ExtraListData } from '../../types/types'

@Component({
  selector: 'eo-ng-application-extra-form',
  templateUrl: './form.component.html',
  styles: [
  ]
})
export class ApplicationExtraFormComponent {
  @Input() validateParamForm:FormGroup = new FormGroup({})
  @Input() editData?:ExtraListData
  @Input()
  set extraList (val) {
    this._extraList = val
    this.extraListChange.emit(this._extraList)
  }

  get extraList () {
    return this._extraList
  }

  @Output() eoNgCloseDrawer: EventEmitter<string> = new EventEmitter()
  @Output() extraListChange : EventEmitter<ExtraListData[]> = new EventEmitter()

  _extraList:ExtraListData[] = []
  positionList:SelectOption[] =[...positionList]
  conflictList:SelectOption[] =[...extraConflictList]
  matchHeaderSet:Set<string> = new Set()
  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  data:ExtraListData|undefined
  nzDisabled:boolean = false
  closeModal: Function | null = null

  constructor (private fb: UntypedFormBuilder, private router:Router) {
  }

  ngOnInit (): void {
    this.validateParamForm = this.fb.group({
      position: [this.data?.position || '', [Validators.required]],
      key: [this.data?.key || '', [Validators.required, Validators.pattern('^[a-zA-Z][a-zA-Z0-9-_]*')]],
      conflict: [this.data?.conflict || '', [Validators.required]],
      value: [this.data?.value || '', [Validators.required]]
    })
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  saveParam () {
    if (this.validateParamForm.valid) {
      if (!this.data) {
        if (this.matchHeaderSet.has(this.validateParamForm.controls['key'].value)) {
          for (const index in this.extraList) {
            if (this.extraList[index].key === this.validateParamForm.controls['key'].value && this.extraList[index].position === this.validateParamForm.controls['position'].value) {
              this.extraList.splice(Number(index), 1)
              break
            }
          }
        }
      } else {
        for (const index in this.extraList) {
          if (this.extraList[index].key === this.editData!.key && this.extraList[index].position === this.editData!.position && this.extraList[index].value === this.editData!.value && this.extraList[index].conflict === this.editData!.conflict) {
            this.extraList.splice(Number(index), 1)
            break
          }
        }
      }
      if (this.validateParamForm.controls['position'].value === 'HEADER') { this.matchHeaderSet.add(this.validateParamForm.controls['key'].value) }
      this.extraList = [{ position: this.validateParamForm.controls['position'].value, key: this.validateParamForm.controls['key'].value, value: this.validateParamForm.controls['value'].value, conflict: this.validateParamForm.controls['conflict'].value }, ...this.extraList]
      this.closeModal && this.closeModal()
    } else {
      Object.values(this.validateParamForm.controls).forEach(control => {
        if (control.invalid) {
          control.markAsDirty()
          control.updateValueAndValidity({ onlySelf: true })
        }
      })
    }
  }
}

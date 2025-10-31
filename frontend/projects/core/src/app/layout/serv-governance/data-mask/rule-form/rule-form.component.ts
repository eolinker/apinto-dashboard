import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { Router } from '@angular/router'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { MatchRules, DataFormatOptions, DataMaskReplaceStrOptions, DataMaskBaseOptions, DataMaskOrderOptions } from '../../types/conf'
import { EoNgMyValidators } from 'projects/core/src/app/constant/eo-ng-validator'
import { v4 as uuidv4 } from 'uuid'
import { MaskRuleData } from '../../types/types'

const InitFormData = {
  match: {
    type: '',
    value: ''
  },
  mask: {
    type: '',
    begin: null,
    length: null,
    replace: {
      type: '',
      value: ''
    }
  }
}

@Component({
  selector: 'eo-ng-data-mask-rule-form',
  templateUrl: './rule-form.component.html',
  styles: [
  ]
})
export class DataMaskRuleFormComponent implements OnInit {
  @Input() validateMatchForm:FormGroup = new FormGroup({})
  @Input() editData?:MaskRuleData
  @Input()
  set ruleList (val) {
    this._ruleList = val
    this.ruleListChange.emit(this._ruleList)
  }

  get ruleList () {
    return this._ruleList
  }

  @Output() eoNgCloseDrawer: EventEmitter<string> = new EventEmitter()
  @Output() ruleListChange : EventEmitter<MaskRuleData[]> = new EventEmitter()

  _ruleList:MaskRuleData[] = []
  nzDisabled:boolean = false

  readonly MatchRules = MatchRules;
  readonly DataFormatOptions = DataFormatOptions;
  readonly DataMaskReplaceStrOptions = DataMaskReplaceStrOptions;
  readonly DataMaskBaseOptions = DataMaskBaseOptions;
  readonly DataMaskOrderOptions = DataMaskOrderOptions;
  readonly autoTips = defaultAutoTips;

  constructor (private fb: UntypedFormBuilder, private router:Router) {

  }

  getMaskingRuleOptions (): any[] {
    const matchRuleType = this.validateMatchForm.get('match.type')?.value
    const matchRuleFormat = this.validateMatchForm.get('match.value')?.value
    return (matchRuleType === 'inner' && ['name', 'phone', 'id-card', 'bank-card'].includes(matchRuleFormat as string))
      ? this.DataMaskOrderOptions
      : this.DataMaskBaseOptions
  }

  ngOnInit (): void {
    this.validateMatchForm = this.fb.group({
      match: this.fb.group({
        type: ['', Validators.required],
        value: ['', Validators.required]
      }),
      mask: this.fb.group({
        type: ['', Validators.required],
        begin: [null],
        length: [null],
        replace: this.fb.group({
          type: [''],
          value: ['']
        })
      })
    })

    // 监听 'match.type' 和 'mask.type' 的值变化，并动态设置校验器
    this.addDynamicValidators()
    if (this.editData) {
      const newEditData = { match: { ...InitFormData.match, ...this.editData.match }, mask: { ...InitFormData.mask, ...this.editData.mask } }
      this.validateMatchForm.setValue(newEditData)
    }
  }

  addDynamicValidators () {
    // 添加 mask.type 的监听
    this.validateMatchForm.get('mask.type')?.valueChanges.subscribe(() => {
      const beginControl = this.validateMatchForm.get('mask.begin')
      const lengthControl = this.validateMatchForm.get('mask.length')
      const replaceTypeControl = this.validateMatchForm.get('mask.replace.type')
      const replaceValueControl = this.validateMatchForm.get('mask.replace.value')

      // 动态设置校验器
      if (['partial-display', 'partial-masking', 'truncation'].indexOf(this.validateMatchForm.get('mask.type')?.value) !== -1) {
        beginControl?.setValidators([Validators.required])
        lengthControl?.setValidators([Validators.required])
      } else {
        beginControl?.clearValidators()
        lengthControl?.clearValidators()
      }

      // 替换类型的校验
      if (this.validateMatchForm.get('mask.type')?.value === 'replacement') {
        replaceTypeControl?.setValidators([Validators.required])
        if (replaceTypeControl?.value === 'custom') {
          replaceValueControl?.setValidators([Validators.required])
        } else {
          replaceValueControl?.clearValidators()
        }
      } else {
        replaceTypeControl?.clearValidators()
        replaceValueControl?.clearValidators()
      }

      // 更新控件的校验状态
      beginControl?.updateValueAndValidity()
      lengthControl?.updateValueAndValidity()
      replaceTypeControl?.updateValueAndValidity()
      replaceValueControl?.updateValueAndValidity()
    })

    this.validateMatchForm.get('mask.replace.type')?.valueChanges.subscribe(() => {
      const replaceValueControl = this.validateMatchForm.get('mask.replace.value')
      if (this.validateMatchForm.get('mask.replace.type')?.value === 'custom') {
        replaceValueControl?.setValidators([Validators.required])
      } else {
        replaceValueControl?.clearValidators()
      }
      replaceValueControl?.updateValueAndValidity()
    })
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  prepareSubmitData (formData: any) {
    const submitData: any = {
      match: {
        type: formData.match.type,
        value: formData.match.value
      },
      mask: {
        type: formData.mask.type
      }
    }

    switch (formData.mask.type) {
      case 'replacement': {
        submitData.mask = {
          ...submitData.mask,
          replace: formData.mask.replace
        }
        break
      }
      case 'shuffling': {
        break
      }
      default: {
        submitData.mask.begin = formData.mask.begin
        submitData.mask.length = formData.mask.length
        break
      }
    }
    return submitData
  }

  save () {
    if (this.validateMatchForm.valid) {
      const formData = this.validateMatchForm.value
      const submitData = this.prepareSubmitData(formData)

      if (this.editData) {
        for (const index in this.ruleList) {
          if (this.ruleList[index].eoKey === this.editData!.eoKey) {
            this.ruleList.splice(Number(index), 1)
            break
          }
        }
      }
      this.ruleList = [{ ...submitData, eoKey: this.editData?.eoKey || uuidv4() }, ...this.ruleList]
      this.closeDrawer()
    } else {
      this.markFormGroupDirty(this.validateMatchForm)
    }
  }

  markFormGroupDirty (formGroup: FormGroup) {
    Object.values(formGroup.controls).forEach(control => {
      if (control instanceof FormGroup) {
        this.markFormGroupDirty(control)
      } else {
        if (control.invalid) {
          control.markAsDirty()
          control.updateValueAndValidity({ onlySelf: true })
        }
      }
    })
  }

  closeDrawer () {
    this.eoNgCloseDrawer.emit('rule')
  }
}

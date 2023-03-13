/* eslint-disable dot-notation */
import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { Router } from '@angular/router'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'

@Component({
  selector: 'eo-ng-match-form',
  templateUrl: './form.component.html',
  styles: [
  ]
})
export class MatchFormComponent implements OnInit {
  @Input() validateMatchForm:FormGroup = new FormGroup({})
  @Output() eoNgCloseDrawer: EventEmitter<any> = new EventEmitter()
  @Input()
  set matchList (val) {
    this._matchList = val
    this.matchListChange.emit(this._matchList)
  }

  get matchList () {
    return this._matchList
  }

  @Output() matchListChange : EventEmitter<any> = new EventEmitter()

  @Input() editData?:any

  _matchList:Array<any> = []
  positionList:Array<{label:string, value:string}>=[
    { label: 'HTTP请求头', value: 'header' },
    { label: '请求参数', value: 'query' },
    { label: 'Cookie', value: 'cookie' }
  ]

  matchTypeList:Array<{label:string, value:string}>=[
    { label: '全等匹配', value: 'EQUAL' },
    { label: '前缀匹配', value: 'PREFIX' },
    { label: '后缀匹配', value: 'SUFFIX' },
    { label: '字串匹配', value: 'SUBSTR' },
    { label: '非等匹配', value: 'UNEQUAL' },
    { label: '空值匹配', value: 'NULL' },
    { label: '存在匹配', value: 'EXIST' },
    { label: '不存在匹配', value: 'UNEXIST' },
    { label: '区分大小写的正则匹配', value: 'REGEXP' },
    { label: '不区分大小写的正则匹配', value: 'REGEXPG' },
    { label: '任意匹配', value: 'ANY' }
  ]

  matchHeaderSet:Set<string> = new Set()

  autoTips: Record<string, Record<string, string>> = defaultAutoTips

  // eslint-disable-next-line camelcase
  data:{position:string, key:string, match_type:string, pattern:string}|undefined
  accessUrl:string = '' // 用来判断权限的url
  nzDisabled:boolean = false

  constructor (private fb: UntypedFormBuilder, private router:Router) {
    switch (this.router.url.split('/')[1]) {
      case 'router':
        this.accessUrl = 'router'
        break
      case 'serv-governance':
        this.accessUrl = 'serv-governance/grey'
        break
    }
  }

  ngOnInit (): void {
    this.validateMatchForm = this.fb.group({
      position: [this.data?.position || '', [Validators.required]],
      key: [this.data?.key || '', [Validators.required, Validators.pattern('^[a-zA-Z][a-zA-Z0-9-_]*')]],
      match_type: [this.data?.match_type || '', [Validators.required]],
      pattern: [this.data?.pattern || '']
    })
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  saveMatch () {
    if (this.validateMatchForm.controls['match_type'].value === 'NULL' ||
    this.validateMatchForm.controls['match_type'].value === 'EXIST' ||
    this.validateMatchForm.controls['match_type'].value === 'UNEXIST' ||
    this.validateMatchForm.controls['match_type'].value === 'ANY') {
      this.validateMatchForm.controls['pattern'].setValue('')
    }
    if (this.validateMatchForm.valid) {
      if (!this.data) {
        if (this.matchHeaderSet.has(this.validateMatchForm.controls['key'].value)) {
          for (const index in this.matchList) {
            if (this.matchList[index].key === this.validateMatchForm.controls['key'].value && this.matchList[index].position === this.validateMatchForm.controls['position'].value) {
              this.matchList.splice(Number(index), 1)
              break
            }
          }
        }
      } else {
        for (const index in this.matchList) {
          if (this.matchList[index].key === this.editData.key && this.matchList[index].position === this.editData.position && this.matchList[index].pattern === this.editData.pattern && this.matchList[index].match_type === this.editData.match_type) {
            this.matchList.splice(Number(index), 1)
            break
          }
        }
      }
      if (this.validateMatchForm.controls['position'].value === 'HEADER') { this.matchHeaderSet.add(this.validateMatchForm.controls['key'].value) }
      this.matchList = [{ position: this.validateMatchForm.controls['position'].value, key: this.validateMatchForm.controls['key'].value, pattern: this.validateMatchForm.controls['pattern'].value, match_type: this.validateMatchForm.controls['match_type'].value }, ...this.matchList]
      this.closeDrawer()
    } else {
      Object.values(this.validateMatchForm.controls).forEach(control => {
        if (control.invalid) {
          control.markAsDirty()
          control.updateValueAndValidity({ onlySelf: true })
        }
      })
    }
  }

  closeDrawer () {
    this.eoNgCloseDrawer.emit('match')
  }
}

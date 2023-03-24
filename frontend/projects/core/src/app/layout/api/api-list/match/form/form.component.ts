/* eslint-disable dot-notation */
import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { Router } from '@angular/router'
import { SelectOption } from 'eo-ng-select'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { positionList, prefixMatchList } from '../../../types/conf'
import { MatchData } from '../../../types/types'

@Component({
  selector: 'eo-ng-match-form',
  templateUrl: './form.component.html',
  styles: [
  ]
})
export class MatchFormComponent implements OnInit {
  @Input() validateMatchForm:FormGroup = new FormGroup({})
  @Input() editData?:MatchData
  @Input()
  set matchList (val) {
    this._matchList = val
    this.matchListChange.emit(this._matchList)
  }

  get matchList () {
    return this._matchList
  }

  @Output() eoNgCloseDrawer: EventEmitter<string> = new EventEmitter()
  @Output() matchListChange : EventEmitter<MatchData[]> = new EventEmitter()

  _matchList:MatchData[] = []
  positionList:SelectOption[] =[...positionList]
  matchTypeList:SelectOption[] =[...prefixMatchList]
  matchHeaderSet:Set<string> = new Set()
  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  data:MatchData|undefined
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
      matchType: [this.data?.matchType || '', [Validators.required]],
      pattern: [this.data?.pattern || '']
    })
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  saveMatch () {
    if (this.validateMatchForm.controls['matchType'].value === 'NULL' ||
    this.validateMatchForm.controls['matchType'].value === 'EXIST' ||
    this.validateMatchForm.controls['matchType'].value === 'UNEXIST' ||
    this.validateMatchForm.controls['matchType'].value === 'ANY') {
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
          if (this.matchList[index].key === this.editData!.key && this.matchList[index].position === this.editData!.position && this.matchList[index].pattern === this.editData!.pattern && this.matchList[index].matchType === this.editData!.matchType) {
            this.matchList.splice(Number(index), 1)
            break
          }
        }
      }
      if (this.validateMatchForm.controls['position'].value === 'HEADER') { this.matchHeaderSet.add(this.validateMatchForm.controls['key'].value) }
      this.matchList = [{ position: this.validateMatchForm.controls['position'].value, key: this.validateMatchForm.controls['key'].value, pattern: this.validateMatchForm.controls['pattern'].value, matchType: this.validateMatchForm.controls['matchType'].value }, ...this.matchList]
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

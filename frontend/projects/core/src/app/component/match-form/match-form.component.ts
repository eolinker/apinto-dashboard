import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { Router } from '@angular/router'
import { SelectOption } from 'eo-ng-select'
import { defaultAutoTips } from '../../constant/conf'
import {
  Position,
  positionList,
  prefixMatchList
} from '../../layout/api/types/conf'
import { MatchData } from '../../layout/api/types/types'

@Component({
  selector: 'eo-ng-match-form',
  template: `
    <form
      nz-form
      [nzNoColon]="true"
      [nzAutoTips]="autoTips"
      [formGroup]="validateMatchForm"
      autocomplete="off"
    >
      <nz-form-item>
        <nz-form-label [nzSpan]="6" nzRequired>参数位置：</nz-form-label>
        <nz-form-control [nzSpan]="14">
          <eo-ng-select
            class="w-INPUT_NORMAL"
            name="position"
            required
            nzShowSearch
            formControlName="position"
            [nzOptions]="positionList"
            nzPlaceHolder="请选择"
            [nzDisabled]="nzDisabled"
          ></eo-ng-select>
        </nz-form-control>
      </nz-form-item>

      <nz-form-item>
        <nz-form-label [nzSpan]="6" nzRequired>参数名：</nz-form-label>
        <nz-form-control [nzSpan]="14" [nzErrorTip]="matchKeyErrorTpl">
          <input
            eo-ng-input
            required
            class="w-INPUT_NORMAL"
            name="key"
            [placeholder]="
              getKeyPlaceholder(validateMatchForm.get('position')?.value)
            "
            formControlName="key"
            [eoNgUserAccess]="accessUrl"
            (disabledEdit)="disabledEdit($event)"
          />
          <ng-template #matchKeyErrorTpl let-control>
            <ng-container *ngIf="control.hasError('pattern')">
              <ng-container
                [ngSwitch]="validateMatchForm.get('position')?.value"
              >
                <ng-container *ngSwitchCase="'header'"
                  >参数名需以字母开头，支持英文、数字、中横线、下划线组合（大小写不敏感）</ng-container
                >
                <ng-container *ngSwitchCase="'body'"
                  >支持 JSONPath
                  匹配（如：$.fieldName）或第一层级字段名</ng-container
                >
                <ng-container *ngSwitchCase="'cookie'"
                  >请输入有效的 Cookie 键名</ng-container
                >
                <ng-container *ngSwitchDefault
                  >参数名需以字母开头，支持英文、数字、中横线、下划线组合</ng-container
                >
              </ng-container>
            </ng-container>
            <ng-container *ngIf="control.hasError('required')"
              >必填项</ng-container
            >
          </ng-template>
        </nz-form-control>
      </nz-form-item>

      <nz-form-item>
        <nz-form-label [nzSpan]="6" nzRequired>参数值匹配：</nz-form-label>
        <nz-form-control [nzSpan]="14">
          <eo-ng-select
            class="mb-2 w-INPUT_NORMAL"
            name="matchType"
            required
            nzShowSearch
            formControlName="matchType"
            [nzOptions]="matchTypeList"
            nzPlaceHolder="请选择匹配类型"
            [nzDisabled]="nzDisabled"
          ></eo-ng-select>
          <input
            *ngIf="
              validateMatchForm.controls['matchType'].value !== 'NULL' &&
              validateMatchForm.controls['matchType'].value !== 'EXIST' &&
              validateMatchForm.controls['matchType'].value !== 'UNEXIST' &&
              validateMatchForm.controls['matchType'].value !== 'ANY'
            "
            eo-ng-input
            required
            class="w-INPUT_NORMAL"
            name="pattern"
            placeholder="请输入参数值"
            formControlName="pattern"
            [eoNgUserAccess]="accessUrl"
          />
        </nz-form-control>
      </nz-form-item>
    </form>
  `,
  styles: []
})
export class MatchFormComponent implements OnInit {
  @Input() validateMatchForm: FormGroup = new FormGroup({})
  @Input() editData?: MatchData
  @Input()
  set matchList(val: MatchData[]) {
    this._matchList = val
    this.matchListChange.emit(this._matchList)
  }

  get matchList(): MatchData[] {
    return this._matchList
  }

  @Output() eoNgCloseDrawer: EventEmitter<string> = new EventEmitter()
  @Output() matchListChange: EventEmitter<MatchData[]> = new EventEmitter()

  _matchList: MatchData[] = []
  positionList: SelectOption[] = [...positionList]
  matchTypeList: SelectOption[] = [...prefixMatchList]
  matchHeaderSet: Set<string> = new Set()
  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  data: MatchData = {} as MatchData
  accessUrl: string = '' // 用来判断权限的url
  nzDisabled: boolean = false

  constructor(private fb: UntypedFormBuilder, private router: Router) {}

  ngOnInit(): void {
    switch (this.router.url.split('/')[1]) {
      case 'router':
        this.accessUrl = 'router/api'
        break
      case 'serv-governance':
        this.accessUrl = 'serv-governance/grey'
        break
    }
    this.validateMatchForm = this.fb.group({
      position: [this.data?.position || '', [Validators.required]],
      key: [this.data?.key || '', [Validators.required]],
      matchType: [this.data?.matchType || '', [Validators.required]],
      pattern: [this.data?.pattern || '']
    })

    // Position-based validation for key field
    this.validateMatchForm
      .get('position')
      ?.valueChanges.subscribe((position) => {
        const keyControl = this.validateMatchForm.get('key')
        if (keyControl) {
          switch (position) {
            case Position.BODY:
              // Body: Support JSONPath and first-level field name matching
              keyControl.setValidators([Validators.required])
              break
            case Position.HEADER:
            case Position.COOKIE:
            default:
              // Query/Path: Starts with letter, allows letters, numbers, hyphen, underscore
              keyControl.setValidators([
                Validators.required,
                Validators.pattern('^[a-zA-Z][a-zA-Z0-9-_]*$')
              ])
          }
          keyControl.updateValueAndValidity()
        }
      })

    this.validateMatchForm.get('matchType')?.valueChanges.subscribe((value) => {
      if (['NULL', 'EXIST', 'UNEXIST', 'ANY'].includes(value)) {
        this.validateMatchForm.get('pattern')?.clearValidators()
      } else {
        this.validateMatchForm
          .get('pattern')
          ?.setValidators([Validators.required])
      }
      this.validateMatchForm.get('pattern')?.updateValueAndValidity()
    })
  }

  getKeyPlaceholder(position: Position): string {
    switch (position) {
      case Position.HEADER:
        return '大小写不敏感'
      case Position.BODY:
        return '支持 JSONPath 匹配和第一层级字段名匹配'
      case Position.COOKIE:
        return '请填写 Cookie 值的键值对名称，例如 userToken'
      default:
        return '支持字母开头、英文数字中横线下划线组合'
    }
  }

  disabledEdit(value: any) {
    this.nzDisabled = value
  }

  saveMatch() {
    if (this.validateMatchForm.valid) {
      if (
        ['NULL', 'EXIST', 'UNEXIST', 'ANY'].includes(
          this.validateMatchForm.controls['matchType'].value
        )
      ) {
        this.validateMatchForm.controls['pattern'].setValue('')
      }

      if (!this.data) {
        if (
          this.matchHeaderSet.has(this.validateMatchForm.controls['key'].value)
        ) {
          for (const index in this.matchList) {
            if (
              this.matchList[index].key ===
                this.validateMatchForm.controls['key'].value &&
              this.matchList[index].position ===
                this.validateMatchForm.controls['position'].value
            ) {
              this.matchList.splice(Number(index), 1)
              break
            }
          }
        }
      } else {
        for (const index in this.matchList) {
          if (
            this.matchList[index].key === this.editData!.key &&
            this.matchList[index].position === this.editData!.position &&
            this.matchList[index].pattern === this.editData!.pattern &&
            this.matchList[index].matchType === this.editData!.matchType
          ) {
            this.matchList.splice(Number(index), 1)
            if (this.matchHeaderSet.has(this.editData!.key)) {
              this.matchHeaderSet.delete(this.editData!.key)
            }
            break
          }
        }
      }
      if (this.validateMatchForm.controls['position'].value === 'HEADER') {
        this.matchHeaderSet.add(this.validateMatchForm.controls['key'].value)
      }
      this.matchList = [
        {
          position: this.validateMatchForm.controls['position'].value,
          key: this.validateMatchForm.controls['key'].value,
          pattern: this.validateMatchForm.controls['pattern'].value,
          matchType: this.validateMatchForm.controls['matchType'].value
        },
        ...this.matchList
      ]
      this.closeDrawer()
    } else {
      Object.values(this.validateMatchForm.controls).forEach((control) => {
        if (control.invalid) {
          control.markAsDirty()
          control.updateValueAndValidity({ onlySelf: true })
        }
      })
    }
  }

  closeDrawer() {
    this.eoNgCloseDrawer.emit('match')
  }
}

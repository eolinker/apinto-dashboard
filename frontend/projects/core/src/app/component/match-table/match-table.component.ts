import { Component, EventEmitter, Input, OnInit, Output, TemplateRef, ViewChild } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { EoNgFeedbackModalService } from 'eo-ng-feedback'
import { THEAD_TYPE, TBODY_TYPE } from 'eo-ng-table'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { MODAL_SMALL_SIZE } from '../../constant/app.config'
import { setFormValue } from '../../constant/form'
import { matchHeaderTableHeadName, matchHeaderTableBody } from '../../layout/api/types/conf'
import { MatchData } from '../../layout/api/types/types'
import { MatchFormComponent } from '../match-form/match-form.component'

@Component({
  selector: 'eo-ng-match-table',
  template: `
    <div>
  <button
    type="button"
    [disabled]="nzDisabled"
    eo-ng-button
    (click)="openDrawer('match')"
  >
    添加配置
  </button>
</div>
<div *ngIf="matchList.length > 0" class="mt-btnybase" style="width: 524px">
  <eo-ng-apinto-table
    [nzTbody]="matchTableBody"
    [nzThead]="matchTableHeadName"
    [nzData]="matchList"
    [nzTrClick]="matchTableClick"
    [nzMaxOperatorButton]="2"
    [nzNoScroll]="true"
  >
  </eo-ng-apinto-table>
</div>

<ng-template #positionTranslateTpl let-item="item">
  <ng-container [ngSwitch]="item.position">
    <span *ngSwitchCase="'header'">HTTP请求头</span>
    <span *ngSwitchCase="'query'">Query 参数</span>
    <span *ngSwitchCase="'body'">Body 参数</span>
    <span *ngSwitchCase="'cookie'">Cookie</span>
  </ng-container>
</ng-template>

<ng-template #matchTypeTranslateTpl let-item="item">
  <ng-container [ngSwitch]="item.matchType">
    <span *ngSwitchCase="'EQUAL'">全等匹配</span>
    <span *ngSwitchCase="'PREFIX'">前缀匹配</span>
    <span *ngSwitchCase="'SUFFIX'">后缀匹配</span>
    <span *ngSwitchCase="'SUBSTR'">子串匹配</span>
    <span *ngSwitchCase="'UNEQUAL'">非等匹配</span>
    <span *ngSwitchCase="'NULL'">空值匹配</span>
    <span *ngSwitchCase="'EXIST'">存在匹配</span>
    <span *ngSwitchCase="'UNEXIST'">不存在匹配</span>
    <span *ngSwitchCase="'REGEXP'">区分大小写的正则匹配</span>
    <span *ngSwitchCase="'REGEXPG'">不区分大小写的正则匹配</span>
    <span *ngSwitchCase="'ANY'">任意匹配</span>
  </ng-container>
</ng-template>

  `,
  styles: [
  ]
})
export class MatchTableComponent implements OnInit {
  @ViewChild('matchTypeTranslateTpl', { read: TemplateRef, static: true }) matchTypeTranslateTpl: TemplateRef<any> | undefined
  @ViewChild('positionTranslateTpl', { read: TemplateRef, static: true }) positionTranslateTpl: TemplateRef<any> | undefined
  @ViewChild('matchForm') matchForm: MatchFormComponent | undefined
  @Input() nzDisabled:boolean = false
  @Input()
  get matchList () {
    return this._matchList
  }

  set matchList (val:MatchData[]) {
    this._matchList = val
    this.matchListChange.emit(val)
  }

  @Output() matchListChange = new EventEmitter()
  _matchList:MatchData[] = []
  editData:MatchData | undefined
  validateMatchForm:FormGroup = new FormGroup({})
  modalRef:NzModalRef | undefined
  matchTableHeadName:THEAD_TYPE[] = [...matchHeaderTableHeadName]
  matchTableBody:TBODY_TYPE[]= [...matchHeaderTableBody]

  constructor (private fb: UntypedFormBuilder,
    private modalService:EoNgFeedbackModalService) {
    this.validateMatchForm = this.fb.group({
      position: ['', [Validators.required]],
      key: ['', [Validators.required, Validators.pattern('^[a-zA-Z][a-zA-Z0-9-_]*')]],
      matchType: ['', [Validators.required]],
      pattern: ['']
    })
  }

  ngOnInit (): void {
    this.matchTableBody[4].btns[0].click = (item:any) => {
      this.openDrawer('match', item.data)
    }
    this.matchTableBody[4].btns[0].disabledFn = () => {
      return this.nzDisabled
    }
    this.matchTableBody[4].btns[1].disabledFn = () => {
      return this.nzDisabled
    }
  }

  ngAfterViewInit () {
    this.matchTableBody[2].title = this.matchTypeTranslateTpl
    this.matchTableBody[0].title = this.positionTranslateTpl
  }

  matchTableClick = (item:{data:MatchData}) => {
    this.openDrawer('match', item.data)
  }

  openDrawer (type:string, data?:any) {
    switch (type) {
      case 'match': {
        if (data) {
          this.editData = data
          setFormValue(this.validateMatchForm, data)
        } else {
          this.validateMatchForm.setValue({
            position: '',
            key: '',
            matchType: '',
            pattern: ''
          }
          )
        }

        this.modalRef = this.modalService.create({
          nzTitle: '配置路由规则',
          nzWidth: MODAL_SMALL_SIZE,
          nzContent: MatchFormComponent,
          nzComponentParams: {
            data: data,
            closeDrawer: this.closeDrawer,
            matchList: this.matchList,
            editData: this.editData
          },
          nzOkDisabled: this.nzDisabled,
          nzOkText: data ? '提交' : '保存',
          nzOnOk: (component:MatchFormComponent) => {
            component.saveMatch()
            this.matchList = component.matchList
            return false
          }
        })
        break
      }
    }
  }

  closeDrawer = () => {
    this.modalRef?.close()
  }
}

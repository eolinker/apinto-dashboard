/* eslint-disable dot-notation */
import { Component, Input, OnInit, Output, TemplateRef, ViewChild, EventEmitter } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { EoNgFeedbackModalService } from 'eo-ng-feedback'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { MODAL_SMALL_SIZE } from 'projects/core/src/app/constant/app.config'
import { setFormValue } from 'projects/core/src/app/constant/form'
import { matchHeaderTableBody, matchHeaderTableHeadName } from '../../../types/conf'
import { MatchData } from '../../../types/types'
import { MatchFormComponent } from '../form/form.component'

@Component({
  selector: 'eo-ng-match-table',
  templateUrl: './table.component.html',
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
            data: data, closeDrawer: this.closeDrawer, matchList: this.matchList, editData: this.editData
          },
          nzOkDisabled: this.nzDisabled,
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

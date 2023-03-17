/* eslint-disable camelcase */
/* eslint-disable dot-notation */
import { Component, Input, OnInit, Output, TemplateRef, ViewChild, EventEmitter } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { EoNgFeedbackModalService } from 'eo-ng-feedback'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { MODAL_SMALL_SIZE } from 'projects/core/src/app/constant/app.config'
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

  set matchList (val) {
    this._matchList = val
    this.matchListChange.emit(val)
  }

@Output() matchListChange = new EventEmitter()

  _matchList:Array<{position:string, match_type:string, key:string, pattern?:string}> = []
  editData:any = null
  validateMatchForm:FormGroup = new FormGroup({})
  modalRef:NzModalRef | undefined

  matchTableHeadName:Array<object> = [
    {
      title: '参数位置'
    },
    { title: '参数名' },
    { title: '匹配类型' },
    { title: '匹配值' },
    {
      title: '操作',
      right: true
    }
  ]

  matchTableBody:Array<any> = [
    { key: 'position' },
    { key: 'key' },
    { key: 'match_type' },
    { key: 'pattern' },
    {
      type: 'btn',
      right: true,
      btns: [
        {
          title: '配置',
          click: (item:any) => {
            this.openDrawer('match', item.data)
          },
          disabledFn: () => {
            return this.nzDisabled
          }
        },
        {
          title: '删除',
          action: 'delete',
          disabledFn: () => {
            return this.nzDisabled
          }
        }
      ]
    }
  ]

  constructor (private fb: UntypedFormBuilder,
    private modalService:EoNgFeedbackModalService) {
    this.validateMatchForm = this.fb.group({
      position: ['', [Validators.required]],
      key: ['', [Validators.required, Validators.pattern('^[a-zA-Z][a-zA-Z0-9-_]*')]],
      match_type: ['', [Validators.required]],
      pattern: ['']
    })
  }

  ngOnInit (): void {
  }

  ngAfterViewInit () {
    this.matchTableBody[2].title = this.matchTypeTranslateTpl
    this.matchTableBody[0].title = this.positionTranslateTpl
  }

  matchTableClick = (item:any) => {
    this.openDrawer('match', item.data)
  }

  openDrawer (type:string, data?:any) {
    switch (type) {
      case 'match': {
        if (data) {
          this.editData = data
          this.validateMatchForm.controls['key'].setValue(data.key)
          this.validateMatchForm.controls['position'].setValue(data.position)
          this.validateMatchForm.controls['match_type'].setValue(data.match_type)
          this.validateMatchForm.controls['pattern'].setValue(data.pattern)
        } else {
          this.validateMatchForm = this.fb.group({
            position: ['', [Validators.required]],
            key: ['', [Validators.required, Validators.pattern('^[a-zA-Z][a-zA-Z0-9-_]*')]],
            match_type: ['', [Validators.required]],
            pattern: ['']
          })
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

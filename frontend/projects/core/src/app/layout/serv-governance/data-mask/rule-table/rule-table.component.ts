import { Component, EventEmitter, Input, OnInit, Output, TemplateRef, ViewChild } from '@angular/core'
import { UntypedFormBuilder } from '@angular/forms'
import { EoNgFeedbackModalService } from 'eo-ng-feedback'
import { THEAD_TYPE, TBODY_TYPE } from 'eo-ng-table'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { MODAL_SMALL_SIZE } from 'projects/apinto-auth/src/app/constant/app.config'
import { DataMaskRuleFormComponent } from '../rule-form/rule-form.component'
import { ruleHeaderTableBody, ruleHeaderTableHeadName } from '../../types/conf'
import { MaskRuleData } from '../../types/types'

@Component({
  selector: 'eo-ng-data-mask-rule-table',
  template: `
    <div>
  <button
    type="button"
    [disabled]="nzDisabled"
    eo-ng-button
    (click)="openDrawer('rule')"
  >
    添加配置
  </button>
</div>
<div *ngIf="ruleList.length > 0" class="mt-btnybase" style="width: 524px">
  <eo-ng-apinto-table
    [nzTbody]="ruleTableBody"
    [nzThead]="ruleTableHeadName"
    [nzData]="ruleList"
    [nzTrClick]="ruleTableClick"
    [nzMaxOperatorButton]="2"
    [nzNoScroll]="true"
  >
  </eo-ng-apinto-table>
</div>

<ng-template #matchTypeTranslateTpl let-item="item">
  <ng-container [ngSwitch]="item.match.type">
    <span *ngSwitchCase="'inner'">数据格式</span>
    <span *ngSwitchCase="'keyword'">关键字</span>
    <span *ngSwitchCase="'regex'">正则表达式</span>
    <span *ngSwitchCase="'json_path'">JSON Path</span>
  </ng-container>
</ng-template>

<ng-template #matchValueTranslateTpl let-item="item">
  <ng-container [ngSwitch]="item.match.value">
    <span *ngSwitchCase="'name'">姓名</span>
    <span *ngSwitchCase="'phone'">手机号</span>
    <span *ngSwitchCase="'id-card'">身份证号</span>
    <span *ngSwitchCase="'bank-card'">银行卡号</span>
    <span *ngSwitchCase="'date'">日期</span>
    <span *ngSwitchCase="'amount'">金额</span>
    <span *ngSwitchDefault>{{item.match.value}}</span>
  </ng-container>
</ng-template>


<ng-template #maskTypeTranslateTpl let-item="item">
  <ng-container [ngSwitch]="item.mask.type">
    <span *ngSwitchCase="'partial-display'">局部显示</span>
    <span *ngSwitchCase="'partial-masking'">局部遮蔽</span>
    <span *ngSwitchCase="'truncation'">截取</span>
    <span *ngSwitchCase="'replacement'">替换</span>
    <span *ngSwitchCase="'shuffling'">乱序</span>
  </ng-container>
</ng-template>


<ng-template #maskRuleTranslateTpl let-item="item">
  <ng-container [ngSwitch]="item.mask.type">
    <span
          *ngSwitchCase="'replacement'"
          eoNgFeedbackTooltip
          [nzTooltipTitle]="'类型：'+(item.mask.replace.type === 'random' ? '随机字符串' : '自定义字符串; 值：')+(item.mask.replace.value)"
          nzTooltipPlacement="top"
          [nzTooltipVisible]="false"
          nzTooltipTrigger="hover">
        类型：{{item.mask.replace.type === 'random' ? '随机字符串' : '自定义字符串; 值：'}}{{item.mask.replace.value}}
    </span>
    <span *ngSwitchCase="'shuffling'">-</span>
    <span *ngSwitchDefault
          eoNgFeedbackTooltip
          [nzTooltipTitle]="'起始位置：'+item.mask.begin + '位；长度：'+item.mask.length+'位'"
          [nzTooltipVisible]="false"
          nzTooltipTrigger="hover">
          起始位置：{{item.mask.begin}}位；长度：{{item.mask.length}}位</span>
    </ng-container>
</ng-template>

  `,
  styles: [
  ]
})
export class DataMaskRuleTableComponent implements OnInit {
  @ViewChild('matchTypeTranslateTpl', { read: TemplateRef, static: true }) matchTypeTranslateTpl: TemplateRef<any> | undefined
  @ViewChild('matchValueTranslateTpl', { read: TemplateRef, static: true }) matchValueTranslateTpl: TemplateRef<any> | undefined
  @ViewChild('maskTypeTranslateTpl', { read: TemplateRef, static: true }) maskTypeTranslateTpl: TemplateRef<any> | undefined
  @ViewChild('maskRuleTranslateTpl', { read: TemplateRef, static: true }) maskRuleTranslateTpl: TemplateRef<any> | undefined
  @Input() nzDisabled:boolean = false
  @Input()
  get ruleList () {
    return this._ruleList
  }

  set ruleList (val:MaskRuleData[]) {
    this._ruleList = val
    this.ruleListChange.emit(val)
  }

  @Output() ruleListChange = new EventEmitter()
  _ruleList:MaskRuleData[] = []
  editData:MaskRuleData | undefined
  modalRef:NzModalRef | undefined
  ruleTableHeadName:THEAD_TYPE[] = [...ruleHeaderTableHeadName]
  ruleTableBody:TBODY_TYPE[]= [...ruleHeaderTableBody]

  constructor (private fb: UntypedFormBuilder,
    private modalService:EoNgFeedbackModalService) {
  }

  ngOnInit (): void {
    this.ruleTableBody[4].btns[0].click = (item:any) => {
      this.openDrawer('rule', item.data)
    }
    this.ruleTableBody[4].btns[0].disabledFn = () => {
      return this.nzDisabled
    }
    this.ruleTableBody[4].btns[1].disabledFn = () => {
      return this.nzDisabled
    }
  }

  ngAfterViewInit () {
    this.ruleTableBody[0].title = this.matchTypeTranslateTpl
    this.ruleTableBody[1].title = this.matchValueTranslateTpl
    this.ruleTableBody[2].title = this.maskTypeTranslateTpl
    this.ruleTableBody[3].title = this.maskRuleTranslateTpl
  }

  ruleTableClick = (item:{data:MaskRuleData}) => {
    this.openDrawer('rule', item.data)
  }

  openDrawer (type:string, data?:any) {
    switch (type) {
      case 'rule': {
        if (data) {
          this.editData = data
        } else {
          this.editData = undefined
        }
        this.modalRef = this.modalService.create({
          nzTitle: '配置脱敏规则',
          nzWidth: MODAL_SMALL_SIZE,
          nzContent: DataMaskRuleFormComponent,
          nzComponentParams: {
            closeDrawer: this.closeDrawer,
            ruleList: this.ruleList,
            editData: this.editData
          },
          nzOkDisabled: this.nzDisabled,
          nzOkText: data ? '提交' : '保存',
          nzOnOk: (component:DataMaskRuleFormComponent) => {
            component.save()
            this.ruleList = component.ruleList
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

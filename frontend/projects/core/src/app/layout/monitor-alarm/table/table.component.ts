/* eslint-disable dot-notation */
import { Component, EventEmitter, Input, OnInit, Output, SimpleChanges } from '@angular/core'
import { EO_NG_DROPDOWN_MENU_ITEM } from 'eo-ng-dropdown'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { MonitorProxyTableConfig, TableQueryData, MonitorData, TableConfigEmitData } from '../types/types'

@Component({
  selector: 'eo-ng-monitor-alarm-table',
  template: `
  <div *ngIf="nzShowSearch"
        class="group-search-large flex justify-end items-center mt-btnybase mb-[6px] ml-btnbase">
        <div class="inline-block">
        <eo-ng-search-input-group  class="ml-label" [eoInputVal]="queryData.keyword" (eoClick)="queryData.keyword = '';changeDisplayData(queryData.keyword)">
          <input
              class="search"
              type="text"
              eo-ng-input
              [placeholder]="'请输入'+nzSearchKeyName+'进行搜索'"
              [(ngModel)]="queryData.keyword"
              (keyup.enter)="changeDisplayData(queryData.keyword)"
            />
        </eo-ng-search-input-group>
      </div>
  </div>
  <div class="monitor-table-block sticky t-[108px] bg-FIX_BG z-10">
    <div
      class="eo-test-table-opr-btns"
      style="
        display: flex;
        margin-top: 6px;
        position: absolute;
        z-index: 4;
        right: 20px;
      "
      dir="rtl"
      >
        <eo-ng-dropdown
          style="cursor: pointer"
          [title]="dropdownTitleTpl"
          [itemTmp]="dropdownMenTpl"
          [menus]="nzDropdownMenu"
          trigger="click"
          overlayClassName="eo-monitor-table-dropdown"
        >
        </eo-ng-dropdown>
        <ng-template #dropdownTitleTpl>
          <svg class="iconpark-icon">
            <use href="#setting"></use>
          </svg>
        </ng-template>
        <ng-template #dropdownMenTpl let-item="item">
          <label
          style="width:100%"
            (click)="$event.stopPropagation()"
            class="eo-test-app-checkbox"
            (ngModelChange)="changeMenu($event, item)"
            eo-ng-checkbox
            [(ngModel)]="item.checked"
            >{{ item.title }}</label
          >
        </ng-template>
    </div>
    <eo-ng-apinto-table
      [nzTbody]="nzTableBody"
      [nzThead]="nzTableHead"
      [nzData]="tableDisplayList"
      [nzTrClick]="trClick"
      [nzMaxOperatorButton]="2"
      [nzTotal]="queryData.total"
      [nzPageIndex]="queryData.pageNum"
      [nzPageSize]="queryData.pageSize"
      [nzFrontPagination]="nzShowPagination"
      [nzShowPagination]="nzShowPagination"
      [nzMonitorDT]="nzMonitorDT"
    >
    </eo-ng-apinto-table>
  </div>
  `,
  styles: [
  ]
})
export class MonitorAlarmTableComponent implements OnInit {
  @Input() nzSearchKeyName:'请求路径'|'API名称'|'目标节点'|'应用名称' = '请求路径'
  @Input() type:string = 'ipTop10'
  @Input() nzTableList:Array<any> = []
  @Input() nzShowPagination:boolean = false
  @Input() nzShowSearch:boolean = false
  @Input() nzMonitorDT:boolean =false
  @Input() nzTableConfig:MonitorProxyTableConfig|undefined
  @Input() nzTableBody:TBODY_TYPE[] = []
  @Input() nzTableHead:THEAD_TYPE[] = []
  @Input() nzDropdownMenu:EO_NG_DROPDOWN_MENU_ITEM[] = []
  @Output() btnClick:EventEmitter<any> = new EventEmitter()
  @Output() tableConfigChange:EventEmitter<TableConfigEmitData> = new EventEmitter()
  tableDisplayList:Array<any> = []

  queryData:TableQueryData = {
    pageNum: 1,
    pageSize: 20,
    total: 0,
    keyword: ''
  }

  ngOnInit (): void {
    this.tableConfigChange.emit({ thead: this.nzTableHead, config: this.nzTableConfig })
    this.refreshMenu()
    this.initTable()
  }

  ngOnChanges (changes:SimpleChanges) {
    if (changes['type'] && localStorage.getItem(this.type)) {
      const storageConfig = localStorage.getItem(this.type)?.split(',')
      const keyList = Object.keys(this.nzTableConfig!)
      for (const index in keyList) {
        this.nzTableConfig![keyList[index] as string] = storageConfig?.indexOf(keyList[index]) !== -1
      }
      this.refreshMenu()
    }
    if (changes['nzTableList']) {
      this.tableDisplayList = [...this.nzTableList]
      this.queryData.total = this.tableDisplayList.length
    }
  }

  initTable () {
    for (let i = 1; i < this.nzTableBody.length; i++) {
      const body = this.nzTableBody[i]
      if (i !== this.nzTableBody.length - 1) {
        body.showFn = () => {
          return this.nzTableConfig![body.key as string] === undefined ? true : this.nzTableConfig![body.key as string]
        }
      } else {
        body.btns = [
          {
            title: '查看',
            click: (item:any) => {
              this.btnClick.emit(item)
            }
          }
        ]
      }
    }

    let descendFlag:boolean = true
    for (let i = 1; i < this.nzTableHead.length - 1; i++) {
      const head = this.nzTableHead[i]
      if (head.showSort) {
        if (descendFlag) {
          head.sortOrder = 'descend'
          descendFlag = false
        } else {
          head.sortOrder = null
        }
        head.sortFn = (a: any, b: any) => a.data[head.key as string] - b.data[head.key as string]
      }
      head.showFn = () => {
        return this.nzTableConfig![head.key as string] === undefined ? true : this.nzTableConfig![head.key as string]
      }
    }
  }

  trClick = ($event:any) => {
    this.btnClick.emit($event)
  }

  changeMenu (event:boolean, item:{title:string, key:string, checked:boolean}) {
    this.nzTableConfig![item.key] = event
    this.tableConfigChange.emit({ thead: this.nzTableHead, config: this.nzTableConfig! })
    this.storeConfig()
  }

  changeDisplayData (keyword:string) {
    const keywordKey = this.nzTableBody[0].key
    if (!keyword) {
      this.tableDisplayList = [...this.nzTableList]
      this.queryData.total = this.tableDisplayList.length
      return
    }
    if (this.nzTableList?.length > 0) {
      this.tableDisplayList = this.nzTableList.filter((item:MonitorData) => (
        item[keywordKey as string].includes(keyword)
      ))
      this.queryData.total = this.tableDisplayList.length
    }
  }

  storeConfig () {
    const keyList = Object.keys(this.nzTableConfig!)
    const res:Array<string> = []
    for (const index in keyList) {
      if (this.nzTableConfig![keyList[index]]) {
        res.push(keyList[index])
      }
    }
    localStorage.setItem(this.type, res.join(','))
  }

  refreshMenu () {
    for (const menu of this.nzDropdownMenu) {
      menu['checked'] = this.nzTableConfig![menu['key'] as string]
    }
  }
}

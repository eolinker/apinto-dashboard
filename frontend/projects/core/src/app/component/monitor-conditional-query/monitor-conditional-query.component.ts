import { Component, EventEmitter, Input, OnInit, Output, TemplateRef, ViewChild } from '@angular/core'
import { TBODY_TYPE } from 'eo-ng-table'
import { rateList, valueUnitMap, compareOptions } from '../../layout/monitor-alarm/types/conf'

interface ruleListData{
  compare:string,
  unit:string,
  value:number|null
}

@Component({
  selector: 'eo-ng-monitor-conditional-query',
  templateUrl: './monitor-conditional-query.component.html',
  styles: [

    `.ant-input:hover,
    .ant-input:focus{
      height:22px;
      border:none !important;
    }

    `
  ]
})
export class MonitorConditionalQueryComponent implements OnInit {
  @ViewChild('alertValueTpl') alertValueTpl:TemplateRef<any>|undefined
  @ViewChild('andColTpl') andColTpl:TemplateRef<any>|undefined
  @Input() deleteFirst:boolean = false
  @Input() quota:'request_fail_count'|'request_fail_rate'|'request_status_4xx'|'request_status_5xx'|'proxy_fail_count'|'proxy_fail_rate'|'proxy_status_4xx'|'proxy_status_5xx'|'request_message'|'response_message'|'avg_resp'|'' = ''
  @Input()
  get rulesData () {
    return this._ruleData
  }

  set rulesData (val:ruleListData[]) {
    this._ruleData = val
    this.rulesDataChange.emit(this._ruleData)
  }

    @Output() rulesDataChange:EventEmitter<ruleListData[]> = new EventEmitter()
  _ruleData:ruleListData[] = []
  ruleDataItem:ruleListData = { compare: '', value: null, unit: '' }
  nzDisabled:boolean = false

  valueUnitMap:Map<string, string> = valueUnitMap
  rateList:Array<string> = rateList
  channel:Array<string> = []
  rulesTableBody:TBODY_TYPE[]|any[] = [
    {
      title: '且',
      seRowspan: () => {
        return this._ruleData.length
      },
      showFn: (item:any) => {
        return this._ruleData.length > 1 && item === this._ruleData[0]
      },
      styleFn: () => {
        return 'width:30px; padding-right:6px !important;vertical-align: middle;padding-bottom:0px !important;'
      }
    },
    {
      key: 'compare',
      type: 'select',
      placeholder: '比较关系',
      disabledFn: () => {
        return this.nzDisabled
      },
      checkMode: 'change',
      check: (item: any) => {
        return !item
      },
      showSearch: false,
      onSearch: () => {},
      errorTip: '必填项',
      opts: compareOptions,
      styleFn: () => {
        return 'width:196px'
      }
    },
    {
      key: 'value',
      styleFn: () => {
        return 'width:142px'
      }
    },
    {
      type: 'btn',
      showFn: (item: any) => {
        return item === this._ruleData[0] && !this.deleteFirst
      },
      btns: [
        {
          title: '添加',
          action: 'add',
          disabledFn: () => {
            return this.nzDisabled
          }
        }
      ]
    },
    {
      type: 'btn',
      showFn: (item: any) => {
        return item !== this._ruleData[0] || (item === this._ruleData[0] && this.deleteFirst)
      },
      btns: [
        {
          title: '添加',
          action: 'add',
          disabledFn: () => {
            return this.nzDisabled
          }
        },
        {
          title: '减少',
          action: 'delete',
          disabledFn: () => {
            return this.nzDisabled
          }
        }
      ]
    }]

  ngOnInit (): void {

  }

  ngAfterViewInit ():void {
    this.rulesTableBody[0].title = this.andColTpl
    this.rulesTableBody[2].title = this.alertValueTpl
  }
}

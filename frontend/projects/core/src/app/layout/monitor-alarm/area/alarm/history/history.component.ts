/* eslint-disable no-useless-constructor */
/* eslint-disable dot-notation */
import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { ActivatedRoute } from '@angular/router'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { differenceInCalendarDays } from 'date-fns'
import { monitorAlarmHistoryTableBody, monitorAlarmHistoryTableHead } from 'projects/core/src/app/constant/table.conf'
import { THEAD_TYPE } from 'eo-ng-table'
import { getTime, timeButtonOptions } from '../../../types/conf'
import { MonitorAlarmHistoryData, StrategyHistoryQueryData } from '../../../types/types'
import { RadioOption } from 'eo-ng-radio'
import { EoNgMessageService } from 'projects/core/src/app/service/eo-ng-message.service'
import { EO_TBODY_TYPE } from 'projects/eo-ng-apinto-table/src/lib/types'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'

@Component({
  selector: 'eo-ng-monitor-alarm-history',
  templateUrl: './history.component.html',
  styles: [
    ''
  ]
})
export class MonitorAlarmHistoryComponent implements OnInit {
  @ViewChild('historyContentTpl') historyContentTpl: TemplateRef<any> | undefined;
  @ViewChild('historyStatusTpl') historyStatusTpl: TemplateRef<any> | undefined;
  timeButton:string = 'hour'
  timeButtonOptions:RadioOption[] =timeButtonOptions

  datePickerValue:Array<Date> = []
  queryData:StrategyHistoryQueryData= {
    strategyName: '',
    startTime: 0,
    endTime: 0,
    partitionId: '',
    total: 0,
    pageNum: 1,
    pageSize: 20
  }

  partitionId:string = ''
  historyTableBody:EO_TBODY_TYPE[] = monitorAlarmHistoryTableBody
  historyTableHead:THEAD_TYPE[] = monitorAlarmHistoryTableHead
  historyList:MonitorAlarmHistoryData[]= []

  constructor (private baseInfo: BaseInfoService, private api:ApiService, private message: EoNgMessageService, private activateInfo:ActivatedRoute) { }

  ngOnInit (): void {
    this.partitionId = this.baseInfo.allParamsInfo.partitionId
    this.queryData.partitionId = this.partitionId
    this.getHistoryList(true)
  }

  ngAfterViewInit ():void {
    this.historyTableBody[2].title = this.historyContentTpl
    this.historyTableBody[3].title = this.historyStatusTpl
  }

  openDatePicker (open:boolean) {
    if (!open && this.datePickerValue.length > 0) {
      this.timeButton = ''
    }
  }

  disabledDate = (current: Date): boolean =>
  // Can not select days after today and today
    differenceInCalendarDays(current, new Date()) > -1;

  clearSearch () {
    this.timeButton = 'hour'
    this.datePickerValue = []
    this.queryData = {
      ...this.queryData,
      strategyName: '',
      startTime: 0,
      endTime: 0
    }
  }

  getHistoryList (init?:boolean) {
    const { startTime, endTime } = getTime(this.timeButton, this.datePickerValue, init)
    const data = { ...this.queryData, startTime: startTime, endTime: endTime }
    this.api.get('warn/history', data).subscribe((resp:{code:number, data:{history:MonitorAlarmHistoryData[], total:number}, msg:string}) => {
      if (resp.code === 0) {
        this.historyList = resp.data.history
        this.queryData.total = resp.data.total
        !init && this.message.success(resp.msg || '获取告警历史列表成功！')
      }
    })
  }
}

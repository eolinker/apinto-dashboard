/* eslint-disable dot-notation */
/* eslint-disable no-useless-constructor */
import { Component, OnInit } from '@angular/core'
import { Router } from '@angular/router'

import { differenceInCalendarDays } from 'date-fns'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { EoNgExcelService } from 'projects/core/src/app/service/eo-ng-excel.service'
import { EoNgMonitorTabsService } from 'projects/core/src/app/service/eo-ng-monitor-tabs.service'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { MonitorAppData, MonitorProxyTableConfig, QueryData, QueryDetailData } from '../../../../types/types'
import { appInitConfig, appInitDropdownMenu, appTableBody, appTableHeadName, getTime, timeButtonOptions } from '../../../../types/conf'
import { RadioOption } from 'eo-ng-radio'
import { ApplicationEnum } from 'projects/core/src/app/layout/application/types/types'
import { SelectOption } from 'eo-ng-select'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { EO_NG_DROPDOWN_MENU_ITEM } from 'eo-ng-dropdown'

@Component({
  selector: 'eo-ng-monitor-alarm-area-others-app-total',
  templateUrl: './total.component.html',
  styles: [
  ]
})
export class MonitorAlarmAreaOthersAppTotalComponent implements OnInit {
  timeButton:string = 'hour'
  timeButtonOptions:RadioOption[] = [...timeButtonOptions]

  env:string = ''
  datePickerValue:Array<Date> = []
  queryData:QueryData= {
    clusters: [],
    appIds: [],
    startTime: 0,
    endTime: 0,
    partitionId: ''
  }

  quaryDataToDetaill:QueryDetailData = { }

  appsList:MonitorAppData[] = []
  listOfClusters:SelectOption[] = []
  listOfApps:SelectOption[] = []
  tableConfig:{[key:string]:boolean} = {}
  initFlag:boolean = true
  partitionId:string = ''
  tableHead:THEAD_TYPE[] = []
  appId:string = ''

  // 表格参数
  appTableName:THEAD_TYPE[] = [...appTableHeadName]
  appTableBody:TBODY_TYPE[] = [...appTableBody]
  appTableConfig:MonitorProxyTableConfig = { ...appInitConfig }
  appTableDropdownMenu:EO_NG_DROPDOWN_MENU_ITEM[] = [...appInitDropdownMenu]

  constructor (
    private baseInfo: BaseInfoService,
    private excel: EoNgExcelService,
    private api:ApiService,
    private tabs:EoNgMonitorTabsService,
    private message: EoNgFeedbackMessageService,
    private router: Router) { }

  ngOnInit (): void {
    this.partitionId = this.baseInfo.allParamsInfo.partitionId
    this.queryData.appIds = this.baseInfo.allParamsInfo.monitorDataId ? [this.baseInfo.allParamsInfo.monitorDataId] : []
    this.appId = this.baseInfo.allParamsInfo.monitorDataId

    this.getAppTableList(true)
    this.getAppsList()
    this.queryData.clusters = this.tabs.getClusters(this.partitionId)
    this.listOfClusters = this.tabs.getClusters(this.partitionId).map((cluster) => {
      return { label: cluster, value: cluster }
    })
    this.initFlag = false
  }

  openDatePicker (open:boolean) {
    if (!open && this.datePickerValue.length > 0) {
      this.timeButton = ''
    }
  }

  // 获取选择器app列表
  getAppsList () {
    this.api.get('application/enum').subscribe((resp:{code:number, data:{applications:ApplicationEnum[]}, msg:string}) => {
      if (resp.code === 0) {
        this.listOfApps = []
        for (const item of resp.data.applications) {
          this.listOfApps = [...this.listOfApps, { label: item.name, value: item.id }]
        }
      } else {
        this.message.error(resp.msg || '获取列表数据失败!')
      }
    })
  }

  getAppTableList (init?:boolean) {
    const { startTime, endTime } = getTime(this.timeButton, this.datePickerValue, init)
    this.queryData.startTime = startTime
    this.queryData.endTime = endTime
    this.queryData.partitionId = this.partitionId
    const data = this.queryData
    this.api.post('monitor/app', data).subscribe((resp:{code:number, data:{statistics:MonitorAppData[], total:number}, msg:string}) => {
      if (resp.code === 0) {
        this.quaryDataToDetaill = {}
        if (this.timeButton) {
          this.quaryDataToDetaill.time = this.timeButton
        } else {
          this.quaryDataToDetaill.startTime = this.queryData.startTime
          this.quaryDataToDetaill.endTime = this.queryData.endTime
        }
        for (const item of resp.data.statistics) {
          item.proxyRate = Number((item.proxyRate * 100).toFixed(2))
          item.requestRate = Number((item.requestRate * 100).toFixed(2))
        }
        this.appsList = resp.data.statistics
      } else {
        this.message.error(resp.msg || '获取应用列表失败，请重试！')
      }
    })
  }

  getNewTableConfig (value:any) {
    this.tableConfig = value.config
    this.tableHead = value.thead
  }

  disabledDate = (current: Date): boolean =>
  // Can not select days before today and today
    differenceInCalendarDays(current, new Date()) > -1;

  clearSearch () {
    this.timeButton = 'hour'
    this.datePickerValue = []
    this.queryData = {
      partitionId: this.partitionId,
      clusters: this.tabs.getClusters(this.partitionId),
      appIds: this.appId ? [this.appId] : [],
      startTime: 0,
      endTime: 0
    }
  }

  export () {
    const { startTime, endTime } = getTime(this.timeButton, this.datePickerValue)
    this.queryData.startTime = startTime
    this.queryData.endTime = endTime
    this.queryData.partitionId = this.partitionId
    const data = this.queryData
    const excelTableConfig:{[key:string]:boolean} = { appName: true, ...this.tableConfig }
    this.api.post('monitor/app', data).subscribe((resp:{code:number, data:{statistics:MonitorAppData[], total:number}, msg:string}) => {
      this.excel.exportExcel('应用调用统计', [this.queryData.startTime, this.queryData.endTime], '应用调用统计', excelTableConfig, this.tableHead, resp.data.statistics)
    })
  }

  goToDetail (value:any) {
    this.router.navigate(['/', 'monitor-alarm', 'area', 'app', this.partitionId, 'detail', value.data.appId], { queryParams: { appName: value.data.appName, ...this.quaryDataToDetaill } })
  }
}

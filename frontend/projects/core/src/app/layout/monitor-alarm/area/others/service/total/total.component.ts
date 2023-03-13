/* eslint-disable dot-notation */
/* eslint-disable no-useless-constructor */
import { Component, OnInit } from '@angular/core'
import { Router, ActivatedRoute } from '@angular/router'
import { differenceInCalendarDays } from 'date-fns'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { EoNgExcelService } from 'projects/core/src/app/service/eo-ng-excel.service'

import { EoNgMonitorTabsService } from 'projects/core/src/app/service/eo-ng-monitor-tabs.service'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { Subscription } from 'rxjs/internal/Subscription'
import { MonitorProxyTableConfig, MonitorUpstreamData, QueryDetailData } from '../../../../types/types'
import { getTime, proxyBaseInitDropdownMenu, proxyInitConfig, serviceTableBody, serviceTableHeadName, timeButtonOptions } from '../../../../types/conf'
import { RadioOption } from 'eo-ng-radio'
import { EO_NG_DROPDOWN_MENU_ITEM } from 'eo-ng-dropdown'
import { THEAD_TYPE, TBODY_TYPE } from 'eo-ng-table'

@Component({
  selector: 'eo-ng-monitor-alarm-area-others-service-total',
  templateUrl: './total.component.html',
  styles: []
})
export class MonitorAlarmAreaOthersUpstreamTotalComponent implements OnInit {
  timeButton: string = 'hour'
  timeButtonOptions:RadioOption[] = [...timeButtonOptions]

  datePickerValue: Array<Date> = []
  quaryDataToDetaill: QueryDetailData = {}

  queryData: {
    partitionId: string
    clusters: Array<string>
    services: Array<string>
    startTime: number
    endTime: number
  } = {
    clusters: [],
    services: [],
    startTime: 0,
    endTime: 0,
    partitionId: ''
  }

  partitionId: string = ''
  tableConfig: { [key: string]: boolean } = {}
  tableHead: Array<any> = []
  servicesList: Array<any> = []
  listOfClusters: Array<any> = []
  listOfServices: Array<any> = []
  initFlag: boolean = true
  service: string = ''

  // 表格参数
  serviceTableName:THEAD_TYPE[] = [...serviceTableHeadName]
  serviceTableBody:TBODY_TYPE[] = [...serviceTableBody]
  serviceTableConfig:MonitorProxyTableConfig = { ...proxyInitConfig }
  serviceTableDropdownMenu:EO_NG_DROPDOWN_MENU_ITEM[] = [...proxyBaseInitDropdownMenu]

  private subscription: Subscription = new Subscription()
  constructor (
    private api: ApiService,
    private message: EoNgFeedbackMessageService,
    private excel: EoNgExcelService,
    private router: Router,
    private tabs: EoNgMonitorTabsService,
    private activateInfo: ActivatedRoute,
    private baseInfo:BaseInfoService
  ) {}

  ngOnInit (): void {
    this.partitionId = this.baseInfo.allParamsInfo.partitionId
    this.queryData.partitionId = this.partitionId
    this.queryData.services = this.baseInfo.allParamsInfo.monitorDataId ? [this.baseInfo.allParamsInfo.monitorDataId] : []
    this.service = this.baseInfo.allParamsInfo.monitorDataId

    this.getServiceList()
    this.getServiceTableList(true)

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

  getServiceTableList (init?: boolean) {
    const { startTime, endTime } = getTime(this.timeButton, this.datePickerValue, init)
    this.queryData.startTime = startTime
    this.queryData.endTime = endTime

    const data = this.queryData
    this.api
      .post('monitor/service', data)
      .subscribe(
        (resp: {
          code: number
          data: { statistics: MonitorUpstreamData[]; total: number }
          msg: string
        }) => {
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
            }
            this.servicesList = resp.data.statistics
          } else {
            this.message.error(resp.msg || '获取上游服务调用列表失败，请重试！')
          }
        }
      )
  }

  getServiceList () {
    this.api.get('service/enum').subscribe((resp: any) => {
      if (resp.code === 0) {
        this.listOfServices = []
        for (const item of resp.data.list) {
          this.listOfServices = [
            ...this.listOfServices,
            { label: item, value: item }
          ]
        }
      } else {
        this.message.error(resp.msg || '获取上游服务列表数据失败!')
      }
    })
  }

  disabledDate = (current: Date): boolean =>
    // Can not select days before today and today
    differenceInCalendarDays(current, new Date()) > -1

  export () {
    const { startTime, endTime } = getTime(this.timeButton, this.datePickerValue)
    this.queryData.startTime = startTime
    this.queryData.endTime = endTime
    this.queryData.partitionId = this.partitionId
    const data = this.queryData

    const excelTableConfig: { [key: string]: boolean } = {
      serviceName: true,
      ...this.tableConfig
    }
    this.api
      .post('monitor/service', data)
      .subscribe(
        (resp: {
          code: number
          data: { statistics: MonitorUpstreamData[]; total: number }
          msg: string
        }) => {
          this.excel.exportExcel(
            '上游调用统计',
            [this.queryData.startTime, this.queryData.endTime],
            '上游调用统计',
            excelTableConfig,
            this.tableHead,
            resp.data.statistics
          )
        }
      )
  }

  getNewTableConfig (value: any) {
    this.tableConfig = value.config
    this.tableHead = value.thead
  }

  goToDetail (value:any) {
    this.router.navigate(['/', 'monitor-alarm', 'area', 'service', this.partitionId, 'detail', value.data.serviceName], { queryParams: { ...this.quaryDataToDetaill } })
  }

  clearSearch () {
    this.timeButton = 'hour'
    this.datePickerValue = []
    this.queryData = {
      clusters: this.tabs.getClusters(this.partitionId),
      services: [],
      startTime: 0,
      endTime: 0,
      partitionId: this.partitionId
    }
  }
}

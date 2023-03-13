/* eslint-disable dot-notation */
/* eslint-disable no-useless-constructor */
/* eslint-disable camelcase */
import { Component, OnInit } from '@angular/core'
import { Router, ActivatedRoute } from '@angular/router'
import { differenceInCalendarDays } from 'date-fns'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { EoNgExcelService } from 'projects/core/src/app/service/eo-ng-excel.service'
import { MonitorUpstreamData } from '../../../total/total.component'
import { Subscription } from 'rxjs'

import { EoNgMonitorTabsService } from 'projects/core/src/app/service/eo-ng-monitor-tabs.service'

@Component({
  selector: 'eo-ng-monitor-alarm-area-others-upstream-total',
  templateUrl: './total.component.html',
  styles: []
})
export class MonitorAlarmAreaOthersUpstreamTotalComponent implements OnInit {
  timeButton: string = 'hour'
  timeButtonOptions: Array<{ label: string; value: any }> = [
    { label: '近1小时', value: 'hour' },
    { label: '近24小时', value: 'today' },
    { label: '近3天', value: 'threeDays' },
    { label: '近7天', value: 'sevenDays' }
  ]

  datePickerValue: Array<string> = []
  quaryDataToDetaill: {
    time?: string
    start_time?: number
    end_time?: number
  } = {}

  queryData: {
    partition_id: string
    clusters: Array<string>
    services: Array<string>
    start_time: number
    end_time: number
  } = {
    clusters: [],
    services: [],
    start_time: 0,
    end_time: 0,
    partition_id: ''
  }

  partitionId: string = ''
  tableConfig: { [key: string]: boolean } = {}
  tableHead: Array<any> = []
  servicesList: Array<any> = []
  listOfClusters: Array<any> = []
  listOfServices: Array<any> = []
  initFlag: boolean = true
  service: string = ''
  private subscription: Subscription = new Subscription()
  constructor (
    private api: ApiService,
    private message: EoNgFeedbackMessageService,
    private excel: EoNgExcelService,
    private router: Router,
    private tabs: EoNgMonitorTabsService,
    private activateInfo: ActivatedRoute
  ) {}

  ngOnInit (): void {
    this.subscription = this.activateInfo.queryParams.subscribe(
      (queryParams: { [x: string]: string }) => {
        if (queryParams['partition_id']) {
          this.partitionId = queryParams['partition_id']
          this.queryData.partition_id = this.partitionId
        }
        if (queryParams['service']) {
          this.queryData.services = [queryParams['service']]
          this.service = queryParams['service']
        }
      }
    )
    this.getServiceList()
    this.getServiceTableList(true)

    this.queryData.clusters = this.tabs.getClusters(this.partitionId)
    this.listOfClusters = this.tabs
      .getClusters(this.partitionId)
      .map((cluster) => {
        return { label: cluster, value: cluster }
      })
    this.initFlag = false
  }

  ngOnDestory () {
    this.subscription.unsubscribe()
  }

  openDatePicker (open:boolean) {
    if (!open && this.datePickerValue.length > 0) {
      this.timeButton = ''
    }
  }

  getServiceTableList (init?: boolean) {
    this.getTime(init)
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
              this.quaryDataToDetaill.start_time = this.queryData.start_time
              this.quaryDataToDetaill.end_time = this.queryData.end_time
            }
            for (const item of resp.data.statistics) {
              item.proxy_rate = Number((item.proxy_rate * 100).toFixed(2))
            }
            this.servicesList = resp.data.statistics
          } else {
            this.message.error(resp.msg || '获取上游服务调用列表失败，请重试！')
          }
        }
      )
  }

  getTime (init?: boolean) {
    const currentSecond = new Date().getTime() // 当前毫秒数时间戳
    let currentMin = currentSecond - (currentSecond % (60 * 1000)) // 当前分钟数时间戳
    let startMin = currentMin - 60 * 60 * 1000
    if (!init && this.timeButton) {
      switch (this.timeButton) {
        case 'hour': {
          startMin = currentMin - 60 * 60 * 1000
          break
        }
        case 'today': {
          startMin = currentMin - 24 * 60 * 60 * 1000
          break
        }
        case 'threeDays': {
          startMin =
            new Date(new Date().setHours(0, 0, 0, 0)).getTime() -
            2 * 24 * 60 * 60 * 1000
          break
        }
        case 'sevenDays': {
          startMin =
            new Date(new Date().setHours(0, 0, 0, 0)).getTime() -
            6 * 24 * 60 * 60 * 1000
          break
        }
      }
    } else if (this.datePickerValue.length === 2) {
      startMin = new Date(
        new Date(this.datePickerValue[0]).setHours(0, 0, 0, 0)
      ).getTime()
    currentMin = new Date(new Date(this.datePickerValue[1]).setHours(23, 59, 59, 0)).getTime()
  }

    this.queryData.start_time = startMin / 1000
    this.queryData.end_time = currentMin / 1000
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
    this.getTime()
    this.queryData.partition_id = this.partitionId
    const data = this.queryData

    const excelTableConfig: { [key: string]: boolean } = {
      service_name: true,
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
            [this.queryData.start_time, this.queryData.end_time],
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

  goToDetail (value: any) {
    this.router.navigate(['/', 'monitor-alarm', 'area', 'upstream', 'detail'], {
      queryParams: {
        partition_id: this.partitionId,
        service: value.data.service_name,
        ...this.quaryDataToDetaill
      }
    })
  }

  clearSearch () {
    this.timeButton = 'hour'
    this.datePickerValue = []
    this.queryData = {
      clusters: this.tabs.getClusters(this.partitionId),
      services: [],
      start_time: 0,
      end_time: 0,
      partition_id: this.partitionId
    }
  }
}

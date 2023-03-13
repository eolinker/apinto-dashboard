/* eslint-disable dot-notation */
/* eslint-disable no-useless-constructor */
import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { ActivatedRoute, Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { MonitorLineGraphComponent } from 'projects/core/src/app/layout/monitor-alarm/area/graph/line/line.component'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { EoNgMonitorTabsService } from 'projects/core/src/app/service/eo-ng-monitor-tabs.service'
import { differenceInCalendarDays } from 'date-fns'
import { getTime, getTimeUnit, initPie, timeButtonOptions, changeNumberUnit, totalTabsList, serviceTableHeadName, apiInitConfig, apiInitDropdownMenu, apiTableBody, apiTableHeadName, appInitConfig, appInitDropdownMenu, appTableBody, appTableHeadName, proxyBaseInitDropdownMenu, proxyInitConfig, serviceTableBody } from '../../types/conf'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { RadioOption } from 'eo-ng-radio'
import { BaseQueryData, InvokeData, MessageData, MonitorApiData, MonitorAppData, MonitorProxyTableConfig, MonitorUpstreamData, QueryDetailData, SummaryPieData, TotalQueryData } from '../../types/types'
import { MonitorPieGraphComponent } from '../graph/pie/pie.component'
import { TabsOptions } from 'eo-ng-tabs'
import { EO_NG_DROPDOWN_MENU_ITEM } from 'eo-ng-dropdown'
import { THEAD_TYPE, TBODY_TYPE } from 'eo-ng-table'

@Component({
  selector: 'eo-ng-monitor-alarm-area-total',
  templateUrl: './total.component.html',
  styles: [
    `
    `
  ]
})
export class MonitorAlarmAreaTotalComponent implements OnInit {
  @ViewChild('appTableTpl') appTableTpl: TemplateRef<any> | undefined;
  @ViewChild('upstreamTableTpl') upstreamTableTpl: TemplateRef<any> | undefined;
  @ViewChild('ipTableTpl') ipTableTpl: TemplateRef<any> | undefined;
  @ViewChild('apiTableTpl') apiTableTpl: TemplateRef<any> | undefined;
  @ViewChild('requestPieRef') requestPieRef: MonitorPieGraphComponent | undefined;
  @ViewChild('proxyPieRef') proxyPieRef: MonitorPieGraphComponent | undefined;
  @ViewChild('invokeLineRef') invokeLineRef: MonitorLineGraphComponent | undefined;
  @ViewChild('trafficLineRef') trafficLineRef: MonitorLineGraphComponent | undefined;
  optionsList:Array<TabsOptions>=[...totalTabsList]

  selectedIndex:number = 0
  tableList:{api:Array<any>, app:Array<any>, service:Array<any>}={
    api: [],
    app: [],
    service: []
  }

  nzDisabled:boolean = false
  partitionId:string = ''
  timeButton:string = 'hour'
  env:string = ''
  requestStaticError:boolean = true
  proxyStaticError:boolean = true
  invokeStaticError:boolean = true
  trafficStaticError:boolean = true
  listOfClusters:Array<{label:string, value:any}>=[]
  timeButtonOptions: RadioOption[] = [...timeButtonOptions]

  queryData:BaseQueryData={
    clusters: [],
    date: []
  }

  timeUnit:string = ''
  timeInterval:string = ''
  time:string = ''

  datePickerValue:Array<Date> = []

  requestStatic:SummaryPieData= { ...initPie }
  proxyStatic:SummaryPieData= { ...initPie }

  requestPie:{[key:string]:number} = {}
  proxyPie:{[key:string]:number} = {}
  proxySucRate:string = '0%'
  requestSucRate:string = '0%'
  invokeStatic:InvokeData={ date: [], requestRate: [], requestTotal: [], proxyRate: [], proxyTotal: [], status4xx: [], status5xx: [] }
  trafficStatic:MessageData = { date: [], request: [], response: [] }
  quaryDataToDetaill:QueryDetailData = { }

  tableNameList:Array<string> = ['apiTop10', 'appTop10', 'upstreamTop10']
  tableTypeList:Array<'api'| 'app'|'service'> = ['api', 'app', 'service']

  // 表格参数
  apiTableName:THEAD_TYPE[] = [...apiTableHeadName]
  apiTableBody:TBODY_TYPE[] = [...apiTableBody]
  apiTableConfig:MonitorProxyTableConfig = { ...apiInitConfig }
  apiTableDropdownMenu:EO_NG_DROPDOWN_MENU_ITEM[] = [...apiInitDropdownMenu]

  appTableName:THEAD_TYPE[] = [...appTableHeadName]
  appTableBody:TBODY_TYPE[] = [...appTableBody]
  appTableConfig:MonitorProxyTableConfig = { ...appInitConfig }
  appTableDropdownMenu:EO_NG_DROPDOWN_MENU_ITEM[] = [...appInitDropdownMenu]

  serviceTableName:THEAD_TYPE[] = [...serviceTableHeadName]
  serviceTableBody:TBODY_TYPE[] = [...serviceTableBody]
  serviceTableConfig:MonitorProxyTableConfig = { ...proxyInitConfig }
  serviceTableDropdownMenu:EO_NG_DROPDOWN_MENU_ITEM[] = [...proxyBaseInitDropdownMenu]

  constructor (private baseInfo:BaseInfoService, private api:ApiService, private message: EoNgFeedbackMessageService, private tabs:EoNgMonitorTabsService, private router:Router, private activateInfo:ActivatedRoute) {
  }

  ngOnInit (): void { // 获取params下的id
    this.partitionId = this.baseInfo.allParamsInfo.partitionId
    this.listOfClusters = this.tabs.getClusters(this.partitionId).map((cluster) => {
      return { label: cluster, value: cluster }
    })
    this.queryData['clusters'] = this.tabs.getClusters(this.partitionId)
    if (this.partitionId) { this.getMonitorData(true) }
  }

  ngAfterViewInit () {
    this.optionsList[0].content = this.apiTableTpl
    this.optionsList[1].content = this.appTableTpl
    this.optionsList[2].content = this.upstreamTableTpl
  }

  changeNumberUnit = changeNumberUnit

  openDatePicker (open:boolean) {
    if (!open && this.datePickerValue.length > 0) {
      this.timeButton = ''
    }
  }

  getMonitorData (init?:boolean) {
    const { startTime, endTime } = getTime(this.timeButton, this.datePickerValue, init)
    const data:TotalQueryData = { uuid: this.partitionId, clusters: this.queryData['clusters'], start: startTime, end: endTime }
    this.queryData.date = [startTime, endTime]
    this.getPieData(data)
    this.getInvokeData(data)
    this.getMessageData(data)
    this.getTablesData(this.tableTypeList[this.selectedIndex])
  }

  getPieData (data:TotalQueryData) {
    this.api.post('monitor/overview/summary', data).subscribe((resp:{code:number, data:{requestSummary:SummaryPieData, proxySummary:SummaryPieData}, msg:string}) => {
      if (resp.code === 0) {
        this.quaryDataToDetaill = {}
        if (this.timeButton) {
          this.quaryDataToDetaill.time = this.timeButton
        } else {
          this.quaryDataToDetaill.startTime = this.queryData['date'][0] as number
          this.quaryDataToDetaill.endTime = this.queryData['date'][1] as number
        }
        this.requestStatic = resp.data.requestSummary
        this.proxyStatic = resp.data.proxySummary
        this.requestPie = { 请求成功数: resp.data.requestSummary.success, 请求失败数: resp.data.requestSummary.fail }
        this.proxyPie = { 转发成功数: resp.data.requestSummary.success, 转发失败数: resp.data.proxySummary.fail }
        this.proxyStaticError = false
        this.requestStaticError = false
        this.requestPieRef?.changePieChart()
        this.proxyPieRef?.changePieChart()
        this.requestSucRate = resp.data.requestSummary.total === 0 ? '0%' : (resp.data.requestSummary.success * 100 / resp.data.requestSummary.total).toFixed(2) + '%'
        this.proxySucRate = resp.data.proxySummary.total === 0 ? '0%' : (resp.data.proxySummary.success * 100 / resp.data.proxySummary.total).toFixed(2) + '%'
      } else {
        this.proxyStaticError = true
        this.requestStaticError = true
        this.message.error(resp.msg || '获取数据失败，请重试！')
      }
    })
  }

  getInvokeData (data:TotalQueryData) {
    this.api.post('monitor/overview/invoke', data).subscribe((resp:{code:number, data:InvokeData, msg:string}) => {
      if (resp.code === 0) {
        const { timeInterval, ...arr } = resp.data
        this.invokeStatic = arr as InvokeData
        this.invokeStaticError = false
        this.timeUnit = getTimeUnit(timeInterval!)
        this.invokeLineRef?.changeLineChart()
      } else {
        this.invokeStaticError = true
        this.message.error(resp.msg || '获取调用量统计数据失败，请重试！')
      }
    })
  }

  getMessageData (data:TotalQueryData) {
    this.api.post('monitor/overview/message', data).subscribe((resp:{code:number, data:MessageData, msg:string}) => {
      if (resp.code === 0) {
        this.trafficStaticError = false
        this.trafficStatic = resp.data
        this.trafficLineRef?.changeLineChart()
      } else {
        this.trafficStaticError = true
        this.message.error(resp.msg || '获取报文量统计数据失败，请重试！')
      }
    })
  }

  getTablesData (type:'app'|'api'|'service') {
    const data:TotalQueryData = { uuid: this.partitionId, clusters: this.queryData['clusters'], start: this.queryData['date'][0] as number, end: this.queryData['date'][1] as number }
    this.api.post('monitor/overview/top', { ...data, data_type: type }).subscribe((resp:{code:number, data:{api:MonitorApiData[], app:MonitorAppData[], service:MonitorUpstreamData[]}, msg:string}) => {
      if (resp.code === 0) {
        switch (type) {
          case 'api': {
            for (const api of resp.data.api) {
              api.proxyRate = Number((api.proxyRate * 100).toFixed(2))
              api.requestRate = Number((api.requestRate * 100).toFixed(2))
            }
            break
          }
          case 'app': {
            for (const app of resp.data.app) {
              app.proxyRate = Number((app.proxyRate * 100).toFixed(2))
              app.requestRate = Number((app.requestRate * 100).toFixed(2))
            }
            break
          }
          case 'service': {
            for (const service of resp.data.service) {
              service.proxyRate = Number((service.proxyRate * 100).toFixed(2))
            }
            break
          }
        }
        this.tableList[type] = resp.data[type]
      } else {
        this.message.error(resp.msg || '获取分区数据失败，请重试！')
      }
    })
  }

  changeMenu (target: string, event: any, item: any) {
    const vm: any = this
    vm[target].listShowRef[item.key] = event
  }

  goToDetail (val:any, type:string) {
    const queryData:{apiId?:string, appId?:string, apiName?:string, appName?:string, service?:string, ip?:string, } = {
    }
    let detailId:string = ''
    switch (type) {
      case 'api':
        detailId = val.data.apiId
        queryData['apiName'] = val.data.apiName
        break
      case 'app':
        detailId = val.data.appId
        queryData['appName'] = val.data.appName
        break
      case 'service':
        detailId = val.data.serviceName
        break
    }
    this.router.navigate(['/', 'monitor-alarm', 'area', type, this.partitionId, 'detail', detailId], { queryParams: { ...queryData, ...this.quaryDataToDetaill } })
  }

  resetQuery () {
    this.timeButton = 'hour'
    this.datePickerValue = []
    this.queryData = {
      clusters: this.tabs.getClusters(this.partitionId),
      date: []
    }
  }

  disabledDate = (current: Date): boolean =>
  // Can not select days before today and today
    differenceInCalendarDays(current, new Date()) > -1;
}

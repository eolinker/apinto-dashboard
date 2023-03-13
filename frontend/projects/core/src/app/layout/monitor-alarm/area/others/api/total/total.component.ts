/* eslint-disable dot-notation */
/* eslint-disable no-useless-constructor */
import { Component, OnInit } from '@angular/core'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { differenceInCalendarDays } from 'date-fns'
import { EoNgExcelService } from 'projects/core/src/app/service/eo-ng-excel.service'
import { ActivatedRoute, Router } from '@angular/router'
import { EoNgMonitorTabsService } from 'projects/core/src/app/service/eo-ng-monitor-tabs.service'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { MonitorApiData, MonitorProxyTableConfig, QueryData, QueryDetailData } from '../../../../types/types'
import { apiInitConfig, apiInitDropdownMenu, apiTableBody, apiTableHeadName, getTime, timeButtonOptions } from '../../../../types/conf'
import { RadioOption } from 'eo-ng-radio'
import { SelectOption } from 'eo-ng-select'
import { RouterEnum } from 'projects/core/src/app/layout/api/types/types'
import { EO_NG_DROPDOWN_MENU_ITEM } from 'eo-ng-dropdown'
import { THEAD_TYPE, TBODY_TYPE } from 'eo-ng-table'

@Component({
  selector: 'eo-ng-monitor-alarm-area-api-total',
  templateUrl: './total.component.html',
  styles: []
})
export class MonitorAlarmAreaOthersApiTotalComponent implements OnInit {
  timeButton:string = 'hour'
  timeButtonOptions: RadioOption[] =[...timeButtonOptions]

  tableConfig:{[key:string]:boolean} = {}
  tableHead:Array<any> = []
  env:string = ''
  datePickerValue:Array<Date> = []
  queryData:QueryData= {
    clusters: [],
    services: [],
    apiIds: [],
    path: '',
    startTime: 0,
    endTime: 0,
    partitionId: ''
  }

  quaryDataToDetaill:QueryDetailData = { }

  listOfClusters:SelectOption[] = []
  listOfServices:SelectOption[] = []
  listOfApis:SelectOption[] = [] // 选择器用
  apisList:Array<any> = [] // 表格用
  partitionId:string = ''
  initFlag:boolean = true
  apiId:string = ''

  // 表格参数
  apiTableName:THEAD_TYPE[] = [...apiTableHeadName]
  apiTableBody:TBODY_TYPE[] = [...apiTableBody]
  apiTableConfig:MonitorProxyTableConfig = { ...apiInitConfig }
  apiTableDropdownMenu:EO_NG_DROPDOWN_MENU_ITEM[] = [...apiInitDropdownMenu]

  constructor (private baseInfo:BaseInfoService, private api:ApiService, private message: EoNgFeedbackMessageService, private tabs:EoNgMonitorTabsService, private excel: EoNgExcelService, private router:Router, private activateInfo:ActivatedRoute) { }

  ngOnInit (): void {
    this.partitionId = this.baseInfo.allParamsInfo.partitionId
    this.queryData.apiIds = this.baseInfo.allParamsInfo.monitorDataId ? [this.baseInfo.allParamsInfo.monitorDataId] : []
    this.apiId = this.baseInfo.allParamsInfo.monitorDataId
    this.getServiceList()
    this.getApisList()
    this.getApiTableList(true)
    this.queryData.clusters = this.tabs.getClusters(this.partitionId)
    this.listOfClusters = this.tabs.getClusters(this.partitionId).map((cluster) => {
      return { label: cluster, value: cluster }
    })
    this.initFlag = false
  }

  ngOnDestroy () {
  }

  openDatePicker (open:boolean) {
    if (!open && this.datePickerValue.length > 0) {
      this.timeButton = ''
    }
  }

  // 获取选择器api列表
  getApisList (services?:Array<string>) {
    this.api.get('router/enum', { service_names: services ? services.join(',') : '' }).subscribe((resp:{code:number, data:{apis:RouterEnum[]}, msg:string}) => {
      if (resp.code === 0) {
        this.listOfApis = []
        for (const item of resp.data.apis) {
          this.listOfApis = [...this.listOfApis, { label: item.name, value: item.apiId }]
        }
      } else {
        this.message.error(resp.msg || '获取列表数据失败!')
      }
    })
  }

  serviceChanged (open:boolean) {
    if (!open) {
      this.getApisList(this.queryData.services)
    }
  }

  getServiceList () {
    this.api.get('service/enum').subscribe((resp:{code:number, data:{list:Array<string>}, msg:string}) => {
      if (resp.code === 0) {
        this.listOfServices = []
        for (const item of resp.data.list) {
          this.listOfServices = [...this.listOfServices, { label: item, value: item }]
        }
      } else {
        this.message.error(resp.msg || '获取上游服务列表数据失败!')
      }
    })
  }

  getApiTableList (init?:boolean) {
    const { startTime, endTime } = getTime(this.timeButton, this.datePickerValue, init)
    this.queryData.startTime = startTime
    this.queryData.endTime = endTime
    this.queryData.partitionId = this.partitionId
    const data = this.queryData
    this.api.post('monitor/api', data).subscribe((resp:{code:number, data:{statistics:MonitorApiData[], total:number}, msg:string}) => {
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
        this.apisList = [...resp.data.statistics]
      } else {
        this.message.error(resp.msg || '获取API调用数据失败，请重试！')
      }
    })
  }

  getNewTableConfig (value:any) {
    this.tableConfig = value.config
    this.tableHead = value.thead
  }

  goToDetail (value:any) {
    this.router.navigate(['/', 'monitor-alarm', 'area', 'api', this.partitionId, 'detail', value.data.apiId], { queryParams: { apiName: value.data.apiName, ...this.quaryDataToDetaill } })
  }

  export () {
    const { startTime, endTime } = getTime(this.timeButton, this.datePickerValue)
    this.queryData.startTime = startTime
    this.queryData.endTime = endTime
    this.queryData.partitionId = this.partitionId
    const data = this.queryData
    const excelTableConfig:{[key:string]:boolean} = { apiName: true, ...this.tableConfig }
    this.api.post('monitor/api', data).subscribe((resp:{code:number, data:{statistics:MonitorApiData[], total:number}, msg:string}) => {
      this.excel.exportExcel('API调用统计', [this.queryData.startTime, this.queryData.endTime], 'API调用统计', excelTableConfig, this.tableHead, resp.data.statistics)
    })
  }

  clearSearch () {
    this.timeButton = 'hour'
    this.datePickerValue = []
    this.queryData = {
      partitionId: this.partitionId,
      clusters: this.tabs.getClusters(this.partitionId),
      services: [],
      apiIds: this.apiId ? [this.apiId] : [],
      startTime: 0,
      endTime: 0,
      path: ''
    }
  }

disabledDate = (current: Date): boolean =>
  // Can not select days before today and today
  differenceInCalendarDays(current, new Date()) > -1;
}

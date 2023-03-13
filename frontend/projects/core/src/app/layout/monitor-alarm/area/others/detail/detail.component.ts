/* eslint-disable dot-notation */
/* eslint-disable no-useless-constructor */
import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { ActivatedRoute, Router } from '@angular/router'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { EoNgExcelService } from 'projects/core/src/app/service/eo-ng-excel.service'
import { MonitorApiData, MonitorApiProxyData, MonitorAppData, MonitorData, MonitorIPData, MonitorNodeData, MonitorPathData, MonitorProxyData, MonitorUpstreamData, InvokeData, QueryData, MonitorProxyTableConfig } from '../../../types/types'
import { differenceInCalendarDays } from 'date-fns'
import { Subscription } from 'rxjs'
import { EoNgMonitorTabsService } from 'projects/core/src/app/service/eo-ng-monitor-tabs.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'
import { MODAL_NORMAL_SIZE } from 'projects/core/src/app/constant/app.config'
import { apiInitConfig, apiInitDropdownMenu, apiProxyTableBody, apiProxyTableHeadName, apiTableBody, apiTableHeadName, appInitConfig, appInitDropdownMenu, appTableBody, appTableHeadName, getTime, getTimeUnit, nodeTableBody, nodeTableHeadName, proxyBaseInitDropdownMenu, proxyInitConfig, timeButtonOptions } from '../../../types/conf'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { RadioOption } from 'eo-ng-radio'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { EO_NG_DROPDOWN_MENU_ITEM } from 'eo-ng-dropdown'
import { TabsOptions } from 'eo-ng-tabs'

interface EoTabOption extends TabsOptions{
  [k:string]:any
}

@Component({
  selector: 'eo-ng-monitor-alarm-area-others-detail',
  templateUrl: './detail.component.html',
  styles: [
    `
    .line-block{
      border:1px solid #e8e8e8;
      margin:0px 24px 6px 12px;
    }
    .table-block{
      border:1px solid #e8e8e8;
      margin:20px 24px 20px 12px;
    }`
  ]
})
export class MonitorAlarmAreaOthersDetailComponent implements OnInit {
  @ViewChild('ipTableTpl') ipTableTpl: TemplateRef<any> | undefined;
  @ViewChild('appTableTpl') appTableTpl: TemplateRef<any> | undefined;
  @ViewChild('apiTableTpl') apiTableTpl: TemplateRef<any> | undefined;
  @ViewChild('apiProxyTableTpl') apiProxyTableTpl: TemplateRef<any> | undefined;
  @ViewChild('nodeTableTpl') nodeTableTpl: TemplateRef<any> | undefined;
  @ViewChild('proxyTableTpl') proxyTableTpl: TemplateRef<any> | undefined;
  @ViewChild('pathTableTpl') pathTableTpl: TemplateRef<any> | undefined;
  @ViewChild('modalTpl') modalTpl: TemplateRef<any> | undefined;

  timeButton:string = 'hour'
  timeButtonOptions:RadioOption[] = timeButtonOptions
  partitionId:string = ''
  tableConfig:{[key:string]:boolean } = {}
  tableHead:Array<any> = []
  tableList: {app:MonitorAppData[], addr: MonitorNodeData[], ip:MonitorIPData[], api:MonitorApiData[]|MonitorApiProxyData[], service:MonitorUpstreamData[], path:MonitorPathData[], proxyPath:MonitorProxyData[], [key:string]:any}= {
    app: [],
    addr: [],
    ip: [],
    api: [],
    service: [],
    path: [],
    proxyPath: []

  }

  compareTotal:boolean = false
  selectedIndex:number = 0
  type:string = ''
  tabsList:EoTabOption[]=[]
  datePickerValue:Array<Date> = []
  lineGraphType:'invoke'|'invoke-service'|'traffic'|'' = 'invoke'
  queryData:QueryData= {
    partitionId: this.partitionId,
    clusters: [],
    startTime: 0,
    endTime: 0
  }

  queryDataFromTotal:{[key:string]:any} = {}
  dataName:string = ''
  originType:string = ''
  invokeStatic:InvokeData={ date: [], requestRate: [], requestTotal: [], proxyRate: [], proxyTotal: [], status4xx: [], status5xx: [] }
  detailInvokeStatic:InvokeData={ date: [], requestRate: [], requestTotal: [], proxyRate: [], proxyTotal: [], status4xx: [], status5xx: [] }
  listOfClusters:Array<any> = []
  timeUnit:string = ''
  tableNameList:Array<string> = ['appDetail', 'upstreamDetail', 'ipDetail']
  breadcrumbName:string = ''
  modalTitle:string = ''
  // 表格参数
  apiTableName:THEAD_TYPE[] = [...apiTableHeadName]
  apiTableBody:TBODY_TYPE[] = [...apiTableBody]
  apiTableConfig:MonitorProxyTableConfig = { ...apiInitConfig }
  apiTableDropdownMenu:EO_NG_DROPDOWN_MENU_ITEM[] = [...apiInitDropdownMenu]

  appTableName:THEAD_TYPE[] = [...appTableHeadName]
  appTableBody:TBODY_TYPE[] = [...appTableBody]
  appTableConfig:MonitorProxyTableConfig = { ...appInitConfig }
  appTableDropdownMenu:EO_NG_DROPDOWN_MENU_ITEM[] = [...appInitDropdownMenu]

  nodeTableName:THEAD_TYPE[] = [...nodeTableHeadName]
  nodeTableBody:TBODY_TYPE[] = [...nodeTableBody]
  nodeTableConfig:MonitorProxyTableConfig = { ...proxyInitConfig }
  nodeTableDropdownMenu:EO_NG_DROPDOWN_MENU_ITEM[] = [...proxyBaseInitDropdownMenu]

  apiProxyTableName:THEAD_TYPE[] = [...apiProxyTableHeadName]
  apiProxyTableBody:TBODY_TYPE[] = [...apiProxyTableBody]
  apiProxyTableConfig:MonitorProxyTableConfig = { ...proxyInitConfig }
  apiProxyTableDropdownMenu:EO_NG_DROPDOWN_MENU_ITEM[] = [...proxyBaseInitDropdownMenu]

  private subscription: Subscription = new Subscription()
  private subscription2: Subscription = new Subscription()
  constructor (private baseInfo:BaseInfoService, private api:ApiService, private appConfigService:AppConfigService, private message: EoNgFeedbackMessageService, private excel: EoNgExcelService, private router:Router, private activateInfo:ActivatedRoute, private modal: EoNgFeedbackModalService, private tabs:EoNgMonitorTabsService) { }

  ngOnInit (): void {
    this.subscription2 = this.router.events.subscribe(() => {
      this.type = this.router.url.split('?')[0].split('/')[3] === 'service' ? 'service' : this.router.url.split('?')[0].split('/')[3]
      this.getTabsList()
    })
    this.type = this.router.url.split('?')[0].split('/')[3] === 'service' ? 'service' : this.router.url.split('?')[0].split('/')[3]

    this.subscription = this.activateInfo.queryParams.subscribe((queryParams: { [x: string]: string }) => {
      this.partitionId = this.baseInfo.allParamsInfo.partitionId
      this.listOfClusters = this.tabs.getClusters(this.partitionId).map((cluster) => {
        return { label: cluster, value: cluster }
      })
      this.queryData.partitionId = this.partitionId
      this.queryDataFromTotal = {}
      if (queryParams['time']) { this.queryDataFromTotal['time'] = queryParams['time'] }
      if (queryParams['startTime']) { this.queryDataFromTotal['startTime'] = queryParams['startTime'] }
      if (queryParams['endTime']) { this.queryDataFromTotal['endTime'] = queryParams['endTime'] }
      switch (this.type) {
        case 'api': {
          this.queryData.apiId = this.baseInfo.allParamsInfo.monitorDataId
          this.dataName = queryParams['apiName'] || 'API'
          this.breadcrumbName = `${queryParams['apiName'] || 'API'}调用详情`
          this.appConfigService.reqFlashBreadcrumb([{ title: '监控告警', routerLink: 'monitor-alarm' }, { title: this.tabs.getTabName(this.partitionId), routerLink: 'monitor-alarm/area/total/' + this.partitionId }, { title: 'API调用统计', routerLink: 'monitor-alarm/area/api/' + this.partitionId }, { title: this.breadcrumbName }])
          break
        }
        case 'app': {
          this.queryData.appId = this.baseInfo.allParamsInfo.monitorDataId
          this.dataName = queryParams['appName'] || '应用'
          this.breadcrumbName = `${queryParams['appName'] || '应用'}调用详情`
          this.appConfigService.reqFlashBreadcrumb([{ title: '监控告警', routerLink: 'monitor-alarm' }, { title: this.tabs.getTabName(this.partitionId), routerLink: 'monitor-alarm/area/total/' + this.partitionId }, { title: '应用调用统计', routerLink: 'monitor-alarm/area/app/' + this.partitionId }, { title: this.breadcrumbName }])
          break
        }
        case 'service': {
          this.queryData.serviceName = this.baseInfo.allParamsInfo.monitorDataId
          this.dataName = this.baseInfo.allParamsInfo.monitorDataId || '上游'
          this.lineGraphType = 'invoke-service'
          this.breadcrumbName = `${this.baseInfo.allParamsInfo.monitorDataId}调用详情`
          this.appConfigService.reqFlashBreadcrumb([{ title: '监控告警', routerLink: 'monitor-alarm' }, { title: this.tabs.getTabName(this.partitionId), routerLink: 'monitor-alarm/area/total/' + this.partitionId }, { title: '上游调用统计', routerLink: 'monitor-alarm/area/service/' + this.partitionId }, { title: this.breadcrumbName }])
          break
        }
      }
    })

    this.queryData.clusters = this.tabs.getClusters(this.partitionId)
  }

  ngAfterViewInit () {
    this.getTabsList()
    this.getMonitorData(true)
  }

  ngOnDestroy () {
    this.subscription.unsubscribe()
    this.subscription2.unsubscribe()
  }

  getNewTableConfig (value:any) {
    this.tableConfig = value.config
    this.tableHead = value.thead
  }

  // 根据类型获取tabsList,tableNameList,tabsKeyList
  getTabsList () {
    if (this.originType !== this.type) {
      switch (this.type) {
        case 'api': {
          this.tabsList = [
            { title: '应用调用', type: 'app' }
          ]
          break
        }
        case 'app': {
          this.tabsList = [
            { title: 'API调用', type: 'api' }
          ]
          break
        }
        case 'service': {
          this.tabsList = [
            { title: 'API转发', content: this.apiProxyTableTpl, type: 'api', lazyLoad: true },
            { title: '目标节点', content: this.nodeTableTpl, type: 'addr', lazyLoad: true }]
          break
        }
      }
      this.originType = this.type
    }
  }

  openDatePicker (open:boolean) {
    if (!open && this.datePickerValue.length > 0) {
      this.timeButton = ''
    }
  }

  getMonitorData (init?:boolean) {
    if (init && this.queryDataFromTotal['time']) {
      this.timeButton = this.queryDataFromTotal['time']

      const { startTime, endTime } = getTime(this.timeButton, this.datePickerValue)
      this.queryData.startTime = startTime
      this.queryData.endTime = endTime
    } else if (init && this.queryDataFromTotal['startTime']) {
      this.queryData.startTime = this.queryDataFromTotal['startTime']
      this.queryData.endTime = this.queryDataFromTotal['endTime']
      this.datePickerValue[0] = new Date(this.queryData.startTime * 1000)
      this.datePickerValue[1] = new Date(this.queryData.endTime * 1000)
      this.timeButton = ''

      const { startTime, endTime } = getTime(this.timeButton, this.datePickerValue)
      this.queryData.startTime = startTime
      this.queryData.endTime = endTime
    } else {
      const { startTime, endTime } = getTime(this.timeButton, this.datePickerValue, init)
      this.queryData.startTime = startTime
      this.queryData.endTime = endTime
    }
    const data:QueryData = { ...this.queryData, partitionId: this.partitionId, clusters: this.queryData.clusters, startTime: this.queryData.startTime, endTime: this.queryData.endTime }
    this.getInvokeData(data)
    if (this.type !== 'service' || !init) {
      this.getTablesData(data, this.tabsList[this.selectedIndex]['type'])
    }
  }

  getInvokeData (data:QueryData) {
    this.api.post(`monitor/${this.type}/details`, data).subscribe((resp:{code:number, data:{tendency:InvokeData, timeInterval:string}, msg:string}) => {
      if (resp.code === 0) {
        this.invokeStatic = resp.data.tendency
        this.timeUnit = getTimeUnit(resp.data.timeInterval)
      } else {
        this.message.error(resp.msg || '获取调用量统计数据失败，请重试！')
      }
    })
  }

  // 当this.type为api,app,ip时，查看调用数据，返回值为api表/app表/ip表/path表
  // 当this.type为service时，查看转发数据，返回值为apiProxy, addr, proxy
  getTablesData (data:QueryData, type:string) {
    this.api.post(`monitor/${this.type}/details/${type}`, data).subscribe((resp:{code:number, data:{statistics?:MonitorData[]|MonitorApiData[]}, msg:string}) => {
      if (resp.code === 0) {
        for (const item of resp.data.statistics!) {
          item.proxyRate = Number((item.proxyRate * 100).toFixed(2))
          item.requestRate = Number((item.requestRate * 100).toFixed(2))
        }
        this.tableList[type] = resp.data.statistics
      } else {
        this.message.error(resp.msg || '获取分区数据失败，请重试！')
      }
    })
  }

  // 根据当前url判断所在区域的参数，如在查看api调用详情总览页面时，清除应用、ip等数据
  getCurrentParam () {
    const newQueryDataTmp:QueryData = {
      partitionId: this.queryData.partitionId,
      clusters: this.queryData.clusters,
      startTime: this.queryData.startTime,
      endTime: this.queryData.endTime
    }
    switch (this.type) {
      case 'api': {
        newQueryDataTmp.apiId = this.queryData.apiId
        break
      }
      case 'app': {
        newQueryDataTmp.appId = this.queryData.appId
        break
      }
      case 'service': {
        newQueryDataTmp.serviceName = this.queryData.serviceName
        break
      }
    }
    this.queryData = newQueryDataTmp
  }

  goToDetail (value:any, type:string) {
    let tabName:string = '' // 被选中的数据名，用于弹窗标题
    let key:string = '' // 被选中的数据id，用于调取接口
    switch (type) {
      case 'api': {
        tabName = value.data.apiName
        key = value.data.apiId
        break
      }
      case 'app': {
        tabName = value.data.appName
        key = value.data.appId
        break
      }
      case 'service': {
        tabName = value.data.serviceName
        key = value.data.serviceName
        break
      }
      case 'addr': {
        tabName = value.data.addr
        key = value.data.addr
        break
      }
      case 'path': {
        tabName = value.data.path
        key = value.data.path
        break
      }
      case 'proxyPath': {
        tabName = value.data.proxyPath
        key = value.data.proxyPath
        break
      }
    }
    this.getDetailStatic(type, key)
    this.modalTitle = `${tabName}-${this.dataName}调用趋势`
    this.modal.create({
      nzTitle: this.modalTitle,
      nzContent: this.modalTpl,
      nzWidth: MODAL_NORMAL_SIZE,
      nzClosable: true,
      nzClassName: 'monitor-modal',
      nzFooter: null,
      nzOnCancel: () => {
        this.getCurrentParam()
        this.compareTotal = false
      }
    })
  }

  // 获取弹窗内折线图数据
  // dataType:api,app,service,ip
  // detailType:api, app, path, proxyPath,addr, ip
  getDetailStatic (detailType:string, key:string) {
    switch (detailType) {
      case 'api': {
        this.queryData.apiId = key
        break
      }
      case 'app': {
        this.queryData.appId = key
        break
      }
      case 'path': {
        this.queryData.path = key
        break
      }
      case 'proxyPath': {
        this.queryData.proxyPath = key
        break
      }
      case 'addr': {
        this.queryData.addr = key
        break
      }
    }
    this.api.post(`monitor/${this.type}/details/${detailType}/trend`, this.queryData).subscribe((resp:{code:number, data:{tendency:InvokeData}, msg:string}) => {
      if (resp.code === 0) {
        this.detailInvokeStatic = resp.data.tendency
      } else {
        this.message.error(resp.msg || '获取调用统计量失败，请重试')
      }
    })
  }

  clearSearch () {
    this.timeButton = 'hour'
    this.datePickerValue = []
    this.queryData.clusters = this.tabs.getClusters(this.partitionId)

    this.queryData.startTime = 0
    this.queryData.endTime = 0
  }

  tableChange (value:number) {
    this.getTablesData(this.queryData, this.tabsList[value]['type'])
  }

disabledDate = (current: Date): boolean =>
  // Can not select days before today and today
  differenceInCalendarDays(current, new Date()) > -1;
}

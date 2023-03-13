/* eslint-disable dot-notation */
/* eslint-disable no-useless-constructor */
import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { Router, ActivatedRoute } from '@angular/router'
import { EO_NG_DROPDOWN_MENU_ITEM } from 'eo-ng-dropdown'
import { EoNgFeedbackModalService } from 'eo-ng-feedback'
import { SelectOption } from 'eo-ng-select'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { MODAL_SMALL_SIZE } from 'projects/core/src/app/constant/app.config'
import { monitorAlarmStrategyTableBody, monitorAlarmStrategyTableConfig, monitorAlarmStrategyTableDropdownMenu, monitorAlarmStrategyTableHead } from 'projects/core/src/app/constant/table.conf'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { EoNgMessageService } from 'projects/core/src/app/service/eo-ng-message.service'
import { EoNgMonitorNoticeService } from 'projects/core/src/app/service/eo-ng-monitor-notice.service'
import { isEnableOptions, warnDimensionOptions } from '../../../../types/conf'
import { MonitorAlarmChannelsData, MonitorAlarmStrategyData, MonitorAlarmStrategyListData, StrategyQueryData } from '../../../../types/types'
import { MonitorAlarmStrategyAlertComponent } from '../alert/alert.component'

@Component({
  selector: 'eo-ng-monitor-alarm-strategy-list',
  templateUrl: './list.component.html',
  styles: [
  ]
})
export class MonitorAlarmStrategyListComponent implements OnInit {
  @ViewChild('switchTpl') switchTpl: TemplateRef<any> | undefined;
  queryData:StrategyQueryData = {
    strategyName: '',
    warnDimension: [],
    partitionId: '',
    status: -1,
    total: 0,
    pageNum: 1,
    pageSize: 20
  }

  partitionId:string = ''
  listOfMetrics:SelectOption[]= warnDimensionOptions
  listOfStatus:SelectOption[] = isEnableOptions

  strategyTableBody:TBODY_TYPE[] = [...monitorAlarmStrategyTableBody]
  strategyTableHead:THEAD_TYPE[] = [...monitorAlarmStrategyTableHead]
  strategyList:MonitorAlarmStrategyListData[] = []
  strategyTableDropdownMenu:EO_NG_DROPDOWN_MENU_ITEM[]= [...monitorAlarmStrategyTableDropdownMenu]
  strategyTableConfig:{[key:string]:boolean} = { ...monitorAlarmStrategyTableConfig }
  nzDisabled:boolean = false
  channelsModalRef:NzModalRef|undefined
  constructor (
    private baseInfo:BaseInfoService,
    private api:ApiService,
    private message: EoNgMessageService,
    private router: Router,
    private modalService:EoNgFeedbackModalService,
    private noticeService:EoNgMonitorNoticeService) { }

  ngOnInit (): void {
    this.partitionId = this.baseInfo.allParamsInfo.partitionId
    this.queryData.partitionId = this.partitionId

    for (const tableBody of this.strategyTableBody) {
      tableBody.showFn = () => {
        return this.strategyTableConfig[tableBody.key as string]
      }
    }

    for (const tableThead of this.strategyTableHead) {
      tableThead.showFn = () => {
        return this.strategyTableConfig[tableThead.key as string]
      }
    }

    this.strategyTableHead.push({ title: '操作', right: true })
    this.strategyTableBody.push({
      type: 'btn',
      right: true,
      btns: [{
        title: '查看',
        click: (item:any) => {
          this.viewStrategy(item)
        }
      }, {
        title: '删除',
        disabledFn: () => { return this.nzDisabled },
        click: (item:any) => {
          this.deleteModal(item.data)
        }
      }]
    })
    this.getStrategyList(true)
  }

  ngAfterViewInit () {
    this.strategyTableBody[5].title = this.switchTpl
  }

  disabledEdit (editAccess:boolean) {
    this.nzDisabled = editAccess
  }

  getStrategyList (init?:boolean) {
    const data:StrategyQueryData = { ...this.queryData }
    if (!init) {
      data.status = data.status === ('' || null) ? -1 : data.status
      if (typeof data.warnDimension === 'object') {
        data.warnDimension = data.warnDimension.join(',')
      }
    }
    this.api.get('warn/strategys', data).subscribe((resp:{code:number, data:{datas:MonitorAlarmStrategyListData[], total:number}, msg:string}) => {
      if (resp.code === 0) {
        this.strategyList = resp.data.datas
        this.queryData.total = resp.data.total
        !init && this.message.success(resp.msg || '获取告警策略列表成功！')
      }
    })
  }

  addNewStrategy () {
    this.noticeService.getLastedNoticeChannels().subscribe((channelsList:MonitorAlarmChannelsData[]) => {
      if (channelsList.length > 0) {
        this.router.navigate(['/', 'monitor-alarm', 'area', 'strategy', this.partitionId, 'config'])
      } else {
        this.channelsModalRef = this.modalService.create({
          nzWrapClassName: 'color-blue',
          nzTitle: '配置告警渠道',
          nzIconType: 'exclamation-circle',
          nzContent: MonitorAlarmStrategyAlertComponent,
          nzComponentParams: {
            closeModal: this.closeModal
          },
          nzClosable: true,
          nzFooter: null
        })
      }
    })
  }

  viewStrategy = (tableItem:{data:MonitorAlarmStrategyListData, [key:string]:any}) => {
    this.router.navigate(['/', 'monitor-alarm', 'area', 'strategy', this.partitionId, 'message', tableItem.data.uuid])
  }

  deleteModal (strategy:MonitorAlarmStrategyListData, e?:Event) {
    e?.stopPropagation()
    this.modalService.create({
      nzTitle: '删除',
      nzContent: `${strategy.strategyTitle}告警策略一旦删除，将不再监控告警。`,
      nzClosable: true,
      nzWidth: MODAL_SMALL_SIZE,
      nzClassName: 'delete-modal',
      nzOkDanger: true,
      nzOnOk: () => {
        this.deleteStrategy(strategy.uuid)
      }
    })
  }

  deleteStrategy (id:string) {
    this.api.delete('warn/strategy', { uuid: id }).subscribe(resp => {
      if (resp.code === 0) {
        this.getStrategyList()
        this.message.success(resp.msg || '删除成功')
      }
    })
  }

  clearSearch () {
    this.queryData = {
      ...this.queryData,
      warnDimension: [],
      status: -1,
      strategyName: ''
    }
  }

  changeMenu (checked:boolean, tableConfig:any) {
    tableConfig.checked = checked
    this.strategyTableConfig[tableConfig.key] = checked
  }

  disabledStrategy (strategy:MonitorAlarmStrategyData, $event:Event) {
    $event.stopPropagation()
    strategy['loading'] = true
    this.api.patch('warn/strategy', { uuid: strategy.uuid, isEnable: !strategy.isEnable }).subscribe((resp:{code:number, data:{}, msg:string}) => {
      strategy['loading'] = false
      if (resp.code === 0) {
        strategy.isEnable = !strategy.isEnable
        this.message.success(resp.msg || '操作成功！')
      }
    })
  }

  closeModal = () => {
    this.channelsModalRef?.close()
  }
}

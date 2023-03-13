import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core'
import { SelectOption } from 'eo-ng-select'
import { ApiService } from 'projects/core/src/app/service/api.service'
import {
  MonitorAlarmChannelsData,
  MonitorAlarmStrategyRuleConditionData,
  MonitorAlarmStrategyRuleData
} from '../../../../types/types'

@Component({
  selector: 'eo-ng-monitor-alarm-strategy-rule',
  templateUrl: './rule.component.html',
  styles: []
})
export class MonitorAlarmStrategyRuleComponent implements OnInit {
  @Input() editPage: boolean = false
  @Input()
  get rulesList () {
    return this._rulesList
  }

  set rulesList (rulesList: MonitorAlarmStrategyRuleData[]) {
    this._rulesList = rulesList
    this.rulesListChange.emit(rulesList)
  }

  @Output() rulesListChange: EventEmitter<MonitorAlarmStrategyRuleData[]> =
    new EventEmitter()

  @Input() listOfChannels: SelectOption[] = []
  @Input() quota:
    | 'request_fail_count'
    | 'request_fail_rate'
    | 'request_status_4xx'
    | 'request_status_5xx'
    | 'proxy_fail_count'
    | 'proxy_fail_rate'
    | 'proxy_status_4xx'
    | 'proxy_status_5xx'
    | 'request_message'
    | 'response_message'
    | 'avg_resp'
    | '' = ''

  _rulesList: MonitorAlarmStrategyRuleData[] = []

  constructor (private api: ApiService) {}

  ngOnInit (): void {
    this.getChannelList()
  }

  getChannelList () {
    this.api
      .get('warn/channels')
      .subscribe(
        (resp: {
          code: number
          data: { channels: MonitorAlarmChannelsData[] }
          msg: string
        }) => {
          if (resp.code === 0) {
            this.listOfChannels = resp.data.channels.map(
              (channelData: MonitorAlarmChannelsData) => {
                if (channelData.type === 2) {
                  return { label: '邮箱通知渠道 - 邮箱', value: channelData.uuid }
                }
                return {
                  label: channelData.title + '通知渠道 - Webhook',
                  value: channelData.uuid
                }
              }
            )
          }
        }
      )
  }

  // 当query内rule数组长度为0时，需要删除当前条件组
  checkRuleEmpty (
    valueArr: MonitorAlarmStrategyRuleConditionData[],
    index: number
  ) {
    if (valueArr.length === 0) {
      this.deleteQueryGroup(index)
    }
    this.rulesListChange.emit(this.rulesList)
  }

  // 告警规则-删除当前条件组
  deleteQueryGroup (index: number) {
    this.rulesList.splice(index, 1)
  }

  // 告警规则-新增条件组
  addQueryGroup () {
    this.rulesList.push({
      channelUuids: [],
      condition: [
        {
          compare: '',
          unit: '',
          value: null
        }
      ]
    })
  }
}

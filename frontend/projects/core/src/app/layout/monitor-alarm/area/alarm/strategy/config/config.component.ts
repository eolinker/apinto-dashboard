/* eslint-disable dot-notation */
import { Component, Input, OnInit } from '@angular/core'
import { AbstractControl, FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { Router } from '@angular/router'
import { EoNgFeedbackModalService } from 'eo-ng-feedback'
import { RadioOption } from 'eo-ng-radio'
import { SelectOption } from 'eo-ng-select'
import { MODAL_LARGE_SIZE } from 'projects/core/src/app/constant/app.config'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { EoNgMyValidators } from 'projects/core/src/app/constant/eo-ng-validator'
import { setFormValue } from 'projects/core/src/app/constant/form'
import { EmptyHttpResponse, UserListData } from 'projects/core/src/app/constant/type'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { EoNgMessageService } from 'projects/core/src/app/service/eo-ng-message.service'
import { EoNgMonitorTabsService } from 'projects/core/src/app/service/eo-ng-monitor-tabs.service'
import { everyOptions, quoteOptions, rateList, targetTypeOptions, valueUnit, valueUnitMap, warnDimensionOptions } from '../../../../types/conf'
import { MonitorAlarmStrategyData, MonitorAlarmStrategyRuleData, MonitorAlarmStrategyTargetValueData } from '../../../../types/types'
import { MonitorAlarmStrategyTransferComponent } from '../transfer/transfer.component'

@Component({
  selector: 'eo-ng-monitor-alarm-strategy-config',
  templateUrl: './config.component.html',
  styles: [
  ]
})
export class MonitorAlarmStrategyConfigComponent implements OnInit {
  @Input() editPage:boolean = false
  @Input() uuid:string = ''
  partitionId:string = ''
  validateForm: FormGroup = new FormGroup({})
  nzDisabled:boolean = false

  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  listOfDimensions:RadioOption[]=[...warnDimensionOptions as RadioOption[]]
  listOfTargetRules:RadioOption[] = [...targetTypeOptions as RadioOption[]]
  listOfQuotas:SelectOption[] = [...quoteOptions]
  listOfEvery:SelectOption[] = [...everyOptions]
  listOfClusters:SelectOption[] = []
  listOfUsers:SelectOption[] = []
  listOfSelectedValue:Array<string> = []
  targetValueMap:{api:MonitorAlarmStrategyTargetValueData,
    service:MonitorAlarmStrategyTargetValueData} = {
      api: {
        contain: [],
        not_contain: []
      },
      service: {
        contain: [],
        not_contain: []
      }
    }

  rulesList:MonitorAlarmStrategyRuleData[] = [
    {
      channelUuids: [],
      condition: [
        {
          compare: '',
          unit: '',
          value: null
        }
      ]
    }
  ]

  rulesFilteredList:MonitorAlarmStrategyRuleData[] = []
  rateList:Array<string> = rateList
  valueUnitMap:Map<string, string> = valueUnitMap
  showClusterValueError:boolean = false
  listOfChannels:SelectOption[] = []
  rulesErrorTip:boolean = false
  valueUnit:{[key:string]:string} = valueUnit

  constructor (
    public api:ApiService,
    private fb: UntypedFormBuilder,
    private tabs:EoNgMonitorTabsService,
    private modalService:EoNgFeedbackModalService,
    private message: EoNgMessageService,
    private router:Router,
    private baseInfo:BaseInfoService,
    private appConfigService:AppConfigService) {
    this.partitionId = this.baseInfo.allParamsInfo.partitionId
    this.uuid = this.baseInfo.allParamsInfo.strategyUuid || ''
    this.listOfClusters = this.tabs.getClusters(this.partitionId).map((clusterName:string) => {
      return { label: clusterName, value: clusterName }
    })

    this.validateForm = this.fb.group({
      title: ['', [Validators.required, EoNgMyValidators.maxByteLength(32), Validators.pattern('^[\u4E00-\u9FA5A-Za-z]+$')]],
      desc: [''],
      isEnable: [true, [Validators.required]],
      dimension: ['api', [Validators.required]],
      target: this.fb.group({
        rule: ['unlimited'],
        values: [null, [this.ruleValidator]]
      }),
      quota: ['', [Validators.required]],
      every: ['', [Validators.required]],
      continuity: [1],
      hourMax: [1],
      users: [[]]
    })
  }

  ngOnInit (): void {
    if (this.editPage) {
      this.appConfigService.reqFlashBreadcrumb([{ title: '监控告警', routerLink: 'monitor-alarm' }, { title: this.tabs.getTabName(this.partitionId), routerLink: 'monitor-alarm/area/total/' + this.partitionId }, { title: '告警策略', routerLink: 'monitor-alarm/area/strategy/' + this.partitionId }, { title: '编辑告警策略' }])
      this.getStrategyMessage()
    } else {
      this.appConfigService.reqFlashBreadcrumb([{ title: '监控告警', routerLink: 'monitor-alarm' }, { title: this.tabs.getTabName(this.partitionId), routerLink: 'monitor-alarm/area/total/' + this.partitionId }, { title: '告警策略', routerLink: 'monitor-alarm/area/strategy/' + this.partitionId }, { title: '新建告警策略' }])
    }
    this.getUsersList()
  }

  getStrategyMessage () {
    this.api.get('warn/strategy', { uuid: this.uuid }).subscribe((resp:{code:number, data:{strategy:MonitorAlarmStrategyData}, msg:string}) => {
      if (resp.code === 0) {
        setFormValue(this.validateForm, resp.data.strategy)
        if ((resp.data.strategy.dimension === 'api' || resp.data.strategy.dimension === 'service') && resp.data.strategy.target.rule !== 'unlimited') {
          // eslint-disable-next-line camelcase
          const temp:{contain?:Array<string>, not_contain?:Array<string>} = {}
          temp[resp.data.strategy.target.rule as 'contain' || 'not_contain'] = resp.data.strategy.target.values
          this.targetValueMap[resp.data.strategy.dimension] = { contain: temp.contain || [], not_contain: temp.not_contain || [] }
        }
        this.rulesList = resp.data.strategy.rule
      }
    })
  }

  getUsersList () {
    this.api.get('user/enum').subscribe((resp:{code:number, data:{users:UserListData[]}}) => {
      if (resp.code === 0) {
        this.listOfUsers = resp.data.users.map((user:UserListData) => {
          return { label: `${user.nickName}(${user.id})`, value: user.id }
        })
      }
    })
  }

  dimensionValueChange (dimension:string) {
    if (dimension === 'api' || dimension === 'service') {
      this.validateForm.get('target.rule')?.setValue('unlimited')
      this.validateForm.get('target.values')?.setValue([])
    } else {
      this.validateForm.get('target.rule')?.setValue('')
      this.validateForm.get('target.values')?.setValue([])
    }
  }

  // 当告警维度选择API或上游，且告警目标选择包含或不包含时，
  // selectedRule不能为空，否则页面需有提示
  checkRadioError () {
    return (this.validateForm.controls['dimension'].value === 'api' ||
               this.validateForm.controls['dimension'].value === 'service') &&
          (this.validateForm.get('target.rule')?.value !== 'unlimited') &&
          this.targetValueMap[this.validateForm.controls['dimension'].value as 'api' || 'service'][this.validateForm.get('target.rule')!.value as 'contain' || 'not_contain']?.length === 0
  }

  // 校验告警规则是否有至少一组填写完成，如有则过滤出该组数据用于保存
  rulesRequired ():boolean {
    this.rulesFilteredList = []
    for (const rule of this.rulesList) {
      if (rule.channelUuids.length !== 0) {
        const conditionFilterList:Array<{compare:string, unit:string, value:number|null}> = []
        for (const cond of rule.condition) {
          if (cond.compare && cond.value) {
            cond.unit = this.valueUnit[(rateList.indexOf(cond.compare) === -1 ? this.valueUnitMap.get(this.validateForm.controls['quota'].value)! : '%')]
            cond.value = cond.unit === 'num' ? Math.floor(cond.value) : (cond.unit === '%' ? Number(cond.value.toFixed(2)) : cond.value)
            conditionFilterList.push({ ...cond })
          }
        }
        if (conditionFilterList.length > 0) {
          this.rulesFilteredList.push({ channelUuids: rule.channelUuids, condition: conditionFilterList })
        }
      }
    }
    this.rulesErrorTip = this.rulesFilteredList.length === 0
    return this.rulesFilteredList.length > 0
  }

  saveStrategy () {
    this.rulesRequired()
    this.showClusterValueError = true
    if (this.validateForm.valid && this.rulesRequired()) {
      const data:MonitorAlarmStrategyData = { ...this.validateForm.value, rule: [...this.rulesFilteredList] }
      // 当维度选择api或上游且告警目标不为不限时，target.value 的值为去除非完整组的rulesFilteredList
      if ((this.validateForm.controls['dimension'].value === 'api' || this.validateForm.controls['dimension'].value === 'service') && this.validateForm.get('target.rule')!.value !== 'unlimited') { data.target.values = this.targetValueMap[this.validateForm.controls['dimension'].value === 'api' ? 'api' : 'service'][this.validateForm.get('target.rule')!.value] }
      // 当维度选择分区时，target送{}
      if (this.validateForm.controls['dimension'].value === 'partition') { data.target = {} }
      // 当维度选择集群时，rule送contain
      if (this.validateForm.controls['dimension'].value === 'cluster') { data.target.rule = 'contain' }
      // 当告警目标为不限时，需删除target.value
      if (this.validateForm.get('target.rule')!.value === 'unlimited') { delete data.target.values }

      data.partitionId = this.partitionId
      if (this.editPage) {
        data.uuid = this.uuid
        this.api.put('warn/strategy', data).subscribe((resp:EmptyHttpResponse) => {
          if (resp.code === 0) {
            this.message.success(resp.msg || '修改告警策略成功！')
            this.backToList()
          }
        })
      } else {
        this.api.post('warn/strategy', data).subscribe((resp:EmptyHttpResponse) => {
          if (resp.code === 0) {
            this.message.success(resp.msg || '创建告警策略成功！')
            this.backToList()
          }
        })
      }
    } else {
      Object.values(this.validateForm.controls).forEach(control => {
        if (control.invalid) {
          control.markAsDirty()
          control.updateValueAndValidity({ onlySelf: true })
        }
      })
    }
  }

  backToList () {
    this.router.navigate(['/', 'monitor-alarm', 'area', 'strategy', this.partitionId])
  }

  // 告警目标，选择对应的api或service
  openTargetValueModal () {
    this.modalService.create({
      nzTitle: `选择${this.validateForm.controls['dimension'].value === 'api' ? 'API' : '上游服务'}数据`,
      nzContent: MonitorAlarmStrategyTransferComponent,
      nzClosable: true,
      nzWidth: MODAL_LARGE_SIZE,
      nzClassName: 'transfer-modal',
      nzCancelText: '取消',
      nzComponentParams: {
        type: this.validateForm.controls['dimension'].value,
        rule: this.validateForm.get('target.rule')!.value as 'contain' || 'not-contain',
        selectedList: this.targetValueMap[this.validateForm.controls['dimension'].value === 'api' ? 'api' : 'service'][this.validateForm.get('target.rule')!.value]
      },
      nzOkText: this.targetValueMap[this.validateForm.controls['dimension'].value === 'api' ? 'api' : 'service'][this.validateForm.get('target.rule')!.value].length > 0 ? '提交' : '保存',
      nzOnOk: () => {
        return this.targetValueMap[this.validateForm.controls['dimension'].value === 'api' ? 'api' : 'service'][this.validateForm.get('target.rule')!.value].length > 0
      }
    })
  }

  // 当告警维度不为分区时，告警目标为必填
  ruleValidator = (control: AbstractControl<any, any>): { [s: string]: boolean } => {
    if (this.validateForm.controls['dimension']?.value !== 'cluster') {
      return {}
    }
    if (!control.value || control.value.length === 0) {
      return { error: true, required: true }
    }
    return {}
  }
}

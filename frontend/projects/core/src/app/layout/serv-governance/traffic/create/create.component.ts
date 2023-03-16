/* eslint-disable camelcase */
/* eslint-disable no-array-constructor */
/* eslint-disable no-prototype-builtins */
/* eslint-disable dot-notation */
/* eslint-disable no-useless-constructor */
import {
  Component,
  Input,
  OnInit,
  Output,
  EventEmitter,
  TemplateRef,
  ViewChild
} from '@angular/core'
import { ActivatedRoute, Router } from '@angular/router'
import {
  EoNgFeedbackMessageService
} from 'eo-ng-feedback'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'
import {
  FormGroup,
  UntypedFormBuilder,
  Validators
} from '@angular/forms'
import { FilterShowData } from '../../filter/footer/footer.component'
import { EoNgMyValidators } from 'projects/core/src/app/constant/eo-ng-validator'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'

interface TrafficStrategyData {
  name: string
  uuid?: string
  desc?: string
  priority?: number | null
  filters: Array<{
    name: string
    values: Array<string>
    type?: string
    label?: string
    title?: string
    [key: string]: any
  }>
  config: {
    metrics: Array<any>
    query: { second: number; minute: number; hour: number }
    traffic: { second: number; minute: number; hour: number }
    response: {
      status_code: number
      content_type: string
      charset: string
      header: Array<{ key: string; value: string }>
      body: string
    }
  }
  [key: string]: any
}

@Component({
  selector: 'eo-ng-traffic-create',
  templateUrl: './create.component.html',
  styles: ['']
})
export class TrafficCreateComponent implements OnInit {
  @ViewChild('checkbox', { read: TemplateRef, static: true }) checkbox:
    | TemplateRef<any>
    | undefined

  @Input() editPage: boolean = false
  @Input() clusterName: string = ''
  @Input() strategyUuid: string = ''
  @Input() fromList: boolean = false
  @Output() changeToList: EventEmitter<any> = new EventEmitter()

  filterNamesSet: Set<string> = new Set() // 用户已选择的筛选条件放入set中,在显示筛选条件的选择器里需要过去set中存在的值

  validatePriority: boolean = true

  filterShowList: FilterShowData[] = [] // 展示在前端页面的筛选条件表格,包含uuid和对应选项名称,实际提交时只需要uuid

  autoTips: Record<string, Record<string, string>> = defaultAutoTips

  metricsList: Array<{ label: string; value: string; disable?: boolean }> = []

  createStrategyForm: TrafficStrategyData = {
    name: '',
    desc: '',
    priority: null,
    filters: [],
    config: {
      metrics: [],
      query: { second: 0, minute: 0, hour: 0 },
      traffic: { second: 0, minute: 0, hour: 0 },
      response: {
        status_code: 500,
        content_type: 'application/json',
        charset: 'UTF-8',
        header: [{ key: '', value: '' }],
        body: ''
      }
    }
  }

  validateForm: FormGroup = new FormGroup({})
  responseForm: FormGroup = new FormGroup({})
  showMetricsError: boolean = false
  nzDisabled: boolean = false

  responseHeaderList: Array<{
    key: string
    value: string
    [key: string]: any
  }> = [{ key: '', value: '' }]

  constructor (
    private baseInfo:BaseInfoService,
    private message: EoNgFeedbackMessageService,
    private api: ApiService,
    private activateInfo: ActivatedRoute,
    private fb: UntypedFormBuilder,
    private router:Router,
    private appConfigService: AppConfigService
  ) {
    this.validateForm = this.fb.group({
      name: [
        '',
        [Validators.required, Validators.pattern('^[a-zA-Z][a-zA-Z0-9_]*')]
      ],
      desc: [],
      priority: [null, [EoNgMyValidators.priority]],
      limitQuerySecond: [0, [Validators.required]],
      limitQueryMinute: [0, [Validators.required]],
      limitQueryHour: [0, [Validators.required]],
      limitTrafficSecond: [0, [Validators.required]],
      limitTrafficMinute: [0, [Validators.required]],
      limitTrafficHour: [0, [Validators.required]]
    })

    this.responseForm = this.fb.group({
      status_code: [200, [Validators.required, Validators.pattern(/^[1-9]{1}\d{2}$/)]],
      content_type: ['application/json', [Validators.required]],
      charset: ['UTF-8', [Validators.required]],
      header: [],
      body: ['']
    })

    this.appConfigService.reqFlashBreadcrumb([
      { title: '流量策略', routerLink: 'serv-governance/traffic' },
      { title: '新建流量策略' }
    ])
  }

  ngOnInit (): void {
    this.clusterName = this.baseInfo.allParamsInfo.clusterName
    if (this.editPage) {
      this.getStrategyMessage()
    }
    this.getMetricsList()
  }

  ngOnChanges (): void {
    if (this.strategyUuid) {
      this.createStrategyForm.uuid = this.strategyUuid
    }
  }

  // 当页面是编辑策略页时,需要根据集群名和策略uuid获取策略信息
  getStrategyMessage () {
    this.api
      .get('strategy/traffic', { uuid: this.createStrategyForm.uuid || '' })
      .subscribe(
        (resp: {
          code: number
          data: { strategy?: TrafficStrategyData; [key: string]: any }
          msg: string
        }) => {
          if (resp.code === 0) {
            this.appConfigService.reqFlashBreadcrumb([
              {
                title: '流量策略',
                routerLink: 'serv-governance/traffic'
              },
              { title: resp.data.strategy!.name }
            ])
            this.validateForm.controls['name'].setValue(resp.data.strategy!.name)
            this.createStrategyForm = resp.data.strategy!
            this.createStrategyForm.uuid = resp.data.strategy!.uuid
            this.createStrategyForm.filters = resp.data.strategy!.filters || []
            this.createStrategyForm.config.metrics =
              resp.data.strategy!.config.metrics || []

            this.validateForm.controls['desc'].setValue(resp.data.strategy!.desc)
            this.validateForm.controls['priority'].setValue(
              resp.data.strategy!.priority
            )
            this.validateForm.controls['limitQuerySecond'].setValue(
              resp.data.strategy!.config.query?.second || 0
            )
            this.validateForm.controls['limitQueryMinute'].setValue(
              resp.data.strategy!.config.query?.minute || 0
            )
            this.validateForm.controls['limitQueryHour'].setValue(
              resp.data.strategy!.config.query?.hour || 0
            )

            this.validateForm.controls['limitTrafficSecond'].setValue(
              resp.data.strategy!.config.traffic?.second || 0
            )
            this.validateForm.controls['limitTrafficMinute'].setValue(
              resp.data.strategy!.config.traffic?.minute || 0
            )
            this.validateForm.controls['limitTrafficHour'].setValue(
              resp.data.strategy!.config.traffic?.hour || 0
            )

            for (const index in resp.data.strategy!.filters) {
              this.filterNamesSet.add(resp.data.strategy!.filters[index].name)
            }
            if (resp.data.strategy!.filters) {
              this.filterShowList = [
                ...(resp.data.strategy!.filters as Array<{
                  title: string
                  name: string
                  label: string
                  values: Array<string>
                  [key: string]: any
                }>)
              ]
            }
            this.responseHeaderList = resp.data.strategy!.config.response
              .header?.length > 0
              ? resp.data.strategy!.config.response
                .header
              : [{ key: '', value: '' }]

            this.responseForm.controls['status_code'].setValue(resp.data.strategy!.config.response.status_code || 200)
            this.responseForm.controls['content_type'].setValue(resp.data.strategy!.config.response.content_type || 'application/json')
            this.responseForm.controls['charset'].setValue(resp.data.strategy!.config.response.charset || 'UTF-8')
            this.responseForm.controls['body'].setValue(resp.data.strategy!.config.response.body || '')
          } else {
            this.message.error(resp.msg || '获取数据失败!')
          }
        }
      )
  }

  // 当用户没有编辑权限时的回调，用于disabled某些只能用nzDisabled的组件
  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  // 获取限流维度的可选选项
  getMetricsList () {
    this.api.get('strategy/metrics-options').subscribe((resp: any) => {
      if (resp.code === 0) {
        this.metricsList = []
        for (const index in resp.data.options) {
          this.metricsList.push({
            label: resp.data.options[index].title,
            value: resp.data.options[index].name
          })
        }
      } else {
        this.message.error(resp.msg || '获取数据失败!')
      }
    })
  }

  // 提交策略
  saveStrategy () {
    this.showMetricsError = this.createStrategyForm.config.metrics.length === 0

    if (
      this.validateForm.valid && this.responseForm.valid &&
      this.createStrategyForm.config.metrics.length > 0
    ) {
      delete this.createStrategyForm['extender']

      this.createStrategyForm.filters = []
      for (const index in this.filterShowList) {
        this.createStrategyForm.filters.push({
          name: this.filterShowList[index].name,
          values:
            this.filterShowList[index].name === 'ip'
              ? this.filterShowList[index].values[0].split(/[\n]/).filter(value => { return !!value })
              : this.filterShowList[index].values
        })
      }

      this.createStrategyForm.config.response.header =
        this.responseHeaderList.filter((item: any) => {
          return item.key && item.value
        })

      const data: TrafficStrategyData = {
        name: this.validateForm.controls['name'].value,
        uuid: this.createStrategyForm.uuid,
        desc: this.validateForm.controls['desc'].value,
        priority: Number(this.validateForm.controls['priority'].value),
        filters: this.createStrategyForm.filters,
        config: {
          metrics: this.createStrategyForm.config.metrics,
          query: {
            second: Number(this.validateForm.controls['limitQuerySecond'].value),
            minute: Number(this.validateForm.controls['limitQueryMinute'].value),
            hour: Number(this.validateForm.controls['limitQueryHour'].value)
          },
          traffic: {
            second: Number(this.validateForm.controls['limitTrafficSecond'].value),
            minute: Number(this.validateForm.controls['limitTrafficMinute'].value),
            hour: Number(this.validateForm.controls['limitTrafficHour'].value)
          },
          response: {
            status_code: this.responseForm.controls['status_code']
              .value
              ? Number(
                this.responseForm.controls['status_code'].value
              )
              : 0,
            content_type: this.responseForm.controls['content_type'].value || '',
            charset: this.responseForm.controls['charset'].value || '',
            header: this.createStrategyForm.config.response.header,
            body: this.responseForm.controls['body'].value || ''
          }
        }
      }

      if (!data.priority) {
        delete data.priority
      }

      if (!this.editPage) {
        this.api
          .post('strategy/traffic', data, { cluster_name: this.clusterName })
          .subscribe((resp: any) => {
            if (resp.code === 0) {
              this.message.success(resp.msg || '创建成功!', { nzDuration: 1000 })
              this.backToList()
            } else {
              this.message.error(resp.msg || '创建失败!')
            }
          })
      } else {
        this.api
          .put('strategy/traffic', data, {
            cluster_name: this.clusterName,
            uuid: this.strategyUuid
          })
          .subscribe((resp: any) => {
            if (resp.code === 0) {
              this.message.success(resp.msg || '修改成功!', { nzDuration: 1000 })
              this.backToList()
            } else {
              this.message.error(resp.msg || '修改失败!')
            }
          })
      }
    } else {
      Object.values(this.validateForm.controls).forEach((control) => {
        if (control.invalid) {
          control.markAsDirty()
          control.updateValueAndValidity({ onlySelf: true })
        }
      })
    }
  }

  // 返回列表页，当fromList为true时，该页面左侧有分组
  backToList () {
    this.router.navigate(['/', 'serv-governance', 'traffic', 'group', 'list', this.clusterName])
  }

  // 限流维度的选择需要保持顺序
  checkMetricOrder () {
    const temMetrics: Array<any> = [...this.createStrategyForm.config.metrics]
    this.createStrategyForm.config.metrics = []
    for (const index in this.metricsList) {
      if (temMetrics.indexOf(this.metricsList[index].value) !== -1) {
        this.createStrategyForm.config.metrics.push(
          this.metricsList[index].value
        )
      }
    }
    this.showMetricsError = this.createStrategyForm.config.metrics.length === 0
  }
}

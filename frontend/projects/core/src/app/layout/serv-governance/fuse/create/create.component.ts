/* eslint-disable dot-notation */
/* eslint-disable camelcase */
import {
  Component,
  EventEmitter,
  Input,
  OnInit,
  Output,
  TemplateRef,
  ViewChild
} from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { ActivatedRoute, Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'
import { FilterShowData } from '../../filter/footer/footer.component'
import { EoNgMyValidators } from 'projects/core/src/app/constant/eo-ng-validator'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'

interface FuseStrategyData {
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
  }>
  config: {
    metric: string
    fuse_condition: { status_codes: Array<number>; count: number }
    fuse_time: { time: number; max_time: number }
    recover_condition: { status_codes: Array<number>; count: number }
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
  selector: 'eo-ng-fuse-create',
  templateUrl: './create.component.html',
  styles: []
})
export class FuseCreateComponent implements OnInit {
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

  metricsList: Array<{ label: string; value: string; disable?: boolean }> = [
    { label: '上游服务', value: '{service}' },
    { label: 'API', value: '{api}' }
  ]

  createStrategyForm: FuseStrategyData = {
    name: '',
    desc: '',
    priority: null,
    filters: [],
    config: {
      metric: '{service}',
      fuse_condition: { status_codes: [500], count: 3 },
      fuse_time: { time: 2, max_time: 300 },
      recover_condition: { status_codes: [200], count: 3 },
      response: {
        status_code: 500,
        content_type: 'application/json',
        charset: 'UTF-8',
        header: [{ key: '', value: '' }],
        body: ''
      }
    }
  }

  responseHeaderList: Array<{
    key: string
    value: string
    [key: string]: any
  }> = [{ key: '', value: '' }]

  validateForm: FormGroup = new FormGroup({})
  responseForm: FormGroup = new FormGroup({})
  showFuseStatusCodeError: boolean = false
  showRecoverStatusCodeError: boolean = false
  nzDisabled: boolean = false

  constructor (
    private baseInfo:BaseInfoService,
    private message: EoNgFeedbackMessageService,
    private api: ApiService,
    private activateInfo: ActivatedRoute,
    private fb: UntypedFormBuilder,
    private appConfigService: AppConfigService,
    private router:Router
  ) {
    this.validateForm = this.fb.group({
      name: [
        '',
        [Validators.required, Validators.pattern('^[a-zA-Z][a-zA-Z0-9_]*')]
      ],
      desc: [],
      priority: [null, [EoNgMyValidators.priority]],
      configFuseCount: [3, [Validators.required]],
      configFuseTime: [2, [Validators.required]],
      configFuseMaxTime: [300, [Validators.required]],
      configRecoverCount: [3, [Validators.required]]
    })

    this.responseForm = this.fb.group({
      status_code: [
        200,
        [Validators.required, Validators.pattern(/^[1-9]{1}\d{2}$/)]
      ],
      content_type: ['application/json', [Validators.required]],
      charset: ['UTF-8', [Validators.required]],
      header: [],
      body: ['']
    })

    this.appConfigService.reqFlashBreadcrumb([
      { title: '熔断策略', routerLink: 'serv-governance/fuse/group/list' },
      { title: '新建熔断策略' }
    ])
  }


  ngOnInit (): void {
    this.clusterName = this.baseInfo.allParamsInfo.clusterName
    if (this.editPage) {
      this.getStrategyMessage()
    }
  }

  ngOnChanges (): void {
    if (this.strategyUuid) {
      this.createStrategyForm.uuid = this.strategyUuid
    }
  }

  // 当页面是编辑策略页时,需要根据集群名和策略uuid获取策略信息
  getStrategyMessage () {
    this.api
      .get('strategy/fuse', { uuid: this.createStrategyForm.uuid || '' })
      .subscribe(
        (resp: {
          code: number
          data: { strategy?: FuseStrategyData; [key: string]: any }
          msg: string
        }) => {
          if (resp.code === 0) {
            this.appConfigService.reqFlashBreadcrumb([
              {
                title: '熔断策略',
                routerLink: 'serv-governance/fuse/group/list'
              },
              { title: resp.data.strategy!.name }
            ])
            this.validateForm.controls['name'].setValue(
              resp.data.strategy!.name
            )
            this.validateForm.controls['desc'].setValue(
              resp.data.strategy!.desc
            )
            this.validateForm.controls['priority'].setValue(
              resp.data.strategy!.priority || null
            )
            this.validateForm.controls['configFuseCount'].setValue(
              resp.data.strategy!.config.fuse_condition.count || 3
            )
            this.validateForm.controls['configFuseTime'].setValue(
              resp.data.strategy!.config.fuse_time.time || 2
            )
            this.validateForm.controls['configFuseMaxTime'].setValue(
              resp.data.strategy!.config.fuse_time.max_time || 300
            )
            this.validateForm.controls['configRecoverCount'].setValue(
              resp.data.strategy!.config.recover_condition.count || 3
            )
            this.createStrategyForm = resp.data.strategy!
            this.createStrategyForm.filters =
              this.createStrategyForm.filters || []
            this.createStrategyForm.config.response.header = this
              .createStrategyForm.config.response.header || [
              { key: '', value: '' }
            ]
            for (const index in resp.data.strategy!.filters) {
              this.filterNamesSet.add(resp.data.strategy!.filters[index].name)
            }
            if (resp.data.strategy!.filters) {
              this.filterShowList = [...resp.data.strategy!.filters]
            }
            this.responseHeaderList =
              resp.data.strategy!.config.response.header.length > 0
                ? resp.data.strategy!.config.response.header
                : [{ key: '', value: '' }]

            this.responseForm.controls['status_code'].setValue(
              resp.data.strategy!.config.response.status_code || 200
            )
            this.responseForm.controls['content_type'].setValue(
              resp.data.strategy!.config.response.content_type ||
                'application/json'
            )
            this.responseForm.controls['charset'].setValue(
              resp.data.strategy!.config.response.charset || 'UTF-8'
            )
            this.responseForm.controls['body'].setValue(
              resp.data.strategy!.config.response.body || ''
            )
          } else {
            this.message.error(resp.msg || '获取数据失败!')
          }
        }
      )
  }

  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  // 提交策略时校验HTTP状态码code，如果状态码数组中包含超过4位的元素或不可转换为数字的元素，则返回true，反之则将所有元素转化为数字后，返回false
  checkStatusCode (type: string) {
    const arrayAfterChange: Array<number> = []
    switch (type) {
      case 'fuse':
        for (const el of this.createStrategyForm.config.fuse_condition
          .status_codes) {
          if (!/^[1-9]{1}\d{2}$/.test(el + '')) {
            return true
          } else {
            if (arrayAfterChange.indexOf(Number(el)) === -1) {
              arrayAfterChange.push(Number(el))
            }
          }
        }
        this.createStrategyForm.config.fuse_condition.status_codes =
          arrayAfterChange
        this.showFuseStatusCodeError = false
        return false
      case 'recover':
        for (const el of this.createStrategyForm.config.recover_condition
          .status_codes) {
          if (!/^[1-9]{1}\d{2}$/.test(el + '')) {
            return true
          } else {
            if (arrayAfterChange.indexOf(Number(el)) === -1) {
              arrayAfterChange.push(Number(el))
            }
          }
        }
        this.createStrategyForm.config.recover_condition.status_codes =
          arrayAfterChange
        this.showRecoverStatusCodeError = false
        return false
      default:
        return false
    }
  }

  // 提交策略
  saveStrategy () {
    this.showFuseStatusCodeError =
      this.createStrategyForm.config.fuse_condition.status_codes.length === 0 ||
      this.checkStatusCode('fuse')
    this.showRecoverStatusCodeError =
      this.createStrategyForm.config.recover_condition.status_codes.length ===
        0 || this.checkStatusCode('recover')
    if (
      this.validateForm.valid &&
      this.responseForm.valid &&
      this.createStrategyForm.config.fuse_condition.status_codes.length > 0 &&
      !this.showFuseStatusCodeError &&
      this.createStrategyForm.config.recover_condition.status_codes.length >
        0 &&
      !this.showRecoverStatusCodeError
    ) {
      delete this.createStrategyForm['extender']

      this.createStrategyForm.filters = []
      for (const index in this.filterShowList) {
        this.createStrategyForm.filters.push({
          name: this.filterShowList[index].name,
          values: this.filterShowList[index].values
        })
      }

      this.createStrategyForm.config.response.header =
        this.responseHeaderList.filter((item: any) => {
          return item.key && item.value
        })

      const data: FuseStrategyData = {
        name: this.validateForm.controls['name'].value || '',
        uuid: this.createStrategyForm.uuid || '',
        desc: this.validateForm.controls['desc'].value || '',
        priority: Number(this.validateForm.controls['priority'].value) || 0,
        filters: this.createStrategyForm.filters,
        config: {
          metric: this.createStrategyForm.config.metric,
          fuse_condition: {
            status_codes:
              this.createStrategyForm.config.fuse_condition.status_codes,
            count: Number(this.validateForm.controls['configFuseCount'].value)
          },
          fuse_time: {
            time: Number(this.validateForm.controls['configFuseTime'].value),
            max_time: Number(
              this.validateForm.controls['configFuseMaxTime'].value
            )
          },
          recover_condition: {
            status_codes:
              this.createStrategyForm.config.recover_condition.status_codes,
            count: Number(
              this.validateForm.controls['configRecoverCount'].value
            )
          },
          response: {
            status_code: Number(
              this.responseForm.controls['status_code'].value
            ),
            content_type: this.responseForm.controls['content_type'].value,
            charset: this.responseForm.controls['charset'].value,
            header: this.createStrategyForm.config.response.header,
            body: this.responseForm.controls['body'].value || ''
          }
        }
      }

      if (!this.editPage) {
        this.api
          .post('strategy/fuse', data, { cluster_name: this.clusterName })
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
          .put('strategy/fuse', data, {
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
    this.router.navigate(['/', 'serv-governance', 'fuse', 'group', 'list', this.clusterName])
  }
}

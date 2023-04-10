/* eslint-disable dot-notation */
import {
  Component,
  Input,
  OnInit,
  TemplateRef,
  ViewChild
} from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { EoNgNavigationService } from 'projects/core/src/app/service/eo-ng-navigation.service'
import { EoNgMyValidators } from 'projects/core/src/app/constant/eo-ng-validator'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { FilterShowData, FuseStrategyData } from '../../types/types'
import { SelectOption } from 'eo-ng-select'
import { metricsList } from '../../types/conf'
import { EmptyHttpResponse } from 'projects/core/src/app/constant/type'

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

  filterNamesSet: Set<string> = new Set() // 用户已选择的筛选条件放入set中,在显示筛选条件的选择器里需要过去set中存在的值
  validatePriority: boolean = true
  filterShowList: FilterShowData[] = [] // 展示在前端页面的筛选条件表格,包含uuid和对应选项名称,实际提交时只需要uuid
  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  metricsList:SelectOption[]= [...metricsList]
  submitButtonLoading:boolean = false
  createStrategyForm: FuseStrategyData = {
    name: '',
    desc: '',
    priority: null,
    filters: [],
    config: {
      metric: '{service}',
      fuseCondition: { statusCodes: [500], count: 3 },
      fuseTime: { time: 2, maxTime: 300 },
      recoverCondition: { statusCodes: [200], count: 3 },
      response: {
        statusCode: 500,
        contentType: 'application/json',
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
    private fb: UntypedFormBuilder,
    private appConfigService: EoNgNavigationService,
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
      statusCode: [
        200,
        [Validators.required, Validators.pattern(/^[1-9]{1}\d{2}$/)]
      ],
      contentType: ['application/json', [Validators.required]],
      charset: ['UTF-8', [Validators.required]],
      header: [],
      body: ['']
    })

    this.appConfigService.reqFlashBreadcrumb([
      { title: '熔断策略', routerLink: 'serv-governance/fuse' },
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
                routerLink: 'serv-governance/fuse'
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
              resp.data.strategy!.config.fuseCondition.count || 3
            )
            this.validateForm.controls['configFuseTime'].setValue(
              resp.data.strategy!.config.fuseTime.time || 2
            )
            this.validateForm.controls['configFuseMaxTime'].setValue(
              resp.data.strategy!.config.fuseTime.maxTime || 300
            )
            this.validateForm.controls['configRecoverCount'].setValue(
              resp.data.strategy!.config.recoverCondition.count || 3
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

            this.responseForm.controls['statusCode'].setValue(
              resp.data.strategy!.config.response.statusCode || 200
            )
            this.responseForm.controls['contentType'].setValue(
              resp.data.strategy!.config.response.contentType ||
                'application/json'
            )
            this.responseForm.controls['charset'].setValue(
              resp.data.strategy!.config.response.charset || 'UTF-8'
            )
            this.responseForm.controls['body'].setValue(
              resp.data.strategy!.config.response.body || ''
            )
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
        for (const el of this.createStrategyForm.config.fuseCondition
          .statusCodes) {
          if (!/^[1-9]{1}\d{2}$/.test(el + '')) {
            return true
          } else {
            if (arrayAfterChange.indexOf(Number(el)) === -1) {
              arrayAfterChange.push(Number(el))
            }
          }
        }
        this.createStrategyForm.config.fuseCondition.statusCodes =
          arrayAfterChange
        this.showFuseStatusCodeError = false
        return false
      case 'recover':
        for (const el of this.createStrategyForm.config.recoverCondition
          .statusCodes) {
          if (!/^[1-9]{1}\d{2}$/.test(el + '')) {
            return true
          } else {
            if (arrayAfterChange.indexOf(Number(el)) === -1) {
              arrayAfterChange.push(Number(el))
            }
          }
        }
        this.createStrategyForm.config.recoverCondition.statusCodes =
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
      this.createStrategyForm.config.fuseCondition.statusCodes.length === 0 ||
      this.checkStatusCode('fuse')
    this.showRecoverStatusCodeError =
      this.createStrategyForm.config.recoverCondition.statusCodes.length ===
        0 || this.checkStatusCode('recover')
    if (
      this.validateForm.valid &&
      this.responseForm.valid &&
      this.createStrategyForm.config.fuseCondition.statusCodes.length > 0 &&
      !this.showFuseStatusCodeError &&
      this.createStrategyForm.config.recoverCondition.statusCodes.length >
        0 &&
      !this.showRecoverStatusCodeError
    ) {
      delete this.createStrategyForm['extender']
      this.createStrategyForm.filters = []
      for (const index in this.filterShowList) {
        this.createStrategyForm.filters.push({
          name: this.filterShowList[index].name,
          values: this.filterShowList[index].name === 'ip'
            ? [...this.filterShowList[index].values[0].split(/[\n]/).filter(value => { return !!value }), ...this.filterShowList[index].values.slice(1)]
            : this.filterShowList[index].values

        })
      }

      this.createStrategyForm.config.response.header =
        this.responseHeaderList.filter((item: {key:string, value:string|number}) => {
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
          fuseCondition: {
            statusCodes:
              this.createStrategyForm.config.fuseCondition.statusCodes,
            count: Number(this.validateForm.controls['configFuseCount'].value)
          },
          fuseTime: {
            time: Number(this.validateForm.controls['configFuseTime'].value),
            maxTime: Number(
              this.validateForm.controls['configFuseMaxTime'].value
            )
          },
          recoverCondition: {
            statusCodes:
              this.createStrategyForm.config.recoverCondition.statusCodes,
            count: Number(
              this.validateForm.controls['configRecoverCount'].value
            )
          },
          response: {
            statusCode: Number(
              this.responseForm.controls['statusCode'].value
            ),
            contentType: this.responseForm.controls['contentType'].value,
            charset: this.responseForm.controls['charset'].value,
            header: this.createStrategyForm.config.response.header,
            body: this.responseForm.controls['body'].value || ''
          }
        }
      }
      this.submitButtonLoading = true

      if (!this.editPage) {
        this.api
          .post('strategy/fuse', data, { clusterName: this.clusterName })
          .subscribe((resp: EmptyHttpResponse) => {
            this.submitButtonLoading = false
            if (resp.code === 0) {
              this.message.success(resp.msg || '创建成功!', { nzDuration: 1000 })
              this.backToList()
            }
          })
      } else {
        this.api
          .put('strategy/fuse', data, {
            clusterName: this.clusterName,
            uuid: this.strategyUuid
          })
          .subscribe((resp: EmptyHttpResponse) => {
            this.submitButtonLoading = false
            if (resp.code === 0) {
              this.message.success(resp.msg || '修改成功!', { nzDuration: 1000 })
              this.backToList()
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

  // 返回列表页
  backToList () {
    this.router.navigate(['/', 'serv-governance', 'fuse', 'group', 'list', this.clusterName])
  }
}

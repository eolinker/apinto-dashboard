import { Component, Input, OnInit } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { EoNgMyValidators } from 'projects/core/src/app/constant/eo-ng-validator'
import { setFormValue } from 'projects/core/src/app/constant/form'
import { EmptyHttpResponse } from 'projects/core/src/app/constant/type'
import { EoNgNavigationService } from 'projects/core/src/app/service/eo-ng-navigation.service'
import { FilterShowData, DataMaskData } from '../../types/types'
import { DataFormatOptions, DataMaskBaseOptions, DataMaskOrderOptions, DataMaskReplaceStrOptions, MatchRules } from '../../types/conf'
import { v4 as uuidv4 } from 'uuid'

@Component({
  selector: 'eo-ng-data-mask-create',
  templateUrl: './create.component.html',
  styles: [
  ]
})
export class DataMaskCreateComponent implements OnInit {
  @Input() editPage: boolean = false
  @Input() clusterName: string = ''
  @Input() strategyUuid: string = ''

  validateForm: FormGroup = new FormGroup({})
  filterShowList: FilterShowData[] = [];
  submitButtonLoading = false;
  nzDisabled = false;

  readonly MatchRules = MatchRules;
  readonly DataFormatOptions = DataFormatOptions;
  readonly DataMaskReplaceStrOptions = DataMaskReplaceStrOptions;
  readonly DataMaskBaseOptions = DataMaskBaseOptions;
  readonly DataMaskOrderOptions = DataMaskOrderOptions;
  readonly autoTips = defaultAutoTips;

  filterNamesSet: Set<string> = new Set() // 用户已选择的筛选条件放入set中,在显示筛选条件的选择器里需要过去set中存在的值
  validatePriority: boolean = true
  createStrategyForm: DataMaskData = {
    name: '',
    desc: '',
    priority: null,
    filters: [],
    config: {
      rules: []
    }
  }

  constructor (
    private baseInfo:BaseInfoService,
    private message: EoNgFeedbackMessageService,
    private api: ApiService,
    private router:Router,
    private fb: UntypedFormBuilder,
    private navigationService: EoNgNavigationService
  ) {
    this.initBreadcrumb()
  }

  private initBreadcrumb (): void {
    this.navigationService.reqFlashBreadcrumb([
      { title: '数据脱敏', routerLink: 'serv-governance/data-mask' },
      { title: '新建数据脱敏策略' }
    ])
  }

  ngOnInit (): void {
    this.initForm()
    this.clusterName = this.baseInfo.allParamsInfo.clusterName
    if (this.editPage) {
      this.getStrategyMessage()
    }
  }

  private initForm (): void {
    this.validateForm = this.fb.group({
      name: ['', [Validators.required, Validators.pattern('^[a-zA-Z][a-zA-Z0-9_]*')]],
      desc: [],
      priority: [null, [EoNgMyValidators.priority]],
      config: {
        rules: []
      }
    })
  }

  ngOnChanges (): void {
    if (this.strategyUuid) {
      this.createStrategyForm.uuid = this.strategyUuid
    }
  }

  // 当页面是编辑策略页时,需要根据集群名和策略uuid获取策略信息
  getStrategyMessage () {
    this.api
      .get('strategy/data-mask', { uuid: this.createStrategyForm.uuid || '' })
      .subscribe(
        (resp: {
          code: number
          data: { strategy?: DataMaskData; [key: string]: any }
          msg: string
        }) => {
          if (resp.code === 0) {
            this.navigationService.reqFlashBreadcrumb([
              {
                title: '数据脱敏',
                routerLink: 'serv-governance/data-mask'
              },
              { title: resp.data.strategy!.name }
            ])

            setFormValue(this.validateForm, {
              name: resp.data.strategy!.name,
              desc: resp.data.strategy!.desc,
              priority: resp.data.strategy!.priority || null,
              config: {
                rules: resp.data.strategy!.config.rules.map((x) => ({ ...x, eoKey: uuidv4() })) || []
              }
            })

            this.createStrategyForm = resp.data.strategy!
            this.createStrategyForm.filters =
              this.createStrategyForm.filters || []

            for (const index in resp.data.strategy!.filters) {
              this.filterNamesSet.add(resp.data.strategy!.filters[index].name)
            }
            if (resp.data.strategy!.filters) {
              this.filterShowList = [...resp.data.strategy!.filters]
            }
          }
        }
      )
  }

  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  // 提交策略
  saveStrategy () {
    if (this.validateForm.valid) {
      const formData = this.validateForm.value

      const data: DataMaskData = {
        name: formData.name,
        uuid: this.strategyUuid || '',
        desc: formData.desc || '',
        priority: Number(formData.priority) || 0,
        filters: this.filterShowList.map(filter => ({
          name: filter.name,
          values: filter.name === 'ip'
            ? [...filter.values[0].split(/[\n]/).filter(Boolean), ...filter.values.slice(1)]
            : filter.values
        })),
        config: this.createStrategyForm.config
      }

      // this.submitButtonLoading = true
      if (!this.editPage) {
        this.api
          .post('strategy/data-mask', data, { clusterName: this.clusterName })
          .subscribe((resp: EmptyHttpResponse) => {
            this.submitButtonLoading = false
            if (resp.code === 0) {
              this.message.success(resp.msg || '创建成功!', { nzDuration: 1000 })
              this.backToList()
            }
          })
      } else {
        this.api
          .put('strategy/data-mask', data, {
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
    this.router.navigate(['/', 'serv-governance', 'data-mask', 'group', 'list', this.clusterName])
  }
}

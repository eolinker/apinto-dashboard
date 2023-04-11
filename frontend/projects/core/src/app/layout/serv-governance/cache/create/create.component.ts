/* eslint-disable dot-notation */
import { Component, Input, OnInit } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { EoNgNavigationService } from 'projects/core/src/app/service/eo-ng-navigation.service'
import { EoNgMyValidators } from 'projects/core/src/app/constant/eo-ng-validator'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { CacheStrategyData, FilterShowData } from '../../types/types'
import { setFormValue } from 'projects/core/src/app/constant/form'
import { EmptyHttpResponse } from 'projects/core/src/app/constant/type'

@Component({
  selector: 'eo-ng-cache-create',
  templateUrl: './create.component.html',
  styles: [
  ]
})
export class CacheCreateComponent implements OnInit {
  @Input() editPage: boolean = false
  @Input() clusterName: string = ''
  @Input() strategyUuid: string = ''

  filterNamesSet: Set<string> = new Set() // 用户已选择的筛选条件放入set中,在显示筛选条件的选择器里需要过去set中存在的值
  validatePriority: boolean = true
  filterShowList: FilterShowData[] = [] // 展示在前端页面的筛选条件表格,包含uuid和对应选项名称,实际提交时只需要uuid
  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  nzDisabled: boolean = false
  validateForm: FormGroup = new FormGroup({})
  submitButtonLoading:boolean = false
  createStrategyForm: CacheStrategyData = {
    name: '',
    desc: '',
    priority: null,
    filters: [],
    config: {
      validTime: 0
    }
  }

  constructor (
    private baseInfo:BaseInfoService,
    private message: EoNgFeedbackMessageService,
    private api: ApiService,
    private router:Router,
    private fb: UntypedFormBuilder,
    private appConfigService: EoNgNavigationService
  ) {
    this.validateForm = this.fb.group({
      name: [
        '',
        [Validators.required, Validators.pattern('^[a-zA-Z][a-zA-Z0-9_]*')]
      ],
      desc: [],
      priority: [null, [EoNgMyValidators.priority]],
      validTime: [null, [Validators.required]]
    })
    this.appConfigService.reqFlashBreadcrumb([
      { title: '缓存策略', routerLink: 'serv-governance/cache' },
      { title: '新建缓存策略' }
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
      .get('strategy/cache', { uuid: this.createStrategyForm.uuid || '' })
      .subscribe(
        (resp: {
          code: number
          data: { strategy?: CacheStrategyData; [key: string]: any }
          msg: string
        }) => {
          if (resp.code === 0) {
            this.appConfigService.reqFlashBreadcrumb([
              {
                title: '缓存策略',
                routerLink: 'serv-governance/cache'
              },
              { title: resp.data.strategy!.name }
            ])

            setFormValue(this.validateForm, {
              name: resp.data.strategy!.name,
              desc: resp.data.strategy!.desc,
              priority: resp.data.strategy!.priority || null,
              validTime: resp.data.strategy!.config.validTime || 1
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

      const data: CacheStrategyData = {
        name: this.validateForm.controls['name'].value,
        uuid: this.createStrategyForm.uuid || '',
        desc: this.validateForm.controls['desc'].value || '',
        priority: Number(this.validateForm.controls['priority'].value) || 0,
        filters: this.createStrategyForm.filters,
        config: {
          validTime: this.validateForm.controls['validTime'].value
        }
      }
      this.submitButtonLoading = true
      if (!this.editPage) {
        this.api
          .post('strategy/cache', data, { clusterName: this.clusterName })
          .subscribe((resp: EmptyHttpResponse) => {
            this.submitButtonLoading = false
            if (resp.code === 0) {
              this.message.success(resp.msg || '创建成功!', { nzDuration: 1000 })
              this.backToList()
            }
          })
      } else {
        this.api
          .put('strategy/cache', data, {
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
    this.router.navigate(['/', 'serv-governance', 'cache', 'group', 'list', this.clusterName])
  }
}

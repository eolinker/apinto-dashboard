/* eslint-disable dot-notation */
import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'
import { EoNgMyValidators } from 'projects/core/src/app/constant/eo-ng-validator'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { FilterShowData, VisitStrategyData } from '../../types/types'
import { SelectOption } from 'eo-ng-select'
import { visitRuleList } from '../../types/conf'
import { setFormValue } from 'projects/core/src/app/constant/form'
import { EmptyHttpResponse } from 'projects/core/src/app/constant/type'

@Component({
  selector: 'eo-ng-visit-create',
  templateUrl: './create.component.html',
  styles: [
  ]
})
export class VisitCreateComponent implements OnInit {
  @Input() editPage: boolean = false
  @Input() clusterName: string = ''
  @Input() strategyUuid: string = ''
  @Output() changeToList: EventEmitter<any> = new EventEmitter()

  filterNamesSet: Set<string> = new Set() // 用户已选择的筛选条件放入set中,在显示筛选条件的选择器里需要过去set中存在的值
  filterShowList: FilterShowData[] = [] // 展示在前端页面的筛选条件表格,包含uuid和对应选项名称,实际提交时只需要uuid
  influenceShowList: FilterShowData[] = [] // 展示在前端页面的生效条件表格,包含uuid和对应选项名称,实际提交时只需要uuid
  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  visitRuleList:SelectOption[]= [...visitRuleList]
  validateForm: FormGroup = new FormGroup({})
  nzDisabled: boolean = false
  createStrategyForm: VisitStrategyData = {
    name: '',
    desc: '',
    priority: null,
    filters: [],
    config: {
      visitRule: 'allow',
      influenceSphere: [],
      continue: false
    }
  }

  constructor (
    private baseInfo:BaseInfoService,
    private message: EoNgFeedbackMessageService,
    private api: ApiService,
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
      continue: [false, [Validators.required]]
    })
    this.appConfigService.reqFlashBreadcrumb([
      { title: '访问策略', routerLink: 'serv-governance/visit' },
      { title: '新建访问策略' }
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
      .get('strategy/visit', { uuid: this.createStrategyForm.uuid || '' })
      .subscribe(
        (resp: {
          code: number
          data: { strategy?: VisitStrategyData; [key: string]: any }
          msg: string
        }) => {
          if (resp.code === 0) {
            this.appConfigService.reqFlashBreadcrumb([
              {
                title: '访问策略',
                routerLink: 'serv-governance/visit'
              },
              { title: resp.data.strategy!.name }
            ])
            setFormValue(this.validateForm, {
              name: resp.data.strategy!.name,
              desc: resp.data.strategy!.desc,
              priority: resp.data.strategy!.priority || null,
              continue: resp.data.strategy!.config.continue || false
            })
            this.createStrategyForm = resp.data.strategy!
            this.createStrategyForm.filters = this.createStrategyForm.filters || []
            this.createStrategyForm.config.influenceSphere = this.createStrategyForm.config.influenceSphere || []
            for (const index in resp.data.strategy!.filters) {
              this.filterNamesSet.add(resp.data.strategy!.filters[index].name)
            }
            if (resp.data.strategy!.filters) {
              this.filterShowList = [...resp.data.strategy!.filters]
            }
            for (const index in resp.data.strategy!.config.influenceSphere) {
              this.filterNamesSet.add(resp.data.strategy!.config.influenceSphere[index].name)
            }
            if (resp.data.strategy!.filters) {
              this.influenceShowList = [...resp.data.strategy!.config.influenceSphere]
            }
          } else {
            this.message.error(resp.msg || '获取数据失败!')
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

      this.createStrategyForm.config.influenceSphere = []
      for (const index in this.influenceShowList) {
        this.createStrategyForm.config.influenceSphere.push({
          name: this.influenceShowList[index].name,
          values: this.influenceShowList[index].name === 'ip'
            ? [...this.influenceShowList[index].values[0].split(/[\n]/).filter(value => { return !!value }), ...this.influenceShowList[index].values.slice(1)]
            : this.influenceShowList[index].values
        })
      }

      const data: VisitStrategyData = {
        name: this.validateForm.controls['name'].value,
        uuid: this.createStrategyForm.uuid || '',
        desc: this.validateForm.controls['desc'].value || '',
        priority: Number(this.validateForm.controls['priority'].value) || 0,
        filters: this.createStrategyForm.filters,
        config: {
          visitRule: this.createStrategyForm.config.visitRule || '',
          influenceSphere: this.createStrategyForm.config.influenceSphere,
          continue: this.validateForm.controls['continue'].value || false
        }
      }

      if (!this.editPage) {
        this.api
          .post('strategy/visit', data, { clusterName: this.clusterName })
          .subscribe((resp: EmptyHttpResponse) => {
            if (resp.code === 0) {
              this.message.success(resp.msg || '创建成功!', { nzDuration: 1000 })
              this.backToList()
            } else {
              this.message.error(resp.msg || '创建失败!')
            }
          })
      } else {
        this.api
          .put('strategy/visit', data, {
            clusterName: this.clusterName,
            uuid: this.strategyUuid
          })
          .subscribe((resp: EmptyHttpResponse) => {
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

  // 返回列表页
  backToList () {
    this.router.navigate(['/', 'serv-governance', 'visit', 'group', 'list', this.clusterName])
  }
}

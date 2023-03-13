/* eslint-disable @typescript-eslint/no-unused-vars */
/* eslint-disable dot-notation */
/* eslint-disable camelcase */
import { Component, EventEmitter, Input, OnInit, Output, TemplateRef, ViewChild } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { ActivatedRoute, Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'
import { FilterShowData } from '../../filter/footer/footer.component'
import { EoNgMyValidators } from 'projects/core/src/app/constant/eo-ng-validator'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'

interface VisitStrategyData {
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
    visit_rule: string
    influence_sphere: Array<{
      name: string
      values: Array<string>
      type?: string
      label?: string
      title?: string
    }>,
    continue:boolean
  }
  [key: string]: any
}

@Component({
  selector: 'eo-ng-visit-create',
  templateUrl: './create.component.html',
  styles: [
  ]
})
export class VisitCreateComponent implements OnInit {
  @ViewChild('filterContentRef', { read: TemplateRef, static: true })
  filterContentRef: TemplateRef<any> | undefined

  @ViewChild('filterFooterRef', { read: TemplateRef, static: true })
  filterFooterRef: TemplateRef<any> | undefined

  @ViewChild('influenceFooterRef', { read: TemplateRef, static: true })
  influenceFooterRef: TemplateRef<any> | undefined

  @ViewChild('filterTableLabel', { read: TemplateRef, static: true })
  filterTableLabel: TemplateRef<any> | undefined

  @Input() editPage: boolean = false
  @Input() clusterName: string = ''
  @Input() strategyUuid: string = ''
  @Input() fromList: boolean = false
  @Output() changeToList: EventEmitter<any> = new EventEmitter()

  filterNamesSet: Set<string> = new Set() // 用户已选择的筛选条件放入set中,在显示筛选条件的选择器里需要过去set中存在的值

  filterShowList: FilterShowData[] = [] // 展示在前端页面的筛选条件表格,包含uuid和对应选项名称,实际提交时只需要uuid

  influenceShowList: FilterShowData[] = [] // 展示在前端页面的生效条件表格,包含uuid和对应选项名称,实际提交时只需要uuid

  autoTips: Record<string, Record<string, string>> = defaultAutoTips

  visitRuleList: Array<{ label: string; value: string; disable?: boolean }> = [
    { label: '允许', value: 'allow' },
    { label: '拒绝', value: 'refuse' }
  ]

  createStrategyForm: VisitStrategyData = {
    name: '',
    desc: '',
    priority: null,
    filters: [],
    config: {
      visit_rule: 'allow',
      influence_sphere: [],
      continue: false
    }
  }

  validateForm: FormGroup = new FormGroup({})
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
      continue: [false, [Validators.required]]
    })
    this.appConfigService.reqFlashBreadcrumb([
      { title: '访问策略', routerLink: 'serv-governance/visit/group/list' },
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
                routerLink: 'serv-governance/visit/group/list'
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
            this.validateForm.controls['continue'].setValue(
              resp.data.strategy!.config.continue || false
            )
            this.createStrategyForm = resp.data.strategy!
            this.createStrategyForm.filters =
              this.createStrategyForm.filters || []

            this.createStrategyForm.config.influence_sphere =
            this.createStrategyForm.config.influence_sphere || []
            for (const index in resp.data.strategy!.filters) {
              this.filterNamesSet.add(resp.data.strategy!.filters[index].name)
            }
            if (resp.data.strategy!.filters) {
              this.filterShowList = [...resp.data.strategy!.filters]
            }

            for (const index in resp.data.strategy!.config.influence_sphere) {
              this.filterNamesSet.add(resp.data.strategy!.config.influence_sphere[index].name)
            }
            if (resp.data.strategy!.filters) {
              this.influenceShowList = [...resp.data.strategy!.config.influence_sphere]
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
          values: this.filterShowList[index].values
        })
      }

      this.createStrategyForm.config.influence_sphere = []
      for (const index in this.influenceShowList) {
        this.createStrategyForm.config.influence_sphere.push({
          name: this.influenceShowList[index].name,
          values: this.influenceShowList[index].values
        })
      }

      const data: VisitStrategyData = {
        name: this.validateForm.controls['name'].value,
        uuid: this.createStrategyForm.uuid || '',
        desc: this.validateForm.controls['desc'].value || '',
        priority: Number(this.validateForm.controls['priority'].value) || 0,
        filters: this.createStrategyForm.filters,
        config: {
          visit_rule: this.createStrategyForm.config.visit_rule || '',
          influence_sphere: this.createStrategyForm.config.influence_sphere,
          continue: this.validateForm.controls['continue'].value || false
        }
      }

      if (!this.editPage) {
        this.api
          .post('strategy/visit', data, { cluster_name: this.clusterName })
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
          .put('strategy/visit', data, {
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
    this.router.navigate(['/', 'serv-governance', 'visit', 'group', 'list', this.clusterName])
  }
}

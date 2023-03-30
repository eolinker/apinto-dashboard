/* eslint-disable dot-notation */
import { Component, Input, OnInit } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'
import { EoNgMyValidators } from 'projects/core/src/app/constant/eo-ng-validator'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { FilterShowData, GreyStrategyData } from '../../types/types'
import { RadioOption } from 'eo-ng-radio'
import { distributionOptions, nodesTableBody } from '../../types/conf'
import { setFormValue } from 'projects/core/src/app/constant/form'
import { EmptyHttpResponse } from 'projects/core/src/app/constant/type'

@Component({
  selector: 'eo-ng-grey-create',
  templateUrl: './create.component.html',
  styles: [
    `
    nz-slider{
      width:318px;
    }

    nz-sider{
      width:318px;
      padding:0px;
      margin:0px;
      display:inline-block;
      vertical-align:middle;
    }

    nz-input-number{
      display:inline-block;
      vertical-align:middle;
      width:80px;
    }
    `
  ]
})
export class GreyCreateComponent implements OnInit {
  @Input() editPage: boolean = false
  @Input() clusterName: string = ''
  @Input() strategyUuid: string = ''
  @Input() fromList: boolean = false

  filterNamesSet: Set<string> = new Set() // 用户已选择的筛选条件放入set中,在显示筛选条件的选择器里需要过去set中存在的值
  filterShowList:FilterShowData[] = [] // 展示在前端页面的筛选条件表格,包含uuid和对应选项名称,实际提交时只需要uuid
  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  distributionOptions: RadioOption[] = [...distributionOptions]
  submitButtonLoading:boolean = false

  createStrategyForm: GreyStrategyData = {
    name: '',
    desc: '',
    priority: null,
    filters: [],
    config: {
      keepSession: false,
      nodes: [],
      distribution: 'percent',
      match: []
    }
  }

  nodesList:Array<{node:string, [key:string]:any}> = [{ node: '' }]

  nodesTableBody: Array<any> = [...nodesTableBody]

  validateForm: FormGroup = new FormGroup({})
  nzDisabled: boolean = false
  showNodesValid:boolean = false

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
      keepSession: [false, [Validators.required]],
      distribution: ['percent', [Validators.required]],
      percent1: [1],
      percent2: [99]
    })

    this.appConfigService.reqFlashBreadcrumb([
      { title: '灰度策略', routerLink: 'serv-governance/grey' },
      { title: '新建灰度策略' }
    ])
  }

  ngOnInit (): void {
    this.clusterName = this.baseInfo.allParamsInfo.clusterName
    if (this.editPage) {
      this.getStrategyMessage()
    }
    this.initTable()
  }

  initTable () {
    this.nodesTableBody[0].check = (item:any) => {
      if (!/^((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}:[0-9]+$/.test(item) &&
      !/^[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+\.?(:[0-9]+)?$/.test(item)) {
        return false
      }
      return true
    }
    this.nodesTableBody[0].change = () => {
      this.showNodesValid = false
    }
    this.nodesTableBody[0].disabledFn = () => {
      return this.nzDisabled
    }
    this.nodesTableBody[1].showFn = (item:any) => {
      return item === this.nodesList[0]
    }
    this.nodesTableBody[1].btns[0].disabledFn = () => {
      return this.nzDisabled
    }
    this.nodesTableBody[2].showFn = (item:any) => {
      return item !== this.nodesList[0]
    }
    this.nodesTableBody[2].btns[0].disabledFn = () => {
      return this.nzDisabled
    }
    this.nodesTableBody[2].btns[1].disabledFn = () => {
      return this.nzDisabled
    }
  }

  ngOnChanges (): void {
    if (this.strategyUuid) {
      this.createStrategyForm.uuid = this.strategyUuid
    }
  }

  drawerTipShowFn = () => {
    return !this.filterNamesSet.has('api') && !this.filterNamesSet.has('service')
  }

  // 当页面是编辑策略页时,需要根据集群名和策略uuid获取策略信息
  getStrategyMessage () {
    this.api
      .get('strategy/grey', { uuid: this.createStrategyForm.uuid || '' })
      .subscribe(
        (resp: {
          code: number
          data: { strategy?: GreyStrategyData; [key: string]: any }
          msg: string
        }) => {
          if (resp.code === 0) {
            this.appConfigService.reqFlashBreadcrumb([
              {
                title: '灰度策略',
                routerLink: 'serv-governance'
              },
              { title: resp.data.strategy!.name }
            ])

            setFormValue(this.validateForm, {
              name: resp.data.strategy!.name,
              desc: resp.data.strategy!.desc,
              priority: resp.data.strategy!.priority || null,
              keepSession: resp.data.strategy!.config.keepSession || false,
              distribution: resp.data.strategy!.config.distribution || 'percent'
            })

            if (resp.data.strategy!.config.distribution === 'percent') {
              this.validateForm.controls['percent1'].setValue(
              resp.data.strategy!.config.percent! / 100 || 1
              )

              this.validateForm.controls['percent2'].setValue(
                100 - (resp.data.strategy!.config.percent! / 100 || 1)
              )
            }
            this.createStrategyForm = resp.data.strategy!
            this.createStrategyForm.filters = this.createStrategyForm.filters || []
            this.createStrategyForm.config.match = this.createStrategyForm.config.match || []
            this.nodesList = []

            for (const index in resp.data.strategy!.filters) {
              this.filterNamesSet.add(resp.data.strategy!.filters[index].name)
            }
            if (resp.data.strategy!.filters) {
              this.filterShowList = [...resp.data.strategy!.filters]
            }

            for (const index in resp.data.strategy?.config.nodes) {
              this.nodesList.push({ node: resp.data.strategy?.config.nodes[index as any] || '' })
            }
            this.nodesList = this.nodesList.length > 0 ? this.nodesList : [{ node: '' }]
          }
        }
      )
  }

  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  // 检查nodeslist是否有值，否则显示提示信息
  checkNodesList () {
    if (this.nodesList.length === 0) {
      this.showNodesValid = true
    } else {
      for (const index in this.nodesList) {
        if (this.nodesList[index].node) {
          this.showNodesValid = false
          return
        }
      }
      this.showNodesValid = true
    }
  }

  // 提交策略
  saveStrategy () {
    this.checkNodesList()

    if (this.validateForm.valid && !this.showNodesValid) {
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

      this.createStrategyForm.config.nodes = []
      for (const index in this.nodesList) {
        this.nodesList[index].node && this.createStrategyForm.config.nodes.push(this.nodesList[index].node)
      }

      const data: GreyStrategyData = {
        name: this.validateForm.controls['name'].value,
        uuid: this.createStrategyForm.uuid || '',
        desc: this.validateForm.controls['desc'].value || '',
        priority: Number(this.validateForm.controls['priority'].value),
        filters: this.createStrategyForm.filters || [],
        config: {
          keepSession: this.validateForm.controls['keepSession'].value || false,
          nodes: this.createStrategyForm.config.nodes,
          distribution: this.validateForm.controls['distribution'].value,
          percent: Number(this.validateForm.controls['percent1'].value) * 100 || 0,
          match: this.createStrategyForm.config.match || []
        }
      }

      this.validateForm.controls['distribution'].value === 'percent' ? delete data.config.match : delete data.config.percent
      if (!this.validateForm.controls['priority'].value) { delete data.priority }
      this.submitButtonLoading = true
      if (!this.editPage) {
        this.api
          .post('strategy/grey', data, { clusterName: this.clusterName })
          .subscribe((resp: EmptyHttpResponse) => {
            this.submitButtonLoading = false
            if (resp.code === 0) {
              this.message.success(resp.msg || '创建成功!', { nzDuration: 1000 })
              this.backToList()
            }
          })
      } else {
        this.api
          .put('strategy/grey', data, {
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

  // 返回列表页，当fromList为true时，该页面左侧有分组
  backToList () {
    this.router.navigate(['/', 'serv-governance', 'grey', 'group', 'list', this.clusterName])
  }

  changePercent (value:any, controlName:string) {
    if (value > 100) {
      value = 100
    } else if (value < 0) {
      value = 0
    }
    this.validateForm.controls[controlName].setValue(value)

    if (controlName === 'percent1') {
      this.validateForm.controls['percent2'].setValue(100 - value)
    } else {
      this.validateForm.controls['percent1'].setValue(100 - value)
    }
  }
}

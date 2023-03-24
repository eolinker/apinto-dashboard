/* eslint-disable dot-notation */
import { ViewportScroller } from '@angular/common'
import { Component, Input, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { NavigationEnd, Router } from '@angular/router'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { Subscription } from 'rxjs'
import { MODAL_NORMAL_SIZE, MODAL_SMALL_SIZE } from '../../../constant/app.config'
import { EmptyHttpResponse } from '../../../constant/type'
import { ApiService } from '../../../service/api.service'
import { BaseInfoService } from '../../../service/base-info.service'
import { ServiceGovernancePublishComponent } from '../publish/publish.component'
import { strategiesTableBody, strategiesTableHeadName } from '../types/conf'
import { StrategyListData } from '../types/types'

@Component({
  selector: 'eo-ng-serv-governance-list',
  templateUrl: './list.component.html',
  styles: [
    `
  input[eo-ng-input].strategy-priority-input.ant-input:not(.w206):not(.w131):not(.w240){
    width:calc(100% - 1px) !important;
    height:38px;
    min-width:auto !important;
    padding:0 16px !important;
    background-color:transparent;
    text-align: center;
  }`
  ]
})
export class ListComponent implements OnInit {
  @ViewChild('switchTpl', { read: TemplateRef, static: true }) switchTpl: TemplateRef<any> | undefined
  @ViewChild('priorityTpl', { read: TemplateRef, static: true }) priorityTpl: TemplateRef<any> | undefined
  @ViewChild('strategiesStatusTpl', { read: TemplateRef, static: true }) strategiesStatusTpl: TemplateRef<any> | undefined
  @ViewChild('isDisableTpl', { read: TemplateRef, static: true }) isDisableTpl: TemplateRef<any> | undefined
  @Input() clusterName:string = ''
  cluster:string = ''
  nzDisabled:boolean = false
  drawerPublishRef:NzModalRef | undefined
  priorityMap : Map<number|string, StrategyListData[]> = new Map() // 优先级与对应策略的map，随着修改而改变
  priorityOldMap: Map<number|string, number> = new Map() // 策略uuid与优先级的map，接口赋值后不变，用于提交优先级修改时做对比
  prioritySaveMap:{[k:string]:number} = {} // 用于提交的优先级，键为uuid，值为优先级
  priorityDanger:Array<string> = [] // 优先级冲突的uuid
  priorityDangerP:Array<number> = [] // 优先级冲突的priority

  strategyType:string = ''
  private subscription: Subscription = new Subscription()

  strategiesTableHeadName:THEAD_TYPE[]= [...strategiesTableHeadName]
  strategiesTableBody:TBODY_TYPE[] = [...strategiesTableBody]

  strategiesList:Array<StrategyListData> = []

  editingPriority:number|string|null = null
  // eslint-disable-next-line no-useless-constructor
  constructor (private baseInfo:BaseInfoService,
    private viewportScroller: ViewportScroller,
                private message: EoNgFeedbackMessageService,
                private modalService:EoNgFeedbackModalService,
                private api:ApiService,
                private router:Router) {
    this.strategyType = this.router.url.split('/')[2]
  }

  ngOnInit (): void {
    this.clusterName = this.baseInfo.allParamsInfo.clusterName
    this.getStrategiesList()
    this.initTable()
    // 当左侧集群被选中时,clusterName参数会发生变化,随之获取新的策略列表
    this.subscription = this.router.events.subscribe((event) => {
      if (event instanceof NavigationEnd) {
        this.clusterName = this.baseInfo.allParamsInfo.clusterName
        this.getStrategiesList()
      }
    })
  }

  ngAfterContentInit () {
    this.strategiesTableBody[3].title = this.switchTpl
    this.strategiesTableBody[1].title = this.priorityTpl
    this.strategiesTableBody[2].title = this.strategiesStatusTpl
    this.strategiesTableHeadName[3].title = this.isDisableTpl
    this.strategiesTableHeadName[5].title = this.strategyConf()
  }

  ngOnDestroy () {
    this.subscription.unsubscribe()
  }

  initTable () {
    this.strategiesTableBody[8].btns[0].click = (item:any) => {
      this.editStrategy(item.data.uuid)
    }
    this.strategiesTableBody[8].btns[1].click = (item:any) => {
      this.deleteModal(item.data)
    }
    this.strategiesTableBody[8].btns[1].disabledFn = () => {
      return this.nzDisabled
    }

    this.strategiesTableBody[9].btns[0].click = (item:any) => {
      this.editStrategy(item.data.uuid)
    }

    this.strategiesTableBody[9].btns[1].click = (item:any) => {
      this.recoverStrategy(item.data)
    }
    this.strategiesTableBody[9].btns[1].disabledFn = () => {
      return this.nzDisabled
    }
  }

  strategiesTableClick = (item:any) => {
    this.editStrategy(item.data.uuid)
  }

  // 修改优先级-当输入框获得焦点时，获取当时输入框内的优先级
  changeEditingPriority (e:Event, value:string|number) {
    e?.stopPropagation()
    if (!value) {
      this.editingPriority = 'NULL'
    } else {
      this.editingPriority = Number(value)
    }
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  deleteModal (item:StrategyListData, e?:Event) {
    e?.stopPropagation()
    this.modalService.create({
      nzTitle: '删除',
      nzContent: '该数据删除后将无法找回，请确认是否删除？',
      nzClosable: true,
      nzClassName: 'delete-modal',
      nzWidth: MODAL_SMALL_SIZE,
      nzOkDanger: true,
      nzOnOk: () => {
        this.deleteStrategy(item)
      }
    })
  }

  // 修改优先级, 当输入框失焦或enter时,
  // 将策略从priorityMap的editingPriority中移除, 放入priority的数组中,并将策略的uuid和对应priority放入prioritySaveMap中,
  // 再检查:
  // 1.输入为空或大于999时, 放入priorityMap中key为'NULL'的数组中,并提示优先级为空或大于999不允许提交
  // 2.输入不为空, 检查priorityMap中相同priority的数组, 如果数组内包含其他1个策略, 提示冲突的策略名, 滚动到相应策略, 不允许提交
  // 3.输入不为空, 检查priorityMap中相同priority的数组, 如果数组内包含其他多个策略, 提示有多个冲突, 不允许提交
  // 4.输入不为空, 检查priorityMap中其他priority的数组, 如果数组内包含其他多个策略, 提示有多个冲突, 不允许提交
  // 5.输入不为空, 检查priorityMap中其他priority的数组, 如果全部数组的策略个数小于等于1, 允许提交
  // 检查priorityMap,将所有冲突策略的优先级放入priorityDangerP中, 以便页面中的input检测状态
  checkPriority (priority:number|string, uuid:string) {
    if (!priority) {
      priority = 'NULL'
    } else {
      priority = Number(priority)
    }
    if (this.editingPriority && this.editingPriority !== priority) {
      this.changePriorityMap(this.editingPriority, priority, uuid)
    }
    // 1.输入为空
    if (priority === 'NULL') {
      this.message.error('优先级不能为空，请填写后提交')
    } else if (Number(priority) > 999) {
      this.message.error('优先级范围在1-999之间，请修改后提交')
    } else {
      // 2.输入不为空, 检查priorityMap中相同priority的数组, 如果数组内包含其他1个策略, 提示冲突的策略名, 滚动到相应策略, 不允许提交
      if (this.priorityMap.get(priority)?.length === 2) {
        const anotherStrategy = this.priorityMap.get(priority)![0].uuid === uuid ? this.priorityMap.get(priority)![1] : this.priorityMap.get(priority)![0]
        this.viewportScroller.scrollToAnchor(anotherStrategy.uuid)
        this.message.error(`修改后的优先级与${anotherStrategy.name}冲突，无法自动提交`)
        if (this.priorityDangerP.indexOf(Number(priority)) === -1) {
          this.priorityDangerP.push(Number(priority))
        }
        // 3.输入不为空, 检查priorityMap中相同priority的数组, 如果数组内包含其他多个策略, 提示有多个冲突, 不允许提交
      } else if (this.priorityMap.get(priority)!.length > 2) {
        this.message.error('优先级存在冲突或数值超出范围，无法自动提交')
        if (this.priorityDangerP.indexOf(Number(priority)) === -1) {
          this.priorityDangerP.push(Number(priority))
        }
        // 4.输入不为空, 检查priorityMap中其他priority的数组, 如果数组内包含其他多个策略, 提示有多个冲突, 不允许提交
        // 5.输入不为空, 检查priorityMap中其他priority的数组, 如果全部数组的策略个数小于等于1, 允许提交
        // 检查priorityMap,将所有冲突策略的优先级放入priorityDangerP中, 以便页面中的input检测状态
      } else {
        if (this.checkPriorityMap()) {
          this.editingPriority !== priority && this.changePriority()
        } else {
          this.message.error('优先级存在冲突或数值超出范围，无法自动提交')
        }
      }
    }
  }

  // 遍历priorityMap,如果所有value数组长度都小于等于1, 返回true, 否则返回false
  // 当数组长度大于1时,将冲突优先级放入priorityDanger中
  checkPriorityMap ():boolean {
    let res:boolean = true
    this.priorityDangerP = []
    for (const [key, value] of this.priorityMap) {
      if (value.length > 1) {
        res = false
        if (key !== 'NULL') {
          this.priorityDangerP.push(Number(key))
        }
      }
      if (key === 'NULL' && value.length > 0) {
        res = false
      }
    }
    return res
  }

  // 将策略从priorityMap的editingPriority中移除, 放入priority的数组中
  // editingPriority:修改前的优先级
  // priority:修改后的优先级
  // uuid:当前策略的uuid
  changePriorityMap (editingPriority:number|string, priority:string|number, uuid:string) {
    let tempStrategy:StrategyListData|undefined
    const tempMapArray:Array<StrategyListData> = this.priorityMap.get(editingPriority)?.filter((item:StrategyListData) => {
      if (item?.uuid === uuid) {
        tempStrategy = item
      }
      return item && item.uuid !== uuid
    }) || []

    this.priorityMap.set(editingPriority, tempMapArray)

    if (this.priorityMap.get(priority)) {
      this.priorityMap.get(priority)!.push(tempStrategy!)
    } else {
      this.priorityMap.set(priority, [tempStrategy!])
    }
  }

  checkListStatus (data:string|number) {
    return this.priorityDangerP.indexOf(Number(data)) !== -1 || !data ? 'error' : ''
  }

  changePriority () {
    this.prioritySaveMap = {}
    for (const [key, value] of this.priorityMap) {
      if (value.length === 1 && key !== 'NULL' && this.priorityOldMap.get(value[0].uuid) !== Number(key)) {
        this.prioritySaveMap[value[0].uuid] = Number(key)
      }
    }
    this.api.post('strategy/' + this.strategyType + '/priority', this.prioritySaveMap, { clusterName: (this.clusterName || '') })
      .subscribe((resp:EmptyHttpResponse) => {
        if (resp.code === 0) {
          this.message.success(resp.msg || '修改优先级成功！', { nzDuration: 1000 })
          this.getStrategiesList()
        } else {
          this.message.error(resp.msg || '修改优先级失败!')
        }
      })
  }

  // 获取策略列表
  getStrategiesList () {
    if (this.clusterName) {
      this.api.get('strategies/' + this.strategyType, { clusterName: this.clusterName || '' })
        .subscribe((resp:{code:number, data:{strategies:StrategyListData[]}, msg:string}) => {
          if (resp.code === 0) {
            this.priorityMap = new Map()
            this.strategiesList = resp.data.strategies
            for (const index in this.strategiesList) {
              this.priorityMap.set(this.strategiesList[index].priority, [this.strategiesList[index]])
              this.priorityOldMap.set(this.strategiesList[index].uuid, this.strategiesList[index].priority)
            }
          } else {
            this.message.error(resp.msg || '获取数据失败!')
          }
        })
    }
  }

  // 新建策略
  addStrategy () {
    this.router.navigate(['/', 'serv-governance', this.strategyType, 'create', this.clusterName])
  }

  // 编辑策略
  editStrategy (uuid:string) {
    this.router.navigate(['/', 'serv-governance', this.strategyType, 'message', this.clusterName, uuid])
  }

  // 当参数为true时,停用策略,false为启用
  stopStrategy (e:Event, item:StrategyListData, isStop:boolean) {
    e.stopPropagation()
    this.api.patch('strategy/' + this.strategyType + '/stop', { isStop: isStop }, { uuid: (item.uuid || ''), clusterName: (this.clusterName || '') })
      .subscribe((resp:EmptyHttpResponse) => {
        if (resp.code === 0) {
          item.isStop = isStop
          this.message.success((isStop ? '停用' : '启用') + '策略成功', { nzDuration: 1000 })
          this.getStrategiesList()
        } else {
          this.message.error(resp.msg || ((isStop ? '停用' : '启用') + '策略失败!'))
        }
      })
  }

  // 删除策略
  deleteStrategy (item:StrategyListData) {
    this.api.delete('strategy/' + this.strategyType, { uuid: (item.uuid || ''), clusterName: (this.clusterName || '') })
      .subscribe((resp:EmptyHttpResponse) => {
        if (resp.code === 0) {
          this.message.success(resp.msg || '删除策略成功!', { nzDuration: 1000 })
          this.getStrategiesList()
        } else {
          this.message.error(resp.msg || '删除策略失败!')
        }
      })
  }

  // 恢复策略
  recoverStrategy (item:StrategyListData, e?:Event) {
    e?.stopPropagation()
    this.api.patch('strategy/' + this.strategyType + '/restore', {}, { uuid: (item.uuid || ''), clusterName: (this.clusterName || '') })
      .subscribe((resp:EmptyHttpResponse) => {
        if (resp.code === 0) {
          this.getStrategiesList()
          this.message.success(resp.msg || '恢复策略成功!', { nzDuration: 1000 })
        } else {
          this.message.error(resp.msg || '恢复策略失败!')
        }
      })
  }

  // 打开弹窗
  openDrawer (type:string) {
    switch (type) {
      case 'publish': {
        this.drawerPublishRef = this.modalService.create({
          nzTitle: '发布策略',
          nzWidth: MODAL_NORMAL_SIZE,
          nzContent: ServiceGovernancePublishComponent,
          nzComponentParams: {
            strategyType: this.strategyType,
            clusterName: this.clusterName,
            strategiesStatusTpl: this.strategiesStatusTpl,
            closeModal: this.cancelDrawer
          },
          nzOkDisabled: this.nzDisabled,
          nzOnOk: (component:ServiceGovernancePublishComponent) => {
            component.publish()
            return false
          }
        }) }
    }
  }

  // 关闭弹窗
  cancelDrawer = () => {
    this.drawerPublishRef?.close()
    this.getStrategiesList()
  }

  // 不同策略在列表展示的规则标题不同
  strategyConf () {
    switch (this.strategyType) {
      case 'traffic':
        return '限流规则'
      case 'grey':
        return '灰度规则'
      case 'fuse':
        return '熔断维度'
      case 'cache':
        return '缓存有效时间'
      case 'visit':
        return '访问规则'
      default:
        return '限流规则'
    }
  }
}

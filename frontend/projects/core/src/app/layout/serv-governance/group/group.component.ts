/* eslint-disable dot-notation */
import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { Router, ActivatedRoute } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { EoNgTreeDefaultComponent } from 'eo-ng-tree'
import { NzTreeNodeOptions, NzTreeNode } from 'ng-zorro-antd/tree'
import { } from 'rxjs'
import { ApiService } from '../../../service/api.service'
import { AppConfigService } from '../../../service/app-config.service'
import { BaseInfoService } from '../../../service/base-info.service'

@Component({
  selector: 'eo-ng-flow-control-group',
  templateUrl: './group.component.html',
  styles: ['']
})
export class GroupComponent implements OnInit {
  @ViewChild('addGroupRef', { read: TemplateRef, static: true }) addGroupRef:
    | TemplateRef<any>
    | string = ''

  @ViewChild('eoNgTreeDefault')
  eoNgTreeDefault!: EoNgTreeDefaultComponent

  public nodesList: NzTreeNodeOptions[] = []
  public apiNodesList: Array<any> = []
  showList: boolean = true // 右侧是否展示列表
  showApiPage: boolean = false // 右侧是否展示api页面
  groupUuid: string = '' // 供右侧list页面用
  strategyUuid: string = '' // 供右侧策略详情页用
  editPage: boolean = false // 右侧策略详情是否是编辑页
  activatedNode: NzTreeNode | null = null
  groupName: string = ''
  clusterName: string = ''
  clusterKey: string = ''
  strategyType: string = ''
  constructor (
    private baseInfo:BaseInfoService,
    private message: EoNgFeedbackMessageService,
    private api: ApiService,
    private activateInfo: ActivatedRoute,
    private router: Router,
    private appConfigService:AppConfigService
  ) {
    this.strategyType = this.router.url.split('/')[2]
    this.getBreadcrumb()
  }

  ngOnInit (): void {
    this.getGroupList()
    this.clusterName = this.baseInfo.allParamsInfo.clusterName
    if (this.clusterName) {
      this.getGroupItemSelected()
    }
  }

  // 获取分组列表, 将api返回的数据通过nodesTransfer转化为分组组件需要的参数格式
  getGroupList () {
    this.api.get('cluster/enum').subscribe((resp: any) => {
      if (resp.code === 0) {
        this.nodesList = []
        this.nodesList = this.nodesTransfer(resp.data.envs)
        this.getGroupItemSelected()
      } else {
        this.message.error(resp.msg || '获取数据失败!')
      }
    })
  }

  // 遍历目录列表，转化成tree组件需要的参数格式
  // input: {name:string, clusters;Array<{name:string}>}
  // output: {title:string, key:string, isLeaf?:boolean, children:[]}
  // key为{集群名}_{环境}
  // 集群列表只有两层,第一层为环境名,第二层为集群名
  // 当分组未有选中节点时,默认第一个集群为选中节点, 并展开当前环境的集群列表
  nodesTransfer (data: any): NzTreeNodeOptions[] {
    const res: NzTreeNodeOptions[] = []
    for (const index in data) {
      data[index].key = data[index].name
      data[index].title = data[index].name
      if (data[index].clusters?.length > 0) {
        data[index].children = this.clustersTransfer(
          data[index].clusters,
          data[index].name
        )
        if (index === '0' && !this.activatedNode && !this.clusterName) {
          data[index].expanded = true
          data[index].children[0].selected = true
          this.clusterName = data[index].children[0].name
          this.clusterKey = data[index].children[0].key
          if (this.clusterName) {
            this.router.navigate(
              ['/', 'serv-governance', this.strategyType, 'group', 'list', this.clusterName]
            )
          }
        }
      } else {
        data[index].isLeaf = true
      }
      res.push(data[index])
    }
    return res
  }

  getGroupItemSelected () {
    for (const index in this.nodesList) {
      for (const indexC in this.nodesList[index].children) {
        const child = this.nodesList[index].children![indexC as any]
        if (this.clusterName && child['name'] === this.clusterName) {
          this.nodesList[index].expanded = true
          this.nodesList[index].children![indexC as any].selected = true
          this.clusterKey = this.nodesList[index].children![indexC as any].key
        } else {
          this.nodesList[index].children![indexC as any].selected = false
        }
      }
    }
  }

  // 将集群列表转化为group的入参格式
  // input:  Array<{name:string}>
  // output: [{title:string, key:string, isLeaf?:boolean ]
  clustersTransfer (data: any, env: string) {
    const res: any = []
    for (const index in data) {
      data[index].key = `${data[index].name}_${env}`
      data[index].title = data[index].name
      data[index].isLeaf = true
      res.push(data[index])
    }
    return res
  }

  // 点击分组节点时，切换activatedNode
  // 当节点是环境名时，收缩或展开对应环境内的集群列表
  // 当节点是集群名时，右侧页面需要展示对应集群的策略列表
  activeNode (data: any): void {
    // console.log(this.eoNgTreeDefault?.getTreeNodeByKey(this.clusterKey))
    // 当选中的节点是集群名且与当前选中节点集群名不同时，取消当前节点的选中
    if (
      data.keys[0] &&
      this.clusterKey !== data.keys[0] &&
      this.eoNgTreeDefault?.getTreeNodeByKey(this.clusterKey)?.isSelected
    ) {
      // @ts-ignore
      this.eoNgTreeDefault.getTreeNodeByKey(this.clusterKey).isSelected = false
    }
    // 节点是集群名
    if (!data.node.origin.children || data.node.origin.children.length === 0) {
      this.clusterKey = data.keys[0]
      this.clusterName = data.node.origin.name
      this.activatedNode = data.node!
      this.showList = true
      this.showApiPage = false
      this.router.navigate(
        ['/', 'serv-governance', this.strategyType, 'group', 'list', data.node.origin.name]
      )
    } else {
      data.node.isExpanded = !data.node.isExpanded
    }
  }

  // 不同流量策略的面包屑
  getBreadcrumb () {
    switch (this.strategyType) {
      case 'traffic':
        this.appConfigService.reqFlashBreadcrumb([{ title: '流量策略' }])
        break
      case 'grey':
        this.appConfigService.reqFlashBreadcrumb([{ title: '灰度策略' }])
        break
      case 'fuse':
        this.appConfigService.reqFlashBreadcrumb([{ title: '熔断策略' }])
        break
      case 'cache':
        this.appConfigService.reqFlashBreadcrumb([{ title: '缓存策略' }])
        break
      case 'visit':
        this.appConfigService.reqFlashBreadcrumb([{ title: '访问策略' }])
        break
      default:
        this.appConfigService.reqFlashBreadcrumb([{ title: '流量策略' }])
    }
  }
}

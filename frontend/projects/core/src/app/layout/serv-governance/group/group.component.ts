/* eslint-disable dot-notation */
import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgTreeDefaultComponent } from 'eo-ng-tree'
import { NzTreeNodeOptions, NzTreeNode } from 'ng-zorro-antd/tree'
import { ClusterEnum } from '../../../constant/type'
import { ApiService } from '../../../service/api.service'
import { EoNgNavigationService } from '../../../service/eo-ng-navigation.service'
import { BaseInfoService } from '../../../service/base-info.service'

@Component({
  selector: 'eo-ng-flow-control-group',
  templateUrl: './group.component.html',
  styles: [
    `
  `]
})
export class GroupComponent implements OnInit {
  @ViewChild('addGroupRef', { read: TemplateRef, static: true }) addGroupRef:
    | TemplateRef<any>
    | string = ''

  @ViewChild('eoNgTreeDefault') eoNgTreeDefault!: EoNgTreeDefaultComponent

  public nodesList: NzTreeNodeOptions[] = []
  groupUuid: string = '' // 供右侧list页面用
  activatedNode: NzTreeNode | null = null
  clusterName: string = ''
  clusterKey: string = ''
  strategyType: string = ''
  constructor (
    private baseInfo:BaseInfoService,
    private api: ApiService,
    private router: Router,
    private navigationService:EoNgNavigationService
  ) {
    this.strategyType = this.router.url.split('/')[this.router.url.split('/').indexOf('serv-governance') + 1]
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
    this.api.get('cluster/enum').subscribe((resp: {code:number, data:{ envs:ClusterEnum[]}, msg:string}) => {
      if (resp.code === 0) {
        this.nodesList = []
        this.nodesList = this.nodesTransfer(resp.data.envs)
        this.getGroupItemSelected()
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
      this.eoNgTreeDefault.getTreeNodeByKey(this.clusterKey)!.isSelected = false
    }
    // 节点是集群名
    if (!data.node.origin.children || data.node.origin.children.length === 0) {
      this.clusterKey = data.keys[0]
      this.clusterName = data.node.origin.name
      this.activatedNode = data.node!
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
        this.navigationService.reqFlashBreadcrumb([{ title: '流量限制' }])
        break
      case 'grey':
        this.navigationService.reqFlashBreadcrumb([{ title: '灰度发布' }])
        break
      case 'fuse':
        this.navigationService.reqFlashBreadcrumb([{ title: '熔断策略' }])
        break
      case 'cache':
        this.navigationService.reqFlashBreadcrumb([{ title: '数据缓存' }])
        break
      case 'visit':
        this.navigationService.reqFlashBreadcrumb([{ title: 'API访问权限' }])
        break
      default:
        this.navigationService.reqFlashBreadcrumb([{ title: '流量限制' }])
    }
  }
}

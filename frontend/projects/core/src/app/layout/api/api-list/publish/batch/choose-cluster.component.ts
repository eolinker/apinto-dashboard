import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { CheckBoxOptionInterface } from 'eo-ng-checkbox'
import { ClusterSimpleOption } from 'projects/core/src/app/constant/type'
import { ApiService } from 'projects/core/src/app/service/api.service'

@Component({
  selector: 'eo-ng-api-batch-publish-choose-cluster',
  templateUrl: './choose-cluster.component.html',
  styles: [
    `label{
      height: 26px;
      line-height: 26px;
      margin-bottom: 16px;
      display: block;
      margin-top: 2px;
    }
  `
  ]
})
export class ApiBatchPublishChooseClusterComponent implements OnInit {
  // disabledCheckboxTpl是前置判断集群状态，v3.1中接口改造不传集群状态，将判断状态后置，所以废弃
  @ViewChild('disabledCheckboxTpl', { read: TemplateRef, static: true }) disabledCheckboxTpl: TemplateRef<any> | undefined
  clustersList:CheckBoxOptionInterface[]=[]
  clustersSet:Set<string> = new Set()
  type:'online'|'offline'|undefined
  loading:boolean = false
  apisSet:Set<string> = new Set()
  constructor (private api:ApiService) {}
  ngOnInit (): void {
    this.getClusterList()
  }

  getClusterList ():void {
    this.loading = true
    this.clustersList = []
    this.api.get('clusters/simple').subscribe((resp:{code:number, data:{clusters:ClusterSimpleOption[]}, msg:string}) => {
      if (resp.code === 0) {
        this.clustersList = []
        for (const index in resp.data.clusters) {
          this.clustersList = resp.data.clusters.map((cluster:ClusterSimpleOption) => {
            return { label: cluster.title, value: cluster.name, checked: this.clustersSet.has(cluster.name) }
          })
          if (this.clustersList.length === 1) {
            this.clustersSet.add(resp.data.clusters[index].name)
          }
        }
      }
      this.loading = false
    })
  }

  // 批量上下线的集群被选中或取消时，clustersSet随之变化
  changeClustersSet (value:CheckBoxOptionInterface[]) {
    for (const index in value) {
      if (value[index].checked && !value[index].disabled) {
        this.clustersSet.add(value[index].value)
      } else {
        this.clustersSet.delete(value[index].value)
      }
    }
  }
}

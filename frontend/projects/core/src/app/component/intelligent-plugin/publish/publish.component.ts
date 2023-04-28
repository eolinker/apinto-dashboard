import { Component, OnInit } from '@angular/core'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { IntelligentPluginService } from '../intelligent-plugin.service'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { ApiService } from '../../../service/api.service'
import { DynamicPublish, DynamicPublishCluster, DynamicPublishData } from '../types/types'

@Component({
  selector: 'eo-ng-intelligent-plugin-publish',
  templateUrl: './publish.component.html',
  styles: [
  ]
})
export class IntelligentPluginPublishComponent implements OnInit {
  name:string = ''
  id:string = ''
  desc:string = ''
  publishTableBody:TBODY_TYPE[] = [...this.service.createPluginTbody(this)]
  publishTableHeadName:THEAD_TYPE[] = [...this.service.createPluginThead(this)]
  publishList:DynamicPublishCluster[] = []
  selectedNum:number = 0
  selectedClusters:Array<string> = []
  moduleName:string = ''
  closeModal:any
  nzDisabled:boolean = false
  constructor (
    private message: EoNgFeedbackMessageService,
    private service:IntelligentPluginService,
    private api:ApiService) {}

  ngOnInit (): void {
    this.getPublishList()
  }

  getPublishList () {
    this.api.get(`dynamic/${this.moduleName}/cluster/${this.id}`).subscribe((resp:{code:number, msg:string, data:DynamicPublishData}) => {
      if (resp.code === 0) {
        this.publishList = resp.data.clusters
      }
    })
  }

  tableClick = (item:any) => {
    item.checked = !item.checked
    item.data.checked = !item.data.checked
    this.checkSelectedCluster()
  }

  // 点击表头全选
  checkAll () {
    this.checkSelectedCluster()
  }

  // 点击单条数据
  clickData () {
    this.checkSelectedCluster()
  }

  checkSelectedCluster () {
    setTimeout(() => {
      this.selectedClusters = this.publishList.filter((item:any) => {
        return item.checked
      }).map((item) => {
        return item.name
      })
      this.selectedNum = this.selectedClusters.length
      this.publishList = [...this.publishList] // 表头的勾选状态需要重载数据才能刷新
    }, 0

    )
  }

  offline () {
    const cluster:Array<string> = this.publishList.filter((item) => {
      return item.checked
    }).map((item) => {
      return item.name
    })
    this.api.put(`dynamic/${this.moduleName}/offline/${this.id}`, { cluster: cluster }).subscribe((resp:DynamicPublish) => {
      if (resp.code === 0) {
        this.message.success(resp.msg)
        this.closeModal && this.closeModal()
      }
    })
  }

  online () {
    const cluster:Array<string> = this.publishList.filter((item) => {
      return item.checked
    }).map((item) => {
      return item.name
    })
    this.api.put(`dynamic/${this.moduleName}/online/${this.id}`, { cluster: cluster }).subscribe((resp:DynamicPublish) => {
      if (resp.code === 0) {
        this.message.success(resp.msg)
        this.closeModal && this.closeModal()
      }
    })
  }
}

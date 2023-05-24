/* eslint-disable dot-notation */
import { Component, OnInit } from '@angular/core'
import { IntelligentPluginPublishComponent } from '../../../component/intelligent-plugin/publish/publish.component'
import { DynamicPublish, DynamicPublishCluster } from '../../../component/intelligent-plugin/types/types'

@Component({
  selector: 'eo-ng-application-publish',
  templateUrl: '../../../component/intelligent-plugin/publish/publish.component.html',
  styles: [
  ]
})
export class ApplicationPublishComponent extends IntelligentPluginPublishComponent implements OnInit {
  override getPublishList () {
    this.api.get('application/onlines', { uuid: this.id }).subscribe((resp:{code:number, msg:string, data:{info:{id:string, name:string, desc:string}, clusters:DynamicPublishCluster[]}}) => {
      if (resp.code === 0) {
        this.publishList = resp.data.clusters
        this.name = resp.data.info.name
        this.id = resp.data.info.id
        this.desc = resp.data.info.desc
      }
    })
  }

  // @ts-ignore
  override offline () {
    try {
      const cluster:Array<string> = this.publishList.filter((item) => {
        return item.checked
      }).map((item) => {
        return item.name
      })
      this.api.put('application/offline', { clusterNames: cluster }, { appId: this.id }).subscribe((resp:DynamicPublish) => {
        if (resp.code === 0) {
          this.message.success(resp.msg)
          return true
        } else {
          return false
        }
      })
    } catch {
      return false
    }
  }

  // @ts-ignore
  override online () {
    try {
      const cluster:Array<string> = this.publishList.filter((item) => {
        return item.checked
      }).map((item) => {
        return item.name
      })
      this.api.put('application/online', { clusterNames: cluster }, { appId: this.id }).subscribe((resp:DynamicPublish) => {
        if (resp.code === 0) {
          this.message.success(resp.msg)
          return true
        } else {
          return false
        }
      })
    } catch {
      return false
    }
  }
}

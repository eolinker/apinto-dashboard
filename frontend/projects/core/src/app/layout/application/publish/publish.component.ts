/* eslint-disable dot-notation */
import { ChangeDetectorRef, Component, OnInit } from '@angular/core'
import { DynamicPublish, DynamicPublishCluster } from '../../../component/intelligent-plugin/types/types'
import { Observable } from 'rxjs'
import { EoNgApplicationService } from '../application.service'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { ApiService } from '../../../service/api.service'
import { EoIntelligentPluginPublishComponent } from '../../../component/intelligent-plugin/publish/publish.component'
import { EoIntelligentPluginService } from '../../../component/intelligent-plugin/intelligent-plugin.service'

@Component({
  selector: 'eo-ng-application-publish',
  templateUrl: '../../../component/intelligent-plugin/publish/publish.component.html',
  styles: [
  ]
})
export class ApplicationPublishComponent extends EoIntelligentPluginPublishComponent implements OnInit {
  constructor (
    message: EoNgFeedbackMessageService,
    service:EoIntelligentPluginService,
    api:ApiService,
    public appService:EoNgApplicationService,
    cdref:ChangeDetectorRef) {
    super(message, service, api, cdref)
  }

  override ngAfterViewInit () {
    this.publishTableBody = [...this.appService.createApplicationPublicTbody(this)]
    this.cdref.detectChanges()
  }

  override getPublishList () {
    this.api.get('application/onlines', { uuid: this.id }).subscribe((resp:{code:number, msg:string, data:{info:{id:string, name:string, desc:string}, clusters:DynamicPublishCluster[]}}) => {
      if (resp.code === 0) {
        this.publishList = resp.data.clusters
        this.name = resp.data.info.name
        this.id = resp.data.info.id
        this.desc = resp.data.info.desc
      } else {
        this.returnToSdk && this.returnToSdk(resp)
      }
    })
  }

  // @ts-ignore
  override offline () {
    return new Observable((observer) => {
      const cluster:Array<string> = this.publishList.filter((item) => {
        return item.checked
      }).map((item) => {
        return item.name
      })
      this.showNoCluster = cluster.length <= 0
      if (this.showNoCluster) {
        observer.next(false)
        return
      }
      this.api.put('application/offline', { clusterNames: cluster }, { appId: this.id }).subscribe((resp:DynamicPublish) => {
        if (resp.code === 0) {
          this.message.success(resp.msg)
          observer.next(true)
        } else {
          observer.next(false)
        }
        this.returnToSdk && this.returnToSdk(resp)
      })
    })
  }

  // @ts-ignore
  override online () {
    return new Observable((observer) => {
      const cluster:Array<string> = this.publishList.filter((item) => {
        return item.checked
      }).map((item) => {
        return item.name
      })
      this.showNoCluster = cluster.length <= 0
      if (this.showNoCluster) {
        observer.next(false)
        return
      }
      this.api.put('application/online', { clusterNames: cluster }, { appId: this.id }).subscribe((resp:DynamicPublish) => {
        if (resp.code === 0) {
          this.message.success(resp.msg)
          observer.next(true)
        } else {
          observer.next(false)
        }
        this.returnToSdk && this.returnToSdk(resp)
      })
    })
  }
}

/* eslint-disable dot-notation */
import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { THEAD_TYPE } from 'eo-ng-table'
import { PublishTableHeadName, PublishTableBody } from 'projects/core/src/app/constant/conf'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { EoNgNavigationService } from 'projects/core/src/app/service/eo-ng-navigation.service'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { PublishFailService } from 'projects/core/src/app/service/publish-fail.service'
import { EO_TBODY_TYPE } from 'projects/eo-ng-apinto-table/src/public-api'

@Component({
  selector: 'eo-ng-service-discovery-publish',
  templateUrl: './publish.component.html',
  styles: [
  ]
})
export class ServiceDiscoveryPublishComponent implements OnInit {
  @ViewChild('clusterStatusTpl', { read: TemplateRef, static: true }) clusterStatusTpl: TemplateRef<any> | undefined
  readonly nowUrl:string = this.router.routerState.snapshot.url
  nzDisabled:boolean = false
  serviceName:string = ''
  errorMessageId:string = ''
  type:string = ''
  failmsg:string = ''
  solutionRouter:string = ''
  solutionParam:any = {}
  clustersList : Array<object> = []
  clustersTableHeadName: THEAD_TYPE[] = [...PublishTableHeadName]
  clustersTableBody: EO_TBODY_TYPE[] = [...PublishTableBody]

  constructor (
    private baseInfo:BaseInfoService,
     private message: EoNgFeedbackMessageService,
      public api:ApiService,
       private router:Router,
       private navigationService:EoNgNavigationService,
       private publishFailModal:PublishFailService) {
    this.navigationService.reqFlashBreadcrumb([{ title: '服务发现', routerLink: 'upstream/discovery' }, { title: '上线管理' }])
  }

  ngOnInit (): void {
    this.serviceName = this.baseInfo.allParamsInfo.discoveryName
    if (!this.serviceName) {
      this.router.navigate(['/'])
    }
    this.getClustersData()
  }

  ngAfterViewInit () {
    this.clustersTableBody[2].title = this.clusterStatusTpl
    this.clustersTableBody[5].btns[0].disabledFn = () => {
      return this.nzDisabled
    }
    this.clustersTableBody[5].btns[0].click = (item:any) => {
      this.updateOrOnline(item.data, '更新')
    }

    this.clustersTableBody[5].btns[1].disabledFn = () => {
      return this.nzDisabled
    }
    this.clustersTableBody[5].btns[1].click = (item:any) => {
      this.offline(item.data)
    }

    this.clustersTableBody[6].btns[0].disabledFn = () => {
      return this.nzDisabled
    }
    this.clustersTableBody[6].btns[0].click = (item:any) => {
      this.offline(item.data)
    }

    this.clustersTableBody[7].btns[0].disabledFn = () => {
      return this.nzDisabled
    }
    this.clustersTableBody[7].btns[0].click = (item:any) => {
      this.updateOrOnline(item.data, '上线')
    }
  }

  getClustersData () {
    this.api.get('discovery/' + this.serviceName + '/onlines').subscribe(resp => {
      if (resp.code === 0) {
        this.clustersList = resp.data.clusters
      }
    })
  }

  // 更新或上线
  updateOrOnline (item:any, type:string) {
    this.solutionRouter = ''
    this.solutionParam = {}
    this.api.put('discovery/' + this.serviceName + '/online', { clusterName: (item.name || '') }).subscribe(resp => {
      if (resp.code === 0) {
        this.message.success(resp.msg || (type + '成功'), { nzDuration: 1000 })
        this.getClustersData()
      } else {
        if (this.errorMessageId) {
          this.message.remove(this.errorMessageId)
          this.errorMessageId = ''
        }
        this.type = type
        this.failmsg = resp.msg
        this.solutionRouter = resp.data?.router?.name ? resp.data.router.name : ''
        this.solutionParam = resp.data?.router?.params ? resp.data.router.params : {}
        // this.errorMessageId = this.message.info(this.error || '获取数据失败!', { nzDuration: 0 }).messageId
        if (this.solutionRouter) {
          this.publishFailModal.openModal(resp.msg, '服务', this.solutionRouter, this.solutionParam)
        } else {
          this.message.error(resp.msg || (type + '失败'))
        }
      }
    })
  }

  // 上线失败的全局提示为强提醒，需要点击才可关闭
  closeErrorBtn () {
    this.message.remove(this.errorMessageId)
    this.errorMessageId = ''
  }

  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  // 下线
  offline (item:any) {
    this.solutionRouter = ''
    this.solutionParam = {}

    this.api.put('discovery/' + this.serviceName + '/offline', { clusterName: item.name || '' }).subscribe(resp => {
      if (resp.code === 0) {
        this.getClustersData()
        this.message.success(resp.msg || '下线成功', { nzDuration: 1000 })
      } else {
        if (this.errorMessageId) {
          this.message.remove(this.errorMessageId)
          this.errorMessageId = ''
        }
        this.type = '下线'
        this.failmsg = resp.msg
        // this.errorMessageId = this.message.info(this.error || '下线失败!', { nzDuration: 0 }).messageId
        this.message.error(resp.msg || '下线失败')
      }
    })
  }
}

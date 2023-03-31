/* eslint-disable dot-notation */
/* eslint-disable no-useless-constructor */
import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { THEAD_TYPE } from 'eo-ng-table'
import { PublishTableBody, PublishTableHeadName } from 'projects/core/src/app/constant/conf'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { PublishFailService } from 'projects/core/src/app/service/publish-fail.service'
import { EO_TBODY_TYPE } from 'projects/eo-ng-apinto-table/src/public-api'

@Component({
  selector: 'eo-ng-upstream-publish',
  templateUrl: './publish.component.html',
  styles: [
  ]
})
export class UpstreamPublishComponent implements OnInit {
  @ViewChild('clusterStatusTpl', { read: TemplateRef, static: true }) clusterStatusTpl: TemplateRef<any> | undefined
  readonly nowUrl:string = this.router.routerState.snapshot.url
  serviceName:string = ''

  errorMessageId:string = ''
  type:string = ''
  failmsg:string = ''
  solutionRouter:string = ''
  solutionParam:any = {}

  nzDisabled:boolean = false
  clustersList : Array<object> = []
  clustersTableHeadName: THEAD_TYPE[] = [...PublishTableHeadName]
  clustersTableBody: EO_TBODY_TYPE[] =[...PublishTableBody]

  constructor (
    private baseInfo:BaseInfoService,
    private message: EoNgFeedbackMessageService,
    public api:ApiService,
    private router:Router,
    private appConfigService:AppConfigService,
    private publishFailModal:PublishFailService) {
    this.appConfigService.reqFlashBreadcrumb([{ title: '上游管理', routerLink: 'upstream/upstream' }, { title: '上线管理' }])
  }

  ngOnInit (): void {
    this.serviceName = this.baseInfo.allParamsInfo.serviceName
    if (!this.serviceName) {
      this.router.navigate(['/', 'upstream', 'upstream'])
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
    this.api.get('service/' + this.serviceName + '/onlines').subscribe((resp: { code: number; data: { clusters: object[] }; msg: any }) => {
      if (resp.code === 0) {
        this.clustersList = resp.data.clusters
      }
    })
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  updateOrOnline (item:any, type:string) {
    this.solutionRouter = ''
    this.solutionParam = {}

    this.api.put('service/' + this.serviceName + '/online', { clusterName: item.name || '' }).subscribe((resp:any) => {
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
        if (this.solutionRouter) {
          this.publishFailModal.openModal(resp.msg, '上游服务', this.solutionRouter, this.solutionParam)
        } else {
          this.message.error(resp.msg || (type + '失败'))
        }
      }
    })
  }

  closeErrorBtn () {
    this.message.remove(this.errorMessageId)
    this.errorMessageId = ''
  }

  offline (item:any) {
    this.solutionRouter = ''
    this.solutionParam = {}

    this.api.put('service/' + this.serviceName + '/offline', { clusterName: item.name || '' }).subscribe((resp:any) => {
      if (resp.code === 0) {
        this.message.success(resp.msg || '下线成功!', { nzDuration: 1000 })
        this.getClustersData()
      } else {
        if (this.errorMessageId) {
          this.message.remove(this.errorMessageId)
          this.errorMessageId = ''
        }
        this.type = '下线'
        this.failmsg = resp.msg
        this.message.error(resp.msg)
      }
    })
  }
}

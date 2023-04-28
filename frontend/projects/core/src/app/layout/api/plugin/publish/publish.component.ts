import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { PublishTableHeadName, PublishTableBody } from 'projects/core/src/app/constant/conf'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { EoNgNavigationService } from 'projects/core/src/app/service/eo-ng-navigation.service'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { PublishFailService } from 'projects/core/src/app/service/publish-fail.service'

@Component({
  selector: 'eo-ng-api-plugin-template-publish',
  templateUrl: './publish.component.html',
  styles: [
  ]
})
export class ApiPluginTemplatePublishComponent implements OnInit {
  @ViewChild('clusterStatusTpl', { read: TemplateRef, static: true }) clusterStatusTpl: TemplateRef<any> | undefined
  readonly nowUrl:string = this.router.routerState.snapshot.url
  uuid:string = ''

  errorMessageId:string = ''
  type:string = ''
  failmsg:string = ''
  solutionRouter:string = ''
  solutionParam:any = {}

  nzDisabled:boolean = false
  clustersList : Array<object> = []
  clustersTableHeadName: Array<object> = [...PublishTableHeadName]
  clustersTableBody: Array<any> =[...PublishTableBody]

  constructor (
    private baseInfo:BaseInfoService,
    private message: EoNgFeedbackMessageService,
    public api:ApiService,
    private router:Router,
    private navigationService:EoNgNavigationService,
    private publishFailModal:PublishFailService) {
    this.navigationService.reqFlashBreadcrumb([{ title: '插件模板', routerLink: 'router/plugin-template' }, { title: '上线管理' }])
  }

  ngOnInit (): void {
    this.uuid = this.baseInfo.allParamsInfo.pluginTemplateId
    if (!this.uuid) {
      this.router.navigate(['/', 'router', 'plugin-template'])
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
    this.api.get('plugin/template/onlines', { uuid: this.uuid }).subscribe((resp: { code: number; data: { clusters: object[] }; msg: any }) => {
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

    this.api.put('plugin/template/online', { clusterName: item.name || '' }, { uuid: this.uuid }).subscribe((resp:any) => {
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
          this.publishFailModal.openModal(resp.msg, '插件模板', this.solutionRouter, this.solutionParam)
        } else {
          this.message.error(resp.msg || (type + '失败'))
        }
      }
    })
  }

  offline (item:any) {
    this.solutionRouter = ''
    this.solutionParam = {}

    this.api.put('plugin/template/offline', { clusterName: item.name || '' }, { uuid: this.uuid }).subscribe((resp:any) => {
      if (resp.code === 0) {
        this.message.success(resp.msg || '下线成功!', { nzDuration: 1000 })
        this.getClustersData()
      } else {
        this.type = '下线'
        this.failmsg = resp.msg
      }
    })
  }
}

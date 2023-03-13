/* eslint-disable dot-notation */
/* eslint-disable no-useless-constructor */
import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { Router, ActivatedRoute } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { PublishFailService } from 'projects/core/src/app/service/publish-fail.service'

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
  clustersTableHeadName: Array<object> = [
    { title: '集群名称' },
    { title: '环境' },
    { title: '状态' },
    { title: '更新者' },
    { title: '更新时间' },
    {
      title: '操作',
      right: true
    }
  ]

  clustersTableBody: Array<any> =[
    { key: 'name' },
    { key: 'env' },
    { key: 'status' },
    { key: 'operator' },
    { key: 'update_time' },
    {
      type: 'btn',
      right: true,
      showFn: (item:any) => {
        return item.status === 'TOUPDATE'
      },
      btns: [
        {
          title: '更新',
          disabledFn: () => {
            return this.nzDisabled
          },
          click: (item:any) => {
            this.updateOrOnline(item.data, '更新')
          }
        },
        {
          title: '下线',
          disabledFn: () => {
            return this.nzDisabled
          },
          click: (item:any) => {
            this.offline(item.data)
          }
        }
      ]
    },
    {
      type: 'btn',
      right: true,
      showFn: (item:any) => {
        return item.status === 'GOONLINE'
      },
      btns: [
        {
          title: '下线',
          disabledFn: () => {
            return this.nzDisabled
          },
          click: (item:any) => {
            this.offline(item.data)
          }
        }
      ]
    },
    {
      type: 'btn',
      right: true,
      showFn: (item:any) => {
        return (item.status === 'OFFLINE' || item.status === 'NOTGOONLINE')
      },
      btns: [
        {
          title: '上线',
          disabledFn: () => {
            return this.nzDisabled
          },
          click: (item:any) => {
            this.updateOrOnline(item.data, '上线')
          }
        }
      ]
    }
  ]

  constructor (
    private baseInfo:BaseInfoService,
    private message: EoNgFeedbackMessageService,
    public api:ApiService,
    private router:Router,
      private activateInfo:ActivatedRoute,
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
  }

  getClustersData () {
    this.api.get('service/' + this.serviceName + '/onlines').subscribe((resp: { code: number; data: { clusters: object[] }; msg: any }) => {
      if (resp.code === 0) {
        this.clustersList = resp.data.clusters
      } else {
        this.message.error(resp.msg || '获取数据失败!')
      }
    })
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  updateOrOnline (item:any, type:string) {
    this.solutionRouter = ''
    this.solutionParam = {}

    this.api.put('service/' + this.serviceName + '/online', { cluster_name: item.name || '' }).subscribe((resp:any) => {
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

    this.api.put('service/' + this.serviceName + '/offline', { cluster_name: item.name || '' }).subscribe((resp:any) => {
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

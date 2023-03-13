/* eslint-disable dot-notation */
/* eslint-disable no-undef */
import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { Router, ActivatedRoute } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'
import { BaseInfoService } from '../../../service/base-info.service'
import { PublishFailService } from '../../../service/publish-fail.service'

@Component({
  selector: 'eo-ng-application-publish',
  templateUrl: './publish.component.html',
  styles: [
  ]
})
export class ApplicationPublishComponent implements OnInit {
  @ViewChild('clusterStatusTpl', { read: TemplateRef, static: true }) clusterStatusTpl: TemplateRef<any> | undefined
  @ViewChild('disabledStatusTpl', { read: TemplateRef, static: true }) disabledStatusTpl: TemplateRef<any> | undefined

  readonly nowUrl:string = this.router.routerState.snapshot.url
  appId:string = ''

  clustersList : Array<object> = []
  clustersTableHeadName: Array<object> = [
    { title: '集群名称' },
    { title: '环境' },
    { title: '状态' },
    { title: '禁用状态' },
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
    { key: 'disable' },
    { key: 'operator' },
    { key: 'update_time' },
    {
      type: 'btn',
      right: true,
      showFn: (item:any) => {
        return item.status === 'TOUPDATE' && !item.disable
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
          click: (item:any) => {
            this.offline(item.data)
          },
          disabledFn: () => {
            return this.nzDisabled
          }
        },
        {
          title: '禁用',
          click: (item:any) => {
            this.disable(item.data)
          },
          disabledFn: () => {
            return this.nzDisabled
          }
        }
      ]
    },
    {
      type: 'btn',
      right: true,
      showFn: (item:any) => {
        return item.status === 'TOUPDATE' && item.disable
      },
      btns: [
        {
          title: '更新',
          click: (item:any) => {
            this.updateOrOnline(item.data, '更新')
          },
          disabledFn: () => {
            return this.nzDisabled
          }
        },
        {
          title: '下线',
          click: (item:any) => {
            this.offline(item.data)
          },
          disabledFn: () => {
            return this.nzDisabled
          }
        },
        {
          title: '启用',
          click: (item:any) => {
            this.enable(item.data)
          },
          disabledFn: () => {
            return this.nzDisabled
          }
        }
      ]
    },
    {
      type: 'btn',
      right: true,
      showFn: (item:any) => {
        return item.status === 'GOONLINE' && !item.disable
      },
      btns: [
        {
          title: '下线',
          click: (item:any) => {
            this.offline(item.data)
          },
          disabledFn: () => {
            return this.nzDisabled
          }
        },
        {
          title: '禁用',
          click: (item:any) => {
            this.disable(item.data)
          },
          disabledFn: () => {
            return this.nzDisabled
          }
        }
      ]
    },
    {
      type: 'btn',
      right: true,
      showFn: (item:any) => {
        return item.status === 'GOONLINE' && item.disable
      },
      btns: [
        {
          title: '下线',
          click: (item:any) => {
            this.offline(item.data)
          },
          disabledFn: () => {
            return this.nzDisabled
          }
        },
        {
          title: '启用',
          click: (item:any) => {
            this.enable(item.data)
          },
          disabledFn: () => {
            return this.nzDisabled
          }
        }
      ]
    },
    {
      type: 'btn',
      right: true,
      showFn: (item:any) => {
        return (item.status === 'OFFLINE' || item.status === 'NOTGOONLINE') && !item.disable
      },
      btns: [
        {
          title: '上线',
          click: (item:any) => {
            this.updateOrOnline(item.data, '上线')
          },
          disabledFn: () => {
            return this.nzDisabled
          }
        },
        {
          title: '禁用',
          click: (item:any) => {
            this.disable(item.data)
          },
          disabledFn: () => {
            return this.nzDisabled
          }
        }
      ]
    },
    {
      type: 'btn',
      right: true,
      showFn: (item:any) => {
        return (item.status === 'OFFLINE' || item.status === 'NOTGOONLINE') && item.disable
      },
      btns: [
        {
          title: '上线',
          click: (item:any) => {
            this.updateOrOnline(item.data, '上线')
          },
          disabledFn: () => {
            return this.nzDisabled
          }
        },
        {
          title: '启用',
          click: (item:any) => {
            this.enable(item.data)
          },
          disabledFn: () => {
            return this.nzDisabled
          }
        }
      ]
    }
  ]

  errorMessageId:string = ''
  type:string = ''
  failmsg:string = ''
  solutionRouter:string = ''
  solutionParam:any = {}
  nzDisabled:boolean = false

  constructor (
                private message: EoNgFeedbackMessageService,
                public api:ApiService,
                private baseInfo:BaseInfoService,
                private router:Router,
                private activateInfo:ActivatedRoute,
                private appConfigService:AppConfigService,
                private publishFailModal:PublishFailService) {
    this.appConfigService.reqFlashBreadcrumb([{ title: '应用管理', routerLink: 'application' }, { title: '上线管理' }])
  }

  ngOnInit (): void {
    this.appId = this.baseInfo.allParamsInfo.appId
    if (!this.appId) {
      this.router.navigate(['/', 'application'])
    }
    this.getClustersData()
  }

  ngAfterViewInit () {
    this.clustersTableBody[2].title = this.clusterStatusTpl
    this.clustersTableBody[3].title = this.disabledStatusTpl
  }

  getClustersData () {
    this.api.get('application/onlines', { app_id: this.appId }).subscribe(resp => {
      if (resp.code === 0) {
        this.clustersList = resp.data.clusters
      } else {
        this.message.error(resp.msg || '刷新列表失败!')
      }
    })
  }

  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  updateOrOnline (item:any, type:string) {
    this.solutionRouter = ''
    this.solutionParam = {}
    this.api.put('application/online', { cluster_name: item.name || '' }, { app_id: this.appId }).subscribe(resp => {
      if (resp.code === 0) {
        this.message.success(resp.msg || (type + '成功'), { nzDuration: 1000 })
        this.getClustersData()
      } else {
        if (this.errorMessageId) {
          this.message.remove(this.errorMessageId)
          this.errorMessageId = ''
        }
        this.failmsg = resp.msg
        this.solutionRouter = resp.data?.router?.name ? resp.data.router.name : ''
        this.solutionParam = resp.data?.router?.params ? resp.data.router.params : {}
        if (this.solutionRouter) {
          this.publishFailModal.openModal(resp.msg, '应用', this.solutionRouter, this.solutionParam)
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
    this.solutionParam = []

    this.api.put('application/offline', { cluster_name: (item.name || '') }, { app_id: this.appId }).subscribe(resp => {
      if (resp.code === 0) {
        this.message.success(resp.msg || '下线成功', { nzDuration: 1000 })
        this.getClustersData()
      } else {
        if (this.errorMessageId) {
          this.message.remove(this.errorMessageId)
          this.errorMessageId = ''
        }
        this.failmsg = resp.msg
        this.message.error(resp.msg || '下线失败')
      }
    })
  }

  enable (item:any) {
    this.solutionRouter = ''
    this.solutionParam = []

    this.api.put('application/enable', { cluster_name: (item.name || '') }, { app_id: this.appId }).subscribe(resp => {
      if (resp.code === 0) {
        this.message.success(resp.msg || '启用成功', { nzDuration: 1000 })
        this.getClustersData()
      } else {
        if (this.errorMessageId) {
          this.message.remove(this.errorMessageId)
          this.errorMessageId = ''
        }
        this.failmsg = resp.msg
        this.message.error(resp.msg || '启用失败')
      }
    })
  }

  disable (item:any) {
    this.solutionRouter = ''
    this.solutionParam = []

    this.api.put('application/disable', { cluster_name: item.name || '' }, { app_id: this.appId }).subscribe(resp => {
      if (resp.code === 0) {
        this.message.success(resp.msg || '停用成功', { nzDuration: 1000 })
        this.getClustersData()
      } else {
        if (this.errorMessageId) {
          this.message.remove(this.errorMessageId)
          this.errorMessageId = ''
        }
        this.failmsg = resp.msg
        this.message.error(resp.msg || '停用失败')
      }
    })
  }
}

/* eslint-disable dot-notation */
/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2022-08-17 23:42:52
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2022-09-20 23:15:51
 * @FilePath: /apinto/src/app/layout/api/api-publish/api-publish.component.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { Router, ActivatedRoute } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { PublishFailService } from 'projects/core/src/app/service/publish-fail.service'
import { Subscription } from 'rxjs'

@Component({
  selector: 'eo-ng-api-publish',
  templateUrl: './publish.component.html',
  styles: [
  ]
})
export class ApiPublishComponent implements OnInit {
  // @ViewChild('error', { static: true }) error!: TemplateRef<void>
  @ViewChild('clusterStatusTpl', { read: TemplateRef, static: true }) clusterStatusTpl: TemplateRef<any> | undefined
  @ViewChild('disabledStatusTpl', { read: TemplateRef, static: true }) disabledStatusTpl: TemplateRef<any> | undefined
  readonly nowUrl:string = this.router.routerState.snapshot.url
  apiUuid:string = ''
  private subscription: Subscription = new Subscription()
  nzDisabled:boolean = false
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
          disabledFn: () => {
            return this.nzDisabled
          },
          click: (item:any) => {
            this.offline(item.data)
          }
        },
        {
          title: '禁用',
          disabledFn: () => {
            return this.nzDisabled
          },
          click: (item:any) => {
            this.disable(item.data)
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
        },
        {
          title: '启用',
          disabledFn: () => {
            return this.nzDisabled
          },
          click: (item:any) => {
            this.enable(item.data)
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
          disabledFn: () => {
            return this.nzDisabled
          },
          click: (item:any) => {
            this.offline(item.data)
          }
        },
        {
          title: '禁用',
          disabledFn: () => {
            return this.nzDisabled
          },
          click: (item:any) => {
            this.disable(item.data)
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
          disabledFn: () => {
            return this.nzDisabled
          },
          click: (item:any) => {
            this.offline(item.data)
          }
        },
        {
          title: '启用',
          disabledFn: () => {
            return this.nzDisabled
          },
          click: (item:any) => {
            this.enable(item.data)
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
          disabledFn: () => {
            return this.nzDisabled
          },
          click: (item:any) => {
            this.updateOrOnline(item.data, '上线')
          }
        },
        {
          title: '禁用',
          disabledFn: () => {
            return this.nzDisabled
          },
          click: (item:any) => {
            this.disable(item.data)
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
          disabledFn: () => {
            return this.nzDisabled
          },
          click: (item:any) => {
            this.updateOrOnline(item.data, '上线')
          }
        },
        {
          title: '启用',
          disabledFn: () => {
            return this.nzDisabled
          },
          click: (item:any) => {
            this.enable(item.data)
          }
        }
      ]
    }
  ]

  constructor (private message: EoNgFeedbackMessageService,
     public api:ApiService,
      private router:Router,
    private publishFailModal:PublishFailService,
    private baseInfo:BaseInfoService,
      private activateInfo:ActivatedRoute,
      private appConfigService:AppConfigService) {
    this.apiUuid = this.baseInfo.allParamsInfo.apiId
    if (!this.apiUuid) {
      this.router.navigate(['/', 'router', 'group'])
    }
    this.appConfigService.reqFlashBreadcrumb([{ title: 'API管理', routerLink: 'router/group/list' }, { title: '上线管理' }])
  }

  ngOnInit (): void {
    this.getClustersData()
  }

  ngAfterViewInit () {
    this.clustersTableBody[2].title = this.clusterStatusTpl
    this.clustersTableBody[3].title = this.disabledStatusTpl
  }

  ngOnDestroy () {
    this.subscription.unsubscribe()
  }

  getClustersData () {
    this.api.get('router/onlines', { uuid: this.apiUuid }).subscribe(resp => {
      if (resp.code === 0) {
        this.clustersList = resp.data.clusters
      } else {
        this.message.error(resp.msg || '获取列表数据失败!')
      }
    })
  }

  errorMessageId:string = ''
  type:string = ''
  failmsg:string = ''
  solutionRouter:string = ''
  solutionParam:any = {}

  updateOrOnline (item:any, type:string) {
    this.solutionRouter = ''
    this.solutionParam = {}
    this.api.put('router/online', { cluster_name: item.name || '' }, { uuid: this.apiUuid }).subscribe(resp => {
      if (resp.code === 0) {
        this.message.success(resp.msg || (type + '成功'), { nzDuration: 1000 })
        this.getClustersData()
      } else {
        this.type = type
        this.solutionRouter = resp.data?.router?.name ? resp.data.router.name : ''
        this.solutionParam = resp.data?.router?.params ? resp.data.router.params : {}
        if (this.solutionRouter) {
          this.publishFailModal.openModal(resp.msg, 'API', this.solutionRouter, this.solutionParam)
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

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  offline (item:any) {
    this.solutionRouter = ''
    this.solutionParam = {}

    this.api.put('router/offline', { cluster_name: (item.name || '') }, { uuid: this.apiUuid }).subscribe(resp => {
      if (resp.code === 0) {
        this.message.success(resp.msg || '下线成功', { nzDuration: 1000 })
        this.getClustersData()
      } else {
        if (this.errorMessageId) {
          this.message.remove(this.errorMessageId)
          this.errorMessageId = ''
        }
        this.type = '下线'
        this.failmsg = resp.msg
        this.message.error(resp.msg || '下线失败')
      }
    })
  }

  enable (item:any) {
    this.solutionRouter = ''
    this.solutionParam = {}

    this.api.put('router/enable', { cluster_name: (item.name || '') }, { uuid: this.apiUuid }).subscribe(resp => {
      if (resp.code === 0) {
        this.message.success(resp.msg || '启用成功', { nzDuration: 1000 })
        this.getClustersData()
      } else {
        if (this.errorMessageId) {
          this.message.remove(this.errorMessageId)
          this.errorMessageId = ''
        }
        // this.type = '下线'
        this.failmsg = resp.msg
        // this.errorMessageId = this.message.info(this.error || '获取列表数据失败!', { nzDuration: 0 }).messageId
        this.message.error(resp.msg || '启用失败')
      }
    })
  }

  disable (item:any) {
    this.solutionRouter = ''
    this.solutionParam = {}

    this.api.put('router/disable', { cluster_name: (item.name || '') }, { uuid: this.apiUuid }).subscribe(resp => {
      if (resp.code === 0) {
        this.message.success(resp.msg || '禁用成功', { nzDuration: 1000 })
        this.getClustersData()
      } else {
        if (this.errorMessageId) {
          this.message.remove(this.errorMessageId)
          this.errorMessageId = ''
        }
        this.message.error(resp.msg || '禁用失败')
      }
    })
  }
}

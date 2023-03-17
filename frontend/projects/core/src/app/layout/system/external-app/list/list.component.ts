/* eslint-disable camelcase */
/* eslint-disable no-useless-constructor */
import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { MODAL_SMALL_SIZE } from 'projects/core/src/app/constant/app.config'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'

@Component({
  selector: 'eo-ng-external-app-list',
  templateUrl: './list.component.html',
  styles: [
  ]
})
export class ExternalAppListComponent implements OnInit {
  @ViewChild('switchTpl', { read: TemplateRef, static: true }) switchTpl: TemplateRef<any> | undefined
  @ViewChild('copyTokenTpl', { read: TemplateRef, static: true }) copyTokenTpl: TemplateRef<any> | undefined
  nzDisabled:boolean = false
  appsTableHeadName:Array<any> = [
    {
      title: '应用名称'
    },
    {
      title: '应用ID'
    },
    {
      title: '鉴权Token'
    },
    {
      title: '关联标签'
    },
    {
      title: '禁用状态',
      width: 84
    },
    {
      title: '更新者'
    },
    {
      title: '更新时间',
      showSort: true
    },
    {
      title: '操作',
      right: true
    }
  ]

  appsTableBody:Array<any> = [
    {
      key: 'name'
    },
    {
      key: 'id'
    },
    {
      key: 'token'
    },
    {
      key: 'tags'
    },
    {
      key: 'status'
    },
    {
      key: 'operator'
    },
    {
      key: 'update_time'
    },
    {
      type: 'btn',
      right: true,
      btns: [{
        title: '更新鉴权Token',
        disabledFn: () => { return this.nzDisabled },
        click: (item:any) => {
          this.updateToken(item.data)
        }
      },
      {
        title: '复制Token'
      },
      {
        title: '查看',
        click: (item:any) => {
          this.getAppMessage(item.data)
        }
      },
      {
        title: '删除',
        disabledFn: () => { return this.nzDisabled },
        click: (item:any) => {
          this.deleteAppModal(item.data)
        }
      }
      ]
    }
  ]

  appsList:Array<any> = []

  constructor (private router:Router, private message: EoNgFeedbackMessageService, private modalService:EoNgFeedbackModalService, private api:ApiService,
    private appConfigService: AppConfigService) {
    this.appConfigService.reqFlashBreadcrumb([
      { title: '外部应用' }
    ])
  }

  ngOnInit (): void {
    this.getAppsData()
  }

  ngAfterViewInit ():void {
    this.appsTableBody[4].title = this.switchTpl
    this.appsTableBody[7].btns[1].type = this.copyTokenTpl
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  getAppsData ():void {
    this.api.get('external-apps').subscribe((resp:{code:number, data:{apps:Array<{id:string, name:string, token:string, tags:string, status:number, operator:string, update_time:string}>}, msg:string}) => {
      if (resp.code === 0) {
        this.appsList = resp.data.apps
        for (const index in this.appsList) {
          this.appsList[index].status_boolean = this.appsList[index].status === 2 // 禁用
        }
      } else {
        this.message.error(resp.msg || '获取API列表数据失败！')
      }
    })
  }

  getAppMessage (item:any) {
    this.router.navigate(['/', 'system', 'ext-app', 'message', item.id])
  }

  addApp () {
    this.router.navigate(['/', 'system', 'ext-app', 'create'])
  }

  updateToken (item:any) {
    this.api.put('external-app/token', null, { id: item.id }).subscribe((resp:any) => {
      if (resp.code === 0) {
        this.message.success(resp.msg || '刷新鉴权Token成功！')
        this.getAppsData()
      } else {
        this.message.error(resp.msg || '刷新鉴权Token失败，请重试！')
      }
    })
  }

  appsTableClick = (item:any) => {
    this.getAppMessage(item.data)
  }

  disabledApp (item:any, e:Event) {
    e.stopPropagation()
    let url = ''
    if (!item.status_boolean) {
      url = 'external-app/disable'
    } else {
      url = 'external-app/enable'
    }

    this.api.put(url, null, { id: item.id }).subscribe((resp:{code:number, data:{}, msg:string}) => {
      if (resp.code === 0) {
        item.status_boolean = !item.status_boolean
        this.message.success(resp.msg || ((url.split('/')[1] === 'disable' ? '禁用' : '启用') + '成功！'))
        this.getAppsData()
      } else {
        this.message.error(resp.msg || ((url.split('/')[1] === 'disable' ? '禁用' : '启用') + '失败！'))
      }
    })
  }

  // 删除api弹窗
  deleteAppModal (item:any) {
    this.modalService.create({
      nzTitle: '删除',
      nzContent: '该数据删除后将无法找回，请确认是否删除？',
      nzClosable: true,
      nzClassName: 'delete-modal',
      nzWidth: MODAL_SMALL_SIZE,
      nzOkDanger: true,
      nzOnOk: () => {
        this.deleteApp(item)
      }
    })
  }

  // 删除单个api
  deleteApp = (items:any) => {
    this.api.delete('external-app', { id: items.id }).subscribe((resp:any) => {
      if (resp.code === 0) {
        this.message.success(resp.msg || '删除成功!')
        this.getAppsData()
      } else {
        this.message.error(resp.msg || '删除失败!')
      }
    })
  }

  callback = () => {
    this.message.success('复制成功')
  };
}

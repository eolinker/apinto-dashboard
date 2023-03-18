/* eslint-disable no-array-constructor */
/* eslint-disable camelcase */
/* eslint-disable dot-notation */
/* eslint-disable no-useless-constructor */
import { Component, Input, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { Router, NavigationEnd } from '@angular/router'
import { NzTreeNodeOptions } from 'ng-zorro-antd/tree'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'
import { Subscription } from 'rxjs'
import { MODAL_SMALL_SIZE } from '../../../constant/app.config'
import { BaseInfoService } from '../../../service/base-info.service'

@Component({
  selector: 'eo-ng-api-management-list',
  templateUrl: './list.component.html',
  styles: [
    `
    div eo-ng-api-method-tag{
      margin: 0 2px;
    }
    div eo-ng-api-method-tag:first-child {
      margin-left: 0;
    }
    div eo-ng-api-method-tag:last-child {
      margin-right: 0;
    }
    .ml107{
      margin-top:4px;
    }
    `
  ]
})
export class ApiManagementListComponent implements OnInit {
  @ViewChild('methodTpl', { read: TemplateRef, static: true }) methodTpl: TemplateRef<any> | undefined
  @Input() groupUuid:string = ''
  public nodesList:NzTreeNodeOptions[] = []
  public apiNodesMap:Map<string, any> = new Map()
  public apiNodesList:Array<any> = []
  apisSet:Set<string> = new Set()
  nzDisabled:boolean = false
  sourcesList:Array<{text:string, value:any, [key:string]:any}> = []
  apiNameForSear:string = ''

  apisForm: {apis:Array<{checked:boolean, group_uuid:string, uuid:string, name:string, method:string, service:string, request_path:string, update_time:string, is_delete:boolean}>, total:number, page_num:number, page_size:number, group_uuid:string, source_ids:Array<string>} = {
    apis: [],
    total: 0,
    page_num: 1,
    page_size: 20,
    group_uuid: '',
    source_ids: []
  }

  apisTableHeadName:Array<any> = [
    {
      type: 'checkbox',
      click: (item:any) => {
        this.changeApisSet(item, 'all')
      },
      showFn: () => {
        return !this.nzDisabled
      },
      resizeable: false
    },
    {
      title: 'API名称'
    },
    {
      title: '协议/方法',
      width: 140,
      resizeable: false
    },
    {
      title: '上游服务名称'
    },
    {
      title: '请求路径'
    },
    {
      title: '来源',
      filterMultiple: true, // 待改
      filterOpts: [{
        text: '自建',
        value: 'build'
      },
      {
        text: '导入',
        value: 'import'
      }
      ],
      filterFn: () => {
        return true
      }
    },
    {
      title: '更新时间'
    },
    {
      title: '操作',
      right: true
    }
  ]

  apisTableBody:Array<any> = [
    {
      key: 'checked',
      type: 'checkbox',
      click: (item:any) => {
        this.changeApisSet(item)
      },
      showFn: () => {
        return !this.nzDisabled
      }
    },
    {
      key: 'name'
    },
    {
      key: 'method'
    },
    {
      key: 'service'
    },
    {
      key: 'request_path'
    },
    {
      key: 'source'
    },
    {
      key: 'update_time'
    },
    {
      type: 'btn',
      right: true,
      btns: [{
        title: '上线管理',
        click: (item:any) => {
          this.router.navigate(['/', 'router', 'content', item.data.uuid, 'publish'])
        }
      },
      {
        title: '查看',
        click: (item:any) => {
          this.router.navigate(['/', 'router', 'content', item.data.uuid])
        }
      },
      {
        title: '删除',
        disabledFn: (data:any, item:any) => { return !item.data.is_delete || this.nzDisabled },
        click: (item:any) => {
          this.deleteApiModal(item.data)
        }
      }
      ]
    }
  ]

  private subscription: Subscription = new Subscription()

  constructor (private message: EoNgFeedbackMessageService,
    private modalService:EoNgFeedbackModalService,
    private api:ApiService, private router:Router,
    private baseInfo:BaseInfoService,
    private appConfigService:AppConfigService) {
    this.appConfigService.reqFlashBreadcrumb([{ title: 'API管理', routerLink: 'router/group/list' }])
  }

  ngOnInit (): void {
    this.apisForm.group_uuid = this.baseInfo.allParamsInfo.apiGroupId
    this.getApisData()
    // 当左侧分组中目录被选中时，group_uuid参数会变化，随之获取新的列表
    this.subscription = this.router.events.subscribe((event) => {
      if (event instanceof NavigationEnd) {
        this.apisForm.group_uuid = this.baseInfo.allParamsInfo.apiGroupId
        this.getApisData()
        this.apisForm.page_num = 1
        this.apisForm.page_size = 20
        this.apiNameForSear = ''
      }
    })

    this.getSourcesList()
  }

  ngAfterViewInit () {
    this.apisTableBody[2].title = this.methodTpl
  }

  ngOnDestroy () {
    this.subscription.unsubscribe()
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  apisTableClick= (item:any) => {
    this.router.navigate(['/', 'router', 'content', item.data.uuid])
  }

  // 根据group_uuid获取新的apis列表
  getApisData () {
    this.api.get('routers', { group_uuid: (this.apisForm.group_uuid || this.groupUuid), search_name: this.apiNameForSear, source_ids: this.apisForm.source_ids.join(','), page_num: this.apisForm.page_num, page_size: this.apisForm.page_size }).subscribe((resp:any) => {
      if (resp.code === 0) {
        this.apisForm.apis = this.apisSet.size > 0
          ? resp.data.apis.map((item:any) => {
            item.checked = this.apisSet.has(item.uuid)
            return item
          })
          : resp.data.apis
        this.apisForm.group_uuid = this.apisForm.group_uuid || this.groupUuid
        this.apisForm.total = resp.data.total || this.apisForm.total
        this.apisForm.page_num = resp.data.page_num || this.apisForm.page_num
        this.apisForm.page_size = resp.data.page_size || this.apisForm.page_size
      } else {
        this.message.error(resp.msg || '获取API列表数据失败！')
      }
    })
  }

  // 获取来源可选列表，供列表筛选用
  getSourcesList () {
    this.api.get('router/source').subscribe((resp:{code:number, data:{list:Array<{id:string, title:string}>}, msg:string}) => {
      if (resp.code === 0) {
        for (const index in resp.data.list) {
          this.sourcesList.push({ text: resp.data.list[index].title, value: resp.data.list[index].id })
          this.apisTableHeadName[5].filterOpts = this.sourcesList
        }
      } else {
        this.message.error(resp.msg || '获取来源列表数据失败！')
      }
    })
  }

  addApi () {
    const url:Array<string> = ['/', 'router', 'create']
    if (this.apisForm.group_uuid) {
      url.push(this.apisForm.group_uuid)
    }
    this.router.navigate(url)
  }

  changeApisSet (item: any, type?:string) {
    if (type === 'all') {
      if (item) {
        for (const index in this.apisForm.apis) {
          this.apisSet.add(this.apisForm.apis[index].uuid)
        }
      } else {
        this.apisSet = new Set()
      }
    } else {
    // 被取消勾选
      if (item?.checked) {
        this.apisSet.delete(item.uuid)
      } else {
      // 被选中
        this.apisSet.add(item.uuid)
      }
    }
  }

  // 删除api弹窗
  deleteApiModal (items:any) {
    this.modalService.create({
      nzTitle: '删除',
      nzContent: '该数据删除后将无法找回，请确认是否删除？',
      nzClosable: true,
      nzClassName: 'delete-modal',
      nzWidth: MODAL_SMALL_SIZE,
      nzOkDanger: true,
      nzOnOk: () => {
        this.deleteApi(items)
      }
    })
  }

  // 删除单个api
  deleteApi = (items:any) => {
    this.api.delete('router', { uuid: items.uuid }).subscribe((resp:any) => {
      if (resp.code === 0) {
        this.message.success(resp.msg || '删除成功!', { nzDuration: 1000 })
        this.apisSet.delete(items.uuid)
        this.getApisData()
      } else {
        this.message.error(resp.msg || '删除失败!')
      }
    })
  }

  // 过滤器内选择的value变化时的回调
  apisFilterChange (value:any) {
    this.apisForm.source_ids = value.keys
    this.getApisData()
  }
}

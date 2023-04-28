/* eslint-disable dot-notation */
import { Component, Input, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { Router, NavigationEnd } from '@angular/router'
import { NzTreeNodeOptions } from 'ng-zorro-antd/tree'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { EoNgNavigationService } from 'projects/core/src/app/service/eo-ng-navigation.service'
import { Subscription } from 'rxjs'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { MODAL_SMALL_SIZE } from 'projects/core/src/app/constant/app.config'
import { apisTableHeadName, apisTableBody } from '../../types/conf'
import { APIList } from '../../types/types'

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

  apisForm: {apis:APIList[], total:number, pageNum:number, pageSize:number, groupUuid:string, sourceIds:Array<string>} = {
    apis: [],
    total: 0,
    pageNum: 1,
    pageSize: 20,
    groupUuid: '',
    sourceIds: []
  }

  apisTableHeadName:THEAD_TYPE[] = [...apisTableHeadName]
  apisTableBody:TBODY_TYPE[] = [...apisTableBody]

  private subscription: Subscription = new Subscription()

  constructor (private message: EoNgFeedbackMessageService,
    private modalService:EoNgFeedbackModalService,
    private api:ApiService,
    private router:Router,
    private baseInfo:BaseInfoService,
    private navigationService:EoNgNavigationService) {
    this.navigationService.reqFlashBreadcrumb([{ title: 'API管理', routerLink: 'router/api/group/list' }])
  }

  ngOnInit (): void {
    this.apisForm.groupUuid = this.baseInfo.allParamsInfo.apiGroupId
    this.initTable()
    this.getApisData()
    // 当左侧分组中目录被选中时，groupUuid参数会变化，随之获取新的列表
    this.subscription = this.router.events.subscribe((event) => {
      if (event instanceof NavigationEnd) {
        this.apisForm.groupUuid = this.baseInfo.allParamsInfo.apiGroupId
        this.getApisData()
        this.apisForm.pageNum = 1
        this.apisForm.pageSize = 20
        this.apiNameForSear = ''
      }
    })

    this.getSourcesList()
  }

  initTable () {
    this.apisTableHeadName[0].click = (item:any) => {
      this.changeApisSet(item, 'all')
    }
    this.apisTableHeadName[0].showFn = () => {
      return !this.nzDisabled
    }
    this.apisTableBody[0].click = (item:any) => {
      this.changeApisSet(item)
    }
    this.apisTableBody[0].showFn = () => {
      return !this.nzDisabled
    }
    this.apisTableBody[7].btns[0].click = (item:any) => {
      this.router.navigate(['/', 'router', 'api', 'content', item.data.uuid, 'publish'])
    }
    this.apisTableBody[7].btns[1].click = (item:any) => {
      this.router.navigate(['/', 'router', 'api', 'content', item.data.uuid])
    }
    this.apisTableBody[7].btns[2].click = (item:any) => {
      this.deleteApiModal(item.data)
    }
    this.apisTableBody[7].btns[2].disabledFn = (data:any, item:any) => { return !item.data.isDelete || this.nzDisabled }
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

  apisTableClick= (item:{data:{uuid:string, [k:string]:any}}) => {
    this.router.navigate(['/', 'router', 'api', 'content', item.data.uuid])
  }

  // 根据groupUuid获取新的apis列表
  getApisData () {
    this.api.get('routers', { groupUuid: (this.apisForm.groupUuid || this.groupUuid), searchName: this.apiNameForSear, sourceIds: this.apisForm.sourceIds.join(','), pageNum: this.apisForm.pageNum, pageSize: this.apisForm.pageSize }).subscribe((resp:any) => {
      if (resp.code === 0) {
        this.apisForm.apis = this.apisSet.size > 0
          ? resp.data.apis.map((item:any) => {
            item.checked = this.apisSet.has(item.uuid)
            return item
          })
          : resp.data.apis

        this.apisForm.groupUuid = this.apisForm.groupUuid || this.groupUuid
        this.apisForm.total = resp.data.total || this.apisForm.total
        this.apisForm.pageNum = resp.data.pageNum || this.apisForm.pageNum
        this.apisForm.pageSize = resp.data.pageSize || this.apisForm.pageSize
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
      }
    })
  }

  addApi () {
    const url:Array<string> = ['/', 'router', 'api', 'create']
    if (this.apisForm.groupUuid) {
      url.push(this.apisForm.groupUuid)
    }
    this.router.navigate(url)
  }

  changeApisSet (item: {uuid:string, [k:string]:any}, type?:string) {
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
      if (item?.['checked']) {
        this.apisSet.delete(item.uuid)
      } else {
      // 被选中
        this.apisSet.add(item.uuid)
      }
    }
  }

  // 删除api弹窗
  deleteApiModal (items:{uuid:string, [k:string]:any}) {
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
  deleteApi = (items:{uuid:string, [k:string]:any}) => {
    this.api.delete('router', { uuid: items.uuid }).subscribe((resp:any) => {
      if (resp.code === 0) {
        this.message.success(resp.msg || '删除成功!', { nzDuration: 1000 })
        this.apisSet.delete(items.uuid)
        this.getApisData()
      }
    })
  }

  // 过滤器内选择的value变化时的回调
  apisFilterChange (value:any) {
    this.apisForm.sourceIds = value.keys
    this.getApisData()
  }
}

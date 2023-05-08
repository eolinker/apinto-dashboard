import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { SelectOption } from 'eo-ng-select'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { IntelligentPluginDefaultThead } from '../types/conf'
import { IntelligentPluginService } from '../intelligent-plugin.service'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { MODAL_NORMAL_SIZE, MODAL_SMALL_SIZE } from '../../../constant/app.config'
import { ApiService } from '../../../service/api.service'
import { IntelligentPluginPublishComponent } from '../publish/publish.component'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { IntelligentPluginCreateComponent } from '../create/create.component'
import { DynamicConfig, DynamicDriverData, DynamicField, DynamicListStatus, DynamicRender } from '../types/types'
import { ClusterSimpleOption, EmptyHttpResponse } from '../../../constant/type'
import { Subscription, forkJoin, map } from 'rxjs'
import { v4 as uuidv4 } from 'uuid'
import { BaseInfoService } from '../../../service/base-info.service'
import { NavigationEnd, Router } from '@angular/router'
import { EoNgNavigationService } from '../../../service/eo-ng-navigation.service'

@Component({
  selector: 'eo-ng-intelligent-plugin-list',
  templateUrl: './list.component.html',
  styles: [
  ]
})
export class IntelligentPluginListComponent implements OnInit {
  @ViewChild('clusterStatusTpl', { read: TemplateRef, static: true }) clusterStatusTpl: TemplateRef<any> | undefined
  @ViewChild('loadingTpl', { read: TemplateRef, static: true }) loadingTpl: TemplateRef<any> | undefined
  moduleName:string = ''
  pluginName:string = ''
  keyword:string = ''
  nzDisabled:boolean = false
  cluster:any = []
  clusterOptions:SelectOption[] = []
  tableBody:TBODY_TYPE[] = [...this.service.createTbody(this)]
  tableHeadName:THEAD_TYPE[] = [...IntelligentPluginDefaultThead]
  tableData:{data:any[], pagination:boolean, total:number, pageNum:number, pageSize:number}
  = { data: [], pagination: true, total: 1, pageSize: 20, pageNum: 1 }

  driverOptions:SelectOption[] = []
  renderSchema:any = {} // 动态渲染数据，是json schema
  modalRef:NzModalRef|undefined
  statusMap:{[k:string]:any} = {}
  tableLoading:boolean = true
  private subscription: Subscription = new Subscription()

  constructor (
    private message: EoNgFeedbackMessageService,
    private service:IntelligentPluginService,
    private modalService:EoNgFeedbackModalService,
    private api:ApiService,
    private router:Router,
    private baseInfo:BaseInfoService,
    private navigationService: EoNgNavigationService) {

  }

  ngOnInit (): void {
    this.moduleName = this.baseInfo.allParamsInfo.moduleName
    this.subscription = this.router.events.subscribe((event) => {
      if (event instanceof NavigationEnd) {
        this.moduleName = this.baseInfo.allParamsInfo.moduleName
        this.pluginName = ''
        this.keyword = ''
        this.cluster = []
        this.clusterOptions = []
        this.tableBody = [...this.service.createTbody(this)]
        this.tableHeadName = [...IntelligentPluginDefaultThead]
        this.tableData = { data: [], pagination: true, total: 1, pageSize: 20, pageNum: 1 }

        this.driverOptions = []
        this.renderSchema = {} // 动态渲染数据，是json schema
        this.modalRef = undefined
        this.statusMap = {}
        this.tableLoading = true
        this.getClusters()
        this.getRender()
        this.getTableData()
      }
    })
    this.getClusters()
    this.getRender()
    this.getTableData()
  }

  ngOnDestroy () {
    this.subscription.unsubscribe()
  }

  getTableData () {
    this.tableLoading = true
    // 表格内的其他数据与状态数据是分别获取的，如果list先返回，需要先展示除了状态数据以外的其他数据
    forkJoin([this.api.get(`dynamic/${this.moduleName}/list`, {
      page: this.tableData.pageNum,
      pageSize: this.tableData.pageSize,
      keyword: this.keyword,
      cluster: JSON.stringify(this.cluster)
    }).pipe(
      map(res => {
        if (res.code === 0) {
          this.getConfig(res.data)
        }
        return res
      })),
    this.api.get(`dynamic/${this.moduleName}/status`, {
      page: this.tableData.pageNum,
      pageSize: this.tableData.pageSize,
      keyword: this.keyword,
      cluster: JSON.stringify(this.cluster)
    })]).subscribe((val:Array<any>) => {
      this.refreshTableData(this.tableData.data, val[1].data)
    })
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  tableClick = () => { }

  getRender () {
    this.api.get(`dynamic/${this.moduleName}/render`).subscribe((resp:{code:number, msg:string, data:DynamicRender}) => {
      if (resp.code === 0) {
        this.renderSchema = resp.data.render
      }
    })
  }

  getClusters () {
    this.api.get('clusters/simple').subscribe((resp:{code:number, msg:string, data:{clusters:ClusterSimpleOption[]}}) => {
      if (resp.code === 0) {
        this.clusterOptions = resp.data.clusters.map((cluster:ClusterSimpleOption) => {
          return { label: cluster.title, value: cluster.name }
        })
        this.cluster = this.clusterOptions.map((cluster:SelectOption) => {
          return cluster.value
        })
      }
    })
  }

  // 获取列表渲染配置、表单渲染配置
  getConfig (data:DynamicConfig) {
    this.navigationService.reqFlashBreadcrumb([
      { title: data.title }
    ])
    this.pluginName = data.title
    this.getTableConfig(data.fields) // 获取列表配置
    this.tableData.data = data.list // 获取列表数据
    this.driverOptions = data.drivers?.map((driver:DynamicDriverData) => {
      return { label: driver.title, value: driver.name }
    }) || []
  }

  refreshTableData (tableData:Array<{[k:string]:any}>, statusData:DynamicListStatus) {
    if (tableData.length && statusData && Object.keys(statusData).length) {
      this.tableData.data = tableData.map((item:any) => {
        return { ...item, ...statusData[item.id] }
      })
      // 将table的loding取消
      this.tableLoading = false
    }
  }

  // table需要设置为loading状态
  private getTableConfig (fields:DynamicField[]) {
    const newTableHeadConfig:THEAD_TYPE[] = []
    const newTableBodyConfig:TBODY_TYPE[] = []
    let statusColFlag:boolean = true
    for (const field of fields) {
      if (field.attr === 'status' && statusColFlag) {
        newTableHeadConfig.push(
          {
            title: '状态',
            showFn: () => {
              return this.tableLoading
            }
          }
        )

        newTableBodyConfig.push(
          {
            title: this.loadingTpl,
            showFn: (item:any) => {
              return item.id === this.tableData.data[0].id && this.tableLoading
            },
            seRowspan: () => {
              return this.tableData.data.length
            }
          }
        )
        statusColFlag = false
      }
      newTableHeadConfig.push(
        {
          title: field.title,
          ...(newTableHeadConfig.length === 0 ? { left: true } : {}),
          ...(field.enum?.length > 0
            ? {
                filterMultiple: true,
                filterOpts: field.enum.map((item:string) => {
                  return { text: item, value: item }
                }),
                filterFn: (list: string[], item: any) => {
                  return list.some((name) => item.data[field.name].indexOf(name) !== -1)
                }
              }
            : {}),
          ...(field.attr === 'status'
            ? {
                showFn: () => {
                  return !this.tableLoading
                }
              }
            : {})

        }
      )

      newTableBodyConfig.push(
        {
          key: field.name,
          ...(field.attr ? { title: this.getTdTpl(field.attr) } : {}),
          ...(newTableHeadConfig.length === 0 ? { left: true } : {}),
          ...(field.attr === 'status'
            ? {
                showFn: () => {
                  return !this.tableLoading
                }
              }
            : {})
        }
      )
    }
    this.tableBody = [...newTableBodyConfig, ...this.service.createTbody(this, 'btn')]
    this.tableHeadName = [...newTableHeadConfig, { title: '操作', right: true }]
  }

  getTdTpl (attr:string) {
    if (attr === 'status') {
      return this.clusterStatusTpl
    }
    return this.clusterStatusTpl
  }

  publish (value:any) {
    this.modalRef = this.modalService.create({
      nzTitle: `${value.data.title}上线管理`,
      nzWidth: MODAL_NORMAL_SIZE,
      nzContent: IntelligentPluginPublishComponent,
      nzComponentParams: {
        name: value.data.title,
        id: value.data.id,
        desc: value.data.description,
        moduleName: this.moduleName,
        closeModal: this.closeModal,
        nzDisabled: this.nzDisabled
      },
      nzFooter: [{
        label: '取消',
        type: 'default',
        onClick: () => {
          this.modalRef?.close()
        }
      },
      {
        label: '下线',
        danger: true,
        onClick: (context:IntelligentPluginPublishComponent) => {
          context.offline()
        },
        disabled: () => {
          return this.nzDisabled
        }
      },
      {
        label: '上线',
        type: 'primary',
        onClick: (context:IntelligentPluginPublishComponent) => {
          context.online()
        },
        disabled: () => {
          return this.nzDisabled
        }
      }]
    })
  }

  addData () {
    this.modalRef = this.modalService.create({
      nzTitle: `新建${this.pluginName}`,
      nzWidth: MODAL_NORMAL_SIZE,
      nzContent: IntelligentPluginCreateComponent,
      nzComponentParams: {
        renderSchema: this.renderSchema,
        editPage: false,
        moduleName: this.moduleName,
        driverSelectOptions: this.driverOptions,
        initFormValue: { driver: this.driverOptions[0].value || '' }
      },
      nzOnOk: (component:IntelligentPluginCreateComponent) => {
        component.form.validate().then((res:any) => {
          if (!res) {
            // eslint-disable-next-line dot-notation
            this.saveData(JSON.parse(JSON.stringify(component.form['values'])))
          }
        })
        return false
      }
    })
  }

  editData (value:any) {
    this.modalRef = this.modalService.create({
      nzTitle: `编辑${this.pluginName}`,
      nzWidth: MODAL_NORMAL_SIZE,
      nzContent: IntelligentPluginCreateComponent,
      nzComponentParams: {
        renderSchema: this.renderSchema,
        editPage: true,
        moduleName: this.moduleName,
        uuid: value.data.id,
        driverSelectOptions: this.driverOptions
      },
      nzOnOk: (component:IntelligentPluginCreateComponent) => {
        component.form.validate().then((res:any) => {
          if (!res) {
            // eslint-disable-next-line dot-notation
            this.saveData(JSON.parse(JSON.stringify(component.form.values)), component.uuid, true)
          }
        })
        return false
      }
    })
  }

  saveData (form:{[k:string]:any}, id:string = uuidv4(), editPage?:boolean) {
    if (editPage) {
      this.api.put(`dynamic/${this.moduleName}/config/${id}`, { ...form }).subscribe((resp:EmptyHttpResponse) => {
        if (resp.code === 0) {
          this.message.success(resp.msg || '操作成功')
          this.getTableData()
          this.modalRef?.close()
        }
      })
    } else {
      this.api.post(`dynamic/${this.moduleName}`, { ...form }).subscribe((resp:EmptyHttpResponse) => {
        if (resp.code === 0) {
          this.message.success(resp.msg || '操作成功')
          this.getTableData()
          this.modalRef?.close()
        }
      })
    }
  }

  closeModal = () => {
    this.modalRef?.close()
    this.getTableData()
  }

  deleteDataModal (items:{id:string, [k:string]:any}) {
    this.modalService.create({
      nzTitle: '删除',
      nzContent: '该数据删除后将无法找回，请确认是否删除？',
      nzClosable: true,
      nzClassName: 'delete-modal',
      nzWidth: MODAL_SMALL_SIZE,
      nzOkDanger: true,
      nzOnOk: () => {
        this.deleteData(items)
      }
    })
  }

  // 删除单条数据
  deleteData = (items:{id:string, [k:string]:any}) => {
    this.api.delete(`dynamic/${this.moduleName}/batch`, { uuids: JSON.stringify([items.id]) }).subscribe((resp:any) => {
      if (resp.code === 0) {
        this.message.success(resp.msg || '删除成功!')
        this.getTableData()
      }
    })
  }
}

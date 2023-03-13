/* eslint-disable dot-notation */
/* eslint-disable no-useless-constructor */
/* eslint-disable camelcase */
import {
  Component,
  EventEmitter,
  Input,
  OnInit,
  Output,
  TemplateRef,
  ViewChild
} from '@angular/core'
import { ActivatedRoute, Router } from '@angular/router'
import {
  EoNgFeedbackMessageService,
  EoNgFeedbackModalService
} from 'eo-ng-feedback'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'
import { UntypedFormBuilder, FormGroup, Validators } from '@angular/forms'
import { Subscription } from 'rxjs'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { MODAL_SMALL_SIZE } from '../../../constant/app.config'
import { defaultAutoTips } from '../../../constant/conf'
import { ApiManagementProxyComponent } from '../proxy/proxy.component'
import { BaseInfoService } from '../../../service/base-info.service'

@Component({
  selector: 'eo-ng-api-create',
  templateUrl: './create.component.html',
  styles: [
    `
      eo-ng-table.ant-table {
        min-width: 508px !important;
      }

      .ant-form-item-control:first-child:not([class^='ant-col-']):not(
          [class*=' ant-col-']
        ) {
        width: auto !important;
      }

      nz-form-item.ant-row.checkbox-group-api.ant-form-item.ant-form-item-has-error {
        margin-bottom: 0 !important;
      }
    `
  ]
})
export class ApiCreateComponent implements OnInit {
  @ViewChild('optTypeTranslateTpl', { read: TemplateRef, static: true })
  optTypeTranslateTpl: TemplateRef<any> | undefined

  @Input() apiUuid: string = ''
  @Input() editPage: boolean = false
  @Input() groupUuid: string = ''
  @Input() fromList: boolean = false
  @Output() changeToList: EventEmitter<any> = new EventEmitter()
  validateForm: FormGroup = new FormGroup({})
  nzDisabled: boolean = false
  private subscription: Subscription = new Subscription()

  createApiForm: {
    uuid: string
    method: Array<string>
    match: Array<{
      position: string
      match_type: string
      key: string
      pattern: string
    }>
    proxy_header: Array<{ opt_type: string; key: string; value: string }>
  } = {
    uuid: '',
    method: [],
    match: [],
    proxy_header: []
  }

  headerList: Array<any> = []
  firstLevelList: Array<string> = []
  serviceList: Array<{ label: string; value: string }> = []
  methodList: Array<{ label: string; value: string; checked: boolean }> = [
    { label: 'GET', value: 'GET', checked: false },
    { label: 'POST', value: 'POST', checked: false },
    { label: 'PUT', value: 'PUT', checked: false },
    { label: 'DELETE', value: 'DELETE', checked: false },
    { label: 'PATCH', value: 'PATCH', checked: false },
    { label: 'HEAD', value: 'HEAD', checked: false }
  ]

  allChecked: boolean = false

  matchHeaderSet: Set<string> = new Set()
  proxyRequestSet: Set<string> = new Set()

  autoTips: Record<string, Record<string, string>> = defaultAutoTips

  proxyHeaderTableHeadName: Array<object> = [
    {
      title: '操作类型'
    },
    { title: '参数名' },
    { title: '参数值' },
    {
      title: '操作',
      right: true
    }
  ]

  proxyHeaderTableBody: Array<any> = [
    { key: 'opt_type' },
    { key: 'key' },
    { key: 'value' },
    {
      type: 'btn',
      right: true,
      btns: [
        {
          title: '配置',
          disabledFn: () => {
            return this.nzDisabled
          },
          click: (item: any) => {
            this.openDrawer('proxyHeader', item.data)
          }
        },
        {
          title: '删除',
          action: 'delete',
          disabledFn: () => {
            return this.nzDisabled
          }
        }
      ]
    }
  ]

  modalRef: NzModalRef | undefined

  proxyEdit: boolean = false
  editData: any = null

  constructor(
    private message: EoNgFeedbackMessageService,
    private baseInfo: BaseInfoService,
    private api: ApiService,
    private activateInfo: ActivatedRoute,
    private appConfigService: AppConfigService,
    private fb: UntypedFormBuilder,
    private router: Router,
    private modalService: EoNgFeedbackModalService
  ) {
    this.appConfigService.reqFlashBreadcrumb([
      { title: 'API管理', routerLink: 'router/group/list' },
      { title: '新建API' }
    ])

    this.validateForm = this.fb.group({
      group_uuid: ['', [Validators.required]],
      name: ['', [Validators.required]],
      desc: [''],
      request_path: ['', [Validators.required, Validators.pattern('^[^?]*')]],
      service: ['', [Validators.required]],
      proxy_path: [''],
      timeout: [10000, [Validators.required]],
      retry: [0, [Validators.required]],
      enable_websocket: [false]
    })
  }

  ngOnInit(): void {
    this.getServiceList()

    if (this.baseInfo.allParamsInfo.apiId) {
      this.appConfigService.reqFlashBreadcrumb([
        { title: 'API管理', routerLink: 'router/group/list' },
        { title: 'API信息' }
      ])
    }
    if (this.baseInfo.allParamsInfo.apiGroupId) {
      this.groupUuid = this.baseInfo.allParamsInfo.apiGroupId
    }
    if (this.editPage) {
      this.getApiMessage()
    } else {
      this.getHeaderList()
    }
  }

  ngAfterViewInit() {
    this.proxyHeaderTableBody[0].title = this.optTypeTranslateTpl
  }

  ngOnDestroy() {
    this.subscription.unsubscribe()
  }

  disabledEdit(value: any) {
    this.nzDisabled = value
  }

  // 当编辑api时，需要获取api信息
  getApiMessage() {
    this.api.get('router', { uuid: this.apiUuid }).subscribe((resp) => {
      if (resp.code === 0) {
        this.validateForm.controls['group_uuid'].setValue(
          resp.data.api.group_uuid
        )
        this.validateForm.controls['name'].setValue(resp.data.api.name)
        this.validateForm.controls['desc'].setValue(resp.data.api.desc)
        this.validateForm.controls['request_path'].setValue(
          resp.data.api.request_path.slice(1)
        )
        this.validateForm.controls['service'].setValue(resp.data.api.service)
        this.validateForm.controls['proxy_path'].setValue(
          resp.data.api.proxy_path
        )
        this.validateForm.controls['timeout'].setValue(resp.data.api.timeout)
        this.validateForm.controls['retry'].setValue(resp.data.api.retry)
        this.validateForm.controls['enable_websocket'].setValue(
          resp.data.api.enable_websocket
        )
        this.createApiForm = resp.data.api
        if (
          !this.createApiForm.method ||
          this.createApiForm.method.length === 0
        ) {
          this.createApiForm.method = [
            'POST',
            'PUT',
            'GET',
            'DELETE',
            'PATCH',
            'HEAD',
            'OPTIONS'
          ]
          this.allChecked = true
          this.updateAllChecked()
        } else {
          this.initCheckbox()
        }
        this.getHeaderList()
      } else {
        this.message.error(resp.msg || '获取数据失败!')
      }
    })
  }

  // 获取请求头部列表参数
  getHeaderList() {
    this.api.get('router/groups').subscribe((resp: any) => {
      if (resp.code === 0) {
        this.headerList = []
        for (const index in resp.data.root.groups) {
          this.firstLevelList.push(resp.data.root.groups[index].uuid)
          if (
            !resp.data.root.groups[index].children ||
            resp.data.root.groups[index].children.length === 0
          ) {
            resp.data.root.groups[index].disabled = true
          } else {
            resp.data.root.groups[index].selectable = false
            this.headerList.push(resp.data.root.groups[index])
          }
        }
        this.headerList = this.transferHeader(this.headerList)
        if (
          this.groupUuid &&
          this.firstLevelList.indexOf(this.groupUuid) === -1
        ) {
          this.validateForm.controls['group_uuid'].setValue(
            this.baseInfo.allParamsInfo.apiGroupId
          )
        }
      } else {
        this.message.error(resp.msg || '获取列表数据失败!')
      }
    })
  }

  transferHeader(header: any) {
    for (const index in header) {
      header[index].title = header[index].name
      header[index].key = header[index].uuid
      if (!header[index].children || header[index].children.length === 0) {
        header[index].isLeaf = true
      } else {
        header[index].children = this.transferHeader(header[index].children)
      }
    }
    return header
  }

  nzTreeClick(value: any) {
    if (value.node.origin.selectable === false) {
      value.node.origin.expanded = !value.node.origin.expanded
    }
    this.headerList = [...this.headerList]
  }

  // 获取上游服务列表
  getServiceList() {
    this.api.get('service/enum').subscribe((resp: any) => {
      if (resp.code === 0) {
        this.serviceList = []
        // this.validateForm.controls['service'].setValue(resp.data.list[0])
        for (const item of resp.data.list) {
          this.serviceList = [...this.serviceList, { label: item, value: item }]
        }
      } else {
        this.message.error(resp.msg || '获取列表数据失败!')
      }
    })
  }

  updateAllChecked(): void {
    if (this.allChecked) {
      this.methodList = this.methodList.map((item: any) => ({
        ...item,
        checked: true
      }))
      this.createApiForm.method = []
      for (const index in this.methodList) {
        if (this.methodList[index].checked) {
          this.createApiForm.method.push(this.methodList[index].value)
        }
      }
      this.showCheckboxGroupValid = false
    } else {
      this.methodList = this.methodList.map((item: any) => ({
        ...item,
        checked: false
      }))
      this.createApiForm.method = []
      this.showCheckboxGroupValid = false
    }
  }

  initCheckbox(): void {
    for (const index in this.methodList) {
      if (
        this.createApiForm.method.indexOf(this.methodList[index].label) !== -1
      ) {
        this.methodList[index].checked = true
      }
    }
  }

  updateSingleChecked(): void {
    if (this.methodList.every((item: any) => !item.checked)) {
      this.allChecked = false
    } else if (this.methodList.every((item: any) => item.checked)) {
      this.allChecked = true
    } else {
      this.allChecked = false
    }
    this.createApiForm.method = []
    for (const index in this.methodList) {
      if (this.methodList[index].checked) {
        this.createApiForm.method.push(this.methodList[index].value)
      }
    }
    if (this.methodList.length > 0) {
      this.showCheckboxGroupValid = false
    }
  }

  proxyTableClick = (item: any) => {
    this.openDrawer('proxyHeader', item.data)
  }

  openDrawer(type: string, data?: any) {
    switch (type) {
      case 'proxyHeader': {
        if (data) {
          this.editData = data
          this.proxyEdit = true
        }

        this.modalRef = this.modalService.create({
          nzTitle: '配置转发上游请求头',
          nzContent: ApiManagementProxyComponent,
          nzComponentParams: {
            data: data || {},
            editPage: !!data
          },
          nzClosable: true,
          nzWidth: MODAL_SMALL_SIZE,
          nzCancelText: '取消',
          nzOkText: this.proxyEdit ? '提交' : '保存',
          nzOnOk: (proxyRef: ApiManagementProxyComponent) => {
            this.saveProxyHeader(proxyRef)
            return false
          }
        })
        this.modalRef.afterClose.subscribe(() => {
          this.proxyEdit = false
        })
        break
      }
    }
  }

  // 返回列表页，当fromList为true时，该页面左侧有分组
  backToList() {
    this.router.navigate(['/', 'router', 'group', 'list'])
  }

  // 提交api数据
  saveApi() {
    if (this.createApiForm.method.length === 0 && !this.allChecked) {
      this.showCheckboxGroupValid = true
    } else {
      this.showCheckboxGroupValid = false
    }
    if (this.validateForm.valid && !this.showCheckboxGroupValid) {
      if (this.allChecked) {
        this.createApiForm.method = []
      }
      if (this.editPage) {
        this.api
          .put(
            'router',
            {
              name: this.validateForm.controls['name'].value,
              uuid: this.createApiForm.uuid,
              group_uuid: this.validateForm.controls['group_uuid'].value,
              desc: this.validateForm.controls['desc'].value,
              request_path:
                '/' + this.validateForm.controls['request_path'].value,
              service: this.validateForm.controls['service'].value,
              method: this.createApiForm.method,
              proxy_path: this.validateForm.controls['proxy_path'].value,
              timeout: Number(this.validateForm.controls['timeout'].value),
              retry: Number(this.validateForm.controls['retry'].value),
              enable_websocket:
                this.validateForm.controls['enable_websocket'].value || false,
              proxy_header: this.createApiForm.proxy_header,
              match: this.createApiForm.match
            },
            { uuid: this.apiUuid }
          )
          .subscribe((resp) => {
            if (resp.code === 0) {
              this.backToList()
              this.message.success(resp.msg || '修改成功！', {
                nzDuration: 1000
              })
            } else {
              this.message.error(resp.msg || '修改失败!')
            }
          })
      } else {
        this.api
          .post('router', {
            ...this.validateForm.value,
            uuid: this.createApiForm.uuid,
            group_uuid: this.validateForm.controls['group_uuid'].value,
            method: this.createApiForm.method,
            proxy_header: this.createApiForm.proxy_header,
            match: this.createApiForm.match,
            request_path:
              '/' + this.validateForm.controls['request_path'].value,
            proxy_path: this.validateForm.controls['proxy_path'].value
          })
          .subscribe((resp) => {
            if (resp.code === 0) {
              this.message.success(resp.msg || '添加成功！', {
                nzDuration: 1000
              })
              this.backToList()
            } else {
              this.message.error(resp.msg || '添加失败!')
            }
          })
      }
    } else {
      Object.values(this.validateForm.controls).forEach((control) => {
        if (control.invalid) {
          control.markAsDirty()
          control.updateValueAndValidity({ onlySelf: true })
        }
      })
    }
  }

  validate = (option: any): boolean => {
    const uuid = option.uuid as string
    return this.firstLevelList.indexOf(uuid) === -1
  }

  showCheckboxGroupValid: boolean = false

  requestPathChange() {
    if (
      !this.validateForm.controls['proxy_path'].value &&
      this.validateForm.controls['request_path'].value
    ) {
      this.validateForm.controls['proxy_path'].setValue(
        '/' + this.validateForm.controls['request_path'].value
      )
    }
  }

  checkTimeout() {
    if (
      this.validateForm.controls['timeout'].value !== null &&
      this.validateForm.controls['timeout'].value < 1
    ) {
      this.validateForm.controls['timeout'].setValue(1)
    }
  }

  // 保存转发上游请求头数据时，如果是新建数据，直接加入tableList，如果是编辑数据，需要删除原先同key的数据再保存
  saveProxyHeader(proxyRef: ApiManagementProxyComponent): void {
    if (proxyRef.validateProxyHeaderForm.valid) {
      if (this.proxyEdit) {
        for (const index in this.createApiForm.proxy_header) {
          if (
            this.createApiForm.proxy_header[index].key === this.editData.key &&
            this.createApiForm.proxy_header[index].opt_type ===
              this.editData.opt_type &&
            this.createApiForm.proxy_header[index].value === this.editData.value
          ) {
            this.createApiForm.proxy_header.splice(Number(index), 1)
            break
          }
        }
      }
      this.createApiForm.proxy_header = [
        {
          opt_type: proxyRef.validateProxyHeaderForm.controls['opt_type'].value,
          key: proxyRef.validateProxyHeaderForm.controls['key'].value,
          value: proxyRef.validateProxyHeaderForm.controls['value'].value
        },
        ...this.createApiForm.proxy_header
      ]
      this.modalRef?.close()
    } else {
      Object.values(proxyRef.validateProxyHeaderForm.controls).forEach(
        (control) => {
          if (control.invalid) {
            control.markAsDirty()
            control.updateValueAndValidity({ onlySelf: true })
          }
        }
      )
    }
  }
}

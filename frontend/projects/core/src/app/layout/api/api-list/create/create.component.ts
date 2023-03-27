/* eslint-disable dot-notation */
/* eslint-disable no-useless-constructor */
import {
  Component,

  Input,
  OnInit,

  TemplateRef,
  ViewChild
} from '@angular/core'
import { Router } from '@angular/router'
import {
  EoNgFeedbackMessageService,
  EoNgFeedbackModalService
} from 'eo-ng-feedback'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'
import { UntypedFormBuilder, FormGroup, Validators } from '@angular/forms'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { ApiManagementProxyComponent } from '../proxy/proxy.component'
import { NzTreeNodeOptions } from 'ng-zorro-antd/tree'
import { SelectOption } from 'eo-ng-select'
import { CheckBoxOptionInterface } from 'eo-ng-checkbox'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { setFormValue } from 'projects/core/src/app/constant/form'
import { ApiGroup, ApiGroupsData } from 'projects/core/src/app/constant/type'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { methodList, proxyHeaderTableHeadName, proxyHeaderTableBody } from '../../types/conf'
import { APINotFormGroupData } from '../../types/types'
import { MODAL_SMALL_SIZE } from 'projects/core/src/app/constant/app.config'
import { cloneDeep } from 'lodash'

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
  @ViewChild('optTypeTranslateTpl', { read: TemplateRef, static: true }) optTypeTranslateTpl: TemplateRef<any> | undefined
  @Input() apiUuid:string = ''
  @Input() editPage:boolean = false
  @Input() groupUuid:string = ''
  nzDisabled:boolean = false
  headerList:NzTreeNodeOptions[]= []
  firstLevelList:Array<string> = []
  serviceList:SelectOption[]= []
  methodList:CheckBoxOptionInterface[]= [...cloneDeep(methodList)]
  allChecked:boolean = false
  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  proxyHeaderTableHeadName:Array<object> = [...proxyHeaderTableHeadName]
  proxyHeaderTableBody:Array<any> = [...proxyHeaderTableBody]
  modalRef:NzModalRef | undefined
  proxyEdit:boolean = false
  editData:any = null
  validateForm:FormGroup = new FormGroup({})
  createApiForm:APINotFormGroupData={
    uuid: '',
    method: [],
    match: [],
    proxyHeader: []
  }

  pluginTemplateList:SelectOption[] = []

  constructor (private message: EoNgFeedbackMessageService,
    private baseInfo:BaseInfoService,
    private api:ApiService,
    private appConfigService:AppConfigService,
    private fb: UntypedFormBuilder,
    private router: Router,
    private modalService: EoNgFeedbackModalService
  ) {
    this.appConfigService.reqFlashBreadcrumb([
      { title: 'API管理', routerLink: 'router/api/group/list' },
      { title: '新建API' }
    ])

    this.validateForm = this.fb.group({
      groupUuid: ['', [Validators.required]],
      name: ['', [Validators.required]],
      desc: [''],
      requestPath: ['', [Validators.required, Validators.pattern('^[^?]*')]],
      service: ['', [Validators.required]],
      proxyPath: [''],
      timeout: [10000, [Validators.required]],
      retry: [0, [Validators.required]],
      enableWebsocket: [false],
      templateUuid: ['']
    })
  }

  ngOnInit (): void {
    this.getServiceList()
    this.getPluginTemplateList()
    if (this.baseInfo.allParamsInfo.apiId) {
      this.appConfigService.reqFlashBreadcrumb([
        { title: 'API管理', routerLink: 'router/api/group/list' },
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

    this.proxyHeaderTableBody[3].btns[0].disabledFn = () => {
      return this.nzDisabled
    }
    this.proxyHeaderTableBody[3].btns[0].click = (item:any) => {
      this.openDrawer('proxyHeader', item.data)
    }
    this.proxyHeaderTableBody[3].btns[1].disabledFn = () => {
      return this.nzDisabled
    }
  }

  ngAfterViewInit () {
    this.proxyHeaderTableBody[0].title = this.optTypeTranslateTpl
  }

  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  // 当编辑api时，需要获取api信息
  getApiMessage () {
    this.api.get('router', { uuid: this.apiUuid }).subscribe((resp) => {
      if (resp.code === 0) {
        setFormValue(this.validateForm, resp.data.api)
        // this.validateForm.controls['groupUuid'].setValue(resp.data.api.groupUuid)
        // this.validateForm.controls['name'].setValue(resp.data.api.name)
        // this.validateForm.controls['desc'].setValue(resp.data.api.desc)
        this.validateForm.controls['requestPath'].setValue(resp.data.api.requestPath.slice(1))
        // this.validateForm.controls['service'].setValue(resp.data.api.service)
        // this.validateForm.controls['proxyPath'].setValue(resp.data.api.proxyPath.slice(1))
        // this.validateForm.controls['timeout'].setValue(resp.data.api.timeout)
        // this.validateForm.controls['retry'].setValue(resp.data.api.retry)
        // this.validateForm.controls['enableWebsocket'].setValue(resp.data.api.enableWebsocket)
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
  getHeaderList () {
    this.api.get('router/groups').subscribe((resp:{code:number, data:ApiGroup, msg:string}) => {
      if (resp.code === 0) {
        const tempList:ApiGroupsData[] = []
        for (const index in resp.data.root.groups) {
          this.firstLevelList.push(resp.data.root.groups[index].uuid)
          if (!resp.data.root.groups[index].children || resp.data.root.groups[index].children.length === 0) {
            resp.data.root.groups[index]['disabled'] = true
          } else {
            resp.data.root.groups[index]['selectable'] = false
            tempList.push(resp.data.root.groups[index])
          }
        }
        this.headerList = this.transferHeader(tempList)
        if (this.groupUuid && this.firstLevelList.indexOf(this.groupUuid) === -1) {
          this.validateForm.controls['groupUuid'].setValue(this.baseInfo.allParamsInfo.apiGroupId)
        }
      } else {
        this.message.error(resp.msg || '获取列表数据失败!')
      }
    })
  }

  transferHeader (header:ApiGroupsData[]):NzTreeNodeOptions[] {
    const resList:NzTreeNodeOptions[] = []
    for (const index in header) {
      const res:NzTreeNodeOptions = {
        title: header[index].name,
        key: header[index].uuid,
        uuid: header[index].uuid,
        isDelete: header[index].isDelete
      }
      if (!header[index].children || header[index].children.length === 0) {
        res['isLeaf'] = true
      } else {
        res.children = this.transferHeader(header[index].children)
      }
      resList.push(res)
    }
    return resList
  }

  nzTreeClick (value: any) {
    if (value.node.origin.selectable === false) {
      value.node.origin.expanded = !value.node.origin.expanded
    }
    this.headerList = [...this.headerList]
  }

  // 获取上游服务列表
  getServiceList () {
    this.api.get('service/enum').subscribe((resp: any) => {
      if (resp.code === 0) {
        this.serviceList = []
        for (const item of resp.data.list) {
          this.serviceList = [...this.serviceList, { label: item, value: item }]
        }
      } else {
        this.message.error(resp.msg || '获取列表数据失败!')
      }
    })
  }

  getPluginTemplateList () {
    this.api.get('plugin/template/enum').subscribe((resp: {code:number, data:{templates:Array<{uuid:string, name:string}>}, msg:string}) => {
      if (resp.code === 0) {
        this.pluginTemplateList = resp.data.templates.map((item) => {
          const data = { label: item.name, value: item.uuid }
          return data
        })
      } else {
        this.message.error(resp.msg || '获取列表数据失败!')
      }
    })
  }

  updateAllChecked (): void {
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

  initCheckbox (): void {
    for (const index in this.methodList) {
      if (
        this.createApiForm.method.indexOf(this.methodList[index].label) !== -1
      ) {
        this.methodList[index].checked = true
      }
    }
  }

  updateSingleChecked (): void {
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

  openDrawer (type: string, data?: any) {
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
  backToList () {
    this.router.navigate(['/', 'router', 'api', 'group', 'list'])
  }

  // 提交api数据
  saveApi () {
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
        this.api.put('router', {
          name: this.validateForm.controls['name'].value,
          uuid: this.createApiForm.uuid,
          groupUuid: this.validateForm.controls['groupUuid'].value,
          desc: this.validateForm.controls['desc'].value,
          requestPath: '/' + this.validateForm.controls['requestPath'].value,
          service: this.validateForm.controls['service'].value,
          method: this.createApiForm.method,
          proxyPath: this.validateForm.controls['proxyPath'].value,
          timeout: Number(this.validateForm.controls['timeout'].value),
          retry: Number(this.validateForm.controls['retry'].value),
          enableWebsocket: this.validateForm.controls['enableWebsocket'].value || false,
          templateUuid: this.validateForm.controls['templateUuid'].value || '',
          proxyHeader: this.createApiForm.proxyHeader,
          match: this.createApiForm.match
        }, { uuid: this.apiUuid }).subscribe(resp => {
          if (resp.code === 0) {
            this.backToList()
            this.message.success(resp.msg || '修改成功！', { nzDuration: 1000 })
          } else {
            this.message.error(resp.msg || '修改失败!')
          }
        })
      } else {
        this.api.post('router', {
          ...this.validateForm.value,
          uuid: this.createApiForm.uuid,
          groupUuid: this.validateForm.controls['groupUuid'].value,
          method: this.createApiForm.method,
          proxyHeader: this.createApiForm.proxyHeader,
          match: this.createApiForm.match,
          requestPath: '/' + this.validateForm.controls['requestPath'].value,
          proxyPath: '/' + this.validateForm.controls['proxyPath'].value
        }).subscribe(resp => {
          if (resp.code === 0) {
            this.message.success(resp.msg || '添加成功！', { nzDuration: 1000 })
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

  requestPathChange () {
    if (!this.validateForm.controls['proxyPath'].value && this.validateForm.controls['requestPath'].value) {
      this.validateForm.controls['proxyPath'].setValue('/' + this.validateForm.controls['requestPath'].value)
    }
  }

  checkTimeout () {
    if (
      this.validateForm.controls['timeout'].value !== null &&
      this.validateForm.controls['timeout'].value < 1
    ) {
      this.validateForm.controls['timeout'].setValue(1)
    }
  }

  // 保存转发上游请求头数据时，如果是新建数据，直接加入tableList，如果是编辑数据，需要删除原先同key的数据再保存
  saveProxyHeader (proxyRef: ApiManagementProxyComponent): void {
    let proxyValid:boolean = false
    if (proxyRef.validateProxyHeaderForm.controls['optType'].value === 'DELETE') {
      proxyValid = !!proxyRef.validateProxyHeaderForm.controls['key'].value
    } else {
      proxyValid = proxyRef.validateProxyHeaderForm.valid
    }
    if (proxyValid) {
      if (this.proxyEdit) {
        for (const index in this.createApiForm.proxyHeader) {
          if (this.createApiForm.proxyHeader[index].key === this.editData.key && this.createApiForm.proxyHeader[index].optType === this.editData.optType && this.createApiForm.proxyHeader[index].value === this.editData.value) {
            this.createApiForm.proxyHeader.splice(Number(index), 1)
            break
          }
        }
      }
      this.createApiForm.proxyHeader = [{ optType: proxyRef.validateProxyHeaderForm.controls['optType'].value, key: proxyRef.validateProxyHeaderForm.controls['key'].value, value: proxyRef.validateProxyHeaderForm.controls['value'].value }, ...this.createApiForm.proxyHeader]
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

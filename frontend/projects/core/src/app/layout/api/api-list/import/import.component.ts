import { Component, Output, EventEmitter, TemplateRef, ViewChild, AfterViewInit } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { SelectOption } from 'eo-ng-select'
import { THEAD_TYPE } from 'eo-ng-table'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { NzFormatEmitEvent } from 'ng-zorro-antd/tree'
import { NzUploadFile } from 'ng-zorro-antd/upload'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { ApiGroup, ApiGroupsData, EmptyHttpResponse } from 'projects/core/src/app/constant/type'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { EO_TBODY_TYPE } from 'projects/eo-ng-apinto-table/src/public-api'
import { MODAL_NORMAL_SIZE } from 'projects/core/src/app/constant/app.config'
import { apiImportCheckResultTableHeadName, apiImportCheckResultTableBody } from '../../types/conf'
import { APIImportData, APIImportItem } from '../../types/types'

type ResultItem = APIImportItem & {checked?:boolean, disabled?:boolean, statusString?:string}

@Component({
  selector: 'eo-ng-api-import',
  templateUrl: './import.component.html',
  styles: [
  ]
})
export class ApiImportComponent implements AfterViewInit {
  @ViewChild('importContentTpl', { read: TemplateRef, static: true }) importContentTpl: TemplateRef<any> | undefined
  @ViewChild('methodTpl', { read: TemplateRef, static: true }) methodTpl: TemplateRef<any> | undefined
  @Output() flashList:EventEmitter<any> = new EventEmitter()
  modalRef:NzModalRef | undefined
  groupList:any[]= []
  upstreamList:SelectOption[]= []
  importFormPage:boolean = true
  fileList: NzUploadFile[] = [];
  authFile:NzUploadFile|undefined
  fileError:boolean = false
  token:string = ''
  resultMap:Map<number, any> = new Map()
  resultTableThead:THEAD_TYPE[] = [...apiImportCheckResultTableHeadName]
  resultTableTbody:EO_TBODY_TYPE[] = [...apiImportCheckResultTableBody]

  apisSet:Set<number> = new Set()
  resultList:Array<ResultItem> = []

  autoTips: Record<string, Record<string, string>> = defaultAutoTips

  validateForm:FormGroup = new FormGroup({})
  nzDisabled:boolean = false
  checkButtonLoading:boolean = false
  importBtnLoading:boolean = false
  constructor (
    private modalService:EoNgFeedbackModalService,
    private message: EoNgFeedbackMessageService,
    private api:ApiService,
    private fb: UntypedFormBuilder) {
  }

  ngAfterViewInit ():void {
    // 表格checkbox
    this.resultTableThead[0].click = (item:any) => {
      this.changeApisSet(item, 'all')
    }
    this.resultTableTbody[0].click = (item:any) => {
      this.changeApisSet(item)
    }
    // 表格name支持修改且为必填项
    this.resultTableTbody[2].check = (value:any) => {
      return !!value
    }
    this.resultTableTbody[3].title = this.methodTpl
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  nzTreeClick (value:NzFormatEmitEvent) {
    if (value.node!.origin.selectable === false) {
      value.node!.origin.expanded = !value.node!.origin.expanded
    }
    this.groupList = [...this.groupList]
  }

  openDrawer () {
    this.importFormPage = true
    this.token = ''
    this.getGroupList()
    this.getUpstreamList()
    this.fileList = []
    this.validateForm = this.fb.group({
      file: [null, [Validators.required]],
      group: ['', [Validators.required]],
      upstream: ['', [Validators.required]],
      requestPrefix: ['', [Validators.pattern('^[^?]*')]]
    })
    this.modalRef = this.modalService.create({
      nzTitle: '导入swagger文件',
      nzWidth: MODAL_NORMAL_SIZE,
      nzContent: this.importContentTpl,
      nzWrapClassName: 'import-drawer-wrap',
      nzOkText: '确定',
      nzCancelText: '取消',
      nzOkDisabled: this.nzDisabled,
      nzFooter: [
        {
          label: '取消',
          show: () => {
            return this.importFormPage
          },
          onClick: () => {
            this.modalRef?.close()
          }
        },
        {
          label: '查重',
          type: 'primary',
          show: () => {
            return this.importFormPage
          },
          disabled: () => {
            return this.nzDisabled
          },
          onClick: () => {
            this.checkConflict()
          },
          loading: () => {
            return this.checkButtonLoading
          }
        },
        {
          label: '返回上级',
          show: () => {
            return !this.importFormPage
          },
          onClick: () => {
            this.importFormPage = true
            this.resultList = []
            this.token = ''
          }
        },
        {
          label: '提交',
          type: 'primary',
          show: () => {
            return !this.importFormPage
          },
          disabled: () => {
            return !this.apisSet.size
          },
          onClick: () => {
            this.importApis()
          },
          loading: () => {
            return this.importBtnLoading
          }
        }
      ]
    })
  }

  // 获取API分组参数
  getGroupList () {
    this.api.get('router/groups').subscribe((resp:{code:number, data:ApiGroup, msg:string}) => {
      if (resp.code === 0) {
        this.groupList = []
        const tempList:ApiGroupsData[] = [...resp.data.root.groups]
        this.groupList = this.transferHeader(tempList)
      }
    })
  }

  // 将数据处理成树选择器需要的参数格式
  transferHeader (group:any): SelectOption[] {
    for (const index in group) {
      group[index]['title'] = group[index].name
      group[index]['key'] = group[index].uuid
      if (!group[index].children || group[index].children.length === 0) {
        group[index]['isLeaf'] = true
      } else {
        group[index].children = this.transferHeader(group[index].children)
      }
    }
    return group
  }

  // 获取上游服务列表
  getUpstreamList () {
    this.api.get('common/provider/Service').subscribe((resp:any) => {
      if (resp.code === 0) {
        this.upstreamList = [...resp.data.Service.map((x:any) => ({ label: x.title, value: x.name }))]
      }
    })
  }

  // 手动上传文件
  beforeUpload = (file: NzUploadFile): boolean => {
    this.fileList = []
    this.fileList = this.fileList.concat(file)
    this.authFile = file
    this.fileError = this.fileList.length === 0
    return false
  }

  // 移除文件
  removeFile () {
    this.fileList = []
    this.authFile = undefined
    this.fileError = true
    return true
  }

  checkConflict () {
    this.validateForm.controls['file'].setValue(this.fileList[0])
    if (this.validateForm.valid) {
      this.apisSet = new Set()
      this.resultList = []
      const formData = new FormData()
      formData.append('file', this.fileList[0] as any)
      formData.append('group', this.validateForm.controls['group'].value)
      formData.append('upstream', this.validateForm.controls['upstream'].value)
      formData.append('request_prefix', this.validateForm.controls['requestPrefix'].value)
      this.checkButtonLoading = true
      this.api.post('router/import', formData).subscribe((resp:{code:number, data:{apis:Array<APIImportItem>, token:string}, msg:string}) => {
        this.checkButtonLoading = false
        if (resp.code === 0) {
          this.importFormPage = false
          const validArray = resp.data.apis.filter((value) => {
            return value.status === 'normal'
          })
          this.resultList = resp.data.apis

          for (const api of this.resultList) {
            api['disabled'] = api.status !== 'normal'
            // 这里全选是因为要让表头上的勾选框被选中，实际提交会过滤掉disabled的api
            api['checked'] = validArray.length > 0
            api['statusString'] = api.status === 'normal' ? '正常' : (api.status === 'conflict' ? '冲突' : '无效path')
            this.resultMap.set(api.id, api)
            if (!api['disabled']) { this.apisSet.add(api.id) }
          }
          this.token = resp.data.token
        }
      })
    } else {
      if (this.validateForm.controls['file'].invalid) {
        this.fileError = true
      }
      Object.values(this.validateForm.controls).forEach(control => {
        if (control.invalid) {
          control.markAsDirty()
          control.updateValueAndValidity({ onlySelf: true })
        }
      })
    }
  }

  // 不插入数据
  nzCheckAddRow = () => {
    return false
  }

  // 勾选或取消勾选数据
  changeApisSet (item: boolean | ResultItem, type?:string) {
    if (type === 'all') {
      if (item) {
        for (const r of this.resultList.filter((x:ResultItem) => (!x.disabled))) {
          this.apisSet.add(r.id)
        }
      } else {
        this.apisSet = new Set()
      }
      return
    }

    // 被取消勾选
    if ((item as ResultItem)?.checked) {
      this.apisSet.delete((item as ResultItem).id)
    } else if (!(item as ResultItem).disabled) {
      // 被选中
      this.apisSet.add((item as ResultItem).id)
    }
  }

  // 导入apis
  importApis () {
    const submitApis: APIImportData[] = []
    for (const id of this.apisSet) {
      if (!this.resultMap.get(Number(id))['name']) {
        document.getElementsByTagName('input')[Number(id)]?.scrollIntoView()
        return
      }
      submitApis.push({ id: Number(id), name: this.resultMap.get(id)['name'], desc: this.resultMap.get(id).desc })
    }
    this.importBtnLoading = true

    this.api.put('router/import', { apis: submitApis, token: this.token }).subscribe((resp:EmptyHttpResponse) => {
      if (resp.code === 0) {
        this.message.success(resp.msg || 'API导入成功！', { nzDuration: 1000 })
        this.modalRef?.close()
        this.flashList.emit()
      }
      this.importBtnLoading = false
    })
  }
}

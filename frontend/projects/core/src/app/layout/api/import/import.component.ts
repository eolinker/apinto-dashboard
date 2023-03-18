/* eslint-disable no-useless-constructor */
/* eslint-disable dot-notation */
import { Component, Output, EventEmitter, TemplateRef, ViewChild } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { NzUploadFile } from 'ng-zorro-antd/upload'
import { MODAL_NORMAL_SIZE } from '../../../constant/app.config'
import { defaultAutoTips } from '../../../constant/conf'
import { ApiService } from '../../../service/api.service'

@Component({
  selector: 'eo-ng-api-import',
  templateUrl: './import.component.html',
  styles: [
  ]
})
export class ApiImportComponent {
  @ViewChild('importContentTpl', { read: TemplateRef, static: true }) importContentTpl: TemplateRef<any> | undefined
  @ViewChild('importFooterTpl', { read: TemplateRef, static: true }) importFooterTpl: TemplateRef<any> | undefined
  @ViewChild('methodTpl', { read: TemplateRef, static: true }) methodTpl: TemplateRef<any> | undefined
  @Output() flashList:EventEmitter<any> = new EventEmitter()
  drawerRef:NzModalRef | undefined
  groupList:Array<any> = []
  upstreamList:Array<{label:string, value:string}> = []
  importFormPage:boolean = true
  fileList: NzUploadFile[] = [];
  authFile:NzUploadFile|undefined
  fileError:boolean = false
  token:string = ''
  resultMap:Map<number, any> = new Map()
  resultTableThead:Array<any> = [
    {
      type: 'checkbox',
      key: 'checked',
      click: (item:any) => {
        this.changeApisSet(item, 'all')
      },
      resizeable: false
    },
    {
      title: '序号',
      width: 58
    },
    {
      title: 'API名称'
    },
    {
      title: '协议/方法',
      width: 99
    },
    {
      title: '请求路径'
    },
    {
      title: '描述'
    },
    {
      title: '状态',
      width: 70,
      resizeable: false
    }
  ]

  resultTableTbody:Array<any> = [
    {
      key: 'checked',
      type: 'checkbox',
      click: (item:any) => {
        this.changeApisSet(item)
      }
    },
    {
      key: 'id',
      styleFn: (item:any) => {
        if (item.disabled) {
          return 'color:var(--disabled-text-color)'
        }
        return ''
      }
    },
    {
      key: 'name',
      type: 'input',
      checkMode: 'change',
      check: (value:any) => {
        return !!value
      },
      errorTip: '必填项',
      placeholder: '请输入'
    },
    {
      key: 'method',
      styleFn: (item:any) => {
        if (item.disabled) {
          return 'color:var(--disabled-text-color)'
        }
        return ''
      }
    },
    {
      key: 'path',
      styleFn: (item:any) => {
        if (item.disabled) {
          return 'color:var(--disabled-text-color)'
        }
        return ''
      }
    },
    {
      key: 'desc',
      styleFn: (item:any) => {
        if (item.disabled) {
          return 'color:var(--disabled-text-color)'
        }
        return ''
      }
    },
    {
      key: 'statusString',
      styleFn: (item:any) => {
        if (item.disabled) {
          return 'color:var(--disabled-text-color)'
        }
        return ''
      }
    }
  ]

  apisSet:Set<number> = new Set()
  resultList:Array<any> = []

  autoTips: Record<string, Record<string, string>> = defaultAutoTips

  validateForm:FormGroup = new FormGroup({})
  nzDisabled:boolean = false
  constructor (
    private modalService:EoNgFeedbackModalService,
    private message: EoNgFeedbackMessageService,
    private api:ApiService,
    private fb: UntypedFormBuilder) {
  }

  ngAfterViewInit ():void {
    this.resultTableTbody[3].title = this.methodTpl
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  nzTreeClick (value:any) {
    if (value.node.origin.selectable === false) {
      value.node.origin.expanded = !value.node.origin.expanded
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
      request_prefix: ['', [Validators.pattern('^[^?]*')]]
    })
    this.drawerRef = this.modalService.create({
      nzTitle: '导入swagger文件',
      nzWidth: MODAL_NORMAL_SIZE,
      nzContent: this.importContentTpl,
      nzWrapClassName: 'import-drawer-wrap',
      nzOkText: '确定',
      nzCancelText: '取消',
      nzOkDisabled: this.nzDisabled,
      nzFooter: this.importFooterTpl
    })
  }

  // 获取API分组参数
  getGroupList () {
    this.api.get('router/groups').subscribe((resp:any) => {
      if (resp.code === 0) {
        this.groupList = []
        for (const index in resp.data.root.groups) {
          if (resp.data.root.groups[index].children && resp.data.root.groups[index].children.length > 0) {
            resp.data.root.groups[index].selectable = false
            this.groupList.push(resp.data.root.groups[index])
          }
        }
        this.groupList = this.transferHeader(this.groupList)
      } else {
        this.message.error(resp.msg || '获取列表数据失败!')
      }
    })
  }

  // 将数据处理成树选择器需要的参数格式
  transferHeader (group:any) {
    for (const index in group) {
      group[index].title = group[index].name
      group[index].key = group[index].uuid
      if (!group[index].children || group[index].children.length === 0) {
        group[index].isLeaf = true
      } else {
        group[index].children = this.transferHeader(group[index].children)
      }
    }
    return group
  }

  // 获取上游服务列表
  getUpstreamList () {
    this.api.get('service/enum').subscribe((resp:any) => {
      if (resp.code === 0) {
        this.upstreamList = []
        for (const item of resp.data.list) {
          this.upstreamList = [...this.upstreamList, { label: item, value: item }]
        }
      } else {
        this.message.error(resp.msg || '获取列表数据失败!')
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
      formData.append('request_prefix', this.validateForm.controls['request_prefix'].value)
      this.api.post('router/import', formData).subscribe((resp:{code:number, data:{apis:Array<{id:number, name:string, method:string, path:string, desc:string, status:string, [key:string]:any}>, token:string}, msg:string}) => {
        if (resp.code === 0) {
          this.importFormPage = false
          const validArray = resp.data.apis.filter((value) => {
            return value.status === 'normal'
          })
          for (const api of resp.data.apis) {
            api['disabled'] = api.status !== 'normal'
            api['checked'] = validArray.length > 0
            switch (api.status) {
              case 'normal':
                api['statusString'] = '正常'
                break
              case 'conflict':
                api['statusString'] = '冲突'
                break
              case 'invalidPath':
                api['statusString'] = '无效path'
                break
            }
            this.resultMap.set(api.id, api)
            if (!api['disabled']) { this.apisSet.add(api.id) }
          }
          this.resultList = resp.data.apis
          this.token = resp.data.token
        } else {
          this.message.error(resp.msg || '查重失败，请重试！')
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
  changeApisSet (item: any, type?:string) {
    if (type === 'all') {
      if (item) {
        for (const index in this.resultList) {
          if (!this.resultList[index].disabled) { this.apisSet.add(this.resultList[index].id) }
        }
      } else {
        this.apisSet = new Set()
      }
    } else {
    // 被取消勾选
      if (item?.checked) {
        this.apisSet.delete(item.id)
      } else if (!item.disabled) {
      // 被选中
        this.apisSet.add(item.id)
      }
    }
  }

  // 导入apis
  importApis () {
    const submitApis: Array<{id:number, name:string, desc:string}> = []
    for (const id of this.apisSet) {
      if (!this.resultMap.get(Number(id))['name']) {
        document.getElementsByTagName('input')[Number(id)]?.scrollIntoView()
        return
      }
      submitApis.push({ id: Number(id), name: this.resultMap.get(id)['name'], desc: this.resultMap.get(id).desc })
    }
    this.api.put('router/import', { apis: submitApis, token: this.token }).subscribe((resp:{code:number, data:{}, msg:string}) => {
      if (resp.code === 0) {
        this.message.success(resp.msg || 'API导入成功！', { nzDuration: 1000 })
        this.drawerRef?.close()
        this.flashList.emit()
        return true
      } else {
        this.message.error(resp.msg || 'API导入失败，请重试！')
        return false
      }
    })
  }
}

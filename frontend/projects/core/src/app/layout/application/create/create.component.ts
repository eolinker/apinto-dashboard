/* eslint-disable dot-notation */
/* eslint-disable camelcase */
import { Component, Input, OnInit } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'
import { EO_TBODY_TYPE } from 'projects/eo-ng-apinto-table/src/public-api'
import { defaultAutoTips } from '../../../constant/conf'
import { ArrayItemData, EmptyHttpResponse, RandomId } from '../../../constant/type'
import { customAttrTableBody, extraHeaderTableBody } from '../types/conf'
import { ApplicationData } from '../types/types'

@Component({
  selector: 'eo-ng-application-create',
  templateUrl: './create.component.html',
  styles: [
    `
      td input.ant-input {
        width: 140px !important;
      }
    `
  ]
})
export class ApplicationCreateComponent implements OnInit {
  @Input() editPage: boolean = false
  @Input() appId: string = ''
  validateForm: FormGroup = new FormGroup({})
  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  customAttrTableBody: EO_TBODY_TYPE[] = [...customAttrTableBody]
  extraHeaderTableBody:EO_TBODY_TYPE[]= [...extraHeaderTableBody]

  createApplicationForm: {
    customAttrList: ArrayItemData[]
    extraParamList: ArrayItemData[]
  } = {
    customAttrList: [],
    extraParamList: []
  }

  customAttrList:ArrayItemData[] = [
    { key: '', value: '', disabled: false }
  ]

  extraHeaderList: ArrayItemData[] = [
    { key: '', value: '', disabled: false }
  ]

  nzDisabled: boolean = false

  constructor (
    private message: EoNgFeedbackMessageService,
    private api: ApiService,
    private router: Router,
    private fb: UntypedFormBuilder,
    private appConfigService: AppConfigService
  ) {
    this.appConfigService.reqFlashBreadcrumb([
      { title: '应用管理', routerLink: 'application' },
      { title: '新建应用' }
    ])

    this.validateForm = this.fb.group({
      name: ['', [Validators.required]],
      id: [''],
      desc: ['']
    })
  }

  ngOnInit (): void {
    this.initTable()
    if (this.editPage) {
      this.getApplicationMessage()
      this.appConfigService.reqFlashBreadcrumb([
        { title: '应用管理', routerLink: 'application' },
        { title: '应用信息' }
      ])
    } else {
      this.getApplicationId()
    }
  }

  initTable () {
    this.customAttrTableBody[0].disabledFn = () => { return this.nzDisabled }
    this.customAttrTableBody[1].disabledFn = () => { return this.nzDisabled }
    this.customAttrTableBody[2].showFn = (item: any) => { return item === this.customAttrList[0] }
    this.customAttrTableBody[2].btns[0].click = (item: any) => { this.editArray(item.data, 'addCustom') }
    this.customAttrTableBody[2].btns[0].disabledFn = () => { return this.nzDisabled }
    this.customAttrTableBody[3].showFn = (item: any) => { return item !== this.customAttrList[0] }
    this.customAttrTableBody[3].btns[0].click = (item: any) => { this.editArray(item.data, 'addCustom') }
    this.customAttrTableBody[3].btns[0].disabledFn = () => { return this.nzDisabled }
    this.customAttrTableBody[3].btns[1].click = (item: any) => { this.editArray(item.data, 'deleteCustom') }
    this.customAttrTableBody[3].btns[1].disabledFn = () => { return this.nzDisabled }

    this.extraHeaderTableBody[0].disabledFn = () => { return this.nzDisabled }
    this.extraHeaderTableBody[1].disabledFn = () => { return this.nzDisabled }
    this.extraHeaderTableBody[2].showFn = (item: any) => { return item === this.extraHeaderList[0] }
    this.extraHeaderTableBody[2].btns[0].click = (item: any) => { this.editArray(item.data, 'addHeader') }
    this.extraHeaderTableBody[2].btns[0].disabledFn = () => { return this.nzDisabled }
    this.extraHeaderTableBody[3].showFn = (item: any) => { return item !== this.extraHeaderList[0] }
    this.extraHeaderTableBody[3].btns[0].click = (item: any) => { this.editArray(item.data, 'addHeader') }
    this.extraHeaderTableBody[3].btns[0].disabledFn = () => { return this.nzDisabled }
    this.extraHeaderTableBody[3].btns[1].click = (item: any) => { this.editArray(item.data, 'deleteHeader') }
    this.extraHeaderTableBody[3].btns[1].disabledFn = () => { return this.nzDisabled }
  }

  getApplicationMessage () {
    this.api
      .get('application', { appId: this.appId })
      .subscribe((resp: {code:number, data:{application:ApplicationData}, msg:string}) => {
        if (resp.code === 0) {
          this.createApplicationForm = resp.data.application
          this.validateForm.controls['name'].setValue(
            resp.data.application.name
          )
          this.validateForm.controls['id'].setValue(resp.data.application.id)
          this.validateForm.controls['desc'].setValue(
            resp.data.application.desc
          )
          if (resp.data.application.name === '匿名应用') {
            this.validateForm.controls['name'].disable()
          }
          this.validateForm.controls['id'].disable()

          this.customAttrList =
            this.createApplicationForm?.customAttrList?.length > 0
              ? this.createApplicationForm.customAttrList
              : [{ key: '', value: '', disabled: false }]
          this.extraHeaderList =
            this.createApplicationForm?.extraParamList?.length > 0
              ? this.createApplicationForm.extraParamList
              : [{ key: '', value: '', disabled: false }]
        }
      })
  }

  getApplicationId () {
    this.api.get('random/application/id')
      .subscribe((resp: {code:number, data:RandomId, msg:string}) => {
        if (resp.code === 0) {
          this.validateForm.controls['id'].setValue(resp.data.id)
        }
      })
  }

  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  // 添加或删除自定义属性或Header额外参数
  editArray (item: ArrayItemData, type: string) {
    switch (type) {
      case 'addCustom':
        for (const index in this.customAttrList) {
          if (this.customAttrList[index] === item) {
            this.customAttrList.splice(Number(index) + 1, 0, {
              key: '',
              value: '',
              disabled: false
            })
            break
          }
        }
        this.customAttrList = [...this.customAttrList]
        break
      case 'deleteCustom':
        for (const index in this.customAttrList) {
          if (this.customAttrList[index] === item) {
            this.customAttrList.splice(Number(index), 1)
            break
          }
        }
        this.customAttrList = [...this.customAttrList]
        break
      case 'addHeader':
        for (const index in this.extraHeaderList) {
          if (this.extraHeaderList[index] === item) {
            this.extraHeaderList.splice(Number(index) + 1, 0, {
              key: '',
              value: '',
              disabled: false
            })
            break
          }
        }
        this.extraHeaderList = [...this.extraHeaderList]
        break
      case 'deleteHeader':
        for (const index in this.extraHeaderList) {
          if (this.extraHeaderList[index] === item) {
            this.extraHeaderList.splice(Number(index), 1)
            break
          }
        }
        this.extraHeaderList = [...this.extraHeaderList]
        break
    }
  }

  // 保存鉴权，editPage = true时，表示页面为编辑页，false为新建页
  // custom_attr是创建和编辑鉴权时都会有的数据，需要将object转化为map发给后端
  // extra_header是编辑鉴权时才会有的数据，也需从Object转为map发送给后端
  saveApplication () {
    if (this.validateForm.valid) {
      this.createApplicationForm.customAttrList = this.customAttrList.filter(
        (item: ArrayItemData) => {
          return item.key && item.value
        }
      )

      this.createApplicationForm.extraParamList = this.extraHeaderList.filter(
        (item: ArrayItemData) => {
          return item.key && item.value
        }
      )

      if (!this.editPage) {
        this.api
          .post('application', {
            ...this.createApplicationForm,
            ...this.validateForm.value
          })
          .subscribe((resp: EmptyHttpResponse) => {
            if (resp.code === 0) {
              this.message.success(resp.msg || '添加成功', { nzDuration: 1000 })
              this.backToList()
            }
          })
      } else {
        this.api
          .put('application', {
            ...this.createApplicationForm,
            ...this.validateForm.value
          })
          .subscribe((resp: EmptyHttpResponse) => {
            if (resp.code === 0) {
              this.message.success(resp.msg || '修改成功', { nzDuration: 1000 })
              this.backToList()
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

  backToList () {
    this.router.navigate(['/', 'application'])
  }
}

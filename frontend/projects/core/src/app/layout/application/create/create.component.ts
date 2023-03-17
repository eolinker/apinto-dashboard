/* eslint-disable dot-notation */
/* eslint-disable no-useless-constructor */
/* eslint-disable camelcase */
import { Component, Input, OnInit } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'
import { defaultAutoTips } from '../../../constant/conf'

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

  createApplicationForm: {
    custom_attr_list: Array<{ key: string; value: string; disabled: boolean }>
    extra_param_list: Array<{
      key: string
      value: string
      disabled: boolean
      conflict?: string
      position?: string
    }>
  } = {
    custom_attr_list: [],
    extra_param_list: []
  }

  autoTips: Record<string, Record<string, string>> = defaultAutoTips

  customAttrTableBody: Array<any> = [
    {
      key: 'key',
      type: 'input',
      placeholder: '请输入Key',
      disabledFn: () => {
        return this.nzDisabled
      },
      checkMode: 'change',
      check: (item: any) => {
        return !item || /^[a-zA-Z_][a-zA-Z0-9_]*$/.test(item)
      },
      errorTip: '首字母必须为英文'
    },
    {
      key: 'value',
      type: 'input',
      placeholder: '请输入Value',
      disabledFn: () => {
        return this.nzDisabled
      }
    },
    {
      type: 'btn',
      showFn: (item: any) => {
        return item === this.customAttrList[0]
      },
      btns: [
        {
          title: '添加',
          click: (item: any) => {
            this.editArray(item.data, 'addCustom')
          },
          disabledFn: () => {
            return this.nzDisabled
          }
        }
      ]
    },
    {
      type: 'btn',
      showFn: (item: any) => {
        return item !== this.customAttrList[0]
      },
      btns: [
        {
          title: '添加',
          click: (item: any) => {
            this.editArray(item.data, 'addCustom')
          },
          disabledFn: () => {
            return this.nzDisabled
          }
        },

        {
          title: '减少',
          click: (item: any) => {
            this.editArray(item.data, 'deleteCustom')
          },
          disabledFn: () => {
            return this.nzDisabled
          }
        }
      ]
    }
  ]

  customAttrList: Array<{ key: string; value: string; disabled: boolean }> = [
    { key: '', value: '', disabled: false }
  ]

  extraHeaderTableBody: Array<any> = [
    {
      key: 'key',
      type: 'input',
      placeholder: '请输入Key',
      disabledFn: () => {
        return this.nzDisabled
      }
    },
    {
      key: 'value',
      type: 'input',
      placeholder: '请输入Value',
      disabledFn: () => {
        return this.nzDisabled
      }
    },
    {
      type: 'btn',
      showFn: (item: any) => {
        return item === this.extraHeaderList[0]
      },
      btns: [
        {
          title: '添加',
          click: (item: any) => {
            this.editArray(item.data, 'addHeader')
          },
          disabledFn: () => {
            return this.nzDisabled
          }
        }
      ]
    },
    {
      type: 'btn',
      showFn: (item: any) => {
        return item !== this.extraHeaderList[0]
      },
      btns: [
        {
          title: '添加',
          click: (item: any) => {
            this.editArray(item.data, 'addHeader')
          },
          disabledFn: () => {
            return this.nzDisabled
          }
        },

        {
          title: '减少',
          click: (item: any) => {
            this.editArray(item.data, 'deleteHeader')
          },
          disabledFn: () => {
            return this.nzDisabled
          }
        }
      ]
    }
  ]

  extraHeaderList: Array<{ key: string; value: string; disabled: boolean }> = [
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

  getApplicationMessage () {
    this.api
      .get('application', { app_id: this.appId })
      .subscribe((resp: any) => {
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
            this.createApplicationForm?.custom_attr_list?.length > 0
              ? this.createApplicationForm.custom_attr_list
              : [{ key: '', value: '', disabled: false }]
          this.extraHeaderList =
            this.createApplicationForm?.extra_param_list?.length > 0
              ? this.createApplicationForm.extra_param_list
              : [{ key: '', value: '', disabled: false }]
        } else {
          this.message.error(resp.msg || '获取数据失败!')
        }
      })
  }

  getApplicationId () {
    this.api.get('random/application/id').subscribe((resp: any) => {
      if (resp.code === 0) {
        this.validateForm.controls['id'].setValue(resp.data.id)
      } else {
        this.message.error(resp.msg || '获取应用ID失败!')
      }
    })
  }

  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  // 添加或删除自定义属性或Header额外参数
  editArray (item: any, type: string) {
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
      this.createApplicationForm.custom_attr_list = this.customAttrList.filter(
        (item: any) => {
          return item.key && item.value
        }
      )

      this.createApplicationForm.extra_param_list = this.extraHeaderList.filter(
        (item: any) => {
          return item.key && item.value
        }
      )

      if (!this.editPage) {
        this.api
          .post('application', {
            ...this.createApplicationForm,
            ...this.validateForm.value
          })
          .subscribe((resp: any) => {
            if (resp.code === 0) {
              this.message.success(resp.msg || '添加成功', { nzDuration: 1000 })
              this.backToList()
            } else {
              this.message.error(resp.msg || '添加失败!')
            }
          })
      } else {
        this.api
          .put('application', {
            ...this.createApplicationForm,
            ...this.validateForm.value
          })
          .subscribe((resp: any) => {
            if (resp.code === 0) {
              this.message.success(resp.msg || '修改成功', { nzDuration: 1000 })
              this.backToList()
            } else {
              this.message.error(resp.msg || '修改失败!')
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

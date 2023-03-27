/* eslint-disable dot-notation */
import { Component, Input, OnInit } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { SelectOption } from 'eo-ng-select'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { setFormValue } from 'projects/core/src/app/constant/form'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'
import { PluginData } from '../types/types'

@Component({
  selector: 'eo-ng-deploy-plugin-create',
  templateUrl: './create.component.html',
  styles: [
  ]
})
export class DeployPluginCreateComponent implements OnInit {
  @Input() editPage: boolean = false
  @Input() name?:string
  validateForm: FormGroup = new FormGroup({})
  extendsList: SelectOption[] = []
  relysList:SelectOption[] = []
  nzDisabled:boolean = false
  nodesTableShow = false
  clusterCanBeCreated: boolean = false
  testFlag:boolean = false
  testPassAddr:string = '' // 通过测试的集群地址
  constructor (
    private message: EoNgFeedbackMessageService,
    private api: ApiService,
    private router: Router,
    private fb: UntypedFormBuilder,
    private appConfigService: AppConfigService) {
    this.validateForm = this.fb.group({
      name: ['', [Validators.required, Validators.pattern('^[a-zA-Z][a-zA-Z0-9_]*')]],
      extended: ['', [Validators.required]],
      rely: [''],
      desc: ['']
    })

    this.appConfigService.reqFlashBreadcrumb([{ title: '插件管理', routerLink: 'deploy/plugin' }, { title: '新建插件' }])
  }

  ngOnInit (): void {
    this.getExtendsList()
    this.getRelysList()
    if (this.editPage) {
      this.getPluginMessage()
      this.validateForm.controls['name'].disable()
      this.validateForm.controls['rely'].disable()
      this.validateForm.controls['extended'].disable()
    }
  }

  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  ngAfterViewInit ():void {
  }

  autoTips: Record<string, Record<string, string>> = defaultAutoTips

  getPluginMessage () {
    this.api
      .get('plugin', { name: this.name })
      .subscribe(
        (resp: {
          code: number
          data: { plugin:PluginData }
          msg: string
        }) => {
          if (resp.code === 0) {
            this.appConfigService.reqFlashBreadcrumb([
              {
                title: '插件管理',
                routerLink: 'deploy/plugin'
              },
              { title: '编辑插件' }
            ])

            setFormValue(this.validateForm, {
              name: resp.data.plugin!.name,
              rely: resp.data.plugin!.rely || '',
              extended: resp.data.plugin!.extended,
              desc: resp.data.plugin!.desc || ''
            })
          } else {
            this.message.error(resp.msg || '获取数据失败!')
          }
        }
      )
  }

  getExtendsList () {
    this.api.get('plugin/extendeds').subscribe((resp:{code:number, data:{extendeds:string[]}, msg:string}) => {
      if (resp.code === 0) {
        this.extendsList = resp.data.extendeds.map(
          (name: string) => ({
            label: name,
            value: name
          })
        )
        this.validateForm
          .controls['extended']
          .setValue(this.extendsList[0]?.value)
        this.validateForm.controls['extended'].updateValueAndValidity({
          onlySelf: true
        })
      } else {
        this.message.error(resp.msg || '获取列表数据失败！')
      }
    })
  }

  getRelysList () {
    this.api.get('basic/info/plugins').subscribe((resp:{code:number, data:{plugins:Array<{name:string, extended:string}>}, msg:string}) => {
      if (resp.code === 0) {
        this.relysList = resp.data.plugins.map(
          (plugins: { name: string; extended: string }) => ({
            label: plugins.name,
            value: plugins.name
          })
        )
        this.validateForm
          .controls['rely']
          .setValue(this.relysList[0]?.value)
        this.validateForm.controls['rely'].updateValueAndValidity({
          onlySelf: true
        })
      } else {
        this.message.error(resp.msg || '获取列表数据失败！')
      }
    })
  }

  // 新建集群
  savePlugin () {
    if (this.validateForm.valid) {
      const params = {
        name: this.validateForm.controls['name'].value || '',
        extended: this.validateForm.controls['extended'].value || '',
        rely: this.validateForm.controls['rely'].value || '',
        desc: this.validateForm.controls['desc'].value || ''
      }
      if (!this.editPage) {
        this.api.post('plugin', params).subscribe((resp) => {
          if (resp.code === 0) {
            this.router.navigate(['/', 'deploy', 'plugin'])
          } else {
            this.message.error(resp.msg || '操作失败！')
          }
        })
      } else {
        this.api.put('plugin', { name: this.validateForm.controls['name'].value || '', desc: this.validateForm.controls['desc'].value || '' }).subscribe((resp) => {
          if (resp.code === 0) {
            this.router.navigate(['/', 'deploy', 'plugin'])
          } else {
            this.message.error(resp.msg || '操作失败！')
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

  // 取消新建集群
  backToList () {
    this.router.navigate(['/', 'deploy', 'plugin'])
  }
}

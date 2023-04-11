/* eslint-disable dot-notation */
import { Component, Input, OnInit } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { Router } from '@angular/router'
import { SelectOption } from 'eo-ng-select'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { setFormValue } from 'projects/core/src/app/constant/form'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { EoNgNavigationService } from 'projects/core/src/app/service/eo-ng-navigation.service'
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
  submitButtonLoading:boolean = false
  constructor (
    private api: ApiService,
    private router: Router,
    private fb: UntypedFormBuilder,
    private appConfigService: EoNgNavigationService) {
    this.validateForm = this.fb.group({
      name: ['', [Validators.required, Validators.pattern('^[a-zA-Z][a-zA-Z0-9_]*')]],
      extended: ['', [Validators.required]],
      rely: [''],
      desc: ['']
    })

    this.appConfigService.reqFlashBreadcrumb([{ title: '插件管理', routerLink: 'deploy/plugin' }, { title: '新建插件' }])
  }

  ngOnInit (): void {
    if (this.editPage) {
      this.getPluginMessage()
      this.validateForm.controls['name'].disable()
      this.validateForm.controls['rely'].disable()
      this.validateForm.controls['extended'].disable()
    } else {
      this.getExtendsList()
    }
    this.getRelysList()
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
            this.getExtendsList()
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
        if (!this.editPage) {
          this.validateForm
            .controls['extended']
            .setValue(this.extendsList[0]?.value)
        } else {
          this.extendsList = [
            ...this.extendsList,
            {
              label: this.validateForm.controls['extended'].value,
              value: this.validateForm.controls['extended'].value
            }]
        }
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
      this.submitButtonLoading = true
      if (!this.editPage) {
        this.api.post('plugin', params).subscribe((resp) => {
          this.submitButtonLoading = false
          if (resp.code === 0) {
            this.router.navigate(['/', 'deploy', 'plugin'])
          }
        })
      } else {
        this.api.put('plugin', { name: this.validateForm.controls['name'].value || '', desc: this.validateForm.controls['desc'].value || '' }).subscribe((resp) => {
          this.submitButtonLoading = false
          if (resp.code === 0) {
            this.router.navigate(['/', 'deploy', 'plugin'])
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

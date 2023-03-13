/* eslint-disable dot-notation */
import { Component, EventEmitter, Inject, Input, OnInit, Output } from '@angular/core'
import { FormGroup, UntypedFormBuilder } from '@angular/forms'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { Subscription } from 'rxjs'
import { AppServiceAdapter, APP_SERVICE_ADAPTER } from '../../../public-api'
import { API_SERVICE_ADAPTER, ApiServiceAdapter } from '../../constant/api-service-adapter'
import { defaultAutoTips } from '../../constant/conf'
import { EoNgMyValidators } from '../user-profile/user-profile.component'

@Component({
  selector: 'eo-ng-apinto-role-profile',
  templateUrl: './role-profile.component.html'
})
export class RoleProfileComponent implements OnInit {
  @Input() type:string =''
  @Input() roleId:string = ''
  @Input() nzDisabled:boolean = false // 权限控制
  @Output() eoCloseModal:EventEmitter<any> = new EventEmitter()
  @Input() accessLink:string = ''
  @Input() closeModal?:(value?:any)=>void
  validateForm:FormGroup = new FormGroup({})

  autoTips: Record<string, Record<string, string>> = defaultAutoTips

  rolesList:Array<any> = []
  editPage:boolean = false
  private subscription: Subscription = new Subscription()

  constructor (private message: EoNgFeedbackMessageService,
    private fb: UntypedFormBuilder,
    @Inject(API_SERVICE_ADAPTER) private apiService: ApiServiceAdapter,
     @Inject(APP_SERVICE_ADAPTER) private appService: AppServiceAdapter) {
    const { required } = EoNgMyValidators
    this.validateForm = this.fb.group({
      title: ['', [required]],
      desc: [''],
      access: [new Set(), [required]]
    })
  }

  ngOnInit (): void {
    if (this.type === 'editRole') {
      this.getRoleProfile()
    }
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  getRoleProfile () {
    this.apiService.get('role', { id: this.roleId }).subscribe((resp:any) => {
      if (resp.code === 0) {
        this.validateForm.get('title')?.setValue(resp.data.role.title)
        this.validateForm.get('desc')?.setValue(resp.data.role.desc)
        this.validateForm.get('access')?.setValue(new Set(resp.data.role.access))
      } else {
        this.message.error(resp.msg || '获取角色信息失败!')
      }
    })
  }

  backToList (value:any) {
    this.eoCloseModal.emit(value)
  }

  // 当表单通过验证后,根据父组件data传来的type提交表单
  saveRoleProfile () {
    if (this.validateForm.valid) {
      switch (this.type) {
        case 'addRole':
          return this.addRole()
        case 'editRole':
          return this.editRole()
      }
    } else {
      Object.values(this.validateForm.controls).forEach(control => {
        if (control.invalid) {
          control.markAsDirty()
          control.updateValueAndValidity({ onlySelf: true })
        }
      })
    }
    return false
  }

  addRole () {
    this.apiService.post('role', {
      title: this.validateForm.value.title,
      desc: this.validateForm.value.desc || '',
      access: this.validateForm.value.access.size > 0 ? Array.from(this.validateForm.value.access) : []
    }).subscribe((resp:any) => {
      if (resp.code === 0) {
        this.message.success(resp.msg || '创建角色成功!', { nzDuration: 1000 })
        this.eoCloseModal.emit(true)
        this.closeModal && this.closeModal()
      } else {
        this.message.error(resp.msg || '创建角色失败!')
      }
    })
  }

  editRole () {
    this.apiService.put('role', {
      title: this.validateForm.value.title,
      desc: this.validateForm.value.desc || '',
      access: this.validateForm.value.access.size > 0 ? Array.from(this.validateForm.value.access) : []
    }, { id: this.roleId }).subscribe((resp:any) => {
      if (resp.code === 0) {
        this.message.success(resp.msg || '编辑角色信息成功!', { nzDuration: 1000 })
        if (this.roleId === this.appService.getUserRoleId()) {
          this.subscription = this.appService.getMenuList().subscribe(() => {})
          this.subscription.unsubscribe()
        }
        this.eoCloseModal.emit(true)
        this.closeModal && this.closeModal()
      } else {
        this.message.error(resp.msg || '编辑角色信息失败!')
      }
    })
  }
}

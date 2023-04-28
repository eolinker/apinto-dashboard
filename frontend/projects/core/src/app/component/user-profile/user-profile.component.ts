/* eslint-disable dot-notation */
/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2022-09-21 22:19:44
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2022-09-21 23:14:13
 * @FilePath: /apinto/projects/eo-ng-apinto-user/src/lib/component/user-profile/user-profile.component.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { Component, Input, OnInit, Output, EventEmitter } from '@angular/core'
import { AbstractControl, FormGroup, UntypedFormBuilder, ValidatorFn, Validators } from '@angular/forms'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { NzSafeAny } from 'ng-zorro-antd/core/types'
import { defaultAutoTips } from '../../constant/conf'
import { ApiService } from '../../service/api.service'
import { EoNgNavigationService } from '../../service/app-config.service'
import { UserData } from '../../constant/type'
import { setFormValue } from '../../constant/form'

@Component({
  selector: 'eo-ng-apinto-user-profile',
  templateUrl: './user-profile.component.html'
})
export class UserProfileComponent implements OnInit {
  @Input() userId:string = '' // 用户id
  @Input() type:string = '' // 操作类型 editCurrentUser
  @Input() accessLink:string = ''
  @Input() nzDisabled:boolean = false // 权限
  @Output() eoCloseDrawer:EventEmitter<any> = new EventEmitter()
  validateForm:FormGroup = new FormGroup({})
  autoTips: Record<string, Record<string, string>> = defaultAutoTips

  rolesList:Array<any> = []
  editPage:boolean = false

  constructor (private message: EoNgFeedbackMessageService,
    private appService: EoNgNavigationService,
    private fb: UntypedFormBuilder,
     private apiService: ApiService) {
    const { required, email } = EoNgMyValidators
    this.validateForm = this.fb.group({
      userName: ['', [required, Validators.pattern('^[a-zA-Z][a-zA-Z0-9/_]*')]],
      nickName: ['', [required]],
      noticeUserId: [''],
      email: ['', [required, email]],
      role: [''],
      desc: ['']
    })
  }

  ngOnInit (): void {
    switch (this.type) {
      // 用户设置(修改当前用户信息)
      case 'editCurrentUser':
        this.getCurrentUserProfile()
        this.validateForm.controls['userName'].disable()
        this.validateForm.controls['role'].disable()
        break
      case 'editUser':
        this.getOtherUserProfile(this.userId)
        this.validateForm.controls['userName'].disable()
        if (this.userId === this.appService.getUserId()) {
          this.validateForm.controls['role'].disable()
        }
        this.getRolesList(false)
        break
      case 'addUser':
        this.getRolesList(false)
        break
    }
  }

  getCurrentUserProfile () {
    this.apiService.get('my/profile').subscribe((resp:any) => {
      if (resp.code === 0) {
        console.log(resp)
        setFormValue(this.validateForm, resp.data.profile)
        this.validateForm.controls['desc'].setValue(resp.data.describe)
        this.appService.setUserId(resp.data.profile.id)
      } else {
        this.message.error(resp.msg || '获取用户信息失败!')
      }
    })
  }

  getOtherUserProfile (id:string) {
    this.apiService.get('user/profile', { id: id || '' }).subscribe((resp:{code:number, data:{profile:UserData}, msg:string}) => {
      if (resp.code === 0) {
        this.validateForm.controls['userName'].setValue(resp.data.profile.user_name)
        this.validateForm.controls['nickName'].setValue(resp.data.profile.nick_name)
        this.validateForm.controls['noticeUserId'].setValue(resp.data.profile.notice_user_id)
        this.validateForm.controls['email'].setValue(resp.data.profile.email)
        this.validateForm.controls['role'].setValue(resp.data.profile.role_ids[0])
        this.validateForm.controls['desc'].setValue(resp.data.profile.desc)
      } else {
        this.message.error(resp.msg || '获取用户信息失败!')
      }
    })
  }

  // 获取角色id与title对应值, 传入list时,需要为该list的角色id与角色名匹配
  // 传入参数为true时,展示超管角色
  getRolesList (showM:boolean) {
    this.apiService.get('role/options').subscribe((resp:any) => {
      if (resp.code === 0) {
        this.rolesList = showM
          ? resp.data.roles
          : resp.data.roles.filter((item:any) => {
            return item.title !== '超级管理员'
          })
        for (const index in this.rolesList) {
          this.rolesList[index].label = this.rolesList[index].title
          this.rolesList[index].value = this.rolesList[index].id
        }
        this.rolesList.push({ label: '未分配', value: '' })
      } else {
        this.message.error(resp.msg || '获取角色列表失败!')
      }
    })
  }

  backToList (value:any) {
    this.closeModal(value)
    this.eoCloseDrawer.emit(value)
  }

  // 当表单通过验证后,根据父组件data传来的type提交表单
  saveUserProfile () {
    if (this.validateForm.valid) {
      switch (this.type) {
        case 'editCurrentUser':
          this.apiService.put('my/profile', {
            nick_name: this.validateForm.value.nickName,
            notice_user_id: this.validateForm.value.noticeUserId,
            email: this.validateForm.value.email,
            desc: this.validateForm.value.desc || ''
          }).subscribe((resp:any) => {
            if (resp.code === 0) {
              this.message.success(resp.msg || '修改成功!', { nzDuration: 1000 })
              this.closeModal()
            } else {
              this.message.error(resp.msg || '修改失败!')
            }
          })
          break
        case 'addUser':
          this.apiService.post('user/profile', {
            user_name: this.validateForm.value.userName,
            nick_name: this.validateForm.value.nickName,
            notice_user_id: this.validateForm.value.noticeUserId,
            email: this.validateForm.value.email,
            desc: this.validateForm.value.desc || '',
            role_ids: [this.validateForm.value.role]
          }).subscribe((resp:any) => {
            if (resp.code === 0) {
              this.message.success(resp.msg || '新增用户成功!', { nzDuration: 1000 })
              this.closeModal(true)
            } else {
              this.message.error(resp.msg || '新增用户失败!')
            }
          })
          break
        case 'editUser':
          this.apiService.put('user/profile', {
            user_name: this.validateForm.controls['userName'].value,
            nick_name: this.validateForm.value.nickName,
            notice_user_id: this.validateForm.value.noticeUserId,
            email: this.validateForm.value.email,
            desc: this.validateForm.value.desc || '',
            role_ids: [this.validateForm.value.role]
          }, { id: this.userId }).subscribe((resp:any) => {
            if (resp.code === 0) {
              this.message.success(resp.msg || '编辑用户信息成功!', { nzDuration: 1000 })
              this.closeModal(true)
            } else {
              this.message.error(resp.msg || '编辑用户信息失败!')
            }
          })
          break
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

  closeModal:(value?:any)=>void = () => {

  }
}

// current locale is key of the MyErrorsOptions
export type EoNgMyErrorsOptions = { 'zh-cn': string; en: string } & Record<string, NzSafeAny>;
export type EoNgMyValidationErrors = Record<string, EoNgMyErrorsOptions>;

export class EoNgMyValidators extends Validators {
  static override minLength (minLength: number): ValidatorFn {
    return (control: AbstractControl): EoNgMyValidationErrors | null => {
      if (Validators.minLength(minLength)(control) === null) {
        return null
      }
      return { minlength: { 'zh-cn': `最小长度为 ${minLength}`, en: `MinLength is ${minLength}` } }
    }
  }

  static override maxLength (maxLength: number): ValidatorFn {
    return (control: AbstractControl): EoNgMyValidationErrors | null => {
      if (Validators.maxLength(maxLength)(control) === null) {
        return null
      }
      return { maxlength: { 'zh-cn': `最大长度为 ${maxLength}`, en: `MaxLength is ${maxLength}` } }
    }
  }

  static roleAccess (control:AbstractControl): EoNgMyValidationErrors | null {
    const value = control.value
    if (value.size > 0) {
      return null
    } else {
      return { roleAccess: { 'zh-cn': '角色权限不能为空', en: 'Not Empty' } }
    }
  }
}

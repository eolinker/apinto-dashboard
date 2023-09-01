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
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { defaultAutoTips } from '../../constant/conf'
import { ApiService } from '../../service/api.service'
import { setFormValue } from '../../constant/form'
import { EoNgMyValidators } from '../../constant/eo-ng-validator'
import { EoNgNavigationService } from '../../service/eo-ng-navigation.service'

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
    this.getCurrentUserProfile()
    this.validateForm.controls['userName'].disable()
    this.validateForm.controls['role'].disable()
  }

  getCurrentUserProfile () {
    this.apiService.get('my/profile').subscribe((resp:any) => {
      if (resp.code === 0) {
        setFormValue(this.validateForm, resp.data.profile)
        this.validateForm.controls['desc'].setValue(resp.data.describe)
        this.appService.setUserId(resp.data.profile.id)
      } else {
        this.message.error(resp.msg || '获取用户信息失败!')
      }
    })
  }

  // 当表单通过验证后,根据父组件data传来的type提交表单
  saveUserProfile () {
    if (this.validateForm.valid) {
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

/* eslint-disable dot-notation */
import {
  Component,
  EventEmitter,
  Input,
  OnInit,
  Output
} from '@angular/core'
import {
  FormGroup,
  UntypedFormBuilder,
  UntypedFormControl
} from '@angular/forms'
import {
  EoNgFeedbackMessageService
} from 'eo-ng-feedback'
import { CryptoService } from '../../service/crypto.service'
import { defaultAutoTips } from '../../constant/conf'
import { ApiService } from '../../service/api.service'
import { EoNgMyValidators } from '../../constant/eo-ng-validator'

@Component({
  selector: 'eo-ng-apinto-reset-psw',
  templateUrl: './reset-psw.component.html'
})
export class ResetPswComponent implements OnInit {
  @Input() type:string = ''
  @Input() closeModal?:(value?:any)=>void
  @Output() closeDrawer: EventEmitter<any> = new EventEmitter()
  userName: string = '' // 用户账户
  validateForm: FormGroup = new FormGroup({})
  autoTips: Record<string, Record<string, string>> = defaultAutoTips

  constructor (
    private message: EoNgFeedbackMessageService,
    private crypto: CryptoService,
    private fb: UntypedFormBuilder,
    private apiService: ApiService
  ) {
    const { required, minLength, maxLength } = EoNgMyValidators
    this.validateForm = this.fb.group({
      userName: [{ value: '', disabled: true }, [required]],
      old: ['', [required, maxLength(32), minLength(8)]],
      password: [
        '',
        [required, maxLength(32), minLength(8), this.pswValidator]
      ],
      confirm: [
        '',
        [required, maxLength(32), minLength(8), this.confirmValidator]
      ]
    })
  }

  // 校验确认新密码和新密码是否一致
  confirmValidator = (
    control: UntypedFormControl
  ): { [s: string]: boolean } => {
    if (!control.value) {
      return { error: true, required: true }
    } else if (control.value !== this.validateForm.controls['password'].value) {
      return { confirm: true, error: true }
    }
    return {}
  }

  strength: number = 0
  // 校验新密码强度
  pswValidator = (control: UntypedFormControl): { [s: string]: boolean } => {
    this.strength = this.getPswStrength(control.value)
    if (!control.value) {
      return { error: true, required: true }
    } else if (this.strength === 1) {
      return { strength: true, error: true }
    } else {
      return {}
    }
  }

  // 获得密码强度, 当密码为 纯数字,小写字母,大写字母,特殊字符 中任三种组合时,密码强度为强
  // 任两种组合,密码强度为中
  // 只一种,密码强度为弱
  getPswStrength (value: string): number {
    const pswRegNum: RegExp = /[0-9]/
    const pswRegLowercase: RegExp = /[a-z]/
    const pswRegUppercase: RegExp = /[A-Z]/
    const pswRegSymbol: RegExp = /!@#$%^&*`~()-+=/
    let strength: number = 0
    if (pswRegNum.test(value)) {
      strength++
    }
    if (pswRegLowercase.test(value)) {
      strength++
    }
    if (pswRegUppercase.test(value)) {
      strength++
    }
    if (pswRegSymbol.test(value)) {
      strength++
    }

    return strength
  }

  ngOnInit (): void {
    this.validateForm.controls['userName'].setValue(this.userName)
  }

  resetPsw () :void {
    if (this.validateForm.valid) {
      this.apiService
        .post('my/password', {
          old: this.crypto.encryptByEnAES(this.validateForm.controls['userName'].value, this.validateForm.controls['old'].value),
          password: this.crypto.encryptByEnAES(this.validateForm.controls['userName'].value, this.validateForm.controls['password'].value)
        })
        .subscribe((resp: any) => {
          if (resp.code === 0) {
            this.message.success(resp.msg || '修改密码成功!')
            this.backToList(true)
          } else {
            this.message.error(resp.msg || '修改密码失败!')
          }
        })
    } else {
      Object.values(this.validateForm.controls).forEach((control) => {
        if (control.invalid) {
          control.markAsDirty()
          control.updateValueAndValidity({ onlySelf: true })
        }
      })
    }
  }

  backToList (value:any) {
    this.closeDrawer.emit(value)
    this.closeModal && this.closeModal()
  }

  validateConfirmPassword (): void {
    setTimeout(() =>
      this.validateForm.controls['confirm'].updateValueAndValidity()
    )
  }
}

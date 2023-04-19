/* eslint-disable brace-style */
/* eslint-disable dot-notation */
/* eslint-disable no-useless-constructor */
/* eslint-disable no-undef */
import { Component, ElementRef, OnInit, ViewChild } from '@angular/core'
import { FormGroup, FormBuilder, Validators } from '@angular/forms'
import { ActivatedRoute, Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { ApiService } from '../../../service/api.service'
import { EoNgNavigationService } from '../../../service/eo-ng-navigation.service'
import { CryptoService } from '../../../service/crypto.service'

@Component({
  selector: 'eo-ng-password',
  templateUrl: './password.component.html',
  styleUrls: ['./password.component.scss']
})
export class PasswordComponent implements OnInit {
  showCaptcha: boolean = false
  validateForm!: FormGroup
  verifyCode!: string
  isAutoFocus!: boolean
  loginLoading!: boolean
  isShowTooltip!: boolean
  // @ts-ignore
  routeQuery = this.route.queryParams.value
  @ViewChild('needAutoFocus') autoFocusInput!: ElementRef
  constructor (
    private api: ApiService,
    private fb: FormBuilder,
    private router: Router,
    private route: ActivatedRoute,
    private message: EoNgFeedbackMessageService,
    private appConfig: EoNgNavigationService,
    private crypto: CryptoService
  ) {}

  ngOnInit (): void {
    this.validateForm = this.fb.group({
      username: [null, [Validators.required]],
      password: [null, [Validators.required]]
    })
  }

  /* 使input 自动获取焦点  */
  ngAfterViewInit (): void {
    this.autoFocusInput.nativeElement.focus()
  }

  async login () {
    // await this.apikitService.cleanCookies()
    if (this.validateForm.valid) {
      this.loginLoading = true
      try {
        this.api
          .login({
            username: this.validateForm.controls['username'].value,
            password: this.crypto.encryptByEnAES(this.validateForm.controls['username'].value, this.validateForm.controls['password'].value),
            client: 1,
            type: 1,
            app_type: 4
          })
          .subscribe((resp: any) => {
            if (resp.code === 0) {
              this.appConfig.reqFlashMenu()
              this.message.create('success', '登录成功')
              const callbackUrl:string | null = this.route.snapshot.queryParams['callback']
              if (callbackUrl) {
                // this.router.navigate([callbackUrl])
              }
              else {
                // this.router.navigate([''])
              }
            }
          })
      } catch (err) {
        console.warn(err)
      } finally {
        this.loginLoading = false
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
}

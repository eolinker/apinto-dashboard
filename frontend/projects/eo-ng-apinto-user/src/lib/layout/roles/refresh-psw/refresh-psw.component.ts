import { Component, Inject, Input, OnInit } from '@angular/core'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { API_SERVICE_ADAPTER, ApiServiceAdapter } from '../../../constant/api-service-adapter'
import { APP_SERVICE_ADAPTER, AppServiceAdapter } from '../../../constant/app-service-adapter'

@Component({
  selector: 'eo-ng-apinto-roles-refresh-psw',
  templateUrl: './refresh-psw.component.html',
  styles: ['']
})
export class RefreshPswComponent implements OnInit {
  @Input() userId:string = '' // 需要重置密码的用户id
  @Input() closeModal?:(value?:any)=>void

  newRandomPsw: string = '' // 随机密码(从接口获取)
  constructor (
    private message: EoNgFeedbackMessageService,
    @Inject(API_SERVICE_ADAPTER) private apiService: ApiServiceAdapter,
    @Inject(APP_SERVICE_ADAPTER) public appService: AppServiceAdapter) { }

  ngOnInit (): void {
    this.generateNewPsw()
  }

  // 随机生成六位数密码(可含数字,符号,英文)
  generateNewPsw () {
    this.apiService.get('random/password/id').subscribe((resp: any) => {
      if (resp.code === 0) {
        this.newRandomPsw = resp.data.id
      } else {
        this.message.error(resp.msg || '生成随机密码失败！')
      }
    })
  }

  refreshPsw () {
    this.apiService
      .post('user/password-reset', {
        id: this.userId || '',
        password: this.newRandomPsw
      })
      .subscribe((resp: any) => {
        if (resp.code === 0) {
          this.closeModal && this.closeModal()
          this.message.success(resp.msg || '重置密码成功！', { nzDuration: 1000 })
        } else {
          this.message.error(resp.msg || '重置密码失败！')
        }
      })
  }
}

/* eslint-disable no-useless-constructor */
import { Injectable, Injector, TemplateRef } from '@angular/core'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { NzSingletonService } from 'ng-zorro-antd/core/services'
import { NzMessageDataOptions, NzMessageRef } from 'ng-zorro-antd/message'
import { Overlay } from '@angular/cdk/overlay'

@Injectable({
  providedIn: 'root'
})
export class EoNgMessageService extends EoNgFeedbackMessageService {
  constructor (
    nzSingletonService: NzSingletonService,
    overlay: Overlay,
    injector: Injector, private message: EoNgFeedbackMessageService) {
    super(nzSingletonService, overlay, injector)
  }

  override success (message:string|TemplateRef<any>, options?:NzMessageDataOptions):NzMessageRef {
    const option:NzMessageDataOptions = { ...options }
    if (!option.nzDuration) {
      option.nzDuration = 1000
    }
    return this.message.success(message, option)
  }
}

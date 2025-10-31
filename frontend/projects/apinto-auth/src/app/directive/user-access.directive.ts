/* eslint-disable no-useless-constructor */
import { Directive, ElementRef, EventEmitter, Input, Output, Renderer2 } from '@angular/core'
import { Subscription } from 'rxjs'
import { BaseInfoService } from '../service/base-info.service'

@Directive({
  selector: '[eoNgUserAccess]'
})
export class UserAccessDirective {
  @Input() eoNgUserAccess:string = ''
  @Output() disabledEdit:EventEmitter<any> = new EventEmitter()

  private userUpdateRight:boolean = false // 默认用户无权限编辑
  private userUpdateRightList:Array<string> = [] // 用户编辑权限路由列表
  private subscription: Subscription = new Subscription()

  constructor (
    private baseInfo:BaseInfoService,
    private el:ElementRef,
    private renderer:Renderer2) {
  }

  ngOnInit (): void {
    this.disableEdit()
  }

  ngDoCheck () {
    this.disableEdit()
  }

  ngOnDestroy () {
    this.subscription.unsubscribe()
  }

  disableEdit () {
    if (this.baseInfo.userInfoUpdated) {
      this.userUpdateRight = this.baseInfo.userModuleAccess === 'edit' // 用户在当前模块是否有编辑权限
      if (!this.userUpdateRight) {
        if (this.el.nativeElement.localName === 'eo-ng-dropdown' || this.el.nativeElement.localName === 'a') {
          this.renderer.setStyle(this.el.nativeElement, 'visibility', 'hidden')
        } else {
          this.renderer.setAttribute(this.el.nativeElement, 'disabled', 'true')
        }
        this.disabledEdit.emit(true)
      } else {
        if (this.el.nativeElement.localName === 'eo-ng-dropdown' || this.el.nativeElement.localName === 'a') {
          this.renderer.setStyle(this.el.nativeElement, 'visibility', 'visible')
        } else {
          this.renderer.removeAttribute(this.el.nativeElement, 'disabled')
        }
        this.disabledEdit.emit(false)
      }
    }
  }
}

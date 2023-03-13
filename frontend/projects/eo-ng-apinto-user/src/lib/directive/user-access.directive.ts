/* eslint-disable no-useless-constructor */
import { Directive, ElementRef, EventEmitter, Inject, Input, Output, Renderer2 } from '@angular/core'
import { Subscription } from 'rxjs'
import { APP_SERVICE_ADAPTER, AppServiceAdapter } from '../constant/app-service-adapter'

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
    @Inject(APP_SERVICE_ADAPTER) private appService:AppServiceAdapter,
    private el:ElementRef,
    private renderer:Renderer2) {
  }

  ngOnInit (): void {
    this.subscription = this.appService.repUpdateRightList().subscribe(() => {
      this.disableEdit()
    })
  }

  ngDoCheck () {
    this.disableEdit()
  }

  ngOnDestroy () {
    this.subscription.unsubscribe()
  }

  disableEdit () {
    if (this.appService.dataUpdated) {
      this.userUpdateRightList = this.appService.getUpdateRightsRouter()
      this.userUpdateRight = this.userUpdateRightList.indexOf(this.eoNgUserAccess) !== -1
      if (!this.userUpdateRight && this.eoNgUserAccess) {
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

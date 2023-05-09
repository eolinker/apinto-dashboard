/* eslint-disable dot-notation */
/* eslint-disable no-useless-constructor */
import { Directive, ElementRef, EventEmitter, Input, OnInit, Output, Renderer2 } from '@angular/core'
import { Router } from '@angular/router'
import { Subscription } from 'rxjs'
import { EoNgNavigationService } from '../service/eo-ng-navigation.service'

@Directive({
  selector: '[eoNgUserAccess]'
})
export class UserAccessDirective implements OnInit {
  @Input() eoNgUserAccess:string = ''
  @Input() viewAccess:boolean = false
  @Output() disabledEdit:EventEmitter<any> = new EventEmitter()
  oldAccessRouter:string = ''
  private userRight:boolean = false // 默认用户无权限或查看编辑
  private subscription: Subscription = new Subscription()
  private subscription1: Subscription = new Subscription()
  constructor (
    private navigationService:EoNgNavigationService,
    private el:ElementRef,
    private renderer:Renderer2,
    private router:Router) {
  }

  ngOnInit (): void {
    this.subscription = this.navigationService.repUpdateRightList().subscribe(() => {
      this.disableEdit()
    })
  }

  ngDoCheck () {
    this.disableEdit()
  }

  ngOnDestroy () {
    this.subscription.unsubscribe()
    this.subscription1.unsubscribe()
  }

  disableEdit () {
    if (this.navigationService.dataUpdated) {
      const moduleName = this.navigationService.routerNameMap.get(this.eoNgUserAccess) || this.eoNgUserAccess
      this.userRight = this.viewAccess ? !!this.navigationService.accessMap.get(moduleName) : this.navigationService.accessMap.get(moduleName) === 'edit'
      if (!this.userRight) {
        if (this.el.nativeElement.localName === 'eo-ng-dropdown' || this.el.nativeElement.localName === 'a') {
          this.renderer.setStyle(this.el.nativeElement, 'visibility', 'hidden')
        } else {
          this.disabledEdit.emit(true)
          this.renderer.setProperty(this.el.nativeElement, 'disabled', true)
        }
      } else {
        if (this.el.nativeElement.localName === 'eo-ng-dropdown' || this.el.nativeElement.localName === 'a') {
          this.renderer.setStyle(this.el.nativeElement, 'visibility', 'none')
        } else {
          this.disabledEdit.emit(false)
        }
      }
    }
  }
}

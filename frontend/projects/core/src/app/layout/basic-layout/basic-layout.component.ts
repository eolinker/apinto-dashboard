/* eslint-disable camelcase */
/* eslint-disable dot-notation */
/* eslint-disable no-useless-constructor */
/*
 * @Author:
 * @Date: 2022-07-12 00:19:11
 * @LastEditors:
 * @LastEditTime: 2022-07-29 02:56:25
 * @FilePath: /apinto/src/app/basic-layout/basic-layout.component.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { Component, OnInit } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgBreadcrumbOptions } from 'eo-ng-breadcrumb'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { MenuOptions } from 'eo-ng-menu'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'
import { Subscription } from 'rxjs'
import { ApiService } from '../../service/api.service'

@Component({
  selector: 'basic-layout',
  templateUrl: './basic-layout.component.html',
  styleUrls: ['./basic-layout.component.scss']
})
export class BasicLayoutComponent implements OnInit {
  sideMenuOptions :MenuOptions[] = []
  breadcrumbOptions: EoNgBreadcrumbOptions[] = []
  currentRouter:string = '' // 当前路由
  openMap:{[name:string]:boolean} = {}
  modalRef:NzModalRef | undefined
  showEmpty:boolean = false
  showSideLine:boolean = true

  authInfo:{title:string, infos:Array<{key:string, value:string}>}
  = { title: '', infos: [] }

  private subscription1: Subscription = new Subscription()
  private subscription2: Subscription = new Subscription()
  private subscription3: Subscription = new Subscription()
  private subscription4: Subscription = new Subscription()

  constructor (
              private message: EoNgFeedbackMessageService,
              private router: Router,
              private api:ApiService,
              private appConfigService: AppConfigService
  ) {
    this.subscription1 = this.appConfigService.repFlashBreadcrumb().subscribe((data:any) => {
      this.breadcrumbOptions = data
    })

    this.subscription2 = this.appConfigService.repFlashMenu().subscribe(() => {
      this.sideMenuOptions = [...this.appConfigService.getCurrentMenuList()]
      for (const menu of this.sideMenuOptions) {
        menu.open = this.openMap[menu.title as string]
      }
    })

    this.subscription3 = this.router.events.subscribe(() => {
      if (this.router.url !== this.currentRouter) {
        this.selectOrOpenMenu(this.router.url)
      }
    })
  }

  ngOnInit () {
    this.getSideMenu()
  }

  ngOnDestroy () {
    this.subscription1.unsubscribe()
    this.subscription2.unsubscribe()
    this.subscription3.unsubscribe()
    this.subscription4.unsubscribe()
  }

  getSideMenu () {
    this.subscription4 = this.appConfigService.getMenuList()
      .subscribe((res:MenuOptions[]) => {
        this.sideMenuOptions = [...res]
        for (const index in this.sideMenuOptions) {
          this.sideMenuOptions[index].openChange = (value:MenuOptions) => {
            this.openHandler(value['id']!)
          }
        }
        this.getAccess()
      })
  }

  updateAuth = () => {
    this.modalRef?.close()
    this.router.navigate(['/', 'auth-update'])
  }

  getAccess () {
    if (this.appConfigService.getUserAccess()) {
      this.showEmpty = true
      this.showSideLine = false
    } else {
      this.showEmpty = false
      this.showSideLine = true
      if (this.router.routerState.snapshot.url === '/' || this.router.routerState.snapshot.url === '/login') {
        this.router.navigate([this.appConfigService.getPageRoute()])
      }

      if (this.router.url !== this.currentRouter) {
        setTimeout(() => {
          this.selectOrOpenMenu(this.router.url)
        }, 0)
      }
    }
  }

  // 根据路由选中并打开对应menu
  selectOrOpenMenu (router:string):void {
    if (this.sideMenuOptions.length > 0) {
      for (const index in this.sideMenuOptions) {
        if (router.split('/')[1] === this.sideMenuOptions[index]['router']?.split('/')[0]) {
          if (this.sideMenuOptions[index].children?.length) {
            this.openHandler(this.sideMenuOptions[index]['id']!)
          }
          break
        }
      }
      this.sideMenuOptions = [...this.sideMenuOptions]
      this.currentRouter = router
    }
  }

  openHandler (id: string): void {
    for (const index in this.sideMenuOptions) {
      if (this.sideMenuOptions[index]['id'] !== id) {
        this.sideMenuOptions[index].open = false
      } else {
        this.sideMenuOptions[index].open = true
      }
      this.openMap[this.sideMenuOptions[index]['title'] as string] = !!this.sideMenuOptions[index].open
    }
  }

  goToGithub () {
    window.open('https://github.com/eolinker/apinto')
  }
}

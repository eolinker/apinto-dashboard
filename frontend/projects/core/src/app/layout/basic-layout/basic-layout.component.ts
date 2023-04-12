/* eslint-disable camelcase */
/* eslint-disable dot-notation */
/* eslint-disable no-useless-constructor */
/*
 * @Author: maggieyyy im.ymj@hotmail.com
 * @Date: 2022-07-12 00:19:11
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2022-07-29 02:56:25
 * @FilePath: /apinto/src/app/basic-layout/basic-layout.component.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { Component, OnInit } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgBreadcrumbOptions } from 'eo-ng-breadcrumb'
import { MenuOptions } from 'eo-ng-menu'
import { NzModalRef, NzModalService } from 'ng-zorro-antd/modal'
import { Subscription } from 'rxjs'
import { MODAL_SMALL_SIZE } from '../../constant/app.config'
import { ApiService } from '../../service/api.service'
import { EoNgNavigationService } from '../../service/eo-ng-navigation.service'
import { AuthInfoDetailComponent } from '../auth/info/detail/detail.component'

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
  seletedGuide:boolean = false
  authInfo:{title:string, infos:Array<{key:string, value:string}>}
  = { title: '', infos: [] }

  private subscription1: Subscription = new Subscription()
  private subscription2: Subscription = new Subscription()
  private subscription3: Subscription = new Subscription()
  private subscription4: Subscription = new Subscription()

  constructor (
              private router: Router,
              private api:ApiService,
              private navigationService: EoNgNavigationService,
              private modalService: NzModalService
  ) {
    this.subscription1 = this.navigationService.repFlashBreadcrumb().subscribe((data:any) => {
      this.breadcrumbOptions = data
    })

    this.subscription2 = this.navigationService.repFlashMenu().subscribe(() => {
      this.sideMenuOptions = [...this.navigationService.getCurrentMenuList()]
      for (const menu of this.sideMenuOptions) {
        menu.open = this.openMap[menu['titleString']! as string]
      }
    })

    this.subscription3 = this.router.events.subscribe(() => {
      if (this.router.url !== this.currentRouter) {
        this.selectOrOpenMenu(this.router.url)
      }
      this.seletedGuide = this.router.url.includes('guide')
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
    this.subscription4 = this.navigationService.getMenuList()
      .subscribe((res:MenuOptions[]) => {
        this.sideMenuOptions = [...res]
        for (const index in this.sideMenuOptions) {
          this.sideMenuOptions[index].openChange = (value:MenuOptions) => {
            this.openHandler(value['key']!)
          }
        }
        // this.getAccess()
      })
  }

  getAuthInfo () {
    if (this.navigationService.getUserAuthAccess()) {
      this.api.authGet('activation/info')
        .subscribe((resp:{code:number, data:{infos:Array<{key:string, value:string}>, title:string}, msg:string}) => {
          if (resp.code === 0) {
            this.authInfo = resp.data
          }
        })
    }
  }

  openAuthDialog () {
    this.modalRef = this.modalService.create({
      nzWrapClassName: 'auth-modal-header',
      nzTitle: `${this.authInfo.title}授权`,
      nzContent: AuthInfoDetailComponent,
      nzComponentParams: { eoInfos: this.authInfo.infos, updateAuth: this.updateAuth },
      nzClosable: true,
      nzFooter: null,
      nzWidth: MODAL_SMALL_SIZE
    })
  }

  updateAuth = () => {
    this.modalRef?.close()
    this.router.navigate(['/', 'auth-update'])
  }

  getAccess () {
    if (this.navigationService.getUserAccess()) {
      this.showEmpty = true
      this.showSideLine = false
    } else {
      this.showEmpty = false
      this.showSideLine = true
      if (this.router.routerState.snapshot.url === '/' || this.router.routerState.snapshot.url === '/login') {
        this.router.navigate([this.navigationService.getPageRoute()])
      }

      if (this.router.url !== this.currentRouter) {
        setTimeout(() => {
          this.selectOrOpenMenu(this.router.url)
        }, 0)
      }
    }
    this.getAuthInfo()
  }

  // 根据路由选中并打开对应menu
  selectOrOpenMenu (router:string):void {
    if (this.sideMenuOptions.length > 0) {
      for (const index in this.sideMenuOptions) {
        if (router.split('/')[1] === this.sideMenuOptions[index]['router']?.split('/')[0]) {
          if (this.sideMenuOptions[index].children?.length) {
            this.openHandler(this.sideMenuOptions[index]['key']!)
          }
          break
        }
      }
      this.sideMenuOptions = [...this.sideMenuOptions]
      this.currentRouter = router
    }
  }

  openHandler (key: string): void {
    for (const index in this.sideMenuOptions) {
      if (this.sideMenuOptions[index]['key'] !== key) {
        this.sideMenuOptions[index].open = false
      } else {
        this.sideMenuOptions[index].open = true
      }
      this.openMap[this.sideMenuOptions[index]['key'] as string] = !!this.sideMenuOptions[index].open
    }
  }

  goToGithub () {
    window.open('https://github.com/eolinker/apinto')
  }

  goToGuide () {
    this.router.navigate(['/', 'guide'])
  }
}

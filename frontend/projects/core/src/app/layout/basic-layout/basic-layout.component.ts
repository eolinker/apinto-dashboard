/* eslint-disable dot-notation */
/*
 * @Author: maggieyyy im.ymj@hotmail.com
 * @Date: 2022-07-12 00:19:11
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2023-04-13 23:40:55
 * @FilePath: /apinto/src/app/basic-layout/basic-layout.component.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { Component, ElementRef, OnInit, ViewChild } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgBreadcrumbOptions } from 'eo-ng-breadcrumb'
import { MenuOptions } from 'eo-ng-menu'
import { NzModalRef, NzModalService } from 'ng-zorro-antd/modal'
import { Subscription } from 'rxjs'
import { ApiService } from '../../service/api.service'
import { EoNgNavigationService } from '../../service/eo-ng-navigation.service'
import { BaseInfoService } from '../../service/base-info.service'
import { IframeHttpService } from '../../service/iframe-http.service'
import { environment } from '../../../environments/environment'
@Component({
  selector: 'basic-layout',
  templateUrl: './basic-layout.component.html',
  styleUrls: ['./basic-layout.component.scss']
})
export class BasicLayoutComponent implements OnInit {
  @ViewChild('breadcrumbTitleTpl', { static: true }) breadcrumbTitleTpl!: ElementRef

  sideMenuOptions: MenuOptions[] = []
  breadcrumbOptions: EoNgBreadcrumbOptions[] = []
  currentRouter: string = '' // 当前路由
  openMap: { [name: string]: boolean } = {}
  modalRef: NzModalRef | undefined
  showEmpty: boolean = false
  showSideLine: boolean = true
  authInfo: { title: string; infos: Array<{ key: string; value: string }> } = {
    title: '',
    infos: []
  }

  guideMenu: MenuOptions = {
    matchRouter: true,
    matchRouterExact: false,
    menu: true,
    name: 'guide',
    routerLink: 'guide',
    title: '🚀 快速入门',
    type: 'built-in',
    menuTitleClassName: 'menu-icon-hidden'
  }

  userAvatar:boolean = false // 是否显示用户头像，取决于是否开启用户权限插件
  isBusiness:boolean = environment.isBusiness

  private subscription1: Subscription = new Subscription()
  private subscription2: Subscription = new Subscription()
  private subscription3: Subscription = new Subscription()
  private subscription4: Subscription = new Subscription()

  constructor (
    private router: Router,
    private api: ApiService,
    private navigationService: EoNgNavigationService,
    private modalService: NzModalService,
    private baseInfo:BaseInfoService,
    private iframeService:IframeHttpService
  ) {
    this.subscription1 = this.navigationService
      .repFlashBreadcrumb()
      .subscribe((data: any) => {
        // data[0].iframe=true时面包屑是iframe的导航，需要特殊处理路由
        if (data && data[0] && data[0].iframe) {
          for (const bc of data) {
            bc.nzContext = { url: bc.routerLink, title: bc.title }
            bc.title = this.breadcrumbTitleTpl
          }
        }
        this.breadcrumbOptions = data
      })

    this.subscription2 = this.navigationService.repFlashMenu().subscribe(() => {
      this.sideMenuOptions = [
        ...(this.isBusiness ? [] : [this.guideMenu]),
        ...this.navigationService.getCurrentMenuList()
      ]
      for (const menu of this.sideMenuOptions) {
        menu.open = this.openMap[menu['titleString']! as string]
      }
      this.userAvatar = this.navigationService.getUserPlugin()
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

  clickIframeBreadcrumb (url:string) {
    const moduleName:string = this.baseInfo.allParamsInfo.moduleName
    window.location.href = `module/${moduleName}#/${url}`
    this.iframeService.reqFlashIframe(url)
  }

  getSideMenu () {
    this.subscription4 = this.navigationService
      .getMenuList()
      .subscribe((res: MenuOptions[]) => {
        this.sideMenuOptions = [this.guideMenu, ...res]
        this.checkOpenMenu()
        this.getAccess()
      })
  }

  checkOpenMenu () {
    for (const sideMenu of this.sideMenuOptions) {
      if (sideMenu?.children?.length) {
        if (sideMenu.children?.filter((menu:any) => {
          return this.router.url.includes(menu.routerLink) && this.router.url.substring(1, menu.routerLink.length + 1) === menu.routerLink
        }).length > 0) {
          this.openMap[sideMenu['titleString']] = true
        }
      }
    }
  }

  getAuthInfo () {
    if (this.navigationService.getUserAuthAccess()) {
      this.api.authGet('activation/info').subscribe(
        (resp: {
          code: number
          data: {
            infos: Array<{ key: string; value: string }>
            title: string
          }
          msg: string
        }) => {
          if (resp.code === 0) {
            this.authInfo = resp.data
          }
        }
      )
    }
  }

  openAuthDialog () {
    this.router.navigate(['/', 'auth-info'])
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
      if (
        this.router.routerState.snapshot.url === '/' ||
        this.router.routerState.snapshot.url === '/login'
      ) {
        this.router.navigate([this.navigationService.getPageRoute()])
      }
    }
    // this.getAuthInfo()
  }

  // 根据路由选中并打开对应menu
  selectOrOpenMenu (router: string): void {
    if (this.sideMenuOptions.length > 0) {
      for (const index in this.sideMenuOptions) {
        if (
          router.split('/')[1] ===
          this.sideMenuOptions[index]['router']?.split('/')[0]
        ) {
          if (this.sideMenuOptions[index].children?.length) {
            // this.openHandler(this.sideMenuOptions[index]['key']!)
          }
          break
        }
      }
      this.sideMenuOptions = [...this.sideMenuOptions]
      this.currentRouter = router
    }
  }

  // openHandler (key: string): void {
  //   for (const index in this.sideMenuOptions) {
  //     if (this.sideMenuOptions[index]['key'] !== key) {
  //       // this.sideMenuOptions[index].open = false
  //     } else {
  //       this.sideMenuOptions[index].open = true
  //     }
  //     this.openMap[this.sideMenuOptions[index]['key'] as string] = !!this.sideMenuOptions[index].open
  //   }
  // }

  goToGithub () {
    window.open('https://github.com/eolinker/apinto')
  }

  goToHelp () {
    window.open('https://help.apinto.com/docs')
  }
}

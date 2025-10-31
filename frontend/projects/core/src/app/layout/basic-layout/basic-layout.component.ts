/* eslint-disable dot-notation */
/*
 * @Author: maggieyyy im.ymj@hotmail.com
 * @Date: 2022-07-12 00:19:11
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-05-30 13:07:41
 * @FilePath: \apinto\projects\core\src\app\layout\basic-layout\basic-layout.component.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { Component, ElementRef, OnInit, ViewChild, ViewContainerRef } from '@angular/core'
import { ActivatedRoute, NavigationEnd, Router } from '@angular/router'
import { EoNgBreadcrumbOptions } from 'eo-ng-breadcrumb'
import { MenuOptions } from 'eo-ng-menu'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { Subscription, filter, map, mergeMap } from 'rxjs'
import { EoNgNavigationService } from '../../service/eo-ng-navigation.service'
import { BaseInfoService } from '../../service/base-info.service'
import { IframeHttpService } from '../../service/iframe-http.service'
import { environment } from '../../../environments/environment'
import { TryBusinessAddr } from '../../constant/conf'
import { PluginSlotHubService } from '../../service/plugin-slot-hub.service'
import { PluginEventHubService } from '../../service/plugin-event-hub.service'
@Component({
  selector: 'basic-layout',
  templateUrl: './basic-layout.component.html',
  styleUrls: ['./basic-layout.component.scss']
})
export class BasicLayoutComponent implements OnInit {
  @ViewChild('breadcrumbTitleTpl', { static: true }) breadcrumbTitleTpl!: ElementRef
  @ViewChild('topContainer', { read: ViewContainerRef, static: true }) topContainer!: any;
  @ViewChild('authButton', { read: ViewContainerRef, static: true }) authButton!: any;

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

  newAvatar:any = undefined
  newAuthButton:any = undefined

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
  showGuide = this.baseInfo.showGuide

  private subscriptions = new Subscription();
  private subAuthCheck: Subscription = new Subscription()

  constructor (
    private router: Router,
    private route: ActivatedRoute,
    private navigationService: EoNgNavigationService,
    private baseInfo:BaseInfoService,
    private iframeService:IframeHttpService,
    private pluginSlotHub:PluginSlotHubService,
    private pluginEventHub:PluginEventHubService
  ) {
    this.subscriptions.add(this.navigationService
      .repFlashBreadcrumb()
      .subscribe((data: any) => {
        // data[0].iframe=true时面包屑是iframe的导航，需要特殊处理路由
        if (data && data[0] && data[0].iframe) {
          for (const bc of data) {
            bc.nzContext = { url: bc.routerLink, title: bc.title }
            bc.title = this.breadcrumbTitleTpl
          }
        }
        while (this.breadcrumbOptions.length > 0) {
          this.breadcrumbOptions.pop()
        }
        for (const newBd of data) {
          this.breadcrumbOptions.push(newBd)
        }
      }))

    this.subscriptions.add(this.navigationService.repFlashMenu().subscribe(() => {
      this.sideMenuOptions = (this.baseInfo.showGuide
        ? [{ ...this.guideMenu },
            ...this.navigationService.getCurrentMenuList()
          ]
        : [...this.navigationService.getCurrentMenuList()]).map((x) => { x.click = (e:any) => { console.log(e) }; return x })
      for (const menu of this.sideMenuOptions) {
        menu.open = this.openMap[menu['titleString']! as string]
      }
      if (this.sideMenuOptions.length === 0) { this.showEmpty = true }
      this.userAvatar = this.navigationService.getUserPlugin()
    }))

    this.subscriptions.add(this.router.events.subscribe(() => {
      if (this.router.url !== this.currentRouter) {
        this.selectOrOpenMenu(this.router.url)
      }
    }))
  }

  ngOnInit () {
    this.getSideMenu()
    this.pluginEventHub.initHub()!.emit('projectInited', {})
  //   this.router.events.pipe(
  //     filter(event => event instanceof NavigationEnd),
  //     map(() => this.route),
  //     map(route => {
  //       while (route.firstChild) route = route.firstChild
  //       return route
  //     }),
  //     filter(route => route.outlet === 'primary'),
  //     mergeMap(route => route.data)
  //   ).subscribe((event) => {
  //     this.updateMenuOptions()
  //   })
  }

  // updateMenuOptions () {
  //   this.sideMenuOptions = this.sideMenuOptions.map(option => {
  //     if (option.children && option.children.length > 0) {
  //       return false
  //     }
  //     if (this.router.url.startsWith('/' + option.routerLink)) {
  //       return { ...option, disabled: true }
  //     } else {
  //       return { ...option, disabled: false }
  //     }
  //   })
  //   console.log(this.router.url, this.sideMenuOptions)
  // }

  ngAfterViewInit () {
    this.newAvatar = this.pluginSlotHub.getSlot('renderAvatar')
    this.topContainer.clear()
    if (this.newAvatar) {
      this.topContainer.createComponent(...this.pluginSlotHub.getSlot('renderAvatar'))
    }
    this.newAuthButton = this.pluginSlotHub.getSlot('renderAuthButton')
    this.authButton.clear()
    if (this.newAuthButton) {
      this.authButton.createComponent(...this.pluginSlotHub.getSlot('renderAuthButton'))
    }
  }

  ngOnDestroy () {
    this.subscriptions.unsubscribe()
    this.subAuthCheck.unsubscribe()
  }

  clickIframeBreadcrumb (url:string) {
    const moduleName:string = this.baseInfo.allParamsInfo.moduleName
    window.location.href = `module/${moduleName}#/${url}`
    this.iframeService.reqFlashIframe(url)
  }

  getSideMenu () {
    const subscription = this.navigationService.getMenuList().subscribe((resp:MenuOptions[]) => {
      if (resp !== undefined) {
        const newMenu = resp
        this.sideMenuOptions = [...(this.baseInfo.showGuide ? [this.guideMenu] : []), ...newMenu]
        this.checkOpenMenu()
      }
      subscription.unsubscribe()
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

  goToGithub () {
    window.open('https://github.com/eolinker/apinto')
  }

  goToHelp () {
    window.open('https://help.apinto.com/docs')
  }

  goToRry () {
    window.open(TryBusinessAddr)
  }
}

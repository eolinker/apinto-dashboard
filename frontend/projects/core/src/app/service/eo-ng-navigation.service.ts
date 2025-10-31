import { Injectable } from '@angular/core'
import { MenuOptions } from 'eo-ng-menu'
import { Subject, Observable, concatMap, from } from 'rxjs'
import { ApiService } from './api.service'
import { v4 as uuidv4 } from 'uuid'
import { environment } from '../../environments/environment'
import { PluginSlotHubService } from './plugin-slot-hub.service'
import { BaseInfoService } from './base-info.service'

@Injectable({
  providedIn: 'root'
})
export class EoNgNavigationService {
  private menuList: MenuOptions[] = [] // 当前用户可显示的菜单
  private updateRightsRouterList: string[] = [] // 当前用户可编辑的菜单rouer列表
  private viewRightsRouterList: string[] = [] // 当前用户可查看的菜单router列表
  private mainPageRouter: string = '' // 首页路由
  dataUpdated: boolean = false // 是否获取过数据，避免组件在ngOnChanges时读取空数组
  userInfo:{[k:string]:any} = {} // 用户数据，由用户管理插件传入
  private userRoleId: string = '' // 当前用户角色id
  private userId: string = '' // 当前用户id
  private navLayoutHidden:boolean = false
  constructor (public api: ApiService, private pluginSlotHub:PluginSlotHubService, private baseInfo:BaseInfoService) {
  }

  iframePrefix:string = 'module' // 与后端约定好的，所有iframe打开的页面都要加该前缀
  isBusiness:boolean = environment.isBusiness
  setUserRoleId (val: string) {
    this.userRoleId = val
  }

  getUserRoleId () {
    return this.userRoleId
  }

  setUserId (id: string) {
    this.userId = id
  }

  getUserId () {
    return this.userId
  }

  // TODO返回控制台是否开启用户角色插件
  getUserPlugin () {
    return false
  }

  // 获取首页路由地址
  getPageRoute (): string {
    return this.baseInfo.showGuide ? '/guide' : '/router/api'
  }

  // 如果用户没有任何除商业授权以外的功能查看权限, 返回true
  getUserAccess () {
    return this.noAccess
  }

  // 用户是否有查看授权的权限，有则返回true
  getUserAuthAccess () {
    return this.viewRightsRouterList.includes('auth-info')
  }

  // 用户是否有某个模块的权限，options目前为空，未来可以用于存储项目/空间id
  // 权限通过插槽实现，如果项目没有安装用户管理插件，则不存在对应插槽，则默认用户拥有所有模块的编辑权限
  getUserModuleAccess (moduleName:string, options?:any) {
    const checkModuleAccess = this.pluginSlotHub.getSlot('checkModuleAccess')
    return checkModuleAccess ? checkModuleAccess(moduleName, options) : 'edit'
  }

  setNavHidden (hidden:boolean) {
    this.navLayoutHidden = hidden
  }

  getNavHidden () {
    return this.navLayoutHidden
  }

  private flashFlag: Subject<boolean> = new Subject<boolean>()

  reqFlashMenu () {
    this.flashFlag.next(true)
  }

  repFlashMenu () {
    return this.flashFlag.asObservable()
  }

  private userUpdateRightList: Subject<boolean> = new Subject<
    boolean
  >()

  reqUpdateRightList () {
    this.userUpdateRightList.next(true)
    this.dataUpdated = true
  }

  repUpdateRightList () {
    return this.userUpdateRightList.asObservable()
  }

  getUpdateRightsRouter () {
    return this.updateRightsRouterList
  }

  getViewRightsRouter () {
    return this.viewRightsRouterList
  }

  // 获取最新目录列表
  getMenuList (): Observable<MenuOptions[]> {
    return this.getRightsList()
  }

  // 获取当前目录列表
  getCurrentMenuList (): MenuOptions[] {
    return this.menuList
  }

  menuMap: Map<number, any> = new Map()
  routerNameMap: Map<string, string> = new Map()
  modulesMap: Map<string, any> = new Map()
  accessList: Array<string> = []
  firstModulesId: number | null = null
  findFirstModulesId: boolean = false
  noAccess: boolean = true
  originAccessData:{[k:string]:string} = {} // 从接口获取的access数据
  count:number = 0
  // 获得最新的权限列表和菜单
  getRightsList (): Observable<MenuOptions[]> {
    return new Observable((observer) => {
      this.noAccess = true
      const renewAccess = this.pluginSlotHub.getSlot('renewModuleAccess')
      const renewAccess$ = renewAccess ? from(renewAccess()) : from([null])

      renewAccess$.pipe(
        concatMap(() => {
          return this.api.get('system/modules')
        })
      ).subscribe((resp: any) => {
        if (resp.code === 0) {
          this.generateMenuList(resp)
        }
        observer.next(this.menuList)
        this.reqFlashMenu()
        this.reqUpdateRightList()
      })
    })
  }

  generateMenuList (resp:any) {
    if (resp.code === 0) {
      this.mainPageRouter = ''
      this.modulesMap = new Map()
      this.menuList = []
      this.routerNameMap = new Map()
      this.originAccessData = resp.data.access
      for (const navigation of resp.data.navigation) {
        const menu:any = {
          title: navigation.title,
          titleString: navigation.title,
          menu: true,
          key: uuidv4(),
          icon: navigation.icon || 'daohang',
          ...(navigation.modules?.length > 0 && !navigation.default
            ? {
                children: navigation.modules.filter((module:any) => {
                  // TODO 插件化一期不做权限
                  return this.getUserModuleAccess(module.name)
                }).map((module: any) => {
                  this.routerNameMap.set(module.path, module.name)
                  return {
                    title: module.title,
                    titleString: navigation.title,
                    name: module.name,
                    type: module.type,
                    routerLink: module.path,
                    matchRouter: true,
                    matchRouterExact: false
                  }
                })
              }
            : navigation.modules?.length > 0 &&
              this.getDefaultModule(navigation).path
              ? {
                  name: this.getDefaultModule(navigation).name,
                  routerLink: this.getDefaultModule(navigation).path,
                  matchRouter: true,
                  matchRouterExact: false,
                  type: this.getDefaultModule(navigation).type
                }
              : {
                  children: [{
                    menu: true,
                    group: true,
                    title: '暂无内容',
                    menuTitleClassName: 'menu-no-content'
                  }]
                })
        }
        if (navigation.name && navigation.path) {
          this.routerNameMap.set(navigation.path, navigation.name)
        }
        if ((menu as any).routerLink && !this.routerNameMap.get((menu as any).routerLink)) {
          this.routerNameMap.set((menu as any).routerLink, (menu as any).name)
        }

        // TODO 插件化一期不做权限
        if (this.getUserModuleAccess(menu.name) || (!menu.name && menu.children?.length > 0)) {
          if (menu.title === '系统管理' && environment.isBusiness) {
            menu.children.unshift({
              name: 'auth',
              routerLink: 'auth-info',
              title: '授权管理',
              titleString: '授权管理',
              matchRouter: true,
              matchRouterExact: false
            })
          }
          this.menuList.push(menu)
        }
      }

      if (!this.mainPageRouter) {
        this.findMainPage()
      }
    }
  }

  getDefaultModule (nav: any): { name: string; path: string; type: string } {
    let res = { name: '', path: '', type: '' }
    if (!nav.default) {
      return res
    }
    for (const module of nav.modules) {
      if (module.name === nav.default) {
        res = { ...res, ...module }
        return res
      }
    }
    return res
  }

  findMainPage () {
    for (const menu of this.menuList) {
      // eslint-disable-next-line dot-notation
      // TODO 插件化一期不做权限
      // if (menu.routerLink && menu['name'] && this.getUserModuleAccess(menu['name']) && !menu.routerLink.includes('module/')) {
      if (menu.routerLink && menu['name'] && !menu.routerLink.includes('module/') && this.getUserModuleAccess(menu['name'])) {
        this.mainPageRouter = menu.routerLink
        return
      } else if (menu.children) {
        for (const child of menu.children) {
          if (
            // eslint-disable-next-line dot-notation
          // TODO 插件化一期不做权限
            // child.routerLink && this.getUserModuleAccess(child['name']) && !child.routerLink.includes('module/')
            child.routerLink && !child.routerLink.includes('module/') && this.getUserModuleAccess(child['name'])
          ) {
            this.mainPageRouter = child.routerLink
            return
          }
        }
      }
    }
  }

  // 检查用户是否有编辑该路由页面下内容的权限,若有返回true
  // checkUpdateRight (menuRouter: string) {
  //   return new Observable((observer) => {
  //     if (this.updateRightsRouterList?.length === 0) {
  //       this.getMenuList().subscribe(() => {
  //         observer.next(
  //           this.getUserModuleAccess(menuRouter)?.filter((x) => {
  //             return x.includes('edit')
  //           }).length
  //         )
  //       })
  //     } else {
  //       observer.next(
  //         this.getUserModuleAccess(menuRouter)?.filter((x) => {
  //           return x.includes('edit')
  //         }).length
  //       )
  //       // return of(this.updateRightsRouterList.indexOf(menuRouter) !== -1)
  //     }
  //   })
  // }

  private breadcrumb: Subject<any> = new Subject<any>()
  private breadcrumbList: MenuOptions[] = []
  getLatestBreadcrumb () {
    return this.breadcrumbList
  }

  reqFlashBreadcrumb (value: any, type?:string) {
    if (type && type === 'iframe') {
      value[0].iframe = true
    }
    this.breadcrumbList = value
    this.breadcrumb.next(value)
  }

  repFlashBreadcrumb () {
    return this.breadcrumb.asObservable()
  }
}

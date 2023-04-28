import { Injectable } from '@angular/core'
import { MenuOptions } from 'eo-ng-menu'
import { Subject, Observable, forkJoin, map } from 'rxjs'
import { ApiService } from './api.service'
import { v4 as uuidv4 } from 'uuid'

@Injectable({
  providedIn: 'root'
})
export class EoNgNavigationService {
  private menuList: MenuOptions[] = [] // 当前用户可显示的菜单
  private updateRightsRouterList: string[] = [] // 当前用户可编辑的菜单rouer列表
  private viewRightsRouterList: string[] = [] // 当前用户可查看的菜单router列表
  private mainPageRouter: string = '' // 首页路由
  dataUpdated: boolean = false // 是否获取过数据，避免组件在ngOnChanges时读取空数组
  private userRoleId: string = '' // 当前用户角色id
  private userId: string = '' // 当前用户id
  accessMap: Map<string, string> = new Map()
  constructor (public api: ApiService) {}
  iframePrefix:string = 'module' // 与后端约定好的，所有iframe打开的页面都要加该前缀
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
    return '/guide'
    // return this.mainPageRouter
  }

  // 如果用户没有任何除商业授权以外的功能查看权限, 返回true
  getUserAccess () {
    return this.noAccess
  }

  // 用户是否有查看授权的权限，有则返回true
  getUserAuthAccess () {
    return this.viewRightsRouterList.includes('auth-info')
  }

  private flashFlag: Subject<boolean> = new Subject<boolean>()

  reqFlashMenu () {
    this.flashFlag.next(true)
  }

  repFlashMenu () {
    return this.flashFlag.asObservable()
  }

  private userUpdeteRightList: Subject<boolean> = new Subject<
    boolean
  >()

  reqUpdateRightList () {
    this.userUpdeteRightList.next(true)
    this.dataUpdated = true
  }

  repUpdateRightList () {
    return this.userUpdeteRightList.asObservable()
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
      forkJoin([this.api.get('system/modules').pipe(map((resp:any) => {
        this.generateMenuList(resp)
      })),
      this.api.get('my/access')]).subscribe((resp:Array<any>) => {
        this.accessMap = new Map()
        const accessListBackend:Array<{name:string, access:string}> = resp[1].data.access
        for (const access of accessListBackend) {
          this.accessMap.set(access.name, access.access)
          if (access.access && this.noAccess) {
            this.noAccess = false
          }
        }
        observer.next(this.menuList)

        this.reqFlashMenu()
        this.reqUpdateRightList()
      })
    }
    )
  }

  generateMenuList (resp:any) {
    if (resp.code === 0) {
      this.mainPageRouter = ''
      this.modulesMap = new Map()
      this.menuList = []
      this.routerNameMap = new Map()
      this.noAccess = true
      this.originAccessData = resp.data.access
      for (const navigation of resp.data.navigation) {
        const menu = {
          title: navigation.title,
          titleString: navigation.title,
          menu: true,
          key: uuidv4(),
          icon: navigation.icon || 'daohang',
          ...(navigation.modules?.length > 0 && !navigation.default
            ? {
                children: navigation.modules.map((module: any) => {
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
              : (navigation.modules?.length > 0
                  ? { // TODO似乎没用，后续排查
                    }
                  : {
                      children: [{
                        menu: true,
                        group: true,
                        title: '暂无内容',
                        menuTitleClassName: 'menu-no-content'
                      }]
                    }))
        }
        if (navigation.name && navigation.path) {
          this.routerNameMap.set(navigation.path, navigation.name)
        }
        if ((menu as any).routerLink && !this.routerNameMap.get((menu as any).routerLink)) {
          this.routerNameMap.set((menu as any).routerLink, (menu as any).name)
        }
        this.menuList.push(menu)
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
      if (menu.routerLink && this.accessMap?.get(menu?.routerLink)?.length) {
        this.mainPageRouter = menu.routerLink
        return
      } else if (menu.children) {
        for (const child of menu.children) {
          if (
            child.routerLink &&
            this.accessMap?.get(child?.routerLink)?.length
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
  //           this.accessMap.get(menuRouter)?.filter((x) => {
  //             return x.includes('edit')
  //           }).length
  //         )
  //       })
  //     } else {
  //       observer.next(
  //         this.accessMap.get(menuRouter)?.filter((x) => {
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

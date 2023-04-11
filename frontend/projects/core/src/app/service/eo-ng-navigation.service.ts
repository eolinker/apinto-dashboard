import { Injectable } from '@angular/core'
import { MenuOptions } from 'eo-ng-menu'
import { Subject, Observable } from 'rxjs'
import { ApiService } from './api.service'

@Injectable({
  providedIn: 'root'
})
export class EoNgNavigationService {
  private menuList:MenuOptions[] = [] // 当前用户可显示的菜单
  private updateRightsRouterList:string[] = [] // 当前用户可编辑的菜单rouer列表
  private viewRightsRouterList:string[] = [] // 当前用户可查看的菜单router列表
  private mainPageRouter:string = '' // 首页路由
  dataUpdated:boolean = false // 是否获取过数据，避免组件在ngOnChanges时读取空数组
  private userRoleId:string = '' // 当前用户角色id
  private userId:string = '' // 当前用户id
  private accessMap:Map<string, Array<string>> = new Map()
  constructor (
    public api:ApiService) {
  }

  setUserRoleId (val:string) {
    this.userRoleId = val
  }

  getUserRoleId () {
    return this.userRoleId
  }

  setUserId (id:string) {
    this.userId = id
  }

  getUserId () {
    return this.userId
  }

  // 获取首页路由地址
  getPageRoute ():string {
    return this.mainPageRouter
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

  private userUpdeteRightList: Subject<Array<string>> = new Subject<Array<string>>()

  reqUpdateRightList () {
    this.userUpdeteRightList.next(this.updateRightsRouterList)
    this.dataUpdated = true
  }

  repUpdateRightList () {
    return this.userUpdeteRightList.asObservable()
  }

  private userViewUpdeteRightList: Subject<Array<string>> = new Subject<Array<string>>()

  reqViewRightList () {
    this.userViewUpdeteRightList.next(this.viewRightsRouterList)
    this.dataUpdated = true
  }

  repViewRightList () {
    return this.userViewUpdeteRightList.asObservable()
  }

  getUpdateRightsRouter () {
    return this.updateRightsRouterList
  }

  getViewRightsRouter () {
    return this.viewRightsRouterList
  }

  // 获取最新目录列表
  getMenuList ():Observable<MenuOptions[]> {
    return this.getRightsList()
  }

  // 获取当前目录列表
  getCurrentMenuList ():MenuOptions[] {
    return this.menuList
  }

  menuMap:Map<number, any> = new Map()
  routerNameMap:Map<string, string> = new Map()
  modulesMap:Map<string, any> = new Map()
  accessList:Array<string> = []
  firstModulesId:number|null = null
  findFirstModulesId:boolean = false
  noAccess:boolean = true
  // 获得最新的权限列表和菜单
  getRightsList () :Observable<MenuOptions[]> {
    return new Observable(observer => {
      this.api.get('system/modules').subscribe((resp:any) => {
        if (resp.code === 0) {
          this.mainPageRouter = ''
          this.modulesMap = new Map()
          this.menuList = []
          this.accessMap = new Map()
          this.routerNameMap = new Map()
          this.noAccess = true
          for (const navigation of resp.data.navigation) {
            const menu = {
              title: navigation.title,
              icon: navigation.icon,
              iconType: navigation.iconType,
              menu: true,
              name: navigation.name,
              ...(navigation.modules?.length > 0
                ? {
                    children: navigation.modules.map((module:any) => {
                      this.routerNameMap.set(module.name, module.path)
                      return {
                        title: module.title,
                        name: module.name,
                        type: module.type,
                        ...(module.type === 'buildin'
                          ? {
                              routerLink: module.path,
                              matchRouter: true,
                              matchRouterExact: false
                            }
                          : {
                              path: module.path
                            })
                      }
                    })
                  }
                : {
                    routerLink: navigation.path,
                    matchRouter: true,
                    matchRouterExact: false
                  })
            }

            if (navigation.name && navigation.path) {
              this.routerNameMap.set(navigation.name, navigation.path)
            }
            this.menuList.push(menu)
          }

          for (const acc of Object.keys(resp.data.access)) {
            // accessMap 存的是router-access
            this.accessMap.set(this.routerNameMap.get(acc) || acc, resp.data.access[acc])
            if (resp.data.access[acc]?.length > 0 && this.noAccess) {
              this.noAccess = false
            }
          }

          if (!this.mainPageRouter) {
            this.findMainPage()
          }

          observer.next(this.menuList)
          this.reqFlashMenu()
          this.reqUpdateRightList()
        }
      })
    })
  }

  findMainPage () {
    for (const menu of this.menuList) {
      if (menu.routerLink && this.accessMap?.get(menu?.routerLink)?.length) {
        this.mainPageRouter = menu.routerLink
        return
      } else if (menu.children) {
        for (const child of menu.children) {
          if (child.routerLink && this.accessMap?.get(child?.routerLink)?.length) {
            this.mainPageRouter = child.routerLink
            return
          }
        }
      }
    }
  }

  // 检查用户是否有编辑该路由页面下内容的权限,若有返回true
  checkUpdateRight (menuRouter:string) {
    return new Observable(observer => {
      if (this.updateRightsRouterList?.length === 0) {
        this.getMenuList().subscribe(() => {
          observer.next(this.accessMap.get(menuRouter)?.filter(x => {
            return x.includes('edit')
          }).length)
        })
      } else {
        observer.next(this.accessMap.get(menuRouter)?.filter(x => {
          return x.includes('edit')
        }).length)
      // return of(this.updateRightsRouterList.indexOf(menuRouter) !== -1)
      }
    })
  }

  private breadcrumb: Subject<any> = new Subject<any>()
  private breadcrumbList:MenuOptions[] = []
  getLatestBreadcrumb () {
    return this.breadcrumbList
  }

  reqFlashBreadcrumb (value:any) {
    this.breadcrumbList = value
    this.breadcrumb.next(value)
  }

  repFlashBreadcrumb () {
    return this.breadcrumb.asObservable()
  }
}

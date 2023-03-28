/* eslint-disable dot-notation */
import { Injectable } from '@angular/core'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { MenuOptions } from 'eo-ng-menu/public-api'
import { Observable, Subject } from 'rxjs'
import { AppConfig } from '../constant/app.config'
import { ApiService } from './api.service'

@Injectable({
  providedIn: 'root'
})
export class AppConfigService {
  private menuList:MenuOptions[] = [] // 当前用户可显示的菜单
  private updateRightsRouterList:string[] = [] // 当前用户可编辑的菜单rouer列表
  private viewRightsRouterList:string[] = [] // 当前用户可查看的菜单router列表
  private mainPageRouter:string = '' // 首页路由
  dataUpdated:boolean = false // 是否获取过数据，避免组件在ngOnChanges时读取空数组
  private userRoleId:string = '' // 当前用户角色id
  private userId:string = '' // 当前用户id
  // eslint-disable-next-line no-useless-constructor
  constructor (
    private message: EoNgFeedbackMessageService,
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
    return this.viewRightsRouterList.length === 0 || (this.viewRightsRouterList.length === 1 && this.viewRightsRouterList[0] === 'auth-info')
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
  modulesMap:Map<string, any> = new Map()
  accessList:Array<string> = []
  firstModulesId:number|null = null
  findFirstModulesId:boolean = false
  // 获得最新的权限列表和菜单
  getRightsList () :Observable<MenuOptions[]> {
    return new Observable(observer => {
      this.api.get('my/modules').subscribe((resp:any) => {
        if (resp.code === 0) {
          this.updateRightsRouterList = []
          this.viewRightsRouterList = []
          this.mainPageRouter = ''
          this.firstModulesId = null
          this.modulesMap = new Map()
          this.firstModulesId = resp.data.modules[0]?.id
          for (const index in resp.data.modules) {
            if (resp.data.modules[index].access.length > 0) {
              this.modulesMap.set(resp.data.modules[index].id, resp.data.modules[index].access)
              if (resp.data.modules[index].access.join().includes('view') && !this.findFirstModulesId) {
                this.firstModulesId = resp.data.modules[index].id
                this.findFirstModulesId = true
              }
            }
            this.accessList.push(...resp.data.modules[index].access)
          }
          this.menuList = []
          const originMenuList = AppConfig.menuList
          for (const indexMenu in originMenuList) {
            const moduleAccess = this.modulesMap.get(originMenuList[indexMenu].id)
            // 有子菜单时,需要匹配子菜单的view和edit字段, 需要确认该一级菜单的路由指向
            if (originMenuList[indexMenu].children) {
              const childTemp = this.childMenu(originMenuList[indexMenu].children)
              if (childTemp.length > 0) {
                originMenuList[indexMenu].routerLink = childTemp[0].routerLink
                const tempMenu : MenuOptions = {
                  title: originMenuList[indexMenu].title,
                  id: originMenuList[indexMenu].id,
                  router: originMenuList[indexMenu].router,
                  menu: originMenuList[indexMenu].menu,
                  children: childTemp,
                  icon: originMenuList[indexMenu].icon
                }
                this.menuMap.set(tempMenu['id'], tempMenu)
                this.viewRightsRouterList.push(originMenuList[indexMenu].routerLink)
              }
            } else if (moduleAccess?.length) {
              // 无子菜单时, 匹配一级菜单的view和edit字段,有view时才可显示
              if (moduleAccess.indexOf(originMenuList[indexMenu].view) !== -1) {
                const tempMenu : MenuOptions = {
                  title: originMenuList[indexMenu].title,
                  routerLink: originMenuList[indexMenu].routerLink,
                  matchRouter: originMenuList[indexMenu].matchRouter,
                  matchRouterExact: originMenuList[indexMenu].matchRouterExact,
                  id: originMenuList[indexMenu].id,
                  menu: originMenuList[indexMenu].menu,
                  icon: originMenuList[indexMenu].icon
                }
                this.menuMap.set(tempMenu['id'], tempMenu)
                this.viewRightsRouterList.push(originMenuList[indexMenu].routerLink)
                if (this.firstModulesId === originMenuList[indexMenu].id) {
                  this.mainPageRouter = originMenuList[indexMenu].routerLink
                }
              }
              // 判断是否有编辑权限,有则将路由放入updateRightsRouterList
              if (moduleAccess.indexOf(originMenuList[indexMenu].edit) !== -1) {
                this.updateRightsRouterList.push(originMenuList[indexMenu].routerLink)
              }
            }
          }

          for (const index in resp.data.modules) {
            if (this.menuMap.has(resp.data.modules[index].id) && this.menuMap.get(resp.data.modules[index].id).menu) {
              this.menuList.push(this.menuMap.get(resp.data.modules[index].id))
            }
          }
          observer.next(this.menuList)
          this.reqFlashMenu()
          this.reqUpdateRightList()
        }
      })
    })
  }

  // 子菜单,根据获得的权限列表,选择是否显示菜单
  childMenu (data:any):MenuOptions[] {
    const childMenuList:MenuOptions[] = []
    for (const index in data) {
      const accessList = this.modulesMap.get(data[index].id)
      if (accessList?.length > 0 && accessList.indexOf(data[index].view) !== -1) {
        const tempMenu : MenuOptions = {
          title: data[index].title,
          routerLink: data[index].routerLink,
          matchRouter: data[index].matchRouter,
          matchRouterExact: data[index].matchRouterExact,
          id: data[index].id
        }
        childMenuList.push(tempMenu)
        this.viewRightsRouterList.push(data[index].routerLink)
        if (this.firstModulesId === data[index].id) {
          this.mainPageRouter = data[index].routerLink
        }
      }
      if (accessList?.length > 0 && accessList.indexOf(data[index].edit) !== -1) {
        this.updateRightsRouterList.push(data[index].routerLink)
      }
    }
    return childMenuList
  }

  openMap:{[name:string]:boolean}={}

  openHandler = (title:string) => {
    for (const key in this.openMap) {
      if (key !== title) {
        this.openMap[key] = false
      } else {
        this.openMap[key] = true
      }
    }
  }

  // 检查用户是否有编辑该路由页面下内容的权限,若有返回true
  checkUpdateRight (menuRouter:string) {
    return new Observable(observer => {
      if (this.updateRightsRouterList?.length === 0) {
        this.getMenuList().subscribe(() => {
          observer.next(this.updateRightsRouterList.indexOf(menuRouter) !== -1)
        })
      } else {
        observer.next(this.updateRightsRouterList.indexOf(menuRouter) !== -1)
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

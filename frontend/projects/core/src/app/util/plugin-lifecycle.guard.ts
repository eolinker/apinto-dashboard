import { Injectable } from '@angular/core'
import { ActivatedRouteSnapshot, CanActivate, CanActivateChild, CanDeactivate, CanLoad, NavigationEnd, Router, RouterStateSnapshot, UrlTree } from '@angular/router'
import { Observable, Subscription } from 'rxjs'
import { PluginLoaderService } from '../service/plugin-loader.service'

@Injectable({
  providedIn: 'root'
})
export class PluginLifecycleGuard implements CanActivate, CanActivateChild, CanDeactivate<unknown>, CanLoad {
  private navigationEndSubscription: Subscription | null = null;
  private isRedirected: boolean = false; // 添加一个标志来防止重定向循环

  constructor (private pluginLoader: PluginLoaderService, private router: Router) {}

  private async handleBeforeMount (route: ActivatedRouteSnapshot, state:RouterStateSnapshot): Promise<string | null> {
    if (this.isRedirected) {
      this.isRedirected = false // 如果已经重定向过了，重置标志并允许路由继续
      return null
    }
    const modulePath = route.routeConfig!.path as string
    const module = this.pluginLoader.getModule(route.routeConfig!.path as string)
    if (!this.isRedirected && !this.isPathFullyRecognized(modulePath, state.url)) {
      this.isRedirected = true
      this.router.navigate([modulePath])
      this.pluginLoader.redirectUrl = state.url
      // return state.url
    }
    if (module && module.beforeMount) {
      await module.beforeMount({ redirectUrl: encodeURI(state.url) })
    }
    return null
  }

  private handleMounted (): void {
    this.navigationEndSubscription = this.router.events.subscribe(event => {
      if (event instanceof NavigationEnd) {
        const module = this.pluginLoader.getModule(event.urlAfterRedirects.split('?')[0])
        if (module && module.mount) {
          module.mount()
        }
        this.navigationEndSubscription?.unsubscribe()
      }
    })
  }

  private async handleBeforeUnMount (route: ActivatedRouteSnapshot): Promise<void> {
    const module = this.pluginLoader.getModule(route.routeConfig!.path as string)
    this.pluginLoader.redirectUrl = ''
    if (module && module.beforeUnmount) {
      await module.beforeUnmount()
    }
  }

  private handleAfterUnMount (previousRoute: ActivatedRouteSnapshot): void {
    this.navigationEndSubscription?.unsubscribe() // Unsubscribe from previous subscription
    this.navigationEndSubscription = this.router.events.subscribe(event => {
      if (event instanceof NavigationEnd && !event.urlAfterRedirects.startsWith(previousRoute.routeConfig!.path as string)) {
        const module = this.pluginLoader.getModule(previousRoute.routeConfig!.path as string)
        if (module && module.unmount) {
          module.unmount()
        }
        this.navigationEndSubscription?.unsubscribe() // Unsubscribe to avoid memory leaks
      }
    })
  }

  async canActivate (route: ActivatedRouteSnapshot, state:RouterStateSnapshot): Promise<boolean | UrlTree> {
    try {
      await this.handleBeforeMount(route, state)
      this.handleMounted()
      return true
    } catch (error) {
      console.error('Error in lifecycle guard:', error)
      return false
    }
  }

  canActivateChild (): Observable<boolean | UrlTree> | Promise<boolean | UrlTree> | boolean | UrlTree {
    return true
  }

  async canDeactivate (component: unknown, currentRoute: ActivatedRouteSnapshot, currentState: RouterStateSnapshot, nextState?: RouterStateSnapshot): Promise<boolean | UrlTree> {
    await this.handleBeforeUnMount(currentRoute)
    if (nextState) {
      this.handleAfterUnMount(currentRoute)
    }
    return true
  }

  canLoad (): Observable<boolean | UrlTree> | Promise<boolean | UrlTree> | boolean | UrlTree {
    return true
  }

  private isPathFullyRecognized (basePath:string, currentPath:string): boolean {
    // 此函数判断当前的路由是否是远程模块的子路由
    return currentPath.split('?')[0] === `/${basePath}`
  }
}

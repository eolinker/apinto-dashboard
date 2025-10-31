/*
 * @Author: maggieyyy im.ymj@hotmail.com
 * @Date: 2022-07-11 23:20:14
 * @LastEditors: maggieyyy
 * @LastEditTime: 2023-12-14 20:56:39
 * @FilePath: \apinto\projects\core\src\app\app.component.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { Component } from '@angular/core'
import { CookieService } from 'ngx-cookie-service'
import { BaseInfoService } from './service/base-info.service'
import { AsIframeService } from './service/as-iframe.service'

@Component({
  selector: 'app-root',
  template: `
  <router-outlet></router-outlet>
  `,
  styles: []
})
export class AppComponent {
  titleDom:HTMLElement|null = document.querySelector('#appTitle')
  iconDom:HTMLLinkElement |null= document.querySelector('#appIcon')

  constructor (private cookieService: CookieService, public baseInfo:BaseInfoService, private asIframeService:AsIframeService) {
  }

  ngOnInit () {
    const time: number = 200 * 60 * 60 * 1000// cookie过期时间200个小时 200*60*60*1000
    this.cookieService.set('namespace', 'default', new Date(new Date().getTime() + time))
  }

  ngAfterViewInit () {
    if (this.titleDom) {
      this.titleDom.innerHTML = this.baseInfo.product
    }
  }

  ngOnDestroy () {
    this.asIframeService.removeReceiveMessage()
  }
}

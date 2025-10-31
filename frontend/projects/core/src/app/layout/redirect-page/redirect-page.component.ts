/* eslint-disable dot-notation */
/*
 * @Author: maggieyyy im.ymj@hotmail.com
 * @Date: 2022-07-12 00:19:11
 * @LastEditors: maggieyyy
 * @LastEditTime: 2023-12-22 15:30:01
 * @FilePath: \apinto\projects\core\src\app\layout\redirect-page\redirect-page.component.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { Component, OnInit } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgNavigationService } from '../../service/eo-ng-navigation.service'
@Component({
  selector: 'redirect-page',
  template: '<div></div>',
  styles: ['']
})
export class RedirectPage implements OnInit {
  constructor (private router:Router, private navigation:EoNgNavigationService) {

  }

  ngOnInit (): void {
    this.router.navigate([this.navigation.getPageRoute()])
  }
}

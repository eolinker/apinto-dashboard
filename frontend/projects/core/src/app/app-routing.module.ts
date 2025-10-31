/*
 * @Author: maggieyyy im.ymj@hotmail.com
 * @Date: 2022-07-11 23:20:14
 * @LastEditors: maggieyyy
 * @LastEditTime: 2023-11-17 09:46:22
 * @FilePath: \apinto\projects\core\src\app\app-routing.module.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { NgModule } from '@angular/core'
import { RouterModule, Routes } from '@angular/router'
import { RedirectPageService } from './service/redirect-page.service'
const routes: Routes = [
]

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
  providers: [RedirectPageService]
})
export class AppRoutingModule { }

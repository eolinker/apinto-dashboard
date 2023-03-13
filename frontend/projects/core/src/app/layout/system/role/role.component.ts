/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2022-07-30 00:40:51
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2022-09-03 20:05:12
 * @FilePath: /apinto/src/app/layout/system/system-role/system-role.component.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { Component } from '@angular/core'
import { AppConfigService } from '../../../service/app-config.service'

@Component({
  selector: 'app-system-role',
  templateUrl: './role.component.html',
  styles: [
  ]
})
export class SystemRoleComponent {
  constructor (private appConfigService: AppConfigService) {
    this.appConfigService.reqFlashBreadcrumb([{ title: '用户角色', routerLink: 'system/role' }])
  }
}

import { Component, OnInit } from '@angular/core'
import { Router } from '@angular/router'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'

@Component({
  selector: 'eo-ng-deploy-plugin-message',
  template: `
    <eo-ng-deploy-plugin-create [editPage]="true" [name]="pluginName"></eo-ng-deploy-plugin-create>
  `,
  styles: [
  ]
})
export class DeployPluginMessageComponent implements OnInit {
  pluginName:string = ''
  constructor (
    private baseInfo:BaseInfoService,
     private router:Router) {
  }

  ngOnInit (): void {
    this.pluginName = this.baseInfo.allParamsInfo.pluginName
    if (!this.pluginName) {
      this.router.navigate(['/', 'deploy', 'plugin'])
    }
  }
}

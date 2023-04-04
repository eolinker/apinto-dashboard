import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { ApiService } from 'projects/core/src/app/service/api.service';
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service';

@Component({
  selector: 'eo-ng-api-plugin-template-message',
  template: `
    <eo-ng-api-plugin-template-create [editPage]="true" [uuid]="uuid"></eo-ng-api-plugin-template-create>
  `,
  styles: [
  ]
})
export class ApiPluginTemplateMessageComponent {
  readonly nowUrl:string = this.router.routerState.snapshot.url
  uuid:string = ''

  constructor (
    private baseInfo:BaseInfoService,
    public api:ApiService,
    private router:Router
  ) {
    this.uuid = this.baseInfo.allParamsInfo.pluginTemplateId
    if (!this.uuid) {
      this.router.navigate(['/'])
    }
  }
}
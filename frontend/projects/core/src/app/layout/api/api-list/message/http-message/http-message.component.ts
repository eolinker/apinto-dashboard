import { Component } from '@angular/core'
import { ApiMessageComponent } from '../message.component'

@Component({
  selector: 'eo-ng-api-http-message',
  template: `
    <eo-ng-api-http-create [editPage]="true" [apiUuid]="apiUuid"></eo-ng-api-http-create>
  `,
  styles: [
  ]
})
export class ApiHttpMessageComponent extends ApiMessageComponent {
}

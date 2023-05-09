import { Component } from '@angular/core'
import { ApiMessageComponent } from '../message.component'

@Component({
  selector: 'eo-ng-api-websocket-message',
  template: `
       <eo-ng-api-websocket-create [editPage]="true" [apiUuid]="apiUuid"></eo-ng-api-websocket-create>

  `,
  styles: [
  ]
})
export class ApiWebsocketMessageComponent extends ApiMessageComponent {

}

import { Component } from '@angular/core'

@Component({
  selector: 'eo-ng-system',
  template: `
  <router-outlet></router-outlet>`,
  styles: [
    `
    :host ::ng-deep{
        td {
          eo-ng-select.ant-select,
          eo-ng-select-top-control.ant-select-selector {
            width: 100% !important;
          }
        }
    }`
  ]
})
export class SystemComponent {
}

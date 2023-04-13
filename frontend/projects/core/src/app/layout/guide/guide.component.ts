import { Component } from '@angular/core'

@Component({
  selector: 'eo-ng-guide',
  templateUrl: './guide.component.html',
  styles: [
    `
    :host ::ng-deep{
      height: 100%;
      width: 100%;
      display: block;
      background-color: #f5f7fa;
      overflow: hidden;
    }`
  ]
})
export class GuideComponent {

}

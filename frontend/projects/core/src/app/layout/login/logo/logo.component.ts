import { Component, Input } from '@angular/core'

type LogoType = 'light' | 'dark'

@Component({
  selector: 'eo-ng-apinto-logo',
  templateUrl: './logo.component.html',
  styleUrls: ['./logo.component.scss']
})
export class LogoComponent {
  @Input() type: LogoType = 'light'
}

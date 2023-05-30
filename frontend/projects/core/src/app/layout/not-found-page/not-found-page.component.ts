import { Component } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgNavigationService } from '../../service/eo-ng-navigation.service'

@Component({
  selector: 'eo-ng-not-found-page',
  template: '',
  styles: [
  ]
})
export class NotFoundPageComponent {
  constructor (private router:Router, private navigation:EoNgNavigationService) {
    this.router.navigate([this.navigation.getPageRoute()])
  }
}

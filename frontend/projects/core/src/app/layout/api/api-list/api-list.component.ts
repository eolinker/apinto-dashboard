import { Component } from '@angular/core'
import { Router } from '@angular/router'

@Component({
  selector: 'eo-ng-api-list',
  template: `
  <router-outlet></router-outlet>
  `,
  styles: [
  ]
})
export class ApiListComponent {
  constructor (private router:Router) {
    if (this.router.url === '/router/api') {
      this.router.navigate(['/', 'router', 'api', 'group', 'list'])
    }
  }

  ngDoCheck () {
    if (this.router.url === '/router/api') {
      this.router.navigate(['/', 'router', 'api', 'group', 'list'])
    }
  }
}

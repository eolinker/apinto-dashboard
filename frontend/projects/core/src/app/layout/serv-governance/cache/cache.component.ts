import { Component } from '@angular/core'
import { Router } from '@angular/router'

@Component({
  selector: 'eo-ng-cache',
  template: `
  <router-outlet></router-outlet>
  `,
  styles: [
  ]
})
export class CacheComponent {
  constructor (private router:Router) {
    if (this.router.url === '/serv-governance/cache') {
      this.router.navigate(['/', 'serv-governance', 'cache', 'group'])
    }
  }

  ngDoCheck () {
    if (this.router.url === '/serv-governance/cache') {
      this.router.navigate(['/', 'serv-governance', 'cache', 'group'])
    }
  }
}

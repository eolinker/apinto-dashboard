import { Component } from '@angular/core'
import { Router } from '@angular/router'

@Component({
  selector: 'eo-ng-traffic',
  template: `
  <router-outlet></router-outlet>
  `,
  styles: [
  ]
})
export class TrafficComponent {
  constructor (private router:Router) {
    if (this.router.url === '/serv-governance/traffic') {
      this.router.navigate(['/', 'serv-governance', 'traffic', 'group'])
    }
  }

  ngDoCheck () {
    if (this.router.url === '/serv-governance/traffic') {
      this.router.navigate(['/', 'serv-governance', 'traffic', 'group'])
    }
  }
}

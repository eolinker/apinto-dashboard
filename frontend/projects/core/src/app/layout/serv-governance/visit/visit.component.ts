import { Component } from '@angular/core'
import { Router } from '@angular/router'

@Component({
  selector: 'eo-ng-visit',
  template: `
  <router-outlet></router-outlet>
  `,
  styles: [
  ]
})
export class VisitComponent {
  constructor (private router:Router) {
    if (this.router.url === '/serv-governance/visit') {
      this.router.navigate(['/', 'serv-governance', 'visit', 'group'])
    }
  }

  ngDoCheck () {
    if (this.router.url === '/serv-governance/visit') {
      this.router.navigate(['/', 'serv-governance', 'visit', 'group'])
    }
  }
}

import { Component } from '@angular/core'
import { Router } from '@angular/router'

@Component({
  selector: 'eo-ng-grey',
  template: `
  <router-outlet></router-outlet>
  `,
  styles: [
  ]
})
export class GreyComponent {
  constructor (private router:Router) {
    if (this.router.url === '/serv-governance/grey') {
      this.router.navigate(['/', 'serv-governance', 'grey', 'group'])
    }
  }

  ngDoCheck () {
    if (this.router.url === '/serv-governance/grey') {
      this.router.navigate(['/', 'serv-governance', 'grey', 'group'])
    }
  }
}

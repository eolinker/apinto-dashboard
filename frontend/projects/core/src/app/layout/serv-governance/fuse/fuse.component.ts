import { Component } from '@angular/core'
import { Router } from '@angular/router'

@Component({
  selector: 'eo-ng-fuse',
  template: `
  <router-outlet></router-outlet>
  `,
  styles: [
  ]
})
export class FuseComponent {
  constructor (private router:Router) {
    if (this.router.url === '/serv-governance/fuse') {
      this.router.navigate(['/', 'serv-governance', 'fuse', 'group'])
    }
  }

  ngDoCheck () {
    if (this.router.url === '/serv-governance/fuse') {
      this.router.navigate(['/', 'serv-governance', 'fuse', 'group'])
    }
  }
}

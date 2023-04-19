import { Component } from '@angular/core'
import { Router } from '@angular/router'

@Component({
  selector: 'eo-ng-plugin-management',
  template: `
  <router-outlet></router-outlet>`,
  styles: [
  ]
})
export class PluginManagementComponent {
  constructor (private router:Router) {
    if (this.router.url === '/plugin') {
      this.router.navigate(['/', 'plugin', 'group', 'list'])
    }
  }

  ngDoCheck () {
    if (this.router.url === '/plugin') {
      this.router.navigate(['/', 'plugin', 'group', 'list'])
    }
  }
}

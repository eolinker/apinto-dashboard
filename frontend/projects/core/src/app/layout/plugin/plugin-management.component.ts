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
    if (this.router.url === '/plugin' || this.router.url === '/plugin/list') {
      this.router.navigate(['/', 'plugin', 'list', ''])
    }
  }

  ngDoCheck () {
    if (this.router.url === '/plugin' || this.router.url === '/plugin/list') {
      this.router.navigate(['/', 'plugin', 'list', ''])
    }
  }
}

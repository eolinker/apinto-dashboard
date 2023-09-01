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
    if (this.router.url.split('?')[0] === '/module-plugin') {
      this.router.navigate(['/', 'module-plugin', 'group', 'list'])
    }
  }

  ngDoCheck () {
    if (this.router.url.split('?')[0] === '/module-plugin') {
      this.router.navigate(['/', 'module-plugin', 'group', 'list'])
    }
  }
}

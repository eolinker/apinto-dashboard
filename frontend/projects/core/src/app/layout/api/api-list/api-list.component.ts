import { Component } from '@angular/core'
import { NavigationEnd, Router } from '@angular/router'
import { Subscription } from 'rxjs'

@Component({
  selector: 'eo-ng-api-list',
  template: `
  <router-outlet></router-outlet>
  `,
  styles: [
  ]
})
export class ApiListComponent {
  private subscription: Subscription = new Subscription()

  constructor (private router:Router) {
    if (this.router.url.split('?')[0] === '/router/api') {
      this.router.navigate(['/', 'router', 'api', 'group', 'list'])
    }
    this.subscription = this.router.events.subscribe((event) => {
      if (event instanceof NavigationEnd && this.router.url.split('?')[0] === '/router/api') {
        this.router.navigate(['/', 'router', 'api', 'group', 'list'])
      }
    })
  }

  ngOnDestroy () {
    this.subscription.unsubscribe()
  }
}

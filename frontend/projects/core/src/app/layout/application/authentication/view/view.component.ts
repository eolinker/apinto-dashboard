/*
 * @Date: 2023-05-30 17:04:07
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-04-16 16:41:12
 * @FilePath: \apinto\projects\core\src\app\layout\application\authentication\view\view.component.ts
 */
import { Component, OnInit } from '@angular/core'
import { ApiService } from 'projects/core/src/app/service/api.service'

@Component({
  selector: 'eo-ng-application-authentication-view',
  template: `
      <div class="w-[100%] overflow-x-hidden">
          <section >
            <ng-container *ngFor="let detail of detailList">
              <div *ngIf="detail.key" class="flex">
                <label class="label inline-block w-[220px] min-w-[220px] text-right font-bold">{{detail.key}}：</label
                ><span class="overflow-hidden break-all leading-[32px]">{{ getValue(detail.value) }}</span>
              </div>
            </ng-container>
          </section>
        </div>
  `,
  styles: [
  ]
})
export class ApplicationAuthenticationViewComponent implements OnInit {
  detailList:Array<{key:string, value:string}> = []
  authId:string = ''
  appId:string = ''
  constructor (private api:ApiService) {}

  ngOnInit (): void {
    this.getAuthData()
  }

  getValue (value:string | string[]) {
    if (typeof value === 'string') {
      return value
    }
    if (value instanceof Array) {
      return value.join(',')
    }
    return '-'
  }

  getAuthData () {
    this.api.get('application/auth/details', { uuid: this.authId, appId: this.appId })
      .subscribe((resp:{code:number, data:{details:Array<{key:string, value:string}>}, msg:string}) => {
        if (resp.code === 0) {
          this.detailList = resp.data.details
        }
      })
  }
}

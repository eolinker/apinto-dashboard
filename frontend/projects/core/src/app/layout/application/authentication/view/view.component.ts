import { Component, OnInit } from '@angular/core'
import { ApiService } from 'projects/core/src/app/service/api.service'

@Component({
  selector: 'eo-ng-application-authentication-view',
  template: `
      <div class="w-[100%] overflow-x-hidden">
          <section >
            <ng-container *ngFor="let detail of detailList">
              <div *ngIf="detail.key" >
                <label class="label inline-block w-[120px] text-right font-bold">{{detail.key}}：</label
                ><span>{{ detail.value }}</span>
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

  getAuthData () {
    this.api.get('application/auth/details', { uuid: this.authId, appId: this.appId })
      .subscribe((resp:{code:number, data:{details:Array<{key:string, value:string}>}, msg:string}) => {
        if (resp.code === 0) {
          this.detailList = resp.data.details
        }
      })
  }
}

import { Component, OnInit } from '@angular/core'
import { ApiService } from 'projects/core/src/app/service/api.service'

@Component({
  selector: 'eo-ng-application-authentication-view',
  template: `
      <div class="w-[100%] overflow-x-hidden">
          <section class=" mb-formtop">
            <ng-container *ngFor="let detail of detailList">
              <div>
                <label class="label inline-blockw-[42px] text-right">{{detail.key}}：</label
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
  constructor (private api:ApiService) {}

  ngOnInit (): void {
    this.getAuthData()
  }

  getAuthData () {
    this.api.get('application/auth/details')
      .subscribe((resp:{code:number, data:{details:Array<{key:string, value:string}>}, msg:string}) => {
        if (resp.code === 0) {
          this.detailList = resp.data.details
        }
      })
  }
}

import { Component, EventEmitter, Inject, Input, Output, TemplateRef } from '@angular/core'
import { API_URL } from '../../service/api.service'

export type CardItem = {title:string|TemplateRef<any>, enable:boolean, desc:string, iconAddr?:string, isInner?:boolean, id:string, [k:string]:any}
@Component({
  selector: 'eo-ng-card-list',
  template: `
    <div class="p-[20px] grid gap-6 py-5 2xl:grid-cols-4 xl:grid-cols-3 lg:grid-cols-2 md:grid-cols-1">
      <ng-container *ngFor="let card of cardList">
      <div>
      <div class="relative min-w-[250px] max-w-[400px] h-[100%]">
        <nz-card [nzHoverable]="true" class='min-w-[250px] max-w-[400px] min-h-[120px] h-[100%]' (click)="handlerCardClick(card)">
          <nz-card-meta
              [nzTitle]="card.title"
              [nzDescription]="cardStatusTml"
              [nzAvatar]="avatarTemplate"
            ></nz-card-meta>


            <ng-template #cardStatusTml>
              <ng-container *ngIf="type === 'plugin'">
              <span class="mr-[8px] text-[12px] font-medium text-[#00785A] bg-[#00785A1A] px-[4px] py-[2px] leading-[20px] rounded" *ngIf="card?.isInner">Apinto 内置</span>
              </ng-container>
            </ng-template>
            <ng-template #avatarTemplate>
            <div style="height:45px; width:45px;">
              <img [src]="card.iconAddr? (urlPrefix + 'plugin/icon/' + card.id + '/'+card.iconAddr) : './assets/default-plugin-icon.svg'" alt="icon" width="35px" height="35px">
              </div>
            </ng-template>
            <p class="mt-[20px] card-desc-text">{{card.desc}}</p>
          </nz-card>

          <div class="absolute top-[14px] right-[35px]">
          <span class="text-[12px] font-medium text-[#0A89FF] bg-[#0A89FF1A] px-[4px] py-[2px] leading-[20px] rounded" *ngIf="card.enable">已启用</span>
          <span  class="text-[12px] font-medium text-[#bbbbbb] bg-[#bbbbbb1A] px-[4px] py-[2px] leading-[20px] rounded" *ngIf="!card.enable"> 未启用</span>
          </div>
          </div>
          </div>
      </ng-container>
    </div>
  `,
  styles: [
    `
    :host ::ng-deep{
      nz-card {
        &.ant-card-hoverable:hover{
          border-color:var(--primary-color);
          box-shadow:none;
        }
        > .ant-card-body{
          padding:18px 22px;
        }
        .ant-card-meta-title{
          font-size:14px;
          font-weight:bold;
          max-width:calc(100% - 50px);
        }
    }

    .ant-card-meta-detail>div:not(:last-child){
      margin-bottom:0px;
    }
      p.card-desc-text{
        font-size:12px;
        word-break: break-all;
        display: -webkit-box;
        -webkit-line-clamp: 3;
        -webkit-box-orient: vertical;
        overflow: hidden;
      }
    }`
  ]
})
export class CardListComponent {
  @Input() cardList:CardItem[] = []

  @Input() type:string = 'plugin'
  @Output() cardClick:EventEmitter<CardItem> = new EventEmitter()

  constructor (
    @Inject(API_URL) public urlPrefix:string) {
  }

  handlerCardClick (card:CardItem) {
    this.cardClick.emit(card)
  }
}

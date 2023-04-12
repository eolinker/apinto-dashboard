import { Component, EventEmitter, Inject, Input, Output } from '@angular/core'
import { API_URL } from '../../service/api.service'

export type CardItem = {title:string, enable:boolean, desc:string, iconAddr?:string, isInner?:boolean, id:string, [k:string]:any}
@Component({
  selector: 'eo-ng-card-list',
  template: `
    <div class="pl-btnbase flex flex-wrap">
      <ng-container *ngFor="let card of cardList">
        <nz-card [nzHoverable]="true" class='w-[280px] h-[150px] mr-btnrbase mb-btnrbase' (click)="handlerCardClick(card)">
          <nz-card-meta
              [nzTitle]="card.title"
              [nzDescription]="cardStatusTml"
              [nzAvatar]="avatarTemplate"
            ></nz-card-meta>

            <ng-template #cardStatusTml>
              <ng-container *ngIf="type === 'plugin'">
              <span class="mr-[8px] font-medium text-[#00785A] bg-[#00785A1A] px-[4px] py-[2px] leading-[20px] rounded" *ngIf="card?.isInner">内置</span>
                <span class="mr-[8px] font-medium text-[#0A89FF] bg-[#0A89FF1A] px-[4px] py-[2px] leading-[20px] rounded" *ngIf="card.enable">已启用</span>
                <span  class="mr-[8px] font-medium text-[#bbbbbb] bg-[#bbbbbb1A] px-[4px] py-[2px] leading-[20px] rounded" *ngIf="!card.enable"> 未启用</span>
              </ng-container>
            </ng-template>
            <ng-template #avatarTemplate>
            <div style="height:50px; width:50px;">
              <img [src]="card.iconAddr? (urlPrefix + 'plugin/icon/' + card.id + '/'+card.iconAddr) : './assets/default-plugin-icon.svg'" alt="icon" width="50px" height="50px">
              </div>
            </ng-template>
            <p class="mt-[20px] card-desc-text">{{card.desc}}</p>
          </nz-card>
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
        }
    }
      p.card-desc-text{
        word-break: break-all;
        display: -webkit-box;
        -webkit-line-clamp: 2;
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

import { Component, EventEmitter, Input, Output } from '@angular/core'

export type CardItem = {title:string, enable:boolean, desc:string, iconAddr?:string, isInner?:boolean, [k:string]:any}
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
              <span class="mr-[16px]" *ngIf="card?.isInner">内置</span>
                <span *ngIf="card.enable"> 已启用</span>
                <span *ngIf="!card.enable"> 未启用</span>
              </ng-container>
            </ng-template>
            <ng-template #avatarTemplate>
              <img *ngIf="card.iconAddr" [src]="card.iconAddr" alt="plugin icon" width="64px" height="50px">
              <svg *ngIf="!card.iconAddr"  class="w-[64px] h-[50px]"><use href="#bug"></use></svg>
            </ng-template>
            <p class="mt-[20px] card-desc-text">{{card.desc}}</p>
          </nz-card>
      </ng-container>
    </div>
  `,
  styles: [
    `
    :host ::ng-deep{
      nz-card > .ant-card-body{
        padding:18px 22px;
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

  handlerCardClick (card:CardItem) {
    this.cardClick.emit(card)
  }
}

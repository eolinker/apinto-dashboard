import { Component, ElementRef, Input, OnInit, ViewChild } from '@angular/core'
import { IframeHttpService } from '../../service/iframe-http.service'

@Component({
  selector: 'eo-ng-iframe-page',
  template: `
    <iframe [src] ="iframeSrc | safe:'resourceUrl'" #iframe (load)="onLoad($event)" scrolling="yes">
    <p>Your browser does not support iframe.</p>
    </iframe>
  `,
  styles: [
    `
    iframe{
      width:100%;
      height:100%;
      border:none;
    }`
  ]
})
export class IframePageComponent implements OnInit {
  @ViewChild('iframe') iframe: ElementRef|undefined;
  @Input() path:string ='monitor-alarm'
  iframeSrc:string = ''
  constructor (private iframeService:IframeHttpService) {}
  ngOnInit (): void {
    // console.log(window.frames.ifr1.window)
    // console.dir(document.getElementById('ifr1')?.contentWindow)
    this.getIframeSrc()
  }

  // 打开的iframe可能需要传入header
  getIframeSrc () {
    this.iframeService.openIframe(this.path, { headers: { test: 'test' } }).subscribe((blob) => {
      this.iframeSrc = blob
    })
  }

  // ngAfterViewInit () {
  //   const doc = this.iframe!.nativeElement.contentDocument || this.iframe!.nativeElement.contentWindow
  // }

  onLoad (iframe:any) {
    console.log('---------')
    const doc = iframe!.nativeElement?.contentDocument || iframe!.nativeElement?.contentWindow
    console.log('---------', iframe, doc)
  }
}

import { Component, ViewChild } from '@angular/core'
import { WebsocketService, WsRef } from '../../../service/websocket.service'
import { EoNgCodeboxComponent } from 'eo-ng-codebox'
import { Subscription } from 'rxjs'
import { saveAs } from 'file-saver'

@Component({
  selector: 'eo-ng-tail-log',
  template: `
    <eo-ng-codebox
        #codeboxRef
        [code]="log"
        mode="toml"
        theme="monokai"
        [nzIsResizeHeight]="true"
        [oprBtns]="[]"
        [autoFormat]="false"
        [readonly]="true"
      ></eo-ng-codebox>
  `,
  styles: [
  ]
})
export class EoNgLogRetrievalTailComponent {
  @ViewChild('codeboxRef') codeboxRef:EoNgCodeboxComponent | undefined
  outputName:string = ''
  tail:boolean = true
  editPage:boolean = false
  log:string = ''
  lineCount:number = 0
  wsRef:WsRef|undefined
  tailKey:string = ''
  connected:boolean = false
  // TODO 调试地址
  url:string = 'ws://172.18.166.219:9400/apinto/log/node/tail'
  private subscription: Subscription = new Subscription()
  constructor (private ws:WebsocketService) {}

  ngOnInit () {
    this.connectWs()
  }

  ngOnDestroy () {
    this.closeConnect()
  }

  closeConnect () {
    this.log += '\n[...已中断连接...]\n'
    this.connected = false
    this.wsRef?.ws.close()
    this.subscription.unsubscribe()
  }

  clear () {
    this.log = ''
  }

  connectWs (reConnect?:boolean) {
    this.wsRef = this.ws.create(this.url)
    this.connected = true
    if (this.wsRef && reConnect) {
      this.log += '\n[...已恢复连接...]\n\n'
    }
    this.subscription = this.wsRef.wsRef.subscribe({
      next: (resp:any) => {
        if (resp !== 'connected') {
          this.log += resp
          this.lineCount++
          this.tail && this.codeboxRef?.aceEditor.renderer.scrollToLine(this.lineCount)
        }
      },
      error: (e:Event) => {
        console.log('ws连接出现错误：', e)
        this.connected = false
      },
      complete: () => {
        console.log('ws连接结束')
        this.connected = false
      }
    })
  }

  download () {
    const vDate = new Date()
    const fileName: string = `${this.outputName}_${vDate.getFullYear() + '-' + (vDate.getMonth() + 1) + '-' + vDate.getDate()}`
    saveAs(new Blob([this.log], { type: 'text/plain;charset=utf-8' }), `${fileName}.txt`)
  }
}

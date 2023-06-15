import { Injectable } from '@angular/core'
import { Subject } from 'rxjs'

export type WsRef = {
  wsRef: Subject<MessageEvent>,
  ws:WebSocket,
  send:Function,
}
@Injectable({
  providedIn: 'root'
})
export class WebsocketService {
  private subject: Subject<MessageEvent> |undefined;

  public create (url:string):WsRef {
    const ws = new WebSocket(url)
    this.subject = new Subject()

    ws.onmessage = (e:MessageEvent) => {
      this.subject?.next(e.data)
    }
    ws.onerror = (ev:Event):any => {
      this.subject?.error(ev)
    }

    ws.onclose = ():any => {
      this.subject?.complete()
    }

    return {
      wsRef: this.subject,
      ws: ws,
      send: (msg:string) => {
        if (ws.readyState === WebSocket.OPEN) {
          ws.send(JSON.stringify(msg))
        }
      }
    }
  }
}

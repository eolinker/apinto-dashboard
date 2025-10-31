import { Injectable } from '@angular/core'
import { CookieService } from 'ngx-cookie-service'

@Injectable({
  providedIn: 'root'
})
export class AsIframeService {
  private messageListener:((event: MessageEvent) => void) | undefined
  constructor (private cookieService: CookieService) {
  }

  startReceiveMessage (): void {
    this.messageListener = this.receiveToken.bind(this)
    window.addEventListener('message', this.messageListener, false)
  }

  receiveToken (event: MessageEvent): void {
    try {
      // 解析 JSON 字符串
      const data = JSON.parse(event.data)
      console.log('Received JSON object:', data)

      if (data && data.driver === 'quanzhi' && data.token) {
        this.cookieService.set('quanzhi_token', data.token, { path: '/' })
      }
      // 处理 JSON 对象
    } catch (error) {
      console.error('Failed to parse message:', error)
    }
  }

  removeReceiveMessage (): void {
    if (this.messageListener) {
      window.removeEventListener('message', this.messageListener, false)
    }
  }
}

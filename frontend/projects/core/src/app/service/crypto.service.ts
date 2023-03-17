import { Injectable } from '@angular/core'
import * as CryptoJS from 'crypto-js'

@Injectable({
  providedIn: 'root'
})
export class CryptoService {
  private key:string = '1e42=7838a1vfc6n'

  // AES加密，登录、修改密码用, 向量默认 1e42=7838a1vfc6n
  encryptByEnAES (inputKey:string, inputData:string, inputIv?:string):string {
    const tmpIv = CryptoJS.enc.Latin1.parse(inputIv || this.key)
    const tmpKey = CryptoJS.enc.Latin1.parse(CryptoJS.MD5(inputKey).toString() || '')
    return CryptoJS.enc.Base64.stringify(CryptoJS.AES.encrypt(inputData, tmpKey, {
      iv: tmpIv,
      mode: CryptoJS.mode.CBC,
      padding: CryptoJS.pad.Pkcs7
    }).ciphertext)
  }
}

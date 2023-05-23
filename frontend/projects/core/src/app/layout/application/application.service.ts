import { Injectable } from '@angular/core'

@Injectable({
  providedIn: 'root'
})
export class EoNgApplicationService {
  appName:string = ''
  appDesc:string = ''
  constructor () { }
}

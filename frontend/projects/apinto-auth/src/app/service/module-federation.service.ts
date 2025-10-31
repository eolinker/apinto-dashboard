import { Injectable } from '@angular/core'

@Injectable({
  providedIn: 'root'
})
export class ModuleFederationService {
  providerFromCore:any
  initialized:boolean = false
  constructor () { }
}

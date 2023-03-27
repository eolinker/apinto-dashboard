import { Injectable } from '@angular/core'

@Injectable({
  providedIn: 'root'
})
export class EoNgJsonService {
  constructor () { }

  dataType = ['string', 'boolean', 'number']
  jsonKey = ['eo:type', 'dependencies', 'skill', 'switch', 'label', 'ui:sort', 'properties', 'type'] // json schema自定义的关键字

  handleJsonSchema2Json (data:any) {
    const obj:{[k:string]:any} = {}
    for (const key in data.properties) {
      const param = data.properties[key]
      const type = param.type
      if (!this.jsonKey.includes(key)) {
        obj[key] = {}
        if (this.dataType.includes(type)) {
          obj[key] = param.default || ''
        } else if (type === 'array') {
          const items = param.items
          if (this.dataType.includes(items.type)) {
            obj[key] = [items.default || items.type === 'string' ? '' : (items.type === 'number' ? 0 : (items.type === 'boolean' ? false : null))]
          } else {
            obj[key] = [this.handleJsonSchema2Json(items)]
          }
        } else if (type === 'object') {
          obj[key] = this.handleJsonSchema2Json(param)
        } else if (type === 'number') {
          obj[key] = 0
        } else if (type === 'boolean') {
          obj[key] = false
        }
      }
    }
    return obj
  }
}

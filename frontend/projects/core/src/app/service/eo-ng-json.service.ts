import { Injectable } from '@angular/core'

@Injectable({
  providedIn: 'root'
})
export class EoNgJsonService {
  dataType = ['string', 'boolean', 'number', 'integer']
  jsonKey = ['eo:type', 'dependencies', 'skill', 'switch', 'label', 'ui:sort', 'properties', 'type'] // json schema自定义的关键字

  handleJsonSchema2Json (data:any) {
    const obj:{[k:string]:any} = {}
    for (const key in data.properties) {
      const param = data.properties[key]
      const type = param.type
      if (!this.jsonKey.includes(key)) {
        obj[key] = {}
        if (!this.dataType.includes(type) && type !== 'object' && type !== 'array') {
          obj[key] = param.default || ''
        } else if (type === 'array') {
          const items = param.items
          if (this.dataType.includes(items.type)) {
            obj[key] = items.required === false
              ? items.default || null
              : items.default || [(
                items.type === 'string'
                  ? ((items.required !== false && items?.enum?.length > 0) ? items.enum[0] : (items.required === false ? null : ''))
                  : (items.type === 'number' ? ((items.required !== false && items.enum.length > 0) ? items.enum[0] : (items.required === false ? null : 0)) : (items.type === 'boolean' ? ((items.required !== false && items.enum.length > 0) ? items.enum[0] : (items.required === false ? null : false)) : null)))]
          } else {
            obj[key] = [this.handleJsonSchema2Json(items)]
          }
        } else if (type === 'object') {
          obj[key] = this.handleJsonSchema2Json(param)
        } else if (type === 'number' || type === 'integer') {
          obj[key] = ((param.required !== false && param?.enum?.length > 0) ? param.enum[0] : (param.required === false ? null : 0))
        } else if (type === 'boolean') {
          obj[key] = ((param.required !== false && param?.enum?.length > 0) ? param.enum[0] : (param.required === false ? null : false))
        } else if (type === 'string') {
          obj[key] = ((param.required !== false && param?.enum?.length > 0) ? param.enum[0] : (param.required === false ? null : ''))
        }
      }
    }
    return obj
  }
}

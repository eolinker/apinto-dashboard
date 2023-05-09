import { Pipe, PipeTransform } from '@angular/core'

@Pipe({
  name: 'changeWordColor'
})
export class ChangeWordColorPipe implements PipeTransform {
  // value:需转换的字符串，args[0]：需改变颜色的字符串，args[1]：改变后的颜色
  transform (value: string, ...args:Array<string>): unknown {
    const wordArray:Array<string> = value.split(args[0])
    return wordArray.join(`<span style="color:${args[1]}">${args[0]}</span>`)
  }
}

import { Injectable } from '@angular/core'
import * as ExcelJS from 'exceljs'
import * as FileSaver from 'file-saver'

@Injectable({
  providedIn: 'root'
})
export class EoNgExcelService {
  createExcel (sheetTitle:string, columns:Array<{header:string, key:string, width:number}>, tableData:any) {
    const workBook = new ExcelJS.Workbook()
    const sheet = workBook.addWorksheet(sheetTitle || '默认工作表')
    sheet.columns = columns
    sheet.addRows(tableData)
    return workBook
  }

  exportExcel (fileTitle:string, date:Array<number>, sheetTitle:string, tableConfig:{[key:string]:boolean}, tableHead:Array<any>, tableData:any) {
    tableData = tableData || []
    const columns = this.getColumns(tableConfig, tableHead)
    const workBook = this.createExcel(sheetTitle, columns, tableData)
    const fileName = this.getFileName(fileTitle, date)
    workBook.xlsx.writeBuffer().then(function (buffer) {
      FileSaver(
        new Blob([buffer], {
          type: 'application/octet-stream'
        }),
        `${fileName}.xlsx`
      )
    })
  }

  getColumns (tableConfig:{[key:string]:boolean}, tableHead:Array<any>) {
    const columns = []
    for (const index in tableHead) {
      if (tableConfig[tableHead[index].key]) {
        if (tableHead[index].title === 'request_rate' || tableHead[index].title === 'proxy_rate') {
          columns.push({ header: tableHead[index].title, key: tableHead[index].key, width: tableHead[index].title > 5 ? tableHead[index].title * 3 : 15, style: { numFmt: '0.00%' } })
        } else {
          columns.push({ header: tableHead[index].title, key: tableHead[index].key, width: tableHead[index].title > 5 ? tableHead[index].title * 3 : 15 })
        }
      }
    }
    return columns
  }

  getFileName (fileTitle:string, date:Array<number>):string {
    return `${fileTitle}-${this.getDateFormat(date[0])}至${this.getDateFormat(date[1])}`
  }

  getDateFormat (time:number):string {
    const date = new Date(time * 1000)
    return `${date.getFullYear()}${date.getMonth() < 9 ? '0' + (date.getMonth() + 1) : (date.getMonth() + 1)}${date.getDate() < 10 ? '0' + date.getDate() : date.getDate()}-${date.getHours() < 10 ? '0' + date.getHours() : date.getHours()}${date.getMinutes() < 10 ? '0' + date.getMinutes() : date.getMinutes()}`
  }
}

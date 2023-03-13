/* eslint-disable camelcase */
/* eslint-disable dot-notation */
import { Component, Input, OnInit, SimpleChanges, TemplateRef, ViewChild } from '@angular/core'
import { EChartsOption } from 'echarts'
// 监控告警中用到的折线图封装
// 监控总览：y轴左右两侧各有标题，标题会随x轴变化而更新（时间跨度）
// 调用量统计，共六条线，默认出现四条实线
// 报文量统计，共两条线
// 调用趋势：六条线，并且会出现需要加入对比的可能
// 预留数据给标题，标题根据需求
export interface InvokeData{
  date:Array<string>
  request_total:Array<string>
  request_rate:Array<string>
  proxy_total:Array<string>
  proxy_rate:Array<string>
  status_4xx:Array<string>
  status_5xx:Array<string>
  time_interval?:string
  [key:string]:any
}

export interface MessageData{
  date:Array<string>
  request:Array<string>
  response:Array<string>
  [key:string]:any
}

@Component({
  selector: 'eo-ng-monitor-line-graph',
  template: `
    <div echarts (chartInit)="chartInit($event)"  [options]="lineChartOption" (chartLegendSelectChanged)="legendselectchanged($event)" style="min-width: 100%;" [ngStyle]="{'height': compare? '672px':'336px'}"></div>
  `,
  styles: [
  ]
})
export class MonitorLineGraphComponent implements OnInit {
  @ViewChild('lineChart') lineChart: TemplateRef<any> | undefined;
  @Input() lineData:InvokeData | MessageData| {[key:string]:any} = {}
  @Input() compareData:InvokeData | MessageData | {} = {}
  @Input() titles:Array<string> = ['']
  @Input() yAxisTitle:string = ''// 调用量统计时，y周显示的统计粒度
  @Input() compare:boolean = false
  @Input() type:'invoke'|'invoke-upstream'|'traffic'|'' = 'invoke'
  @Input() modalTitle:string = '' // 折线图所在modal的名称，用在对比图中的标签名
  @Input() dataTitle:string = '' // 折线图所代表的维度的名称，比如查看上游调用某应用的调用趋势，这里就是某上游的名称

   maxYNumber:number = 1
   lineChartOption: EChartsOption = {
   }

   lineNameMap:{[key:string]:string} = {
     request_total: '请求总数',
     request_rate: '请求成功率',
     proxy_total: '转发总数',
     proxy_rate: '转发成功率',
     status_4xx: '状态码4xx数',
     status_5xx: '状态码5xx数',
     request: '请求报文量',
     response: '响应报文量'
   }

  lineChartOptionConfig: EChartsOption = {
    legend: {
      orient: 'horizontal',
      top: '50',
      selected: {
        请求总数: true,
        请求成功率: true,
        转发总数: true,
        转发成功率: true,
        状态码4xx数: false,
        状态码5xx数: false
      }
    },
    tooltip: {
      trigger: 'axis',
      // 为了失败率显示成百分比，所以自定义了formatter
      formatter: (params:any) => {
        const startHtml = params[0].axisValue + '<br/>'
        const listArr = []
        for (let i = 0; i < params.length; i++) {
          const item = params[i]
          // echarts会根据你定义的颜色返回一个生成好的带颜色的标记，直接实用即可
          let str = '<div><section style="align-items: center;display:flex; justify-content: space-between;flex-wrap: nowrap;"><span> ' + item.marker
          if (item.seriesName === '请求成功率' || item.seriesName === '转发成功率') {
            str += (item.seriesName + '&nbsp&nbsp&nbsp </span><span style="font-weight:bold"> ' + item.value + '% </span></section></div>')
          } else {
            str += (item.seriesName + '&nbsp&nbsp&nbsp </span><span style="font-weight:bold"> ' + item.value + '</span></section></div>')
          }
          listArr.push(str)
        }
        return startHtml + listArr.join('')
      }
    }
  };

  echartsIntance: any

  ngOnInit (): void {
    this.changeLineChart()
    this.getMaxY()
    this.echartsIntance?.setOption(this.lineChartOption)
  }

  ngOnChanges (changes:SimpleChanges):void {
    if (changes['lineData'] || changes['titles'] || changes['yAxisTitle']) {
      this.changeLineChart()
      this.getMaxY()
      this.echartsIntance?.setOption(this.lineChartOption)
    }
    if (changes['compare']) {
      this.changeLineChart()
      this.lineChartOption.legend = { ...this.lineChartOption.legend, selected: { ...this.echartsIntance.getOption().legend[0].selected } }
      this.getMaxY()
      this.echartsIntance?.setOption(this.lineChartOption)
    }
  }

  ngAfterViewInit () {
  }

  chartInit (event:any) {
    this.echartsIntance = event
  }

  legendselectchanged (value:any) {
    // 当勾选请求成功率或转发成功率其中之一时，显示右侧y轴
    if (this.lineData.date && this.lineData.date.length > 0) {
      if (!value.selected['转发成功率'] && !value.selected['请求成功率'] && (this.lineChartOption.yAxis as Array<any>)?.length > 1 && (this.lineChartOption.yAxis as Array<any>)[1].show !== false) {
        (this.lineChartOption.yAxis as Array<any>)[1].show = false
      } else if ((value.selected['转发成功率'] || value.selected['请求成功率']) && (this.lineChartOption.yAxis as Array<any>)?.length > 1 && (this.lineChartOption.yAxis as Array<any>)[1].show !== true) {
        (this.lineChartOption.yAxis as Array<any>)[1].show = true
      }
    }
    this.lineChartOption.legend = { ...this.lineChartOption.legend, selected: value.selected }
    this.getMaxY()
    this.echartsIntance?.setOption(this.lineChartOption)
  }

  getTimeFormatter (time?:any):string {
    switch (this.yAxisTitle) {
      case '分钟': {
        return `${new Date(time).getMonth() < 9 ? '0' + (new Date(time).getMonth() + 1) : (new Date(time).getMonth() + 1)}/${new Date(time).getDate() < 10 ? '0' + new Date(time).getDate() : new Date(time).getDate()} ${new Date(time).getHours() < 10 ? '0' + new Date(time).getHours() : new Date(time).getHours()}:${new Date(time).getMinutes() < 10 ? '0' + new Date(time).getMinutes() : new Date(time).getMinutes()}`
      }
      case '5分钟': {
        return `${new Date(time).getMonth() < 9 ? '0' + (new Date(time).getMonth() + 1) : (new Date(time).getMonth() + 1)}/${new Date(time).getDate() < 10 ? '0' + new Date(time).getDate() : new Date(time).getDate()} ${new Date(time).getHours() < 10 ? '0' + new Date(time).getHours() : new Date(time).getHours()}:${new Date(time).getMinutes() < 10 ? '0' + new Date(time).getMinutes() : new Date(time).getMinutes()}   `
      }
      case '1小时': {
        return `${new Date(time).getMonth() < 9 ? '0' + (new Date(time).getMonth() + 1) : (new Date(time).getMonth() + 1)}/${new Date(time).getDate() < 10 ? '0' + new Date(time).getDate() : new Date(time).getDate()} ${new Date(time).getHours() < 10 ? '0' + new Date(time).getHours() : new Date(time).getHours()}     `
      }
      case '1天': {
        return `${new Date(time).getFullYear().toString().slice(2)}年-${new Date(time).getMonth() < 9 ? '0' + (new Date(time).getMonth() + 1) : (new Date(time).getMonth() + 1)}/${new Date(time).getDate() < 10 ? '0' + new Date(time).getDate() : new Date(time).getDate()}      `
      }
      case '1周': {
        return `${new Date(time).getFullYear().toString().slice(2)}年-${new Date(time).getMonth() < 9 ? '0' + (new Date(time).getMonth() + 1) : (new Date(time).getMonth() + 1)}/${new Date(time).getDate() < 10 ? '0' + new Date(time).getDate() : new Date(time).getDate()}      `
      }
    }
    return `${new Date(time).getMonth() + 1}/${new Date(time).getDate() < 10 ? '0' + new Date(time).getDate() : new Date(time).getDate()} ${new Date(time).getHours() < 10 ? '0' + new Date(time).getHours() : new Date(time).getHours()}:${new Date(time).getMinutes() < 10 ? '0' + new Date(time).getMinutes() : new Date(time).getMinutes()}`
  }

  getMaxY (legendSeleted?:any) {
    (this.lineChartOption.yAxis as any[])[0] = this.changeInterval(this.lineData, (this.lineChartOption.yAxis as any[])[0], legendSeleted)
    if (this.compare) {
      (this.lineChartOption.yAxis as any[])[2] = this.changeInterval(this.compareData, (this.lineChartOption.yAxis as any[])[2], legendSeleted)
    }
  }

  yUnitFormatter (value:number):string {
    let res:string = ''
    if (value > 100000000) {
      res = (value / 100000000).toFixed(2) + '亿'
    } else if (value > 1000000) {
      res = (value / 10000).toFixed(0) + '万'
    } else if (value > 100000) {
      res = (value / 10000).toFixed(2) + '万'
    } else {
      res = value.toFixed(0)
    }
    return res
  }

  changeInterval (data:any, yAxis:any, legendSeleted?:any) {
    const maxNumberObj:{[key:string]:number} = {}
    for (const key of Object.keys(data)) {
      if (key !== 'date' && key !== 'request_rate' && key !== 'proxy_rate') {
        maxNumberObj[key] = Math.max(...(data[key as any] || []))
      }
    }
    const { legend } = this.lineChartOption
    const selected = legendSeleted || (legend as any).selected
    const numArr:Array<number> = []
    for (const key of Object.keys(maxNumberObj)) {
      if ((selected as any)[this.lineNameMap[key]]) {
        numArr.push(maxNumberObj[key])
      }
    }
    // if (numArr.length > 0) {
    //   return Math.max(...numArr)
    // }
    const maxNum:number = Math.max(...numArr) || 0
    yAxis = { ...yAxis, interval: maxNum > 5 ? (maxNum / 5) : 1, max: maxNum > 5 ? maxNum : 5 }
    return yAxis
  }

  changeLineChart () {
    switch (this.type) {
      case 'invoke': {
        if (!this.compare) {
          this.lineChartOption = {
            ...this.lineChartOptionConfig,
            grid: {
              left: '16',
              right: '16',
              bottom: '16',
              top: '120',
              containLabel: true
            },
            xAxis: {
              type: 'category',
              data: (this.lineData as InvokeData).date?.map((x:string) => {
                return `${new Date(x).getFullYear()}/${new Date(x).getMonth() + 1}/${new Date(x).getDate()} ${new Date(x).getHours()}:${new Date(x).getMinutes()}              `
              }) || [],
              axisLabel: {
                showMaxLabel: true,
                formatter: (value:any) => {
                  return this.getTimeFormatter(value)
                }
              },
              axisTick: {
                show: false
              },
              boundaryGap: false
            },
            yAxis: [{
              type: 'value',
              name: (this.lineData as InvokeData).date.length > 0 ? `${this.yAxisTitle}调用量` : '',
              nameLocation: 'end',
              nameTextStyle: {
                align: 'left'
              },
              min: 0,
              max: 'dataMax',
              axisLabel: {
                formatter: (value:any) => {
                  return this.yUnitFormatter(value)
                }
              }
            },
            {
              type: 'value',
              name: (this.lineData as InvokeData).date.length > 0 ? `${this.yAxisTitle}调用成功率` : '',
              position: 'right',
              min: 0,
              max: 100,
              show: (this.lineData as InvokeData).date.length > 0,
              interval: 20,
              axisLabel: {
                formatter: '{value} %'
              },
              axisLine: {
                show: false
              },
              axisTick: {
                show: false
              }
            }],
            series: [
              { type: 'line', symbol: 'none', name: '请求总数', data: (this.lineData as InvokeData).request_total, yAxisIndex: 0 },
              { type: 'line', symbol: 'none', name: '请求成功率', data: (this.lineData as InvokeData).request_rate?.map((x) => Number((Number(x) * 100).toFixed(2))) || [], yAxisIndex: 1 },
              { type: 'line', symbol: 'none', name: '转发总数', data: (this.lineData as InvokeData).proxy_total, yAxisIndex: 0 },
              { type: 'line', symbol: 'none', name: '转发成功率', data: (this.lineData as InvokeData).proxy_rate?.map((x) => Number((Number(x) * 100).toFixed(2))) || [], yAxisIndex: 1 },
              { type: 'line', symbol: 'none', lineStyle: { type: 'dashed' }, name: '状态码4xx数', data: (this.lineData as InvokeData).status_4xx, yAxisIndex: 0 },
              { type: 'line', symbol: 'none', lineStyle: { type: 'dashed' }, name: '状态码5xx数', data: (this.lineData as InvokeData).status_5xx, yAxisIndex: 0 }
            ]
          }
        } else {
          this.lineChartOption = {
            ...this.lineChartOptionConfig,
            tooltip: {
              formatter: (params:any) => {
                const startHtml = '<div><section style="align-items: center;display:flex; justify-content: space-between;flex-wrap: nowrap;"><span> ' + (params[0].seriesIndex === 0 ? this.modalTitle : this.dataTitle + '调用总体趋势') + '</span>&nbsp&nbsp&nbsp<span>' + params[0].axisValue + '</span></div>'
                const listArr = []
                for (let i = 0; i < params.length; i++) {
                  const item = params[i]
                  // echarts会根据你定义的颜色返回一个生成好的带颜色的标记，直接实用即可
                  let str = ''
                  if (i === Math.floor(params.length / 2)) {
                    str = '<br/><div><section style="align-items: center;display:flex; justify-content: space-between;flex-wrap: nowrap;"><span> ' + (params[0].seriesIndex === 0 ? this.dataTitle + '调用总体趋势' : this.modalTitle) + '</span>&nbsp&nbsp&nbsp<span>' + params[0].axisValue + '</span></div><div><section style="align-items: center;display:flex; justify-content: space-between;flex-wrap: nowrap;"><span> ' + item.marker
                  } else {
                    str = '<div><section style="align-items: center;display:flex; justify-content: space-between;flex-wrap: nowrap;"><span> ' + item.marker
                  }
                  if (item.seriesName === '请求成功率' || item.seriesName === '转发成功率') {
                    str += (item.seriesName + '&nbsp&nbsp&nbsp </span><span style="font-weight:bold"> ' + item.value + '% </span></section></div>')
                  } else {
                    str += (item.seriesName + '&nbsp&nbsp&nbsp </span><span style="font-weight:bold"> ' + item.value + '</span></section></div>')
                  }
                  listArr.push(str)
                }
                return startHtml + listArr.join('')
              },
              trigger: 'axis',
              axisPointer: {
                link: [{
                  xAxisIndex: 'all'
                }]
              }
            },
            axisPointer: {
              link: [{
                xAxisIndex: 'all'
              }]
            },
            grid: [
              {
                left: '16',
                right: '16',
                top: '120',
                containLabel: true,
                bottom: '50%'
              },
              {
                left: '16',
                right: '16',
                bottom: '16',
                containLabel: true,
                top: '60%'
              }
            ],
            xAxis: [{
              type: 'category',
              data: (this.lineData as InvokeData).date?.map((x:string) => {
                return `${new Date(x).getFullYear()}/${new Date(x).getMonth() + 1}/${new Date(x).getDate()} ${new Date(x).getHours()}:${new Date(x).getMinutes()}              `
              }) || [],
              axisLabel: {
                showMaxLabel: true,
                formatter: (value:any) => {
                  return this.getTimeFormatter(value)
                }
              },
              axisTick: {
                show: false
              }
            },
            {
              gridIndex: 1,
              type: 'category',
              data: (this.lineData as InvokeData).date?.map((x:string) => {
                return `${new Date(x).getFullYear()}/${new Date(x).getMonth() + 1}/${new Date(x).getDate()} ${new Date(x).getHours()}:${new Date(x).getMinutes()}              `
              }) || [],
              axisLabel: {
                showMaxLabel: true,
                formatter: (value:any) => {
                  return this.getTimeFormatter(value)
                }
              }
            }
            ],
            yAxis: [
              {
                type: 'value',
                name: (this.lineData as InvokeData).date.length > 0 ? `${this.yAxisTitle}调用量` : '',
                nameLocation: 'end',
                nameTextStyle: {
                  align: 'left'
                },
                min: 0,
                max: 'dataMax',
                axisLabel: {

                  formatter: (value:any) => {
                    return this.yUnitFormatter(value)
                  }
                }
              },
              {
                type: 'value',
                name: (this.lineData as InvokeData).date.length > 0 ? `${this.yAxisTitle}调用成功率` : '',
                position: 'right',
                min: 0,
                max: 100,
                show: (this.lineData as InvokeData).date.length > 0,

                interval: 20,
                axisLabel: {
                  formatter: '{value} %'
                },
                axisLine: {
                  show: false
                },
                axisTick: {
                  show: false
                }
              },
              {
                gridIndex: 1,
                type: 'value',
                name: (this.lineData as InvokeData).date.length > 0 ? `${this.yAxisTitle}调用量` : '',
                nameLocation: 'end',
                nameTextStyle: {
                  align: 'left'
                },
                min: 0,
                max: 'dataMax',
                axisLabel: {

                  formatter: (value:any) => {
                    return this.yUnitFormatter(value)
                  }
                }
              },
              {
                gridIndex: 1,
                type: 'value',
                name: (this.lineData as InvokeData).date.length > 0 ? `${this.yAxisTitle}调用成功率` : '',
                position: 'right',
                min: 0,
                max: 100,
                show: (this.lineData as InvokeData).date.length > 0,

                interval: 20,
                axisLabel: {
                  formatter: '{value} %'
                },
                axisLine: {
                  show: false
                },
                axisTick: {
                  show: false
                }
              }],
            series: [
              { type: 'line', symbol: 'none', name: '请求总数', data: (this.lineData as InvokeData).request_total, xAxisIndex: 0, yAxisIndex: 0 },
              { type: 'line', symbol: 'none', name: '请求成功率', data: (this.lineData as InvokeData).request_rate?.map((x) => Number((Number(x) * 100).toFixed(2))) || [], xAxisIndex: 0, yAxisIndex: 1 },
              { type: 'line', symbol: 'none', name: '转发总数', data: (this.lineData as InvokeData).proxy_total, xAxisIndex: 0, yAxisIndex: 0 },
              { type: 'line', symbol: 'none', name: '转发成功率', data: (this.lineData as InvokeData).proxy_rate?.map((x) => Number((Number(x) * 100).toFixed(2))) || [], xAxisIndex: 0, yAxisIndex: 1 },
              { type: 'line', lineStyle: { type: 'dashed' }, symbol: 'none', name: '状态码4xx数', data: (this.lineData as InvokeData).status_4xx, xAxisIndex: 0, yAxisIndex: 0 },
              { type: 'line', lineStyle: { type: 'dashed' }, symbol: 'none', name: '状态码5xx数', data: (this.lineData as InvokeData).status_5xx, xAxisIndex: 0, yAxisIndex: 0 },
              { type: 'line', symbol: 'none', name: '请求总数', data: (this.compareData as InvokeData).request_total, xAxisIndex: 1, yAxisIndex: 2 },
              { type: 'line', symbol: 'none', name: '请求成功率', data: (this.compareData as InvokeData).request_rate?.map((x) => Number((Number(x) * 100).toFixed(2))) || [], xAxisIndex: 1, yAxisIndex: 3 },
              { type: 'line', symbol: 'none', name: '转发总数', data: (this.compareData as InvokeData).proxy_total, xAxisIndex: 1, yAxisIndex: 2 },
              { type: 'line', symbol: 'none', name: '转发成功率', data: (this.compareData as InvokeData).proxy_rate?.map((x) => Number((Number(x) * 100).toFixed(2))) || [], xAxisIndex: 1, yAxisIndex: 3 },
              { type: 'line', lineStyle: { type: 'dashed' }, symbol: 'none', name: '状态码4xx数', data: (this.compareData as InvokeData).status_4xx, xAxisIndex: 1, yAxisIndex: 2 },
              { type: 'line', lineStyle: { type: 'dashed' }, symbol: 'none', name: '状态码5xx数', data: (this.compareData as InvokeData).status_5xx, xAxisIndex: 1, yAxisIndex: 2 }
            ]
          }
        }
        break
      }
      case 'invoke-upstream': {
        if (!this.compare) {
          this.lineChartOption = {
            ...this.lineChartOptionConfig,
            legend: {
              orient: 'horizontal',
              top: '50',
              selected: {
                转发总数: true,
                转发成功率: true,
                状态码4xx数: false,
                状态码5xx数: false
              }
            },
            grid: {
              left: '16',
              right: '16',
              bottom: '16',
              top: '120',
              containLabel: true
            },
            xAxis: {
              type: 'category',
              data: (this.lineData as InvokeData).date?.map((x:string) => {
                return `${new Date(x).getFullYear()}/${new Date(x).getMonth() + 1}/${new Date(x).getDate()} ${new Date(x).getHours()}:${new Date(x).getMinutes()}              `
              }) || [],
              axisLabel: {
                showMaxLabel: true,
                formatter: (value:any) => {
                  return this.getTimeFormatter(value)
                }
              },
              axisTick: {
                show: false
              }
            },
            yAxis: [
              {
                type: 'value',
                name: (this.lineData as InvokeData).date.length > 0 ? `${this.yAxisTitle}调用量` : '',
                nameLocation: 'end',
                nameTextStyle: {
                  align: 'left'
                },

                min: 0,
                max: 'dataMax',
                axisLabel: {

                  formatter: (value:any) => {
                    return this.yUnitFormatter(value)
                  }
                }
              },
              {
                type: 'value',
                name: (this.lineData as InvokeData).date.length > 0 ? `${this.yAxisTitle}调用成功率` : '',
                position: 'right',
                min: 0,
                max: 100,
                show: (this.lineData as InvokeData).date.length > 0,
                interval: 20,
                axisLabel: {
                  formatter: '{value} %'
                },
                axisLine: {
                  show: false
                },
                axisTick: {
                  show: false
                }
              }],
            series: [
              { type: 'line', symbol: 'none', name: '转发总数', data: (this.lineData as InvokeData).proxy_total, yAxisIndex: 0 },
              { type: 'line', symbol: 'none', name: '转发成功率', data: (this.lineData as InvokeData).proxy_rate?.map((x) => Number((Number(x) * 100).toFixed(2))) || [], yAxisIndex: 1 },
              { type: 'line', lineStyle: { type: 'dashed' }, symbol: 'none', name: '状态码4xx数', data: (this.lineData as InvokeData).status_4xx, yAxisIndex: 0 },
              { type: 'line', lineStyle: { type: 'dashed' }, symbol: 'none', name: '状态码5xx数', data: (this.lineData as InvokeData).status_5xx, yAxisIndex: 0 }
            ]
          }
        } else {
          this.lineChartOption = {
            ...this.lineChartOptionConfig,
            legend: {
              orient: 'horizontal',
              top: '50',
              selected: {
                转发总数: true,
                转发成功率: true,
                状态码4xx数: false,
                状态码5xx数: false
              }
            },
            grid: [
              {
                left: '16',
                right: '16',
                top: '120',
                containLabel: true,
                bottom: '50%'
              },
              {
                left: '16',
                right: '16',
                bottom: '16',
                containLabel: true,
                top: '60%'
              }
            ],
            tooltip: {
              trigger: 'axis',
              formatter: (params:any) => {
                const startHtml = '<div><section style="align-items: center;display:flex; justify-content: space-between;flex-wrap: nowrap;"><span> ' + (params[0].seriesIndex === 0 ? this.modalTitle : this.dataTitle + '调用总体趋势') + '</span>&nbsp&nbsp&nbsp<span>' + params[0].axisValue + '</span></div>'
                const listArr = []
                for (let i = 0; i < params.length; i++) {
                  const item = params[i]
                  // echarts会根据你定义的颜色返回一个生成好的带颜色的标记，直接实用即可
                  let str = ''
                  if (i === Math.floor(params.length / 2)) {
                    str = '<br/><div><section style="align-items: center;display:flex; justify-content: space-between;flex-wrap: nowrap;"><span> ' + (params[0].seriesIndex === 0 ? this.dataTitle + '调用总体趋势' : this.modalTitle) + '</span>&nbsp&nbsp&nbsp<span>' + params[0].axisValue + '</span></div><div><section style="align-items: center;display:flex; justify-content: space-between;flex-wrap: nowrap;"><span> ' + item.marker
                  } else {
                    str = '<div><section style="align-items: center;display:flex; justify-content: space-between;flex-wrap: nowrap;"><span> ' + item.marker
                  }
                  if (item.seriesName === '请求成功率' || item.seriesName === '转发成功率') {
                    str += (item.seriesName + '&nbsp&nbsp&nbsp </span><span style="font-weight:bold"> ' + item.value + '% </span></section></div>')
                  } else {
                    str += (item.seriesName + '&nbsp&nbsp&nbsp </span><span style="font-weight:bold"> ' + item.value + '</span></section></div>')
                  }
                  listArr.push(str)
                }
                return startHtml + listArr.join('')
              }
            },
            axisPointer: {
              link: [{
                xAxisIndex: 'all'
              }]
            },
            xAxis: [{
              type: 'category',
              data: (this.lineData as InvokeData).date?.map((x:string) => {
                return `${new Date(x).getFullYear()}/${new Date(x).getMonth() + 1}/${new Date(x).getDate()} ${new Date(x).getHours()}:${new Date(x).getMinutes()}              `
              }) || [],
              axisLabel: {
                showMaxLabel: true,
                formatter: (value:any) => {
                  return this.getTimeFormatter(value)
                }
              },
              axisTick: {
                show: false
              }
            },
            {
              gridIndex: 1,
              type: 'category',
              data: (this.lineData as InvokeData).date?.map((x:string) => {
                return `${new Date(x).getFullYear()}/${new Date(x).getMonth() + 1}/${new Date(x).getDate()} ${new Date(x).getHours()}:${new Date(x).getMinutes()}              `
              }) || [],
              splitNumber: 5,
              axisLabel: {
                showMaxLabel: true,
                formatter: (value:any) => {
                  return this.getTimeFormatter(value)
                }
              },
              axisTick: {
                show: false
              }
            }],
            yAxis: [
              {
                type: 'value',
                name: (this.lineData as InvokeData).date.length > 0 ? `${this.yAxisTitle}调用量` : '',
                nameLocation: 'end',
                nameTextStyle: {
                  align: 'left'
                },
                min: 0,
                max: 'dataMax',
                axisLabel: {

                  formatter: (value:any) => {
                    return this.yUnitFormatter(value)
                  }
                }
              },
              {
                type: 'value',
                name: (this.lineData as InvokeData).date.length > 0 ? `${this.yAxisTitle}调用成功率` : '',
                position: 'right',
                min: 0,
                max: 100,
                show: (this.lineData as InvokeData).date.length > 0,

                interval: 20,
                axisLabel: {
                  formatter: '{value} %'
                },
                axisLine: {
                  show: false
                },
                axisTick: {
                  show: false
                }
              },
              {
                gridIndex: 1,
                type: 'value',
                name: (this.lineData as InvokeData).date.length > 0 ? `${this.yAxisTitle}调用量` : '',
                nameLocation: 'end',
                nameTextStyle: {
                  align: 'left'
                },
                min: 0,
                max: 'dataMax',
                axisLabel: {

                  formatter: (value:any) => {
                    return this.yUnitFormatter(value)
                  }
                }
              },
              {
                gridIndex: 1,
                type: 'value',
                name: (this.lineData as InvokeData).date.length > 0 ? `${this.yAxisTitle}调用成功率` : '',
                position: 'right',
                min: 0,
                max: 100,
                show: (this.lineData as InvokeData).date.length > 0,

                interval: 20,
                axisLabel: {
                  formatter: '{value} %'
                },
                axisLine: {
                  show: false
                },
                axisTick: {
                  show: false
                }
              }],
            series: [
              { type: 'line', symbol: 'none', name: '转发总数', data: (this.lineData as InvokeData).proxy_total, xAxisIndex: 0, yAxisIndex: 0 },
              { type: 'line', symbol: 'none', name: '转发成功率', data: (this.lineData as InvokeData).proxy_rate, xAxisIndex: 0, yAxisIndex: 1 },
              { type: 'line', symbol: 'none', lineStyle: { type: 'dashed' }, name: '状态码4xx数', data: (this.lineData as InvokeData).status_4xx, xAxisIndex: 0, yAxisIndex: 0 },
              { type: 'line', symbol: 'none', lineStyle: { type: 'dashed' }, name: '状态码5xx数', data: (this.lineData as InvokeData).status_5xx, xAxisIndex: 0, yAxisIndex: 0 },
              { type: 'line', symbol: 'none', name: '转发总数', data: (this.compareData as InvokeData).proxy_total, xAxisIndex: 1, yAxisIndex: 2 },
              { type: 'line', symbol: 'none', name: '转发成功率', data: (this.compareData as InvokeData).proxy_rate, xAxisIndex: 1, yAxisIndex: 3 },
              { type: 'line', symbol: 'none', lineStyle: { type: 'dashed' }, name: '状态码4xx数', data: (this.compareData as InvokeData).status_4xx, xAxisIndex: 1, yAxisIndex: 2 },
              { type: 'line', symbol: 'none', lineStyle: { type: 'dashed' }, name: '状态码5xx数', data: (this.compareData as InvokeData).status_5xx, xAxisIndex: 1, yAxisIndex: 2 }
            ]
          }
        }
        break
      }
      case 'traffic':
        this.lineChartOption = {
          ...this.lineChartOptionConfig,
          legend: {
            orient: 'horizontal',
            top: '50',
            selected: {
              请求报文量: true,
              响应报文量: true
            }
          },
          grid: {
            left: '16',
            right: '16',
            bottom: '16',
            top: '120',
            containLabel: true
          },
          xAxis: {
            type: 'category',
            data: (this.lineData as MessageData).date?.map((x:string) => {
              return `${new Date(x).getFullYear()}/${new Date(x).getMonth() + 1}/${new Date(x).getDate()} ${new Date(x).getHours()}:${new Date(x).getMinutes()}              `
            }) || [],
            axisLabel: {
              showMaxLabel: true,
              formatter: (value:any) => {
                return this.getTimeFormatter(value)
              }
            },
            axisTick: {
              show: false
            }
          },
          yAxis: [{
            type: 'value',
            name: (this.lineData as MessageData).date.length > 0 ? `${this.yAxisTitle}报文量（KB）` : '',
            nameLocation: 'end',
            nameTextStyle: {
              align: 'left'
            },
            axisTick: {
              length: 6
            },

            min: 0,
            max: 'dataMax',
            axisLabel: {

              formatter: (value:any) => {
                return this.yUnitFormatter(value)
              }
            }
          },
          {
            type: 'value',
            name: (this.lineData as InvokeData).date.length > 0 ? `${this.yAxisTitle}调用成功率` : '',
            position: 'right',
            min: 0,
            max: 100,
            show: false,

            interval: 20,
            axisLabel: {
              formatter: '{value} %'
            },
            axisLine: {
              show: false
            },
            axisTick: {
              show: false
            }
          }],
          series: [
            { type: 'line', symbol: 'none', name: '请求报文量', data: (this.lineData as MessageData).request, yAxisIndex: 0 },
            { type: 'line', symbol: 'none', name: '响应报文量', data: (this.lineData as MessageData).response, yAxisIndex: 0 }
          ]
        }
        break
    }
    this.lineChartOption.title = {
      text: this.titles[0],
      left: 'center',
      top: '0',
      textStyle: {
        fontSize: 16,
        fontWeight: 500
      }
    }
  }
}

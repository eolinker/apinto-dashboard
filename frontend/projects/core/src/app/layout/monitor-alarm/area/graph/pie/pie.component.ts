/* eslint-disable dot-notation */
import { Component, Input, OnInit, SimpleChanges } from '@angular/core'
import { EChartsOption } from 'echarts'
import { changeNumberUnit } from '../../../types/conf'

@Component({
  selector: 'eo-ng-monitor-pie-graph',
  template: `
    <div echarts [options]="pieChartOption"  style="min-width: 186px;min-height:232px;"></div>

  `,
  styles: [
  ]
})
export class MonitorPieGraphComponent implements OnInit {
  @Input() title:string = ''
  @Input() pieData:{[key:string]:any} = {}
  @Input() labelName:string = ''
  @Input() labelValue:string = '0%'
  @Input() subtext:string = ''
  @Input() subvalue:string = ''
  pieChartOption: EChartsOption = {};
  echartsIntance: any

  ngOnInit (): void {
    this.changePieChart()
  }

  ngOnChanges (changes:SimpleChanges):void {
    if (changes['pieData'] || changes['labelValue']) {
      this.changePieChart()
    }
  }

  ngAfterViewInit () {
  }

  chartInit (event:any) {
    this.echartsIntance = event
  }

  changePieChart () {
    this.pieChartOption = {
      tooltip: {
        trigger: 'item'
      },
      title: [
        {
          right: '10',
          subtext: `{title|${this.subtext}}{percent|${this.subvalue}}`,
          top: '15%',
          subtextStyle: {
            rich: {
              title: { fontSize: 14, color: '#666666', lineHeight: 22, padding: [8, 0, 8, 0] },
              percent: { fontSize: 14, color: '#666666', width: 60, lineHeight: 22, align: 'right', padding: [8, 0, 8, 8] }
            },
            fontSize: 14,
            color: '#666666',
            lineHeight: 22,
            padding: [8, 0, 8, 0]
          }
        }],
      legend: [
        {
          top: 'center',
          right: '10',
          orient: 'vertical',
          formatter: (name) => {
            let value = this.pieData[name]
            value = changeNumberUnit(value)
            return `{title|${name}}{percent|${value || '0'}}`
          },
          textStyle: {
            rich: {
              title: { fontSize: 14, color: '#666666', lineHeight: 22, padding: [8, 0, 8, 0] },
              percent: { fontSize: 14, color: '#666666', width: 60, lineHeight: 22, align: 'right', padding: [8, 0, 8, 8] }
            }
          }
        }],
      series: [
        {
          center: ['25%', '50%'],
          name: this.title,
          type: 'pie',
          color: ['#1890FF', '#13c2c2'],
          radius: ['45%', '75%'],
          avoidLabelOverlap: false,
          label: {
            show: true,
            position: 'center',
            formatter: '{text|' + this.labelName + '}\n{value|' + this.labelValue + '}',
            rich: {
              text: { fontSize: 14, color: '#666666', lineHeight: 22, padding: [0, 0, 6, 0] },
              value: { fontSize: 20, color: '#333333', lineHeight: 32, padding: [0, 0, 6, 0] }
            }
          },
          labelLine: {
            show: false
          },
          data: this.transferData(this.pieData)
        }
      ]
    }
    this.echartsIntance?.setOption(this.pieChartOption)
  }

  transferData (value:{[key:string]:any}):Array<{name:string, value:number}> {
    const res:Array<{name:string, value:number}> = []
    const keys = Object.keys(value)
    for (const item of keys) {
      res.push({ name: item, value: value[item] })
    }
    return res
  }
}

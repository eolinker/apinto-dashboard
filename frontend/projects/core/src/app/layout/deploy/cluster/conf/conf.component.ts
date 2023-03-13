/* eslint-disable dot-notation */
/* eslint-disable no-useless-constructor */
/* eslint-disable camelcase */
import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { ActivatedRoute, Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { NzDrawerRef } from 'ng-zorro-antd/drawer'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'

interface RedisData{
  addrs:string,
  username:string,
  password:string,
  enable:boolean|null,
  [key:string]:any
}

interface InfluxdbData{
  addr:string,
  org:string,
  token:string,
  enable:boolean|null,
  [key:string]:any
}
@Component({
  selector: 'eo-ng-deploy-cluster-conf',
  templateUrl: './conf.component.html',
  styles: [
    `
    .intro-box .title-box .title{
      font-size:14px;
      font-weight:500;
      line-height:26px;
      color:var(--MAIN_TEXT);
    }

    .deploy-cluster-redis-box{
      border:1px solid var(--BORDER);
    }
    .redis-box-expand,
    .influxdb-box-expand{
      border-bottom:none;
    }

    .title-box{
      p{
        display:flex;
        align-items:center;
      }
    }


    .icon-tishizhongxin{
      font-size:20px;
    }



   `
  ]
})
export class DeployClusterConfComponent implements OnInit {
  @ViewChild('switchRedisTpl', { read: TemplateRef, static: true }) switchRedisTpl: TemplateRef<any> | undefined
  @ViewChild('switchInfluxdbTpl', { read: TemplateRef, static: true }) switchInfluxdbTpl: TemplateRef<any> | undefined
  @ViewChild('redisTestBtnTpl', { read: TemplateRef, static: true }) redisTestBtnTpl: TemplateRef<any> | undefined
  @ViewChild('influxdbTestBtnTpl', { read: TemplateRef, static: true }) influxdbTestBtnTpl: TemplateRef<any> | undefined
  clusterName:string=''
  redisList:RedisData[]=[{
    username: '', password: '', addrs: '', enable: null
  }
  ]

  influxdbList:InfluxdbData[]=[{
    addr: '', org: '', token: '', enable: null
  }
  ]

  drawerRef:NzDrawerRef | undefined
  drawerEditRef:NzDrawerRef | undefined
  redisExpand:boolean = true
  influxdbExpand:boolean = true
  redisTableHeadName: Array<object> = [
    {
      title: '地址',
      resizeable: true,
      required: true
    },
    {
      title: '用户名',
      resizeable: true
    },
    {
      title: '密码',
      resizeable: true
    },
    {
      title: '启用',
      width: 90,
      resizeable: true
    },
    {
      title: '操作',
      width: 60,
      resizeable: false,
      right: true
    }
  ]

  redisTableBody: Array<any> =[
    {
      key: 'addrs',
      type: 'input',
      placeholder: '请输入域名/IP：端口，多个以逗号分隔',
      checkMode: 'change',
      check: (value:any) => {
        return this.checkAddr(value, 'redisAddr')
      },
      errorTip: '请输入域名/IP：端口，多个以逗号分隔',
      disabledFn: () => {
        return this.nzDisabled
      }
    },
    {
      key: 'username',
      type: 'input',
      placeholder: '请输入用户名',
      change: (item:any) => {
        this.checkValueForEdit(item, 'redis')
      },
      disabledFn: () => {
        return this.nzDisabled
      }
    },
    {
      key: 'password',
      type: 'input',
      placeholder: '请输入密码',
      change: (item:any) => {
        this.checkValueForEdit(item, 'redis')
      },
      disabledFn: () => {
        return this.nzDisabled
      }
    },
    { key: 'enable' },
    {
      type: 'btn',
      right: true,
      btns: [
        {
          title: '测试',
          click: (item:any) => {
            this.testAndSave(item.data, 'redis')
          },
          disabledFn: (data:any, item:any) => {
            return !item.data.addrs || this.nzDisabled
          }
        }
      ]
    }
  ]

  influxdbTableHeadName: Array<object> = [
    {
      title: '数据源地址',
      resizeable: true,
      required: true
    },
    {
      title: 'Organization',
      resizeable: true,
      required: true
    },
    {
      title: '鉴权token',
      resizeable: true
    },
    {
      title: '启用',
      width: 90,
      resizeable: true
    },
    {
      title: '操作',
      width: 60,
      resizeable: false,
      right: true
    }
  ]

  influxdbTableBody: Array<any> =[
    {
      key: 'addr',
      type: 'input',
      placeholder: '请输入数据源地址',
      checkMode: 'change',
      check: (value:any) => {
        return this.checkAddr(value, 'influxdbAddr')
      },
      errorTip: '请输入数据源地址',
      disabledFn: () => {
        return this.nzDisabled
      },
      change: (item:any) => {
        this.checkValueForEdit(item, 'influxdb')
      }
    },
    {
      key: 'org',
      type: 'input',
      placeholder: '请输入Organization',
      change: (item:any) => {
        this.checkValueForEdit(item, 'influxdb')
      },
      disabledFn: () => {
        return this.nzDisabled
      }
    },
    {
      key: 'token',
      type: 'input',
      placeholder: '请输入鉴权Token',
      change: (item:any) => {
        this.checkValueForEdit(item, 'influxdb')
      },
      disabledFn: () => {
        return this.nzDisabled
      }
    },
    { key: 'enable' },
    {
      type: 'btn',
      right: true,
      btns: [
        {
          title: '测试',
          click: (item:any) => {
            this.testAndSave(item.data, 'influxdbv2')
          },
          disabledFn: (data:any, item:any) => {
            return !item.data.addr || !item.data.org || this.nzDisabled
          }
        }
      ]
    }
  ]

  autoTips: Record<string, Record<string, string>> = defaultAutoTips

  redisErrorTip:string = ''
  influxdbErrorTip:string = ''
  nzDisabled:boolean = false
  showRedisError:boolean = false
  showInfluxdbError:boolean = false
  constructor (
    private api:ApiService,
    private baseInfo:BaseInfoService,
    private message: EoNgFeedbackMessageService,
    private activateInfo:ActivatedRoute,
    private router:Router,
    private appConfigService:AppConfigService
  ) {
    this.appConfigService.reqFlashBreadcrumb([{ title: '网关集群', routerLink: 'deploy/cluster' }, { title: '配置管理' }])
  }

  ngOnInit (): void {
    this.clusterName = this.baseInfo.allParamsInfo.clusterName
    if (!this.clusterName) {
      this.router.navigate(['/'])
    }
    this.getRedisList()
    this.getInfluxdbList()
  }

  ngAfterViewInit ():void {
    this.redisTableBody[3].title = this.switchRedisTpl
    this.influxdbTableBody[3].title = this.switchInfluxdbTpl
    this.redisTableBody[4].btns[0].title = this.redisTestBtnTpl
    this.influxdbTableBody[4].btns[0].title = this.influxdbTestBtnTpl
  }

  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  getRedisList () {
    this.api.get('cluster/' + this.clusterName + '/configuration/redis').subscribe((resp:{code:number, data:{redis?:RedisData}, msg:string}) => {
      if (resp.code === 0) {
        if (resp.data.redis) {
          resp.data.redis['origin_addrs'] = resp.data.redis.addrs
          resp.data.redis['origin_username'] = resp.data.redis.username
          resp.data.redis['origin_password'] = resp.data.redis.password
          this.redisList = [resp.data.redis]
        }
      } else {
        this.message.error(resp.msg || '获取redis配置列表失败！')
      }
    })
  }

  getInfluxdbList () {
    this.api.get('cluster/' + this.clusterName + '/configuration/influxdbv2').subscribe((resp:{code:number, data:{influxdbv2?:InfluxdbData}, msg:string}) => {
      if (resp.code === 0) {
        if (resp.data.influxdbv2) {
          resp.data.influxdbv2!['origin_addr'] = resp.data.influxdbv2!.addr
          resp.data.influxdbv2!['origin_org'] = resp.data.influxdbv2!.org
          resp.data.influxdbv2!['origin_token'] = resp.data.influxdbv2!.token
          this.influxdbList = [resp.data.influxdbv2]
        }
      } else {
        this.message.error(resp.msg || '获取influxdb配置列表失败！')
      }
    })
  }

  stopOrStart (item:any, type:string) {
    this.api.put(`cluster/${this.clusterName}/configuration/${type}/${item.enable ? 'enable' : 'disable'}`).subscribe((resp:{code:number, data:any, msg:string}) => {
      if (resp.code === 0) {
        this.message.success(resp.msg || (item.enable ? '启用' : '禁用') + '成功！', { nzDuration: 1000 })
      } else {
        this.message.error(resp.msg || (item.enable ? '启用' : '禁用') + '失败！')
        item.enable = !item.enable
      }
    })
  }

  testAndSave (item:any, type:string) {
    switch (type) {
      case 'redis': {
        this.redisErrorTip = ''
        for (const index in item.addrs.split(',')) {
          if (!/^((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}:[0-9]+$/.test(item.addrs.split(',')[index])) {
            return
          }
        }
        break
      }
      case 'influxdbv2': {
        this.influxdbErrorTip = ''
        if (!/[a-zA-z]+:\/\/[^\s]/.test(item.addr)) {
          return
        }
        break
      }
    }
    const data = type === 'redis'
      ? { addrs: item.addrs || '', username: item.username || '', password: item.password || '' }
      : { addr: item.addr || '', org: item.org || '', token: item.token || '' }
    this.api.put('cluster/' + this.clusterName + '/configuration/' + type, data).subscribe((resp:{code:number, data:any, msg:string}) => {
      if (resp.code === 0) {
        this.message.success(resp.msg || '配置成功！', { nzDuration: 1000 })
        if (type === 'redis') {
          this.getRedisList()
        } else {
          this.getInfluxdbList()
        }
      } else {
        this.message.error(resp.msg || '配置失败！')
        if (type === 'redis') {
          this.redisErrorTip = resp.msg
        } else {
          this.influxdbErrorTip = resp.msg
        }
      }
    })
  }

  nzCheckAddRow = () => {
    return false
  }

  checkValueForEdit (item:any, type:string):void {
    switch (type) {
      case 'redis':
        if (item.username !== item.origin_username || item.addrs !== item.origin_addrs || item.password !== item.origin_password) {
          item.edit = true
        } else {
          item.edit = false
        }
        break
      case 'influxdb':
        if (item.addr !== item.origin_addr || item.org !== item.origin_org || item.token !== item.origin_token) {
          item.edit = true
        } else {
          item.edit = false
        }
        break
    }
  }

  checkAddr (value:any, type:string):boolean {
    switch (type) {
      case 'redisAddr': {
        if (value) {
          for (const index in value.split(',')) {
            if (!/^((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3}:[0-9]+$/.test(value.split(',')[index])) {
              this.showRedisError = true
              return false
            }
          }
          this.showRedisError = false
          return true
        }
        this.showRedisError = true
        return false
      }
      case 'influxdbAddr': {
        if (value) {
          if (!/[a-zA-z]+:\/\/[^\s]/.test(value)) {
            this.showInfluxdbError = true
            return false
          }
          this.showInfluxdbError = false
          return true
        }
        this.showInfluxdbError = true
        return false
      }
    }
    return false
  }
}

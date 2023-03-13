/* eslint-disable dot-notation */
import { Component, Input, OnInit } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { AppConfigService } from '../../../service/app-config.service'
import { EoNgMonitorTabsService } from '../../../service/eo-ng-monitor-tabs.service'
import { ApiService } from '../../../service/api.service'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { ActivatedRoute, Router } from '@angular/router'
import { Subscription } from 'rxjs'
import { MODAL_SMALL_SIZE } from '../../../constant/app.config'
import { defaultAutoTips } from '../../../constant/conf'
import { BaseInfoService } from '../../../service/base-info.service'
export interface ConfigData{
  addr:string
  org:string
  token:string
}
@Component({
  selector: 'eo-ng-monitor-alarm-config',
  templateUrl: './config.component.html',
  styles: [
  ]
})
export class MonitorAlarmConfigComponent implements OnInit {
  @Input() editPage: boolean = false
  validateForm: FormGroup = new FormGroup({})
  nzDisabled: boolean = false
  @Input() partitionId:string = ''
  autoTips: Record<string, Record<string, string>> = defaultAutoTips

  serviceTypeList:Array<{label:string, value:any}> =
    [{ label: 'influxdb', value: 'influxV2' }]

  envList:Array<any> = []
  clustersList:Array<any> = []
  saveTab:boolean = false
  cancelTab:boolean = false
  private subscription: Subscription = new Subscription()

  constructor (private fb: UntypedFormBuilder,
    private appConfigService: AppConfigService,
    private api: ApiService,
    private router: Router,
    public tabs:EoNgMonitorTabsService,
    private modalService:EoNgFeedbackModalService,
    private message: EoNgFeedbackMessageService,
    private baseInfo:BaseInfoService,
    private activateInfo:ActivatedRoute) {
    this.validateForm = this.fb.group({
      name: ['', [this.nameValidator]],
      sourceType: ['influxV2', [Validators.required]],
      addr: ['', [Validators.required, Validators.pattern('[a-zA-z]+://[a-zA-Z0-9.:/]*')]],
      org: [null, [Validators.required]],
      token: [null],
      env: [null, [Validators.required]],
      clusterNames: [null, [Validators.required]]
    })

    this.appConfigService.reqFlashBreadcrumb([
      { title: '监控告警' }
    ])
  }

  ngOnInit (): void {
    this.partitionId = this.baseInfo.allParamsInfo.partitionId
    if (this.editPage) {
      this.getConfig()
    } else {
      this.getClustersList()
    }
  }

  ngOnDestroy () {
    if (!this.saveTab && !this.cancelTab && this.tabs.addTabFlag) {
      this.tabs.addTabFlag = false
      // this.deleteUnsavedTab()
    }
    this.subscription.unsubscribe()
  }

  deleteUnsavedTab () {
    if (!this.editPage && this.tabs.list.length > 1) {
      this.tabs.list.pop()
      this.tabs.index = this.tabs.prevIndex
      this.tabs.prevIndex = 0
    }
  }

  getConfig () {
    this.api.get('monitor/partition', { uuid: this.partitionId }).subscribe((resp:{code:number, data:{info:{name:string, sourceType:string, config:ConfigData, env:string, clusterNames:Array<string>}}, msg:string}) => {
      if (resp.code === 0) {
        this.validateForm.controls['name'].setValue(resp.data.info.name)
        this.validateForm.controls['sourceType'].setValue(resp.data.info.sourceType)
        this.validateForm.controls['addr'].setValue(resp.data.info.config.addr)
        this.validateForm.controls['org'].setValue(resp.data.info.config.org)
        this.validateForm.controls['token'].setValue(resp.data.info.config.token)
        this.validateForm.controls['env'].setValue(resp.data.info.env)
        this.getClustersList(resp.data.info.clusterNames)
        this.validateForm.controls['clusterNames'].setValue(resp.data.info.clusterNames)
      } else {
        this.message.error(resp.msg || '获取分区配置信息失败，请重试!')
      }
    })
  }

  getClustersList (clustersName?:Array<string>) {
    this.api.get('cluster/enum').subscribe((resp:any) => {
      if (resp.code === 0) {
        this.envList = []
        for (const env of resp.data.envs) {
          this.envList.push({ label: env.name, value: env.name, clusters: env.clusters })
        }
        if (this.validateForm.controls['env'].value !== '') {
          this.changeClustersList(this.validateForm.controls['env'].value, clustersName)
        }
      } else {
        this.message.error(resp.msg || '获取集群列表失败，请重试！')
      }
    })
  }

  // 选择不同环境，对应的集群列表随之更新
  changeClustersList (val:string, clustersName?:Array<string>) {
    this.clustersList = []
    this.validateForm.controls['clusterNames'].setValue(null)
    for (const env of this.envList) {
      if (env.value === val) {
        this.clustersList = env.clusters.map((cluster:{name:string}) => {
          return { label: cluster.name, value: cluster.name }
        })
        break
      }
    }
    if (clustersName && clustersName.length > 0) {
      this.validateForm.controls['clusterNames'].setValue(clustersName)
    }
  }

  testLink () {
    if (this.validateForm.valid) {
      const data:{name:string, sourceType:string, config:{addr:string, org:string, token:string}, env:string, clusterNames:Array<string>} =
      {
        name: this.validateForm.controls['name'].value,
        sourceType: this.validateForm.controls['sourceType'].value,
        config: {
          addr: this.validateForm.controls['addr'].value,
          org: this.validateForm.controls['org'].value,
          token: this.validateForm.controls['token'].value || ''
        },
        env: this.validateForm.controls['env'].value,
        clusterNames: this.validateForm.controls['clusterNames'].value
      }

      if (!this.editPage) {
        this.api.post('monitor/partition', data).subscribe((resp:{code:number, data:{info:{uuid:string, name:string, clusterNames:Array<string>}}, msg:string}) => {
          if (resp.code === 0) {
            this.saveTab = true
            this.tabs.getTabsData(resp.data.info.uuid)
          } else {
            this.message.error(resp.msg || '连接失败，请检查输入后重试！')
          }
        })
      } else {
        this.api.put('monitor/partition', data, { uuid: this.partitionId }).subscribe((resp:{code:number, data:{info:{clusterNames:Array<string>}}, msg:string}) => {
          if (resp.code === 0) {
            this.saveTab = true
            this.tabs.getTabsData(this.partitionId)
          } else {
            this.message.error(resp.msg || '连接失败，请检查输入后重试！')
          }
        })
      }
    } else {
      Object.values(this.validateForm.controls).forEach(control => {
        if (control.invalid) {
          control.markAsDirty()
          control.updateValueAndValidity({ onlySelf: true })
        }
      })
    }
  }

  deleteModal () {
    this.modalService.create({
      nzTitle: '删除',
      nzContent: `请确认是否删除${this.tabs.getTabName(this.partitionId)}分区？`,
      nzClosable: true,
      nzClassName: 'delete-modal',
      nzWidth: MODAL_SMALL_SIZE,
      nzOkDanger: true,
      nzOnOk: () => {
        this.deleteArea()
      }
    })
  }

  deleteArea () {
    this.api.delete('monitor/partition', { uuid: this.partitionId }).subscribe((resp:any) => {
      if (resp.code === 0) {
        this.tabs.deleteTab(this.partitionId)
        this.router.navigate(['/', 'monitor-alarm', 'area', 'total', this.tabs.list[this.tabs.index].uuid])
      } else {
        this.message.error(resp.msg || '删除失败，请重试！')
      }
    })
  }

  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  // 返回列表页，当fromList为true时，该页面左侧有分组
  backToList () {
    this.cancelTab = true
    this.tabs.showList = this.tabs.list
    this.deleteUnsavedTab()
    this.router.navigate(['/', 'monitor-alarm', 'area', 'total', this.tabs.list[this.tabs.index].uuid])
  }

  // 分区名称为必填，长度在16个字符以内，不可与已有分区名重复
  nameValidator = (control: any): { [s: string]: boolean } => {
    if (!control.value) {
      return { error: true, required: true }
    } else if (this.getByteLen(control.value) > 16) {
      return { error: true, length: true }
    } else if (this.checkNameOverlap(control.value)) {
      return { error: true, overlap: true }
    }
    return {}
  }

  getByteLen (val:string) {
    let len = 0
    for (let i = 0; i < val.length; i++) {
      const length = val.charCodeAt(i)
      if (length >= 0 && length <= 128) {
        len += 1
      } else {
        len += 2
      }
    }
    return len
  }

  checkNameOverlap (val:string) {
    for (const tab of this.tabs.list) {
      if (tab.title === val && tab.uuid !== this.partitionId) {
        return true
      }
    }
    return false
  }
}

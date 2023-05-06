/* eslint-disable dot-notation */
/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2022-07-20 22:34:58
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2022-09-28 00:20:26
 * @FilePath: /apinto/src/app/layout/deploy/deploy-cluster-list/deploy-cluster-list.component.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import {
  Component,
  ElementRef,
  OnInit,
  TemplateRef,
  ViewChild
} from '@angular/core'
import {
  FormGroup
} from '@angular/forms'
import { Router } from '@angular/router'
import {
  EoNgFeedbackModalService,
  EoNgFeedbackMessageService
} from 'eo-ng-feedback'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { EoNgNavigationService } from 'projects/core/src/app/service/eo-ng-navigation.service'
import { DeployService } from '../../deploy.service'
import { ClustersThead } from '../types/conf'
import { MODAL_SMALL_SIZE } from 'projects/core/src/app/constant/app.config'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { DeployClusterCreateComponent } from '../create/create.component'
import { ClusterEnum } from 'projects/core/src/app/constant/type'
import { SelectOption } from 'eo-ng-select'

@Component({
  selector: 'eo-ng-deploy-cluster-list',
  templateUrl: './list.component.html',
  styles: []
})
export class DeployClusterListComponent implements OnInit {
  @ViewChild('clusterStatusTpl', { read: TemplateRef, static: true })
  clusterStatusTpl: TemplateRef<any> | undefined

  readonly nowUrl: string = this.router.routerState.snapshot.url
  nzDisabled:boolean = false
  validateForm: FormGroup = new FormGroup({})
  source: string = '' // 集群地址通过测试后得到的source, 有source的情况才能新建集群成功
  @ViewChild('clusterNameInput', { static: false }) clusterNameInput:
    | ElementRef
    | undefined

  clustersList: Array<object> = []
  clustersDisplayList: Array<object> = []
  clustersTableHeadName:THEAD_TYPE[] = [...ClustersThead]
  clustersTableBody: TBODY_TYPE[] = []

  environmentList: SelectOption[]= []

  autoTips: Record<string, Record<string, string>> = defaultAutoTips

  nodesTableShow = false
  clusterCanBeCreated: boolean = false
  modalRef:NzModalRef|undefined

  env:Array<string> = []

  // eslint-disable-next-line no-useless-constructor
  constructor (
    private message: EoNgFeedbackMessageService,
    private modalService: EoNgFeedbackModalService,
    private api: ApiService,
    public router: Router,
    private navigationService: EoNgNavigationService,
    private service:DeployService
  ) {
    this.navigationService.reqFlashBreadcrumb([{ title: '网关集群', routerLink: 'deploy/cluster' }])
  }

  ngOnInit (): void {
    this.getClustersData()
    this.getEnvList()
  }

  ngAfterViewInit () {
    this.clustersTableBody = [...this.service.createClusterTbody(this)]
  }

  getClustersData () {
    this.api.get('clusters').subscribe((resp) => {
      if (resp.code === 0) {
        this.clustersList = resp.data.clusters
        this.clustersDisplayList = [...this.clustersList]
      }
    })
  }

  getEnvList () {
    this.api.get('cluster/enum').subscribe((resp: {code:number, data:{ envs:ClusterEnum[]}, msg:string}) => {
      if (resp.code === 0) {
        this.environmentList = resp.data.envs.map((env:ClusterEnum) => {
          return { label: env.name, value: env.name }
        })
      }
    })
  }

  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  delete (item:any) {
    this.modalService.create({
      nzTitle: '删除',
      nzContent: '该数据删除后将无法找回，请确认是否删除？',
      nzClosable: true,
      nzClassName: 'delete-modal',
      nzWidth: MODAL_SMALL_SIZE,
      nzOkDanger: true,
      nzOnOk: () => {
        this.deleteCluster(item)
      }
    })
  }

  deleteCluster (item: any) {
    this.api
      .delete('cluster', { clusterName: item.name })
      .subscribe((resp) => {
        if (resp.code === 0) {
          this.message.success(resp.msg || '删除成功', { nzDuration: 1000 })
          this.getClustersData()
        }
      })
  }

  clusterTableClick= (item:any) => {
    this.router.navigate(['/', 'deploy', 'cluster', 'content', item.data.name])
    this.service.clusterName = item.data.title
    this.service.clusterDesc = item.data.desc
  }

  addCluster () {
    this.modalRef = this.modalService.create({
      nzTitle: '新建集群',
      nzWidth: MODAL_SMALL_SIZE,
      nzContent: DeployClusterCreateComponent,
      nzComponentParams: {
        closeModal: () => { this.modalRef?.close() }
      },
      nzOnOk: (component:DeployClusterCreateComponent) => {
        console.log(component)
      },
      nzFooter: [{
        label: '取消',
        type: 'default',
        onClick: () => {
          this.modalRef?.close()
        }
      },
      {
        label: '下一步，检查集群',
        type: 'primary',
        onClick: (context:DeployClusterCreateComponent) => {
          context.testCluster()
        },
        disabled: () => {
          return this.nzDisabled
        },
        show: (context:any) => {
          return !context.nodesTableShow && !context.checkClusterError
        },
        loading: (context:any) => {
          return context.testButtonLoading
        }
      },
      {
        label: '重新检查',
        type: 'primary',
        onClick: (context:DeployClusterCreateComponent) => {
          context.testCluster()
        },
        disabled: () => {
          return this.nzDisabled
        },
        show: (context:any) => {
          return !context.nodesTableShow && context.checkClusterError
        },
        loading: (context:any) => {
          return context.testButtonLoading
        }
      },
      {
        label: '上一步',
        type: 'default',
        onClick: (context:DeployClusterCreateComponent) => {
          context.cancel()
        },
        show: (context:any) => {
          return context.nodesTableShow
        }
      },
      {
        label: '完成',
        type: 'primary',
        onClick: (context:DeployClusterCreateComponent) => {
          context.saveCluster()
        },
        disabled: () => {
          return this.nzDisabled
        },
        show: (context:any) => {
          return context.nodesTableShow
        },
        loading: (context:any) => {
          return context.submitButtonLoading
        }
      }]

    })
  }

  handleEnvChange () {
    this.clustersDisplayList = this.clustersList.filter((cluster:any) => {
      return this.env.indexOf(cluster.env) !== -1
    })
  }
}

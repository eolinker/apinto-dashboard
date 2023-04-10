/* eslint-disable dot-notation */
import { Component } from '@angular/core'
import { FormArray, FormControl, FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { NzUploadFile } from 'ng-zorro-antd/upload'
import { defaultAutoTips } from '../../../constant/conf'
import { EmptyHttpResponse } from '../../../constant/type'
import { ApiService } from '../../../service/api.service'
import { EoNgMessageService } from '../../../service/eo-ng-message.service'
import { CdkDragDrop, moveItemInArray } from '@angular/cdk/drag-drop'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { NavigationData } from '../types/types'

const getBase64 = (file: File): Promise<string | ArrayBuffer | null> =>
  new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.readAsDataURL(file)
    reader.onload = () => resolve(reader.result)
    reader.onerror = error => reject(error)
  })

@Component({
  selector: 'eo-ng-navigation-create',
  templateUrl: './create.component.html',
  styles: [
    `:host ::ng-deep{
      .ant-upload.ant-upload-select-picture-card{
        height:32px;
        width:32px;

        [nz-upload-btn]{
          display:block;
        }
      }
    }`
  ]
})
export class NavigationCreateComponent {
  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  validateForm: FormGroup = new FormGroup({ })

  fileList: NzUploadFile[] = [];
  editPage:boolean = false
  navigationUuid:string = ''
  listOfControl: Array<{ id: number; controlInstance: string }> = [];
  imageUrl:string = ''

  modalRef:NzModalRef|undefined
  get childArray () {
    return this.validateForm.get('children') as FormArray
  }

  constructor (
    private fb: UntypedFormBuilder, private api: ApiService, private message:EoNgMessageService
  ) {
  }

  // 手动上传文件
  beforeUpload = (file: NzUploadFile): boolean => {
    this.fileList = []
    this.fileList = this.fileList.concat(file)
    this.handlerFileChange()
    return false
  }

  // 移除文件
  removeFile () {
    this.fileList = []
    return true
  }

  ngOnInit (): void {
    this.validateForm = this.fb.group({
      name: ['', [Validators.required]],
      children: this.fb.array([])
    })
    if (this.editPage) {
      this.getNavMessage()
    }
  }

  dataURLtoFile (dataUrl:string, fileName:string = '') {
    const arr = dataUrl.split(',')
    const mime = arr[0].match(/:(.*?);/)![1]
    const bstr = atob(arr[1]); let n = bstr.length
    const u8arr = new Uint8Array(n)
    while (n--) {
      u8arr[n] = bstr.charCodeAt(n)
    }
    return new File([u8arr], fileName, { type: mime })
  }

  getNavMessage () {
    this.api.get(`system/navigation/${this.navigationUuid}`).subscribe((resp:{code:number, msg:string, data:{navigation:NavigationData}}) => {
      if (resp.code === 0) {
        this.validateForm.patchValue(
          {
            name: resp.data.navigation.title
          }
        )
        this.fileList = [this.dataURLtoFile(resp.data.navigation.icon) as any]
        this.handlerFileChange()
        const childArray = this.validateForm.controls['children'] as FormArray
        resp.data.navigation.modules.forEach((module:{id:string, title:string}) => {
          const childGroup = new FormGroup({})
          childGroup.addControl('id', new FormControl(module.id))
          childGroup.addControl('title', new FormControl(module.title))
          childArray.push(childGroup)
        })
      }
    })
  }

  drop (event: CdkDragDrop<string[]>) {
    const childList = this.validateForm.controls['children'].value
    moveItemInArray(childList, event.previousIndex, event.currentIndex)
    this.validateForm.controls['children'].setValue(childList)
  }

  async submit () {
    if (this.validateForm.controls['name'].invalid) {
      this.validateForm.controls['name'].markAsDirty()
      this.validateForm.controls['name'].updateValueAndValidity({ onlySelf: true })
    }

    if (this.editPage && this.validateForm.controls['children'].invalid) {
      this.validateForm.controls['children'].markAsDirty()
      this.validateForm.controls['children'].updateValueAndValidity({ onlySelf: true })
      return
    }

    if (this.validateForm.controls['name'].invalid) {
      return
    }
    const iconBase64 = await getBase64(this.fileList[0] as any)
    if (this.editPage) {
      this.api.put(`system/navigation/${this.navigationUuid}`, { name: this.validateForm.controls['name'].value, icon: iconBase64, modules: this.validateForm.controls['children'].value }).subscribe((resp:EmptyHttpResponse) => {
        if (resp.code === 0) {
          this.message.success(resp.msg || '编辑导航成功')
          this.modalRef?.close()
        }
      })
    } else {
      this.api.post('system/navigation', { uuid: this.navigationUuid, name: this.validateForm.controls['name'].value, icon: iconBase64 }).subscribe((resp:EmptyHttpResponse) => {
        if (resp.code === 0) {
          this.message.success(resp.msg || '新建导航成功')
          this.modalRef?.close()
        }
      })
    }
  }

  async handlerFileChange () {
    this.imageUrl = await getBase64(this.fileList[0] as any) as string
  }
}

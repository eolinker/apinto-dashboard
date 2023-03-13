/* eslint-disable camelcase */
/**
 * 处理 Api Tabs状态
 */

// import { createAction, props } from '@ngrx/store'
export interface APITabParams {
    hashkey?: string
    groupId?: string
    timestamp?: string
    apiId?: string
    /** 组件类型 */
    type?: string
    // 实际
    partition_id?:string // 分区的uuid
  }

export interface APITabType {
   /** 标题名称 */
   name: string
   /** 路由 */
   path: string[]
   /** Params参数 */
   params: APITabParams
   /** 渲染组件 */
   component: unknown
   /** 避免重复加载, loadComponent之后根据该标识不再重复加载 */
   loaded?: boolean
   /** 当前 API 方法 Tab上展示用 */
  //  method?: APIMethodsType
   /** 预留设计 - 该值注入组件内部 */
   payload?: unknown
   /** 修改状态 */
   modified?: boolean
   /** 组件唯一标识 */
   componentId?: string
 }

// /** 索引值修改 */
// export const currentIndexChange = createAction('[API Tabs] Current index changed', props<{ index: number }>())

// /** 编辑状态修改 */
// export const currentEditableChange = createAction(
//   '[API Tabs] Current tab editable status changed',
//   props<{ index: number; editable: boolean }>()
// )

// /** 新增标签页 */
// export const addTab = createAction('[API Tabs] add tab', props<{ tab: APITabType }>())

// /** 标记动态组件已缓存 */
// export const tabLoaded = createAction('[API Tabs] tab has loaded', props<{ index: number }>())

// /** 删除标签页 */
// export const delTabByIndex = createAction('[API Tabs] del tab', props<{ index: number }>())

// /** 批量删除标签页 */
// export const batchDelTabByIndex = createAction('[API Tabs] batch del tab', props<{ ids: number[] }>())

// /** 清空标签页列表 */
// export const clearTabList = createAction('[API Tabs] clear tab list')

// /** 更改Tab的路由 */
// export const resetTabRoute = createAction('[API Tabs] reset tab route', props<{ index?: number; path?: string[] }>())

// /** 更改Tab的方法 */
// export const resetTabMethod = createAction(
//   '[API Tabs] reset tab method',
//   props<{ index?: number; method?: APIMethodsType }>()
// )

// /** 更改Tab的标题 */
// export const resetTabTitle = createAction('[API Tabs] reset tab title', props<{ index?: number; title: string }>())

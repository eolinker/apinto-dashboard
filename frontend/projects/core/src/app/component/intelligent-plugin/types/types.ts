
export type DynamicField = {
    name: string,
    title: string,
    attr: string,
    enum: Array<string>
}

export type DynamicDriverData = {
    name:string, title:string
}

export type DynamicConfig = {
    id:string,
    name: string,
    title: string,
    drivers: Array<DynamicDriverData>,
    fields: Array<DynamicField>,
    list: Array<{[k:string]:any}>
}

export type DynamicRender = {
    render:any,
    id:string,
    name:string,
    title:string
}

export type DynamicListStatus = {
    [k:string]:{
        [id:string]:string
    }
}

export type DynamicPublish = {
    code:number,
    msg:string,
    data:{
        success:Array<string>,
        fail:Array<string>
    }
}

export type DynamicPublishCluster = {
    name:string,
    title:string,
    status:string,
    updater:string,
    update_time:string,
    checked?:boolean
}

export type DynamicPublishData = {
    id:string,
    name:string,
    title:string,
    description:string
    clusters:DynamicPublishCluster[]
}

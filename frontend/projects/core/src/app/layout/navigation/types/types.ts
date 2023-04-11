export type NavigationItem = {
    uuid:string
    title:string
    icon:string
    canDelete:boolean
}

export type NavigationData = {
    uuid:string
    title:string
    icon:string
    modules:Array<{
        id:string
        title:string
    }>
}

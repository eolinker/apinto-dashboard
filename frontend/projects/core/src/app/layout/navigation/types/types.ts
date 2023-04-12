export type NavigationItem = {
    uuid:string
    title:string
    icon:string
    canDelete:boolean
    iconType:string
}

export type NavigationData = {
    uuid:string
    title:string
    icon:string
    iconType:string
    modules:Array<{
        id:string
        title:string
    }>
}

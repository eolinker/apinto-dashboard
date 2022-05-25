

    const RootId = "FormRender"
    function CheckBySchema(id, schema,value){
        let validator = validate.djv()
        if (!validator.resolved.hasOwnProperty(id)) {
            validator.addSchema(id, schema);
        }
        let err = validator.validate(id, value)
        if(err){
            console.log(err)
            return false
        }
        return true
    }
    function ValidHandler(schema){
        let value =$(this).val()
        if (schema["type"] === "integer" || schema["type"] === "number"){
            value = Number(value)
        }
        let rs = CheckBySchema(this.id, schema, value)
        if( rs === true) {
            $(this).removeClass("is-invalid")
            $(this).addClass("is-valid")
            return true
        }else{
            $(this).removeClass("is-valid")
            $(this).addClass("is-invalid")
            return false
        }
    }
    function InputValid(schema,Id){
        $("#"+Id).on("change",function (){
            ValidHandler.apply(this,[schema])
        })
    }
    class BaseEnumRender {
        constructor(panel,schema,path) {

            const Id = path
            this.Id = Id

            $(panel).append(createEnum(schema,Id))
            // InputValid(schema,Id)
        }
        get Value(){
            return $("#"+this.Id).val()
        }
        set Value(v){
            $("#"+this.Id).val(v)
        }
    }
    class SwitchRender {
        constructor(panel,schema,path) {

            const Id = path
            this.Id = Id

            $(panel).append('<input id="switch_'+Id+'}" type="checkbox" data-toggle="toggle" >')

            InputValid(schema,Id)

        }
        get Value(){
            return $("#switch_"+this.Id).checked
        }
        set Value(v){
            if (v === true || v === "true"){
                $("#switch_"+this.Id).attr('checked', 'checked');
            }else{
                $("#switch_"+this.Id).removeAttr('checked');
            }

        }
    }
    class BaseInputRender {
        constructor(panel,schema,path) {

            const Id = path
            this.Id = Id

            $(panel).append(createInput(schema,Id))

            InputValid(schema,Id)

        }
        get Value(){
            return $("#"+this.Id).val()
        }
        set Value(v){
            $("#"+this.Id).val(v)
        }
        check(){

        }
    }
    function getLabel(schema){
        let label = schema["label"]
        if (!label || label.trim() === ""){
            label = schema["name"]
        }
        return label
    }
    function createLabel(schema,id,appendAttr){
        if (!appendAttr){
            appendAttr =""
        }

        let labelHtml = '<label class="col-sm-2 col-form-label text-right"'
        if (id){
            labelHtml +=' for="'+id+'"'
        }
        if (appendAttr){
            labelHtml += ' '+appendAttr
        }
        labelHtml += '>'
        if (schema["required"]){
            labelHtml += '<span style="color: red">*</span>'
        }
        labelHtml += getLabel(schema)+'</label>'
        return labelHtml
    }

    function createFieldPanel(panel,schema,id,appendAttr){
        $(panel).append('<div class="form-group row mb-3""></div>')
        let fieldPanel = $(panel).children().last()
        fieldPanel.append(createLabel(schema,id,appendAttr))
        fieldPanel.append('<div class="col-sm-10"></div>')
        return fieldPanel.children().last()
    }
    function createEnum(schema,id,appendAttr){
        let readOnly = ""
        this.schema = schema
        if (schema["readonly"] === true){
            readOnly = "readonly"
        }

        if (schema["enum"]){
            let enums = schema["enum"]

            let select = '<select '+readOnly+' '+appendAttr+' class="form-control" id="'+id+'"'
            if (schema["required"]){
                select += " required"
            }
            select += ">"

            for (let i in enums){
                if (typeof value != "undefined" &&  value === enums[i]){
                    select += '<option>'+enums[i]+'</option>'
                }else{
                    select += '<option selected>'+enums[i]+'</option>'
                }
            }
            select += '</select>'

            return select
        }
    }
    function createInput(schema,id,appendAttr){
        let readOnly = ""
        this.schema = schema
        if (schema["readonly"] === true){
            readOnly = "readonly"
        }

        let input = '<input '+readOnly +' class="form-control" id="'+id+'" aria-describedby="validation_'+id+'" ';
        if (appendAttr){
            input += appendAttr
        }

        function readFormatForString(format){

            switch (format){

                case "email","password","date","time","number":{
                    return format
                }
                case "idn-email":{
                    return "email"
                }
                case "date-time":{
                    return 'datetime-local'
                }
                case "boolean":{
                    return "checkbox"
                }
                default:{
                    return "text"
                }
            }
        }
        switch (schema["type"]){
            case "string":{
                input += ' type="'+readFormatForString(schema["format"])+'"'
                if (schema["format"] === 'password'){
                    input += ' autocomplete="on"'
                }
                break
            }
            case "integer","number":{
                input += ' type="number"'
                break
            }
        }
        if (schema["maxLength"]){
            input += ' maxlength="'+schema["maxLength"]+'"'
        }

        if (schema["minLength"]){
            input += ' minLength="'+schema["minLength"]+'"'
        }
        if (typeof value != "undefined"){
            input += 'value="'+value+'"'
        }
        if(schema["required"]){
            input += ' required'
        }
        input+='/>'
        return input
    }
    class MapRender{
        constructor(panel,schema,generator,path) {
            this.Schema = schema
            this.Id = path;
            $(panel).append("<div id='"+path+"_panel' class='container-fluid'></div>")
            this.PanelId = "#"+path+"_panel"
            this.GeneratorHandler = generator

        }
        set Value(v){

            const Items = this.Schema["items"]
            const Id = this.Id
            const keySchema = {type:"string"}

            $(this.PanelId).empty()
            switch (Items["type"]){
                case "object","map":{
                    break
                }
                default:{
                    $(this.PanelId).append('<div class="input-group input-group-sm m-2">' +
                        '<div class="input-group-prepend">' +
                        '<div class="input-group-text  btn" id="btnGroupAddon_'+Id+'_new">+</div>' +
                        '</div>'+
                        createInput(keySchema,Id+'_key','aria-describedby="btnGroupAddon_'+Id+'_new" placeholder="Input new key" ')+
                        '<div class="input-group-prepend ">' +
                        '<div class="input-group-text  btn" id="btnGroupAddon_'+Id+'_eq">=</div>' +
                        '</div>'+
                        createInput(Items,Id+'_value','aria-describedby="btnGroupAddon_'+Id+'_eq" placeholder="Input new value" ')+

                        '</div>')

                for(let k in v){
                    $(this.PanelId).append('<div class="input-group input-group-sm m-2"></div>')
                    let itemPanel =  $(this.PanelId).children().last()
                    itemPanel.append(
                        '<div class="input-group-prepend">' +
                        '<button class="btn btn-danger" id="btnGroupAddon_'+Id+'_key_'+k+'" type="button" data-itemId="'+Id+'_key"> - </button>\n' +
                        '</div>')
                    let childItemKey = new BaseInputRender(itemPanel,keySchema,Id+"_key_"+k)
                    childItemKey.Value = k
                    itemPanel.append('<div class="input-group-prepend"><div class="input-group-text  btn" >=</div></div>')
                    let childItemValue = new BaseInputRender(itemPanel,Items,Id+"_value_"+k)
                    childItemValue.Value = v[k]
                }
                break
                }
            }

        }
        get Value(){
            return {}
        }
        check(){

        }
    }
    class InterObjRender{
        constructor(panel,schema,generator,path) {

            const Id = path;
            this.Id = Id;
            this.Schema = schema
            const items = schema["items"]
            console.assert(items["type"]==="object")
            const p = $(panel)
            this.panelId = Id+'_items'
            p.append('<div id="'+Id+'_toolbar">\n' +
                '  <button id="'+Id+'_AddButton" class="btn btn-secondary">Add</button>\n' +
                '</div>')
            $("#"+Id+"_AddButton").on("click",function (event){
                Table.bootstrapTable('append', {})
                Table.bootstrapTable('scrollTo', 'bottom')
                return false
            })
            p.append('<table  id="'+this.panelId +'"></table>')
            const Table = $("#"+this.panelId)
            Table.delegate("a.remove","click",function (event){
                let rowIndex = $(this).attr("array-row")
                Table.bootstrapTable('remove', {
                    field: '$index',
                    values: [Number(rowIndex)]
                })
            })
            this.Table = Table
            const properties = items["properties"]
            let lastDetailRow = undefined
            let lastField = undefined
            function DetailFormatterHandler(fieldIndex){

                this.detailFormatter = function (index, row, $element){

                    if(typeof lastDetailRow !== "undefined" && lastDetailRow !== index ){
                        Table.bootstrapTable('collapseRow', lastDetailRow)
                    }
                    lastDetailRow = index
                    lastField = fieldIndex
                    let item = properties[fieldIndex]
                    let child =  generator($element,item,generator,path+"_"+item["name"])
                    child.Value = row[item.name]
                    return ""
                }
                return this
            }

            function NotDetailFormatterMap(index, row, $element){
                if (typeof lastDetailRow !== "undefined"){
                    Table.bootstrapTable('collapseRow')
                    lastDetailRow = undefined
                    lastField = undefined
                }
                Table.bootstrapTable('collapseAllRows')

                return ''
            }
            function formatterKV(v){
                let html =""
                for (let k in v){
                    html+="<span class='btn btn-outline-secondary btn-sm'>"+k+"="+v[k]+"</span>\n"
                }
                html+= ""
                return html
            }


            const columns =[]
            columns.push({
                title:"",
                field:"__index",
                sortable:false,
                editable:false,
                detailFormatter:NotDetailFormatterMap,
                formatter: function (v,row,index){
                    return index
                }
            })

            for(let i in properties){
                const item = properties[i]
                switch (item["type"]){
                    case "object":{
                        columns.push({
                            title:getLabel(item),
                            field:item["name"],
                            sortable:false,
                            editable:false,
                            formatter: formatterKV,
                            detailFormatter:new DetailFormatterHandler(i).detailFormatter
                        })
                        break
                    }
                    case "map":{
                        columns.push({
                            title:getLabel(item),
                            field:item["name"],
                            sortable:false,
                            editable:false,
                            formatter:formatterKV,
                            detailFormatter:new DetailFormatterHandler(i).detailFormatter
                        })
                        break
                    }
                    case "array":{
                        columns.push({
                            title:getLabel(item),
                            field:item["name"],
                            sortable:false,
                            editable:false,
                            formatter: formatterKV,
                            detailFormatter:new DetailFormatterHandler(i).detailFormatter,
                        })
                        break
                    }
                    default:{
                        if (item["enum"]){
                            columns.push({
                                title:getLabel(item),
                                field:item["name"],
                                sortable:true,
                                detailFormatter:NotDetailFormatterMap,
                                editable: {
                                    type:"select",
                                    options:{
                                        items:item["enum"]
                                    }
                                },
                            })
                        }else{
                            let typeInput = "text"
                            if (item["type"] === "number" || item["type"] === "integer"){
                                typeInput = "number"
                            }
                            columns.push({

                                title:getLabel(item),
                                field:item["name"],
                                sortable:true,
                                width:200,
                                detailFormatter:NotDetailFormatterMap,
                                editable:{
                                    type:typeInput
                                }
                            })
                        }
                        break
                    }
                }
            }
            columns.push({
                title:"操作",
                field:"",
                sortable:false,
                editable:false,
                formatter: function (v,row,index){
                    return "<a class=\"remove\" href=\"javascript:void(0)\" array-row=\""+index+"\" title=\"remove\">删除</a>"
                }
            })
            const _This = this
            const tableOptions = {
                columns:columns,
                editable: true,
                // toolbar:'#'+Id+'toolbar',

                detailView:true,
                detailViewByClick:true,
                detailViewIcon: false,
                width: "100%",
                onEditorShown:function (field, row, $el, editable){
                    Table.bootstrapTable('collapseAllRows')

                    return true;
                },
                onEditorSave:function (field, row, oldValue, $el){
                    if (field !=="__index", _This.Data){
                        const rowIndex = $el.parent().data("index")
                        _This.Data[rowIndex][field] = row[field]
                    }
                    return true;
                },
            }
            Table.bootstrapTable(tableOptions);
        }
        set Value(v){
            if (!v){
                v = []
            }
            this.Table.bootstrapTable("load",v)
            this.Data = v
        }
        get Value(){
            return  this.Data
        }
    }
    class InterMapRender{
        constructor(panel,schema,generator,path) {}
        set Value(v){

        }
        get Value(){
            return {}
        }
    }
    class ArrayRenderEnum {
        constructor(panel,schema,generator,path) {
            const Id = path;
            this.Id = Id;
            const items = schema["items"]
            this.Enum = items["enum"]

            let p =$(panel);
            p.append('<div id="'+Id+'_items" class="border p-sm-1 btn-toolbar " role="toolbar"></div>')
            const itemPanel = p.children('#'+Id+'_items')
            for (let i in items["enum"]){
                let e  = items["enum"][i]
                let itemId = Id+'_'+e
                itemPanel.append('<div class="custom-control custom-checkbox custom-control-inline">\n' +
                    '  <input type="checkbox" id="'+itemId+'" value="'+e+'" name="'+Id+'" class="custom-control-input">\n' +
                    '  <label class="custom-control-label" for="'+itemId+'">'+e+'</label>\n' +
                    '</div>')
            }
        }
        get Value(){
            const list = []
            $('input[name="'+this.Id +'"]').each(function (){
                if( $(this).is(':checked')){
                    list.push($(this).val())
                }
            })
            return list
        }
        set Value(vs){
            if (!vs){
                vs = []
            }
            const list = vs

            for (let i in this.Enum){
                let v = this.Enum[i]
                const itemId = this.Id+'_'+v
                $("#"+itemId).attr("checked",list.includes(v))
            }
        }

    }
    class ArrayRenderSimple {
        constructor(panel,schema,generator,path) {
            this.Schema = schema
            const items = schema["items"]
            const Id = path;
            this.Id = Id;
            let p =$(panel);

            p.append('<div id="'+Id+'_items" class="border p-sm-1 btn-toolbar " role="toolbar"></div>')

            const itemPanel = p.children('#'+Id+'_items')
            itemPanel.append('<div class="input-group input-group-sm m-2">' +
                '<div class="input-group-prepend ">' +
                '<div class="input-group-text  btn" id="btnGroupAddon_'+Id+'_new">+</div>' +
                '</div>'+
                createInput(items,Id+'_new','aria-describedby="btnGroupAddon_'+Id+'_new" placeholder="Input new" ')+
                '</div>')

            $('#'+Id+"_new").on("change",function (){
                let v = $(this).val()
                if (items["type"] === "integer" || items["type"] === "number"){
                    v = Number(value)
                }
                if (v !== ""){
                    if (CheckBySchema(Id,items,v)){
                        add(v)
                        $(this).val("")
                    }
                }
            })

            let lastIndex = 0

            function add(value){
                const itemId = Id+"_item_"+(lastIndex++)
                const appendAtt = 'array-for="'+Id+'" aria-describedby="btnGroupAddon_'+itemId+'"'
                itemPanel.append('<div class="input-group input-group-sm m-2" id="array-item_'+itemId+'">\n' +
                    '<div class="input-group-prepend">\n' +
                    '<button class="btn btn-danger" id="btnGroupAddon_'+itemId+'" type="button" aria-describedby="btnGroupAddon_'+itemId+'" data-itemId="'+itemId+'"> - </button>\n' +
                    '</div>\n' +
                    createInput(items,itemId,appendAtt) +
                    '</div>')
                $("#"+itemId).val(value)
                return false
            }
            this.add = add
            itemPanel.delegate('button',"click",function (){
                let itemId =  $(this).attr('data-itemId')
                itemPanel.children('#array-item_'+itemId).remove()
            })

        }
        check(){

            const arrayId = "[array-for='"+this.Id+"']"
            if(schema["minLength"]>0 &&  $(arrayId).length === 0){
                return false
            }
            let rel = true
            const items = this.Schema
            const Id = this.Id
            $(arrayId).each(function (){
               if( CheckBySchema(Id,items,$(this).val()) !== true){
                   rel = false
                   return false
               }
            })

            return rel

        }
        get Value(){
            let val = []
            $("[array-for='"+this.Id+"']").each(function (){
                val.push($(this).val())
            })
            return val
        }
        set Value(vs){
            if (!Array.isArray(vs)){
                return
            }
            const arrayId = "[array-for='"+this.Id+"']"
            let list =  $(arrayId)
            if (list.length < vs.length){
                let num = vs.length - list.length
                for (let i = 0; i < num; i++) {
                    this.add()
                }
            }else if (list.length > vs.length){
                let num = vs.length - list.length
                for (let i = 0; i < num; i++) {
                    list.last().remove()
                }
            }
            let index = 0;
            $(arrayId).each(function (){
                $(this).val(vs[index])
                index ++;
            })
        }
    }
    class ObjectRender {
        constructor(panel,schema,generator,path) {
            this.Fields = {}
            this.Id = path
            if (schema["type"] !== "object"){
                return
            }
            let properties = schema["properties"]
            for (let i in properties){
                let sub =  properties[i]
                let name = sub["name"]
                const subId = path+"_"+name
                const fieldPanel = createFieldPanel(panel,sub,subId)
                this.Fields[name] = generator(fieldPanel,sub,generator,subId)
            }
        }

        get Value(){
            let v = {}
            for (let k in this.Fields){
                let fi = this.Fields[k]
                v[k]=fi.Value
            }
            return v
        }
        set Value(v){
            for (let k in this.Fields){
                let fi = this.Fields[k]
                fi.Value = v[k]
            }
        }
        check(){
            for (let k in this.Fields){
              if(! this.Fields[k].check()){
                  return false
              }
            }
        }
    }


function BaseGenerator(panel,schema,generator,path){
    switch (schema["type"]){
        case "object":{
            return new ObjectRender(panel,schema,generator,path)
        }
        case "array":{
            const items = schema["items"]
            switch (items["type"]){
                case "object":{
                    return new InterObjRender(panel,schema,generator,path)
                }
                case "map":{
                    return new InterMapRender(panel,schema,generator,path)
                }
            }
            if (items["enum"]){
                return new ArrayRenderEnum(panel,schema,generator,path)
            }
            return new ArrayRenderSimple(panel,schema,generator,path)
        }
        case "map":{
            return new MapRender(panel,schema,generator,path)
        }
        case "boolean":{
            return new SwitchRender(panel,schema,path)
        }
    }
    if (schema["enum"]){
        return new BaseEnumRender(panel,schema,path)
    }
    return new BaseInputRender(panel,schema,path)
}

class FormRender {

    constructor(panel,schema,generator) {

        if(!generator || typeof generator !== "function"){
            generator = BaseGenerator
        }
        $(panel).empty()
        this.Object = generator($(panel),schema,generator,RootId)
    }
    check(){
        return this.Object.check()
    }
    get Value(){
        return this.Object.Value
    }
    set Value(v){
        this.Object.Value = v
    }
}


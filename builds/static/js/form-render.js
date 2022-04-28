
function FormRender(panel,schema,generator){
    const RootId = "FormRender"
    function CheckBySchema(schema,value){
        return false
    }
    function ValidHandler(schema){
        if(CheckBySchema(schema,$(this).val())){
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
            ValidHandler.apply(this,schema)
        })
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
    }

    function createLabel(schema,id,appendAttr){
        if (!appendAttr){
            appendAttr =""
        }
        let label = schema["label"]
        if (!label || label.trim() === ""){
            label = schema["name"]
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
        labelHtml += label+'</label>'
        return labelHtml
    }

    function createFieldPanel(panel,schema,id,appendAttr){
        $(panel).append('<div class="form-group row"></div>')
        let fieldPanel = $(panel).children().last()
        fieldPanel.append(createLabel(schema,id,appendAttr))
        fieldPanel.append('<div class="col-sm-10"></div>')
        return fieldPanel.children().last()
    }

    function createInput(schema,id,appendAttr){
        let readOnly = ""
        this.schema = schema
        if (schema["readonly"] === true){
            readOnly = "readonly"
        }

        if (schema["enum"]){
            let enums = schema["enum"]

            let select = '<select '+readOnly+'  class="form-control" id="'+id+'">'
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
        let input = '<input '+readOnly +' class="form-control" id="'+id+'" ';
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
            case "integer":{
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
        input+=' />'

        return input
    }
    class MapRender{
        constructor(panel,schema,generator,path) {}
        set Value(v){

        }
        get Value(){
            return {}
        }
    }
    class InterObjRender{
        constructor(panel,schema,generator,path) {}
        set Value(v){

        }
        get Value(){
            return {}
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
            this.Enum = schema["enum"]

            let p =$(panel);
            p.append('<div id="'+Id+'_items" class="border p-sm-1 btn-toolbar " role="toolbar"></div>')
            const itemPanel = p.children('#'+Id+'_items')
            for (let i in schema["enum"]){
                let e  = schema["enum"][i]
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

            const Id = path;
            this.Id = Id;
            let p =$(panel);

            p.append('<div id="'+Id+'_items" class="border p-sm-1 btn-toolbar " role="toolbar"></div>')

            const itemPanel = p.children('#'+Id+'_items')
            itemPanel.append('<div class="input-group input-group-sm m-2">' +
                '<div class="input-group-prepend ">' +
                '<div class="input-group-text  btn" id="btnGroupAddon_'+Id+'_new">+</div>' +
                '</div>'+
                createInput(schema,Id+'_new','aria-describedby="btnGroupAddon_'+Id+'_new" placeholder="Input new" ')+
                '</div>')

            $('#'+Id+"_new").on("change",function (){
                let v = $(this).val()

                if (v !== ""){
                    if (ValidHandler.apply(this,schema)){
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
                    createInput(schema,itemId,appendAtt) +
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
            let list =  $("[array-for='"+this.Id+"']")
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
            $("[array-for='"+this.Id+"']").each(function (){
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
                        return new InterObjRender(panel,items,generator,path)
                    }
                    case "map":{
                        return new InterMapRender(panel,items,generator,path)
                    }
                }
                if (items["enum"]){
                    return new ArrayRenderEnum(panel,items,generator,path)
                }
                return new ArrayRenderSimple(panel,items,generator,path)
            }
            case "map":{
                return new MapRender(panel,schema,generator,path)
            }
        }
        return new BaseInputRender(panel,schema,path)
    }
    class TopFormRender {

        constructor(panel,schema,generator) {

            if(!generator || typeof generator !== "function"){
                generator = BaseGenerator
            }
            $(panel).html('<form></form>')

            this.Object = generator($(panel).children("form"),schema,generator,RootId)
        }
        get Value(){
            return this.Object.Value
        }
        set Value(v){
            this.Object.Value = v
        }
    }

    return new TopFormRender(panel,schema,generator);
}



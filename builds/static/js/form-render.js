
function FormRender(panel,schema,generator){
    const RootId = "FormRender"
    function CheckBySchema(schema,value){
        let validator = validate.djv()
        if (!validator.resolved.hasOwnProperty(this.id)) {
            validator.addSchema(this.id, schema);
        }
        let err = validator.validate(this.id, value)
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
        let rs = CheckBySchema.apply(this,[schema,value])
        // let validPanel = $('validation_'+$(this).attr("id"))
        if( rs === true) {
            $(this).removeClass("is-invalid")
            $(this).addClass("is-valid")
            // validPanel.html("")
            // validPanel.remove("invalid-feedback")
            // validPanel.add("valid-feedback")

            return true
        }else{
            $(this).removeClass("is-valid")
            $(this).addClass("is-invalid")

            // validPanel.html(rs)
            // validPanel.remove("valid-feedback")
            // validPanel.add("invalid-feedback")

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

        let input = '<input '+readOnly +' class="form-control is-valid" id="'+id+'" aria-describedby="validation_'+id+'" ';
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
        if(schema["required"]){
            input += ' required'
        }
        input+='/>'
            // '<div id="validation_'+id+'" class="valid-feedback">\n</div>'

        return input
    }
    class MapRender{
        constructor(panel,schema,generator,path) {}
        set Value(v){

        }
        get Value(){
            return {}
        }
        check(){

        }
    }
    class InterObjRender{
        constructor(panel,schema,generator,path) {

        }
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

                if (v !== ""){
                    if (ValidHandler.apply(this,[items])){
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
            $(arrayId).each(function (){
               if( CheckBySchema(this.id,items,$(this).val()) !== true){
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

    class TopFormRender {

        constructor(panel,schema,generator) {

            if(!generator || typeof generator !== "function"){
                generator = BaseGenerator
            }
            $(panel).html('<form class=""></form>')

            this.Object = generator($(panel).children("form"),schema,generator,RootId)
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
        }
        if (schema["enum"]){
            return new BaseEnumRender(panel,schema,path)
        }
        return new BaseInputRender(panel,schema,path)
    }
    return new TopFormRender(panel,schema,generator);
}



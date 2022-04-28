
function HtmlHandler(htmlOpen,close){
    if (!htmlOpen || typeof htmlOpen !== "string"  ){
        $.error("need html open")
    }
    this.openHtml = htmlOpen
    if(!close){
        close =""
    }
    close += "\n"

    this.close = close
    this.items=[]
    this.Html=function (){
        let html = this.openHtml
        if (this.items){
            for(let i in this.items){
                html += this.items[i].Html()
            }
        }

        html += this.close
        return html
    }
    this.Append = function (handler){
        if(!handler){
            return
        }
        if(typeof handler === "string"){
            this.items.push(new HtmlHandler(handler))
        }else{
            if (handler.hasOwnProperty("Html")){
                this.items.push(handler)
            }
        }
    }
}

const defaultRenderHandler={
    "string":{
        "*":{
            pre:function (schema){

            },
            label:function (schema,data,id){

            },
            form:function (schema,data,id){

            },
            data:function (id){

            }
        }
    },
    "object":{

    }
}

function FormRender(panel,filter){
    this.myPanel = panel
    this.myFilter = $.extend(defaultRenderHandler,filter)
    let RootHandler = new HtmlHandler("<form>","</form>")

    this.RenderData = function (schema,data) {

        if (schema["type"] !== "object"){
            $.error("root schema must object")
        }

        let divRow = new HtmlHandler("<div class=\"form-group\">","</div>")
        divRow.Append("<label for=\"exampleInputEmail1\">Email address</label>")
        divRow.Append("<input type=\"email\" class=\"form-control\" id=\"exampleInputEmail1\" aria-describedby=\"emailHelp\"/>")
        divRow.Append("<small id=\"emailHelp\" class=\"form-text text-muted\">We'll never share your email with anyone else.</small>")


        RootHandler.Append(divRow)
        $(this.myPanel).html(RootHandler.Html())
    }

    this.doRender = function (filters,schema,data,id){
        const typeName = schema["type"]
        switch (typeName){
            case "object":{

            }
        }
    }
}









(function( $ ){
    const render = function (filter,schema,data){

        const renderField = function(schemaItem,path){
           const typeName = schemaItem["type"];
           let currentFilter = filter[typeName];
           if (!currentFilter){
               $.error("unknown type:"+typeName)
           }

        }

    }



    const renderDefault = {

        "string":{
            "*":function (schema,value,path){
                let label = schema["name"]
                if (schema["title"]){
                    label = schema["title"]
                }
                switch (schema["format"]){
                    case "password":{
                        if (!label){
                            label = "Password"
                        }
                        return ' <div class="form-group row">\n' +
                               '    <label for="inputPassword" class="col-sm-2 col-form-label">'+label+'</label>\n' +
                               '    <div class="col-sm-10">\n' +
                               '      <input type="password" class="form-control" id="'+path+'">\n' +
                               '    </div>\n' +
                               '  </div>'
                    }
                    case "email":{
                        if (!label){
                            label = "Email"
                        }
                        return ' <div class="form-group row">\n' +
                            '    <label for="inputPassword" class="col-sm-2 col-form-label">'+label+'</label>\n' +
                            '    <div class="col-sm-10">\n' +
                            '      <input type="password" class="form-control" id="'+path+'">\n' +
                            '    </div>\n' +
                            '  </div>'
                    }
                }
                return '<'
            }
        }}
    var methods = {
        init : function( options ) {
            let filter = $.extend(renderDefault,options["filter"])
            let initData = options["data"]
            let schema = options["schema"]

            return this.each(function(){
                $(this).data("render-filter",filter)
                $(this).data("render-initData",initData)
                $(this).data("render-schema",schema)
                if (schema["type"] !== "object"){
                    $.error("root schema must object")
                }

                let html = ""
                html+="<form>"
                html += render(schema,".")
                html+="</form>"
                return html

               $(this).html(
                   render(filter,schema,initData,".")
                )
            });
        },
        data : function() {
            let set =  $(this).data("render-initData")
            if (!set){
                return {};
            }
            return set;
        }
    };
    $.fn.formRender = function(method) {
        if ( methods[method] ) {
            return methods[method].apply( this, Array.prototype.slice.call( arguments, 1 ));
        } else if ( typeof method === 'object' || ! method ) {
            return methods.init.apply( this, arguments );
        } else {
            $.error( 'Method ' +  method + ' does not exist on jQuery.formRender' );
        }
    };

})( jQuery );


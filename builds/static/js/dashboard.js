let http = {
    post: function (url, data, success, error){
        this.ajax("POST", url, data, success, error)
    },
    put: function (url, data, success, error){
        this.ajax("PUT", url, data, success, error)
    },
    get: function (url, data, success, error){
        this.ajax("GET", url, data, success, error)
    },
    delete: function (url, data, success, error){
        this.ajax("DELETE", url, data, success, error)
    },
    patch: function (url, data, success, error){
        this.ajax("PATCH", url, data, success, error)
    },
    ajax: function (type, url, data, success, error, complete, async) {
        let options = {
            url: url,
            type: type,
            contentType: "application/json; charset=utf-8",
            dataType: 'json',
            success: function (res) {
                if (success) {
                    success(res)
                }
            },
            error: function (res) {
                if (error) {
                    error(res)
                }
            },
            complete: function (res) {
                if (complete) {
                    complete(res)
                }
            },
        }
        if (data) {
            options["data"] = JSON.stringify(data)
        }
        if(!async){
            options["async"] = false
        }
        $.ajax(options);
    },
    handleError: function (res, msg){
        if (res['msg']){
            common.message(res['msg'], "danger")
        }else {
            common.message(msg, "danger")
        }
    }
}
let dashboard = {
    get: function (url, success, error){
        http.get(url, null,  success, error)
    },
    update: function (url, data, success, error) {
        http.put(url, data,  success, error)
    },
    patch: function (url, data, success, error) {
        http.patch(url, data,  success, error)
    },
    delete:function (url, success, error){
        http.delete(url, success, error)
    },
    create: function (url, data, success, error){
        http.post(url, data,  success, error)
    },
    getWithAsync: function (url, success, error){
        http.ajax("GET", url, null, success, error, null, false)
    },
    getExtenders: function (success, error){
        this.get("/api/extenders/", success, error)
    },
    getExtenderInfo: function (id, success, error){
        this.get("/api/extenders/"+id, success, error)
    }
}
let common = {
    /**
     * 弹出确认框
     * @param title
     * @param msg
     * @param success
     * @param cancel
     */
    confirm :function (title, msg, success, cancel){
        let model = $("#confirmModel")
        if (model.length > 0) {
            model.remove();
        }
        let html = `<div class="modal fade" id="confirmModel" tabindex="-1" aria-labelledby="confirmModelLabel" aria-hidden="true">
                <div class="modal-dialog">
                    <div class="modal-content">
                        <div class="modal-header">
                            <h5 class="modal-title" id="confirmModelLabel">`+title+`</h5>
                            <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                                <span aria-hidden="true">&times;</span>
                            </button>
                        </div>
                        <div class="modal-body">`+msg+`</div>
                        <div class="modal-footer">
                            <button type="button" id="cancel_btn" class="btn btn-secondary">取消</button>
                            <button type="button" id="confirm_btn" class="btn btn-danger" >删除</button>
                        </div>
                    </div>
                </div>
            </div>`;
        $("body").append(html);
        model = $("#confirmModel")
        model.modal("show");
        model.on("click", "button#cancel_btn",function() {
            if (cancel) {
                cancel()
            }
            $("#confirmModel").modal("hide");
        });
        model.on("click", "button#confirm_btn",function() {
            if (success) {
                success()
            }
            $("#confirmModel").modal("hide");
        });
    },
    /**
     * 弹出消息框
     * @param msg 消息内容
     * @param type 消息框类型（参考bootstrap的alert） danger,success,warning,info...
     */
    alert: function(msg, type){
        if(typeof(type) =="undefined") { // 未传入type则默认为success类型的消息框
            type = "success";
        }
        // 创建bootstrap的alert元素
        let divElement = $("<div></div>").addClass('alert').addClass('alert-'+type).addClass('alert-dismissible').addClass('col-md-4').addClass('col-md-offset-4');
        let scroll = document.body.scrollTop || document.documentElement.scrollTop;
        divElement.css({ // 消息框的定位样式
            "position": "fixed",
            "right":"50px",
            "top": scroll + 80 +"px"
        });
        divElement.text(msg); // 设置消息框的内容
        // 消息框添加可以关闭按钮
        let closeBtn = $('<button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">×</span></button>');
        $(divElement).append(closeBtn);
        // 消息框放入到页面中
        $('body').append(divElement);
        divElement.css("z-index","999999")
        return divElement;
    },

    /**
     * 短暂显示后上浮消失的消息框
     * @param msg 消息内容
     * @param type 消息框类型
     */
    message: function(msg, type) {
        let divElement = this.alert(msg, type); // 生成Alert消息框
        let isIn = false; // 鼠标是否在消息框中

        divElement.on({ // 在setTimeout执行之前先判定鼠标是否在消息框中
            mouseover : function(){isIn = true;},
            mouseout  : function(){isIn = false;}
        });

        // 短暂延时后上浮消失
        setTimeout(function() {
            let IntervalMS = 20; // 每次上浮的间隔毫秒
            let floatSpace = 60; // 上浮的空间(px)
            let scroll = document.body.scrollTop || document.documentElement.scrollTop;
            let nowTop = divElement.offset().top - scroll; // 获取元素当前的top值
            let stopTop = nowTop - floatSpace;    // 上浮停止时的top值
            divElement.fadeOut(IntervalMS * floatSpace); // 设置元素淡出

            let upFloat = setInterval(function(){ // 开始上浮
                if (nowTop >= stopTop) { // 判断当前消息框top是否还在可上升的范围内
                    divElement.css({"top": nowTop--}); // 消息框的top上升1px
                } else {
                    clearInterval(upFloat); // 关闭上浮
                    divElement.remove();    // 移除元素
                }
            }, IntervalMS);

            if (isIn) { // 如果鼠标在setTimeout之前已经放在的消息框中，则停止上浮
                clearInterval(upFloat);
                divElement.stop();
            }

            divElement.hover(function() { // 鼠标悬浮时停止上浮和淡出效果，过后恢复
                clearInterval(upFloat);
                divElement.stop();
            },function() {
                divElement.fadeOut(IntervalMS * (nowTop - stopTop)); // 这里设置元素淡出的时间应该为：间隔毫秒*剩余可以上浮空间
                upFloat = setInterval(function(){ // 继续上浮
                    if (nowTop >= stopTop) {
                        divElement.css({"top": nowTop--});
                    } else {
                        clearInterval(upFloat); // 关闭上浮
                        divElement.remove();    // 移除元素
                    }
                }, IntervalMS);
            });
        }, 1500);
    }
}
let aceEditor = {
    InitEditor: function (id){
        let editor = ace.edit(id);
        editor.setFontSize(14)
        editor.setTheme("ace/theme/crimson_editor");
        editor.session.setMode("ace/mode/json");
        editor.renderer.setScrollMargin(10, 10);
        editor.setOptions({
            autoScrollEditorIntoView: true
        });
        return editor
    }
}
let modal = {
    table: function (id){
        let target = $("#"+id)
        target.removeClass("pop_window").removeClass("pop_window_small").html("")
        target.addClass("pop_window").addClass("pop_window_small").append(`<div class="pop_window_header">
            <span class="pop_window_title" id="${id}_title"></span>
            <button class="pop_window_button btn btn_default" id="${id}_close" >关闭</button>
            <br>
        </div>
        <div class="card pop_window_body">
            <div class="card-body">
                <table class="table table-bordered">
                    <thead class="thead-light">
                    <tr>
                        <th>属性</th>
                        <th>配置</th>
                    </tr>
                    </thead>
                    <tbody id="${id}_body">
                    </tbody>
                </table>
    
            </div>
        </div>`
        )
        let close = $(`#${id}_close`)
        close.unbind("click");
        close.click(function (){
            target.animate({
                width:'toggle'
            }, "fast", function () {
                $(`#${id}_title`).text("")
                $(`#${id}_body`).html("")
                if($("body").hasClass("modal-open")){
                    $('body').removeClass("modal-open")
                    $("div.modal-backdrop.fade.show").remove()
                }
            });
        })
        $("body").on('click', 'div.modal-backdrop.fade.show', function () {
            $(`#${id}_close`).click()
        })
        target.show = function (title, body) {
            $(`#${id}_title`).text(title)
            $(`#${id}_body`).html(body)
            $("body").append("<div class='modal-backdrop fade show'></div>").addClass("modal-open")
            target.animate({
                width:'toggle'
            }, "fast");
        }
        target.hide = function () {
            close.click()
        }
        target.destroy = function () {
            target.hide()
            close.unbind("click")
            target.removeClass("pop_window").removeClass("pop_window_small").html("")
        }
        return target
    }
}

let validate = {
    _validator: null,
    djv: function () {
        if(this._validator){
            return this._validator
        }
        this._validator = new djv({
            errorHandler(type) {
                return `errors.push({
                  type: '${type}',
                  schema: '${this.schema[this.schema.length - 1]}',
                  data: '${this.data[this.data.length - 1]}'
                });`;
            }
        })
        return this._validator
    }
}

class Render {
    constructor(id, schema, data, callback) {
        let target = $("#"+id)
        let renderHandler = new FormRender(target, schema)
        if(data){
            renderHandler.Value = data
        }
        if (callback){
            callback()
        }
        this.pane = renderHandler
        this.id = id
    }
    Reset(schema){
        this.Destroy()
        this.pane = new FormRender($("#"+this.id), schema)
    }
    ResetVal(){
        if (this.pane){
            this.pane.Value = {}
        }
    }
    Value(){
        if (this.pane){
            return this.pane.Value
        }
        return {}
    }
    Check() {
        if (this.pane){
            return this.pane.check()
        }
        return false
    }
    Submit(success, error){
        if(this.Check()){
            success(this.Value())
        }else {
            error()
        }
    }
    Destroy(){
        this.id = null
        this.pane = null
    }
}
class ProfessionRender {
    constructor(module, name, uiID){
        this.module = module
        this.name = name
        this.uiId = uiID
        this.buttonId = this.uiId+"_button"
        this.ui = null
        this.generateBtn()
    }
    generateBtn(){
        let o = this
        $(`#${this.uiId}`).after(`<div id="${this.buttonId}" class="row justify-content-between" style="display: none">
                                    <div class="col-4">
                                        <button type="button" class="btn btn-outline-secondary ${this.buttonId}_reset">重置</button>
                                    </div>
                                    <div class="col-4" style="text-align: right">
                                        <button type="button" class="btn btn-primary ${this.buttonId}_submit">提交</button>
                                    </div>
                                </div>`)
        $(`#${this.buttonId}`).on("click", `button.${this.buttonId}_reset`, function () {
            o.resetEvent()
        }).on("click", `button.${this.buttonId}_submit`, function () {
            o.submitEvent()
        })
    }
    resetEvent(){
        if (this.ui){
            this.ui.ResetVal()
        }
    }
    submitEvent(){
        if (this.ui){
            this.ui.Submit(function (data) {
                dashboard.update(`/api/${this.module}/${this.name}`,  data, function (res){
                    if(res.code !== 200){
                        http.handleError(res, "更新失败")
                        return
                    }
                    common.message("更新成功", "success")
                }, function (res){
                    http.handleError(res, "更新失败")
                })
            }, function () {
                common.message("config format error", "danger")
            })
        }

    }
}
class ProfessionCreator extends ProfessionRender{
    constructor(module, name, uiID, callback){
        super(module, name, uiID)
        this.initRender()
        if (callback){
            callback()
        }
    }

}
class ProfessionEditor extends ProfessionRender{
    constructor(module, name, uiID, callback){
        super(module, name, uiID)
        this.initRender()
        if (callback){
            callback()
        }
    }
    initRender() {
        let o = this
        dashboard.get(`/api/${this.module}/${this.name}`,function (res) {
            if(res.code !== 200){
                return http.handleError(res, "获取详情失败")
            }
            if(!res.data["driver"]){
                return http.handleError(res, "获取driver失败")
            }
            o.initRenderForm(res.data["driver"], res.data)
        },function (res) {
            return http.handleError(res, "获取详情失败")
        })
    }
    initRenderForm(driver, data){
        let o = this
        dashboard.get(`/profession/${this.module}/${driver}`, function (res) {
            if(res.code !== 200){
                return http.handleError(res, "获取render失败")
            }
            o.ui = new Render(o.uiId, res.data["render"], data, function () {
                let btn = "#" + o.buttonId
                if (!$(btn).is(":visible")){
                    $(btn).show()
                }
            })
        }, function (res){
            return http.handleError(res, "获取render失败")
        })
    }
}
let demoSchema = {
    "type": "object",
    "properties": [
        {
            "name": "driver",
            "readonly": true,
            "type": "string",
            "enum": [
                "http"
            ],
            "required": true
        },
        {
            "name": "host",
            "type": "array",
            "items": {
                "type": "string"
            },
            "minLength": 1
        },
        {
            "name": "listen",
            "type": "integer",
            "format": "int32",
            "minimum": 1,
            "required": true
        },
        {
            "name": "method",
            "type": "array",
            "items": {
                "type": "string",
                "enum": [
                    "GET",
                    "POST",
                    "PATH",
                    "DELETE"
                ]
            }
        },
        {
            "name": "plugins",
            "type": "map",
            "items": {
                "type": "object",
                "properties": [
                    {
                        "name": "config",
                        "required": true
                    },
                    {
                        "name": "disable",
                        "type": "boolean",
                        "required": true
                    }
                ]
            }
        },
        {
            "name": "protocol",
            "type": "string",
            "enum": [
                "http",
                "https"
            ],
            "default": "http",
            "required": true
        },
        {
            "name": "rules",
            "type": "array",
            "items": {
                "type": "object",
                "properties": [
                    {
                        "name": "header",
                        "type": "map",
                        "items": {
                            "type": "string"
                        }
                    },
                    {
                        "name": "location",
                        "type": "string",
                        "minLength": 1
                    },
                    {
                        "name": "query",
                        "type": "map",
                        "items": {
                            "type": "string"
                        }
                    }
                ]
            }
        },
        {
            "name": "target",
            "type": "string",
            "format": "text",
            "minLength": 1,
            "required": true
        }
    ]
}
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
    getRender: function (url, success, error){
        http.ajax("GET", url, null, success, error, null, false)
    },
}
let common = {
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
            console.log(nowTop)
            console.log(stopTop)
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
let util = {
    filter: function (high, low){
        let data = {}
        Object.keys(high).forEach(function (key){
            if(low[key]){
                data[key] = low[key]
            }else {
                data[key] = null
            }
        })
        return data
    }
}
let JsonEditor = {
    // default_Schema: {"type":"object","properties":{"cert":{"type":"array","items":{"type":"object","properties":{"crt":{"type":"string"},"key":{"type":"string"}},"additionalProperties":false,"required":["key","crt"]}},"driver":{"type":"string","enum":["http"]},"host":{"type":"array","items":{"type":"string"},"minLength":1},"listen":{"type":"integer","format":"int32","minimum":1},"method":{"type":"array","items":{"type":"string","enum":["GET","POST","PATH","DELETE"]}},"plugins":{"type":"object","additionalProperties":{"type":"object","properties":{"config":{},"disable":{"type":"boolean"}},"additionalProperties":false,"required":["disable","config"]}},"protocol":{"type":"string","enum":["http","https"],"default":"http"},"rules":{"type":"array","items":{"type":"object","properties":{"header":{"type":"object","additionalProperties":{"type":"string"}},"location":{"type":"string","minLength":1},"query":{"type":"object","additionalProperties":{"type":"string"}}},"additionalProperties":false}},"target":{"type":"string","minLength":1}},"additionalProperties":false,"required":["driver","listen","protocol","target"]},
    default_Schema: {
        "type": "info",
        "title": "Apinto",
        "description": "operation error!",
        "properties":{},
    },
    init: function (theme, iconlib, callbacks){
        JSONEditor.defaults.options.theme = theme;
        JSONEditor.defaults.options.iconlib = iconlib;
        JSONEditor.defaults.callbacks = callbacks
    },
    getEditorWithData: function (id, title, schemaUrl, options, dataUrl) {
        let editor = new JSONEditor(document.getElementById(id), this.getOptions(schemaUrl, title, options));
        this.setValue(editor, dataUrl)
        return editor
    },
    getEditor: function (id, title, schemaUrl, options) {
        return new JSONEditor(document.getElementById(id), this.getOptions(schemaUrl, title, options));
    },
    getOptions: function (url, title, options) {
        options["schema"] = this.getSchema(url, title)
        return options
    },
    getSchema: function (url, title){
        let schema = this.default_Schema
        dashboard.getRender(url, function (res) {
            if(res.code === 200){
                schema = res.data
                schema["properties"]["operation"] = {
                    "type": "button",
                    "title": "Submit",
                    "options": {
                        "button": {
                            "action": "submit",
                            "validated": true
                        },
                        "inputAttributes": {
                            "class": "btn btn-primary"
                        }
                    }
                }
                schema["title"]= title
                return schema
            }
            http.handleError(res, "获取render失败")
        }, function (res){
            http.handleError(res, "获取render失败")
        })
        schema["title"]= title
        return schema
    },
    updateSchema: function(editor, url, title){
        let id = editor.element.id
        let options = editor.options
        editor.destroy()
        return this.getEditor(id, title, url, options)
    },
    setValue: function (editor, url){
        dashboard.get(url, function (res) {
            if(res.code === 200){
                if (editor["schema"]["properties"]){
                    editor.setValue(util.filter(editor["schema"]["properties"], res.data))
                }else {
                    editor.setValue(res.data)
                }
                return
            }
            http.handleError(res, "获取详情失败")
            editor.disable()
        }, function (res) {
            http.handleError(res, "获取详情失败")
            editor.disable()
        })
    },
    update: function (url, data){
        dashboard.update(url, data, function (res) {
            if(res.code === 200){
                common.message("success", "success")
                return
            }
            http.handleError(res, "update error")

        }, function (res) {
            http.handleError(res, "update error")
        })
    },
    create: function (url, data){
        dashboard.create(url, data, function (res) {
            if(res.code === 200){
                common.message("success", "success")
                return
            }
            http.handleError(res, "create error")

        }, function (res) {
            http.handleError(res, "create error")
        })
    }
}
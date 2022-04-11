let http = {
    post: function (url, data, success, error){
        this.ajax("POST", url, data, success, error)
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
    ajax: function (type, url, data, success, error, complete) {
        $.ajax({
            url: url,
            type: type,
            data: data,
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
        });
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
        divElement.css({ // 消息框的定位样式
            "position": "absolute",
            "right":"50px",
            "top": "80px"
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
            let nowTop = divElement.offset().top; // 获取元素当前的top值
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

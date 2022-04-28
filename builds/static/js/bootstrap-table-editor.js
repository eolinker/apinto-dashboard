/**
 * author: netcy
 */
!function ($) {
    'use strict';

    const getItemField = $.fn.bootstrapTable.utils.getItemField;
    const calculateObjectValue  = $.fn.bootstrapTable.utils.calculateObjectValue;

    function createTextEditor(){
        const editor = $('<input class = "form-control"></input>');
        editor.css("width",`100%`);
        return editor;
    }

    function createNumberEditor(){
        var editor = $('<input class = "form-control"></input>');
        editor.attr('type', 'number');
        editor.attr('step', '0.1');
        editor.css("width",`100%`);
        return editor;
    }

    function createSelectEditor(){
        var select = $('<select class="form-control"></select>');
        return select;
    }

    const textEditor = createTextEditor();
    const numberEditor = createNumberEditor();

    const editorCache = {
        text:textEditor,
        number:numberEditor,
        select:createSelectEditor(),
    };

    $.extend($.fn.bootstrapTable.defaults, {
        editable: false,
        reInitOnEdit:true,
        onEditorInit: function () {
            return false;
        },
        onEditorSave: function (field, row, oldValue, $el) {
            return false;
        },
        onEditorShown: function (field, row, $el, editable) {
            return false;
        },
    });

    $.extend($.fn.bootstrapTable.Constructor.EVENTS, {
        'editor-init.bs.table': 'onEditorInit',
        'editor-save.bs.table': 'onEditorSave',
        'editor-shown.bs.table': 'onEditorShown',
    });

    var BootstrapTable = $.fn.bootstrapTable.Constructor,
        _initTable = BootstrapTable.prototype.initTable,
        _initBody = BootstrapTable.prototype.initBody;

    BootstrapTable.prototype.getEditor = function(key){
        return editorCache[key] || textEditor;
    }

    BootstrapTable.prototype.initTable = function () {
        var that = this;
        _initTable.apply(this, Array.prototype.slice.apply(arguments));

        if (!this.options.editable) {
            return;
        }

        $.each(this.columns, function (i, column) {
            if (!column.editable) {
                return;
            }
            that._initColumnEditor(column);
        });

        this.trigger('editor-init');
    };

    BootstrapTable.prototype._initColumnEditor = function(column){
        var editable = column.editable;
        var editorType = editable.type;
        var editor = this.getEditor(editorType);
        column._editor = editor;
    };

    BootstrapTable.prototype._setColumnEditorOptions = function(column){
        var editable = column.editable;
        var editor = column._editor;
        if(!editor) {
            return;
        }

        var options = editable.options || {},
            attributes = options.attributes,
            styles = options.styles;
        if(editor._options == options) { //避免重复
            return;
        }
        if(attributes){
            Object.keys(attributes).forEach(function(key){
                editor.attr(key,attributes[key]);
            });
        }
        if(styles){
            Object.keys(styles).forEach(function(key){
                editor.css(key,styles[key]);
            });
        }
        if(editable.type == "select" && options.items){
            editor.empty();
            $.each(options.items, function (index, val) {
                let label = val.label || val.value || val;
                let value = val.value || val;
                var opt = $('<option value="' + value + '">' + label + '</option>');
                editor.append(opt);
            });
        }
        editor._options = options;
    }

    BootstrapTable.prototype._startEdit = function($td){
        var that = this;
        var $tr = $td.parent(),
            item = that.data[$tr.data('index')],
            index = $td[0].cellIndex,
            fields = that.getVisibleFields(),
            field = fields[that.options.detailView && !that.options.cardView ? index - 1 : index],
            column = that.columns[that.fieldsColumnsIndex[field]],
            value = getItemField(item, field, that.options.escape),
            dataIndex = $td.parent().data('index');
        var editor = column._editor;
        if(!editor){
            return;
        }
        if(editor._row != item || editor._index != index){
            that._setColumnEditorOptions(column);
            editor._row = item;
            editor._index = index;

            editor.keydown(function(e) {
                if (e.keyCode == 13) {
                    that._saveEditData($td,editor,padding);
                }
            });
            editor.blur(function(e){
                that._saveEditData($td,editor,padding);
            });
            editor.change(function(e){
                that._saveEditData($td,editor,padding);
            });
            var padding = $td.css("padding");
            $td.html("");
            $td.css("padding","0px");
            $td.append(editor);
            editor.focus();
            editor.val(value);
            that.trigger('editor-shown', column.field, item, $td);
        }
    }

    BootstrapTable.prototype._saveEditData = function($td,editor,padding){
        let that = this;
        var $tr = $td.parent(),
            item = that.data[$tr.data('index')],
            index = $td[0].cellIndex,
            fields = that.getVisibleFields(),
            field = fields[that.options.detailView && !that.options.cardView ? index - 1 : index],
            column = that.columns[that.fieldsColumnsIndex[field]],
            value = getItemField(item, field, that.options.escape),
            dataIndex = $td.parent().data('index');

        let tdValue = editor.val();
        if(tdValue != value){
            this.updateCell({
                index: dataIndex,
                field: field,
                value: tdValue,
                reinit:this.options.reInitOnEdit,
            });
            if(!this.options.reInitOnEdit) {
                let args = [tdValue,item,index,field];
                let formatText = calculateObjectValue(column,column.formatter,args,tdValue);
                $td.html(formatText);
                $td.css("padding",padding);
            }
            that.trigger('editor-save', column.field, item, value, $td);
        } else {
            let args = [value,item,index,field];
            let formatText = calculateObjectValue(column,column.formatter,args,value);
            $td.html(formatText);
            $td.css("padding",padding);
        }
        editor._row = undefined;
        editor._index = undefined;
    };

    BootstrapTable.prototype.initBody = function () {
        var that = this;
        _initBody.apply(this, Array.prototype.slice.apply(arguments));

        if (!this.options.editable) {
            return;
        }

        this.$body.find('> tr[data-index] > td').on('click', function (e) {
            var $td = $(this);
            that._startEdit($td);
        });
    };

}(jQuery);
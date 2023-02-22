'use strict';

class FileItem {

    constructor(name,size,type,data) {
        this.$name = name;
        this.$size = size;
        this.$data = data;
        this.$type = type;
    }
    get name(){
        return this.$name
    }
    get size(){
        return this.$size
    }
    get index(){
        return this.$size+"_"+this.name
    }
    get type(){
        return this.$type
    }
    get Value(){
        let data = this.$dataCache
        if (!data){
            data = window.btoa(window.pako.gzip(this.$data).reduce((data, byte) => data + String.fromCharCode(byte), ''));
            this.$dataCache = data
        }

        return {
            name:this.$name,
            size:this.$size,
            type:this.$type,
            data:data
        }
    }
    set Data(v){
        this.$data = v
    }
    get DataUrl(){
        if (this.$dataUrlCache){
            return this.$dataUrlCache;
        }

        let data = window.btoa(new Uint8Array(this.$data)
            .reduce((data, byte) => data + String.fromCharCode(byte), '')
        )
        this.$dataUrlCache = `data:${this.$type};base64,${data}`
        return this.$dataUrlCache
    }
}

function DecodeFile(obj){

   return new FileItem(obj.name,obj.size,obj.type, window.pako.ungzip(_base64ToArrayBuffer(obj.data)))

}
function _base64ToArrayBuffer(base64) {
    let binary_string = window.atob(base64);
    let len = binary_string.length;
    let bytes = new Uint8Array(len);
    for (let i = 0; i < len; i++) {
        bytes[i] = binary_string.charCodeAt(i);
    }
    return bytes.buffer;
}

function _getFileName(file){
    if (!file) {
        return null;
    }
    let relativePath = String(file.relativePath || file.webkitRelativePath || (file ? (file.fileName || file.name || '') : '') || null);
    if (!relativePath) {
        return null;
    }

    return encodeURIComponent(relativePath).replace(/[%\- ]/g, '_');
}

!function(v,w){"object"==typeof exports&&"object"==typeof module?module.exports=w():"function"==typeof define&&define.amd?define([],w):"object"==typeof exports?exports.ClipboardJS=w():v.ClipboardJS=w()}(this,function(){return w={686:function(f,n,t){"use strict";t.d(n,{default:function(){return R}}),n=t(279);var c=t.n(n),y=(n=t(370),t.n(n)),l=(n=t(817),t.n(n));function p(o){try{return document.execCommand(o)}catch{return}}var g=function(o){return o=l()(o),p("cut"),o};function h(s,i){var r,a;return r=s,a="rtl"===document.documentElement.getAttribute("dir"),(s=document.createElement("textarea")).style.fontSize="12pt",s.style.border="0",s.style.padding="0",s.style.margin="0",s.style.position="absolute",s.style[a?"right":"left"]="-9999px",a=window.pageYOffset||document.documentElement.scrollTop,s.style.top="".concat(a,"px"),s.setAttribute("readonly",""),s.value=r,i.container.appendChild(s),i=l()(s),p("copy"),s.remove(),i}var m=function(o){var i=1<arguments.length&&void 0!==arguments[1]?arguments[1]:{container:document.body},r="";return"string"==typeof o?r=h(o,i):o instanceof HTMLInputElement&&!["text","search","url","tel","password"].includes(o?.type)?r=h(o.value,i):(r=l()(o),p("copy")),r};function x(o){return(x="function"==typeof Symbol&&"symbol"==typeof Symbol.iterator?function(i){return typeof i}:function(i){return i&&"function"==typeof Symbol&&i.constructor===Symbol&&i!==Symbol.prototype?"symbol":typeof i})(o)}function S(o){return(S="function"==typeof Symbol&&"symbol"==typeof Symbol.iterator?function(i){return typeof i}:function(i){return i&&"function"==typeof Symbol&&i.constructor===Symbol&&i!==Symbol.prototype?"symbol":typeof i})(o)}function E(o,i){for(var r=0;r<i.length;r++){var a=i[r];a.enumerable=a.enumerable||!1,a.configurable=!0,"value"in a&&(a.writable=!0),Object.defineProperty(o,a.key,a)}}function A(o,i){return(A=Object.setPrototypeOf||function(r,a){return r.__proto__=a,r})(o,i)}function k(o){return(k=Object.setPrototypeOf?Object.getPrototypeOf:function(i){return i.__proto__||Object.getPrototypeOf(i)})(o)}function L(o,i){if(o="data-clipboard-".concat(o),i.hasAttribute(o))return i.getAttribute(o)}var R=function(){!function(u,d){if("function"!=typeof d&&null!==d)throw new TypeError("Super expression must either be null or a function");u.prototype=Object.create(d&&d.prototype,{constructor:{value:u,writable:!0,configurable:!0}}),d&&A(u,d)}(s,c());var o,i,r,a=function C(o){var i=function(){if(typeof Reflect>"u"||!Reflect.construct||Reflect.construct.sham)return!1;if("function"==typeof Proxy)return!0;try{return Date.prototype.toString.call(Reflect.construct(Date,[],function(){})),!0}catch{return!1}}();return function(){var r,a=k(o);return r=i?(r=k(this).constructor,Reflect.construct(a,arguments,r)):a.apply(this,arguments),a=this,!r||"object"!==S(r)&&"function"!=typeof r?function(s){if(void 0!==s)return s;throw new ReferenceError("this hasn't been initialised - super() hasn't been called")}(a):r}}(s);function s(u,d){var b;return function(j){if(!(j instanceof s))throw new TypeError("Cannot call a class as a function")}(this),(b=a.call(this)).resolveOptions(d),b.listenClick(u),b}return o=s,r=[{key:"copy",value:function(u){var d=1<arguments.length&&void 0!==arguments[1]?arguments[1]:{container:document.body};return m(u,d)}},{key:"cut",value:function(u){return g(u)}},{key:"isSupported",value:function(){var u="string"==typeof(u=0<arguments.length&&void 0!==arguments[0]?arguments[0]:["copy","cut"])?[u]:u,d=!!document.queryCommandSupported;return u.forEach(function(b){d=d&&!!document.queryCommandSupported(b)}),d}}],(i=[{key:"resolveOptions",value:function(){var u=0<arguments.length&&void 0!==arguments[0]?arguments[0]:{};this.action="function"==typeof u.action?u.action:this.defaultAction,this.target="function"==typeof u.target?u.target:this.defaultTarget,this.text="function"==typeof u.text?u.text:this.defaultText,this.container="object"===S(u.container)?u.container:document.body}},{key:"listenClick",value:function(u){var d=this;this.listener=y()(u,"click",function(b){return d.onClick(b)})}},{key:"onClick",value:function(j){var d=j.delegateTarget||j.currentTarget,b=this.action(d)||"copy";j=function(){var o=void 0===(r=(a=0<arguments.length&&void 0!==arguments[0]?arguments[0]:{}).action)?"copy":r,i=a.container,r=a.target,a=a.text;if("copy"!==o&&"cut"!==o)throw new Error('Invalid "action" value, use either "copy" or "cut"');if(void 0!==r){if(!r||"object"!==x(r)||1!==r.nodeType)throw new Error('Invalid "target" value, use a valid Element');if("copy"===o&&r.hasAttribute("disabled"))throw new Error('Invalid "target" attribute. Please use "readonly" instead of "disabled" attribute');if("cut"===o&&(r.hasAttribute("readonly")||r.hasAttribute("disabled")))throw new Error('Invalid "target" attribute. You can\'t cut text from elements with "readonly" or "disabled" attributes')}return a?m(a,{container:i}):r?"cut"===o?g(r):m(r,{container:i}):void 0}({action:b,container:this.container,target:this.target(d),text:this.text(d)}),this.emit(j?"success":"error",{action:b,text:j,trigger:d,clearSelection:function(){d&&d.focus(),window.getSelection().removeAllRanges()}})}},{key:"defaultAction",value:function(u){return L("action",u)}},{key:"defaultTarget",value:function(u){if(u=L("target",u))return document.querySelector(u)}},{key:"defaultText",value:function(u){return L("text",u)}},{key:"destroy",value:function(){this.listener.destroy()}}])&&E(o.prototype,i),r&&E(o,r),s}()},828:function(f){var e;typeof Element>"u"||Element.prototype.matches||((e=Element.prototype).matches=e.matchesSelector||e.mozMatchesSelector||e.msMatchesSelector||e.oMatchesSelector||e.webkitMatchesSelector),f.exports=function(t,c){for(;t&&9!==t.nodeType;){if("function"==typeof t.matches&&t.matches(c))return t;t=t.parentNode}}},438:function(f,e,t){var c=t(828);function y(n,l,p,g,h){var m=function(x,T,S,E){return function(A){A.delegateTarget=c(A.target,T),A.delegateTarget&&E.call(x,A)}}.apply(this,arguments);return n.addEventListener(p,m,h),{destroy:function(){n.removeEventListener(p,m,h)}}}f.exports=function(n,l,p,g,h){return"function"==typeof n.addEventListener?y.apply(null,arguments):"function"==typeof p?y.bind(null,document).apply(null,arguments):("string"==typeof n&&(n=document.querySelectorAll(n)),Array.prototype.map.call(n,function(m){return y(m,l,p,g,h)}))}},879:function(f,e){e.node=function(t){return void 0!==t&&t instanceof HTMLElement&&1===t.nodeType},e.nodeList=function(t){var c=Object.prototype.toString.call(t);return void 0!==t&&("[object NodeList]"===c||"[object HTMLCollection]"===c)&&"length"in t&&(0===t.length||e.node(t[0]))},e.string=function(t){return"string"==typeof t||t instanceof String},e.fn=function(t){return"[object Function]"===Object.prototype.toString.call(t)}},370:function(f,e,t){var c=t(879),y=t(438);f.exports=function(n,l,p){if(!n&&!l&&!p)throw new Error("Missing required arguments");if(!c.string(l))throw new TypeError("Second argument must be a String");if(!c.fn(p))throw new TypeError("Third argument must be a Function");if(c.node(n))return(x=n).addEventListener(T=l,S=p),{destroy:function(){x.removeEventListener(T,S)}};if(c.nodeList(n))return g=n,h=l,m=p,Array.prototype.forEach.call(g,function(E){E.addEventListener(h,m)}),{destroy:function(){Array.prototype.forEach.call(g,function(E){E.removeEventListener(h,m)})}};if(c.string(n))return y(document.body,n,l,p);throw new TypeError("First argument must be a String, HTMLElement, HTMLCollection, or NodeList");var g,h,m,x,T,S}},817:function(f){f.exports=function(e){var t,c="SELECT"===e.nodeName?(e.focus(),e.value):"INPUT"===e.nodeName||"TEXTAREA"===e.nodeName?((t=e.hasAttribute("readonly"))||e.setAttribute("readonly",""),e.select(),e.setSelectionRange(0,e.value.length),t||e.removeAttribute("readonly"),e.value):(e.hasAttribute("contenteditable")&&e.focus(),c=window.getSelection(),(t=document.createRange()).selectNodeContents(e),c.removeAllRanges(),c.addRange(t),c.toString());return c}},279:function(f){function e(){}e.prototype={on:function(t,c,y){var n=this.e||(this.e={});return(n[t]||(n[t]=[])).push({fn:c,ctx:y}),this},once:function(t,c,y){var n=this;function l(){n.off(t,l),c.apply(y,arguments)}return l._=c,this.on(t,l,y)},emit:function(t){for(var c=[].slice.call(arguments,1),y=((this.e||(this.e={}))[t]||[]).slice(),n=0,l=y.length;n<l;n++)y[n].fn.apply(y[n].ctx,c);return this},off:function(t,c){var y=this.e||(this.e={}),n=y[t],l=[];if(n&&c)for(var p=0,g=n.length;p<g;p++)n[p].fn!==c&&n[p].fn._!==c&&l.push(n[p]);return l.length?y[t]=l:delete y[t],this}},f.exports=e,f.exports.TinyEmitter=e}},O={},v.n=function(f){var e=f&&f.__esModule?function(){return f.default}:function(){return f};return v.d(e,{a:e}),e},v.d=function(f,e){for(var t in e)v.o(e,t)&&!v.o(f,t)&&Object.defineProperty(f,t,{enumerable:!0,get:e[t]})},v.o=function(f,e){return Object.prototype.hasOwnProperty.call(f,e)},v(686).default;function v(f){if(O[f])return O[f].exports;var e=O[f]={exports:{}};return w[f](e,e.exports,v),e.exports}var w,O});
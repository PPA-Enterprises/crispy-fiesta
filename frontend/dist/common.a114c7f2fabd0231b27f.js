(window.webpackJsonp=window.webpackJsonp||[]).push([[1],{oR8h:function(t,e,n){"use strict";n.d(e,"a",(function(){return r}));var p=n("hy6k"),i=n("67Y/"),c=n("CcnG"),o=n("t/Na"),r=function(){function t(t,e){this.http=t,this.jobService=e}return t.prototype.getAllClients=function(){return this.http.get("http://ppaenterprises.com:8888/api/v1/clients/").pipe(Object(i.a)((function(t){return t.success?t.payload:[]})))},t.prototype.getClientById=function(t){return this.http.get("http://ppaenterprises.com:8888/api/v1/clients/id/"+t).pipe(Object(i.a)((function(t){return t.success?t.payload:null})))},t.prototype.editClientById=function(t,e){return this.http.patch("http://ppaenterprises.com:8888/api/v1/clients/"+t,e).pipe(Object(i.a)((function(t){return t.success?t.payload:null})))},t.\u0275prov=c["\u0275\u0275defineInjectable"]({factory:function(){return new t(c["\u0275\u0275inject"](o.c),c["\u0275\u0275inject"](p.a))},token:t,providedIn:"root"}),t}()}}]);
(window.webpackJsonp=window.webpackJsonp||[]).push([[1],{oR8h:function(t,n,e){"use strict";e.d(n,"a",(function(){return o}));var i=e("hy6k"),c=e("67Y/"),p=e("CcnG"),r=e("t/Na"),o=function(){function t(t,n){this.http=t,this.jobService=n,this.url="http://ppaenterprises.com/"}return t.prototype.getAllClients=function(){return this.http.get(this.url+"api/v1/clients/").pipe(Object(c.a)((function(t){return t.success?t.payload:[]})))},t.prototype.getClientById=function(t){return this.http.get(this.url+"api/v1/clients/id/"+t).pipe(Object(c.a)((function(t){return t.success?t.payload:null})))},t.prototype.editClientById=function(t,n){return this.http.patch(this.url+"api/v1/clients/"+t,n).pipe(Object(c.a)((function(t){return t.success?t.payload:null})))},t.\u0275prov=p["\u0275\u0275defineInjectable"]({factory:function(){return new t(p["\u0275\u0275inject"](r.c),p["\u0275\u0275inject"](i.a))},token:t,providedIn:"root"}),t}()}}]);
"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[1549],{3905:(e,t,n)=>{n.d(t,{Zo:()=>u,kt:()=>f});var r=n(7294);function a(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function o(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function i(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?o(Object(n),!0).forEach((function(t){a(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):o(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function l(e,t){if(null==e)return{};var n,r,a=function(e,t){if(null==e)return{};var n,r,a={},o=Object.keys(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||(a[n]=e[n]);return a}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(a[n]=e[n])}return a}var s=r.createContext({}),p=function(e){var t=r.useContext(s),n=t;return e&&(n="function"==typeof e?e(t):i(i({},t),e)),n},u=function(e){var t=p(e.components);return r.createElement(s.Provider,{value:t},e.children)},d="mdxType",c={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},g=r.forwardRef((function(e,t){var n=e.components,a=e.mdxType,o=e.originalType,s=e.parentName,u=l(e,["components","mdxType","originalType","parentName"]),d=p(n),g=a,f=d["".concat(s,".").concat(g)]||d[g]||c[g]||o;return n?r.createElement(f,i(i({ref:t},u),{},{components:n})):r.createElement(f,i({ref:t},u))}));function f(e,t){var n=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var o=n.length,i=new Array(o);i[0]=g;var l={};for(var s in t)hasOwnProperty.call(t,s)&&(l[s]=t[s]);l.originalType=e,l[d]="string"==typeof e?e:a,i[1]=l;for(var p=2;p<o;p++)i[p]=n[p];return r.createElement.apply(null,i)}return r.createElement.apply(null,n)}g.displayName="MDXCreateElement"},9030:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>s,contentTitle:()=>i,default:()=>c,frontMatter:()=>o,metadata:()=>l,toc:()=>p});var r=n(7462),a=(n(7294),n(3905));const o={sidebar_position:1,title:"Basic Configuration"},i="Basic Configuration",l={unversionedId:"configuration/storage/index",id:"configuration/storage/index",title:"Basic Configuration",description:"Storage is the place where terediX will store the discovered data. You can add multiple storage engines in the",source:"@site/docs/configuration/storage/index.md",sourceDirName:"configuration/storage",slug:"/configuration/storage/",permalink:"/docs/configuration/storage/",draft:!1,editUrl:"https://github.com/shaharia-lab/teredix/tree/master/website/docs/configuration/storage/index.md",tags:[],version:"current",lastUpdatedAt:1700477589,formattedLastUpdatedAt:"Nov 20, 2023",sidebarPosition:1,frontMatter:{sidebar_position:1,title:"Basic Configuration"},sidebar:"tutorialSidebar",previous:{title:"Storage",permalink:"/docs/category/storage"},next:{title:"PostgreSQL",permalink:"/docs/configuration/storage/postgresql"}},s={},p=[{value:"Supported Storage Engines",id:"supported-storage-engines",level:3}],u={toc:p},d="wrapper";function c(e){let{components:t,...n}=e;return(0,a.kt)(d,(0,r.Z)({},u,n,{components:t,mdxType:"MDXLayout"}),(0,a.kt)("h1",{id:"basic-configuration"},"Basic Configuration"),(0,a.kt)("p",null,"Storage is the place where terediX will store the discovered data. You can add multiple storage engines in the\nconfiguration file. You can also add multiple storage engines."),(0,a.kt)("p",null,"If you add multiple storage engines, you must need to define ",(0,a.kt)("inlineCode",{parentName:"p"},"default_engine")),(0,a.kt)("table",null,(0,a.kt)("thead",{parentName:"table"},(0,a.kt)("tr",{parentName:"thead"},(0,a.kt)("th",{parentName:"tr",align:null},"option"),(0,a.kt)("th",{parentName:"tr",align:null},"type"),(0,a.kt)("th",{parentName:"tr",align:"left"},"description"))),(0,a.kt)("tbody",{parentName:"table"},(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:null},"batch_size"),(0,a.kt)("td",{parentName:"tr",align:null},"number"),(0,a.kt)("td",{parentName:"tr",align:"left"},"Number of data to store in a single batch. You can increase the number to speed up the storage process.")),(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:null},"engines"),(0,a.kt)("td",{parentName:"tr",align:null},"object"),(0,a.kt)("td",{parentName:"tr",align:"left"},"List of storage engines. You can add multiple storage engines.")),(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:null},"default_engine"),(0,a.kt)("td",{parentName:"tr",align:null},"text"),(0,a.kt)("td",{parentName:"tr",align:"left"},"Name of the default storage engine. If you add multiple storage engines, you must need to define ",(0,a.kt)("inlineCode",{parentName:"td"},"default_engine"))))),(0,a.kt)("h3",{id:"supported-storage-engines"},"Supported Storage Engines"),(0,a.kt)("p",null,"list of supported storage engines"),(0,a.kt)("table",null,(0,a.kt)("thead",{parentName:"table"},(0,a.kt)("tr",{parentName:"thead"},(0,a.kt)("th",{parentName:"tr",align:null},"storage engine"),(0,a.kt)("th",{parentName:"tr",align:"left"},"description"))),(0,a.kt)("tbody",{parentName:"table"},(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:null},"postgresql"),(0,a.kt)("td",{parentName:"tr",align:"left"},"Store data in PostgreSQL database. You can use this storage engine to store data in PostgreSQL.")))))}c.isMDXComponent=!0}}]);
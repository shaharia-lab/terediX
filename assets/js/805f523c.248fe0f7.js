"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[9665],{3905:(e,t,n)=>{n.d(t,{Zo:()=>u,kt:()=>m});var r=n(7294);function o(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function i(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function a(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?i(Object(n),!0).forEach((function(t){o(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):i(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function l(e,t){if(null==e)return{};var n,r,o=function(e,t){if(null==e)return{};var n,r,o={},i=Object.keys(e);for(r=0;r<i.length;r++)n=i[r],t.indexOf(n)>=0||(o[n]=e[n]);return o}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(r=0;r<i.length;r++)n=i[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(o[n]=e[n])}return o}var s=r.createContext({}),c=function(e){var t=r.useContext(s),n=t;return e&&(n="function"==typeof e?e(t):a(a({},t),e)),n},u=function(e){var t=c(e.components);return r.createElement(s.Provider,{value:t},e.children)},p="mdxType",d={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},f=r.forwardRef((function(e,t){var n=e.components,o=e.mdxType,i=e.originalType,s=e.parentName,u=l(e,["components","mdxType","originalType","parentName"]),p=c(n),f=o,m=p["".concat(s,".").concat(f)]||p[f]||d[f]||i;return n?r.createElement(m,a(a({ref:t},u),{},{components:n})):r.createElement(m,a({ref:t},u))}));function m(e,t){var n=arguments,o=t&&t.mdxType;if("string"==typeof e||o){var i=n.length,a=new Array(i);a[0]=f;var l={};for(var s in t)hasOwnProperty.call(t,s)&&(l[s]=t[s]);l.originalType=e,l[p]="string"==typeof e?e:o,a[1]=l;for(var c=2;c<i;c++)a[c]=n[c];return r.createElement.apply(null,a)}return r.createElement.apply(null,n)}f.displayName="MDXCreateElement"},9813:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>s,contentTitle:()=>a,default:()=>d,frontMatter:()=>i,metadata:()=>l,toc:()=>c});var r=n(7462),o=(n(7294),n(3905));const i={title:"File System"},a="File System",l={unversionedId:"configuration/scanner/file_system",id:"configuration/scanner/file_system",title:"File System",description:"Configuration",source:"@site/docs/configuration/scanner/file_system.md",sourceDirName:"configuration/scanner",slug:"/configuration/scanner/file_system",permalink:"/docs/configuration/scanner/file_system",draft:!1,editUrl:"https://github.com/shaharia-lab/teredix/tree/master/website/docs/configuration/scanner/file_system.md",tags:[],version:"current",lastUpdatedAt:1696875037,formattedLastUpdatedAt:"Oct 9, 2023",frontMatter:{title:"File System"},sidebar:"tutorialSidebar",previous:{title:"AWS S3",permalink:"/docs/configuration/scanner/aws_s3"},next:{title:"GitHub Repository",permalink:"/docs/configuration/scanner/github_repository"}},s={},c=[{value:"Configuration",id:"configuration",level:2},{value:"Type",id:"type",level:3},{value:"Configuration",id:"configuration-1",level:3},{value:"Fields",id:"fields",level:3},{value:"Schedule",id:"schedule",level:3},{value:"Example",id:"example",level:2}],u={toc:c},p="wrapper";function d(e){let{components:t,...n}=e;return(0,o.kt)(p,(0,r.Z)({},u,n,{components:t,mdxType:"MDXLayout"}),(0,o.kt)("h1",{id:"file-system"},"File System"),(0,o.kt)("img",{src:"/img/file_system_icon.png",alt:"File System",width:"250"}),(0,o.kt)("h2",{id:"configuration"},"Configuration"),(0,o.kt)("h3",{id:"type"},"Type"),(0,o.kt)("p",null,"Resource type. In this case it would be ",(0,o.kt)("inlineCode",{parentName:"p"},"file_system"),"."),(0,o.kt)("h3",{id:"configuration-1"},"Configuration"),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},(0,o.kt)("strong",{parentName:"li"},"root_directory"),": Provide the absolute path of the directory to scan.")),(0,o.kt)("h3",{id:"fields"},"Fields"),(0,o.kt)("p",null,"List of available fields to add with resource add metadata. During scanning resources, scanner will only fetch data\nfrom the following fields."),(0,o.kt)("ul",null,(0,o.kt)("li",{parentName:"ul"},"machineHost"),(0,o.kt)("li",{parentName:"ul"},"rootDirectory")),(0,o.kt)("h3",{id:"schedule"},"Schedule"),(0,o.kt)("p",null,(0,o.kt)("strong",{parentName:"p"},"\ud83d\udd17 ",(0,o.kt)("a",{parentName:"strong",href:"/docs/configuration/scanner/overview#schedule-format"},"Check schedule format")),"."),(0,o.kt)("h2",{id:"example"},"Example"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-yaml"},'source:\n  file_system_source_one:\n      type: file_system\n      configuration:\n        root_directory: "/file/path"\n      fields:\n        - machineHost\n        - rootDirectory\n      schedule: "@every 24h"\n')),(0,o.kt)("p",null,"In the above example, we have added a source named ",(0,o.kt)("inlineCode",{parentName:"p"},"file_system_source_one")," with type ",(0,o.kt)("inlineCode",{parentName:"p"},"file_system"),". We have added some fields to add with each resource.\nWe have also set the schedule to run this source every 24 hours."),(0,o.kt)("p",null,"Based on the above example, scanner_name would be ",(0,o.kt)("inlineCode",{parentName:"p"},"file_system_source_one")," and scanner_type would be ",(0,o.kt)("inlineCode",{parentName:"p"},"file_system"),". This is\nimportant to filter resources in Grafana dashboard."))}d.isMDXComponent=!0}}]);
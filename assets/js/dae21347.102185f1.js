"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[1549],{37382:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>a,contentTitle:()=>o,default:()=>g,frontMatter:()=>r,metadata:()=>d,toc:()=>l});var s=n(85893),i=n(11151);const r={sidebar_position:1,title:"Basic Configuration"},o="Basic Configuration",d={id:"configuration/storage/index",title:"Basic Configuration",description:"Storage is the place where terediX will store the discovered data. You can add multiple storage engines in the",source:"@site/docs/configuration/storage/index.md",sourceDirName:"configuration/storage",slug:"/configuration/storage/",permalink:"/docs/configuration/storage/",draft:!1,unlisted:!1,editUrl:"https://github.com/shaharia-lab/teredix/tree/master/website/docs/configuration/storage/index.md",tags:[],version:"current",lastUpdatedAt:1711041757,formattedLastUpdatedAt:"Mar 21, 2024",sidebarPosition:1,frontMatter:{sidebar_position:1,title:"Basic Configuration"},sidebar:"tutorialSidebar",previous:{title:"Storage",permalink:"/docs/category/storage"},next:{title:"PostgreSQL",permalink:"/docs/configuration/storage/postgresql"}},a={},l=[{value:"Supported Storage Engines",id:"supported-storage-engines",level:3}];function c(e){const t={code:"code",h1:"h1",h3:"h3",p:"p",table:"table",tbody:"tbody",td:"td",th:"th",thead:"thead",tr:"tr",...(0,i.a)(),...e.components};return(0,s.jsxs)(s.Fragment,{children:[(0,s.jsx)(t.h1,{id:"basic-configuration",children:"Basic Configuration"}),"\n",(0,s.jsx)(t.p,{children:"Storage is the place where terediX will store the discovered data. You can add multiple storage engines in the\nconfiguration file. You can also add multiple storage engines."}),"\n",(0,s.jsxs)(t.p,{children:["If you add multiple storage engines, you must need to define ",(0,s.jsx)(t.code,{children:"default_engine"})]}),"\n",(0,s.jsxs)(t.table,{children:[(0,s.jsx)(t.thead,{children:(0,s.jsxs)(t.tr,{children:[(0,s.jsx)(t.th,{children:"option"}),(0,s.jsx)(t.th,{children:"type"}),(0,s.jsx)(t.th,{style:{textAlign:"left"},children:"description"})]})}),(0,s.jsxs)(t.tbody,{children:[(0,s.jsxs)(t.tr,{children:[(0,s.jsx)(t.td,{children:"batch_size"}),(0,s.jsx)(t.td,{children:"number"}),(0,s.jsx)(t.td,{style:{textAlign:"left"},children:"Number of data to store in a single batch. You can increase the number to speed up the storage process."})]}),(0,s.jsxs)(t.tr,{children:[(0,s.jsx)(t.td,{children:"engines"}),(0,s.jsx)(t.td,{children:"object"}),(0,s.jsx)(t.td,{style:{textAlign:"left"},children:"List of storage engines. You can add multiple storage engines."})]}),(0,s.jsxs)(t.tr,{children:[(0,s.jsx)(t.td,{children:"default_engine"}),(0,s.jsx)(t.td,{children:"text"}),(0,s.jsxs)(t.td,{style:{textAlign:"left"},children:["Name of the default storage engine. If you add multiple storage engines, you must need to define ",(0,s.jsx)(t.code,{children:"default_engine"})]})]})]})]}),"\n",(0,s.jsx)(t.h3,{id:"supported-storage-engines",children:"Supported Storage Engines"}),"\n",(0,s.jsx)(t.p,{children:"list of supported storage engines"}),"\n",(0,s.jsxs)(t.table,{children:[(0,s.jsx)(t.thead,{children:(0,s.jsxs)(t.tr,{children:[(0,s.jsx)(t.th,{children:"storage engine"}),(0,s.jsx)(t.th,{style:{textAlign:"left"},children:"description"})]})}),(0,s.jsx)(t.tbody,{children:(0,s.jsxs)(t.tr,{children:[(0,s.jsx)(t.td,{children:"postgresql"}),(0,s.jsx)(t.td,{style:{textAlign:"left"},children:"Store data in PostgreSQL database. You can use this storage engine to store data in PostgreSQL."})]})})]})]})}function g(e={}){const{wrapper:t}={...(0,i.a)(),...e.components};return t?(0,s.jsx)(t,{...e,children:(0,s.jsx)(c,{...e})}):c(e)}},11151:(e,t,n)=>{n.d(t,{Z:()=>d,a:()=>o});var s=n(67294);const i={},r=s.createContext(i);function o(e){const t=s.useContext(r);return s.useMemo((function(){return"function"==typeof e?e(t):{...t,...e}}),[t,e])}function d(e){let t;return t=e.disableParentContext?"function"==typeof e.components?e.components(i):e.components||i:o(e.components),s.createElement(r.Provider,{value:t},e.children)}}}]);
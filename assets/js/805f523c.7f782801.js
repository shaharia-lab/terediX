"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[9665],{1664:(e,n,i)=>{i.r(n),i.d(n,{assets:()=>c,contentTitle:()=>o,default:()=>u,frontMatter:()=>r,metadata:()=>l,toc:()=>a});var s=i(5893),t=i(1151);const r={title:"Configure File System Resource Scanner for TerediX",sidebar_label:"File System"},o="File System",l={id:"configuration/scanner/file_system",title:"Configure File System Resource Scanner for TerediX",description:"Configuration",source:"@site/docs/configuration/scanner/file_system.md",sourceDirName:"configuration/scanner",slug:"/configuration/scanner/file_system",permalink:"/docs/configuration/scanner/file_system",draft:!1,unlisted:!1,editUrl:"https://github.com/shaharia-lab/teredix/tree/master/website/docs/configuration/scanner/file_system.md",tags:[],version:"current",lastUpdatedAt:1704148309,formattedLastUpdatedAt:"Jan 1, 2024",frontMatter:{title:"Configure File System Resource Scanner for TerediX",sidebar_label:"File System"},sidebar:"tutorialSidebar",previous:{title:"AWS S3",permalink:"/docs/configuration/scanner/aws_s3"},next:{title:"GitHub Repository",permalink:"/docs/configuration/scanner/github_repository"}},c={},a=[{value:"Configuration",id:"configuration",level:2},{value:"Type",id:"type",level:3},{value:"Configuration",id:"configuration-1",level:3},{value:"Fields",id:"fields",level:3},{value:"Schedule",id:"schedule",level:3},{value:"Example",id:"example",level:2}];function d(e){const n={a:"a",code:"code",h1:"h1",h2:"h2",h3:"h3",li:"li",p:"p",pre:"pre",strong:"strong",ul:"ul",...(0,t.a)(),...e.components};return(0,s.jsxs)(s.Fragment,{children:[(0,s.jsx)(n.h1,{id:"file-system",children:"File System"}),"\n",(0,s.jsx)("img",{src:"/img/file_system_icon.png",alt:"File System",width:"250"}),"\n",(0,s.jsx)(n.h2,{id:"configuration",children:"Configuration"}),"\n",(0,s.jsx)(n.h3,{id:"type",children:"Type"}),"\n",(0,s.jsxs)(n.p,{children:["Resource type. In this case it would be ",(0,s.jsx)(n.code,{children:"file_system"}),"."]}),"\n",(0,s.jsx)(n.h3,{id:"configuration-1",children:"Configuration"}),"\n",(0,s.jsxs)(n.ul,{children:["\n",(0,s.jsxs)(n.li,{children:[(0,s.jsx)(n.strong,{children:"root_directory"}),": Provide the absolute path of the directory to scan."]}),"\n"]}),"\n",(0,s.jsx)(n.h3,{id:"fields",children:"Fields"}),"\n",(0,s.jsx)(n.p,{children:"List of available fields to add with resource add metadata. During scanning resources, scanner will only fetch data\nfrom the following fields."}),"\n",(0,s.jsxs)(n.ul,{children:["\n",(0,s.jsx)(n.li,{children:"machineHost"}),"\n",(0,s.jsx)(n.li,{children:"rootDirectory"}),"\n"]}),"\n",(0,s.jsx)(n.h3,{id:"schedule",children:"Schedule"}),"\n",(0,s.jsxs)(n.p,{children:[(0,s.jsxs)(n.strong,{children:["\ud83d\udd17 ",(0,s.jsx)(n.a,{href:"/docs/configuration/scanner/overview#schedule-format",children:"Check schedule format"})]}),"."]}),"\n",(0,s.jsx)(n.h2,{id:"example",children:"Example"}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-yaml",children:'source:\n  file_system_source_one:\n      type: file_system\n      configuration:\n        root_directory: "/file/path"\n      fields:\n        - machineHost\n        - rootDirectory\n      schedule: "@every 24h"\n'})}),"\n",(0,s.jsxs)(n.p,{children:["In the above example, we have added a source named ",(0,s.jsx)(n.code,{children:"file_system_source_one"})," with type ",(0,s.jsx)(n.code,{children:"file_system"}),". We have added some fields to add with each resource.\nWe have also set the schedule to run this source every 24 hours."]}),"\n",(0,s.jsxs)(n.p,{children:["Based on the above example, scanner_name would be ",(0,s.jsx)(n.code,{children:"file_system_source_one"})," and scanner_type would be ",(0,s.jsx)(n.code,{children:"file_system"}),". This is\nimportant to filter resources in Grafana dashboard."]})]})}function u(e={}){const{wrapper:n}={...(0,t.a)(),...e.components};return n?(0,s.jsx)(n,{...e,children:(0,s.jsx)(d,{...e})}):d(e)}},1151:(e,n,i)=>{i.d(n,{Z:()=>l,a:()=>o});var s=i(7294);const t={},r=s.createContext(t);function o(e){const n=s.useContext(r);return s.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function l(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(t):e.components||t:o(e.components),s.createElement(r.Provider,{value:n},e.children)}}}]);
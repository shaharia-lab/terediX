"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[3549],{13524:(e,n,r)=>{r.r(n),r.d(n,{assets:()=>a,contentTitle:()=>c,default:()=>d,frontMatter:()=>i,metadata:()=>o,toc:()=>l});var s=r(85893),t=r(11151);const i={sidebar_position:2,title:"Technical Architecture"},c="Technical Architecture",o={id:"technical_architecture",title:"Technical Architecture",description:"Components",source:"@site/docs/technical_architecture.md",sourceDirName:".",slug:"/technical_architecture",permalink:"/docs/technical_architecture",draft:!1,unlisted:!1,editUrl:"https://github.com/shaharia-lab/teredix/tree/master/website/docs/technical_architecture.md",tags:[],version:"current",lastUpdatedAt:1713174558e3,sidebarPosition:2,frontMatter:{sidebar_position:2,title:"Technical Architecture"},sidebar:"tutorialSidebar",previous:{title:"Introduction",permalink:"/docs/intro"},next:{title:"Getting Started",permalink:"/docs/getting_started"}},a={},l=[{value:"Components",id:"components",level:2},{value:"Configuration file",id:"configuration-file",level:3},{value:"Scanner",id:"scanner",level:3},{value:"Scheduler",id:"scheduler",level:3},{value:"Processor",id:"processor",level:3},{value:"Storage",id:"storage",level:3},{value:"Metrics",id:"metrics",level:3}];function h(e){const n={code:"code",h1:"h1",h2:"h2",h3:"h3",li:"li",p:"p",pre:"pre",strong:"strong",ul:"ul",...(0,t.a)(),...e.components};return(0,s.jsxs)(s.Fragment,{children:[(0,s.jsx)(n.h1,{id:"technical-architecture",children:"Technical Architecture"}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{children:"                                    +----------+\n                                    | Scheduler |\n                                    +----------+\n                                          |\n                                          |\n                                          v\n                +--------------+     +--------------+\n                | Resource     |     | Resource     |\n                | Scanner 1    |     | Scanner 2    |\n                +--------------+     +--------------+\n                                          |\n                                          |\n                                          v\n                                  +----------+\n                                  | Processor |\n                                  +----------+\n                                          |\n                                          |\n                                          v\n                                  +----------+\n                                  | Storage   |\n                                  +----------+\n"})}),"\n",(0,s.jsx)(n.h2,{id:"components",children:"Components"}),"\n",(0,s.jsx)(n.h3,{id:"configuration-file",children:"Configuration file"}),"\n",(0,s.jsx)(n.p,{children:"The configuration file is a JSON file that contains all the necessary settings for TerediX to run. The configuration file specifies the following:"}),"\n",(0,s.jsxs)(n.ul,{children:["\n",(0,s.jsx)(n.li,{children:"The list of scanners to use."}),"\n",(0,s.jsx)(n.li,{children:"The schedule for each scanner."}),"\n",(0,s.jsx)(n.li,{children:"The storage configuration."}),"\n",(0,s.jsx)(n.li,{children:"The metrics configuration."}),"\n"]}),"\n",(0,s.jsx)(n.h3,{id:"scanner",children:"Scanner"}),"\n",(0,s.jsx)(n.p,{children:"Scanners are responsible for fetching resources from their respective sources. Scanners can be implemented in any programming language.\nTerediX includes a number of built-in scanners, such as a scanner for AWS S3 and a scanner for Google Cloud Storage."}),"\n",(0,s.jsx)(n.h3,{id:"scheduler",children:"Scheduler"}),"\n",(0,s.jsx)(n.p,{children:"The scheduler is responsible for scheduling scanners to run at regular intervals. The scheduler can be implemented\nin any programming language. TerediX includes a built-in scheduler that uses cron expressions."}),"\n",(0,s.jsx)(n.h3,{id:"processor",children:"Processor"}),"\n",(0,s.jsx)(n.p,{children:"The processor is responsible for processing resources as they are fetched by scanners. The processor can be implemented in any programming language.\nTerediX includes a built-in processor that stores resources and metadata in the storage."}),"\n",(0,s.jsx)(n.h3,{id:"storage",children:"Storage"}),"\n",(0,s.jsx)(n.p,{children:"The storage is responsible for storing resources and metadata. The storage can be implemented in any database or file system.\nTerediX includes a built-in storage that uses a PostgreSQL database."}),"\n",(0,s.jsx)(n.h3,{id:"metrics",children:"Metrics"}),"\n",(0,s.jsx)(n.p,{children:"The metrics component collects and exposes metrics about TerediX's operation. The metrics component can be implemented in any programming language.\nTerediX includes a built-in metrics component that exposes metrics to Prometheus."}),"\n",(0,s.jsx)(n.p,{children:(0,s.jsx)(n.strong,{children:"The components are connected as follows:"})}),"\n",(0,s.jsxs)(n.ul,{children:["\n",(0,s.jsx)(n.li,{children:"The processor reads the configuration file to determine which scanners to use and their schedules."}),"\n",(0,s.jsx)(n.li,{children:"The processor starts the scanners and processes the resources that they fetch."}),"\n",(0,s.jsx)(n.li,{children:"The processor stores the resources and metadata in the storage."}),"\n",(0,s.jsx)(n.li,{children:"The metrics component collects and exposes metrics about the processor's operation."}),"\n"]}),"\n",(0,s.jsx)(n.p,{children:"The system is designed to be scalable and reliable. The processor can be scaled horizontally to handle more resources. The storage is designed to be highly available and durable."})]})}function d(e={}){const{wrapper:n}={...(0,t.a)(),...e.components};return n?(0,s.jsx)(n,{...e,children:(0,s.jsx)(h,{...e})}):h(e)}},11151:(e,n,r)=>{r.d(n,{Z:()=>o,a:()=>c});var s=r(67294);const t={},i=s.createContext(t);function c(e){const n=s.useContext(i);return s.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function o(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(t):e.components||t:c(e.components),s.createElement(i.Provider,{value:n},e.children)}}}]);
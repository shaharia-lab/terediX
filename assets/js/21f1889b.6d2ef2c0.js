"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[6854],{91297:(e,t,r)=>{r.r(t),r.d(t,{assets:()=>d,contentTitle:()=>o,default:()=>h,frontMatter:()=>n,metadata:()=>c,toc:()=>l});var s=r(85893),i=r(11151);const n={sidebar_position:1,title:"Available Scanners and Configurations",sidebar_label:"Overview"},o="Source",c={id:"configuration/scanner/overview",title:"Available Scanners and Configurations",description:"Source is the place where terediX will discover the data. You can add multiple sources in the configuration file.",source:"@site/docs/configuration/scanner/overview.md",sourceDirName:"configuration/scanner",slug:"/configuration/scanner/overview",permalink:"/docs/configuration/scanner/overview",draft:!1,unlisted:!1,editUrl:"https://github.com/shaharia-lab/teredix/tree/master/website/docs/configuration/scanner/overview.md",tags:[],version:"current",lastUpdatedAt:1713216937e3,sidebarPosition:1,frontMatter:{sidebar_position:1,title:"Available Scanners and Configurations",sidebar_label:"Overview"},sidebar:"tutorialSidebar",previous:{title:"Scanner",permalink:"/docs/category/scanner"},next:{title:"AWS EC2",permalink:"/docs/configuration/scanner/aws_ec2"}},d={},l=[{value:"Common source configuration",id:"common-source-configuration",level:3},{value:"Supported Source Types",id:"supported-source-types",level:3},{value:"Schedule Format",id:"schedule-format",level:3}];function a(e){const t={code:"code",h1:"h1",h3:"h3",li:"li",p:"p",table:"table",tbody:"tbody",td:"td",th:"th",thead:"thead",tr:"tr",ul:"ul",...(0,i.a)(),...e.components};return(0,s.jsxs)(s.Fragment,{children:[(0,s.jsx)(t.h1,{id:"source",children:"Source"}),"\n",(0,s.jsx)(t.p,{children:"Source is the place where terediX will discover the data. You can add multiple sources in the configuration file.\nEvery source will build a scanner based on the source kind and configuration."}),"\n",(0,s.jsx)(t.h3,{id:"common-source-configuration",children:"Common source configuration"}),"\n",(0,s.jsxs)(t.table,{children:[(0,s.jsx)(t.thead,{children:(0,s.jsxs)(t.tr,{children:[(0,s.jsx)(t.th,{children:"option"}),(0,s.jsx)(t.th,{children:"type"}),(0,s.jsx)(t.th,{style:{textAlign:"left"},children:"description"})]})}),(0,s.jsxs)(t.tbody,{children:[(0,s.jsxs)(t.tr,{children:[(0,s.jsx)(t.td,{children:"source_name"}),(0,s.jsx)(t.td,{children:"text"}),(0,s.jsx)(t.td,{style:{textAlign:"left"},children:"Key/name of each source. e.g: github_repositories, aws_resources_rds, aws_s3_one"})]}),(0,s.jsxs)(t.tr,{children:[(0,s.jsx)(t.td,{children:"[source_name].type"}),(0,s.jsx)(t.td,{children:"text"}),(0,s.jsx)(t.td,{style:{textAlign:"left"},children:"Type of the source. See the full list of supported source types"})]}),(0,s.jsxs)(t.tr,{children:[(0,s.jsx)(t.td,{children:"[source_name].configuration"}),(0,s.jsx)(t.td,{children:"key value pair"}),(0,s.jsx)(t.td,{style:{textAlign:"left"},children:"Configuration of the source. This configuration is different for different type of source"})]}),(0,s.jsxs)(t.tr,{children:[(0,s.jsx)(t.td,{children:"[source_name].fields"}),(0,s.jsx)(t.td,{children:"list"}),(0,s.jsx)(t.td,{style:{textAlign:"left"},children:"Additional data to store with each resource"})]}),(0,s.jsxs)(t.tr,{children:[(0,s.jsx)(t.td,{children:"[source_name].schedule"}),(0,s.jsx)(t.td,{children:"interval"}),(0,s.jsxs)(t.td,{style:{textAlign:"left"},children:["Set interval to schedule this source. e.g: ",(0,s.jsx)(t.code,{children:"@every 10s"}),", ",(0,s.jsx)(t.code,{children:"@every <br/>24h"})," or any valid cron expression ",(0,s.jsx)(t.code,{children:"*/10 * * * * *"})]})]})]})]}),"\n",(0,s.jsx)(t.h3,{id:"supported-source-types",children:"Supported Source Types"}),"\n",(0,s.jsxs)(t.table,{children:[(0,s.jsx)(t.thead,{children:(0,s.jsxs)(t.tr,{children:[(0,s.jsx)(t.th,{children:"source type"}),(0,s.jsx)(t.th,{style:{textAlign:"left"},children:"description"})]})}),(0,s.jsxs)(t.tbody,{children:[(0,s.jsxs)(t.tr,{children:[(0,s.jsx)(t.td,{children:"aws_s3"}),(0,s.jsxs)(t.td,{style:{textAlign:"left"},children:["Discover data from AWS S3. You can use this source to discover data from AWS S3. See configuration for ",(0,s.jsx)(t.code,{children:"aws_s3"})," source type."]})]}),(0,s.jsxs)(t.tr,{children:[(0,s.jsx)(t.td,{children:"aws_ec2"}),(0,s.jsxs)(t.td,{style:{textAlign:"left"},children:["Discover data from AWS EC2. You can use this source to discover data from AWS EC2. See the configuration for ",(0,s.jsx)(t.code,{children:"aws_ec2"})," source type."]})]}),(0,s.jsxs)(t.tr,{children:[(0,s.jsx)(t.td,{children:"aws_rds"}),(0,s.jsxs)(t.td,{style:{textAlign:"left"},children:["Discover data from AWS RDS. You can use this source to discover data from AWS RDS. See the configuration for ",(0,s.jsx)(t.code,{children:"aws_rds"})," source type."]})]}),(0,s.jsxs)(t.tr,{children:[(0,s.jsx)(t.td,{children:"aws_ecr"}),(0,s.jsxs)(t.td,{style:{textAlign:"left"},children:["Discover data from AWS ECR. You can use this source to discover data from AWS ECR repository. See the configuration for ",(0,s.jsx)(t.code,{children:"aws_ecr"})," source type."]})]}),(0,s.jsxs)(t.tr,{children:[(0,s.jsx)(t.td,{children:"file_system"}),(0,s.jsxs)(t.td,{style:{textAlign:"left"},children:["Discover data from local file system. You can use this source to discover data from local file system. See the configuration for ",(0,s.jsx)(t.code,{children:"file_system"})," source type."]})]}),(0,s.jsxs)(t.tr,{children:[(0,s.jsx)(t.td,{children:"github_repository"}),(0,s.jsxs)(t.td,{style:{textAlign:"left"},children:["List of GitHub repositories. See the configuration for ",(0,s.jsx)(t.code,{children:"github_repository"})," source type"]})]})]})]}),"\n",(0,s.jsx)(t.h3,{id:"schedule-format",children:"Schedule Format"}),"\n",(0,s.jsx)(t.p,{children:"You can set the schedule for each source. The schedule format is similar to the cron expression."}),"\n",(0,s.jsx)(t.p,{children:"Valid schedule formats are:"}),"\n",(0,s.jsxs)(t.ul,{children:["\n",(0,s.jsx)(t.li,{children:(0,s.jsx)(t.code,{children:"@every 10s"})}),"\n",(0,s.jsx)(t.li,{children:(0,s.jsx)(t.code,{children:"@every 1m"})}),"\n",(0,s.jsx)(t.li,{children:(0,s.jsx)(t.code,{children:"@every 1h"})}),"\n",(0,s.jsx)(t.li,{children:(0,s.jsx)(t.code,{children:"@every 1d"})}),"\n",(0,s.jsx)(t.li,{children:(0,s.jsx)(t.code,{children:"@every 1w"})}),"\n",(0,s.jsxs)(t.li,{children:[(0,s.jsx)(t.code,{children:"*/10 * * * * *"}),"  # cron expression that will run in every 10 seconds"]}),"\n"]})]})}function h(e={}){const{wrapper:t}={...(0,i.a)(),...e.components};return t?(0,s.jsx)(t,{...e,children:(0,s.jsx)(a,{...e})}):a(e)}},11151:(e,t,r)=>{r.d(t,{Z:()=>c,a:()=>o});var s=r(67294);const i={},n=s.createContext(i);function o(e){const t=s.useContext(n);return s.useMemo((function(){return"function"==typeof e?e(t):{...t,...e}}),[t,e])}function c(e){let t;return t=e.disableParentContext?"function"==typeof e.components?e.components(i):e.components||i:o(e.components),s.createElement(n.Provider,{value:t},e.children)}}}]);
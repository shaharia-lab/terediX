"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[7358],{71414:(e,n,s)=>{s.r(n),s.d(n,{assets:()=>t,contentTitle:()=>c,default:()=>u,frontMatter:()=>o,metadata:()=>a,toc:()=>d});var r=s(85893),i=s(11151);const o={title:"Configure AWS RDS Resource Scanner for TerediX",sidebar_label:"AWS RDS"},c="AWS RDS",a={id:"configuration/scanner/aws_rds",title:"Configure AWS RDS Resource Scanner for TerediX",description:"Configuration",source:"@site/docs/configuration/scanner/aws_rds.md",sourceDirName:"configuration/scanner",slug:"/configuration/scanner/aws_rds",permalink:"/docs/configuration/scanner/aws_rds",draft:!1,unlisted:!1,editUrl:"https://github.com/shaharia-lab/teredix/tree/master/website/docs/configuration/scanner/aws_rds.md",tags:[],version:"current",lastUpdatedAt:1713134108,formattedLastUpdatedAt:"Apr 14, 2024",frontMatter:{title:"Configure AWS RDS Resource Scanner for TerediX",sidebar_label:"AWS RDS"},sidebar:"tutorialSidebar",previous:{title:"AWS ECR",permalink:"/docs/configuration/scanner/aws_ecr"},next:{title:"AWS S3",permalink:"/docs/configuration/scanner/aws_s3"}},t={},d=[{value:"Configuration",id:"configuration",level:2},{value:"Type",id:"type",level:3},{value:"Configuration",id:"configuration-1",level:3},{value:"Fields",id:"fields",level:3},{value:"Schedule",id:"schedule",level:3},{value:"Example",id:"example",level:2}];function l(e){const n={a:"a",code:"code",h1:"h1",h2:"h2",h3:"h3",li:"li",p:"p",pre:"pre",strong:"strong",ul:"ul",...(0,i.a)(),...e.components};return(0,r.jsxs)(r.Fragment,{children:[(0,r.jsx)(n.h1,{id:"aws-rds",children:"AWS RDS"}),"\n",(0,r.jsx)("img",{src:"/img/aws_rds_icon.png",alt:"AWS RDS",width:"250"}),"\n",(0,r.jsx)(n.h2,{id:"configuration",children:"Configuration"}),"\n",(0,r.jsx)(n.h3,{id:"type",children:"Type"}),"\n",(0,r.jsxs)(n.p,{children:["Resource type. In this case it would be ",(0,r.jsx)(n.code,{children:"aws_rds"}),"."]}),"\n",(0,r.jsx)(n.h3,{id:"configuration-1",children:"Configuration"}),"\n",(0,r.jsxs)(n.ul,{children:["\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.strong,{children:"access_key"}),": AWS access key"]}),"\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.strong,{children:"secret_key"}),": AWS secret key"]}),"\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.strong,{children:"region"}),": AWS region. e.g: us-west-1"]}),"\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.strong,{children:"session_token"}),": AWS session token"]}),"\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.strong,{children:"account_id"}),": AWS account ID"]}),"\n"]}),"\n",(0,r.jsx)(n.h3,{id:"fields",children:"Fields"}),"\n",(0,r.jsx)(n.p,{children:"List of available fields to add with resource add metadata. During scanning resources, scanner will only fetch data\nfrom the following fields."}),"\n",(0,r.jsxs)(n.ul,{children:["\n",(0,r.jsx)(n.li,{children:"instance_id"}),"\n",(0,r.jsx)(n.li,{children:"region"}),"\n",(0,r.jsx)(n.li,{children:"arn"}),"\n",(0,r.jsx)(n.li,{children:"tags"}),"\n"]}),"\n",(0,r.jsx)(n.h3,{id:"schedule",children:"Schedule"}),"\n",(0,r.jsxs)(n.p,{children:[(0,r.jsxs)(n.strong,{children:["\ud83d\udd17 ",(0,r.jsx)(n.a,{href:"/docs/configuration/scanner/overview#schedule-format",children:"Check schedule format"})]}),"."]}),"\n",(0,r.jsx)(n.h2,{id:"example",children:"Example"}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-yaml",children:'source:\n  aws_rds_source_one:\n      type: aws_rds\n      configuration:\n        access_key: "xxxx"\n        secret_key: "xxxx"\n        session_token: "xxxx"\n        region: "x"\n        account_id: "xxx"\n      fields:\n        - instance_id\n        - region\n        - arn\n        - tags\n      schedule: "@every 24h"\n'})}),"\n",(0,r.jsxs)(n.p,{children:["In the above example, we have added a source named ",(0,r.jsx)(n.code,{children:"aws_rds_source_one"})," with type ",(0,r.jsx)(n.code,{children:"aws_rds"}),". We have added some fields to add with each resource.\nWe have also set the schedule to run this source every 24 hours."]}),"\n",(0,r.jsxs)(n.p,{children:["Based on the above example, scanner_name would be ",(0,r.jsx)(n.code,{children:"aws_rds_source_one"})," and scanner_type would be ",(0,r.jsx)(n.code,{children:"aws_rds"}),". This is\nimportant to filter resources in Grafana dashboard."]})]})}function u(e={}){const{wrapper:n}={...(0,i.a)(),...e.components};return n?(0,r.jsx)(n,{...e,children:(0,r.jsx)(l,{...e})}):l(e)}},11151:(e,n,s)=>{s.d(n,{Z:()=>a,a:()=>c});var r=s(67294);const i={},o=r.createContext(i);function c(e){const n=r.useContext(o);return r.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function a(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(i):e.components||i:c(e.components),r.createElement(o.Provider,{value:n},e.children)}}}]);
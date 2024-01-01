"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[5575],{4102:(e,n,s)=>{s.r(n),s.d(n,{assets:()=>c,contentTitle:()=>t,default:()=>u,frontMatter:()=>o,metadata:()=>a,toc:()=>l});var i=s(5893),r=s(1151);const o={title:"Configure AWS S3 Resource Scanner for TerediX",sidebar_label:"AWS S3"},t="AWS S3",a={id:"configuration/scanner/aws_s3",title:"Configure AWS S3 Resource Scanner for TerediX",description:"Configuration",source:"@site/docs/configuration/scanner/aws_s3.md",sourceDirName:"configuration/scanner",slug:"/configuration/scanner/aws_s3",permalink:"/docs/configuration/scanner/aws_s3",draft:!1,unlisted:!1,editUrl:"https://github.com/shaharia-lab/teredix/tree/master/website/docs/configuration/scanner/aws_s3.md",tags:[],version:"current",lastUpdatedAt:1704116149,formattedLastUpdatedAt:"Jan 1, 2024",frontMatter:{title:"Configure AWS S3 Resource Scanner for TerediX",sidebar_label:"AWS S3"},sidebar:"tutorialSidebar",previous:{title:"AWS RDS",permalink:"/docs/configuration/scanner/aws_rds"},next:{title:"File System",permalink:"/docs/configuration/scanner/file_system"}},c={},l=[{value:"Configuration",id:"configuration",level:2},{value:"Type",id:"type",level:3},{value:"Configuration",id:"configuration-1",level:3},{value:"Fields",id:"fields",level:3},{value:"Schedule",id:"schedule",level:3},{value:"Example",id:"example",level:2}];function d(e){const n={a:"a",code:"code",h1:"h1",h2:"h2",h3:"h3",li:"li",p:"p",pre:"pre",strong:"strong",ul:"ul",...(0,r.a)(),...e.components};return(0,i.jsxs)(i.Fragment,{children:[(0,i.jsx)(n.h1,{id:"aws-s3",children:"AWS S3"}),"\n",(0,i.jsx)("img",{src:"/img/aws_s3_icon.png",alt:"AWS S3",width:"250"}),"\n",(0,i.jsx)(n.h2,{id:"configuration",children:"Configuration"}),"\n",(0,i.jsx)(n.h3,{id:"type",children:"Type"}),"\n",(0,i.jsxs)(n.p,{children:["Resource type. In this case it would be ",(0,i.jsx)(n.code,{children:"aws_s3"}),"."]}),"\n",(0,i.jsx)(n.h3,{id:"configuration-1",children:"Configuration"}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.strong,{children:"access_key"}),": AWS access key"]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.strong,{children:"secret_key"}),": AWS secret key"]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.strong,{children:"region"}),": AWS region. e.g: us-west-1"]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.strong,{children:"session_token"}),": AWS session token"]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.strong,{children:"account_id"}),": AWS account ID"]}),"\n"]}),"\n",(0,i.jsx)(n.h3,{id:"fields",children:"Fields"}),"\n",(0,i.jsx)(n.p,{children:"List of available fields to add with resource add metadata. During scanning resources, scanner will only fetch data\nfrom the following fields."}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsx)(n.li,{children:"bucket_name"}),"\n",(0,i.jsx)(n.li,{children:"region"}),"\n",(0,i.jsx)(n.li,{children:"arn"}),"\n",(0,i.jsx)(n.li,{children:"tags"}),"\n"]}),"\n",(0,i.jsx)(n.h3,{id:"schedule",children:"Schedule"}),"\n",(0,i.jsxs)(n.p,{children:[(0,i.jsxs)(n.strong,{children:["\ud83d\udd17 ",(0,i.jsx)(n.a,{href:"/docs/configuration/scanner/overview#schedule-format",children:"Check schedule format"})]}),"."]}),"\n",(0,i.jsx)(n.h2,{id:"example",children:"Example"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-yaml",children:'source:\n  aws_s3_source_one:\n      type: aws_s3\n      configuration:\n        access_key: "xxxx"\n        secret_key: "xxxx"\n        session_token: "xxxx"\n        region: "x"\n        account_id: "xxx"\n      fields:\n        - bucket_name\n        - region\n        - arn\n        - tags\n      schedule: "@every 24h"\n'})}),"\n",(0,i.jsxs)(n.p,{children:["In the above example, we have added a source named ",(0,i.jsx)(n.code,{children:"aws_s3_source_one"})," with type ",(0,i.jsx)(n.code,{children:"aws_s3"}),". We have added some fields to add with each resource.\nWe have also set the schedule to run this source every 24 hours."]}),"\n",(0,i.jsxs)(n.p,{children:["Based on the above example, scanner_name would be ",(0,i.jsx)(n.code,{children:"aws_s3_source_one"})," and scanner_type would be ",(0,i.jsx)(n.code,{children:"aws_s3"}),". This is\nimportant to filter resources in Grafana dashboard."]})]})}function u(e={}){const{wrapper:n}={...(0,r.a)(),...e.components};return n?(0,i.jsx)(n,{...e,children:(0,i.jsx)(d,{...e})}):d(e)}},1151:(e,n,s)=>{s.d(n,{Z:()=>a,a:()=>t});var i=s(7294);const r={},o=i.createContext(r);function t(e){const n=i.useContext(o);return i.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function a(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(r):e.components||r:t(e.components),i.createElement(o.Provider,{value:n},e.children)}}}]);
"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[7731],{5616:(e,n,r)=>{r.r(n),r.d(n,{assets:()=>a,contentTitle:()=>o,default:()=>u,frontMatter:()=>c,metadata:()=>t,toc:()=>d});var s=r(5893),i=r(1151);const c={title:"Configure AWS ECR Resource Scanner for TerediX",sidebar_label:"AWS ECR"},o="AWS ECR",t={id:"configuration/scanner/aws_ecr",title:"Configure AWS ECR Resource Scanner for TerediX",description:"Configuration",source:"@site/docs/configuration/scanner/aws_ecr.md",sourceDirName:"configuration/scanner",slug:"/configuration/scanner/aws_ecr",permalink:"/docs/configuration/scanner/aws_ecr",draft:!1,unlisted:!1,editUrl:"https://github.com/shaharia-lab/teredix/tree/master/website/docs/configuration/scanner/aws_ecr.md",tags:[],version:"current",lastUpdatedAt:1705842729,formattedLastUpdatedAt:"Jan 21, 2024",frontMatter:{title:"Configure AWS ECR Resource Scanner for TerediX",sidebar_label:"AWS ECR"},sidebar:"tutorialSidebar",previous:{title:"AWS EC2",permalink:"/docs/configuration/scanner/aws_ec2"},next:{title:"AWS RDS",permalink:"/docs/configuration/scanner/aws_rds"}},a={},d=[{value:"Configuration",id:"configuration",level:2},{value:"Type",id:"type",level:3},{value:"Configuration",id:"configuration-1",level:3},{value:"Fields",id:"fields",level:3},{value:"Schedule",id:"schedule",level:3},{value:"Example",id:"example",level:2}];function l(e){const n={a:"a",code:"code",h1:"h1",h2:"h2",h3:"h3",li:"li",p:"p",pre:"pre",strong:"strong",ul:"ul",...(0,i.a)(),...e.components};return(0,s.jsxs)(s.Fragment,{children:[(0,s.jsx)(n.h1,{id:"aws-ecr",children:"AWS ECR"}),"\n",(0,s.jsx)("img",{src:"/img/aws_ecr_icon.png",alt:"AWS ECR",width:"250"}),"\n",(0,s.jsx)(n.h2,{id:"configuration",children:"Configuration"}),"\n",(0,s.jsx)(n.h3,{id:"type",children:"Type"}),"\n",(0,s.jsxs)(n.p,{children:["Resource type. In this case it would be ",(0,s.jsx)(n.code,{children:"aws_ecr"}),"."]}),"\n",(0,s.jsx)(n.h3,{id:"configuration-1",children:"Configuration"}),"\n",(0,s.jsxs)(n.ul,{children:["\n",(0,s.jsxs)(n.li,{children:[(0,s.jsx)(n.strong,{children:"access_key"}),": AWS access key"]}),"\n",(0,s.jsxs)(n.li,{children:[(0,s.jsx)(n.strong,{children:"secret_key"}),": AWS secret key"]}),"\n",(0,s.jsxs)(n.li,{children:[(0,s.jsx)(n.strong,{children:"region"}),": AWS region. e.g: us-west-1"]}),"\n",(0,s.jsxs)(n.li,{children:[(0,s.jsx)(n.strong,{children:"session_token"}),": AWS session token"]}),"\n",(0,s.jsxs)(n.li,{children:[(0,s.jsx)(n.strong,{children:"account_id"}),": AWS account ID"]}),"\n"]}),"\n",(0,s.jsx)(n.h3,{id:"fields",children:"Fields"}),"\n",(0,s.jsx)(n.p,{children:"List of available fields to add with resource add metadata. During scanning resources, scanner will only fetch data\nfrom the following fields."}),"\n",(0,s.jsxs)(n.ul,{children:["\n",(0,s.jsx)(n.li,{children:"repository_name"}),"\n",(0,s.jsx)(n.li,{children:"repository_uri"}),"\n",(0,s.jsx)(n.li,{children:"registry_id"}),"\n",(0,s.jsx)(n.li,{children:"arn"}),"\n",(0,s.jsx)(n.li,{children:"tags"}),"\n"]}),"\n",(0,s.jsx)(n.h3,{id:"schedule",children:"Schedule"}),"\n",(0,s.jsxs)(n.p,{children:[(0,s.jsxs)(n.strong,{children:["\ud83d\udd17 ",(0,s.jsx)(n.a,{href:"/docs/configuration/scanner/overview#schedule-format",children:"Check schedule format"})]}),"."]}),"\n",(0,s.jsx)(n.h2,{id:"example",children:"Example"}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-yaml",children:'source:\n  aws_ecr_source_one:\n      type: aws_ecr\n      configuration:\n        access_key: "xxxx"\n        secret_key: "xxxx"\n        session_token: "xxxx"\n        region: "x"\n        account_id: "xxx"\n      fields:\n        - repository_name\n        - repository_uri\n        - registry_id\n        - arn\n        - tags\n      schedule: "@every 24h"\n'})}),"\n",(0,s.jsxs)(n.p,{children:["In the above example, we have added a source named ",(0,s.jsx)(n.code,{children:"aws_ecr_source_one"})," with type ",(0,s.jsx)(n.code,{children:"aws_ecr"}),". We have added some fields to add with each resource.\nWe have also set the schedule to run this source every 24 hours."]}),"\n",(0,s.jsxs)(n.p,{children:["Based on the above example, scanner_name would be ",(0,s.jsx)(n.code,{children:"aws_ecr_source_one"})," and scanner_type would be ",(0,s.jsx)(n.code,{children:"aws_ecr"}),". This is\nimportant to filter resources in Grafana dashboard."]})]})}function u(e={}){const{wrapper:n}={...(0,i.a)(),...e.components};return n?(0,s.jsx)(n,{...e,children:(0,s.jsx)(l,{...e})}):l(e)}},1151:(e,n,r)=>{r.d(n,{Z:()=>t,a:()=>o});var s=r(7294);const i={},c=s.createContext(i);function o(e){const n=s.useContext(c);return s.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function t(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(i):e.components||i:o(e.components),s.createElement(c.Provider,{value:n},e.children)}}}]);
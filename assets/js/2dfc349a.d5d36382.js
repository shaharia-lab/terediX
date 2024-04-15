"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[6404],{80309:(e,n,s)=>{s.r(n),s.d(n,{assets:()=>o,contentTitle:()=>t,default:()=>h,frontMatter:()=>r,metadata:()=>a,toc:()=>l});var i=s(85893),c=s(11151);const r={title:"Configure AWS EC2 Resource Scanner for TerediX",sidebar_label:"AWS EC2"},t="AWS EC2",a={id:"configuration/scanner/aws_ec2",title:"Configure AWS EC2 Resource Scanner for TerediX",description:"Configuration",source:"@site/docs/configuration/scanner/aws_ec2.md",sourceDirName:"configuration/scanner",slug:"/configuration/scanner/aws_ec2",permalink:"/docs/configuration/scanner/aws_ec2",draft:!1,unlisted:!1,editUrl:"https://github.com/shaharia-lab/teredix/tree/master/website/docs/configuration/scanner/aws_ec2.md",tags:[],version:"current",lastUpdatedAt:1713216937e3,frontMatter:{title:"Configure AWS EC2 Resource Scanner for TerediX",sidebar_label:"AWS EC2"},sidebar:"tutorialSidebar",previous:{title:"Overview",permalink:"/docs/configuration/scanner/overview"},next:{title:"AWS ECR",permalink:"/docs/configuration/scanner/aws_ecr"}},o={},l=[{value:"Configuration",id:"configuration",level:2},{value:"Type",id:"type",level:3},{value:"Configuration",id:"configuration-1",level:3},{value:"Fields",id:"fields",level:3},{value:"Schedule",id:"schedule",level:3},{value:"Example",id:"example",level:2}];function d(e){const n={a:"a",code:"code",h1:"h1",h2:"h2",h3:"h3",li:"li",p:"p",pre:"pre",strong:"strong",ul:"ul",...(0,c.a)(),...e.components};return(0,i.jsxs)(i.Fragment,{children:[(0,i.jsx)(n.h1,{id:"aws-ec2",children:"AWS EC2"}),"\n",(0,i.jsx)("img",{src:"/img/aws_ec2_icon.png",alt:"AWS EC2",width:"250"}),"\n",(0,i.jsx)(n.h2,{id:"configuration",children:"Configuration"}),"\n",(0,i.jsx)(n.h3,{id:"type",children:"Type"}),"\n",(0,i.jsxs)(n.p,{children:["Resource type. In this case it would be ",(0,i.jsx)(n.code,{children:"aws_ec2"}),"."]}),"\n",(0,i.jsx)(n.h3,{id:"configuration-1",children:"Configuration"}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.strong,{children:"access_key"}),": AWS access key"]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.strong,{children:"secret_key"}),": AWS secret key"]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.strong,{children:"region"}),": AWS region. e.g: us-west-1"]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.strong,{children:"session_token"}),": AWS session token"]}),"\n",(0,i.jsxs)(n.li,{children:[(0,i.jsx)(n.strong,{children:"account_id"}),": AWS account ID"]}),"\n"]}),"\n",(0,i.jsx)(n.h3,{id:"fields",children:"Fields"}),"\n",(0,i.jsx)(n.p,{children:"List of available fields to add with resource add metadata. During scanning resources, scanner will only fetch data\nfrom the following fields."}),"\n",(0,i.jsxs)(n.ul,{children:["\n",(0,i.jsx)(n.li,{children:"instance_id"}),"\n",(0,i.jsx)(n.li,{children:"image_id"}),"\n",(0,i.jsx)(n.li,{children:"private_dns_name"}),"\n",(0,i.jsx)(n.li,{children:"instance_type"}),"\n",(0,i.jsx)(n.li,{children:"architecture"}),"\n",(0,i.jsx)(n.li,{children:"instance_lifecycle"}),"\n",(0,i.jsx)(n.li,{children:"instance_state"}),"\n",(0,i.jsx)(n.li,{children:"vpc_id"}),"\n",(0,i.jsx)(n.li,{children:"tags"}),"\n"]}),"\n",(0,i.jsx)(n.h3,{id:"schedule",children:"Schedule"}),"\n",(0,i.jsxs)(n.p,{children:[(0,i.jsxs)(n.strong,{children:["\ud83d\udd17 ",(0,i.jsx)(n.a,{href:"/docs/configuration/scanner/overview#schedule-format",children:"Check schedule format"})]}),"."]}),"\n",(0,i.jsx)(n.h2,{id:"example",children:"Example"}),"\n",(0,i.jsx)(n.pre,{children:(0,i.jsx)(n.code,{className:"language-yaml",children:'source:\n  aws_ec2_source_one:\n      type: aws_ec2\n      configuration:\n        access_key: "xxxx"\n        secret_key: "xxxx"\n        session_token: "xxxx"\n        region: "x"\n        account_id: "xxx"\n      fields:\n        - instance_id\n        - image_id\n        - private_dns_name\n        - instance_type\n        - architecture\n        - instance_lifecycle\n        - instance_state\n        - vpc_id\n        - tags\n      schedule: "@every 24h"\n'})}),"\n",(0,i.jsxs)(n.p,{children:["In the above example, we have added a source named ",(0,i.jsx)(n.code,{children:"aws_ec2_source_one"})," with type ",(0,i.jsx)(n.code,{children:"aws_ec2"}),". We have added some fields to add with each resource.\nWe have also set the schedule to run this source every 24 hours."]}),"\n",(0,i.jsxs)(n.p,{children:["Based on the above example, scanner_name would be ",(0,i.jsx)(n.code,{children:"aws_ec2_source_one"})," and scanner_type would be ",(0,i.jsx)(n.code,{children:"aws_ec2"}),". This is\nimportant to filter resources in Grafana dashboard."]})]})}function h(e={}){const{wrapper:n}={...(0,c.a)(),...e.components};return n?(0,i.jsx)(n,{...e,children:(0,i.jsx)(d,{...e})}):d(e)}},11151:(e,n,s)=>{s.d(n,{Z:()=>a,a:()=>t});var i=s(67294);const c={},r=i.createContext(c);function t(e){const n=i.useContext(r);return i.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function a(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(c):e.components||c:t(e.components),i.createElement(r.Provider,{value:n},e.children)}}}]);
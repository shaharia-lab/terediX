"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[4817],{54268:(e,n,t)=>{t.r(n),t.d(n,{assets:()=>c,contentTitle:()=>s,default:()=>h,frontMatter:()=>i,metadata:()=>a,toc:()=>l});var r=t(85893),o=t(11151);const i={sidebar_position:2,title:"Prepare Local Environment"},s=void 0,a={id:"contribution/prepare-local-env",title:"Prepare Local Environment",description:"We have made it easy to get started working on TerediX codebase. You can follow the steps below to get started.",source:"@site/docs/contribution/prepare-local-env.md",sourceDirName:"contribution",slug:"/contribution/prepare-local-env",permalink:"/docs/contribution/prepare-local-env",draft:!1,unlisted:!1,editUrl:"https://github.com/shaharia-lab/teredix/tree/master/website/docs/contribution/prepare-local-env.md",tags:[],version:"current",lastUpdatedAt:1713174558e3,sidebarPosition:2,frontMatter:{sidebar_position:2,title:"Prepare Local Environment"},sidebar:"tutorialSidebar",previous:{title:"Overview",permalink:"/docs/contribution/overview"},next:{title:"Helm Chart for Local Development",permalink:"/docs/contribution/helm-chart-for-local-dev"}},c={},l=[{value:"Prerequisites",id:"prerequisites",level:2},{value:"Clone the repository",id:"clone-the-repository",level:2},{value:"Run TerediX inside Docker",id:"run-teredix-inside-docker",level:2},{value:"Access the development environment",id:"access-the-development-environment",level:2},{value:"Test the development environment",id:"test-the-development-environment",level:2},{value:"Run website",id:"run-website",level:2}];function d(e){const n={a:"a",code:"code",h2:"h2",li:"li",p:"p",pre:"pre",ul:"ul",...(0,o.a)(),...e.components};return(0,r.jsxs)(r.Fragment,{children:[(0,r.jsx)(n.p,{children:"We have made it easy to get started working on TerediX codebase. You can follow the steps below to get started."}),"\n",(0,r.jsx)(n.h2,{id:"prerequisites",children:"Prerequisites"}),"\n",(0,r.jsxs)(n.ul,{children:["\n",(0,r.jsx)(n.li,{children:"Docker"}),"\n"]}),"\n",(0,r.jsx)(n.h2,{id:"clone-the-repository",children:"Clone the repository"}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-bash",children:"git clone git@github.com:shaharia-lab/terediX.git\ncd terediX\n"})}),"\n",(0,r.jsx)(n.h2,{id:"run-teredix-inside-docker",children:"Run TerediX inside Docker"}),"\n",(0,r.jsx)(n.p,{children:"We have a Docker image for development purpose. You can run the following command to start the development server."}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-bash",children:"docker run -i -d --name teredix-dev \\\n  -v $(pwd):/home/app/src \\\n  -p 3000:3000 \\\n  -p 2112:2112 \\\n  ghcr.io/shaharia-lab/teredix:dev\n"})}),"\n",(0,r.jsx)(n.h2,{id:"access-the-development-environment",children:"Access the development environment"}),"\n",(0,r.jsx)(n.p,{children:"You can access development environment in Docker container by running the following command:"}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-bash",children:"docker exec -it teredix-dev bash\n"})}),"\n",(0,r.jsx)(n.p,{children:"Then, you can go to the project root"}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-bash",children:"su app\ncd ~/src/\n"})}),"\n",(0,r.jsx)(n.h2,{id:"test-the-development-environment",children:"Test the development environment"}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-bash",children:"make test-unit\n"})}),"\n",(0,r.jsx)(n.h2,{id:"run-website",children:"Run website"}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-bash",children:"cd website\nyarn install\nyarn start --host=0.0.0.0 --port=3000\n"})}),"\n",(0,r.jsxs)(n.p,{children:["Now in your browser, you can access the website at ",(0,r.jsx)(n.a,{href:"http://localhost:3000",children:"http://localhost:3000"})]}),"\n",(0,r.jsx)(n.p,{children:"Voila! You are ready to contribute to TerediX."})]})}function h(e={}){const{wrapper:n}={...(0,o.a)(),...e.components};return n?(0,r.jsx)(n,{...e,children:(0,r.jsx)(d,{...e})}):d(e)}},11151:(e,n,t)=>{t.d(n,{Z:()=>a,a:()=>s});var r=t(67294);const o={},i=r.createContext(o);function s(e){const n=r.useContext(i);return r.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function a(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(o):e.components||o:s(e.components),r.createElement(i.Provider,{value:n},e.children)}}}]);
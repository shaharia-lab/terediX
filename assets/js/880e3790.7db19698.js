"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[7791],{37358:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>i,contentTitle:()=>o,default:()=>h,frontMatter:()=>s,metadata:()=>l,toc:()=>c});var r=n(85893),a=n(11151);const s={sidebar_position:3,title:"Helm Chart for Local Development"},o=void 0,l={id:"contribution/helm-chart-for-local-dev",title:"Helm Chart for Local Development",description:"We have made it easy to test deploying terediX in your local Kubernetes cluster. You can follow the steps below to get started.",source:"@site/docs/contribution/helm-chart-for-local-dev.md",sourceDirName:"contribution",slug:"/contribution/helm-chart-for-local-dev",permalink:"/docs/contribution/helm-chart-for-local-dev",draft:!1,unlisted:!1,editUrl:"https://github.com/shaharia-lab/teredix/tree/master/website/docs/contribution/helm-chart-for-local-dev.md",tags:[],version:"current",lastUpdatedAt:1713134108,formattedLastUpdatedAt:"Apr 14, 2024",sidebarPosition:3,frontMatter:{sidebar_position:3,title:"Helm Chart for Local Development"},sidebar:"tutorialSidebar",previous:{title:"Prepare Local Environment",permalink:"/docs/contribution/prepare-local-env"}},i={},c=[{value:"Prerequisites",id:"prerequisites",level:2},{value:"Create a namespace",id:"create-a-namespace",level:2},{value:"Install PostgreSQL",id:"install-postgresql",level:2},{value:"Install terediX helm chart",id:"install-teredix-helm-chart",level:2}];function d(e){const t={a:"a",code:"code",h2:"h2",li:"li",p:"p",pre:"pre",ul:"ul",...(0,a.a)(),...e.components};return(0,r.jsxs)(r.Fragment,{children:[(0,r.jsx)(t.p,{children:"We have made it easy to test deploying terediX in your local Kubernetes cluster. You can follow the steps below to get started."}),"\n",(0,r.jsx)(t.h2,{id:"prerequisites",children:"Prerequisites"}),"\n",(0,r.jsxs)(t.ul,{children:["\n",(0,r.jsx)(t.li,{children:"Docker"}),"\n",(0,r.jsxs)(t.li,{children:["Kubernetes Cluster (Minikube, Kind, K3s, etc). You can easily start a local Kubernetes cluster using ",(0,r.jsx)(t.a,{href:"https://github.com/shaharia-lab/k8s-dev-cluster",children:"KiND"})]}),"\n"]}),"\n",(0,r.jsx)(t.p,{children:"Go to the next step when your Kubernetes cluster is ready."}),"\n",(0,r.jsx)(t.h2,{id:"create-a-namespace",children:"Create a namespace"}),"\n",(0,r.jsx)(t.p,{children:"Create a namespace for terediX"}),"\n",(0,r.jsx)(t.pre,{children:(0,r.jsx)(t.code,{className:"language-bash",children:"kubectl create namespace teredix\n"})}),"\n",(0,r.jsx)(t.h2,{id:"install-postgresql",children:"Install PostgreSQL"}),"\n",(0,r.jsxs)(t.p,{children:["Because terediX need a storage solution and currently ",(0,r.jsx)(t.a,{href:"/docs/configuration/storage#supported-storage-engines",children:"we support"})," only PostgreSQL,\nso you need to install PostgreSQL in your Kubernetes cluster. You can install PostgreSQL by helm chart."]}),"\n",(0,r.jsx)(t.pre,{children:(0,r.jsx)(t.code,{className:"language-bash",children:'helm repo add bitnami https://charts.bitnami.com/bitnami \nhelm repo update\nhelm upgrade --install postgresql bitnami/postgresql --namespace "teredix" \\\n        --set auth.username="app" \\\n        --set auth.password="pass" \\\n        --set auth.database="app"\n'})}),"\n",(0,r.jsx)(t.h2,{id:"install-teredix-helm-chart",children:"Install terediX helm chart"}),"\n",(0,r.jsxs)(t.p,{children:["Create a local values file in ",(0,r.jsx)(t.code,{children:"helm-chart/teredix/values-local.yaml"})," for terediX helm chart to override few values for local development."]}),"\n",(0,r.jsx)(t.pre,{children:(0,r.jsx)(t.code,{className:"language-yaml",children:'# helm-chart/teredix/values-local.yaml\nimage:\n  repository: teredix\n  tag: "prod"\n\nappConfig:\n  organization:\n    name: Your Organization\n    logo: https://your-org-url.com/logo.png\n  discovery:\n    name: Name of the discovery\n    description: Some description about the discovery\n    worker_pool_size: 1\n  storage:\n    batch_size: 2\n    engines:\n      postgresql:\n        host: "postgresql"\n        port: 5432\n        user: "app"\n        password: "pass"\n        db: "app"\n    default_engine: postgresql\n  source:\n    fs_one:\n      type: file_system\n      configuration:\n        root_directory: "/config"\n      fields:\n        - machineHost\n        - rootDirectory\n      schedule: "@every 300s"\n  relations:\n    criteria:\n      - name: "file-system-rule1"\n        source:\n          kind: "FilePath"\n          meta_key: "rootDirectory"\n          meta_value: "/some/path"\n        target:\n          kind: "FilePath"\n          meta_key: "rootDirectory"\n          meta_value: "/some/path"\n\nservice:\n  type: ClusterIP\n  port: 2112\n\ningress:\n  enabled: true\n  hosts:\n    - host: teredix.dev.local\n      paths:\n        - path: /\n          pathType: ImplementationSpecific\n'})}),"\n",(0,r.jsx)(t.p,{children:"Now install terediX helm chart using the following command:"}),"\n",(0,r.jsx)(t.pre,{children:(0,r.jsx)(t.code,{className:"language-bash",children:"helm upgrade --install teredix ./helm-chart/teredix --namespace teredix \\\n        -f ./helm-chart/teredix/values.yaml \\\n        -f ./helm-chart/teredix/values-local.yaml\n"})})]})}function h(e={}){const{wrapper:t}={...(0,a.a)(),...e.components};return t?(0,r.jsx)(t,{...e,children:(0,r.jsx)(d,{...e})}):d(e)}},11151:(e,t,n)=>{n.d(t,{Z:()=>l,a:()=>o});var r=n(67294);const a={},s=r.createContext(a);function o(e){const t=r.useContext(s);return r.useMemo((function(){return"function"==typeof e?e(t):{...t,...e}}),[t,e])}function l(e){let t;return t=e.disableParentContext?"function"==typeof e.components?e.components(a):e.components||a:o(e.components),r.createElement(s.Provider,{value:t},e.children)}}}]);
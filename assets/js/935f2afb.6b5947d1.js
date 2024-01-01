"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[53],{1109:e=>{e.exports=JSON.parse('{"pluginId":"default","version":"current","label":"Next","banner":null,"badge":false,"noIndex":false,"className":"docs-version-current","isLast":true,"docsSidebars":{"tutorialSidebar":[{"type":"link","label":"Introduction","href":"/docs/intro","docId":"intro","unlisted":false},{"type":"link","label":"Technical Architecture","href":"/docs/technical_architecture","docId":"technical_architecture","unlisted":false},{"type":"link","label":"Getting Started","href":"/docs/getting_started","docId":"getting_started","unlisted":false},{"type":"category","label":"Configuration","collapsible":true,"collapsed":true,"items":[{"type":"link","label":"Overview","href":"/docs/configuration/general","docId":"configuration/general","unlisted":false},{"type":"link","label":"Organization","href":"/docs/configuration/organization","docId":"configuration/organization","unlisted":false},{"type":"link","label":"Discovery","href":"/docs/configuration/discovery","docId":"configuration/discovery","unlisted":false},{"type":"category","label":"Storage","collapsible":true,"collapsed":true,"items":[{"type":"link","label":"Basic Configuration","href":"/docs/configuration/storage/","docId":"configuration/storage/index","unlisted":false},{"type":"link","label":"PostgreSQL","href":"/docs/configuration/storage/postgresql","docId":"configuration/storage/postgresql","unlisted":false}],"href":"/docs/category/storage"},{"type":"category","label":"Scanner","collapsible":true,"collapsed":true,"items":[{"type":"link","label":"Overview","href":"/docs/configuration/scanner/overview","docId":"configuration/scanner/overview","unlisted":false},{"type":"link","label":"AWS EC2","href":"/docs/configuration/scanner/aws_ec2","docId":"configuration/scanner/aws_ec2","unlisted":false},{"type":"link","label":"AWS ECR","href":"/docs/configuration/scanner/aws_ecr","docId":"configuration/scanner/aws_ecr","unlisted":false},{"type":"link","label":"AWS RDS","href":"/docs/configuration/scanner/aws_rds","docId":"configuration/scanner/aws_rds","unlisted":false},{"type":"link","label":"AWS S3","href":"/docs/configuration/scanner/aws_s3","docId":"configuration/scanner/aws_s3","unlisted":false},{"type":"link","label":"File System","href":"/docs/configuration/scanner/file_system","docId":"configuration/scanner/file_system","unlisted":false},{"type":"link","label":"GitHub Repository","href":"/docs/configuration/scanner/github_repository","docId":"configuration/scanner/github_repository","unlisted":false}],"href":"/docs/category/scanner"}],"href":"/docs/category/configuration"},{"type":"category","label":"Installation","collapsible":true,"collapsed":true,"items":[{"type":"link","label":"Docker","href":"/docs/installation/docker","docId":"installation/docker","unlisted":false},{"type":"link","label":"Go Install","href":"/docs/installation/go_install","docId":"installation/go_install","unlisted":false},{"type":"link","label":"Kubernetes","href":"/docs/installation/kubernetes","docId":"installation/kubernetes","unlisted":false},{"type":"link","label":"Standalone Binary","href":"/docs/installation/standalone_binary","docId":"installation/standalone_binary","unlisted":false}],"href":"/docs/category/installation"},{"type":"category","label":"Contribution","collapsible":true,"collapsed":true,"items":[{"type":"link","label":"Overview","href":"/docs/contribution/overview","docId":"contribution/overview","unlisted":false},{"type":"link","label":"Prepare Local Environment","href":"/docs/contribution/prepare-local-env","docId":"contribution/prepare-local-env","unlisted":false},{"type":"link","label":"Helm Chart for Local Development","href":"/docs/contribution/helm-chart-for-local-dev","docId":"contribution/helm-chart-for-local-dev","unlisted":false}],"href":"/docs/category/contribution"}]},"docs":{"configuration/discovery":{"id":"configuration/discovery","title":"Discovery","description":"The entire terediX process is called a discovery.","sidebar":"tutorialSidebar"},"configuration/general":{"id":"configuration/general","title":"Overview","description":"terediX uses a configuration file to run. You can create a configuration file with the following command:","sidebar":"tutorialSidebar"},"configuration/organization":{"id":"configuration/organization","title":"Organization","description":"For reporting and visualization, you can add your organization details in the configuration file.","sidebar":"tutorialSidebar"},"configuration/scanner/aws_ec2":{"id":"configuration/scanner/aws_ec2","title":"Configure AWS EC2 Resource Scanner for TerediX","description":"Configuration","sidebar":"tutorialSidebar"},"configuration/scanner/aws_ecr":{"id":"configuration/scanner/aws_ecr","title":"Configure AWS ECR Resource Scanner for TerediX","description":"Configuration","sidebar":"tutorialSidebar"},"configuration/scanner/aws_rds":{"id":"configuration/scanner/aws_rds","title":"Configure AWS RDS Resource Scanner for TerediX","description":"Configuration","sidebar":"tutorialSidebar"},"configuration/scanner/aws_s3":{"id":"configuration/scanner/aws_s3","title":"Configure AWS S3 Resource Scanner for TerediX","description":"Configuration","sidebar":"tutorialSidebar"},"configuration/scanner/file_system":{"id":"configuration/scanner/file_system","title":"Configure File System Resource Scanner for TerediX","description":"Configuration","sidebar":"tutorialSidebar"},"configuration/scanner/github_repository":{"id":"configuration/scanner/github_repository","title":"Configure GitHub Repository Resource Scanner for TerediX","description":"Configuration","sidebar":"tutorialSidebar"},"configuration/scanner/overview":{"id":"configuration/scanner/overview","title":"Available Scanners and Configurations","description":"Source is the place where terediX will discover the data. You can add multiple sources in the configuration file.","sidebar":"tutorialSidebar"},"configuration/storage/index":{"id":"configuration/storage/index","title":"Basic Configuration","description":"Storage is the place where terediX will store the discovered data. You can add multiple storage engines in the","sidebar":"tutorialSidebar"},"configuration/storage/postgresql":{"id":"configuration/storage/postgresql","title":"PostgreSQL","description":"Here is the configuration for PostgreSQL storage engine.","sidebar":"tutorialSidebar"},"contribution/helm-chart-for-local-dev":{"id":"contribution/helm-chart-for-local-dev","title":"Helm Chart for Local Development","description":"We have made it easy to test deploying terediX in your local Kubernetes cluster. You can follow the steps below to get started.","sidebar":"tutorialSidebar"},"contribution/overview":{"id":"contribution/overview","title":"Overview","description":"Contributing to TerediX is a great way to help the open source community. We appreciate contributions of all kinds, from code to documentation, from small tweaks to significant features. Here are some of the ways you can contribute:","sidebar":"tutorialSidebar"},"contribution/prepare-local-env":{"id":"contribution/prepare-local-env","title":"Prepare Local Environment","description":"We have made it easy to get started working on TerediX codebase. You can follow the steps below to get started.","sidebar":"tutorialSidebar"},"getting_started":{"id":"getting_started","title":"Getting Started","description":"Getting Started","sidebar":"tutorialSidebar"},"installation/docker":{"id":"installation/docker","title":"Run terediX using Docker","description":"Docker","sidebar":"tutorialSidebar"},"installation/go_install":{"id":"installation/go_install","title":"Install terediX with Go Install","description":"Go Install","sidebar":"tutorialSidebar"},"installation/kubernetes":{"id":"installation/kubernetes","title":"Deploy terediX in Kubernetes","description":"Deploy in Kubernetes using Helm Chart","sidebar":"tutorialSidebar"},"installation/standalone_binary":{"id":"installation/standalone_binary","title":"Use terediX as standalone binary","description":"Standalone Binary","sidebar":"tutorialSidebar"},"intro":{"id":"intro","title":"Introduction","description":"TeReDiX (Tech Resource Discover &amp; Exploration) is a tool to discover tech resources for an organization from","sidebar":"tutorialSidebar"},"technical_architecture":{"id":"technical_architecture","title":"Technical Architecture","description":"Components","sidebar":"tutorialSidebar"}}}')}}]);
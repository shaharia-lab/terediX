// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require('prism-react-renderer').themes.github;
const darkCodeTheme = require('prism-react-renderer').themes.dracula;

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'TerediX',
  tagline: 'Tech Resource Discovery - an open source initiative by Shaharia Lab OÜ',
  favicon: 'img/favicon.ico',

  // Set the production url of your site here
  url: 'https://teredix.shaharialab.com',
  // Set the /<baseUrl>/ pathname under which your site is served
  // For GitHub pages deployment, it is often '/<projectName>/'
  baseUrl: '/',

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: 'shaharia-lab', // Usually your GitHub org/user name.
  projectName: 'teredix', // Usually your repo name.
  deploymentBranch: "gh-pages",
  trailingSlash: false,

  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',

  // Even if you don't use internalization, you can use this field to set useful
  // metadata like html lang. For example, if your site is Chinese, you may want
  // to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          breadcrumbs: true,
          showLastUpdateTime: true,
          editUrl: 'https://github.com/shaharia-lab/teredix/tree/master/website',
        },
        blog: {
          showReadingTime: true,
          // Please change this to your repo.
          // Remove this to remove the "edit this page" links.
          editUrl: 'https://github.com/shaharia-lab/teredix/tree/master/website',
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      }),
    ],
  ],

  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      // Replace with your project's social card
      image: 'img/docusaurus-social-card.jpg',
      navbar: {
        title: 'terediX',
        logo: {
          alt: 'terediX - Tech Resource Discovery',
          src: 'img/teredix_logo.png',
        },
        items: [
          {
            type: 'docSidebar',
            sidebarId: 'tutorialSidebar',
            position: 'left',
            label: 'Documentations',
          },
          {
            to: 'https://github.com/shaharia-lab/terediX/releases',
            label: 'Releases',
            position: 'left'
          },
          {
            to: 'https://github.com/sponsors/shaharia-lab',
            label: 'Sponsor',
            position: 'left'
          },
          {
            href: 'https://github.com/shaharia-lab/teredix',
            label: 'GitHub',
            position: 'right',
          },
        ],
      },
      footer: {
        style: 'dark',
        links: [
          {
            title: 'Docs',
            items: [
              {
                label: 'Documentations',
                to: '/docs/intro',
              },
            ],
          },
          {
            title: 'Community',
            items: [
              {
                label: 'Stack Overflow',
                href: 'https://stackoverflow.com/questions/tagged/teredix',
              },
              {
                label: 'Discord',
                href: 'https://discordapp.com/invite/teredix',
              },
            ],
          },
          {
            title: 'More',
            items: [
              {
                label: 'Releases',
                to: 'https://github.com/shaharia-lab/terediX/releases',
              },
              {
                label: 'GitHub',
                href: 'https://github.com/shaharia-lab/teredix',
              },
            ],
          },
        ],
        copyright: `Copyright © ${new Date().getFullYear()} Shaharia Lab OÜ. Built with Docusaurus.`,
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
      },
      tableOfContents: {
        minHeadingLevel: 2,
        maxHeadingLevel: 5,
      },
    }),
};

module.exports = config;

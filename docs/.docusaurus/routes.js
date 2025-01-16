import React from 'react';
import ComponentCreator from '@docusaurus/ComponentCreator';

export default [
  {
    path: '/pocket-id/__docusaurus/debug',
    component: ComponentCreator('/pocket-id/__docusaurus/debug', '31b'),
    exact: true
  },
  {
    path: '/pocket-id/__docusaurus/debug/config',
    component: ComponentCreator('/pocket-id/__docusaurus/debug/config', '7c4'),
    exact: true
  },
  {
    path: '/pocket-id/__docusaurus/debug/content',
    component: ComponentCreator('/pocket-id/__docusaurus/debug/content', '751'),
    exact: true
  },
  {
    path: '/pocket-id/__docusaurus/debug/globalData',
    component: ComponentCreator('/pocket-id/__docusaurus/debug/globalData', '402'),
    exact: true
  },
  {
    path: '/pocket-id/__docusaurus/debug/metadata',
    component: ComponentCreator('/pocket-id/__docusaurus/debug/metadata', 'dce'),
    exact: true
  },
  {
    path: '/pocket-id/__docusaurus/debug/registry',
    component: ComponentCreator('/pocket-id/__docusaurus/debug/registry', '520'),
    exact: true
  },
  {
    path: '/pocket-id/__docusaurus/debug/routes',
    component: ComponentCreator('/pocket-id/__docusaurus/debug/routes', '758'),
    exact: true
  },
  {
    path: '/pocket-id/',
    component: ComponentCreator('/pocket-id/', '5d3'),
    exact: true
  },
  {
    path: '/pocket-id/',
    component: ComponentCreator('/pocket-id/', 'e51'),
    routes: [
      {
        path: '/pocket-id/',
        component: ComponentCreator('/pocket-id/', '164'),
        routes: [
          {
            path: '/pocket-id/',
            component: ComponentCreator('/pocket-id/', 'c2f'),
            routes: [
              {
                path: '/pocket-id/examples/jellyfin',
                component: ComponentCreator('/pocket-id/examples/jellyfin', '04d'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/pocket-id/examples/oauth2-proxy',
                component: ComponentCreator('/pocket-id/examples/oauth2-proxy', 'ff0'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/pocket-id/help-out/contribute',
                component: ComponentCreator('/pocket-id/help-out/contribute', 'aea'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/pocket-id/introduction',
                component: ComponentCreator('/pocket-id/introduction', '519'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/pocket-id/setup/configuration',
                component: ComponentCreator('/pocket-id/setup/configuration', '024'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/pocket-id/setup/installation',
                component: ComponentCreator('/pocket-id/setup/installation', 'b47'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/pocket-id/setup/nginx-reverse-proxy',
                component: ComponentCreator('/pocket-id/setup/nginx-reverse-proxy', '2f4'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/pocket-id/setup/upgrading',
                component: ComponentCreator('/pocket-id/setup/upgrading', '9ea'),
                exact: true,
                sidebar: "docsSidebar"
              },
              {
                path: '/pocket-id/troubleshooting/account-recovery',
                component: ComponentCreator('/pocket-id/troubleshooting/account-recovery', '788'),
                exact: true,
                sidebar: "docsSidebar"
              }
            ]
          }
        ]
      }
    ]
  },
  {
    path: '*',
    component: ComponentCreator('*'),
  },
];

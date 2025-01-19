import type { SidebarsConfig } from "@docusaurus/plugin-content-docs";

const sidebars: SidebarsConfig = {
  docsSidebar: [
    {
      type: "doc",
      id: "introduction",
    },
    {
      type: "category",
      label: "Getting Started",
      items: [
        {
          type: "doc",
          id: "setup/installation",
        },
        {
          type: "doc",
          id: "setup/nginx-reverse-proxy",
        },
        {
          type: "doc",
          id: "setup/upgrading",
        },
      ],
    },
    {
      type: "category",
      label: "Configuration",
      items: [
        {
          type: "doc",
          id: "configuration/environment-variables",
        },
        {
          type: "doc",
          id: "configuration/ldap",
        },
      ],
    },
    {
      type: "category",
      label: "Troubleshooting",
      items: [
        {
          type: "doc",
          id: "troubleshooting/account-recovery",
        },
      ],
    },
    {
      type: "category",
      label: "Examples",
      items: [
        {
          type: "doc",
          id: "examples/caddy",
        },
        {
          type: "doc",
          id: "examples/oauth2-proxy",
        },
        {
          type: "doc",
          id: "examples/jellyfin",
        },
      ],
    },
    {
      type: "category",
      label: "Helping Out",
      items: [
        {
          type: "doc",
          id: "help-out/contribute",
        },
      ],
    },
    {
      type: "link",
      label: "Demo",
      href: "https://pocket-id.eliasschneider.com/",
    },
  ],
};

export default sidebars;

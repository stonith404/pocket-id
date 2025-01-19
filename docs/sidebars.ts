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
          id: "setup/configuration",
        },
        {
          type: "doc",
          id: "setup/ldap",
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
    {
      type: "link",
      label: "Discord",
      href: "https://discord.gg/HutpbfB59Q",
    },
  ],
};

export default sidebars;
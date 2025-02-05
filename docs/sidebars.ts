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
          id: "setup/user-management",
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
          id: "configuration/allowed-groups",
        },
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
      label: "Guides",
      items: [
        {
          type: "doc",
          id: "guides/proxy-services",
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
        {
          type: "doc",
          id: "troubleshooting/common-issues",
        },
      ],
    },
    {
      type: "category",
      label: "Client Examples",
      link: {
        type: "generated-index",
        title: "Client Examples",
        description:
          "Examples of how to setup Pocket ID with different clients",
        slug: "client-examples",
      },
      items: [
        "client-examples/cloudflare-zero-trust",
        "client-examples/freshrss",
        "client-examples/gitea",
        "client-examples/grist",
        "client-examples/headscale",
        "client-examples/hoarder",
        "client-examples/immich",
        "client-examples/jellyfin",
        "client-examples/memos",
        "client-examples/netbox",
        "client-examples/open-webui",
        "client-examples/pgadmin",
        "client-examples/portainer",
        "client-examples/proxmox",
        "client-examples/semaphore-ui",
        "client-examples/vikunja",
      ],
    },
    {
      type: "link",
      label: "Contributing",
      href: "https://github.com/pocket-id/pocket-id/blob/main/CONTRIBUTING.md",
    },
    {
      type: "link",
      label: "Demo",
      href: "https://demo.pocket-id.org",
    },
  ],
};

export default sidebars;

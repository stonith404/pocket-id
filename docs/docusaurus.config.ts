import type * as Preset from "@docusaurus/preset-classic";
import type { Config } from "@docusaurus/types";
import { themes as prismThemes } from "prism-react-renderer";

const config: Config = {
  title: "Pocket ID",
  tagline:
    "Pocket ID is a simple OIDC provider that allows users to authenticate with their passkeys to your services.",
  favicon: "img/pocket-id.png",

  url: "https://stonith404.github.io",
  baseUrl: "/pocket-id/",
  organizationName: "stonith404",
  projectName: "pocket-id",

  onBrokenLinks: "warn",
  onBrokenMarkdownLinks: "warn",

  i18n: {
    defaultLocale: "en",
    locales: ["en"],
  },

  presets: [
    [
      "classic",
      {
        docs: {
          routeBasePath: "/",
          sidebarPath: "./sidebars.ts",
          editUrl: "https://github.com/stonith404/pocket-id/edit/main/docs",
        },
        blog: false,
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    image: "img/pocket-id.png",
    colorMode: {
      respectPrefersColorScheme: true,
    },
    navbar: {
      title: "Pocket ID",
      logo: {
        alt: "Pocket ID Share Logo",
        src: "img/pocket-id.png",
      },
      items: [
        {
          label: "v0.27.2",
          position: "right",
          to: "https://github.com/stonith404/pocket-id/releases/latest",
        },
        {
          href: "https://github.com/stonith404/pocket-id",
          label: "GitHub",
          position: "right",
        },
      ],
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
    },
  } satisfies Preset.ThemeConfig,
};

export default config;

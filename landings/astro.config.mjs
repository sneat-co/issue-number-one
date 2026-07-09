// @ts-check
import { defineConfig } from "astro/config";
import sitemap from "@astrojs/sitemap";

// https://astro.build/config
export default defineConfig({
  // Custom domain (issuenumber.one) attach is handled by the main session on the
  // issuenumber.one zone; the worker also serves on *.workers.dev meanwhile.
  site: "https://issuenumber.one",
  output: "static",
  outDir: "./dist",
  integrations: [sitemap()],
});

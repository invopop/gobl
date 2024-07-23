// vite.config.ts
import { resolve } from "path";
import { defineConfig } from "vite";
import dts from "vite-plugin-dts";

// https://vitejs.dev/guide/build.html#library-mode
export default defineConfig({
  plugins: [dts()],
  define: {
    GOBL_CLI_VERSION: JSON.stringify(process.env.npm_package_version),
  },
  build: {
    lib: {
      entry: resolve(__dirname, "src/index.ts"),
      name: "gobl-worker",
      fileName: "gobl-worker",
    },
  },
});

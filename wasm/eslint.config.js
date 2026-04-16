const js = require("@eslint/js");
const globals = require("globals");
const cypress = require("eslint-plugin-cypress/flat");
const prettier = require("eslint-plugin-prettier/recommended");

module.exports = [
  {
    ignores: ["wasm_exec.js"],
  },
  js.configs.recommended,
  prettier,
  {
    languageOptions: {
      ecmaVersion: "latest",
      sourceType: "module",
      globals: {
        ...globals.browser,
        ...globals.node,
      },
    },
    rules: {
      "no-unused-vars": "warn",
    },
  },
  {
    ...cypress.configs.recommended,
    files: ["cypress/**/*.js"],
    rules: {
      ...cypress.configs.recommended.rules,
      "cypress/no-unnecessary-waiting": "warn",
    },
  },
];

# GOBL.js (WASM)

This folder contains the files necessary to execute the GOBL client library in your web browser. The core library is written in Go, and compiled to WebAssembly.  The file `gobl.js` provides a thin JavaScript wrapper around the compiled WebAssembly running in a web worker.

To execute a simple demo of the GOBL library in your browser, you will need to:

1. [Install Go](https://go.dev/dl/) (1.17 or newer)
2. From this directory, execute the command `./build.sh` which will compile the WASM target and start a bare web server.
3. In your browser, navigate to `http://localhost:9999/`
4. Open the JavaScript console in your browser to see the test output.

---

## Running The GOBL.js Playground

### Development

To start the local server via browser-sync, including hot reloading (no browser refresh required to see changes):

```bash
npm run dev
```

To start the interactive cyprus.io browser test runner. You must have a server running (`npm run dev` or `npm run start` for prod). This interactive UI runs tests against whichever server you have running. (Note: the Lighthouse performance tests will fail against the development server, because performance is much slower there.):

```bash
npm run test
```

To run the testing suite as it would run on a ci server, against the production build:

```bash
npm run ci
```

To run the code formatting checking tools `eslint` (and `prettier`):

```bash
npm run check-formatting
```

### Production

For deployment, build the project and serve it on a simple http server.

```bash
npm start
```

## Packages

| Motivation | Package |
| -- | -- |
| development HTTP server, with hot reloading | [browser-sync](https://browsersync.io/) |
| production HTTP Server | [http-server](https://github.com/http-party/http-server) |
| Test runner | [Cypress](https://www.cypress.io/) (note: we're just using the free and open source runner, not the paid/hosted product) |
| Application performance tool | [Lighthouse](https://developers.google.com/web/tools/lighthouse) |
| Standardized Code Formatting (syntax gotchas, etc) | [ESLint](https://eslint.org/) |
| Standardized Code formatting (for cypress test code) | [eslint-plugin-cypress](https://github.com/cypress-io/eslint-plugin-cypress) |
| Standardized Code Formatting (indentation, etc) | [Prettier](https://prettier.io/) |
| Getting eslint and prettier to play well together (avoiding conflicting rules) | [eslint-config-prettier](https://github.com/prettier/eslint-config-prettier) |
| Running prettier rules as part of eslint | [eslint-plugin-prettier](https://github.com/prettier/eslint-plugin-prettier) |
| For CI, wait for the server to start before running tests | [wait-on](https://github.com/jeffbski/wait-on) |
| Keeping performance/accessibility in check via Lighthouse | [cypress-audit](https://mfrachet.github.io/cypress-audit/) which uses [lighthouse-ci](https://github.com/GoogleChrome/lighthouse-ci) |
| CSS Library | [Tailwind CSS](https://tailwindcss.com/) |

## Recommended Text Editor Extensions

* Prettier, to auto-format files (reads the project `.prettierrc.json` file) [list of editor plugins](https://prettier.io/docs/en/editors.html)
* ESLint, to auto-format files (reads the project `.eslintrc.json` file) [list of editor plugins](https://eslint.org/docs/user-guide/integrations#editors)
* MarkdownLint, to call out common markdown mistakes (reads the project `.markdownlint.json`). [list of editor plugins](https://github.com/DavidAnson/markdownlint#related)
* A [Tailwind CSS editor plugin](https://tailwindcss.com/docs/editor-setup), to do css class autocompletion from this library.

## Notes

* We're using tailwindcss via a CDN for now to try it out. If we decide to stick with it, we should add it in via a css build tool (like postcss) instead of using the CDN directly.
* For now everything (dev, prod, testing) uses port 8080. If/when that gets in the way of anything, we can set it up to use a different port for dev/testing vs prod.

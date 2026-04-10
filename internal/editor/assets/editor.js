// editor.js -- ES module for GOBL Editor
// Initializes CodeMirror with JSON schema support.

import { basicSetup, EditorView } from "codemirror";
import { EditorState, Compartment } from "@codemirror/state";

// Theme compartment allows dynamic reconfiguration.
const themeCompartment = new Compartment();

const lightTheme = EditorView.theme(
  {
    "&": { backgroundColor: "#ffffff", color: "#1e293b" },
    ".cm-content": { caretColor: "#334155" },
    ".cm-cursor, .cm-dropCursor": { borderLeftColor: "#334155" },
    "&.cm-focused .cm-selectionBackground, .cm-selectionBackground, .cm-content ::selection":
      { backgroundColor: "#dbeafe" },
    ".cm-panels": { backgroundColor: "#f8fafc", color: "#1e293b" },
    ".cm-panels.cm-panels-top": { borderBottom: "1px solid #e2e8f0" },
    ".cm-panels.cm-panels-bottom": { borderTop: "1px solid #e2e8f0" },
    ".cm-searchMatch": { backgroundColor: "#fef08a" },
    ".cm-searchMatch.cm-searchMatch-selected": { backgroundColor: "#fde047" },
    ".cm-activeLine": { backgroundColor: "#f8fafc" },
    ".cm-selectionMatch": { backgroundColor: "#e0f2fe" },
    ".cm-matchingBracket": { color: "#16a34a", backgroundColor: "#dcfce7" },
    ".cm-gutters": {
      backgroundColor: "#f8fafc",
      color: "#94a3b8",
      borderRight: "1px solid #e2e8f0",
    },
    ".cm-activeLineGutter": { backgroundColor: "#f1f5f9" },
    ".cm-foldPlaceholder": {
      backgroundColor: "#e2e8f0",
      color: "#64748b",
      border: "none",
    },
  },
  { dark: false },
);

const darkTheme = EditorView.theme(
  {
    "&": { backgroundColor: "#1e1e2e", color: "#cdd6f4" },
    ".cm-content": { caretColor: "#f5e0dc" },
    ".cm-cursor, .cm-dropCursor": { borderLeftColor: "#f5e0dc" },
    "&.cm-focused .cm-selectionBackground, .cm-selectionBackground, .cm-content ::selection":
      { backgroundColor: "#45475a" },
    ".cm-panels": { backgroundColor: "#181825", color: "#cdd6f4" },
    ".cm-panels.cm-panels-top": { borderBottom: "1px solid #313244" },
    ".cm-panels.cm-panels-bottom": { borderTop: "1px solid #313244" },
    ".cm-searchMatch": { backgroundColor: "#a6adc844" },
    ".cm-searchMatch.cm-searchMatch-selected": { backgroundColor: "#585b7044" },
    ".cm-activeLine": { backgroundColor: "#181825" },
    ".cm-selectionMatch": { backgroundColor: "#585b7044" },
    ".cm-matchingBracket": { color: "#a6e3a1", backgroundColor: "#45475a" },
    ".cm-gutters": {
      backgroundColor: "#181825",
      color: "#6c7086",
      borderRight: "1px solid #313244",
    },
    ".cm-activeLineGutter": { backgroundColor: "#1e1e2e" },
    ".cm-foldPlaceholder": {
      backgroundColor: "#45475a",
      color: "#cdd6f4",
      border: "none",
    },
  },
  { dark: true },
);

function isDark() {
  return document.documentElement.classList.contains("dark");
}

const defaultDoc = JSON.stringify(
  {
    $schema: "https://gobl.org/draft-0/bill/invoice",
    currency: "USD",
    issue_date: new Date().toISOString().slice(0, 10),
    supplier: {
      name: "Acme Inc.",
      tax_id: {
        country: "US",
      },
    },
    customer: {
      name: "Sample Customer",
    },
    lines: [
      {
        quantity: "10",
        item: {
          name: "Development Services",
          price: "100.00",
        },
        taxes: [
          {
            cat: "ST",
            percent: "8.25%",
          },
        ],
      },
    ],
  },
  null,
  2,
);

const { jsonSchema, updateSchema } = await import("codemirror-json-schema");

const SCHEMA_PREFIX = "https://gobl.org/draft-0/";
let activeSchemaURL = null;
let debounceTimer;

async function loadSchemaFromDoc(view) {
  try {
    const text = view.state.doc.toString();
    const m = text.match(/"\$schema"\s*:\s*"([^"]+)"/);
    const url = m?.[1];
    if (!url || url === activeSchemaURL) return;
    if (!url.startsWith(SCHEMA_PREFIX)) return;

    activeSchemaURL = url;
    const path = url.slice(SCHEMA_PREFIX.length);
    const res = await fetch("/v0/schemas/" + path + "?bundle");
    if (res.ok) updateSchema(view, await res.json());
  } catch (e) {
    console.warn("Schema loading failed:", e);
  }
}

const container = document.getElementById("editor-container");
container.replaceChildren();

const editor = new EditorView({
  state: EditorState.create({
    doc: defaultDoc,
    extensions: [
      basicSetup,
      themeCompartment.of(isDark() ? darkTheme : lightTheme),
      EditorView.lineWrapping,
      jsonSchema(),
      EditorView.updateListener.of((update) => {
        if (update.docChanged) {
          clearTimeout(debounceTimer);
          debounceTimer = setTimeout(
            () => loadSchemaFromDoc(update.view),
            500,
          );
        }
      }),
    ],
  }),
  parent: container,
});

window._cmEditor = editor;
window._cmSetDark = (dark) => {
  editor.dispatch({
    effects: themeCompartment.reconfigure(dark ? darkTheme : lightTheme),
  });
};
loadSchemaFromDoc(editor);


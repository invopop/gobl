// editor-data.js -- Alpine.js data components for the GOBL editor.
// Loaded as a synchronous script (before Alpine.js which is deferred)
// so the alpine:init listener is registered in time.
document.addEventListener("alpine:init", () => {
  Alpine.data("editor", () => ({
    loading: false,
    envelop: false,
    error: null,
    // Counter so each successful build creates a fresh FlashMessage.
    success: 0,

    async build() {
      const ed = window._cmEditor;
      if (!ed) return;

      this.loading = true;
      this.error = null;
      this.success = 0;

      try {
        const content = ed.state.doc.toString();
        const payload = JSON.parse(content);

        const res = await fetch("/v0/build", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ data: payload, envelop: this.envelop }),
        });

        const result = await res.json();

        if (!res.ok) {
          this.error = result;
          return;
        }

        const pretty = JSON.stringify(result, null, 2);
        ed.dispatch({
          changes: { from: 0, to: ed.state.doc.length, insert: pretty },
        });
        this.success++;
      } catch (e) {
        this.error = { message: e.message };
      } finally {
        this.loading = false;
      }
    },
  }));

  Alpine.data("darkModeToggle", () => ({
    dark: false,
    init() {
      const stored = localStorage.getItem("dark-mode");
      this.dark =
        stored !== null
          ? stored === "true"
          : document.documentElement.classList.contains("dark");
      document.documentElement.classList.toggle("dark", this.dark);
      this.applyEditorTheme();
    },
    toggle() {
      this.dark = !this.dark;
      document.documentElement.classList.toggle("dark", this.dark);
      localStorage.setItem("dark-mode", this.dark);
      this.applyEditorTheme();
    },
    applyEditorTheme() {
      if (window._cmSetDark) {
        window._cmSetDark(this.dark);
      }
    },
  }));
});

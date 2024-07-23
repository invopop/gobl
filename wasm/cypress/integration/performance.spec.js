describe("Website", () => {
  beforeEach(() => {
    cy.visit("http://localhost:8080");
  });
  it("Performs well", () => {
    // Config used to match the default Lighthouse settings from Chrome
    // https://mfrachet.github.io/cypress-audit/guides/lighthouse/good-to-know.html#lighthouse-scores-may-be-different-between-local-run-and-cypress-audit
    const customThresholds = {
      performance: 90,
      pwa: false,
    };

    const desktopConfig = {
      formFactor: "desktop",
      screenEmulation: {
        width: 1350,
        height: 940,
        deviceScaleRatio: 1,
        mobile: false,
        disable: false,
      },
      throttling: {
        rttMs: 40,
        throughputKbps: 11024,
        cpuSlowdownMultiplier: 1,
        requestLatencyMs: 0,
        downloadThroughputKbps: 0,
        uploadThroughputKbps: 0,
      },
    };

    cy.lighthouse(customThresholds, desktopConfig);
  });
});

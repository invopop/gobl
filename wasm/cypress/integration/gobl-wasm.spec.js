describe("Keygen", () => {
  beforeEach(() => {
    cy.visit("http://localhost:8080");
  });
  it("Produces a key", () => {
    cy.window().then((win) => {
      win.gobl.keygen().then((result) => {
        expect(result).to.not.be.empty;
        console.log(`RESULT: ${result}`);
      });
    });
  });
});

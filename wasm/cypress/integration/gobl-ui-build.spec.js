context("Inputs", () => {
  beforeEach(() => {
    cy.visit("http://localhost:8080");
  });

  describe("Valid Input", () => {
    beforeEach(() => {
      cy.get("#input-file").clear().paste(`{
        "doc": { 
            "$schema": "https://gobl.org/draft-0/note/message",
            "title": "Test Message",
            "content": "test content"
        }
      }`);
      cy.wait(2000); // build is async, it needs some time
    });
    it("displays success", () => {
      cy.get("#status").contains("Success");
    });
    it("shows a built file", () => {
      cy.get("#output-file").invoke("val").should("contain", "schema");
    });
  });

  describe("Invalid Input", () => {
    beforeEach(() => {
      cy.get("#input-file").clear().paste(`{
        "doczzz": { 
            "$schema": "https://gobl.org/draft-0/note/message",
            "title": "Test Message",
            "content": "test content"
        }
      }`);
      cy.wait(2000); // build is async, it needs some time
    });
    it("displays an error", () => {
      cy.get("#status").contains("Error");
    });
    it("shows no built file", () => {
      cy.get("#output-file").invoke("val").should("be.empty");
    });
  });
});

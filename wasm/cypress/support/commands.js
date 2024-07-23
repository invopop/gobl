// ***********************************************
// This example commands.js shows you how to
// create various custom commands and overwrite
// existing commands.
//
// For more comprehensive examples of custom
// commands please read more here:
// https://on.cypress.io/custom-commands
// ***********************************************
//
//
// -- This is a parent command --
// Cypress.Commands.add('login', (email, password) => { ... })
//
//
// -- This is a child command --
// Cypress.Commands.add('drag', { prevSubject: 'element'}, (subject, options) => { ... })
//
//
// -- This is a dual command --
// Cypress.Commands.add('dismiss', { prevSubject: 'optional'}, (subject, options) => { ... })
//
//
// -- This will overwrite an existing command --
// Cypress.Commands.overwrite('visit', (originalFn, url, options) => { ... })

import "@cypress-audit/lighthouse/commands";

// paste from https://gist.github.com/nickytonline/bcdef8ef00211b0faf7c7c0e7777aaf6
// note: this triggers two events from typing space then backspace

const paste = (subject, text) => {
  subject[0].value = text;
  return cy.get(subject).type(" {backspace}"); // the use of type to type a space and delete it after changing the value ensures that change detection kicks in
};

Cypress.Commands.add("paste", { prevSubject: "element" }, paste);

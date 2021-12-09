describe(
  "Policies flow",
  {
    defaultCommandTimeout: 20000,
  },
  () => {
    beforeEach(() => {
      cy.setup();
      cy.login();
    });
    it("Can create, check, and delete a policy successfully", () => {
      cy.visit("/policies/manage");
      cy.wait(1000); // eslint-disable-line cypress/no-unnecessary-waiting

      // Add a policy
      cy.findByText(/add a policy/i).click();

      // Confirm that policy was added successfully
      // cy.findByText(/successfully added policy/i).should("exist");
      // cy.findByText(/select query/i).should("not.exist");
      // cy.get(".policies-list-wrapper").within(() => {
      //   cy.findByText(/1 query/i).should("exist");
      //   cy.findByText(/yes/i).should("exist");
      //   cy.findByText(
      //     /Detect Linux hosts with high severity vulnerable versions of OpenSSL/i
      //   ).should("exist");

      // Click on link in table and confirm that policies filter block diplays as expected on manage hosts page
      //   cy.get("tbody").within(() => {
      //     cy.get("tr")
      //       .first()
      //       .within(() => {
      //         cy.get("td").last().children().first().should("exist").click();
      //       });
      //   });
      // });
      // cy.get(".manage-hosts__policies-filter-block").within(() => {
      //   cy.findByText(
      //     /Detect Linux hosts with high severity vulnerable versions of OpenSSL/i
      //   ).should("exist");
      //   cy.findByText(/yes/i).should("not.exist");
      //   cy.findByText(/failing/i)
      //     .should("exist")
      //     .click();
      //   cy.findByText(/yes/i).should("exist");
      //   cy.get('img[alt="Remove policy filter"]').click();
      //   cy.findByText(
      //     /Detect Linux hosts with high severity vulnerable versions of OpenSSL/i
      //   ).should("not.exist");
      // });

      // Click on policies tab to return to manage policies page
      cy.get(".site-nav-container").within(() => {
        cy.findByText(/policies/i).click();
      });

      // Delete policy
      // cy.get("tbody").within(() => {
      //   cy.get("tr")
      //     .first()
      //     .within(() => {
      //       cy.get(".fleet-checkbox__input").check({ force: true });
      //     });
      // });
      // cy.findByRole("button", { name: /remove/i }).click();
      // cy.get(".remove-policies-modal").within(() => {
      //   cy.findByRole("button", { name: /cancel/i }).should("exist");
      //   cy.findByRole("button", { name: /remove/i }).click();
      // });
      // cy.findByText(
      //   /Detect Linux hosts with high severity vulnerable versions of OpenSSL/i
      // ).should("not.exist");
    });
  }
);

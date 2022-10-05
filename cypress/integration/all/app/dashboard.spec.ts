import dashboardPage from "../../pages/dashboardPage";

describe("Dashboard", () => {
  before(() => {
    Cypress.session.clearAllSavedSessions();
    cy.setup();
    cy.loginWithCySession();
    cy.viewport(1200, 660);
  });

  after(() => {
    cy.logout();
  });

  describe("Operating systems card", () => {
    beforeEach(() => {
      cy.loginWithCySession();
      dashboardPage.visitsDashboardPage();
    });

    it("displays operating systems card if macOS platform is selected", () => {
      dashboardPage.switchesPlatform("macOS");
    });

    it("displays operating systems card if Windows platform is selected", () => {
      dashboardPage.switchesPlatform("Windows");
    });

    it("displays operating systems card if Linux platform is selected", () => {
      dashboardPage.switchesPlatform("Linux");
    });
  });
});

import "@testing-library/cypress/add-commands";

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
// Cypress.Commands.add("login", (email, password) => { ... })
//
//
// -- This is a child command --
// Cypress.Commands.add("drag", { prevSubject: 'element'}, (subject, options) => { ... })
//
//
// -- This is a dual command --
// Cypress.Commands.add("dismiss", { prevSubject: 'optional'}, (subject, options) => { ... })
//
//
// -- This will overwrite an existing command --
// Cypress.Commands.overwrite("visit", (originalFn, url, options) => { ... })

Cypress.Commands.add("setup", () => {
  cy.exec("make e2e-reset-db e2e-setup", { timeout: 20000 });
});

Cypress.Commands.add("login", (username, password) => {
  username ||= "admin";
  password ||= "user123#";
  cy.request("POST", "/api/v1/fleet/login", { username, password }).then(
    (resp) => {
      window.localStorage.setItem("FLEET::auth_token", resp.body.token);
    }
  );
});

Cypress.Commands.add("logout", () => {
  cy.request({
    url: "/api/v1/fleet/logout",
    method: "POST",
    body: {},
    auth: {
      bearer: window.localStorage.getItem("FLEET::auth_token"),
    },
  }).then(() => {
    window.localStorage.removeItem("FLEET::auth_token");
  });
});

Cypress.Commands.add("setupSMTP", () => {
  const body = {
    smtp_settings: {
      authentication_type: "authtype_none",
      enable_smtp: true,
      port: 1025,
      sender_address: "fleet@example.com",
      server: "localhost",
    },
  };

  cy.request({
    url: "/api/v1/fleet/config",
    method: "PATCH",
    body,
    auth: {
      bearer: window.localStorage.getItem("FLEET::auth_token"),
    },
  });
});

Cypress.Commands.add("setupSSO", (enable_idp_login = false) => {
  const body = {
    sso_settings: {
      enable_sso: true,
      enable_sso_idp_login: enable_idp_login,
      entity_id: "https://localhost:8080",
      idp_name: "SimpleSAML",
      issuer_uri: "http://localhost:8080/simplesaml/saml2/idp/SSOService.php",
      metadata_url: "http://localhost:9080/simplesaml/saml2/idp/metadata.php",
    },
  };

  cy.request({
    url: "/api/v1/fleet/config",
    method: "PATCH",
    body,
    auth: {
      bearer: window.localStorage.getItem("FLEET::auth_token"),
    },
  });
});

Cypress.Commands.add("loginSSO", () => {
  // Note these requests set cookies that are required for the SSO flow to
  // work properly. This is handled automatically by the browser.
  cy.request({
    method: "GET",
    url:
      "http://localhost:9080/simplesaml/saml2/idp/SSOService.php?spentityid=https://localhost:8080",
    followRedirect: false,
  }).then((firstResponse) => {
    const redirect = firstResponse.headers.location;

    cy.request({
      method: "GET",
      url: redirect,
      followRedirect: false,
    }).then((secondResponse) => {
      const el = document.createElement("html");
      el.innerHTML = secondResponse.body;
      const authState = el.getElementsByTagName("input").namedItem("AuthState")
        .defaultValue;

      cy.request({
        method: "POST",
        url: redirect,
        body: `username=sso_user&password=user123#&AuthState=${authState}`,
        form: true,
        followRedirect: false,
      }).then((finalResponse) => {
        el.innerHTML = finalResponse.body;
        const saml = el.getElementsByTagName("input").namedItem("SAMLResponse")
          .defaultValue;

        // Load the callback URL with the response from the IdP
        cy.visit({
          url: "/api/v1/fleet/sso/callback",
          method: "POST",
          body: {
            SAMLResponse: saml,
          },
        });
      });
    });
  });
});

Cypress.Commands.add("getEmails", () => {
  return cy
    .request("http://localhost:8025/api/v2/messages")
    .then((response) => {
      expect(response.status).to.eq(200);
      return response;
    });
});

Cypress.Commands.add("seedCore", () => {
  const authToken = window.localStorage.getItem("FLEET::auth_token");
  cy.exec("bash ./tools/api/fleet/teams/create_core", {
    env: {
      TOKEN: authToken,
      CURL_FLAGS: "-k",
      SERVER_URL: Cypress.config().baseUrl,
      // clear any value for FLEET_ENV_PATH since we set the environment explicitly just above
      FLEET_ENV_PATH: "",
    },
  });
});

Cypress.Commands.add("seedBasic", () => {
  const authToken = window.localStorage.getItem("FLEET::auth_token");
  cy.exec("bash ./tools/api/fleet/teams/create_basic", {
    env: {
      TOKEN: authToken,
      CURL_FLAGS: "-k",
      SERVER_URL: Cypress.config().baseUrl,
      // clear any value for FLEET_ENV_PATH since we set the environment explicitly just above
      FLEET_ENV_PATH: "",
    },
  });
});

Cypress.Commands.add("seedFigma", () => {
  const authToken = window.localStorage.getItem("FLEET::auth_token");
  cy.exec("bash ./tools/api/fleet/teams/create_figma", {
    env: {
      TOKEN: authToken,
      CURL_FLAGS: "-k",
      SERVER_URL: Cypress.config().baseUrl,
      // clear any value for FLEET_ENV_PATH since we set the environment explicitly just above
      FLEET_ENV_PATH: "",
    },
  });
});

Cypress.Commands.add("addUser", (username, options = {}) => {
  let { password, email, globalRole } = options;
  password ||= "test123#";
  email ||= `${username}@example.com`;
  globalRole ||= "admin";

  cy.exec(
    `./build/fleetctl user create --context e2e --username "${username}" --password "${password}" --email "${email}" --global-role "${globalRole}"`,
    { timeout: 20000 }
  );
});

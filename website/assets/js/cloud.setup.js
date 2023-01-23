/**
 * cloud.setup.js
 *
 * Configuration for this Sails app's generated browser SDK ("Cloud").
 *
 * Above all, the purpose of this file is to provide endpoint definitions,
 * each of which corresponds with one particular route+action on the server.
 *
 * > This file was automatically generated.
 * > (To regenerate, run `sails run rebuild-cloud-sdk`)
 */

Cloud.setup({

  /* eslint-disable */
  methods: {"downloadSitemap":{"verb":"GET","url":"/sitemap.xml","args":[]},"receiveUsageAnalytics":{"verb":"POST","url":"/api/v1/webhooks/receive-usage-analytics","args":["anonymousIdentifier","fleetVersion","licenseTier","numHostsEnrolled","numUsers","numTeams","numPolicies","numLabels","softwareInventoryEnabled","vulnDetectionEnabled","systemUsersEnabled","hostStatusWebhookEnabled","numWeeklyActiveUsers","numWeeklyPolicyViolationDaysActual","numWeeklyPolicyViolationDaysPossible","hostsEnrolledByOperatingSystem","hostsEnrolledByOrbitVersion","hostsEnrolledByOsqueryVersion","storedErrors","numHostsNotResponding","organization"]},"receiveFromGithub":{"verb":"GET","url":"/api/v1/webhooks/github","args":["botSignature","action","sender","repository","changes","issue","comment","pull_request","label"]},"receiveFromStripe":{"verb":"POST","url":"/api/v1/webhooks/receive-from-stripe","args":["id","type","data","webhookSecret"]},"deliverContactFormMessage":{"verb":"POST","url":"/api/v1/deliver-contact-form-message","args":["emailAddress","topic","firstName","lastName","message"]},"sendPasswordRecoveryEmail":{"verb":"POST","url":"/api/v1/entrance/send-password-recovery-email","args":["emailAddress"]},"signup":{"verb":"POST","url":"/api/v1/customers/signup","args":["emailAddress","password","organization","firstName","lastName","signupReason"]},"updateProfile":{"verb":"POST","url":"/api/v1/account/update-profile","args":["firstName","lastName","organization","emailAddress"]},"updatePassword":{"verb":"POST","url":"/api/v1/account/update-password","args":["oldPassword","newPassword"]},"updateBillingCard":{"verb":"POST","url":"/api/v1/account/update-billing-card","args":["stripeToken","billingCardLast4","billingCardBrand","billingCardExpMonth","billingCardExpYear"]},"login":{"verb":"POST","url":"/api/v1/customers/login","args":["emailAddress","password","rememberMe"]},"logout":{"verb":"GET","url":"/api/v1/account/logout","args":[]},"createQuote":{"verb":"POST","url":"/api/v1/customers/create-quote","args":["numberOfHosts"]},"saveBillingInfoAndSubscribe":{"verb":"POST","url":"/api/v1/customers/save-billing-info-and-subscribe","args":["quoteId","organization","firstName","lastName","paymentSource"]},"updatePasswordAndLogin":{"verb":"POST","url":"/api/v1/entrance/update-password-and-login","args":["password","token"]},"deliverDemoSignup":{"verb":"POST","url":"/api/v1/deliver-demo-signup","args":["emailAddress"]},"createOrUpdateOneNewsletterSubscription":{"verb":"POST","url":"/api/v1/create-or-update-one-newsletter-subscription","args":["emailAddress","subscribeTo"]},"unsubscribeFromAllNewsletters":{"verb":"GET","url":"/api/v1/unsubscribe-from-all-newsletters","args":["emailAddress"]},"generateLicenseKey":{"verb":"POST","url":"/api/v1/admin/generate-license-key","args":["numberOfHosts","organization","expiresAt"]},"createVantaAuthorizationRequest":{"verb":"POST","url":"/api/v1/create-vanta-authorization-request","args":["emailAddress","fleetInstanceUrl","fleetApiKey"]},"deliverMdmBetaSignup":{"verb":"POST","url":"/api/v1/deliver-mdm-beta-signup","args":["emailAddress","fullName","jobTitle"]},"deliverAppleCsr":{"verb":"POST","url":"/api/v1/deliver-apple-csr","args":["csr"]}}
  /* eslint-enable */

});

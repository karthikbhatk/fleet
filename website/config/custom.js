/**
 * Custom configuration
 * (sails.config.custom)
 *
 * One-off settings specific to your application.
 *
 * For more information on custom configuration, visit:
 * https://sailsjs.com/config/custom
 */

module.exports.custom = {

  /**************************************************************************
  *                                                                         *
  * The base URL to use during development.                                 *
  *                                                                         *
  * • No trailing slash at the end                                          *
  * • `http://` or `https://` at the beginning.                             *
  *                                                                         *
  * > This is for use in custom logic that builds URLs.                     *
  * > It is particularly handy for building dynamic links in emails,        *
  * > but it can also be used for user-uploaded images, webhooks, etc.      *
  *                                                                         *
  **************************************************************************/
  baseUrl: 'http://localhost:2024',

  /**************************************************************************
  *                                                                         *
  * The TTL (time-to-live) for various sorts of tokens before they expire.  *
  *                                                                         *
  **************************************************************************/
  passwordResetTokenTTL: 24*60*60*1000,// 24 hours
  emailProofTokenTTL:    24*60*60*1000,// 24 hours

  /**************************************************************************
  *                                                                         *
  * The extended length that browsers should retain the session cookie      *
  * if "Remember Me" was checked while logging in.                          *
  *                                                                         *
  **************************************************************************/
  rememberMeCookieMaxAge: 30*24*60*60*1000, // 30 days

  /**************************************************************************
  *                                                                         *
  * Automated email configuration                                           *
  *                                                                         *
  * Sandbox Sendgrid credentials for use during development, as well as any *
  * other default settings related to "how" and "where" automated emails    *
  * are sent.                                                               *
  *                                                                         *
  * (https://app.sendgrid.com/settings/api_keys)                            *
  *                                                                         *
  **************************************************************************/
  // sendgridSecret: 'SG.fake.3e0Bn0qSQVnwb1E4qNPz9JZP5vLZYqjh7sn8S93oSHU',
  //--------------------------------------------------------------------------
  // /\  Configure this to enable support for automated emails.
  // ||  (Important for password recovery, verification, contact form, etc.)
  //--------------------------------------------------------------------------

  // The sender that all outgoing emails will appear to come from.
  fromEmailAddress: 'noreply@example.com',
  fromName: 'The NEW_APP_NAME Team',

  // Email address for receiving support messages & other correspondences.
  // > If you're using the default privacy policy, this will be referenced
  // > as the contact email of your "data protection officer" for the purpose
  // > of compliance with regulations such as GDPR.
  internalEmailAddress: 'support+development@example.com',

  // Whether to require proof of email address ownership any time a new user
  // signs up, or when an existing user attempts to change their email address.
  verifyEmailAddresses: false,

  /**************************************************************************
  *                                                                         *
  * Billing & payments configuration                                        *
  *                                                                         *
  * (https://dashboard.stripe.com/account/apikeys)                          *
  *                                                                         *
  **************************************************************************/
  // stripePublishableKey: 'pk_test_Zzd814nldl91104qor5911gjald',
  // stripeSecret: 'sk_test_Zzd814nldl91104qor5911gjald',
  //--------------------------------------------------------------------------
  // /\  Configure these to enable support for billing features.
  // ||  (Or if you don't need billing, feel free to remove them.)
  //--------------------------------------------------------------------------

  /***************************************************************************
  *                                                                          *
  * Directly responsible individuals (DRIs) whose changes to areas of the    *
  * code respository (outside of the core product code) are auto-approved,   *
  * even during code freezes.                                                *
  *                                                                          *
  * See api/controllers/webhooks/receive-from-github.js for context.         *
  *                                                                          *
  ***************************************************************************/
  githubRepoDRIByPath: {
    'README.md': ['chris-mcgillicuddy'],// (github brandfront)
    'tools/fleetctl-npm/README.md': ['chris-mcgillicuddy', 'mike-j-thomas'],//« brandfront for fleetctl package on npm

    'CODEOWNERS': ['zwass', 'mikermcneil'],

    'articles': ['chris-mcgillicuddy', 'mike-j-thomas', 'eashaw', 'zwass', 'mikermcneil'],

    'handbook': ['chris-mcgillicuddy', 'mike-j-thomas', 'eashaw', 'mikermcneil', 'zwass'],// (default for handbook)
    'handbook/company': 'mikermcneil',
    'handbook/business-operations': ['hollidayn', 'charlottechance'],
    'handbook/engineering': 'zwass',
    'handbook/product': ['noahtalerman', 'zhumo'],
    'handbook/security': 'guillaumeross',
    'handbook/customers': ['alexmitchelliii','zayhanlon'],
    'handbook/marketing': ['mike-j-thomas','chris-mcgillicuddy'],

    'website': 'mikermcneil',// (default for website)
    'website/views': 'eashaw',
    'website/assets': 'eashaw',
    'website/config/routes.js': ['eashaw', 'mike-j-thomas'],// (for managing website URLs)
    'website/package.json': 'eashaw',

    'docs': ['chris-mcgillicuddy', 'zwass', 'mikermcneil'],// (default for docs)
    'docs/images': ['chris-mcgillicuddy', 'noahtalerman', 'eashaw', 'mike-j-thomas'],
    'docs/Using-Fleet/REST-API.md': 'ksatter',
    'docs/Contributing/API-for-contributors.md': 'ksatter',
    'docs/Deploying/FAQ.md': ['ksatter'],
    'docs/Contributing/FAQ.md': ['ksatter'],
    'docs/Using-Fleet/FAQ.md': ['ksatter'],

    'docs/01-Using-Fleet/standard-query-library/standard-query-library.yml': ['guillaumeross','zhumo','eashaw','zwass'],// (standard query library)
    'schema/': ['guillaumeross','zhumo','eashaw','zwass'],// (standard query library)
  },
  confidentialGithubRepoDRIByPath: {
    'readme.md': ['mikermcneil', 'zwass'],
    // GitHub issue templates
    '.github/ISSUE_TEMPLATE': ['mikermcneil', 'zwass'],
    '.github/ISSUE_TEMPLATE/2-website-changes.md': 'mike-j-thomas',
    '.github/ISSUE_TEMPLATE/3-opportunity Fleet Premium PoV.md': 'alexmitchelliii',
    '.github/ISSUE_TEMPLATE/3-sale.md': 'alexmitchelliii',
    '.github/ISSUE_TEMPLATE/4-release.md': 'noahtalerman',
    '.github/ISSUE_TEMPLATE/5-monthly-accounting.md': 'hollidayn',
    '.github/ISSUE_TEMPLATE/6-speaking-event.md': 'mike-j-thomas',
    '.github/ISSUE_TEMPLATE/9-renewal-process.md': ['hollidayn', 'alexmitchelliii'],
    '.github/ISSUE_TEMPLATE/hiring.md': 'charlottechance',
    '.github/ISSUE_TEMPLATE/onboarding.md': 'charlottechance',
    '.github/ISSUE_TEMPLATE/y-offboarding.md': 'charlottechance',
    '.github/ISSUE_TEMPLATE/equity-grants.md': ['charlottechance','hollidayn'],
    'ceo-shadow-onboarding': ['mikermcneil', 'charlottechance'],

    'mdm_profiles': 'guillaumeross',
    'vpn': 'rfairburn',
  },


  /***************************************************************************
  *                                                                          *
  * Any other custom config this Sails app should use during development.    *
  * (and possibly in ALL environments, if not overridden in config/env/)     *
  *                                                                          *
  ***************************************************************************/
  // Contact form:
  // slackWebhookUrlForContactForm: '…',

  // GitHub bot:
  // githubAccessToken: '…',
  // githubBotWebhookSecret: '…',
  // slackWebhookUrlForGithubBot: '…',
  // mergeFreezeAccessToken: '…',

  //…

};

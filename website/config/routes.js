/**
 * Route Mappings
 * (sails.config.routes)
 *
 * Your routes tell Sails what to do each time it receives a request.
 *
 * For more information on configuring custom routes, check out:
 * https://sailsjs.com/anatomy/config/routes-js
 */

module.exports.routes = {

  //  ╦ ╦╔═╗╔╗ ╔═╗╔═╗╔═╗╔═╗╔═╗
  //  ║║║║╣ ╠╩╗╠═╝╠═╣║ ╦║╣ ╚═╗
  //  ╚╩╝╚═╝╚═╝╩  ╩ ╩╚═╝╚═╝╚═╝
  'GET /':                   { action: 'view-homepage-or-redirect' },

  '/blog':           'https://medium.com/fleetdm',

  '/pricing':        (req, res)=>{
    // FUTURE: pricing page explaining commercial support and EE, w/ our subscription terms
    res.redirect('mailto:todo@example.com?subject=Pricing question&body=Please do not send this email!\n\nWe are a very young company and still working on our processes.  For now, if you have a pricing question or would like to know Fleet\'s latest pricing and support tiers, please create an issue at https://github.com/fleetdm/fleet/issues.  Thank you!');
  },

  '/legal/terms':        (req, res)=>{
    // FUTURE: page w/ our subscription terms
    res.redirect('mailto:todo@example.com?subject=Legal question&body=Please do not send this email!\n\nWe are a very young company and still working on our processes.  For now, if you have a legal question or would like to review Fleet\'s subscription terms, please send a DM to mikermcneil or Zach Wasserman in the "osquery" Slack workspace.  Thank you!');
  },

  '/security':       (req, res)=>{
    // FUTURE: make a page- check out how Sails does it, and also https://about.gitlab.com/security/
    res.redirect('mailto:todo@example.com?subject=Security vulnerability&body=Please do not send this email!\n\nWe are a very young company and still working on our processes.  For now, if you have a security vulnerability to report, please send a DM to mikermcneil or Zach Wasserman in the "osquery" Slack workspace.  Thank you for letting us know!');
  },

  '/company/about':          '/blog', // FUTURE: brief "about" page explaining the origins of the company
  '/company/stewardship':    'https://github.com/fleetdm/fleet', // FUTURE: page about how we approach open source and our commitments to the community
  '/company/contact':        (req, res)=>{
    // FUTURE: Page about this.  See “Contacting Fleet” in google drive for verbiage (updated holiday hours) - better to just make it all open source and check it in here.  For inspiration, see also: https://about.gitlab.com/company/contact/
    res.redirect('mailto:todo@example.com?subject=I have a question&body=Please do not send this email!\n\nWe are a very young company and still working on our processes.  For now, if you have questions or feedback for us, please create an issue at https://github.com/fleetdm/fleet/issues.  Thank you!');
  },

  // 'GET /welcome/:unused?':   { action: 'dashboard/view-welcome' },

  // 'GET /faq':                { action:   'view-faq' },
  // 'GET /legal/terms':        { action:   'legal/view-terms' },
  // 'GET /legal/privacy':      { action:   'legal/view-privacy' },
  // 'GET /contact':            { action:   'view-contact' },

  // 'GET /signup':             { action: 'entrance/view-signup' },
  // 'GET /email/confirm':      { action: 'entrance/confirm-email' },
  // 'GET /email/confirmed':    { action: 'entrance/view-confirmed-email' },

  // 'GET /login':              { action: 'entrance/view-login' },
  // 'GET /password/forgot':    { action: 'entrance/view-forgot-password' },
  // 'GET /password/new':       { action: 'entrance/view-new-password' },

  // 'GET /account':            { action: 'account/view-account-overview' },
  // 'GET /account/password':   { action: 'account/view-edit-password' },
  // 'GET /account/profile':    { action: 'account/view-edit-profile' },


  //  ╔╦╗╦╔═╗╔═╗  ╦═╗╔═╗╔╦╗╦╦═╗╔═╗╔═╗╔╦╗╔═╗   ┬   ╔╦╗╔═╗╦ ╦╔╗╔╦  ╔═╗╔═╗╔╦╗╔═╗
  //  ║║║║╚═╗║    ╠╦╝║╣  ║║║╠╦╝║╣ ║   ║ ╚═╗  ┌┼─   ║║║ ║║║║║║║║  ║ ║╠═╣ ║║╚═╗
  //  ╩ ╩╩╚═╝╚═╝  ╩╚═╚═╝═╩╝╩╩╚═╚═╝╚═╝ ╩ ╚═╝  └┘   ═╩╝╚═╝╚╩╝╝╚╝╩═╝╚═╝╩ ╩═╩╝╚═╝
  // '/logout':                  '/api/v1/account/logout',
  '/company':                    '/company/about',
  '/support':                    '/company/contact',
  '/contact':                    '/company/contact',
  '/legal':                      '/legal/terms',
  '/terms':                      '/legal/terms',


  //  ╦ ╦╔═╗╔╗ ╦ ╦╔═╗╔═╗╦╔═╔═╗
  //  ║║║║╣ ╠╩╗╠═╣║ ║║ ║╠╩╗╚═╗
  //  ╚╩╝╚═╝╚═╝╩ ╩╚═╝╚═╝╩ ╩╚═╝
  // …


  //  ╔═╗╔═╗╦  ╔═╗╔╗╔╔╦╗╔═╗╔═╗╦╔╗╔╔╦╗╔═╗
  //  ╠═╣╠═╝║  ║╣ ║║║ ║║╠═╝║ ║║║║║ ║ ╚═╗
  //  ╩ ╩╩  ╩  ╚═╝╝╚╝═╩╝╩  ╚═╝╩╝╚╝ ╩ ╚═╝
  // Note that, in this app, these API endpoints may be accessed using the `Cloud.*()` methods
  // from the Parasails library, or by using those method names as the `action` in <ajax-form>.
  // '/api/v1/account/logout':                           { action: 'account/logout' },
  // 'PUT   /api/v1/account/update-password':            { action: 'account/update-password' },
  // 'PUT   /api/v1/account/update-profile':             { action: 'account/update-profile' },
  // 'PUT   /api/v1/account/update-billing-card':        { action: 'account/update-billing-card' },
  // 'PUT   /api/v1/entrance/login':                        { action: 'entrance/login' },
  // 'POST  /api/v1/entrance/signup':                       { action: 'entrance/signup' },
  // 'POST  /api/v1/entrance/send-password-recovery-email': { action: 'entrance/send-password-recovery-email' },
  // 'POST  /api/v1/entrance/update-password-and-login':    { action: 'entrance/update-password-and-login' },
  // 'POST  /api/v1/deliver-contact-form-message':          { action: 'deliver-contact-form-message' },

};

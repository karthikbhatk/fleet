module.exports = {


  friendlyName: 'View basic documentation',


  description: 'Display "Basic documentation" page.',


  urlWildcardSuffix: 'pageUrlSuffix',


  inputs: {
    pageUrlSuffix : {
      description: 'The relative path to the doc page from within this route.  (i.e. the URL wildcard suffix)',
      example: 'using-fleet/supported-browsers',
      type: 'string',
      defaultsTo: ''
    }
  },


  exits: {
    success: { viewTemplatePath: 'pages/docs/basic-documentation' },
    badConfig: { responseType: 'badConfig' },
    notFound: { responseType: 'notFound' },
    redirect: { responseType: 'redirect' },
  },


  fn: async function ({pageUrlSuffix}) {

    if (!_.isObject(sails.config.builtStaticContent) || !_.isArray(sails.config.builtStaticContent.markdownPages) || !sails.config.builtStaticContent.compiledPagePartialsAppPath) {
      throw {badConfig: 'builtStaticContent.markdownPages'};
    }

    let SECTION_URL_PREFIX = '/docs';

    // Serve appropriate page content.
    // Note that this action also serves the '/docs' landing page, as well as individual doc pages.
    // > Inspired by https://github.com/sailshq/sailsjs.com/blob/b53c6e6a90c9afdf89e5cae00b9c9dd3f391b0e7/api/controllers/documentation/view-documentation.js
    let thisPage = _.find(sails.config.builtStaticContent.markdownPages, { url: SECTION_URL_PREFIX + '/' + pageUrlSuffix });
    if (!thisPage) {// If there's no matching page, try a revised version of the URL suffix that's lowercase, with internal slashes deduped, and any trailing slash or whitespace trimmed
      let revisedPageUrlSuffix = pageUrlSuffix.toLowerCase().replace(/\/+/g, '/').replace(/\/+\s*$/,'');
      thisPage = _.find(sails.config.builtStaticContent.markdownPages, { url: SECTION_URL_PREFIX + '/' + revisedPageUrlSuffix });
      if (thisPage) {// If we matched a page with the revised suffix, then redirect to that rather than rendering it, so the URL gets cleaned up.
        throw {redirect: thisPage.url};
      } else {// If no page could be found even with the revised suffix, then throw a 404 error.
        throw 'notFound';
      }
    }

    // Respond with view.
    return {
      path: require('path'),
      thisPage: thisPage,
      markdownPages: sails.config.builtStaticContent.markdownPages,
      compiledPagePartialsAppPath: sails.config.builtStaticContent.compiledPagePartialsAppPath,
      pageTitleForMeta: (
        thisPage.title !== 'Readme.md' ? thisPage.title + ' | Fleet documentation'// « custom meta title for this page, if provided in markdown
        : 'Documentation | Fleet for osquery' // « otherwise we're on the landing page for this section of the site, so we'll follow the title format of other top-level pages
      ),
      pageDescriptionForMeta: (
        thisPage.meta.description ? thisPage.meta.description // « custom meta description for this page, if provided in markdown
        : 'Documentation for Fleet for osquery.'// « otherwise use the generic description
      ),
    };

  }


};

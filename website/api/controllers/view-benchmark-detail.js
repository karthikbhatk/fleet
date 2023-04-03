module.exports = {


  friendlyName: 'View benchmark detail',


  description: 'Display "Benchmark detail" page.',


  inputs: {
    slug: { type: 'string', required: true, description: 'A slug uniquely identifying this benchmark in the library.', example: 'cis-ensure-all-apple-provided-software-is-current-fleetd-required' },
  },


  exits: {
    success: { viewTemplatePath: 'pages/benchmark-detail' },
    notFound: { responseType: 'notFound' },
    badConfig: { responseType: 'badConfig' },
  },


  fn: async function ({ slug }) {

    if (!_.isObject(sails.config.builtStaticContent) || !_.isArray(sails.config.builtStaticContent.cisBenchmarks)) {
      console.log("line 25");
      throw {badConfig: 'builtStaticContent.cisBenchmarks'};
    } else if (!_.isString(sails.config.builtStaticContent.cisBenchmarkLibraryMacYmlRepoPath)) {
      console.log("line 28");
      throw {badConfig: 'builtStaticContent.cisBenchmarkLibraryMacLibraryYmlRepoPath'};
    }

    // Serve appropriate content for query.
    // > Inspired by https://github.com/sailshq/sailsjs.com/blob/b53c6e6a90c9afdf89e5cae00b9c9dd3f391b0e7/api/controllers/documentation/view-documentation.js
    let benchmark = _.find(sails.config.builtStaticContent.cisBenchmarks, { slug: slug });
    if (!benchmark) {
      console.log("line 36");
      throw 'notFound';
    }

    // Setting the meta title and description of this page using the query object, and falling back to a generic title or description if query.name or query.description are missing.
    let pageTitleForMeta = benchmark.name ? benchmark.name + ' | Benchmark details' : 'Benchmark details | Fleet for osquery';
    let pageDescriptionForMeta = benchmark.description ? benchmark.description : 'View more information about a benchmark in Fleet\'s CIS benchmark library';
    // Respond with view.
    return {
      benchmark,
      cisBenchmarkLibraryMacYmlRepoPath: sails.config.builtStaticContent.cisBenchmarkLibraryMacYmlRepoPath,
      pageTitleForMeta,
      pageDescriptionForMeta,
    };

  }


};

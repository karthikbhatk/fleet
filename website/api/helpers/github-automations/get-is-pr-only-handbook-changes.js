module.exports = {


  friendlyName: 'Get "Is PR only handbook changes"',


  description: 'Checks each file in a GitHub pull request, and returns true if the PR only changes files in the handbook folder',


  inputs: {
    prNumber: { type: 'number', example: 382, required: true },
  },

  exits: {

    success: {
      outputFriendlyName: 'Is PR only handbook changes',
      outputDescription: 'Whether the provided pull request only makes changes to the handbook',
      outputType: 'boolean',
    },

  },


  fn: async function ({prNumber}) {

    require('assert')(sails.config.custom.githubAccessToken);
    let DRI_BY_PATH = sails.config.custom.githubRepoDRIByPath;
    let owner = 'fleetdm';
    let repo = 'fleet';
    let baseHeaders = {
      'User-Agent': 'Fleet labels',
      'Authorization': `token ${sails.config.custom.githubAccessToken}`
    };

    // Check the path of each file that this PR makes changes to
    return await sails.helpers.flow.build(async()=>{

      let isHandbookOnlyPR = false;

      // [?] https://docs.github.com/en/rest/reference/pulls#list-pull-requests-files
      let changedPaths = _.pluck(await sails.helpers.http.get(`https://api.github.com/repos/${owner}/${repo}/pulls/${prNumber}/files`, {
        per_page: 100,//eslint-disable-line camelcase
      }, baseHeaders).retry(), 'filename');// (don't worry, it's the whole path, not the filename)

      isHandbookOnlyPR = _.all(changedPaths, (changedPath)=>{
        changedPath = changedPath.replace(/\/+$/,'');// « trim trailing slashes, just in case (b/c otherwise could loop forever)

        let numRemainingPathsToCheck = changedPath.split('/').length;
        while (numRemainingPathsToCheck > 0) {
          if(_.startsWith(changedPath, 'handbook/')) {
            return true;
          }
          numRemainingPathsToCheck--;
        }//∞
      });//∞

      if (isHandbookOnlyPR) {
        return true;
      } else {
        return false;
      }
    });

  }


};


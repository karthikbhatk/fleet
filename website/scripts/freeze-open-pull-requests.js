module.exports = {


  friendlyName: 'Freeze open pull requests',


  description: 'Freeze existing pull requests open on https://github.com/fleetdm/fleet, except those that consist exclusively of changes to files where the author is the DRI, according to auto-approval rules.',


  inputs: {
    dry: { type: 'boolean', defaultsTo: false, description: 'Whether to do a dry run, and not actually freeze anything.' },
  },


  fn: async function ({dry: isDryRun}) {

    sails.log('Running custom shell script... (`sails run freeze-open-pull-requests`)');

    let owner = 'fleetdm';
    let repo = 'fleet';
    let baseHeaders = {
      'User-Agent': 'sails run freeze-open-pull-requests',
      'Authorization': `token ${sails.config.custom.githubAccessToken}`
    };

    // Fetch open pull requests
    // [?] https://docs.github.com/en/rest/pulls/pulls#list-pull-requests
    let openPullRequests = await sails.helpers.http.get(`https://api.github.com/repos/${owner}/${repo}/pulls`, {
      state: 'open',
      per_page: 100,//eslint-disable-line camelcase
    }, baseHeaders);

    let SECONDS_TO_WAIT = 5;
    sails.log(`Examining and potentially freezing ${openPullRequests.length} PRs in ${SECONDS_TO_WAIT} seconds…  (To cancel, press CTRL+C quickly!)`);
    await sails.helpers.flow.pause(SECONDS_TO_WAIT*1000);

    // For all open pull requests…
    await sails.helpers.flow.simultaneouslyForEach(openPullRequests, async(pullRequest)=>{

      let prNumber = pullRequest.number;
      let prAuthor = pullRequest.user.login;
      require('assert')(prAuthor !== undefined);

      // Freeze, if appropriate.
      // (Check the PR's author versus the intersection of DRIs for all changed files.)
      let isAuthorPreapproved = await sails.helpers.githubAutomations.getIsPrPreapproved.with({
        prNumber: prNumber,
        githubUserToCheck: prAuthor,
        isGithubUserMaintainerOrDoesntMatter: true// « doesn't matter here because no auto-approval is happening.  Worst case, a community PR to an area with a "*" in the DRI mapping remains unfrozen.
      });

      if (isDryRun) {
        sails.log(`#${prNumber} by @${prAuthor}:`, isAuthorPreapproved ? 'Would have skipped freeze…' : 'Would have frozen…');
      } else {
        sails.log(`#${prNumber} by @${prAuthor}:`, isAuthorPreapproved ? 'Skipping freeze…' : 'Freezing…');
        // TODO: uncomment when ready
        // [?] https://docs.github.com/en/rest/reference/pulls#create-a-review-for-a-pull-request
        // await sails.helpers.http.post(`https://api.github.com/repos/${owner}/${repo}/pulls/${prNumber}/reviews`, {
        //   event: 'REQUEST_CHANGES',
        //   body: 'The repository has been frozen for an upcoming release.  In case of emergency, you can dismiss this review and merge.'
        // }, baseHeaders);
      }//
    });

  }


};


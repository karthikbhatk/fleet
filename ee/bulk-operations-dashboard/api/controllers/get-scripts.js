module.exports = {


  friendlyName: 'Get scripts',


  description: '',


  inputs: {

  },


  exits: {

  },


  fn: async function (inputs) {

    let teamsResponseData = await sails.helpers.http.get.with({
      url: '/api/v1/fleet/teams',
      baseUrl: sails.config.custom.fleetBaseUrl,
      headers: {
        Authorization: `Bearer ${sails.config.custom.fleetApiToken}`
      }
    })
    .timeout(120000)
    .retry(['requestFailed', {name: 'TimeoutError'}]);

    let allTeams = teamsResponseData.teams;

    let teamApids = _.pluck(allTeams, 'id');
    let teams = [];
    let teamsInformation = [];
    for(let team of allTeams) {
      teams.push({
        fleetApid: team.id,
        teamName: team.name,
      });
    }
    // Add the "team" for hosts with no team
    teams.push({
      fleetApid: 0,
      teamName: 'No team',
    });


    let allScripts = [];

    for(let teamApid of teamApids){
      let scriptsResponseData = await sails.helpers.http.get.with({
        url: `/api/v1/fleet/scripts?team_id=${teamApid}`,
        baseUrl: sails.config.custom.fleetBaseUrl,
        headers: {
          Authorization: `Bearer ${sails.config.custom.fleetApiToken}`
        }
      })
      .timeout(120000)
      .retry(['requestFailed', {name: 'TimeoutError'}]);
      let scriptsForThisTeam = scriptsResponseData.scripts;
      if(scriptsForThisTeam !== null) {
        allScripts = allScripts.concat(scriptsForThisTeam);
      }
    }

    // Grab all of the configuration scripts on the Fleet instance.
    let noTeamConfigurationScriptsResponseData = await sails.helpers.http.get.with({
      url: '/api/v1/fleet/scripts',
      baseUrl: sails.config.custom.fleetBaseUrl,
      headers: {
        Authorization: `Bearer ${sails.config.custom.fleetApiToken}`
      }
    })
    .timeout(120000)
    .retry(['requestFailed', {name: 'TimeoutError'}]);
    let scriptsForThisTeam = noTeamConfigurationScriptsResponseData.scripts;

    if(scriptsForThisTeam !== null){
      allScripts = allScripts.concat(scriptsForThisTeam);
    }

    if(allScripts === [ null ]){
      return {scripts: [], teams};
    }
    let scriptsOnThisFleetInstance = [];

    let allScriptsByIdentifier = _.groupBy(allScripts, 'name');
    for(let scriptIdentifier in allScriptsByIdentifier) {
      if(scriptIdentifier === null){
        continue;
      }
      let teamIdsForThisProfile = _.pluck(allScriptsByIdentifier[scriptIdentifier], 'team_id');
      let teamsForThisProfile = [];
      for(let script of allScriptsByIdentifier[scriptIdentifier]){
        let informationAboutThisScript = {
          scriptFleetApid: script.id,
          fleetApid: script.team_id ? script.team_id : 0,
          teamName: script.team_id ? _.find(teams, {fleetApid: script.team_id}).teamName : 'No team',
        }
        teamsForThisProfile.push(informationAboutThisScript);
      }
      let script = allScriptsByIdentifier[scriptIdentifier][0];// Grab the first script returned in the api repsonse to build our script configuration.
      let scriptInformation = {
        name: script.name,
        identifier: scriptIdentifier,
        platform: _.endsWith(script.name, 'sh') ? 'macOS & Linux' : 'Windows',
        createdAt: new Date(script.created_at).getTime(),
        teams: teamsForThisProfile
      }
      scriptsOnThisFleetInstance.push(scriptInformation)
    }
    scriptsOnThisFleetInstance = _.sortByOrder(scriptsOnThisFleetInstance, 'name', 'asc');
    return scriptsOnThisFleetInstance;

  }


};

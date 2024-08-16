module.exports = {


  friendlyName: 'Edit script',


  description: '',

  files: ['newScript'],

  inputs: {
    script: {
      type: {},
      description: 'The script that is being editted',
      required: true,
    },
    newTeamIds: {
      type: ['ref'],
      description: 'An array of teams that this script will be added to.'
    },
    newScript: {
      type: 'ref',
      description: 'A file that will be replacing the script.'
    },
  },


  exits: {
    scriptNameDoesNotMatch: {
      description: 'The provided replacement script\'s filename does not match the name of the script on the Fleet instance.',
      statusCode: 400,
    },
  },


  fn: async function ({script, newTeamIds, newScript}) {
    if(newScript.isNoop){
      newScript.noMoreFiles();
      newScript = undefined;
    }
    let scriptContents; // The raw text contents of a script file.
    let filename;
    let extension;
    // If there is not a new script, and the script is deployed (has teams array === deployed), download the script to be able to add it to other teams.
    if(!newScript){
      let scriptFleetApid = script.teams[0].scriptFleetApid;
      let scriptDownloadResponse = await sails.helpers.http.sendHttpRequest.with({
        method: 'GET',
        url: `${sails.config.custom.fleetBaseUrl}/api/v1/fleet/scripts/${scriptFleetApid}?alt=media`,
        headers: {
          Authorization: `Bearer ${sails.config.custom.fleetApiToken}`
        }
      });
      let contentDispositionHeader = scriptDownloadResponse.headers['content-disposition'];
      let filenameMatch = contentDispositionHeader.match(/filename="(.+?)"/);
      filename = filenameMatch[1];
      extension = '.'+filename.split('.').pop();
      filename = filename.replace(/^\d{4}-\d{2}-\d{2}[_|\s]?/, '');
      scriptContents = scriptDownloadResponse.body;
    } else if(newScript) {
      let file = await sails.reservoir(newScript);
      scriptContents = file[0].contentBytes;
      let scriptFilename = file[0].name;
      filename = scriptFilename.replace(/^\d{4}-\d{2}-\d{2}[_|\s]?/, '').replace(/\.[^/.]+$/, '');
      extension = '.'+scriptFilename.split('.').pop();
      if(script.name !== filename+extension){
        throw 'scriptNameDoesNotMatch';
      }
    }

    if(!newScript){
      let currentScriptTeamIds = _.pluck(script.teams, 'fleetApid');
      let addedTeams = _.difference(newTeamIds, currentScriptTeamIds);
      let removedTeams = _.difference(currentScriptTeamIds, newTeamIds);
      let removedTeamsInfo = _.filter(script.teams, (team)=>{
        return removedTeams.includes(team.fleetApid);
      });
      for(let script of removedTeamsInfo){
        await sails.helpers.http.sendHttpRequest.with({
          method: 'DELETE',
          baseUrl: sails.config.custom.fleetBaseUrl,
          url: `/api/v1/fleet/scripts/${script.scriptFleetApid}`,
          headers: {
            Authorization: `Bearer ${sails.config.custom.fleetApiToken}`,
          }
        });
      }
      for(let teamApid of addedTeams){
        // Build a request body for the team.
        let requestBodyForThisTeam = {
          script: {
            options: {
              filename: filename+extension,
              contentType: 'application/octet-stream'
            },
            value: scriptContents,
          }
        };
        let addScriptUrl;
        // If the script is being added to the "no team" team, then we need to include the team ID of the no team team in the request URL
        if(Number(teamApid) === 0){
          addScriptUrl = `/api/v1/fleet/scripts?team_id=${teamApid}`;
        } else {
          // Otherwise, the team_id needs to be included in the request's formData.
          addScriptUrl = `/api/v1/fleet/scripts`;
          requestBodyForThisTeam.team_id = Number(teamApid);// eslint-disable-line camelcase
        }
        await sails.helpers.http.sendHttpRequest.with({
          method: 'POST',
          baseUrl: sails.config.custom.fleetBaseUrl,
          url: addScriptUrl,
          enctype: 'multipart/form-data',
          body: requestBodyForThisTeam,
          headers: {
            Authorization: `Bearer ${sails.config.custom.fleetApiToken}`,
          },
        });
      }
    } else {
      for(let scriptId of script.teams){
        await sails.helpers.http.sendHttpRequest.with({
          method: 'DELETE',
          baseUrl: sails.config.custom.fleetBaseUrl,
          url: `/api/v1/fleet/scripts/${scriptId.scriptFleetApid}`,
          headers: {
            Authorization: `Bearer ${sails.config.custom.fleetApiToken}`,
          }
        });
      }
      for(let teamApid of newTeamIds){
        // Build a request body for the team.
        let requestBodyForThisTeam = {
          script: {
            options: {
              filename: filename+extension,
              contentType: 'application/octet-stream'
            },
            value: scriptContents,
          }
        };
        let addScriptUrl;
        // If the script is being added to the "no team" team, then we need to include the team ID of the no team team in the request URL
        if(Number(teamApid) === 0){
          addScriptUrl = `/api/v1/fleet/scripts?team_id=${teamApid}`;
        } else {
          // Otherwise, the team_id needs to be included in the request's formData.
          addScriptUrl = `/api/v1/fleet/scripts`;
          requestBodyForThisTeam.team_id = Number(teamApid);// eslint-disable-line camelcase
        }
        await sails.helpers.http.sendHttpRequest.with({
          method: 'POST',
          baseUrl: sails.config.custom.fleetBaseUrl,
          url: addScriptUrl,
          enctype: 'multipart/form-data',
          body: requestBodyForThisTeam,
          headers: {
            Authorization: `Bearer ${sails.config.custom.fleetApiToken}`,
          },
        });
      }
    }

    // All done.
    return;

  }


};

import endpoints from "kolide/endpoints";
import { INewMembersBody, IRemoveMembersBody, ITeam } from "interfaces/team";
import { ICreateTeamFormData } from "pages/admin/TeamManagementPage/components/CreateTeamModal/CreateTeamModal";

interface ITeamsResponse {
  teams: ITeam[];
}

export default (client: any) => {
  return {
    create: (formData: ICreateTeamFormData) => {
      const { TEAMS } = endpoints;

      return client
        .authenticatedPost(client._endpoint(TEAMS), JSON.stringify(formData))
        .then((response: ITeam) => response);
    },

    destroy: (teamId: number) => {
      const { TEAMS } = endpoints;
      const endpoint = `${client._endpoint(TEAMS)}/${teamId}`;
      return client.authenticatedDelete(endpoint);
    },
    load: (teamId: number) => {
      const { TEAMS } = endpoints;
      const endpoint = client._endpoint(`${TEAMS}/${teamId}`);

      return client
        .authenticatedGet(endpoint)
        .then((response: any) => response.team);
    },
    loadAll: ({ page = 0, perPage = 100, globalFilter = "" }) => {
      const { TEAMS } = endpoints;

      // TODO: add this query param logic to client class
      const pagination = `page=${page}&per_page=${perPage}`;

      let searchQuery = "";
      if (globalFilter !== "") {
        searchQuery = `&query=${globalFilter}`;
      }

      const teamsEndpoint = `${TEAMS}?${pagination}${searchQuery}`;
      return client
        .authenticatedGet(client._endpoint(teamsEndpoint))
        .then((response: ITeamsResponse) => {
          const { teams } = response;
          return teams;
        });
    },
    update: (teamId: number, updateParams: ITeam) => {
      const { TEAMS } = endpoints;
      const updateTeamEndpoint = `${client.baseURL}${TEAMS}/${teamId}`;

      return client
        .authenticatedPatch(updateTeamEndpoint, JSON.stringify(updateParams))
        .then((response: ITeam) => response);
    },
    addMembers: (teamId: number, newMembers: INewMembersBody) => {
      const { TEAMS_MEMBERS } = endpoints;
      return client
        .authenticatedPatch(
          client._endpoint(TEAMS_MEMBERS(teamId)),
          JSON.stringify(newMembers)
        )
        .then((response: ITeam) => response);
    },
    removeMembers: (teamId: number, removeMembers: IRemoveMembersBody) => {
      const { TEAMS_MEMBERS } = endpoints;
      return client.authenticatedDelete(
        client._endpoint(TEAMS_MEMBERS(teamId)),
        {},
        JSON.stringify(removeMembers)
      );
    },
  };
};

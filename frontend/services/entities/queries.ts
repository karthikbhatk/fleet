/* eslint-disable  @typescript-eslint/explicit-module-boundary-types */
import sendRequest, { getError } from "services";
import endpoints from "utilities/endpoints";
import { IQueryFormData } from "interfaces/query";
import { ISelectedTargets } from "interfaces/target";
import { AxiosResponse } from "axios";
import { buildQueryStringFromParams } from "utilities/url";

// Mock API requests to be used in developing FE for #7765 in parallel with BE development
// import { sendRequest } from "services/mock_service/service/service";

// Mock API requests to be used in developing FE for #7765 in parallel with BE development
// import { sendRequest } from "services/mock_service/service/service";

export default {
  create: ({ description, name, query, observer_can_run }: IQueryFormData) => {
    const { QUERIES } = endpoints;

    return sendRequest("POST", QUERIES, {
      description,
      name,
      query,
      observer_can_run,
    });
  },
  destroy: (id: string | number) => {
    const { QUERIES } = endpoints;
    const path = `${QUERIES}/id/${id}`;

    return sendRequest("DELETE", path);
  },
  bulkDestroy: (ids: number[]) => {
    const { QUERIES } = endpoints;
    const path = `${QUERIES}/delete`;
    return sendRequest("POST", path, { ids });
  },
  load: (id: number) => {
    const { QUERIES } = endpoints;
    const path = `${QUERIES}/${id}`;

    return sendRequest("GET", path);
  },
  loadAll: (teamId?: number) => {
    const { QUERIES } = endpoints;
    const queryString = buildQueryStringFromParams({ team_id: teamId });
    const path = `${QUERIES}`;

    return sendRequest(
      "GET",
      queryString ? path.concat(`?${queryString}`) : path
    );
  },
  run: async ({
    query,
    queryId,
    selected,
  }: {
    query: string;
    queryId: number | null;
    selected: ISelectedTargets;
  }) => {
    const { RUN_QUERY } = endpoints;

    try {
      const { campaign } = await sendRequest("POST", RUN_QUERY, {
        query,
        query_id: queryId,
        selected,
      });
      return Promise.resolve({
        ...campaign,
        hosts_count: {
          successful: 0,
          failed: 0,
          total: 0,
        },
      });
    } catch (response) {
      throw new Error(getError(response as AxiosResponse));
    }
  },
  update: (id: number, updateParams: IQueryFormData) => {
    const { QUERIES } = endpoints;
    const path = `${QUERIES}/${id}`;

    return sendRequest("PATCH", path, updateParams);
  },
};

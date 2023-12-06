import { snakeCase, reduce } from "lodash";

import sendRequest from "services";
import endpoints from "utilities/endpoints";
import {
  ISoftwareResponse,
  ISoftwareCountResponse,
  IGetSoftwareByIdResponse,
} from "interfaces/software";
import { buildQueryStringFromParams, QueryParams } from "utilities/url";
import {
  createMockSoftwareTitle,
  createMockSoftwareTitleReponse,
} from "__mocks__/softwareMock";

interface ISoftwareApiParams {
  page?: number;
  perPage?: number;
  orderKey?: string;
  orderDirection?: "asc" | "desc";
  query?: string;
  vulnerable?: boolean;
  teamId?: number;
}

export interface ISoftwareVersions {
  id: number;
  version: string;
  vulnerabilities: string[] | null;
}

export interface ISoftwareTitle {
  id: number;
  // bundle_identifier: string; TODO: include this?
  name: string;
  versions_count: number;
  source: string;
  hosts_count: number;
  versions: ISoftwareVersions[];
}

export interface ISoftwareTitlesResponse {
  counts_updated_at: string;
  count: number;
  software_titles: ISoftwareTitle[];
  // TODO: include this?
  // meta: {
  //   has_next_results: boolean;
  //   has_previous_results: boolean;
  // }
}

export interface ISoftwareQueryKey extends ISoftwareApiParams {
  scope: "software";
}

export interface ISoftwareCountQueryKey
  extends Pick<ISoftwareApiParams, "query" | "vulnerable" | "teamId"> {
  scope: "softwareCount";
}

const ORDER_KEY = "name";
const ORDER_DIRECTION = "asc";

const convertParamsToSnakeCase = (params: ISoftwareApiParams) => {
  return reduce<typeof params, QueryParams>(
    params,
    (result, val, key) => {
      result[snakeCase(key)] = val;
      return result;
    },
    {}
  );
};

export default {
  load: async ({
    page,
    perPage,
    orderKey = ORDER_KEY,
    orderDirection: orderDir = ORDER_DIRECTION,
    query,
    vulnerable,
    teamId,
  }: ISoftwareApiParams): Promise<ISoftwareResponse> => {
    const { SOFTWARE } = endpoints;
    const queryParams = {
      page,
      perPage,
      orderKey,
      orderDirection: orderDir,
      teamId,
      query,
      vulnerable,
    };

    const snakeCaseParams = convertParamsToSnakeCase(queryParams);
    const queryString = buildQueryStringFromParams(snakeCaseParams);
    const path = `${SOFTWARE}?${queryString}`;

    try {
      return sendRequest("GET", path);
    } catch (error) {
      throw error;
    }
  },

  getCount: async ({
    query,
    teamId,
    vulnerable,
  }: Pick<
    ISoftwareApiParams,
    "query" | "teamId" | "vulnerable"
  >): Promise<ISoftwareCountResponse> => {
    const { SOFTWARE } = endpoints;
    const path = `${SOFTWARE}/count`;
    const queryParams = {
      query,
      teamId,
      vulnerable,
    };
    const snakeCaseParams = convertParamsToSnakeCase(queryParams);
    const queryString = buildQueryStringFromParams(snakeCaseParams);

    return sendRequest("GET", path.concat(`?${queryString}`));
  },

  getSoftwareById: async (
    softwareId: string
  ): Promise<IGetSoftwareByIdResponse> => {
    const { SOFTWARE } = endpoints;
    const path = `${SOFTWARE}/${softwareId}`;

    return sendRequest("GET", path);
  },

  getSoftwareTitles: (params: ISoftwareApiParams) => {
    const { SOFTWARE_TITLES } = endpoints;
    const snakeCaseParams = convertParamsToSnakeCase(params);
    const queryString = buildQueryStringFromParams(snakeCaseParams);
    const path = `${SOFTWARE_TITLES}?${queryString}`;

    // TODO: integrate with API.
    return new Promise<ISoftwareTitlesResponse>((resolve, reject) => {
      resolve(
        createMockSoftwareTitleReponse({
          software_titles: [
            createMockSoftwareTitle({
              versions: [
                { id: 1, version: "1.0.0", vulnerabilities: null },
                { id: 1, version: "1.0.0", vulnerabilities: null },
              ],
            }),
          ],
        })
      );
    });
    // return sendRequest("GET", path);
  },
};

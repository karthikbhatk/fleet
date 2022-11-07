import endpoints from "utilities/endpoints";
import { IActivity } from "interfaces/activity";
import sendRequest from "services";
import { buildQueryStringFromParams } from "utilities/url";

const DEFAULT_PAGE = 0;
const DEFAULT_PAGE_SIZE = 8;
const ORDER_KEY = "created_at";
const ORDER_DIRECTION = "desc";

export interface IActivitiesResponse {
  activities: IActivity[];
}

export default {
  loadNext: (
    page = DEFAULT_PAGE,
    perPage = DEFAULT_PAGE_SIZE
  ): Promise<IActivitiesResponse> => {
    const { ACTIVITIES } = endpoints;

    const queryParams = {
      page,
      per_page: perPage,
      order_key: ORDER_KEY,
      order_direction: ORDER_DIRECTION,
    };

    const queryString = buildQueryStringFromParams(queryParams);

    const path = `${ACTIVITIES}?${queryString}`;

    return sendRequest("GET", path);
  },
};

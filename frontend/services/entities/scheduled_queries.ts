import sendRequest from "services";
import { omit } from "lodash";

import endpoints from "fleet/endpoints";
import {
  IPackQueryFormData,
  IScheduledQuery,
} from "interfaces/scheduled_query";
import { IPack } from "interfaces/pack";
import helpers from "fleet/helpers";

// interface ICreateProps {
//   name: string;
//   description: string;
//   targets: ITargets;
// }

// TODO: Restructure this on the front end passing the data through the form
// logging types and parse ints at the page level -- please see entities/scheduled_queries
export default {
  create: (packQueryFormData: IPackQueryFormData) => {
    const { SCHEDULED_QUERIES } = endpoints;

    return sendRequest("POST", SCHEDULED_QUERIES, packQueryFormData);
  },
  destroy: (packId: number) => {
    const { SCHEDULED_QUERY } = endpoints;
    const path = SCHEDULED_QUERY(packId);

    return sendRequest("DELETE", path);
  },
  loadAll: (packId: number) => {
    const { SCHEDULED_QUERY } = endpoints;
    const path = SCHEDULED_QUERY(packId);

    return sendRequest("GET", path);
  },
  update: (scheduledQuery: IScheduledQuery, updatedAttributes: any) => {
    // TODO: new interface for updated attributes
    const { SCHEDULED_QUERIES } = endpoints;
    const path = `${SCHEDULED_QUERIES}/${scheduledQuery.id}`;
    const params = helpers.formatScheduledQueryForServer(updatedAttributes);

    return sendRequest("PATCH", path, params);
  },
};

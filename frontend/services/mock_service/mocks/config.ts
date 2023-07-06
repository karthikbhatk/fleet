/*
 * NOTE: This is an example of how to configure your mock service.
 * Be sure to copy this file into `../mocks` and only edit that copy!
 * Also please check the README for how to use the mock service :)
 */

import RESPONSES from "./responses";

type IResponses = Record<string, Record<string, Record<string, unknown>>>;

const DELAY = 1000;

const ENDPOINT = "/latest/fleet";

const WILDCARDS: string[] = [":", "*", "{", "}"];

const REQUEST_RESPONSE_MAPPINGS: IResponses = {
  GET: {
    // response is list of all labels excluding any expensive data operations (UI only needs label
    // name and id for this page)
    "labels?summary=true": RESPONSES.labels,
    // request query string is hostname, uuid, or mac address; response is host detail excluding any
    // expensive data operations
    "targets?query={*}": RESPONSES.hosts,
    queries: RESPONSES.queries,
    "queries/1": RESPONSES.query1,
    "queries/2": RESPONSES.query2,
    "queries/3": RESPONSES.query3,
  },
  POST: {
    // request body is ISelectedTargets
    "targets/count": {
      targets_count: 1,
      targets_online: 0,
      targets_offline: 1,
      targets_missing_in_action: 0,
    },
    queries: {
      description: "Ok",
      name: "New query name",
      observer_can_run: false,
      query: "SELECT * FROM osquery_info;",
      team_id: null,
      platform: "linux",
    },
  },
} as IResponses;

export default { DELAY, ENDPOINT, WILDCARDS, REQUEST_RESPONSE_MAPPINGS };

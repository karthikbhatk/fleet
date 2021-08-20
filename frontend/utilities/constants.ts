import URL_PREFIX from "router/url_prefix";

const { origin } = global.window.location;
export const BASE_URL = `${origin}${URL_PREFIX}/api`;

export const FREQUENCY_DROPDOWN_OPTIONS = [
  { value: 3600, label: "Every hour" },
  { value: 21600, label: "Every 6 hours" },
  { value: 43200, label: "Every 12 hours" },
  { value: 86400, label: "Every day" },
  { value: 604800, label: "Every week" },
];

export const PLATFORM_OPTIONS = [
  { label: "All", value: "" },
  { label: "Windows", value: "windows" },
  { label: "Linux", value: "linux" },
  { label: "macOS", value: "darwin" },
];

export const LOGGING_TYPE_OPTIONS = [
  { label: "Snapshot", value: "snapshot" },
  { label: "Differential", value: "differential" },
  {
    label: "Differential (Ignore Removals)",
    value: "differential_ignore_removals",
  },
];

export const MIN_OSQUERY_VERSION_OPTIONS = [
  { label: "All", value: "" },
  { label: "4.7.0 +", value: "4.7.0" },
  { label: "4.6.0 +", value: "4.6.0" },
  { label: "4.5.1 +", value: "4.5.1" },
  { label: "4.5.0 +", value: "4.5.0" },
  { label: "4.4.0 +", value: "4.4.0" },
  { label: "4.3.0 +", value: "4.3.0" },
  { label: "4.2.0 +", value: "4.2.0" },
  { label: "4.1.2 +", value: "4.1.2" },
  { label: "4.1.1 +", value: "4.1.1" },
  { label: "4.1.0 +", value: "4.1.0" },
  { label: "4.0.2 +", value: "4.0.2" },
  { label: "4.0.1 +", value: "4.0.1" },
  { label: "4.0.0 +", value: "4.0.0" },
  { label: "3.4.0 +", value: "3.4.0" },
  { label: "3.3.2 +", value: "3.3.2" },
  { label: "3.3.1 +", value: "3.3.1" },
  { label: "3.2.6 +", value: "3.2.6" },
  { label: "2.2.1 +", value: "2.2.1" },
  { label: "2.2.0 +", value: "2.2.0" },
  { label: "2.1.2 +", value: "2.1.2" },
  { label: "2.1.1 +", value: "2.1.1" },
  { label: "2.0.0 +", value: "2.0.0" },
  { label: "1.8.2 +", value: "1.8.2" },
  { label: "1.8.1 +", value: "1.8.1" },
];

export const QUERIES_PAGE_STEPS = {
  1: "EDITOR",
  2: "TARGETS",
  3: "RUN",
};

export const DEFAULT_QUERY = {
  description: "",
  name: "New query",
  query: "SELECT * FROM osquery_info",
  id: 0,
  interval: 0,
  last_excuted: "",
  observer_can_run: false,
  author_name: "",
  updated_at: "",
};

export const DEFAULT_CAMPAIGN = {
  created_at: "",
  errors: [],
  hosts: [],
  hosts_count: {
    total: 0,
    successful: 0,
    failed: 0,
  },
  id: 0,
  query_id: 0,
  query_results: [],
  status: "",
  totals: {
    count: 0,
    missing_in_action: 0,
    offline: 0,
    online: 0,
  },
  updated_at: "",
  user_id: 0,
};

export const DEFAULT_CAMPAIGN_STATE = {
  observerShowSql: false,
  queryIsRunning: false,
  queryPosition: {},
  queryResultsToggle: null,
  runQueryMilliseconds: 0,
  selectRelatedHostTarget: false,
  targetsCount: 0,
  targetsError: null,
  campaign: { ...DEFAULT_CAMPAIGN },
};

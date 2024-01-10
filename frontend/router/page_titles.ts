// Note: Dynamic page titles are constructed for host, software, query, and policy details on their respective *DetailsPage.tsx file

import { DOCUMENT_TITLE_SUFFIX } from "utilities/constants";
import PATHS from "router/paths";

// Note: Order matters for use of array.find() (specific subpaths must be listed before their parent path)
export default [
  { path: PATHS.DASHBOARD, title: `Dashboard | ${DOCUMENT_TITLE_SUFFIX}` },
  { path: "/hosts/manage", title: `Hosts | ${DOCUMENT_TITLE_SUFFIX}` },
  {
    path: "/controls",
    title: `Controls | ${DOCUMENT_TITLE_SUFFIX}`,
  },
  {
    path: "/software/",
    title: `Software | ${DOCUMENT_TITLE_SUFFIX}`,
  },
  {
    path: PATHS.MANAGE_QUERIES,
    title: `Queries | ${DOCUMENT_TITLE_SUFFIX}`,
  },
  { path: PATHS.NEW_QUERY(), title: `New query | ${DOCUMENT_TITLE_SUFFIX}` },
  {
    path: PATHS.MANAGE_POLICIES,
    title: `Policies | ${DOCUMENT_TITLE_SUFFIX}`,
  },
  { path: PATHS.NEW_POLICY, title: `New policy | ${DOCUMENT_TITLE_SUFFIX}` },
  {
    path: "/settings/",
    title: `Settings | ${DOCUMENT_TITLE_SUFFIX}`,
  },
  {
    path: PATHS.USER_SETTINGS,
    title: `Settings | My account | ${DOCUMENT_TITLE_SUFFIX}`,
  },
];

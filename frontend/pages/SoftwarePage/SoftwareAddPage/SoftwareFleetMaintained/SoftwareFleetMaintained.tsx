import React from "react";
import { InjectedRouter } from "react-router";
import { Location } from "history";
import { useQuery } from "react-query";
import { AxiosError } from "axios";

import softwareAPI, {
  ISoftwareFleetMaintainedAppsQueryParams,
  ISoftwareFleetMaintainedAppsResponse,
} from "services/entities/software";
import { DEFAULT_USE_QUERY_OPTIONS } from "utilities/constants";

import Spinner from "components/Spinner";
import DataError from "components/DataError";

import FleetMaintainedAppsTable from "./FleetMaintainedAppsTable";
import { ISoftwareAddPageQueryParams } from "../SoftwareAddPage";

const baseClass = "software-fleet-maintained";

interface IQueryKey extends ISoftwareFleetMaintainedAppsQueryParams {
  scope?: "fleet-maintained-apps";
}

interface ISoftwareFleetMaintainedProps {
  currentTeamId: number;
  router: InjectedRouter;
  location: Location<ISoftwareAddPageQueryParams>;
}

// default values for query params used on this page if not provided
const DEFAULT_SORT_DIRECTION = "asc";
const DEFAULT_SORT_HEADER = "name";
const DEFAULT_PAGE_SIZE = 20;
const DEFAULT_PAGE = 0;

const SoftwareFleetMaintained = ({
  currentTeamId,
  router,
  location,
}: ISoftwareFleetMaintainedProps) => {
  const {
    order_key = DEFAULT_SORT_HEADER,
    order_direction = DEFAULT_SORT_DIRECTION,
    query = "",
    page,
  } = location.query;
  const currentPage = page ? parseInt(page, 10) : DEFAULT_PAGE;

  const { data, isLoading, isError } = useQuery<
    ISoftwareFleetMaintainedAppsResponse,
    AxiosError,
    ISoftwareFleetMaintainedAppsResponse,
    [IQueryKey]
  >(
    [
      {
        scope: "fleet-maintained-apps",
        page: currentPage,
        per_page: DEFAULT_PAGE_SIZE,
        query,
        order_direction,
        order_key,
        team_id: currentTeamId,
      },
    ],
    (options) => {
      delete options.queryKey[0].scope;
      return softwareAPI.getFleetMaintainedApps(options.queryKey[0]);
    },
    {
      ...DEFAULT_USE_QUERY_OPTIONS,
    }
  );

  if (isLoading) {
    return <Spinner />;
  }

  if (isError) {
    return <DataError className={`${baseClass}__table-error`} />;
  }

  return (
    <div className={baseClass}>
      <FleetMaintainedAppsTable
        data={data}
        isLoading={false}
        router={router}
        query={query}
        teamId={currentTeamId}
        orderDirection={order_direction}
        orderKey={order_key}
        perPage={DEFAULT_PAGE_SIZE}
        currentPage={currentPage}
      />
    </div>
  );
};

export default SoftwareFleetMaintained;

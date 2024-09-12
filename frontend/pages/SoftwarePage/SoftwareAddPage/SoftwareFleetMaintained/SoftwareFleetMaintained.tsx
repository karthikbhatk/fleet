import React from "react";
import { InjectedRouter } from "react-router";
import { Location } from "history";
import { useQuery } from "react-query";

import softwareAPI from "services/entities/software";
import { DEFAULT_USE_QUERY_OPTIONS } from "utilities/constants";

import Spinner from "components/Spinner";
import DataError from "components/DataError";

import FleetMaintainedAppsTable from "./FleetMaintainedAppsTable";
import { ISoftwareAddPageQueryParams } from "../SoftwareAddPage";

const baseClass = "software-fleet-maintained";

interface ISoftwareFleetMaintainedProps {
  currentTeamId: number;
  router: InjectedRouter;
  location: Location<ISoftwareAddPageQueryParams>;
}

// default values for query params used on this page if not provided
const DEFAULT_SORT_DIRECTION = "desc";
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

  const { data, isLoading, isError } = useQuery(
    ["fleet-maintained", currentTeamId],
    () => softwareAPI.getFleetMaintainedApps(currentTeamId),
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

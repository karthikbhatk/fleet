import React, { useEffect, useMemo, useState } from "react";
import { useDebouncedCallback } from "use-debounce";
import { InjectedRouter } from "react-router";
import { Row } from "react-table";
import PATHS from "router/paths";

import { ISoftware } from "interfaces/software";
import { VULNERABLE_DROPDOWN_OPTIONS } from "utilities/constants";
import { buildQueryStringFromParams } from "utilities/url";

// @ts-ignore
import Dropdown from "components/forms/fields/Dropdown";
import TableContainer from "components/TableContainer";

import EmptyState from "../EmptyState";
import SoftwareVulnCount from "./SoftwareVulnCount";

import {
  generateSoftwareTableHeaders,
  generateSoftwareTableData,
} from "./SoftwareTableConfig";

const baseClass = "host-details";

export interface ITableSoftware extends Omit<ISoftware, "vulnerabilities"> {
  vulnerabilities: string[]; // for client-side search purposes, we only want an array of cve strings
}

interface ISoftwareTableProps {
  isLoading: boolean;
  software: ISoftware[];
  deviceUser?: boolean;
  deviceType?: string;
  isSoftwareEnabled?: boolean;
  router?: InjectedRouter;
}

interface IRowProps extends Row {
  original: {
    id?: number;
  };
  isSoftwareEnabled?: boolean;
}

const SoftwareTable = ({
  isLoading,
  software,
  deviceUser,
  deviceType,
  isSoftwareEnabled,
  router,
}: ISoftwareTableProps): JSX.Element => {
  const [searchString, setSearchString] = useState("");
  const [filterVuln, setFilterVuln] = useState(false);
  const [filters, setFilters] = useState({
    global: searchString,
    vulnerabilities: filterVuln,
  });

  useEffect(() => {
    setFilters({ global: searchString, vulnerabilities: filterVuln });
  }, [searchString, filterVuln]);

  const onQueryChange = useDebouncedCallback(
    ({ searchQuery }: { searchQuery: string }) => {
      setSearchString(searchQuery);
    },
    300
  );

  const tableSoftware = useMemo(() => generateSoftwareTableData(software), [
    software,
  ]);
  const tableHeaders = useMemo(
    () => generateSoftwareTableHeaders(deviceUser, router),
    [deviceUser, router]
  );

  const onVulnFilterChange = (value: boolean) => {
    setFilterVuln(value);
  };

  const handleRowSelect = (row: IRowProps) => {
    if (deviceUser || !router) {
      return;
    }

    const queryParams = { software_id: row.original.id };

    const path = queryParams
      ? `${PATHS.MANAGE_HOSTS}?${buildQueryStringFromParams(queryParams)}`
      : PATHS.MANAGE_HOSTS;

    router.push(path);
  };

  const renderVulnFilterDropdown = () => {
    return (
      <Dropdown
        value={filters.vulnerabilities}
        className={`${baseClass}__vuln_dropdown`}
        options={VULNERABLE_DROPDOWN_OPTIONS}
        searchable={false}
        onChange={onVulnFilterChange}
      />
    );
  };

  const EmptySoftwareSearch = () => (
    <EmptyState title="software" reason="empty-search" />
  );

  // if (softwareInventoryEnabled === false) {
  //   return (
  //     <div className="section section--software">
  //       <p className="section__header">Software</p>
  //       <EmptyState title="software" reason="disabled" />
  //     </div>
  //   );
  // }

  // TODO: Refactor this component to show the correct search
  // // TODO: Rework after backend is adjusted to differentiate empty search/filter results from
  // // collecting inventory
  // const isCollectingInventory =
  //   !searchQuery &&
  //   !filterVuln &&
  //   !currentTeam?.id &&
  //   !pageIndex &&
  //   !software?.software &&
  //   software?.counts_updated_at === null;

  // const emptyHeaderText = () => {
  //   if (!isSoftwareEnabled) {
  //     return "Software inventory disabled";
  //   }
  //   if (isCollectingInventory) {
  //     return "No software detected";
  //   }
  //   if (currentTeam && filterVuln) {
  //     return "No vulnerable software detected";
  //   }
  //   return "No software matches the current search criteria";
  // };

  // const emptyInfoText = () => {
  //   if (!isSoftwareEnabled) {
  //     return (
  //       <>
  //         Users with the admin role can{" "}
  //         <CustomLink
  //           url="https://fleetdm.com/docs/using-fleet/vulnerability-processing#configuration"
  //           text="turn on software inventory"
  //           newTab
  //         />
  //         .
  //       </>
  //     );
  //   }
  //   if (isCollectingInventory || (currentTeam && filterVuln)) {
  //     return "This report is updated every hour to protect the performance of your devices.";
  //   }
  //   return "Try again in about 1 hour as the system catches up.";
  // };

  return (
    <div className="section section--software">
      <p className="section__header">Software</p>

      {software?.length ? (
        <>
          {software && (
            <SoftwareVulnCount
              softwareList={software}
              deviceUser={deviceUser}
            />
          )}
          {software && (
            <div className={deviceType || ""}>
              <TableContainer
                columns={tableHeaders}
                data={tableSoftware || []}
                filters={filters}
                isLoading={isLoading}
                defaultSortHeader={"name"}
                defaultSortDirection={"asc"}
                inputPlaceHolder={
                  "Search software by name or vulnerabilities ( CVEs)"
                }
                onQueryChange={onQueryChange}
                resultsTitle={"software items"}
                emptyComponent={EmptySoftwareSearch}
                showMarkAllPages={false}
                isAllPagesSelected={false}
                searchable
                customControl={renderVulnFilterDropdown}
                isClientSidePagination
                pageSize={20}
                isClientSideFilter
                disableMultiRowSelect={!deviceUser && !!router} // device user cannot view hosts by software
                onSelectSingleRow={handleRowSelect}
              />
            </div>
          )}
        </>
      ) : (
        <EmptyState title="software" />
      )}
    </div>
  );
};
export default SoftwareTable;

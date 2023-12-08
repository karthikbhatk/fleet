import React, { useMemo } from "react";
import { InjectedRouter } from "react-router";

import { ISoftwareTitleVersion } from "interfaces/software";

import TableContainer from "components/TableContainer";

import generateSoftwareTitleDetailsTableConfig from "./SoftwareTitleDetailsTableConfig";

const baseClass = "software-title-details-table";

interface ISoftwareTitleDetailsTableProps {
  router: InjectedRouter;
  data: ISoftwareTitleVersion[];
  isLoading: boolean;
}

const SoftwareTitleDetailsTable = ({
  router,
  data,
  isLoading,
}: ISoftwareTitleDetailsTableProps) => {
  const softwareTableHeaders = useMemo(
    () => generateSoftwareTitleDetailsTableConfig(router),
    [router]
  );

  return (
    <TableContainer
      className={baseClass}
      resultsTitle="version" // TODO: dynamic based on number of results
      columns={softwareTableHeaders}
      data={data}
      isLoading={isLoading}
      emptyComponent={() => <p>nothing</p>}
      showMarkAllPages={false}
      isAllPagesSelected={false}
      disablePagination
      // TODO: add row click handler
    />
  );
};

export default SoftwareTitleDetailsTable;

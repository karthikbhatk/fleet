import React from "react";
import { useQuery } from "react-query";

import mdmAPI, { IDiskEncryptionSummaryResponse } from "services/entities/mdm";

import TableContainer from "components/TableContainer";
import EmptyTable from "components/EmptyTable";
import DataError from "components/DataError";

import {
  generateTableHeaders,
  generateTableData,
} from "./DiskEncryptionTableConfig";

const baseClass = "disk-encryption-table";

interface IDiskEncryptionTableProps {
  currentTeamId?: number;
}

const DEFAULT_SORT_HEADER = "hosts";
const DEFAULT_SORT_DIRECTION = "asc";

const DiskEncryptionTable = ({ currentTeamId }: IDiskEncryptionTableProps) => {
  const {
    data: diskEncryptionStatusData,
    error: diskEncryptionStatusError,
  } = useQuery<IDiskEncryptionSummaryResponse, Error>(
    ["disk-encryption-summary", currentTeamId],
    () => mdmAPI.getDiskEncryptionSummary(currentTeamId),
    {
      refetchOnWindowFocus: false,
      retry: false,
    }
  );

  const tableHeaders = generateTableHeaders();

  const tableData = generateTableData(diskEncryptionStatusData, currentTeamId);

  if (diskEncryptionStatusError) {
    return <DataError />;
  }

  if (!diskEncryptionStatusData) return null;

  return (
    <div className={baseClass}>
      <TableContainer
        columns={tableHeaders}
        data={tableData}
        resultsTitle="" // TODO: make optional
        isLoading={false}
        showMarkAllPages={false}
        isAllPagesSelected={false}
        manualSortBy
        // defaultSortHeader={DEFAULT_SORT_HEADER}
        // defaultSortDirection={DEFAULT_SORT_DIRECTION}
        disableTableHeader
        disablePagination
        disableCount
        emptyComponent={() => (
          <EmptyTable
            header="No Disk Encryption Status"
            info="Expecting to status data? Try again in a few seconds as the system
              catches up."
          />
        )}
      />
    </div>
  );
};

export default DiskEncryptionTable;

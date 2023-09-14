import React from "react";

import { DiskEncryptionStatus } from "interfaces/mdm";
import {
  IDiskEncryptionStatusAggregate,
  IDiskEncryptionSummaryResponse,
} from "services/entities/mdm";

import TextCell from "components/TableContainer/DataTable/TextCell";
import HeaderCell from "components/TableContainer/DataTable/HeaderCell";
import StatusIndicatorWithIcon from "components/StatusIndicatorWithIcon";
import ViewAllHostsLink from "components/ViewAllHostsLink";
import { IndicatorStatus } from "components/StatusIndicatorWithIcon/StatusIndicatorWithIcon";

interface IStatusCellValue {
  displayName: string;
  statusName: IndicatorStatus;
  value: DiskEncryptionStatus;
  tooltip?: string | JSX.Element;
}

interface IStatusCellProps {
  cell: {
    value: IStatusCellValue;
  };
}

interface ICellProps {
  cell: {
    value: string;
  };
  row: {
    original: {
      status: IStatusCellValue;
      teamId: number;
    };
  };
}

interface IHeaderProps {
  column: {
    title: string;
    isSortedDesc: boolean;
  };
}

type IDataColumn = {
  title: string;
  Header: ((props: IHeaderProps) => JSX.Element) | string;
  accessor: string;
  disableHidden?: boolean;
  disableSortBy?: boolean;
  Cell:
    | ((props: ICellProps) => JSX.Element)
    | ((props: IStatusCellProps) => JSX.Element);
};

const defaultTableHeaders: IDataColumn[] = [
  {
    title: "Status",
    Header: "Status",
    disableSortBy: true,
    accessor: "status",
    Cell: ({ cell: { value } }: IStatusCellProps) => {
      const tooltipProp = value.tooltip
        ? { tooltipText: value.tooltip }
        : undefined;
      return (
        <StatusIndicatorWithIcon
          status={value.statusName}
          value={value.displayName}
          tooltip={tooltipProp}
        />
      );
    },
  },
  {
    title: "macOS hosts",
    Header: (cellProps: IHeaderProps) => (
      <HeaderCell
        value={cellProps.column.title}
        isSortedDesc={cellProps.column.isSortedDesc}
        disableSortBy={false}
      />
    ),
    accessor: "macosHosts",
    Cell: ({ cell: { value: aggregateCount } }: ICellProps) => {
      return (
        <div className="disk-encryption-table__aggregate-table-data">
          <TextCell value={aggregateCount} formatter={(val) => <>{val}</>} />
        </div>
      );
    },
  },
  {
    title: "Windows hosts",
    Header: (cellProps: IHeaderProps) => (
      <HeaderCell
        value={cellProps.column.title}
        isSortedDesc={cellProps.column.isSortedDesc}
        disableSortBy={false}
      />
    ),
    accessor: "windowsHosts",
    Cell: ({
      cell: { value: aggregateCount },
      row: { original },
    }: ICellProps) => {
      return (
        <div className="disk-encryption-table__aggregate-table-data">
          <TextCell value={aggregateCount} formatter={(val) => <>{val}</>} />
          <ViewAllHostsLink
            className="view-hosts-link"
            queryParams={{
              macos_settings_disk_encryption: original.status.value,
              team_id: original.teamId,
            }}
          />
        </div>
      );
    },
  },
];

export const generateTableHeaders = (): IDataColumn[] => {
  return defaultTableHeaders;
};

const STATUS_CELL_VALUES: Record<DiskEncryptionStatus, IStatusCellValue> = {
  verified: {
    displayName: "Verified",
    statusName: "success",
    value: "verified",
    tooltip:
      "These hosts turned disk encryption on and sent their key to Fleet. Fleet verified with osquery.",
  },
  verifying: {
    displayName: "Verifying",
    statusName: "successPartial",
    value: "verifying",
    tooltip:
      "These hosts acknowledged the MDM command to turn on disk encryption. Fleet is verifying with " +
      "osquery and retrieving the disk encryption key.This may take up to one hour.",
  },
  action_required: {
    displayName: "Action required (pending)",
    statusName: "pendingPartial",
    value: "action_required",
    tooltip: (
      <>
        Ask the end user to follow <b>Disk encryption</b> instructions on their{" "}
        <b>My device</b> page.
      </>
    ),
  },
  enforcing: {
    displayName: "Enforcing (pending)",
    statusName: "pendingPartial",
    value: "enforcing",
    tooltip:
      "These hosts will receive the MDM command to turn on the disk encryption when the hosts come online.",
  },
  failed: {
    displayName: "Failed",
    statusName: "error",
    value: "failed",
  },
  removing_enforcement: {
    displayName: "Removing enforcement (pending)",
    statusName: "pendingPartial",
    value: "removing_enforcement",
    tooltip:
      "These hosts will receive the MDM command to turn off the disk encryption when the hosts come online.",
  },
};

type StatusEntry = [DiskEncryptionStatus, IDiskEncryptionStatusAggregate];

export const generateTableData = (
  data?: IDiskEncryptionSummaryResponse,
  currentTeamId?: number
) => {
  if (!data) return [];

  // type cast here gives a little more typing information than Object.entries alone.
  const entries = Object.entries(data) as StatusEntry[];

  return entries.map(([status, statusAggregate]) => ({
    status: STATUS_CELL_VALUES[status],
    macosHosts: statusAggregate.macos,
    windowsHosts: statusAggregate.windows,
    teamId: currentTeamId,
  }));
};

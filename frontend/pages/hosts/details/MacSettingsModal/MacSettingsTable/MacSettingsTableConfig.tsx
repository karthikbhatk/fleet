import TextCell from "components/TableContainer/DataTable/TextCell";
import React from "react";
import { IMacMdmProfile } from "interfaces/mdm";
import MacSettingsIndicator from "../../MacSettingsIndicator";

interface IHeaderProps {
  column: {
    title: string;
    isSortedDesc: boolean;
  };
}

interface ICellProps {
  cell: {
    value: string;
  };
  row: {
    original: IMacMdmProfile;
  };
}

interface IDataColumn {
  Header: ((props: IHeaderProps) => JSX.Element) | string;
  Cell: (props: ICellProps) => JSX.Element;
  id?: string;
  title?: string;
  accessor?: string;
  disableHidden?: boolean;
  disableSortBy?: boolean;
  sortType?: string;
}

const getStatusDisplayOptions = (
  profile: IMacMdmProfile
): {
  statusText: string;
  iconName: "pending" | "success" | "error";
  tooltipText: string | null;
} => {
  const SETTING_STATUS_OPTIONS = {
    pending: {
      "Action required":
        "Follow Disk encryption instructions on your My device page.",
      Enforcing: "Setting will be enforced when the host comes online.",
      "Removing enforcement":
        "Enforcement will be removed when the host comes online.",
      "": "",
    },
    applied: {
      iconName: "success",
      tooltipText:
        "Disk encryption on and disk encryption key stored in Fleet.",
    },
    failed: { iconName: "error", tooltipText: null },
  } as const;

  if (profile.status === "pending") {
    return {
      statusText: `${profile.detail} (pending)`,
      iconName: "pending",
      tooltipText: SETTING_STATUS_OPTIONS.pending[profile.detail],
    };
  }
  return {
    statusText:
      profile.status.charAt(0).toUpperCase() + profile.status.slice(1),
    iconName: SETTING_STATUS_OPTIONS[profile.status].iconName,
    tooltipText: SETTING_STATUS_OPTIONS[profile.status].tooltipText,
  };
};

const tableHeaders: IDataColumn[] = [
  {
    title: "Name",
    Header: "Name",
    disableSortBy: true,
    accessor: "name",
    Cell: (cellProps: ICellProps): JSX.Element => (
      <TextCell value={cellProps.cell.value} />
    ),
  },
  {
    title: "Status",
    Header: "Status",
    disableSortBy: true,
    accessor: "statusText",
    Cell: (cellProps: ICellProps) => {
      const { statusText, iconName, tooltipText } = getStatusDisplayOptions(
        cellProps.row.original
      );
      return (
        <MacSettingsIndicator
          indicatorText={statusText}
          iconName={iconName}
          tooltip={{ tooltipText }}
        />
      );
    },
  },
  {
    title: "Error",
    Header: "Error",
    disableSortBy: true,
    accessor: "error",
    Cell: (cellProps: ICellProps): JSX.Element => (
      <TextCell value={cellProps.row.original.error} />
    ),
  },
];

export default tableHeaders;

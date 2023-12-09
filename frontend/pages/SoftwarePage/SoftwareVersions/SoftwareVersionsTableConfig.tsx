import React from "react";
import { Column } from "react-table";
import { InjectedRouter } from "react-router";

import {
  formatSoftwareType,
  ISoftwareVersion,
  ISoftwareVulnerability,
} from "interfaces/software";
import PATHS from "router/paths";

import HeaderCell from "components/TableContainer/DataTable/HeaderCell";
import TextCell from "components/TableContainer/DataTable/TextCell";
import LinkCell from "components/TableContainer/DataTable/LinkCell/LinkCell";
import ViewAllHostsLink from "components/ViewAllHostsLink";
import VulnerabilitiesCell from "../components/VulnerabilitiesCell";

// NOTE: cellProps come from react-table
// more info here https://react-table.tanstack.com/docs/api/useTable#cell-properties
interface ICellProps {
  cell: {
    value: number | string | ISoftwareVulnerability[];
  };
  row: {
    original: ISoftwareVersion;
  };
}
interface IStringCellProps extends ICellProps {
  cell: {
    value: string;
  };
}

interface IVersionCellProps extends ICellProps {
  cell: {
    value: string;
  };
}

interface INumberCellProps extends ICellProps {
  cell: {
    value: number;
  };
}

interface IVulnCellProps extends ICellProps {
  cell: {
    value: ISoftwareVulnerability[];
  };
}
interface IHeaderProps {
  column: {
    title: string;
    isSortedDesc: boolean;
  };
}

const generateTableHeaders = (
  router: InjectedRouter,
  isPremiumTier?: boolean,
  isSandboxMode?: boolean,
  teamId?: number
): Column[] => {
  const softwareTableHeaders = [
    {
      title: "Name",
      Header: (cellProps: IHeaderProps): JSX.Element => (
        <HeaderCell
          value={cellProps.column.title}
          isSortedDesc={cellProps.column.isSortedDesc}
        />
      ),
      disableSortBy: false,
      accessor: "name",
      Cell: (cellProps: IStringCellProps): JSX.Element => {
        const { id, name } = cellProps.row.original;

        const onClickSoftware = (e: React.MouseEvent) => {
          // Allows for button to be clickable in a clickable row
          e.stopPropagation();
          router?.push(PATHS.SOFTWARE_VERSION_DETAILS(id.toString()));
        };

        return (
          <LinkCell
            path={PATHS.SOFTWARE_VERSION_DETAILS(id.toString())}
            customOnClick={onClickSoftware}
            value={name}
          />
        );
      },
      sortType: "caseInsensitive",
    },
    {
      title: "Version",
      Header: "Version",
      disableSortBy: true,
      accessor: "version",
      Cell: (cellProps: IVersionCellProps): JSX.Element => (
        <TextCell value={cellProps.cell.value} />
      ),
    },
    {
      title: "Type",
      Header: "Type",
      disableSortBy: true,
      accessor: "source",
      Cell: (cellProps: IStringCellProps): JSX.Element => (
        <TextCell formatter={formatSoftwareType} value={cellProps.cell.value} />
      ),
    },
    {
      title: "Vulnerabilities",
      Header: (cellProps: IHeaderProps): JSX.Element => (
        <HeaderCell
          value={cellProps.column.title}
          disableSortBy={false}
          isSortedDesc={cellProps.column.isSortedDesc}
        />
      ),
      disableSortBy: false,
      // the "vulnerabilities" accessor is used but the data is actually coming
      // from the version attribute. We do this as we already have a "versions"
      // attribute used for the "Version" column and we cannot reuse. This is a
      // limitation of react-table.
      // With the versions data, we can sum up the vulnerabilities to get the
      // total number of vulnerabilities for the software title
      accessor: "vulnerabilities",
      Cell: (cellProps: IVulnCellProps): JSX.Element => (
        <VulnerabilitiesCell vulnerabilities={cellProps.cell.value} />
      ),
    },
    {
      title: "Hosts",
      Header: (cellProps: IHeaderProps): JSX.Element => (
        <HeaderCell
          value={cellProps.column.title}
          disableSortBy={false}
          isSortedDesc={cellProps.column.isSortedDesc}
        />
      ),
      disableSortBy: false,
      accessor: "hosts_count",
      Cell: (cellProps: INumberCellProps): JSX.Element => (
        <TextCell value={cellProps.cell.value} />
      ),
    },
  ];

  return softwareTableHeaders;
};

export default generateTableHeaders;

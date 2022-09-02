import React from "react";
import { capitalize } from "lodash";
import ReactTooltip from "react-tooltip";

import formatDistanceToNowStrict from "date-fns/formatDistanceToNowStrict";
import { abbreviateTimeUnits } from "utilities/helpers";

import PATHS from "router/paths";
import HeaderCell from "components/TableContainer/DataTable/HeaderCell/HeaderCell";
import TextCell from "components/TableContainer/DataTable/TextCell";
import TruncatedTextCell from "components/TableContainer/DataTable/TruncatedTextCell";
import TooltipWrapper from "components/TooltipWrapper";
import { IMunkiIssue } from "interfaces/host";

interface IHeaderProps {
  column: {
    title: string;
    isSortedDesc: boolean;
  };
}
interface ICellProps {
  cell: {
    value: number | string | string[];
  };
  row: {
    original: IMunkiIssue;
    index: number;
  };
}

interface IStringCellProps extends ICellProps {
  cell: {
    value: string;
  };
}

interface IDataColumn {
  title: string;
  Header: ((props: IHeaderProps) => JSX.Element) | string;
  accessor: string;
  Cell: (props: IStringCellProps) => JSX.Element;
  disableHidden?: boolean;
  disableSortBy?: boolean;
  disableGlobalFilter?: boolean;
  sortType?: string;
  // Filter can be used by react-table to render a filter input inside the column header
  Filter?: () => null | JSX.Element;
  filter?: string; // one of the enumerated `filterTypes` for react-table
  // (see https://github.com/tannerlinsley/react-table/blob/master/src/filterTypes.js)
  // or one of the custom `filterTypes` defined for the `useTable` instance (see `DataTable`)
}

interface IMunkiIssueTableData extends IMunkiIssue {
  time: string;
}

export const generateMunkiIssuesTableData = (
  munkiIssues: IMunkiIssue[] | undefined
): IMunkiIssueTableData[] => {
  if (!munkiIssues) {
    return [];
  }
  return munkiIssues.map((i) => {
    return {
      ...i,
      time: i.created_at,
    };
  });
};

// NOTE: cellProps come from react-table
// more info here https://react-table.tanstack.com/docs/api/useTable#cell-properties
export const generateMunkiIssuesTableHeaders = (): IDataColumn[] => {
  const tableHeaders: IDataColumn[] = [
    {
      title: "Issue",
      Header: (headerProps: IHeaderProps): JSX.Element => {
        const titleWithToolTip = (
          <TooltipWrapper
            tipContent={`
            Issues reported the last time Munki ran on each host.
          `}
          >
            Issue
          </TooltipWrapper>
        );
        return (
          <HeaderCell
            value={titleWithToolTip}
            isSortedDesc={headerProps.column.isSortedDesc}
          />
        );
      },
      disableSortBy: false,
      accessor: "name",
      Cell: (cellProps: IStringCellProps) => (
        <TruncatedTextCell value={cellProps.cell.value} />
      ),
      sortType: "caseInsensitive",
    },
    {
      title: "Type",
      Header: "Type",
      disableSortBy: true,
      accessor: "type",
      Cell: (cellProps: IStringCellProps) => (
        <TextCell value={capitalize(cellProps.cell.value)} />
      ),
    },
    {
      title: "Time",
      Header: (headerProps: IHeaderProps): JSX.Element => {
        const titleWithToolTip = (
          <TooltipWrapper
            tipContent={`
            The first time Munki reported this issue.
          `}
          >
            Time
          </TooltipWrapper>
        );
        return (
          <HeaderCell
            value={titleWithToolTip}
            isSortedDesc={headerProps.column.isSortedDesc}
          />
        );
      },
      disableSortBy: false,
      accessor: "time",
      Cell: (cellProps: IStringCellProps) => {
        const time = abbreviateTimeUnits(
          formatDistanceToNowStrict(new Date(cellProps.cell.value), {
            addSuffix: true,
          })
        );

        return <TextCell value={time} />;
      },
    },
  ];

  return tableHeaders;
};

export default {
  generateMunkiIssuesTableHeaders,
  generateMunkiIssuesTableData,
};

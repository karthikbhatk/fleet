import React, { useState, useEffect } from "react";
import { Dispatch } from "redux";
import { push } from "react-router-redux";
import { Tab, Tabs, TabList, TabPanel } from "react-tabs";
import moment from "moment";
import classnames from "classnames";
import FileSaver from "file-saver";
import { filter, get, keys, omit } from "lodash";

import PATHS from "router/paths"; // @ts-ignore
import convertToCSV from "utilities/convert_to_csv"; // @ts-ignore
import filterArrayByHash from "utilities/filter_array_by_hash";
import { ICampaign, ICampaignQueryResult } from "interfaces/campaign";

import Button from "components/buttons/Button"; // @ts-ignore
import FleetIcon from "components/icons/FleetIcon"; // @ts-ignore
import InputField from "components/forms/fields/InputField";
import QueryResultsRow from "components/queries/QueryResultsRow";
import Spinner from "components/loaders/Spinner";

interface IQueryResultsProps {
  campaign: ICampaign;
  isQueryFinished: boolean;
  onRunQuery: (evt: React.MouseEvent<HTMLButtonElement>) => void;
  onStopQuery: (evt: React.MouseEvent<HTMLButtonElement>) => void;
  dispatch: Dispatch;
}

const baseClass = "query-results";
const CSV_QUERY_TITLE = "Query Results";
const PAGE_TITLES = {
  RUNNING: "Querying selected hosts",
  FINISHED: "Query finished",
};
const NAV_TITLES = {
  RESULTS: "Results",
  ERRORS: "Errors",
};

const QueryResults = ({
  campaign,
  isQueryFinished,
  onRunQuery,
  onStopQuery,
  dispatch,
}: IQueryResultsProps) => {
  const { hosts_count: hostsCount, query_results: queryResults, errors } =
    campaign || {};

  const totalHostsOnline = get(campaign, ["totals", "online"], 0);
  const totalHostsOffline = get(campaign, ["totals", "offline"], 0);
  const totalRowsCount = get(campaign, ["query_results", "length"], 0);
  const onlineTotalText = `${totalRowsCount} result${
    totalRowsCount === 1 ? "" : "s"
  }`;
  const errorsTotalText = `${errors?.length || 0} result${
    errors?.length === 1 ? "" : "s"
  }`;

  const [pageTitle, setPageTitle] = useState<string>(PAGE_TITLES.RUNNING);
  const [resultsFilter, setResultsFilter] = useState<{ [key: string]: string }>(
    {}
  );
  const [activeColumn, setActiveColumn] = useState<string>("");
  const [navTabIndex, setNavTabIndex] = useState(0);

  useEffect(() => {
    if (isQueryFinished) {
      setPageTitle(PAGE_TITLES.FINISHED);
    } else {
      setPageTitle(PAGE_TITLES.RUNNING);
    }
  }, [isQueryFinished]);

  const onFilterAttribute = (attribute: string) => {
    return (value: string) => {
      setResultsFilter({
        ...resultsFilter,
        [attribute]: value,
      });

      return false;
    };
  };

  const onExportQueryResults = (evt: React.MouseEvent<HTMLButtonElement>) => {
    evt.preventDefault();

    if (queryResults) {
      const csv = convertToCSV(queryResults, (fields: string[]) => {
        const result = filter(fields, (f) => f !== "host_hostname");
        result.unshift("host_hostname");

        return result;
      });

      const formattedTime = moment(new Date()).format("MM-DD-YY hh-mm-ss");
      const filename = `${CSV_QUERY_TITLE} (${formattedTime}).csv`;
      const file = new global.window.File([csv], filename, {
        type: "text/csv",
      });

      FileSaver.saveAs(file);
    }
  };

  const onExportErrorsResults = (evt: React.MouseEvent<HTMLButtonElement>) => {
    evt.preventDefault();

    if (errors) {
      const csv = convertToCSV(errors, (fields: string[]) => {
        const result = filter(fields, (f) => f !== "host_hostname");
        result.unshift("host_hostname");

        return result;
      });

      const formattedTime = moment(new Date()).format("MM-DD-YY hh-mm-ss");
      const filename = `${CSV_QUERY_TITLE} Errors (${formattedTime}).csv`;
      const file = new global.window.File([csv], filename, {
        type: "text/csv",
      });

      FileSaver.saveAs(file);
    }
  };

  const renderTableHeaderColumn = (column: string, index: number) => {
    const filterable = column === "hostname" ? "host_hostname" : column;
    const filterIconClassName = classnames(`${baseClass}__filter-icon`, {
      [`${baseClass}__filter-icon--is-active`]: activeColumn === column,
    });

    return (
      <th key={`query-results-table-header-${index}`}>
        <span>
          <FleetIcon className={filterIconClassName} name="filter" />
          {column}
        </span>
        <InputField
          name={column}
          onChange={onFilterAttribute(filterable)}
          onFocus={() => setActiveColumn(column)}
          value={resultsFilter[filterable]}
        />
      </th>
    );
  };

  const renderTableHeaderRow = (rows: ICampaignQueryResult[]) => {
    if (!rows) {
      return false;
    }

    const queryAttrs = omit(rows[0], ["host_hostname"]);
    const queryResultColumns = keys(queryAttrs);

    return (
      <tr>
        {renderTableHeaderColumn("hostname", -1)}
        {queryResultColumns.map((column, i) => {
          return renderTableHeaderColumn(column, i);
        })}
      </tr>
    );
  };

  const renderTableRows = (rows: ICampaignQueryResult[]) => {
    const filteredRows = filterArrayByHash(rows, resultsFilter);

    return filteredRows.map((row: ICampaignQueryResult) => {
      return (
        <QueryResultsRow
          key={row.uuid || row.host_hostname}
          queryResult={row}
        />
      );
    });
  };

  const renderTable = () => {
    return (
      <div className={`${baseClass}__results-table-container`}>
        <Button
          className={`${baseClass}__export-btn`}
          onClick={onExportQueryResults}
          variant="text-link"
        >
          Export results
        </Button>
        <div className={`${baseClass}__results-table-wrapper`}>
          <table className={`${baseClass}__table`}>
            <thead>{renderTableHeaderRow(queryResults)}</thead>
            <tbody>{renderTableRows(queryResults)}</tbody>
          </table>
        </div>
      </div>
    );
  };

  const renderErrorsTable = () => {
    return (
      <div className={`${baseClass}__error-table-container`}>
        <Button
          className={`${baseClass}__export-btn`}
          onClick={onExportErrorsResults}
          variant="text-link"
        >
          Export errors
        </Button>
        <div className={`${baseClass}__error-table-wrapper`}>
          <table className={`${baseClass}__table`}>
            <thead>{renderTableHeaderRow(errors)}</thead>
            <tbody>{renderTableRows(errors)}</tbody>
          </table>
        </div>
      </div>
    );
  };

  const renderFinishedButtons = () => (
    <div className={`${baseClass}__btn-wrapper`}>
      <Button
        className={`${baseClass}__done-btn`}
        onClick={() => dispatch(push(PATHS.MANAGE_QUERIES))}
        variant="brand"
      >
        Done
      </Button>
      <Button
        className={`${baseClass}__run-btn`}
        onClick={onRunQuery}
        variant="blue-green"
      >
        Run again
      </Button>
    </div>
  );

  const renderStopQueryButton = () => (
    <div className={`${baseClass}__btn-wrapper`}>
      <Button
        className={`${baseClass}__stop-btn`}
        onClick={onStopQuery}
        variant="alert"
      >
        <>
          <Spinner isInButton />
          Stop
        </>
      </Button>
    </div>
  );

  const hasNoResults =
    isQueryFinished &&
    (!hostsCount.successful || !queryResults || !queryResults.length);
  return (
    <div className={baseClass}>
      <div className={`${baseClass}__wrapper`}>
        <h1>{pageTitle}</h1>
        <div className={`${baseClass}__text-wrapper`}>
          <span className={`${baseClass}__text-online`}>
            Online: {totalHostsOnline} hosts / {onlineTotalText}
          </span>
          <span className={`${baseClass}__text-offline`}>
            Offline: {totalHostsOffline} hosts / 0 results
          </span>
          <span className={`${baseClass}__text-error`}>
            Errors: {hostsCount.failed} hosts / {errorsTotalText}
          </span>
        </div>
      </div>
      {isQueryFinished ? renderFinishedButtons() : renderStopQueryButton()}
      <div className={`${baseClass}__nav-header`}>
        <Tabs selectedIndex={navTabIndex} onSelect={(i) => setNavTabIndex(i)}>
          <TabList>
            <Tab className="react-tabs__tab no-count">{NAV_TITLES.RESULTS}</Tab>
            <Tab disabled={!errors?.length}>
              <span className="count">{errors?.length || 0}</span>
              {NAV_TITLES.ERRORS}
            </Tab>
          </TabList>
          <TabPanel>
            {isQueryFinished && hasNoResults ? (
              <p className="no-results-message">
                Your live query returned no results.
                <br />
                <span>
                  Expecting to see results? Check to see if the hosts you
                  targeted reported &ldquo;Online&rdquo; or check out the
                  &ldquo;Errors&rdquo; table.
                </span>
              </p>
            ) : (
              renderTable()
            )}
          </TabPanel>
          <TabPanel>{renderErrorsTable()}</TabPanel>
        </Tabs>
      </div>
    </div>
  );
};

export default QueryResults;

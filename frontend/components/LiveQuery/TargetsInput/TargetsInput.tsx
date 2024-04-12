import React from "react";
import { Row } from "react-table";
import { isEmpty, pullAllBy } from "lodash";

import { IHost } from "interfaces/host";
import { HOSTS_SEARCH_BOX_PLACEHOLDER } from "utilities/constants";

import DataError from "components/DataError";
// @ts-ignore
import InputFieldWithIcon from "components/forms/fields/InputFieldWithIcon/InputFieldWithIcon";
import TableContainer from "components/TableContainer";
import { generateTableHeaders } from "./TargetsInputHostsTableConfig";

interface ITargetsInputProps {
  tabIndex?: number;
  searchText: string;
  searchResults: IHost[];
  isTargetsLoading: boolean;
  hasFetchError: boolean;
  targetedHosts: IHost[];
  searchResultsTableConfig: any; // TODO: add typing;
  selectedHostsTableConifg: any; // TODO: add typing;
  label?: string;
  placeholder?: string;
  setSearchText: (value: string) => void;
  handleRowSelect: (value: Row<IHost>) => void;
}

const baseClass = "targets-input";

const DEFAULT_LABEL = "Target specific hosts";

const TargetsInput = ({
  tabIndex,
  searchText,
  searchResults,
  isTargetsLoading,
  hasFetchError,
  targetedHosts,
  searchResultsTableConfig,
  selectedHostsTableConifg,
  label = DEFAULT_LABEL,
  placeholder = HOSTS_SEARCH_BOX_PLACEHOLDER,
  handleRowSelect,
  setSearchText,
}: ITargetsInputProps): JSX.Element => {
  const dropdownHosts =
    searchResults && pullAllBy(searchResults, targetedHosts, "display_name");
  const isActiveSearch =
    !isEmpty(searchText) && (!hasFetchError || isTargetsLoading);
  const isSearchError = !isEmpty(searchText) && hasFetchError;

  return (
    <div>
      <div className={baseClass}>
        <InputFieldWithIcon
          autofocus
          type="search"
          iconSvg="search"
          value={searchText}
          tabIndex={tabIndex}
          iconPosition="start"
          label={label}
          placeholder={placeholder}
          onChange={setSearchText}
        />
        {isActiveSearch && (
          <div className={`${baseClass}__hosts-search-dropdown`}>
            <TableContainer<Row<IHost>>
              columnConfigs={searchResultsTableConfig}
              data={dropdownHosts}
              isLoading={isTargetsLoading}
              resultsTitle=""
              emptyComponent={() => (
                <div className="empty-search">
                  <div className="empty-search__inner">
                    <h4>No hosts match the current search criteria.</h4>
                    <p>
                      Expecting to see hosts? Try again in a few seconds as the
                      system catches up.
                    </p>
                  </div>
                </div>
              )}
              showMarkAllPages={false}
              isAllPagesSelected={false}
              disableCount
              disablePagination
              disableMultiRowSelect
              onClickRow={handleRowSelect}
              // onSelectSingleRow={handleRowSelect}
            />
          </div>
        )}
        {isSearchError && (
          <div className={`${baseClass}__hosts-search-dropdown`}>
            <DataError />
          </div>
        )}
        <div className={`${baseClass}__hosts-selected-table`}>
          <TableContainer
            columnConfigs={selectedHostsTableConifg}
            data={targetedHosts}
            isLoading={false}
            resultsTitle=""
            showMarkAllPages={false}
            isAllPagesSelected={false}
            disableCount
            disablePagination
            emptyComponent={() => <></>}
          />
        </div>
      </div>
    </div>
  );
};

export default TargetsInput;

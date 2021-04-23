import React, { useState } from "react";
import PropTypes from "prop-types";

import Button from "components/buttons/Button";
import Modal from "components/modals/Modal";
import KolideIcon from "components/icons/KolideIcon";
import InputField from "components/forms/fields/InputField";

import queryInterface from "interfaces/query";
import hostInterface from "interfaces/host";

import helpers from "../helpers";

import OpenNewTabIcon from "../../../../../assets/images/open-new-tab-12x12@2x.png";
import ErrorIcon from "../../../../../assets/images/icon-error-16x16@2x.png";

const baseClass = "select-query-modal";

const onQueryHostCustom = (host, dispatch) => {
  const { queryHostCustom } = helpers;

  queryHostCustom(dispatch, host);

  return false;
};

const onQueryHostSaved = (host, selectedQuery, dispatch) => {
  const { queryHostSaved } = helpers;

  queryHostSaved(dispatch, host, selectedQuery);

  return false;
};

const SelectQueryModal = (props) => {
  const { host, toggleQueryHostModal, queries, dispatch } = props;

  const [queriesFilter, setQueriesFilter] = useState("");

  const onFilterQueries = () => {
    setQueriesFilter(queriesFilter);

    return false;
  };

  const queriesCount = queries.length;
  const disabled = !queriesFilter && queriesCount === 0;

  const results = () => {
    // if (errorStateHere) {
    //   return (
    //     <div className={`${baseClass}__no-query-results`}>
    //       <span className="info__header">
    //         <img src={ErrorIcon} alt="error icon" id="error-icon" />
    //         Something's gone wrong.
    //       </span>
    //       <span className="info__data">
    //         Refresh the page or log in again.
    //       </span>
    //       <span className="info__data">
    //         If this keeps happening, please{" "}
    //         <a
    //           href="https://github.com/fleetdm/fleet/issues"
    //           target="_blank"
    //           rel="noopener noreferrer"
    //         >
    //           file an issue.
    //           <img
    //             src={OpenNewTabIcon}
    //             alt="open new tab"
    //             id="new-tab-icon"
    //           />
    //         </a>
    //       </span>
    //     </div>
    //   );
    // }

    if (queriesCount > 0) {
      const queryList = queries.map((query) => {
        return (
          <Button
            key={query.id}
            variant="unstyled-modal-query"
            className="modal-query-button"
            onClick={() => onQueryHostSaved(host, query, dispatch)}
          >
            <span className="info__header">{query.name}</span>
            <span className="info__data">{query.description}</span>
          </Button>
        );
      });
      return <div>{queryList}</div>;
    }

    if (!queriesFilter && queriesCount === 0) {
      return (
        <div className={`${baseClass}__no-query-results`}>
          <span className="info__header">You have no saved queries.</span>
          <span className="info__data">
            Expecting to see queries? Try again in a few seconds as the system
            catches up.
          </span>
        </div>
      );
    }

    if (queriesFilter && queriesCount === 0) {
      return (
        <div className={`${baseClass}__no-query-results`}>
          <span className="info__header">
            No queries match the current search criteria.
          </span>
          <span className="info__data">
            Expecting to see queries? Try again in a few seconds as the system
            catches up.
          </span>
        </div>
      );
    }
  };

  return (
    <Modal
      title="Select a query"
      onExit={toggleQueryHostModal(null)}
      className={`${baseClass}__modal`}
    >
      <div className={`${baseClass}__filter-queries`}>
        <InputField
          name="query-filter"
          onChange={onFilterQueries}
          placeholder="Filter queries"
          value={queriesFilter}
          disabled={disabled}
        />
        <KolideIcon name="search" />
      </div>
      {results()}
      <p>
        <Button
          onClick={() => onQueryHostCustom(host, dispatch)}
          variant="unstyled"
          className={`${baseClass}__custom-query-button`}
        >
          Custom Query
        </Button>
      </p>
    </Modal>
  );
};

SelectQueryModal.propTypes = {
  dispatch: PropTypes.func,
  host: hostInterface,
  queries: PropTypes.arrayOf(queryInterface),
  toggleQueryHostModal: PropTypes.func,
};

export default SelectQueryModal;

import React, { useState, useCallback } from "react"; //, { useEffect }
import { useDispatch, useSelector } from "react-redux";

import { push } from "react-router-redux";
// @ts-ignore
// import { IUser } from "interfaces/user";
import { IQuery } from "interfaces/query";

// Will I need this? 5/28
// import permissionUtils from "utilities/permissions";
import paths from "router/paths";
import Button from "components/buttons/Button";
import NoSchedule from "./components/NoSchedule";
import ScheduleError from "./components/ScheduleError";
import ScheduleEditorModal from "./components/ScheduleEditorModal";

const baseClass = "manage-schedule-page";

// FAKE DATA ALERT

const fakeData = {
  scheduled: [
    {
      id: 1,
      query_id: 4,
      interval: 172800,
      last_executed: "2021-06-23T20:26:51Z",
    },
    {
      id: 2,
      query_id: 7,
      interval: 14400,
      last_executed: "2021-06-24T20:26:51Z",
    },
    {
      id: 3,
      query_id: 8,
      interval: 86400,
      last_executed: "2021-06-23T20:26:51Z",
    },
    {
      id: 4,
      query_id: 20,
      interval: 604800,
      last_executed: "2021-06-21T20:26:51Z",
    },
  ],
};

// 1. SET TO EMPTY IF YOU WANT TO SEE THE EMPTY STATE RENDER
const fakeDataLength0 = [1];

// 2. SET TO TRUE IF YOU WANT TO SEE THE ERROR STATE RENDER;
const fakeDataError = false;

// END FAKE DATA ALERT

const renderTable = (): JSX.Element => {
  // Schedule has not been set up for this instance yet.
  if (fakeDataLength0.length === 0) {
    return <NoSchedule />;
  }

  // Schedule has an error retrieving data.
  else if (fakeDataError) {
    return <ScheduleError />;
  }

  const scheduledCount = fakeData.scheduled.length;
  const scheduledQueriesTotalDisplay =
    scheduledCount === 1 ? "1 query" : `${scheduledCount} queries`;

  return (
    <div>
      <p className={`${baseClass}__scheduled-query-count`}>
        {scheduledQueriesTotalDisplay}
      </p>
    </div>
  );
};

const onScheduleEditorSubmit = () => {
  //   const { toggleScheduleEditorModal } = this;
  //   const { dispatch } = this.props;
  //   const { selectedHostIds } = this.state;
  //   const teamId = team.id === "no-team" ? null : team.id;
  //   dispatch(hostActions.transferToTeam(teamId, selectedHostIds))
  //     .then(() => {
  //       const successMessage =
  //         teamId === null
  //           ? `Hosts successfully removed from teams.`
  //           : `Hosts successfully transferred to  ${team.name}.`;
  //       dispatch(renderFlash("success", successMessage));
  //       dispatch(getHosts());
  //     })
  //     .catch(() => {
  //       dispatch(
  //         renderFlash("error", "Could not transfer hosts. Please try again.")
  //       );
  //     });
  //   toggleTransferHostModal();
  //   this.setState({ selectedHostIds: [] });
  // };
  // clearHostUpdates() {
  //   if (this.timeout) {
  //     global.window.clearTimeout(this.timeout);
  //     this.timeout = null;
  //   }
};

//  DON'T NEED!
// const renderScheduleEditorModal = (): JSX.Element | null => {

//   // toggleScheduleEditorModal,
//   // onScheduleEditorSubmit
//   // const { showScheduleEditorModal } = this.state;

//   if (!showScheduleEditorModal) return null;

//   return (
//     <ScheduleEditorModal
//       onSubmit={onScheduleEditorSubmit}
//       onCancel={toggleScheduleEditorModal}
//     />
//   );
// };
interface IRootState {
  // auth: {
  //   user: IUser;
  // };
  entities: {
    queries: IQuery[];
  };
}

const ManageSchedulePage = (): JSX.Element => {
  // Links to packs page
  const dispatch = useDispatch();
  const { MANAGE_PACKS } = paths;
  const handleAdvanced = (event: any) => dispatch(push(MANAGE_PACKS));

  // State to show modal
  const [showScheduleEditorModal, setShowScheduleEditorModal] = useState(false);

  // Toggle state to show modal
  const toggleScheduleEditorModal = useCallback(() => {
    setShowScheduleEditorModal(!showScheduleEditorModal);
  }, [showScheduleEditorModal, setShowScheduleEditorModal]);

  // Will I need this? 5/28
  // const user = useSelector((state: IRootState) => state.auth.user);

  const queries = useSelector((state: IRootState) => state.entities.queries);

  return (
    <div className={baseClass}>
      <div className={`${baseClass}__wrapper body-wrap`}>
        <div className={`${baseClass}__header-wrap`}>
          <div className={`${baseClass}__header`}>
            <div className={`${baseClass}__text`}>
              <h1 className={`${baseClass}__title`}>
                <span>Schedule</span>
              </h1>
              <div className={`${baseClass}__description`}>
                <p>
                  Schedule recurring queries for your hosts. Fleet’s query
                  schedule lets you add queries which are executed at regular
                  intervals.
                </p>
              </div>
            </div>
          </div>
          {/* Hides CTA Buttons if no schedule or schedule error */}
          {fakeDataLength0.length !== 0 && !fakeDataError && (
            <div className={`${baseClass}__action-button-container`}>
              <Button
                variant="inverse"
                onClick={handleAdvanced}
                className={`${baseClass}__advanced-button`}
              >
                Advanced
              </Button>
              {/* TODO: SCHEDULE A QUERY MODAL */}
              <Button
                variant="brand"
                className={`${baseClass}__schedule-button`}
                onClick={toggleScheduleEditorModal}
              >
                Schedule a query
              </Button>
            </div>
          )}
        </div>
        <div>{renderTable()}</div>
        {showScheduleEditorModal ? (
          <ScheduleEditorModal
            onCancel={toggleScheduleEditorModal}
            onSubmit={onScheduleEditorSubmit}
            queries={queries}
            // Modify onSubmit
          />
        ) : null}
      </div>
    </div>
  );
};

export default ManageSchedulePage;

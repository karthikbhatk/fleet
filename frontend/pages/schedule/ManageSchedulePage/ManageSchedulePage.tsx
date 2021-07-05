import React, { useState, useCallback } from "react"; //, { useEffect }
import { useDispatch, useSelector } from "react-redux";

import { push } from "react-router-redux";
// @ts-ignore
// import { IUser } from "interfaces/user";
import { IQuery } from "interfaces/query";
// @ts-ignore
import scheduledQueryActions from "redux/nodes/entities/scheduled_queries/actions";

// Will I need this? 5/28
// import permissionUtils from "utilities/permissions";
import paths from "router/paths";
import Button from "components/buttons/Button";
import ScheduleError from "./components/ScheduleError";
import ScheduleListWrapper from "./components/ScheduleListWrapper";
import ScheduleEditorModal from "./components/ScheduleEditorModal";
// @ts-ignore
import { renderFlash } from "redux/nodes/notifications/actions";

const baseClass = "manage-schedule-page";

// FAKE DATA ALERT

const fakeData = {
  scheduled: [
    {
      id: 1,
      query_id: 4,
      query_name: "Get crashes",
      interval: 172800,
      last_executed: "2021-06-23T20:26:51Z",
    },
    {
      id: 2,
      query_id: 7,
      query_name: "Detect machines with Gatekeeper disabled",
      interval: 14400,
      last_executed: "2021-06-24T20:26:51Z",
    },
    {
      id: 3,
      query_id: 8,
      query_name: "Detect fake data",
      interval: 86400,
      last_executed: "2021-06-23T20:26:51Z",
    },
    {
      id: 4,
      query_id: 20,
      query_name: "Detect a shit ton of work",
      interval: 604800,
      last_executed: "2021-06-21T20:26:51Z",
    },
  ],
};

// 1. SET TO TRUE IF YOU WANT TO SEE THE ERROR STATE RENDER;
const fakeDataError = false;

// END FAKE DATA ALERT

const renderTable = (): JSX.Element => {
  // Schedule has an error retrieving data.
  if (fakeDataError) {
    return <ScheduleError />;
  }

  // Schedule table
  return <ScheduleListWrapper fakeData={fakeData} />;
};

interface IRootState {
  entities: {
    queries: {
      data: IQuery[];
    };
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

  const allQueries = useSelector(
    (state: IRootState) => state.entities.queries.data
  );
  // Turn object of objects into array of objects
  const allQueriesList = Object.values(allQueries);

  // SIMILAR TO TEAMMANAGEMENTPAGE onCreateSubmit Line 85
  // SIMILAR TO MEMBERSPAGE onAddMemberSubmit Line 141
  // TODO: FUNCTIONALITY OF ONSUBMIT FORM 6/30, 7/2 WORK ON THIS
  const onScheduleSubmit = useCallback(
    (formData: IQuery) => {
      dispatch(
        // TODO: This is how to send formData to a new pack, is it the same for schedule? 7/2
        scheduledQueryActions({ ...formData, pack_id: 2 })
      )
        .then(() => {
          dispatch(
            renderFlash(
              "success",
              `Successfully added ${formData.name} to the schedule.`
            )
          );
          // Updates page
          dispatch(scheduledQueryActions.loadAll({}));
        })
        .catch(() => {
          dispatch(
            renderFlash("error", "Could not schedule query. Please try again.")
          );
        });
      toggleScheduleEditorModal();
    },
    [dispatch, toggleScheduleEditorModal]
  );

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
          {/* Hide CTA Buttons if no schedule or schedule error */}
          {fakeData.scheduled.length !== 0 && !fakeDataError && (
            <div className={`${baseClass}__action-button-container`}>
              <Button
                variant="inverse"
                onClick={handleAdvanced}
                className={`${baseClass}__advanced-button`}
              >
                Advanced
              </Button>
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
            onScheduleSubmit={onScheduleSubmit}
            allQueries={allQueriesList}
            defaultLoggingType={"snapshot"}
          />
        ) : null}
      </div>
    </div>
  );
};

export default ManageSchedulePage;

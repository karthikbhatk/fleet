import React from "react";

import { IActivity, IActivityDetails } from "interfaces/activity";
import { IActivitiesResponse } from "services/entities/activities";

// @ts-ignore
import FleetIcon from "components/icons/FleetIcon";
import DataError from "components/DataError";
import Button from "components/buttons/Button";

import EmptyFeed from "../EmptyFeed/EmptyFeed";
import UpcomingActivity from "../UpcomingActivity/UpcomingActivity";

const baseClass = "upcoming-activity-feed";

interface IUpcomingActivityFeedProps {
  activities?: IActivitiesResponse;
  isError?: boolean;
  onDetailsClick: (details: IActivityDetails) => void;
  onNextPage: () => void;
  onPreviousPage: () => void;
}

const testActivity = {
  created_at: "2021-07-27T13:25:21Z",
  id: 4,
  actor_full_name: "Rachael",
  actor_id: 1,
  actor_gravatar: "",
  actor_email: "rachael@example.com",
  type: "ran_script",
  details: {
    host_id: 1,
    host_display_name: "Steve's MacBook Pro",
    script_name: "",
    script_execution_id: "y3cffa75-b5b5-41ef-9230-15073c8a88cf",
  },
};

const UpcomingActivityFeed = ({
  activities,
  isError = false,
  onDetailsClick,
  onNextPage,
  onPreviousPage,
}: IUpcomingActivityFeedProps) => {
  if (isError) {
    return <DataError />;
  }

  if (!activities) {
    return null;
  }

  const { activities: activitiesList, meta } = activities;

  if (activitiesList.length === 0) {
    return (
      <EmptyFeed
        title="No pending activity "
        message="When you run a script on an offline host, it will appear here."
      />
    );
  }

  return (
    <div className={baseClass}>
      {activitiesList.map((activity: IActivity) => (
        <UpcomingActivity activity={activity} onDetailsClick={onDetailsClick} />
      ))}
      <div className={`${baseClass}__pagination`}>
        <Button
          disabled={!meta.has_previous_results}
          onClick={onPreviousPage}
          variant="unstyled"
          className={`${baseClass}__load-activities-button`}
        >
          <>
            <FleetIcon name="chevronleft" /> Previous
          </>
        </Button>
        <Button
          disabled={!meta.has_next_results}
          onClick={onNextPage}
          variant="unstyled"
          className={`${baseClass}__load-activities-button`}
        >
          <>
            Next <FleetIcon name="chevronright" />
          </>
        </Button>
      </div>
    </div>
  );
};

export default UpcomingActivityFeed;

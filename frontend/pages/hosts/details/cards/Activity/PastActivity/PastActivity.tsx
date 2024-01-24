import React from "react";
import ReactTooltip from "react-tooltip";
import { formatDistanceToNowStrict } from "date-fns";

import Avatar from "components/Avatar";
import Icon from "components/Icon";
import Button from "components/buttons/Button";

import { COLORS } from "styles/var/colors";
import { DEFAULT_GRAVATAR_LINK } from "utilities/constants";
import {
  addGravatarUrlToResource,
  internationalTimeFormat,
} from "utilities/helpers";
import { IActivity, IActivityDetails } from "interfaces/activity";
import { ShowActivityDetailsHandler } from "../Activity";

const baseClass = "past-activity";

interface IPastActivityProps {
  activity: IActivity; // TODO: type
  onDetailsClick: ShowActivityDetailsHandler;
}

const PastActivity = ({ activity, onDetailsClick }: IPastActivityProps) => {
  const { actor_email } = activity;
  const { gravatar_url } = actor_email
    ? addGravatarUrlToResource({ email: actor_email })
    : { gravatar_url: DEFAULT_GRAVATAR_LINK };
  const activityCreatedAt = new Date(activity.created_at);
  const scriptNameDisplay = activity.details?.script_name ? (
    <>
      the <b>{activity.details.script_name}</b> script
    </>
  ) : (
    "a script"
  );

  return (
    <div className={baseClass}>
      <Avatar
        className={`${baseClass}__avatar-image`}
        user={{ gravatar_url }}
        size="small"
        hasWhiteBackground
      />
      <div className={`${baseClass}__details-wrapper`}>
        <div className={"activity-details"}>
          <span className={`${baseClass}__details-topline`}>
            <b>{activity.actor_full_name} </b>
            <>
              {" "}
              told Fleet to run {scriptNameDisplay} on this host.{" "}
              <Button
                className={`${baseClass}__show-query-link`}
                variant="text-link"
                onClick={() => onDetailsClick?.(activity)}
              >
                Show details{" "}
                <Icon className={`${baseClass}__show-query-icon`} name="eye" />
              </Button>
            </>
          </span>
          <br />
          <span
            className={`${baseClass}__details-bottomline`}
            data-tip
            data-for={`activity-${activity.id}`}
          >
            {formatDistanceToNowStrict(activityCreatedAt, {
              addSuffix: true,
            })}
          </span>
          <ReactTooltip
            className="date-tooltip"
            place="top"
            type="dark"
            effect="solid"
            id={`activity-${activity.id}`}
            backgroundColor={COLORS["tooltip-bg"]}
          >
            {internationalTimeFormat(activityCreatedAt)}
          </ReactTooltip>
        </div>
      </div>
      <div className={`${baseClass}__dash`} />
    </div>
  );
};

export default PastActivity;

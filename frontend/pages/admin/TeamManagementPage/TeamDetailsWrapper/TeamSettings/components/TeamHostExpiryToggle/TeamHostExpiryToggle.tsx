import Checkbox from "components/forms/fields/Checkbox";
import Icon from "components/Icon";
import React from "react";
import { Link } from "react-router";

const baseClass = "team-host-expiry-toggle";

interface ITeamHostExpiryToggle {
  globalHostExpiryEnabled: boolean;
  globalHostExpiryWindow: number;
  teamExpiryEnabled: boolean;
  setTeamExpiryEnabled: (value: boolean) => void;
  addingCustomWindow: boolean;
  setAddingCustomWindow: (value: boolean) => void;
}

const TeamHostExpiryToggle = ({
  globalHostExpiryEnabled,
  globalHostExpiryWindow,
  teamExpiryEnabled,
  setTeamExpiryEnabled,
  addingCustomWindow,
  setAddingCustomWindow,
}: ITeamHostExpiryToggle) => {
  const renderHelpText = () =>
    globalHostExpiryEnabled ? (
      <div className="help-text">
        Host expiry is globally enabled in organization settings. By default,
        hosts expire after {globalHostExpiryWindow} days.{" "}
        {!addingCustomWindow && (
          <Link
            to={""}
            onClick={(e: React.MouseEvent) => {
              e.preventDefault();
              setAddingCustomWindow(true);
            }}
            className={`${baseClass}__add-custom-window`}
          >
            <>
              Add custom expiry window
              <Icon name="chevron-right" color="core-fleet-blue" size="small" />
            </>
          </Link>
        )}
      </div>
    ) : (
      <></>
    );
  return (
    <div className={`${baseClass}`}>
      <Checkbox
        name="enableHostExpiry"
        onChange={setTeamExpiryEnabled}
        value={teamExpiryEnabled || globalHostExpiryEnabled}
        wrapperClassName={
          globalHostExpiryEnabled
            ? `${baseClass}__disabled-team-host-expiry-toggle`
            : ""
        }
        helpText={renderHelpText()}
        tooltipContent={
          <>
            When enabled, allows automatic cleanup of
            <br />
            hosts that have not communicated with Fleet in
            <br />
            the number of days specified in the{" "}
            <strong>
              Host expiry
              <br />
              window
            </strong>{" "}
            setting.{" "}
            <em>
              (Default: <strong>Off</strong>)
            </em>
          </>
        }
      >
        Enable host expiry
      </Checkbox>
    </div>
  );
};

export default TeamHostExpiryToggle;

import React from "react";
import { useQuery } from "react-query";

import {
  SoftwareInstallStatus,
  getInstallStatusPredicate,
} from "interfaces/software";
import mdmApi from "services/entities/mdm";

import Modal from "components/Modal";
import Button from "components/buttons/Button";
import Icon from "components/Icon";
import Textarea from "components/Textarea";
import DataError from "components/DataError/DataError";
import Spinner from "components/Spinner/Spinner";
import { IMdmCommandResult } from "interfaces/mdm";
import { IActivityDetails } from "interfaces/activity";

import { IconNames } from "components/icons";
import {
  getInstallDetailsStatusPredicate,
  INSTALL_DETAILS_STATUS_ICONS,
} from "../constants";

const baseClass = "app-install-details";

export type IAppInstallDetails = Pick<
  IActivityDetails,
  | "host_id"
  | "command_uuid"
  | "host_display_name"
  | "software_title"
  | "app_store_id"
  | "status"
>;

export const AppInstallDetails = ({
  status,
  command_uuid = "",
  host_display_name = "",
  software_title = "",
}: IAppInstallDetails) => {
  const { data: result, isLoading, isError } = useQuery<
    IMdmCommandResult,
    Error
  >(
    ["mdm_command_results", command_uuid],
    async () => {
      return mdmApi.getCommandResults(command_uuid).then((response) => {
        const results = response.results?.[0];
        if (!results) {
          return Promise.reject(new Error("No data returned"));
        }
        return {
          ...results,
          payload: atob(results.payload),
          result: atob(results.result),
        };
      });
    },
    {
      refetchOnWindowFocus: false,
      staleTime: 3000,
    }
  );

  if (isLoading) {
    return <Spinner />;
  } else if (isError) {
    return <DataError description="Close this modal and try again." />;
  } else if (!result) {
    // FIXME: Find a better solution for this.
    return <DataError description="No data returned." />;
  }

  // Note: We need to reconcile status values from two different sources. From props, we
  // get the status from the activity item details (which can be "failed", "pending", or
  // "installed"). From the command results API response, we also receive the raw status
  // from the MDM protocol, e.g., "NotNow" or "Acknowledged". We need to display some special
  // messaging for the "NotNow" status, which otherwise would be treated as "pending".
  const isStatusNotNow = result.status === "NotNow";
  let iconName: IconNames;
  let predicate: string;
  let subordinate: string;
  if (isStatusNotNow) {
    iconName = INSTALL_DETAILS_STATUS_ICONS.pending;
    predicate = "tried to install";
    subordinate =
      " but couldn’t because the host was locked or was running on battery power while in Power Nap. Fleet will try again";
  } else {
    iconName = INSTALL_DETAILS_STATUS_ICONS[status as SoftwareInstallStatus];
    predicate = getInstallDetailsStatusPredicate(status);
    subordinate = status === "pending" ? " when it comes online" : "";
  }

  const showCommandResponse = isStatusNotNow || status !== "pending";

  return (
    <>
      <div className={`${baseClass}__software-install-details`}>
        <div className={`${baseClass}__status-message`}>
          {!!iconName && <Icon name={iconName} />}
          <span>
            Fleet {predicate} <b>{software_title}</b> on{" "}
            <b>{host_display_name}</b>
            {subordinate}.
          </span>
        </div>
        <div className={`${baseClass}__script-output`}>
          Request payload:
          <Textarea className={`${baseClass}__output-textarea`}>
            {result.payload}
          </Textarea>
        </div>
        {showCommandResponse && (
          <div className={`${baseClass}__script-output`}>
            The response from <b>{host_display_name}</b>:
            <Textarea className={`${baseClass}__output-textarea`}>
              {result.result}
            </Textarea>
          </div>
        )}
      </div>
    </>
  );
};

export const AppInstallDetailsModal = ({
  details,
  onCancel,
}: {
  details: IAppInstallDetails;
  onCancel: () => void;
}) => {
  return (
    <Modal
      title="Install details"
      onExit={onCancel}
      onEnter={onCancel}
      className={baseClass}
    >
      <>
        <div className={`${baseClass}__modal-content`}>
          <AppInstallDetails {...details} />
        </div>
        <div className="modal-cta-wrap">
          <Button onClick={onCancel} variant="brand">
            Done
          </Button>
        </div>
      </>
    </Modal>
  );
};
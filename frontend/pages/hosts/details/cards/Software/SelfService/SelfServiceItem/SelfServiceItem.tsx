import React, { useCallback, useContext, useEffect, useRef } from "react";
import ReactTooltip from "react-tooltip";

import {
  IDeviceSoftware,
  IHostSoftware,
  ISoftwareLastInstall,
  SoftwareInstallStatus,
} from "interfaces/software";
import deviceApi from "services/entities/device_user";
import { dateAgo } from "utilities/date_format";
import { NotificationContext } from "context/notification";

import Card from "components/Card";
import Button from "components/buttons/Button";
import Icon from "components/Icon";
import SoftwareIcon from "pages/SoftwarePage/components/icons/SoftwareIcon";

import { IStatusDisplayConfig } from "../../InstallStatusCell/InstallStatusCell";

const baseClass = "self-service-item";

const STATUS_CONFIG: Record<SoftwareInstallStatus, IStatusDisplayConfig> = {
  verified: {
    iconName: "success",
    displayText: "Verified",
    tooltip: ({ lastInstalledAt }) => (
      <>
        Software installed successfully ({dateAgo(lastInstalledAt as string)}).
        Currently, if the software is uninstalled, the &quot;Installed&quot;
        status won&apos;t be updated.
      </>
    ),
  },
  verifying: {
    iconName: "success-outline",
    displayText: "Verifying",
    tooltip: ({ lastInstalledAt }) => (
      <>
        Software installed successfully ({dateAgo(lastInstalledAt as string)}).
        Currently, if the software is uninstalled, the &quot;Installed&quot;
        status won&apos;t be updated.
      </>
    ),
  },
  pending: {
    iconName: "pending-outline",
    displayText: "Install in progress...",
    tooltip: () => "Software installation in progress...",
  },
  blocked: {
    iconName: "pending-outline", // TODO
    displayText: "TODO SelfServieItem.tsx", // TODO
    tooltip: () => "", // TODO
  },
  failed: {
    iconName: "error",
    displayText: "Failed",
    tooltip: ({ lastInstalledAt = "" }) => (
      <>
        Software failed to install
        {lastInstalledAt ? ` (${dateAgo(lastInstalledAt)})` : ""}. Select{" "}
        <b>Retry</b> to install again, or contact your IT department.
      </>
    ),
  },
};

interface IInstallerInfoProps {
  software: IDeviceSoftware;
}

const InstallerInfo = ({ software }: IInstallerInfoProps) => {
  const { name, source, software_package: installerPackage } = software;
  return (
    <div className={`${baseClass}__item-topline`}>
      <div className={`${baseClass}__item-icon`}>
        <SoftwareIcon name={name} source={source} size="large" />
      </div>
      <div className={`${baseClass}__item-name-version`}>
        <div className={`${baseClass}__item-name`}>
          {name || installerPackage?.name}
        </div>
        <div className={`${baseClass}__item-version`}>
          {installerPackage?.version || ""}
        </div>
      </div>
    </div>
  );
};

// TODO: update if/when we support self-service app store apps
type IInstallerStatusProps = Pick<IHostSoftware, "id" | "status"> & {
  last_install: ISoftwareLastInstall | null;
};

const InstallerStatus = ({
  id,
  status,
  last_install,
}: IInstallerStatusProps) => {
  const displayConfig = STATUS_CONFIG[status as keyof typeof STATUS_CONFIG];
  if (!displayConfig) {
    // API should ensure this never happens, but just in case
    return null;
  }

  return (
    <div className={`${baseClass}__status-content`}>
      <div
        className={`${baseClass}__status-with-tooltip`}
        data-tip
        data-for={`install-tooltip__${id}`}
      >
        <Icon name={displayConfig.iconName} />
        <span data-testid={`${baseClass}__status--test`}>
          {displayConfig.displayText}
        </span>
      </div>
      <ReactTooltip
        className={`${baseClass}__status-tooltip`}
        effect="solid"
        backgroundColor="#3e4771"
        id={`install-tooltip__${id}`}
        data-html
      >
        <span className={`${baseClass}__status-tooltip-text`}>
          {displayConfig.tooltip({
            lastInstalledAt: last_install?.installed_at,
          })}
        </span>
      </ReactTooltip>
    </div>
  );
};

interface IInstallerStatusActionProps {
  deviceToken: string;
  software: IHostSoftware;
  onInstall: () => void;
}

const getInstallButtonText = (status: SoftwareInstallStatus | null) => {
  switch (status) {
    case null:
      return "Install";
    case "failed":
      return "Retry";
    case "verified":
      return "Reinstall";
    default:
      // we don't show a button for pending installs
      return "";
  }
};

const InstallerStatusAction = ({
  deviceToken,
  software: { id, status, software_package },
  onInstall,
}: IInstallerStatusActionProps) => {
  const { renderFlash } = useContext(NotificationContext);

  // TODO: update this if/when we support self-service app store apps
  const last_install = software_package?.last_install || null;

  // localStatus is used to track the status of the any user-initiated install action
  const [localStatus, setLocalStatus] = React.useState<
    SoftwareInstallStatus | undefined
  >(undefined);

  // displayStatus allows us to display the localStatus (if any) or the status from the list
  // software reponse
  const displayStatus = localStatus || status;
  const installButtonText = getInstallButtonText(displayStatus);

  // if the localStatus is "failed", we don't want our tooltip to include the old installed_at date so we
  // set this to null, which tells the tooltip to omit the parenthetical date
  const lastInstall = localStatus === "failed" ? null : last_install;

  const isMountedRef = useRef(false);
  useEffect(() => {
    isMountedRef.current = true;
    return () => {
      isMountedRef.current = false;
    };
  }, []);

  const onClick = useCallback(async () => {
    setLocalStatus("pending");
    try {
      await deviceApi.installSelfServiceSoftware(deviceToken, id);
      if (isMountedRef.current) {
        onInstall();
      }
    } catch (error) {
      renderFlash("error", "Couldn't install. Please try again.");
      if (isMountedRef.current) {
        setLocalStatus("failed");
      }
    }
  }, [deviceToken, id, onInstall, renderFlash]);

  return (
    <div className={`${baseClass}__item-status-action`}>
      <div className={`${baseClass}__item-status`}>
        <InstallerStatus
          id={id}
          status={displayStatus}
          last_install={lastInstall}
        />
      </div>
      <div className={`${baseClass}__item-action`}>
        {!!installButtonText && (
          <Button
            variant="text-icon"
            type="button"
            className={`${baseClass}__item-action-button${
              localStatus === "pending" ? "--installing" : ""
            }`}
            onClick={onClick}
          >
            <span data-testid={`${baseClass}__item-action-button--test`}>
              {installButtonText}
            </span>
          </Button>
        )}
      </div>
    </div>
  );
};

interface ISelfServiceItemProps {
  deviceToken: string;
  software: IDeviceSoftware;
  onInstall: () => void;
}

const SelfServiceItem = ({
  deviceToken,
  software,
  onInstall,
}: ISelfServiceItemProps) => {
  return (
    <Card
      borderRadiusSize="large"
      paddingSize="medium"
      className={`${baseClass}__item`}
    >
      <div className={`${baseClass}__item-content`}>
        <InstallerInfo software={software} />
        <InstallerStatusAction
          deviceToken={deviceToken}
          software={software}
          onInstall={onInstall}
        />
      </div>
    </Card>
  );
};

export default SelfServiceItem;

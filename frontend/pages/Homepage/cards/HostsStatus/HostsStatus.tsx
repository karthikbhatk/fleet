import React, { useState } from "react";

import WindowsIcon from "../../../../../assets/images/icon-windows-48x48@2x.png";
import LinuxIcon from "../../../../../assets/images/icon-linux-48x48@2x.png";
import MacIcon from "../../../../../assets/images/icon-mac-48x48@2x.png";

const baseClass = "hosts-status";

interface IHostSummaryProps {
  onlineCount: string | undefined;
  offlineCount: string | undefined;
  newCount: string | undefined;
}

const HostsStatus = ({
  onlineCount,
  offlineCount,
  newCount,
}: IHostSummaryProps): JSX.Element => {
  return (
    <div className={baseClass}>
      <div className={`${baseClass}__tile online-tile`}>
        <div>
          <div
            className={`${baseClass}__tile-count ${baseClass}__tile-count--online`}
          >
            {onlineCount}
          </div>
          <div className={`${baseClass}__tile-description`}>Online hosts</div>
        </div>
      </div>
      <div className={`${baseClass}__tile offline-tile`}>
        <div>
          <div
            className={`${baseClass}__tile-count ${baseClass}__tile-count--offline`}
          >
            {offlineCount}
          </div>
          <div className={`${baseClass}__tile-description`}>Offline hosts</div>
        </div>
      </div>
      <div className={`${baseClass}__tile new-tile`}>
        <div>
          <div
            className={`${baseClass}__tile-count ${baseClass}__tile-count--new`}
          >
            {newCount}
          </div>
          <div className={`${baseClass}__tile-description`}>New hosts</div>
        </div>
      </div>
    </div>
  );
};

export default HostsStatus;

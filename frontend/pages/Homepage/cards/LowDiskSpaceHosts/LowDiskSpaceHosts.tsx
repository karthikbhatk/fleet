import React from "react";
import PATHS from "router/paths";

import SummaryTile from "../HostsSummary/SummaryTile";
import LowDiskSpaceIcon from "../../../../../assets/images/icon-low-disk-space-32x19@2x.png";

const baseClass = "low-disk-space";

interface IHostSummaryProps {
  lowDiskSpaceCount: number;
  isLoadingHosts: boolean;
  showHostsUI: boolean;
}

const LowDiskSpaceHosts = ({
  lowDiskSpaceCount,
  isLoadingHosts,
  showHostsUI,
}: IHostSummaryProps): JSX.Element => {
  return (
    <SummaryTile
      icon={LowDiskSpaceIcon}
      count={lowDiskSpaceCount}
      isLoading={isLoadingHosts}
      showUI={showHostsUI}
      title="Low disk space hosts"
      tooltip="Hosts that have 32 GB or less disk space available."
      path={`${PATHS.MANAGE_HOSTS}?low_disk_space=true`}
    />
  );
};

export default LowDiskSpaceHosts;

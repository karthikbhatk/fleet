import React from "react";

import ReactTooltip from "react-tooltip";

interface IDiskSpaceGraphProps {
  baseClass: string;
  gigsDiskSpaceAvailable: number | string;
  percentDiskSpaceAvailable: number;
  id: string;
  platform: string;
}

const DiskSpaceGraph = ({
  baseClass,
  gigsDiskSpaceAvailable,
  percentDiskSpaceAvailable,
  id,
  platform,
}: IDiskSpaceGraphProps): JSX.Element => {
  const diskSpaceIndicator = (): string => {
    // return space-dependent graph colors for mac and windows hosts, green for linux
    if (platform === "darwin" || platform === "windows") {
      if (gigsDiskSpaceAvailable < 16) {
        return "red";
      } else if (gigsDiskSpaceAvailable < 32) {
        return "yellow";
      }
    }
    return "green";
  };

  const diskSpaceTooltipText = ((): string | undefined => {
    if (platform === "darwin" || platform === "windows") {
      if (gigsDiskSpaceAvailable < 16) {
        return "Not enough disk space available to install most small operating systems updates.";
      } else if (gigsDiskSpaceAvailable < 32) {
        return "Not enough disk space available to install most large operating systems updates.";
      }
      return "Enough disk space available to install most operating systems updates.";
    }
    return undefined;
  })();

  if (gigsDiskSpaceAvailable === 0 || gigsDiskSpaceAvailable === "---") {
    return <span className={`${baseClass}__data`}>No data available</span>;
  }

  return (
    <span className={`${baseClass}__data`}>
      <div
        className={`${baseClass}__disk-space-wrapper tooltip`}
        data-tip
        data-for={id}
      >
        <div className={`${baseClass}__disk-space`}>
          <div
            className={`${baseClass}__disk-space--${diskSpaceIndicator()}`}
            style={{
              width: `${100 - percentDiskSpaceAvailable}%`,
            }}
          />
        </div>
      </div>
      {diskSpaceTooltipText && (
        <ReactTooltip
          className={"disk-space-tooltip"}
          place="bottom"
          type="dark"
          effect="solid"
          id={id}
          backgroundColor="#3e4771"
        >
          <span className={`${baseClass}__tooltip-text`}>
            {diskSpaceTooltipText}
          </span>
        </ReactTooltip>
      )}
      {gigsDiskSpaceAvailable} GB{baseClass === "info-flex" && " available"}
    </span>
  );
};

export default DiskSpaceGraph;

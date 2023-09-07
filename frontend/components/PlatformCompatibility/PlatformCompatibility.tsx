import React from "react";

import { OsqueryPlatform } from "interfaces/platform";
import { PLATFORM_DISPLAY_NAMES } from "utilities/constants";

import NewTooltipWrapper from "components/NewTooltipWrapper";
import Icon from "components/Icon";

interface IPlatformCompatibilityProps {
  compatiblePlatforms: OsqueryPlatform[] | null;
  error: Error | null;
}

const baseClass = "platform-compatibility";

const DISPLAY_ORDER = [
  "macOS",
  "Windows",
  "Linux",
  "ChromeOS",
] as OsqueryPlatform[];

const ERROR_NO_COMPATIBLE_TABLES = Error("no tables in query");

const formatPlatformsForDisplay = (
  compatiblePlatforms: OsqueryPlatform[]
): OsqueryPlatform[] => {
  return compatiblePlatforms.map((str) => PLATFORM_DISPLAY_NAMES[str] || str);
};

const displayIncompatibilityText = (err: Error) => {
  switch (err) {
    case ERROR_NO_COMPATIBLE_TABLES:
      return (
        <span>
          No platforms (check your query for invalid tables or tables that are
          supported on different platforms)
        </span>
      );
    default:
      return (
        <span>No platforms (check your query for a possible syntax error)</span>
      );
  }
};

const PlatformCompatibility = ({
  compatiblePlatforms,
  error,
}: IPlatformCompatibilityProps): JSX.Element | null => {
  if (!compatiblePlatforms) {
    return null;
  }
  if (error || !compatiblePlatforms?.length) {
    return (
      <span className={baseClass}>
        <b>
          <NewTooltipWrapper tipContent="Estimated compatiblity based on <br /> the tables used in the query.">
            Compatible with:
          </NewTooltipWrapper>
        </b>

        {displayIncompatibilityText(error || ERROR_NO_COMPATIBLE_TABLES)}
      </span>
    );
  }

  const displayPlatforms = formatPlatformsForDisplay(compatiblePlatforms);
  return (
    <span className={baseClass}>
      <b>
        <NewTooltipWrapper tipContent="Estimated compatiblity based on <br /> the tables used in the query.">
          Compatible with:
        </NewTooltipWrapper>
      </b>
      {DISPLAY_ORDER.map((platform) => {
        const isCompatible = displayPlatforms.includes(platform);
        return (
          <span
            key={`platform-compatibility__${platform}`}
            className="platform"
          >
            <Icon
              name={isCompatible ? "check" : "ex"}
              className={
                isCompatible ? "compatible-platform" : "incompatible-platform"
              }
              color={isCompatible ? "status-success" : "status-error"}
            />
            {platform}
          </span>
        );
      })}
    </span>
  );
};
export default PlatformCompatibility;

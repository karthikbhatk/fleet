import React from "react";
import ReactTooltip from "react-tooltip";
import { IconNames } from "components/icons";
import Icon from "components/Icon";
import Button from "components/buttons/Button";

const baseClass = "settings-indicator";

interface IMacSettingsIndicator {
  indicatorText: string;
  iconName: IconNames;
  onClick?: () => void;
  tooltip?: {
    tooltipText: string | null;
    position?: "top" | "bottom";
  };
}

const MacSettingsIndicator = ({
  indicatorText,
  iconName,
  onClick,
  tooltip,
}: IMacSettingsIndicator): JSX.Element => {
  const getIndicatorTextWrapped = () => {
    if (onClick && tooltip) {
      return (
        <>
          <span
            className="tooltip tooltip__tooltip-icon"
            data-tip
            data-for={`${indicatorText}-tooltip`}
            data-tip-disable={false}
          >
            <Button
              onClick={onClick}
              variant="text-link"
              className={`${baseClass}__button`}
            >
              {indicatorText}
            </Button>
          </span>
          <ReactTooltip
            place={tooltip.position ?? "bottom"}
            effect="solid"
            backgroundColor="#3e4771"
            id={`${indicatorText}-tooltip`}
            data-html
          >
            <span className="tooltip__tooltip-text">{tooltip.tooltipText}</span>
          </ReactTooltip>
        </>
      );
    }

    // onclick without tooltip
    if (onClick) {
      return (
        <Button
          onClick={onClick}
          variant="text-link"
          className={`${baseClass}__button`}
        >
          {indicatorText}
        </Button>
      );
    }

    // tooltip without onclick
    if (tooltip) {
      return (
        <>
          <span
            className="tooltip tooltip__tooltip-icon"
            data-tip
            data-for="settings-indicator"
            data-tip-disable={false}
          >
            {indicatorText}
          </span>
          <ReactTooltip
            place={tooltip.position ?? "bottom"}
            effect="solid"
            backgroundColor="#3e4771"
            id="settings-indicator"
            data-html
          >
            <span className="tooltip__tooltip-text">{tooltip.tooltipText}</span>
          </ReactTooltip>
        </>
      );
    }

    // no tooltip, no onclick
    return indicatorText;
  };

  return (
    <span className={`${baseClass} info-flex__data`}>
      <Icon name={iconName} />
      {getIndicatorTextWrapped()}
    </span>
  );
};

export default MacSettingsIndicator;

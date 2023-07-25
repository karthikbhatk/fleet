import React from "react";
import { Link } from "react-router";
import classnames from "classnames";

import Icon from "components/Icon/Icon";

const baseClass = "stacked-white-boxes";

interface IStackedWhiteBoxesProps {
  children?: JSX.Element;
  headerText?: string;
  className?: string;
  leadText?: string;
  previousLocation?: string;
}

const StackedWhiteBoxes = ({
  children,
  headerText,
  className,
  leadText,
  previousLocation,
}: IStackedWhiteBoxesProps): JSX.Element => {
  const boxClass = classnames(baseClass, className);

  const renderBackButton = () => {
    if (!previousLocation) return false;

    return (
      <div className={`${baseClass}__back`}>
        <Link to={previousLocation} className={`${baseClass}__back-link`}>
          <Icon name="ex" color="core-fleet-black" />
        </Link>
      </div>
    );
  };

  return (
    <div className={boxClass}>
      <div className={`${baseClass}__box`}>
        {renderBackButton()}
        {headerText && (
          <p className={`${baseClass}__header-text`}>{headerText}</p>
        )}
        {leadText && <p className={`${baseClass}__box-text`}>{leadText}</p>}
        {children}
      </div>
    </div>
  );
};

export default StackedWhiteBoxes;

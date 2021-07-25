import React, { useCallback } from "react";
import PropTypes from "prop-types";

// ignore TS error for now until these are rewritten in ts.
// @ts-ignore
import Button from "../../../buttons/Button";

export interface IActionButtonProps {
  buttonText: string;
  onActionButtonClick: (targetIds: number[]) => void | undefined;
  targetIds?: number[]; // TODO figure out undefined case
  variant?: string;
  hideButton?: boolean | ((targetIds: number[]) => boolean);
}

function useActionCallback(
  callbackFn: (targetIds: number[]) => void | undefined
) {
  return useCallback(
    (targetIds) => {
      callbackFn(targetIds);
    },
    [callbackFn]
  );
}

const ActionButton = (props: IActionButtonProps): JSX.Element | null => {
  const {
    buttonText,
    onActionButtonClick,
    targetIds = [],
    variant,
    hideButton,
  } = props;
  const onButtonClick = useActionCallback(onActionButtonClick);

  // hideButton is intended to provide a flexible way to specify show/hide conditions via a boolean or a function that evaluates to a boolean
  // currently it is typed to accept an array of targetIds but this typing could easily be expanded to include other use cases
  const testCondition = (
    hideButtonProp: boolean | ((targetIds: number[]) => boolean) | undefined
  ) => {
    if (typeof hideButtonProp === "function") {
      return hideButtonProp;
    }
    return () => Boolean(hideButtonProp);
  };
  const isHidden = testCondition(hideButton)(targetIds);

  return !isHidden ? (
    <Button onClick={() => onButtonClick(targetIds)} variant={variant}>
      {buttonText}
    </Button>
  ) : null;
};

ActionButton.propTypes = {
  buttonText: PropTypes.string,
  onActionButtonClick: PropTypes.func,
  targetIds: PropTypes.arrayOf(PropTypes.number),
  variant: PropTypes.string,
  hideButton: PropTypes.oneOfType([PropTypes.bool, PropTypes.func]),
};

export default ActionButton;

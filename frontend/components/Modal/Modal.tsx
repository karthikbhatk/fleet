import React, { useEffect } from "react";
import classnames from "classnames";
import Button from "components/buttons/Button/Button";
import Icon from "components/Icon/Icon";

const baseClass = "modal";

type ModalWidth = "medium" | "large" | "xlarge" | "auto";

export interface IModalProps {
  title: string | JSX.Element;
  children: JSX.Element;
  onExit: () => void;
  onEnter?: () => void;
  /**     default 650px, large 800px, xlarge 850px, auto auto-width */
  width?: ModalWidth;
  /**    isHidden can be set true to hide the modal when opening another modal */
  isHidden?: boolean;
  /**    isLoading can be set true to enable targeting elements by loading state */
  isLoading?: boolean;
  className?: string;
}

const Modal = ({
  title,
  children,
  onExit,
  onEnter,
  width = "medium",
  isHidden = false,
  isLoading = false,
  className,
}: IModalProps): JSX.Element => {
  useEffect(() => {
    const closeWithEscapeKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        onExit();
      }
    };

    document.addEventListener("keydown", closeWithEscapeKey);

    return () => {
      document.removeEventListener("keydown", closeWithEscapeKey);
    };
  }, []);

  useEffect(() => {
    if (onEnter) {
      const closeOrSaveWithEnterKey = (event: KeyboardEvent) => {
        if (event.code === "Enter" || event.code === "NumpadEnter") {
          event.preventDefault();
          onEnter();
        }
      };

      document.addEventListener("keydown", closeOrSaveWithEnterKey);
      return () => {
        document.removeEventListener("keydown", closeOrSaveWithEnterKey);
      };
    }
  }, [onEnter]);

  const backgroundClasses = classnames(`${baseClass}__background`, {
    [`${baseClass}__hidden`]: isHidden,
  });

  const modalContainerClasses = classnames(
    className,
    `${baseClass}__modal_container`,
    `${baseClass}__modal_container__${width}`,
    {
      [`${className}__loading`]: isLoading,
    }
  );

  return (
    <div className={backgroundClasses}>
      <div className={modalContainerClasses}>
        <div className={`${baseClass}__header`}>
          <span>{title}</span>
          <div className={`${baseClass}__ex`}>
            <Button className="button button--unstyled" onClick={onExit}>
              <Icon name="close" color="core-fleet-black" size="medium" />
            </Button>
          </div>
        </div>
        <div className={`${baseClass}__content`}>{children}</div>
      </div>
    </div>
  );
};

export default Modal;

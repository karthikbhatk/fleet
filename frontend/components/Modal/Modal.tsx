import React, { useContext, useEffect } from "react";
import classnames from "classnames";
import { NotificationContext } from "context/notification";

const baseClass = "modal";

type ModalWidth = "m" | "l" | "xl" | "auto";

export interface IModalProps {
  title: string | JSX.Element;
  children: JSX.Element;
  onExit: () => void;
  onEnter?: () => void;
  width?: ModalWidth;
  className?: string;
}

const Modal = ({
  title,
  children,
  onExit,
  onEnter,
  width = "m",
  className,
}: IModalProps): JSX.Element => {
  const { hideFlash } = useContext(NotificationContext);

  useEffect(() => {
    const closeWithEscapeKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        onExit();
      }
    };

    hideFlash();
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

  const modalContainerClassName = classnames(
    `${baseClass}__modal_container`,
    className,
    { [`${baseClass}__modal_container__l`]: width === "l" },
    { [`${baseClass}__modal_container__xl`]: width === "xl" },
    { [`${baseClass}__modal_container__auto`]: width === "auto" }
  );

  return (
    <div className={`${baseClass}__background`}>
      <div className={modalContainerClassName}>
        <div className={`${baseClass}__header`}>
          <span>{title}</span>
          <div className={`${baseClass}__ex`}>
            <button className="button button--unstyled" onClick={onExit} />
          </div>
        </div>
        <div className={`${baseClass}__content`}>{children}</div>
      </div>
    </div>
  );
};

export default Modal;

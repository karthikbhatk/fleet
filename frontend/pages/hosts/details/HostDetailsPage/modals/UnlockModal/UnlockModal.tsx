import React, { useContext } from "react";
import { AxiosError } from "axios";

import { NotificationContext } from "context/notification";
import hostAPI from "services/entities/hosts";

import Modal from "components/Modal";
import Button from "components/buttons/Button";

const baseClass = "unlock-modal";

interface IUnlockModalProps {
  id: number;
  platform: string;
  hostName: string;
  pin?: number;
  onClose: () => void;
}

const UnlockModal = ({
  id,
  platform,
  hostName,
  pin,
  onClose,
}: IUnlockModalProps) => {
  const { renderFlash } = useContext(NotificationContext);
  const [isUnlocking, setIsUnlocking] = React.useState(false);

  const onUnlock = async () => {
    setIsUnlocking(true);
    try {
      await hostAPI.unlockHost(id);
      renderFlash("success", "Host Unlocked successfully!");
    } catch (error) {
      const err = error as AxiosError;
      renderFlash("error", err.message);
    }
    onClose();
    setIsUnlocking(false);
  };

  const renderModalContent = () => {
    if (platform === "darwin" && pin) {
      return (
        <>
          {/* TODO: replace with DataSet component */}
          <p>
            When the host is returned, use the 6-digit PIN to unlock the host.
          </p>
          <div className={`${baseClass}__pin`}>
            <b>PIN</b>
            <span>{pin}</span>
          </div>
        </>
      );
    }

    return (
      <>
        <p>
          Are you sure you&apos;re ready to unlock <b>{hostName}</b>?
        </p>
      </>
    );
  };

  const renderModalButtons = () => {
    if (platform === "darwin") {
      return (
        <>
          <Button type="button" onClick={onClose} variant="brand">
            Done
          </Button>
        </>
      );
    }

    return (
      <>
        <Button
          type="button"
          onClick={onUnlock}
          variant="brand"
          className="delete-loading"
          isLoading={isUnlocking}
        >
          Unlock
        </Button>
        <Button onClick={onClose} variant="inverse">
          Cancel
        </Button>
      </>
    );
  };

  return (
    <Modal className={baseClass} title="Unlock host" onExit={onClose}>
      <>
        <div className={`${baseClass}__modal-content`}>
          {renderModalContent()}
        </div>

        <div className="modal-cta-wrap">{renderModalButtons()}</div>
      </>
    </Modal>
  );
};

export default UnlockModal;

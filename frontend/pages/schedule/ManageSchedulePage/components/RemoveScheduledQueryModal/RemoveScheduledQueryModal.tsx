import React, { useState, useCallback } from "react";

import Modal from "components/modals/Modal";
import Button from "components/buttons/Button";
import { IScheduledQuery } from "interfaces/scheduled_query";

const baseClass = "remove-scheduled-query-modal";

interface IRemoveScheduledQueryModalProps {
  queries: IScheduledQuery[];
  onCancel: () => void;
  onSubmit: () => void;
}

const RemoveScheduledQueryModal = (
  props: IRemoveScheduledQueryModalProps
): JSX.Element => {
  const { onCancel, onSubmit, queries } = props;

  return (
    <Modal
      title={"Remove scheduled query"}
      onExit={onCancel}
      className={baseClass}
    >
      <div className={baseClass}>
        Are you sure you want to remove the selected queries from the schedule?
        <div className={`${baseClass}__btn-wrap`}>
          <Button
            className={`${baseClass}__btn`}
            onClick={onCancel}
            variant="inverse"
          >
            Cancel
          </Button>
          <Button
            className={`${baseClass}__btn`}
            type="button"
            variant="alert"
            onClick={onSubmit}
          >
            Remove
          </Button>
        </div>
      </div>
    </Modal>
  );
};

export default RemoveScheduledQueryModal;

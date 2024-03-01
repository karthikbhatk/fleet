import React, { useMemo } from "react";

import { IMdmSolution } from "interfaces/mdm";

import Modal from "components/Modal";
import TableContainer from "components/TableContainer";

import {
  generateSolutionsDataSet,
  generateSolutionsTableHeaders,
} from "./MdmSolutionModalTableConfig";

const baseClass = "mdm-solution-modal";

const SOLUTIONS_DEFAULT_SORT_HEADER = "hosts_count";
const DEFAULT_SORT_DIRECTION = "desc";

interface IMdmSolutionModalProps {
  mdmSolutions: IMdmSolution[];
  selectedPlatformLabelId?: number;
  selectedTeamId?: number;
  onCancel: () => void;
}

const MdmSolutionsModal = ({
  mdmSolutions,
  selectedPlatformLabelId,
  selectedTeamId,
  onCancel,
}: IMdmSolutionModalProps) => {
  console.log(mdmSolutions);

  const solutionsTableHeaders = useMemo(
    () => generateSolutionsTableHeaders(selectedTeamId),
    [selectedTeamId]
  );
  const solutionsDataSet = generateSolutionsDataSet(
    mdmSolutions,
    selectedPlatformLabelId
  );

  return (
    <Modal
      className={baseClass}
      title={mdmSolutions[0].name ?? "Mdm solution"}
      width="large"
      onExit={onCancel}
    >
      <TableContainer
        isLoading={false}
        emptyComponent={() => null} // if this modal is shown, this table should never be empty
        columnConfigs={solutionsTableHeaders}
        data={solutionsDataSet}
        defaultSortHeader={SOLUTIONS_DEFAULT_SORT_HEADER}
        defaultSortDirection={DEFAULT_SORT_DIRECTION}
        resultsTitle="MDM"
        showMarkAllPages={false}
        isAllPagesSelected={false}
        disableCount
      />
    </Modal>
  );
};

export default MdmSolutionsModal;

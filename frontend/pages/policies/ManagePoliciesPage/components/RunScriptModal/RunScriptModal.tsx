import React, { useCallback, useState } from "react";
import { useQuery } from "react-query";
import { omit } from "lodash";

import { DEFAULT_USE_QUERY_OPTIONS } from "utilities/constants";

import scriptsAPI, {
  IListScriptsQueryKey,
  IScriptsResponse,
} from "services/entities/scripts";

import { IPolicyStats } from "interfaces/policy";
import { IScript } from "interfaces/script";

// @ts-ignore
import Dropdown from "components/forms/fields/Dropdown";
import Checkbox from "components/forms/fields/Checkbox";
import TooltipTruncatedText from "components/TooltipTruncatedText";
import DataError from "components/DataError";
import Spinner from "components/Spinner";
import CustomLink from "components/CustomLink";
import Button from "components/buttons/Button";
import Modal from "components/Modal";

// name avoids class name conflicts with hosts/details/HostDetailsPage/modals/RunScriptModal
const baseClass = "policy-run-script-modal";

interface IScriptDropdownField {
  name: string;
  value: number;
}

interface IFormPolicy {
  name: string;
  id: number;
  runScriptEnabled: boolean;
  scriptIdToRun?: number;
}

export type IPolicyRunScriptFormData = IFormPolicy[];

interface IRunScriptModal {
  onExit: () => void;
  onSubmit: (formData: IPolicyRunScriptFormData) => void;
  isUpdating: boolean;
  policies: IPolicyStats[];
  teamId: number;
}

const RunScriptModal = ({
  onExit,
  onSubmit,
  isUpdating,
  policies,
  teamId,
}: IRunScriptModal) => {
  const [formData, setFormData] = useState<IPolicyRunScriptFormData>(
    policies.map((policy) => ({
      name: policy.name,
      id: policy.id,
      runScriptEnabled: !!policy.run_script,
      scriptIdToRun: policy.run_script?.id,
    }))
  );

  const anyEnabledWithoutSelection = formData.some(
    (policy) => policy.runScriptEnabled && !policy.scriptIdToRun
  );

  const {
    data: availableScripts,
    isLoading: isLoadingAvailableScripts,
    isError: isAvailableScriptsError,
  } = useQuery<IScriptsResponse, Error, IScript[], [IListScriptsQueryKey]>(
    [
      {
        scope: "scripts",
        team_id: teamId,
      },
    ],
    ({ queryKey: [queryKey] }) =>
      scriptsAPI.getScripts(omit(queryKey, "scope")),
    {
      select: (data) => data.scripts,
      ...DEFAULT_USE_QUERY_OPTIONS,
    }
  );

  const onUpdate = useCallback(() => {
    onSubmit(formData);
  }, [formData, onSubmit]);

  const onChangeEnableRunScript = useCallback(
    (newVal: { policyName: string; value: boolean }) => {
      const { policyName, value } = newVal;
      setFormData(
        formData.map((policy) => {
          if (policy.name === policyName) {
            return {
              ...policy,
              runScriptEnabled: value,
              scriptIdToRun: value ? policy.scriptIdToRun : undefined,
            };
          }
          return policy;
        })
      );
    },
    [formData]
  );

  const onSelectPolicyScript = useCallback(
    ({ name, value }: IScriptDropdownField) => {
      const [policyName, scriptId] = [name, value];
      setFormData(
        formData.map((policy) => {
          if (policy.name === policyName) {
            return { ...policy, scriptIdToRun: scriptId };
          }
          return policy;
        })
      );
    },
    [formData]
  );

  const availableScriptOptions = availableScripts?.map((script) => ({
    label: script.name,
    value: script.id,
  }));

  const renderPolicyRunScriptOption = (policy: IFormPolicy) => {
    const {
      name: policyName,
      id: policyId,
      runScriptEnabled: enabled,
      scriptIdToRun,
    } = policy;

    return (
      <li
        className={`${baseClass}__policy-row policy-row`}
        id={`policy-row--${policyId}`}
        key={policyId}
      >
        <Checkbox
          value={enabled}
          name={policyName}
          onChange={() => {
            onChangeEnableRunScript({
              policyName,
              value: !enabled,
            });
          }}
        >
          <TooltipTruncatedText value={policyName} />
        </Checkbox>
        {enabled && (
          <Dropdown
            options={availableScriptOptions}
            value={scriptIdToRun}
            onChange={onSelectPolicyScript}
            placeholder="Select script"
            className={`${baseClass}__script-dropdown`}
            name={policyName}
            parseTarget
          />
        )}
      </li>
    );
  };

  const renderContent = () => {
    if (isAvailableScriptsError) {
      return <DataError />;
    }
    if (isLoadingAvailableScripts) {
      return <Spinner />;
    }
    if (!availableScripts?.length) {
      return (
        <div className={`${baseClass}__no-scripts`}>
          <b>No scripts available for install</b>
          <span>
            Go to <b>Controls &gt; Scripts</b> to add scripts to this team.
          </span>
        </div>
      );
    }

    return (
      <div className={`${baseClass} form`}>
        <div className="form-field">
          <div className="form-field__label">Policies:</div>
          <ul className="automated-policies-section">
            {formData.map((policyData) =>
              renderPolicyRunScriptOption(policyData)
            )}
          </ul>
          <span className="form-field__help-text">
            Selected script will be run when hosts fail the chosen policy.{" "}
            {/* TODO - confirm link destination */}
            <CustomLink
              url="https://fleetdm.com/learn-more-about/policy-automation-run-script"
              text="Learn more"
              newTab
            />
          </span>
        </div>
        <div className="modal-cta-wrap">
          <Button
            type="submit"
            variant="brand"
            onClick={onUpdate}
            className="save-loading"
            isLoading={isUpdating}
            disabled={anyEnabledWithoutSelection}
          >
            Save
          </Button>
          <Button onClick={onExit} variant="inverse">
            Cancel
          </Button>
        </div>
      </div>
    );
  };
  return (
    <Modal
      title="Run script"
      className={baseClass}
      onExit={onExit}
      onEnter={onUpdate}
      width="large"
      isContentDisabled={isUpdating}
    >
      {renderContent()}
    </Modal>
  );
};

export default RunScriptModal;

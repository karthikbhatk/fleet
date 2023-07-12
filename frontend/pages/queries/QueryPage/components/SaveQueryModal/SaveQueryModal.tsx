import React, { useState, useEffect, useCallback } from "react";
import { pull, size } from "lodash";

import { IQueryFormData } from "interfaces/query";
import useDeepEffect from "hooks/useDeepEffect";

import Checkbox from "components/forms/fields/Checkbox";
// @ts-ignore
import InputField from "components/forms/fields/InputField";
// @ts-ignore
import Dropdown from "components/forms/fields/Dropdown";
import Button from "components/buttons/Button";
import Modal from "components/Modal";
import {
  FREQUENCY_DROPDOWN_OPTIONS,
  LOGGING_TYPE_OPTIONS,
  MIN_OSQUERY_VERSION_OPTIONS,
  SCHEDULE_PLATFORM_DROPDOWN_OPTIONS,
} from "utilities/constants";
import RevealButton from "components/buttons/RevealButton";
import { IPlatformString } from "interfaces/platform";
import {
  ISchedulableQuery,
  QueryLoggingOption,
} from "interfaces/schedulable_query";

export interface ISaveQueryModalProps {
  baseClass: string;
  queryValue: string;
  isLoading: boolean;
  saveQuery: (formData: IQueryFormData) => void;
  toggleSaveQueryModal: () => void;
  backendValidators: { [key: string]: string };
  existingQuery?: ISchedulableQuery;
}

const validateQueryName = (name: string) => {
  const errors: { [key: string]: string } = {};

  if (!name) {
    errors.name = "Query name must be present";
  }

  const valid = !size(errors);
  return { valid, errors };
};

const SaveQueryModal = ({
  baseClass,
  queryValue,
  isLoading,
  saveQuery,
  toggleSaveQueryModal,
  backendValidators,
  existingQuery,
}: ISaveQueryModalProps): JSX.Element => {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [selectedFrequency, setSelectedFrequency] = useState(
    existingQuery?.interval ?? 3600
  );
  const [
    selectedPlatformOptions,
    setSelectedPlatformOptions,
  ] = useState<IPlatformString>(existingQuery?.platform ?? "");
  const [
    selectedMinOsqueryVersionOptions,
    setSelectedMinOsqueryVersionOptions,
  ] = useState(existingQuery?.min_osquery_version ?? "");
  const [
    selectedLoggingType,
    setSelectedLoggingType,
  ] = useState<QueryLoggingOption | null>(existingQuery?.logging ?? "snapshot");
  const [observerCanRun, setObserverCanRun] = useState(false);
  const [errors, setErrors] = useState<{ [key: string]: string }>(
    backendValidators
  );
  const [showAdvancedOptions, setShowAdvancedOptions] = useState(false);

  const toggleAdvancedOptions = () => {
    setShowAdvancedOptions(!showAdvancedOptions);
  };

  useDeepEffect(() => {
    if (name) {
      setErrors({});
    }
  }, [name]);

  useEffect(() => {
    setErrors(backendValidators);
  }, [backendValidators]);

  const onClickSaveQuery = (evt: React.MouseEvent<HTMLFormElement>) => {
    evt.preventDefault();

    const { valid, errors: newErrors } = validateQueryName(name);
    setErrors({
      ...errors,
      ...newErrors,
    });

    if (valid) {
      saveQuery({
        description,
        name,
        query: queryValue,
        observer_can_run: observerCanRun,
      });
    }
  };

  const onChangeSelectPlatformOptions = useCallback(
    (values: string) => {
      const valArray = values.split(",");

      // Remove All if another OS is chosen
      // else if Remove OS if All is chosen
      if (valArray.indexOf("") === 0 && valArray.length > 1) {
        // TODO - inmprove type safety of all 3 options
        setSelectedPlatformOptions(
          pull(valArray, "").join(",") as IPlatformString
        );
      } else if (valArray.length > 1 && valArray.indexOf("") > -1) {
        setSelectedPlatformOptions("");
      } else {
        setSelectedPlatformOptions(values as IPlatformString);
      }
    },
    [setSelectedPlatformOptions]
  );

  return (
    <Modal title={"Save query"} onExit={toggleSaveQueryModal}>
      <>
        <form
          onSubmit={onClickSaveQuery}
          className={`${baseClass}__save-modal-form`}
          autoComplete="off"
        >
          <InputField
            name="name"
            onChange={(value: string) => setName(value)}
            value={name}
            error={errors.name}
            inputClassName={`${baseClass}__query-save-modal-name`}
            label="Name"
            placeholder="What is your query called?"
            autofocus
          />
          <InputField
            name="description"
            onChange={(value: string) => setDescription(value)}
            value={description}
            inputClassName={`${baseClass}__query-save-modal-description`}
            label="Description"
            type="textarea"
            placeholder="What information does your query reveal? (optional)"
          />
          <Dropdown
            searchable={false}
            options={FREQUENCY_DROPDOWN_OPTIONS}
            onChange={(value: number) => {
              setSelectedFrequency(value);
            }}
            placeholder={"Every hour"}
            value={selectedFrequency}
            label="Frequency"
            wrapperClassName={`${baseClass}__form-field ${baseClass}__form-field--frequency`}
          />
          <p>
            If automations are on, this is how often your query collects data.
          </p>
          <Checkbox
            name="observerCanRun"
            onChange={setObserverCanRun}
            value={observerCanRun}
            wrapperClassName={`${baseClass}__query-save-modal-observer-can-run-wrapper`}
          >
            Observers can run
          </Checkbox>
          <p>
            Users with the Observer role will be able to run this query as a
            live query.
          </p>
          <RevealButton
            isShowing={showAdvancedOptions}
            className={baseClass}
            hideText={"Hide advanced options"}
            showText={"Show advanced options"}
            caretPosition={"after"}
            onClick={toggleAdvancedOptions}
          />
          {showAdvancedOptions && (
            <>
              <Dropdown
                options={SCHEDULE_PLATFORM_DROPDOWN_OPTIONS}
                placeholder="Select"
                label="Platforms"
                onChange={onChangeSelectPlatformOptions}
                value={selectedPlatformOptions}
                multi
                wrapperClassName={`${baseClass}__form-field ${baseClass}__form-field--platform`}
              />
              <p>
                If automations are turned on, your query collects data on
                compatible platforms.
                <br />
                If you want more control, override platforms.
              </p>
              <Dropdown
                options={MIN_OSQUERY_VERSION_OPTIONS}
                onChange={setSelectedMinOsqueryVersionOptions}
                placeholder="Select"
                value={selectedMinOsqueryVersionOptions}
                label="Minimum osquery version"
                wrapperClassName={`${baseClass}__form-field ${baseClass}__form-field--osquer-vers`}
              />
              <Dropdown
                options={LOGGING_TYPE_OPTIONS}
                onChange={setSelectedLoggingType}
                placeholder="Select"
                value={selectedLoggingType}
                label="Logging"
                wrapperClassName={`${baseClass}__form-field ${baseClass}__form-field--logging`}
              />
            </>
          )}
          <div className="modal-cta-wrap">
            <Button
              type="submit"
              variant="brand"
              className="save-query-loading"
              isLoading={isLoading}
            >
              Save
            </Button>
            <Button onClick={toggleSaveQueryModal} variant="inverse">
              Cancel
            </Button>
          </div>
        </form>
      </>
    </Modal>
  );
};

export default SaveQueryModal;

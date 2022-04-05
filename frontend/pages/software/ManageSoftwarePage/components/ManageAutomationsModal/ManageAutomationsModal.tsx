import React, { useState } from "react";
import { useQuery } from "react-query";

import { Link } from "react-router";
import PATHS from "router/paths";

import {
  IIntegrations,
  IJiraIntegration,
  IJiraIntegrationIndexed,
  IJiraIntegrationFormData,
  IJiraIntegrationFormErrors,
} from "interfaces/integration";
import { IConfigNested } from "interfaces/config";

import configAPI from "services/entities/config";
import MOCKS from "services/mock_service/mocks/responses";

// @ts-ignore
import Dropdown from "components/forms/fields/Dropdown";
import Modal from "components/Modal";
import Button from "components/buttons/Button";
import Slider from "components/forms/fields/Slider";
import Radio from "components/forms/fields/Radio";
// @ts-ignore
import InputField from "components/forms/fields/InputField";

import { IWebhookSoftwareVulnerabilities } from "interfaces/webhook";
import { useDeepEffect } from "utilities/hooks";
import { size } from "lodash";

import PreviewPayloadModal from "../PreviewPayloadModal";

interface IManageAutomationsModalProps {
  onCancel: () => void;
  onCreateWebhookSubmit: (formData: IWebhookSoftwareVulnerabilities) => void;
  togglePreviewPayloadModal: () => void;
  showPreviewPayloadModal: boolean;
  softwareVulnerabilityWebhookEnabled?: boolean;
  currentDestinationUrl?: string;
}

const validateWebhookURL = (url: string) => {
  const errors: { [key: string]: string } = {};

  if (url === "") {
    errors.url = "Please add a destination URL";
  }

  const valid = !size(errors);
  return { valid, errors };
};

const baseClass = "manage-automations-modal";

const ManageAutomationsModal = ({
  onCancel: onReturnToApp,
  onCreateWebhookSubmit,
  togglePreviewPayloadModal,
  showPreviewPayloadModal,
  softwareVulnerabilityWebhookEnabled,
  currentDestinationUrl,
}: IManageAutomationsModalProps): JSX.Element => {
  const [destination_url, setDestinationUrl] = useState<string>(
    currentDestinationUrl || ""
  );
  const [errors, setErrors] = useState<{ [key: string]: string }>({});
  const [
    softwareAutomationsEnabled,
    setSoftwareAutomationsEnabled,
  ] = useState<boolean>(softwareVulnerabilityWebhookEnabled || false);
  const [jiraEnabled, setJiraEnabled] = useState<boolean>(false);
  const [integrationsIndexed, setIntegrationsIndexed] = useState<
    IJiraIntegrationIndexed[]
  >();
  const [
    selectedIntegration,
    setSelectedIntegration,
  ] = useState<IJiraIntegration>();

  useDeepEffect(() => {
    if (destination_url) {
      setErrors({});
    }
  }, [destination_url]);

  const {
    data: integrations,
    isLoading: isLoadingIntegrations,
    error: loadingIntegrationsError,
    refetch: refetchIntegrations,
  } = useQuery<IConfigNested, Error, IJiraIntegration[]>(
    ["integrations"],
    () => configAPI.loadAll(),
    {
      select: (data: IConfigNested) => {
        return data.integrations.jira;
      },
      onSuccess: (data) => {
        if (data) {
          const addIndex = data.map((integration, index) => {
            return { ...integration, integrationIndex: index };
          });
          setIntegrationsIndexed(addIndex);

          console.log("addIndex", addIndex);
        }
      },
    }
  );

  const onURLChange = (value: string) => {
    setDestinationUrl(value);
  };

  const handleSaveAutomation = (evt: React.MouseEvent<HTMLFormElement>) => {
    evt.preventDefault();

    const { valid, errors: newErrors } = validateWebhookURL(destination_url);
    setErrors({
      ...errors,
      ...newErrors,
    });

    // URL validation only needed if software automation is checked
    if (valid || !softwareAutomationsEnabled) {
      onCreateWebhookSubmit({
        destination_url,
        enable_vulnerabilities_webhook: softwareAutomationsEnabled,
      });

      onReturnToApp();
    }
  };

  const createIntegrationDropdownOptions = () => {
    const integrationOptions = integrationsIndexed?.map((i) => {
      return {
        value: String(i.integrationIndex),
        label: `${i.url} - ${i.project_key}`,
      };
    });
    return integrationOptions;
  };

  const onChangeSelectIntegration = (selectIntegrationIndex: string) => {
    const integrationWithIndex:
      | IJiraIntegrationIndexed
      | undefined = integrationsIndexed?.find(
      (integ: IJiraIntegrationIndexed) =>
        integ.integrationIndex === parseInt(selectIntegrationIndex, 10)
    );
    setSelectedIntegration(integrationWithIndex);
  };

  const onRadioChange = (jira: boolean): ((evt: string) => void) => {
    console.log("onRadioChange formField", jira);
    return (evt: string) => {
      console.log("onRadioChange evt", evt);
      setJiraEnabled(jira);
    };
  };

  const renderTicket = () => {
    return (
      <div className={`${baseClass}__ticket`}>
        <div className={`${baseClass}__software-automation-description`}>
          <p>
            A ticket will be created in your <b>Integration</b> if a detected
            vulnerability (CVE) was published in the last 2 days.
          </p>
        </div>
        {integrationsIndexed && integrationsIndexed.length > 0 ? (
          <Dropdown
            searchable
            options={createIntegrationDropdownOptions()}
            onChange={onChangeSelectIntegration}
            placeholder={"Select Jira integration"}
            value={selectedIntegration?.integrationIndex}
            label={"Integration"}
            wrapperClassName={`${baseClass}__form-field ${baseClass}__form-field--frequency`}
            hint={
              "For each new vulnerability detected, Fleet will create a ticket with a list of the affected hosts."
            }
          />
        ) : (
          <div className={`${baseClass}__no-integrations`}>
            <div>
              <b>You have no integrations.</b>
            </div>
            <div className={`${baseClass}__no-integration--cta`}>
              <Link
                to={PATHS.ADMIN_INTEGRATIONS}
                className={`${baseClass}__add-integration-link`}
              >
                <span>Add integration</span>
              </Link>
            </div>
          </div>
        )}
      </div>
    );
  };

  const renderWebhook = () => {
    return (
      <div className={`${baseClass}__webhook`}>
        <div className={`${baseClass}__software-automation-description`}>
          <p>
            A request will be sent to your configured <b>Destination URL</b> if
            a detected vulnerability (CVE) was published in the last 2 days.
          </p>
        </div>
        <InputField
          inputWrapperClass={`${baseClass}__url-input`}
          name="webhook-url"
          label={"Destination URL"}
          type={"text"}
          value={destination_url}
          onChange={onURLChange}
          error={errors.url}
          hint={
            "For each new vulnerability detected, Fleet will send a JSON payload to this URL with a list of the affected hosts."
          }
          placeholder={"https://server.com/example"}
          tooltip="Provide a URL to deliver a webhook request to."
        />
        <Button
          type="button"
          variant="text-link"
          onClick={togglePreviewPayloadModal}
        >
          Preview payload
        </Button>
      </div>
    );
  };

  if (showPreviewPayloadModal) {
    return <PreviewPayloadModal onCancel={togglePreviewPayloadModal} />;
  }

  return (
    <Modal
      onExit={onReturnToApp}
      title={"Manage automations"}
      className={baseClass}
    >
      <div className={baseClass}>
        <div className={`${baseClass}__software-select-items`}>
          <Slider
            value={softwareAutomationsEnabled}
            onChange={() =>
              setSoftwareAutomationsEnabled(!softwareAutomationsEnabled)
            }
            inactiveText={"Vulnerability automations disabled"}
            activeText={"Vulnerability automations enabled"}
          />
        </div>
        <div className={`${baseClass}__overlay-container`}>
          <div className={`${baseClass}__software-automation-enabled`}>
            <div className={`${baseClass}__workflow`}>
              Workflow
              <Radio
                className={`${baseClass}__radio-input`}
                label={"Ticket"}
                id={"ticket"}
                checked={jiraEnabled}
                value={"ticket"}
                name={"ticket"}
                onChange={onRadioChange(true)}
              />
              <Radio
                className={`${baseClass}__radio-input`}
                label={"Webhook"}
                id={"webhook"}
                checked={!jiraEnabled}
                value={"webhook"}
                name={"webhook"}
                onChange={onRadioChange(false)}
              />
            </div>
            {jiraEnabled ? renderTicket() : renderWebhook()}
          </div>
          {!softwareAutomationsEnabled && (
            <div className={`${baseClass}__overlay`} />
          )}
        </div>
        <div className={`${baseClass}__button-wrap`}>
          <Button
            className={`${baseClass}__btn`}
            onClick={onReturnToApp}
            variant="inverse"
          >
            Cancel
          </Button>
          <Button
            className={`${baseClass}__btn`}
            type="submit"
            variant="brand"
            onClick={handleSaveAutomation}
          >
            Save
          </Button>
        </div>
      </div>
    </Modal>
  );
};

export default ManageAutomationsModal;

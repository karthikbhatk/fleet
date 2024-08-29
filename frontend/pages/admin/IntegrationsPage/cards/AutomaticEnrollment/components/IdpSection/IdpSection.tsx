import React, { useCallback, useContext, useState } from "react";

import configAPI from "services/entities/config";

// @ts-ignore
import InputField from "components/forms/fields/InputField";
import CustomLink from "components/CustomLink/CustomLink";
import Button from "components/buttons/Button/Button";
import SectionHeader from "components/SectionHeader";
import { NotificationContext } from "context/notification";
import { AppContext } from "context/app";
import { expandErrorReasonRequired } from "interfaces/errors";
import { AxiosResponse } from "axios";
import TooltipWrapper from "components/TooltipWrapper";
import {
  IFormDataIdp,
  IFormErrorsIdp,
  newFormDataIdp,
  validateFormDataIdp,
} from "./helpers";

const baseClass = "idp-section";

const IdpSection = () => {
  const { config } = useContext(AppContext);
  const { renderFlash } = useContext(NotificationContext);
  const [formData, setFormData] = useState(
    newFormDataIdp(config?.mdm?.end_user_authentication)
  );
  const [formErrors, setFormErrors] = useState<IFormErrorsIdp | null>(null);

  const onInputChange = useCallback(
    ({ name, value }: { name: keyof IFormDataIdp; value: string }) => {
      const newData = { ...formData, [name]: value };
      setFormData(newData);
      setFormErrors(validateFormDataIdp(newData));
    },
    [formData]
  );

  const onSubmit = useCallback(
    async (e: React.FormEvent<SubmitEvent>) => {
      e.preventDefault();
      const newErrors = validateFormDataIdp(formData);
      if (newErrors) {
        setFormErrors(newErrors);
        return;
      }

      try {
        await configAPI.update({
          mdm: {
            end_user_authentication: {
              ...formData,
            },
          },
        });
        renderFlash("success", "Successfully updated end user authentication!");
      } catch (err) {
        const ae = (typeof err === "object" ? err : {}) as AxiosResponse;
        if (ae.status === 422) {
          renderFlash(
            "error",
            `Couldn’t update: ${expandErrorReasonRequired(err)}.`
          );
          return;
        }
        renderFlash("error", "Couldn’t update. Please try again.");
      }
    },
    [formData, renderFlash]
  );

  return (
    <div className={baseClass}>
      <SectionHeader title="End user authentication" />
      <form>
        <p>
          Connect Fleet to your identity provider to require end users to
          authenticate when they first setup their new macOS hosts.{" "}
          <CustomLink
            url="https://fleetdm.com/docs/using-fleet/mdm-macos-setup-experience##end-user-authentication-and-eula"
            text="Learn more"
            newTab
          />
        </p>
        <InputField
          label="Identity provider name"
          onChange={onInputChange}
          name="idp_name"
          value={formData.idp_name}
          parseTarget
          error={formErrors?.idp_name}
          tooltip="A required human friendly name for the identity provider that will provide single sign-on authentication."
        />
        <InputField
          label="Entity ID"
          onChange={onInputChange}
          name="entity_id"
          value={formData.entity_id}
          parseTarget
          error={formErrors?.entity_id}
          tooltip="The required entity ID is a URI that you use to identify Fleet when configuring the identity provider."
        />
        <InputField
          label="Metadata URL"
          helpText={
            <>
              If both <b>Metadata URL</b> and <b>Metadata</b> are specified,{" "}
              <b>Metadata URL</b> will be used.
            </>
          }
          onChange={onInputChange}
          name="metadata_url"
          value={formData.metadata_url}
          parseTarget
          error={formErrors?.metadata_url}
          tooltip="Metadata URL provided by the identity provider."
        />
        <InputField
          label="Metadata"
          type="textarea"
          onChange={onInputChange}
          name="metadata"
          value={formData.metadata}
          parseTarget
          error={formErrors?.metadata}
          tooltip="Metadata XML provided by the identity provider."
        />
        <TooltipWrapper
          tipContent="Complete all required fields to save end user authentication."
          disableTooltip={!formErrors}
        >
          <Button
            disabled={!!formErrors}
            onClick={onSubmit}
            className="button-wrap"
          >
            Save
          </Button>
        </TooltipWrapper>
      </form>
    </div>
  );
};

export default IdpSection;

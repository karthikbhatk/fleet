import React, { useCallback, useContext, useState } from "react";
import { Params } from "react-router/lib/Router";
import { useQuery } from "react-query";
import { useErrorHandler } from "react-error-boundary";

import { IConfig } from "interfaces/config";
import { IApiError } from "interfaces/errors";
import configAPI from "services/entities/config";
import { AppContext } from "context/app";
import { NotificationContext } from "context/notification";
import deepDifference from "utilities/deep_difference";
import SandboxGate from "components/Sandbox/SandboxGate";
import SandboxDemoMessage from "components/Sandbox/SandboxDemoMessage";
import Spinner from "components/Spinner";

import SideNav from "../components/SideNav";
import ORG_SETTINGS_NAV_ITEMS from "./OrgSettingsNavItems";

interface IAppSettingsPageProps {
  params: Params;
}

export const baseClass = "app-settings";

const AppSettingsPage = ({ params }: IAppSettingsPageProps) => {
  const { section } = params;
  const DEFAULT_SETTINGS_SECTION = ORG_SETTINGS_NAV_ITEMS[0];

  const [isUpdatingSettings, setIsUpdatingSettings] = useState(false);
  const { isFreeTier, isPremiumTier, setConfig } = useContext(AppContext);
  const { renderFlash } = useContext(NotificationContext);
  const handlePageError = useErrorHandler();

  const { data: appConfig, isLoading, refetch: refetchConfig } = useQuery<
    IConfig,
    Error,
    IConfig
  >(["config"], () => configAPI.loadAll(), {
    select: (data: IConfig) => data,
    onSuccess: (data) => {
      setConfig(data);
    },
  });

  const onFormSubmit = useCallback(
    (formData: Partial<IConfig>) => {
      if (!appConfig) {
        return false;
      }

      setIsUpdatingSettings(true);

      const diff = deepDifference(formData, appConfig);
      // send all formData.agent_options because diff overrides all agent options
      diff.agent_options = formData.agent_options;

      configAPI
        .update(diff)
        .then(() => {
          renderFlash("success", "Successfully updated settings.");
          refetchConfig();
        })
        .catch((response: { data: IApiError }) => {
          if (
            response?.data.errors[0].reason.includes("could not dial smtp host")
          ) {
            renderFlash(
              "error",
              "Could not connect to SMTP server. Please try again."
            );
          } else if (response?.data.errors) {
            const agentOptionsInvalid =
              response.data.errors[0].reason.includes(
                "unsupported key provided"
              ) ||
              response.data.errors[0].reason.includes("invalid value type");

            renderFlash(
              "error",
              <>
                Could not update settings. {response.data.errors[0].reason}
                {agentOptionsInvalid && (
                  <>
                    <br />
                    If you’re not using the latest osquery, use the fleetctl
                    apply --force command to override validation.
                  </>
                )}
              </>
            );
          }
        })
        .finally(() => {
          setIsUpdatingSettings(false);
        });
    },
    [appConfig, refetchConfig, renderFlash]
  );

  const currentFormSection =
    ORG_SETTINGS_NAV_ITEMS.find((item) => item.urlSection === section) ??
    DEFAULT_SETTINGS_SECTION;

  const CurrentCard = currentFormSection.Card;

  if (isFreeTier && section === "fleet-desktop") {
    handlePageError({ status: 403 });
    return null;
  }

  return (
    <div className={`${baseClass}`}>
      <p className={`${baseClass}__page-description`}>
        Set your organization information and configure SSO and SMTP
      </p>
      <SandboxGate
        fallbackComponent={() => (
          <SandboxDemoMessage
            message="Organization settings are only available in self-managed Fleet"
            utmSource="fleet-ui-organization-settings-page"
            className={`${baseClass}__sandbox-demo-message`}
          />
        )}
      >
        <SideNav
          className={`${baseClass}__side-nav`}
          navItems={ORG_SETTINGS_NAV_ITEMS}
          activeItem={currentFormSection.urlSection}
          CurrentCard={
            !isLoading && appConfig
              ? () => (
                  <CurrentCard
                    appConfig={appConfig}
                    handleSubmit={onFormSubmit}
                    isUpdatingSettings={isUpdatingSettings}
                    isPremiumTier={isPremiumTier}
                  />
                )
              : () => <Spinner />
          }
        />
      </SandboxGate>
    </div>
  );
};

export default AppSettingsPage;

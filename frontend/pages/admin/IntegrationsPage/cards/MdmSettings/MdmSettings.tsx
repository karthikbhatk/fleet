import React, { useContext } from "react";
import { useQuery } from "react-query";
import { AxiosError } from "axios";
import { InjectedRouter } from "react-router";

import { AppContext } from "context/app";

import mdmAppleAPI from "services/entities/mdm_apple";
import { IMdmApple } from "interfaces/mdm";

import MdmSettingsSection from "./components/MdmSettingsSection";
import AutomaticEnrollmentSection from "./components/AutomaticEnrollmentSection";
import VppSection from "./components/VppSection";
import IdpSection from "./components/IdpSection";
import EulaSection from "./components/EulaSection";
import EndUserMigrationSection from "./components/EndUserMigrationSection";

const baseClass = "mdm-settings";

interface IMdmSettingsProps {
  router: InjectedRouter;
}

const MdmSettings = ({ router }: IMdmSettingsProps) => {
  const { isPremiumTier, config } = useContext(AppContext);

  // Currently the status of this API call is what determines various UI states on
  // this page. Because of this we will not render any of this components UI until this API
  // call has completed.
  const {
    data: appleAPNInfo,
    isLoading: isLoadingMdmApple,
    error: errorMdmApple,
  } = useQuery<IMdmApple, AxiosError, IMdmApple>(
    ["appleAPNInfo"],
    () => mdmAppleAPI.getAppleAPNInfo(),
    {
      retry: (tries, error) =>
        error.status !== 404 && error.status !== 400 && tries <= 3,
      // TODO: There is a potential race condition here immediately after MDM is turned off. This
      // component gets remounted and stale config data is used to determine it this API call is
      // enabled, resulting in a 400 response. The race really should  be fixed higher up the chain where
      // we're fetching and setting the config, but for now we'll just assume that any 400 response
      // means that MDM is not enabled and we'll show the "Turn on MDM" button.
      staleTime: 5000,
      enabled: !!config?.mdm.enabled_and_configured,
    }
  );

  const hasAllData = !isLoadingMdmApple && !errorMdmApple;

  return (
    <div className={baseClass}>
      <MdmSettingsSection
        isLoading={isLoadingMdmApple}
        appleAPNInfo={appleAPNInfo}
        appleAPNError={errorMdmApple}
        router={router}
      />
      {hasAllData && (
        <>
          <AutomaticEnrollmentSection
            router={router}
            isPremiumTier={!!isPremiumTier}
          />
          <VppSection router={router} isPremiumTier={!!isPremiumTier} />
          {isPremiumTier && (
            <>
              <IdpSection />
              <EulaSection />
              <EndUserMigrationSection router={router} />
            </>
          )}
        </>
      )}
    </div>
  );
};

export default MdmSettings;

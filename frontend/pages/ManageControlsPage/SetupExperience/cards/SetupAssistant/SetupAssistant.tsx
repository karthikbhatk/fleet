import React, { useState } from "react";
import { useQuery } from "react-query";

import { IConfig } from "interfaces/config";
import { API_NO_TEAM_ID, ITeam, ITeamConfig } from "interfaces/team";
import configAPI from "services/entities/config";
import teamsAPI, { ILoadTeamResponse } from "services/entities/teams";
import mdmAPI, {
  IAppleSetupEnrollmentProfileResponse,
} from "services/entities/mdm";

import SectionHeader from "components/SectionHeader";
import Spinner from "components/Spinner";
import CustomLink from "components/CustomLink";

import SetupAssistantPreview from "./components/SetupAssistantPreview";
import SetupAssistantProfileUploader from "./components/SetupAssistantProfileUploader";
import SetuAssistantProfileCard from "./components/SetupAssistantProfileCard/SetupAssistantProfileCard";
import DeleteAutoEnrollmentProfile from "./components/DeleteAutoEnrollmentProfile";
import AdvancedOptionsForm from "./components/AdvancedOptionsForm";
import { DEFAULT_USE_QUERY_OPTIONS } from "utilities/constants";
import { Axios, AxiosError } from "axios";

const baseClass = "setup-assistant";

interface ISetupAssistantProps {
  currentTeamId: number;
}

const StartupAssistant = ({ currentTeamId }: ISetupAssistantProps) => {
  const [showDeleteProfileModal, setShowDeleteProfileModal] = useState(false);

  const { data: globalConfig, isLoading: isLoadingGlobalConfig } = useQuery<
    IConfig,
    Error
  >(["config", currentTeamId], () => configAPI.loadAll(), {
    enabled: currentTeamId === API_NO_TEAM_ID,
    refetchOnWindowFocus: false,
    retry: false,
  });

  const { data: teamConfig, isLoading: isLoadingTeamConfig } = useQuery<
    ILoadTeamResponse,
    Error,
    ITeamConfig
  >(["team", currentTeamId], () => teamsAPI.load(currentTeamId), {
    refetchOnWindowFocus: false,
    retry: false,
    enabled: currentTeamId !== API_NO_TEAM_ID,
    select: (res) => res.team,
  });

  const {
    data: enrollmentProfileData,
    isLoading: isLoadingEnrollmentProfile,
    isError: isErrorEnrollmentProfile,
  } = useQuery<IAppleSetupEnrollmentProfileResponse, AxiosError>(
    ["enrollment_profile", currentTeamId],
    () => mdmAPI.getSetupEnrollmentProfile(currentTeamId),
    DEFAULT_USE_QUERY_OPTIONS
  );

  const getReleaseDeviceSetting = () => {
    if (currentTeamId === API_NO_TEAM_ID) {
      return (
        globalConfig?.mdm.macos_setup.enable_release_device_manually || false
      );
    }
    return teamConfig?.mdm?.macos_setup.enable_release_device_manually || false;
  };

  const onUpload = () => {};

  const onDelete = () => {};

  const defaultReleaseDeviceSetting = getReleaseDeviceSetting();

  return (
    <div className={baseClass}>
      <SectionHeader title="Setup assistant" />
      {isLoadingEnrollmentProfile ? (
        <Spinner />
      ) : (
        <div className={`${baseClass}__content`}>
          <div className={`${baseClass}__upload-container`}>
            <p className={`${baseClass}__section-description`}>
              Add an automatic enrollment profile to customize the macOS Setup
              Assistant.
              <CustomLink
                url=" https://fleetdm.com/learn-more-about/setup-assistant"
                text="Learn how"
                newTab
              />
            </p>
            {!enrollmentProfileData ? (
              <SetupAssistantProfileUploader
                currentTeamId={currentTeamId}
                onUpload={() => 1}
              />
            ) : (
              <SetuAssistantProfileCard
                profileMetaData={1}
                currentTeamId={currentTeamId}
                onDelete={() => setShowDeleteProfileModal(true)}
              />
            )}
            <AdvancedOptionsForm
              key={String(defaultReleaseDeviceSetting)}
              currentTeamId={currentTeamId}
              defaultReleaseDevice={defaultReleaseDeviceSetting}
            />
          </div>
          <div className={`${baseClass}__preview-container`}>
            <SetupAssistantPreview />
          </div>
        </div>
      )}
      {showDeleteProfileModal && (
        <DeleteAutoEnrollmentProfile
          onDelete={onDelete}
          onCancel={() => setShowDeleteProfileModal(false)}
        />
      )}
    </div>
  );
};

export default StartupAssistant;

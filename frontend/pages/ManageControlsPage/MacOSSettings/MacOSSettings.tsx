import React, { useContext, useState } from "react";
import { Params } from "react-router/lib/Router";

import { AppContext } from "context/app";
import SideNav from "pages/admin/components/SideNav";
import { useQuery } from "react-query";
import { IMdmProfile, IMdmProfilesResponse } from "interfaces/mdm";
import { IConfig } from "interfaces/config";

import mdmAPI from "services/entities/mdm";
import configAPI from "services/entities/config";

import MAC_OS_SETTINGS_NAV_ITEMS from "./MacOSSettingsNavItems";
import AggregateMacSettingsIndicators from "./AggregateMacSettingsIndicators/AggregateMacSettingsIndicators";

const baseClass = "mac-os-settings";

interface IMacOSSettingsProps {
  params: Params;
  // location field looks like this to integrate with the react router Route component, which
  // renders this one
  location: {
    query: { team_id?: string };
  };
}

const MacOSSettings = ({ params, location }: IMacOSSettingsProps) => {
  const { section } = params;
  const { team_id } = location.query;
  // Avoids possible case where Number(undefined) returns NaN
  const teamId = team_id === undefined ? 0 : Number(team_id); // team_id===0 for 'No teams'

  const { setConfig } = useContext(AppContext);

  const [profiles, setProfiles] = useState<IMdmProfile[] | null>();

  const { data: mdmProfiles, refetch: refetchProfiles } = useQuery<
    IMdmProfilesResponse,
    unknown,
    IMdmProfile[] | null
  >(["profiles", teamId], () => mdmAPI.getProfiles(teamId), {
    select: (data) => data.profiles,
    onSuccess: (data) => {
      setProfiles(data);
    },
    refetchOnWindowFocus: false,
  });

  const { data: config, refetch: refetchConfig } = useQuery<IConfig, Error>(
    ["config"],
    () => {
      return configAPI.loadAll();
    },
    {
      onSuccess: (data) => {
        setConfig(data);
      },
    }
  );

  const DEFAULT_SETTINGS_SECTION = MAC_OS_SETTINGS_NAV_ITEMS[0];

  const currentFormSection =
    MAC_OS_SETTINGS_NAV_ITEMS.find((item) => item.urlSection === section) ??
    DEFAULT_SETTINGS_SECTION;

  const CurrentCard = currentFormSection.Card;

  return (
    <div className={baseClass}>
      <p className={`${baseClass}__description`}>
        Remotely enforce settings on macOS hosts assigned to this team.
      </p>
      {profiles && <AggregateMacSettingsIndicators teamId={teamId} />}
      <SideNav
        className={`${baseClass}__side-nav`}
        navItems={MAC_OS_SETTINGS_NAV_ITEMS}
        activeItem={currentFormSection.urlSection}
        CurrentCard={
          <CurrentCard
            key={team_id}
            currentTeamId={teamId}
            profiles={profiles}
            refetchProfiles={refetchProfiles}
            config={config}
            refetchConfig={refetchConfig}
          />
        }
      />
    </div>
  );
};

export default MacOSSettings;

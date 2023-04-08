import React, { useCallback, useContext } from "react";
import { Tab, Tabs, TabList } from "react-tabs";
import { InjectedRouter } from "react-router";

import PATHS from "router/paths";
import { AppContext } from "context/app";
import useTeamIdParam from "hooks/useTeamIdParam";

import TabsWrapper from "components/TabsWrapper";
import MainContent from "components/MainContent";
import TeamsDropdown from "components/TeamsDropdown";
import EmptyTable from "components/EmptyTable";
import Button from "components/buttons/Button";
import permissionUtils from "utilities/permissions";

interface IControlsSubNavItem {
  name: string;
  pathname: string;
}

const controlsSubNav: IControlsSubNavItem[] = [
  {
    name: "macOS updates",
    pathname: PATHS.CONTROLS_MAC_OS_UPDATES,
  },
  {
    name: "macOS settings",
    pathname: PATHS.CONTROLS_MAC_SETTINGS,
  },
];

interface IManageControlsPageProps {
  children: JSX.Element;
  location: {
    pathname: string;
    search: string;
    hash?: string;
    query: {
      team_id?: string;
    };
  };
  router: InjectedRouter; // v3
}

const getTabIndex = (path: string): number => {
  return controlsSubNav.findIndex((navItem) => {
    // tab stays highlighted for paths that start with same pathname
    return path.startsWith(navItem.pathname);
  });
};

const baseClass = "manage-controls-page";

const ManageControlsPage = ({
  // TODO(sarah): decide on pattern to pass team id to subcomponents.
  // using children makes it difficult to centralize page-level control
  // over team id param
  children,
  location,
  router,
}: IManageControlsPageProps): JSX.Element => {
  const {
    availableTeams,
    config,
    currentUser,
    isFreeTier,
    isOnGlobalTeam,
    isPremiumTier,
    isGlobalAdmin,
  } = useContext(AppContext);

  const { currentTeamId, handleTeamChange } = useTeamIdParam({
    location,
    router,
    includeAllTeams: false,
    includeNoTeam: true,
  });

  // Filter out any teams that the user is not a maintainer or admin of
  // TODO: Abstract this to the AppContext so that availableTeams is already filtered
  // based on the current page
  let userAdminOrMaintainerTeams = availableTeams;
  userAdminOrMaintainerTeams = availableTeams?.filter(
    (team) =>
      permissionUtils.isTeamMaintainerOrTeamAdmin(currentUser, team.id) === true
  );

  const navigateToNav = useCallback(
    (i: number): void => {
      const navPath = controlsSubNav[i].pathname;
      router.replace(
        navPath.concat(location?.search || "").concat(location?.hash || "")
      );
    },
    [location, router]
  );

  const onConnectClick = () => {
    router.push(PATHS.ADMIN_INTEGRATIONS_MDM);
  };

  const renderConnectButton = () => {
    if (isGlobalAdmin) {
      return (
        <Button
          variant="brand"
          onClick={onConnectClick}
          className={`${baseClass}__connectAPC-button`}
        >
          Connect
        </Button>
      );
    }
    return <></>;
  };

  const getInfoText = () => {
    if (isGlobalAdmin) {
      return "Connect Fleet to the Apple Push Certificates Portal to get started.";
    }
    return "Your Fleet administrator must connect Fleet to the Apple Push Certificates Portal to get started.";
  };

  const renderBody = () => {
    return config?.mdm.enabled_and_configured ? (
      <div>
        <TabsWrapper>
          <Tabs
            selectedIndex={getTabIndex(location?.pathname || "")}
            onSelect={navigateToNav}
          >
            <TabList>
              {controlsSubNav.map((navItem) => {
                return (
                  <Tab key={navItem.name} data-text={navItem.name}>
                    {navItem.name}
                  </Tab>
                );
              })}
            </TabList>
          </Tabs>
        </TabsWrapper>
        {children}
      </div>
    ) : (
      <EmptyTable
        header="Manage your macOS hosts"
        info={getInfoText()}
        primaryButton={renderConnectButton()}
      />
    );
  };

  return (
    <MainContent>
      <div className={`${baseClass}__wrapper`}>
        <div className={`${baseClass}__header-wrap`}>
          <div className={`${baseClass}__header-wrap`}>
            <div className={`${baseClass}__header`}>
              <div className={`${baseClass}__text`}>
                <div className={`${baseClass}__title`}>
                  {isFreeTier && <h1>Controls</h1>}
                  {isPremiumTier &&
                    userAdminOrMaintainerTeams &&
                    (userAdminOrMaintainerTeams.length > 1 ||
                      isOnGlobalTeam) && (
                      <TeamsDropdown
                        currentUserTeams={userAdminOrMaintainerTeams}
                        selectedTeamId={currentTeamId}
                        onChange={handleTeamChange}
                        includeAll={false}
                        includeNoTeams
                      />
                    )}
                  {isPremiumTier &&
                    !isOnGlobalTeam &&
                    userAdminOrMaintainerTeams &&
                    userAdminOrMaintainerTeams.length === 1 && (
                      <h1>{userAdminOrMaintainerTeams[0].name}</h1>
                    )}
                </div>
              </div>
            </div>
          </div>
        </div>
        {renderBody()}
      </div>
    </MainContent>
  );
};

export default ManageControlsPage;

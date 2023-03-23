import { useCallback, useContext, useEffect, useMemo } from "react";
import { InjectedRouter } from "react-router";
import { findLastIndex, trimStart } from "lodash";

import { AppContext } from "context/app";
import {
  APP_CONTEXT_ALL_TEAMS_ID,
  isAnyTeamSelected,
  ITeamSummary,
  APP_CONTEXT_NO_TEAM_ID,
} from "interfaces/team";

const coerceAllTeamsId = (s: string) => {
  // URLs for the app represent "All teams" by the absence of the team id param
  // "All teams" is represented in AppContext with -1 as the team id so empty
  // strings are coerced to -1 by this function
  return s.length ? parseInt(s, 10) : APP_CONTEXT_ALL_TEAMS_ID;
};

const splitQueryStringParts = (queryString: string) =>
  trimStart(queryString, "?")
    .split("&")
    .filter((p) => p.includes("="));

const joinQueryStringParts = (parts: string[]) =>
  parts.length ? `?${parts.join("&")}` : "";

const rebuildQueryStringWithTeamId = (
  queryString: string,
  newTeamId: number
) => {
  const parts = splitQueryStringParts(queryString);
  const teamIndex = parts.findIndex((p) => p.startsWith("team_id="));

  // URLs for the app represent "All teams" by the absence of the team id param
  const newTeamPart =
    newTeamId > APP_CONTEXT_ALL_TEAMS_ID ? `team_id=${newTeamId}` : "";

  if (teamIndex < 0) {
    // nothing to remove/replace so add the new part (if any) and rejoin
    return joinQueryStringParts(
      newTeamPart ? parts.concat(newTeamPart) : parts
    );
  }

  if (teamIndex !== findLastIndex(parts, (p) => p.startsWith("team_id="))) {
    console.warn(
      `URL contains more than one team_id parameter: ${queryString}`
    );
  }

  if (newTeamPart) {
    parts.splice(teamIndex, 1, newTeamPart); // remove the old part and replace with the new
  } else {
    parts.splice(teamIndex, 1); // just remove the old team part
  }
  return joinQueryStringParts(parts);
};

const getTeamIdForApi = ({
  currentTeam,
  includeNoTeam = false,
}: {
  currentTeam?: ITeamSummary;
  includeNoTeam?: boolean;
}) => {
  if (includeNoTeam && currentTeam?.id === APP_CONTEXT_NO_TEAM_ID) {
    return APP_CONTEXT_NO_TEAM_ID;
  }
  if (currentTeam && currentTeam.id > APP_CONTEXT_NO_TEAM_ID) {
    return currentTeam.id;
  }
  return undefined; // API will treat undefined as equivalent to "All teams"
};

const getDefaultTeam = (
  availableTeams: ITeamSummary[],
  includeAllTeams: boolean,
  includeNoTeam: boolean
) => {
  let defaultTeam: ITeamSummary | undefined = availableTeams[0]; // TODO(sarah): can this be improved?
  if (includeAllTeams) {
    defaultTeam =
      availableTeams.find((t) => t.id === APP_CONTEXT_ALL_TEAMS_ID) ||
      defaultTeam;
  } else if (includeNoTeam) {
    defaultTeam =
      availableTeams.find((t) => t.id === APP_CONTEXT_NO_TEAM_ID) ||
      defaultTeam;
  } else {
    defaultTeam =
      availableTeams.find((t) => t.id > APP_CONTEXT_NO_TEAM_ID) || defaultTeam;
  }
  return defaultTeam;
};

const isValidTeamId = ({
  availableTeams,
  includeAllTeams,
  includeNoTeam,
  teamId,
}: {
  availableTeams: ITeamSummary[];
  includeAllTeams: boolean;
  includeNoTeam: boolean;
  teamId: number;
}) => {
  if (
    (teamId === APP_CONTEXT_ALL_TEAMS_ID && !includeAllTeams) ||
    (teamId === APP_CONTEXT_NO_TEAM_ID && !includeNoTeam) ||
    !availableTeams?.find((t) => t.id === teamId)
  ) {
    return false;
  }
  return true;
};

const shouldRedirectToDefaultTeam = ({
  availableTeams,
  includeAllTeams,
  includeNoTeam,
  query,
}: {
  availableTeams: ITeamSummary[];
  defaultTeamId: number;
  includeAllTeams: boolean;
  includeNoTeam: boolean;
  query: { team_id?: string };
}) => {
  const teamIdString = query?.team_id || "";
  const parsedTeamId = parseInt(teamIdString, 10);

  // redirect non-numeric strings and negative numbers to default (e.g., `/hosts?team_id=-1` should
  // be redirected to `/hosts`)
  if (teamIdString.length && (isNaN(parsedTeamId) || parsedTeamId < 0)) {
    return true;
  }

  // coerce empty string to -1 (i.e. `ALL_TEAMS_ID`) and test again (this ensures that non-global users will be
  // redirected to their default team when they attempt to access the `/hosts` page and also ensures
  // all users are redirected to their default when they attempt to acess non-existent team ids).
  return !isValidTeamId({
    availableTeams,
    includeAllTeams,
    includeNoTeam,
    teamId: coerceAllTeamsId(teamIdString),
  });
};

export const useTeamIdParam = ({
  location = { pathname: "", search: "", hash: "", query: {} },
  router,
  includeAllTeams,
  includeNoTeam,
}: {
  location?: {
    pathname: string;
    search: string;
    query: { team_id?: string };
    hash?: string;
  };
  router: InjectedRouter;
  includeAllTeams: boolean;
  includeNoTeam: boolean;
}) => {
  const { hash, pathname, query, search } = location;
  const { currentTeam, availableTeams, setCurrentTeam } = useContext(
    AppContext
  );

  const memoizedDefaultTeam = useMemo(() => {
    if (!availableTeams?.length) {
      return undefined;
    }
    return getDefaultTeam(availableTeams, includeAllTeams, includeNoTeam);
  }, [availableTeams, includeAllTeams, includeNoTeam]);

  const handleTeamChange = useCallback(
    (teamId: number) => {
      router.replace(
        pathname
          .concat(rebuildQueryStringWithTeamId(search, teamId))
          .concat(hash || "")
      );
    },
    [pathname, search, hash, router]
  );

  useEffect(() => {
    if (!availableTeams?.length || !memoizedDefaultTeam) {
      return; // skip effect until these values are available
    }

    // first reconcile router location and redirect to default team as applicable
    if (
      shouldRedirectToDefaultTeam({
        availableTeams,
        defaultTeamId: memoizedDefaultTeam.id,
        includeAllTeams,
        includeNoTeam,
        query,
      })
    ) {
      handleTeamChange(memoizedDefaultTeam.id);
      return; // allow router to update first before updating context
    }

    // location is resolved, so proceed to update team context
    const teamId = coerceAllTeamsId(query?.team_id || "");
    if (teamId !== currentTeam?.id) {
      setCurrentTeam(availableTeams?.find((t) => t.id === teamId));
    }
  }, [
    availableTeams,
    currentTeam,
    includeAllTeams,
    includeNoTeam,
    memoizedDefaultTeam,
    query,
    search,
    handleTeamChange,
    setCurrentTeam,
  ]);

  return {
    currentTeamId: currentTeam?.id,
    currentTeamName: currentTeam?.name,
    isAnyTeamSelected: isAnyTeamSelected(currentTeam),
    teamIdForApi: getTeamIdForApi({ currentTeam, includeNoTeam }),
    handleTeamChange,
  };
};

export default useTeamIdParam;

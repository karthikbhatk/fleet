import React, { useState, useEffect, useContext, useCallback } from "react";
import { useQuery } from "react-query";
import { useErrorHandler } from "react-error-boundary";
import { InjectedRouter, Params } from "react-router/lib/Router";

import { AppContext } from "context/app";
import { QueryContext } from "context/query";
import { QUERIES_PAGE_STEPS, DEFAULT_QUERY } from "utilities/constants";
import queryAPI from "services/entities/queries";
import hostAPI from "services/entities/hosts";
import statusAPI from "services/entities/status";
import { IHost, IHostResponse } from "interfaces/host";
import { ILabel } from "interfaces/label";
import { ITeam } from "interfaces/team";
import {
  IGetQueryResponse,
  ISchedulableQuery,
} from "interfaces/schedulable_query";
import { ITarget } from "interfaces/target";

import MainContent from "components/MainContent";
import SelectTargets from "components/LiveQuery/SelectTargets";
import CustomLink from "components/CustomLink";

import EditQueryPage from "pages/queries/edit/EditQueryPage/EditQueryPage";
import RunQuery from "pages/queries/run/screens/RunQuery";
import useTeamIdParam from "hooks/useTeamIdParam";

interface IRunQueryPageProps {
  router: InjectedRouter;
  params: Params;
  location: {
    pathname: string;
    query: { host_ids: string; team_id?: string };
    search: string;
  };
}

const baseClass = "run-query-page";

const RunQueryPage = ({
  router,
  params: { id: paramsQueryId },
  location,
}: IRunQueryPageProps): JSX.Element => {
  const queryId = paramsQueryId ? parseInt(paramsQueryId, 10) : null;
  const {
    currentTeamName: teamNameForQuery,
    teamIdForApi: apiTeamIdForQuery,
  } = useTeamIdParam({
    location,
    router,
    includeAllTeams: true,
    includeNoTeam: false,
  });

  const handlePageError = useErrorHandler();
  const {
    isGlobalAdmin,
    isGlobalMaintainer,
    isAnyTeamMaintainerOrTeamAdmin,
    isObserverPlus,
    isAnyTeamObserverPlus,
  } = useContext(AppContext);
  const {
    selectedOsqueryTable,
    setSelectedOsqueryTable,
    setLastEditedQueryId,
    setLastEditedQueryName,
    setLastEditedQueryDescription,
    setLastEditedQueryBody,
    setLastEditedQueryObserverCanRun,
    setLastEditedQueryFrequency,
    setLastEditedQueryLoggingType,
    setLastEditedQueryMinOsqueryVersion,
    setLastEditedQueryPlatforms,
  } = useContext(QueryContext);

  const [queryParamHostsAdded, setQueryParamHostsAdded] = useState(false);
  const [step, setStep] = useState(QUERIES_PAGE_STEPS[1]);
  const [selectedTargets, setSelectedTargets] = useState<ITarget[]>([]);
  const [targetedHosts, setTargetedHosts] = useState<IHost[]>([]);
  const [targetedLabels, setTargetedLabels] = useState<ILabel[]>([]);
  const [targetedTeams, setTargetedTeams] = useState<ITeam[]>([]);
  const [targetsTotalCount, setTargetsTotalCount] = useState(0);
  const [isLiveQueryRunnable, setIsLiveQueryRunnable] = useState(true);

  // disabled on page load so we can control the number of renders
  // else it will re-populate the context on occasion
  const {
    isLoading: isStoredQueryLoading,
    data: storedQuery,
    error: storedQueryError,
  } = useQuery<IGetQueryResponse, Error, ISchedulableQuery>(
    ["query", queryId],
    () => queryAPI.load(queryId as number),
    {
      enabled: !!queryId,
      refetchOnWindowFocus: false,
      select: (data) => data.query,
      onSuccess: (returnedQuery) => {
        setLastEditedQueryId(returnedQuery.id);
        setLastEditedQueryName(returnedQuery.name);
        setLastEditedQueryDescription(returnedQuery.description);
        setLastEditedQueryBody(returnedQuery.query);
        setLastEditedQueryObserverCanRun(returnedQuery.observer_can_run);
        setLastEditedQueryFrequency(returnedQuery.interval);
        setLastEditedQueryPlatforms(returnedQuery.platform);
        setLastEditedQueryLoggingType(returnedQuery.logging);
        setLastEditedQueryMinOsqueryVersion(returnedQuery.min_osquery_version);
      },
      onError: (error) => handlePageError(error),
    }
  );

  useQuery<IHostResponse, Error, IHost>(
    "hostFromURL",
    () =>
      hostAPI.loadHostDetails(parseInt(location.query.host_ids as string, 10)),
    {
      enabled: !!location.query.host_ids && !queryParamHostsAdded,
      select: (data: IHostResponse) => data.host,
      onSuccess: (host) => {
        setTargetedHosts((prevHosts) =>
          prevHosts.filter((h) => h.id !== host.id).concat(host)
        );
        const targets = selectedTargets;
        host.target_type = "hosts";
        targets.push(host);
        setSelectedTargets([...targets]);
        if (!queryParamHostsAdded) {
          setQueryParamHostsAdded(true);
        }
        router.replace(location.pathname);
      },
    }
  );

  const detectIsFleetQueryRunnable = () => {
    statusAPI.live_query().catch(() => {
      setIsLiveQueryRunnable(false);
    });
  };

  useEffect(() => {
    detectIsFleetQueryRunnable();
    if (!queryId) {
      setLastEditedQueryId(DEFAULT_QUERY.id);
      setLastEditedQueryName(DEFAULT_QUERY.name);
      setLastEditedQueryDescription(DEFAULT_QUERY.description);
      setLastEditedQueryBody(DEFAULT_QUERY.query);
      setLastEditedQueryObserverCanRun(DEFAULT_QUERY.observer_can_run);
      setLastEditedQueryFrequency(DEFAULT_QUERY.interval);
      setLastEditedQueryLoggingType(DEFAULT_QUERY.logging);
      setLastEditedQueryMinOsqueryVersion(DEFAULT_QUERY.min_osquery_version);
      setLastEditedQueryPlatforms(DEFAULT_QUERY.platform);
    }
  }, [queryId]);

  // Updates title that shows up on browser tabs
  useEffect(() => {
    // e.g., Run live query | Discover TLS certificates | Fleet for osquery
    document.title = `Run live query | ${storedQuery?.name} | Fleet for osquery`;
  }, [location.pathname, storedQuery?.name]);

  const goToQueryEditor = useCallback(() => setStep(QUERIES_PAGE_STEPS[1]), []);
  const params = { id: paramsQueryId };

  const renderScreen = () => {
    const step1Props = {
      baseClass,
      queryId,
      selectedTargets,
      targetedHosts,
      targetedLabels,
      targetedTeams,
      targetsTotalCount,
      goToQueryEditor: () => setStep(QUERIES_PAGE_STEPS[1]),
      goToRunQuery: () => setStep(QUERIES_PAGE_STEPS[3]),
      setSelectedTargets,
      setTargetedHosts,
      setTargetedLabels,
      setTargetedTeams,
      setTargetsTotalCount,
    };

    const step2Props = {
      queryId,
      selectedTargets,
      storedQuery,
      setSelectedTargets,
      goToQueryEditor,
      targetsTotalCount,
    };

    switch (step) {
      case QUERIES_PAGE_STEPS[1]:
        return <SelectTargets {...step1Props} />;
      case QUERIES_PAGE_STEPS[2]:
        return <RunQuery {...step2Props} />;
      default:
        return <SelectTargets {...step1Props} />;
    }
  };

  return (
    <>
      <MainContent className={baseClass}>
        <div className={`${baseClass}_wrapper`}>{renderScreen()}</div>
      </MainContent>
    </>
  );
};

export default RunQueryPage;

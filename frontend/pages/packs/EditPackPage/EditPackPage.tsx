import React, { useState, useEffect, useCallback, useContext } from "react";
import { useQuery, useMutation } from "react-query";
import { Params } from "react-router/lib/Router";

import { filter, includes, isEqual, noop, size, find } from "lodash";
import { useDispatch } from "react-redux";
import { push } from "react-router-redux";

// second grouping
// @ts-ignore
import { IConfig } from "interfaces/config";
import { IError } from "interfaces/errors";
import { IHost } from "interfaces/host";
import { ILabel } from "interfaces/label";
import { IPack } from "interfaces/pack";
import { IQuery } from "interfaces/query";
import { IScheduledQuery } from "interfaces/scheduled_query";
import { ITeam } from "interfaces/team";
import { AppContext } from "context/app";

import configAPI from "services/entities/config";
import hostAPI from "services/entities/hosts";
import labelAPI from "services/entities/labels";
import packAPI from "services/entities/packs";
import queryAPI from "services/entities/queries";
import scheduledqueryAPI from "services/entities/scheduled_queries";
import teamAPI from "services/entities/teams";

// @ts-ignore
import deepDifference from "utilities/deep_difference";
// @ts-ignore
import EditPackFormWrapper from "components/packs/EditPackFormWrapper";
import PackQueriesListWrapper from "components/queries/PackQueriesListWrapper";
import PackQueryEditorModal from "./components/PackQueryEditorModal";
import RemovePackQueryModal from "./components/RemovePackQueryModal";
// @ts-ignore
import { renderFlash } from "redux/nodes/notifications/actions";
// @ts-ignore
import debounce from "utilities/debounce";
import PATHS from "router/paths";

interface IEditPacksPageProps {
  router: any;
  params: Params;
  location: any; // TODO: find Location type
}
interface IRootState {
  app: {
    config: IConfig;
  };
  entities: {
    packs: {
      loading: boolean; // done
      data: IPack[];
      errors: IError[];
    };
    hosts: {
      isLoading: boolean;
      data: IHost[];
      errors: IError[];
    };
    queries: {
      isLoading: boolean;
      data: IQuery[];
      errors: IError[];
    };
    teams: {
      isLoading: boolean;
      data: ITeam[];
      errors: IError[];
    };
    labels: {
      isLoading: boolean;
      data: ILabel[];
      errors: IError[];
    };
    scheduled_queries: {
      isLoading: boolean;
      data: IScheduledQuery[];
      errors: IError[];
    };
  };
}

interface IPackQueryFormData {
  interval: number;
  name?: string;
  shard: number;
  query?: string;
  query_id?: number;
  removed: boolean;
  snapshot: boolean;
  pack_id: number;
  platform: string;
  version: string;
}
interface IStoredPackResponse {
  pack: IPack;
}

interface IStoredFleetQueriesResponse {
  queries: IQuery[];
}

interface IStoredScheduledQueriesResponse {
  scheduled: IScheduledQuery[];
}

interface IStoredLabelsResponse {
  labels: ILabel[];
}
interface IStoredHostsResponse {
  hosts: IHost[];
}

interface IStoredTeamsResponse {
  teams: ITeam[];
}

const baseClass = "edit-pack-page";

const EditPacksPage = ({
  router, // only needed if I need to navigate to another page from this page
  params: { id: paramsPackId },
  location: { query: URLQueryString }, // might need this if there's team filters
}: IEditPacksPageProps): JSX.Element => {
  // DATA AND API CALLS
  const { isPremiumTier } = useContext(AppContext);

  const dispatch = useDispatch();
  const packId: number = parseInt(paramsPackId, 10);

  const [targetsCount, setTargetsCount] = useState<number>(0);
  const [
    showPackQueryEditorModal,
    setShowPackQueryEditorModal,
  ] = useState<boolean>(false);
  const [showEditPackQueryModal, setShowEditPackQueryModal] = useState<boolean>(
    false
  );
  const [
    showRemovePackQueryModal,
    setShowRemovePackQueryModal,
  ] = useState<boolean>(false);
  const [selectedPackQuery, setSelectedPackQuery] = useState<IScheduledQuery>();
  const [selectedPackQueryIds, setSelectedPackQueryIds] = useState<
    number[] | never[]
  >([]);

  // react-query uses your own api and gives you different states of loading data
  // can set to retreive data based on different properties
  const {
    isLoading: isStoredPackLoading,
    data: storedPack, // only returns pack and not response wrapping
    error: storedPackError,
  } = useQuery<IStoredPackResponse, Error, IPack>(
    ["pack", packId],
    () => packAPI.load(packId),
    {
      enabled: !!packId, // doesn't run unless ID is given, unneeded but extra precaution
      select: (data: IStoredPackResponse) => data.pack,
    }
  );

  const {
    isLoading: isFleetQueriesLoading,
    data: fleetQueries,
    error: fleetQueriesError,
  } = useQuery<IStoredFleetQueriesResponse, Error, IQuery[]>(
    ["fleet queries"], // use single string or array of strings can be named anything
    () => queryAPI.loadAll(),
    {
      select: (data: IStoredFleetQueriesResponse) => data.queries,
    }
  );

  const {
    isLoading: isScheduledQueriesLoading,
    data: scheduledQueries,
    error: scheduledQueriesError,
  } = useQuery<IStoredScheduledQueriesResponse, Error, IScheduledQuery[]>(
    ["scheduled queries"], // use single string or array of strings can be named anything
    () => scheduledqueryAPI.loadAll(packId), // TODO: help with types
    {
      select: (data: IStoredScheduledQueriesResponse) => data.scheduled,
    }
  );

  const {
    isLoading: isLabelsLoading,
    data: labels,
    error: packLabelsError,
  } = useQuery<IStoredLabelsResponse, Error, ILabel[]>(
    ["pack labels"], // use single string or array of strings can be named anything
    () => labelAPI.loadAll(),
    {
      select: (data: IStoredLabelsResponse) => data.labels,
    }
  );

  const packLabels = storedPack
    ? filter(labels, (label) => {
        return includes(storedPack.label_ids, label.id);
      })
    : [];

  const {
    isLoading: isHostsLoading,
    data: hosts,
    error: hostsError,
  } = useQuery<IStoredHostsResponse, Error, IHost[]>(
    ["all hosts"], // use single string or array of strings can be named anything
    () => hostAPI.loadAll(undefined),
    {
      select: (data: IStoredHostsResponse) => data.hosts,
    }
  );
  const packHosts = storedPack
    ? filter(hosts, (host) => {
        return includes(storedPack.host_ids, host.id);
      })
    : [];

  const {
    isLoading: isTeamsLoading,
    data: teams,
    error: teamsError,
  } = useQuery<IStoredTeamsResponse, Error, ITeam[]>(
    ["pack labels"],
    () => teamAPI.loadAll(),
    {
      select: (data: IStoredTeamsResponse) => data.teams,
    }
  );

  const packTeams = storedPack
    ? filter(teams, (team) => {
        return includes(storedPack.team_ids, team.id);
      })
    : [];

  console.log("packLabels", packLabels);
  console.log("labels", labels);
  console.log("packLabelsError", packLabelsError);
  console.log("isLabelsLoading", isLabelsLoading);
  console.log("packHosts", packHosts);
  console.log("packTeams", packTeams);
  console.log("scheduledQueries", scheduledQueries);
  console.log("scheduledQueriesError", scheduledQueriesError);
  console.log("isScheduledQueriesLoading", isScheduledQueriesLoading);
  console.log("fleetQueries", fleetQueries);
  console.log("fleetQueriesError", fleetQueriesError);
  console.log("isFleetQueriesLoading", isFleetQueriesLoading);

  const packTargets = [...packHosts, ...packLabels, ...packTeams];

  // // FUNCTIONS

  const onCancelEditPack = () => {
    return dispatch(push(PATHS.MANAGE_PACKS));
  };

  const onFetchTargets = (query: IQuery, targetsResponse: any) => {
    // TODO: fix type issue
    const { targets_count: targetsCount } = targetsResponse;

    setTargetsCount(targetsCount);

    return false;
  };

  const onEditPackQueryClick = (selectedQuery: any): void => {
    togglePackQueryEditorModal();
    setSelectedPackQuery(selectedQuery); // edit modal renders
  };

  const onRemovePackQueriesClick = (selectedTableQueryIds: any): void => {
    toggleRemovePackQueryModal();
    setSelectedPackQueryIds(selectedTableQueryIds);
  };

  const togglePackQueryEditorModal = useCallback(() => {
    setSelectedPackQuery(undefined); // create modal renders
    setShowPackQueryEditorModal(!showPackQueryEditorModal);
    console.log("togglePackQueryEditorModal clicked!");
  }, [showPackQueryEditorModal, setShowPackQueryEditorModal]);

  const toggleEditPackQueryModal = useCallback(() => {
    setShowEditPackQueryModal(!showEditPackQueryModal);
  }, [showEditPackQueryModal, setShowEditPackQueryModal]);

  const toggleRemovePackQueryModal = useCallback(() => {
    setShowRemovePackQueryModal(!showRemovePackQueryModal);
  }, [showRemovePackQueryModal, setShowRemovePackQueryModal]);

  const handlePackFormSubmit = (formData: any) => {
    const updatedPack = deepDifference(formData, storedPack);
    packAPI
      .update(packId, updatedPack)
      .then(() => {
        toggleEditPackQueryModal();
      })
      .catch(() => {
        dispatch(
          renderFlash("error", `Could not update pack. Please try again.`)
        );
      });
  };

  const {
    mutateAsync: createPackQuery,
  } = useMutation((formData: IPackQueryFormData) =>
    scheduledqueryAPI.create(formData)
  );

  // const onSavePackQueryFormSubmit = debounce(
  //   async (formData: IPackQueryFormData) => {
  //     try {
  //       const { query }: { query: IScheduledQuery } = await createQuery(formData);
  //       router.push(PATHS.EDIT_QUERY(query));
  //       dispatch(renderFlash("success", "Query created!"));
  //     } catch (createError) {
  //       console.error(createError);
  //       dispatch(
  //         renderFlash(
  //           "error",
  //           "Something went wrong creating your query. Please try again."
  //         )
  //       );
  //     }
  //   }
  // );

  // const onPackQueryEditorSubmit = (formData: IPackQueryFormData) => {
  // const { dispatch } = this.props;
  // const { selectedScheduledQuery } = this.state;
  // const { update } = scheduledQueryActions;
  // const updatedAttrs = deepDifference(formData, selectedScheduledQuery);

  // dispatch(update(selectedScheduledQuery, updatedAttrs))
  //   .then(() => {
  //     this.setState({ selectedScheduledQuery: null, selectedQuery: null });
  //     dispatch(renderFlash("success", "Scheduled Query updated!"));
  //   })
  //   .catch(() => {
  //     dispatch(
  //       renderFlash("error", "Unable to update your Scheduled Query.")
  //     );
  //   });
  // };

  const onPackQueryEditorSubmit = useCallback(
    (formData: IPackQueryFormData, editQuery: IScheduledQuery | undefined) => {
      // if (editQuery) {
      //   const updatedAttributes = deepDifference(formData, editQuery);
      //   dispatch(
      //     globalScheduledQueryActions.update(editQuery, updatedAttributes)
      //   )
      //     .then(() => {
      //       dispatch(
      //         renderFlash(
      //           "success",
      //           `Successfully updated ${formData.name} in the schedule.`
      //         )
      //       );
      //       dispatch(globalScheduledQueryActions.loadAll());
      //     })
      //     .catch(() => {
      //       dispatch(
      //         renderFlash(
      //           "error",
      //           "Could not update scheduled query. Please try again."
      //         )
      //       );
      //     });
      // } else {
      //   dispatch(globalScheduledQueryActions.create({ ...formData }))
      //     .then(() => {
      //       dispatch(
      //         renderFlash(
      //           "success",
      //           `Successfully added ${formData.name} to the schedule.`
      //         )
      //       );
      //       dispatch(globalScheduledQueryActions.loadAll());
      //     })
      //     .catch(() => {
      //       dispatch(
      //         renderFlash(
      //           "error",
      //           "Could not schedule query. Please try again."
      //         )
      //       );
      //     });
      // }
      togglePackQueryEditorModal();
    },
    [dispatch, togglePackQueryEditorModal]
  );

  const onRemovePackQuerySubmit = useCallback(() => {
    const ids = selectedPackQueryIds;
    scheduledqueryAPI
      .destroy(ids[0]) // TODO: ugh, destroy multiple
      .then(() => {
        dispatch(
          renderFlash(
            "success",
            `Successfully removed ${
              ids && ids.length === 1 ? "policy" : "policies"
            }.`
          )
        );
      })
      .catch(() => {
        dispatch(
          renderFlash(
            "error",
            `Unable to remove ${
              ids && ids.length === 1 ? "policy" : "policies"
            }. Please try again.`
          )
        );
      })
      .finally(() => {
        toggleRemovePackQueryModal();
        // getPackQueries();
      });
  }, [
    dispatch,
    // getPackQueries,
    selectedPackQueryIds,
    toggleRemovePackQueryModal,
  ]);

  return (
    <div className={`${baseClass}__content`}>
      {storedPack && (
        <EditPackFormWrapper
          className={`${baseClass}__pack-form body-wrap`}
          handleSubmit={handlePackFormSubmit}
          onCancelEditPack={onCancelEditPack}
          onEditPack={toggleEditPackQueryModal}
          onFetchTargets={onFetchTargets}
          pack={storedPack}
          packTargets={packTargets}
          targetsCount={targetsCount}
          isPremiumTier={isPremiumTier}
        />
      )}
      <PackQueriesListWrapper
        onAddPackQuery={togglePackQueryEditorModal}
        onEditPackQuery={onEditPackQueryClick}
        onRemovePackQueries={onRemovePackQueriesClick}
        onPackQueryFormSubmit={onPackQueryEditorSubmit}
        scheduledQueries={scheduledQueries}
        packId={packId}
        isLoadingPackQueries={isScheduledQueriesLoading}
      />
      {showPackQueryEditorModal && fleetQueries && (
        <PackQueryEditorModal
          onCancel={togglePackQueryEditorModal}
          onPackQueryFormSubmit={onPackQueryEditorSubmit}
          allQueries={fleetQueries}
          editQuery={selectedPackQuery}
          packId={packId}
        />
      )}
      {showRemovePackQueryModal && fleetQueries && (
        <RemovePackQueryModal
          onCancel={toggleRemovePackQueryModal}
          onSubmit={onRemovePackQuerySubmit}
          selectedQueries={selectedPackQuery}
        />
      )}
    </div>
  );
};

export default EditPacksPage;

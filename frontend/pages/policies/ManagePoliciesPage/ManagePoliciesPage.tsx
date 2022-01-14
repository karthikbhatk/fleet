import React, { useCallback, useContext, useEffect, useState } from "react";
import { useQuery } from "react-query";
import { useDispatch } from "react-redux";
import { noop } from "lodash";

import { AppContext } from "context/app";
import { PolicyContext } from "context/policy";
import { inMilliseconds, secondsToHms } from "fleet/helpers";
import { IPolicyStats, ILoadAllPoliciesResponse } from "interfaces/policy";
import { IWebhookFailingPolicies } from "interfaces/webhook";
// @ts-ignore
import { renderFlash } from "redux/nodes/notifications/actions";
import PATHS from "router/paths";
import configAPI from "services/entities/config";
import globalPoliciesAPI from "services/entities/global_policies";
import teamPoliciesAPI from "services/entities/team_policies";
import { DEFAULT_POLICY } from "utilities/constants";

import Button from "components/buttons/Button";
import InfoBanner from "components/InfoBanner/InfoBanner";
import IconToolTip from "components/IconToolTip";
import Spinner from "components/Spinner";
import TeamsDropdown from "components/TeamsDropdown";
import TableDataError from "components/TableDataError";
import PoliciesListWrapper from "./components/PoliciesListWrapper";
import ManageAutomationsModal from "./components/ManageAutomationsModal";
import AddPolicyModal from "./components/AddPolicyModal";
import RemovePoliciesModal from "./components/RemovePoliciesModal";

interface IManagePoliciesPageProps {
  router: any;
  location: any;
  // params: any;
}

const baseClass = "manage-policies-page";

const DOCS_LINK =
  "https://fleetdm.com/docs/deploying/configuration#osquery-policy-update-interval";

const ManagePolicyPage = ({
  router,
  location,
}: IManagePoliciesPageProps): JSX.Element => {
  const dispatch = useDispatch();

  const {
    availableTeams,
    config,
    isGlobalAdmin,
    isGlobalMaintainer,
    isOnGlobalTeam,
    isFreeTier,
    isPremiumTier,
    isTeamAdmin,
    isTeamMaintainer,
    currentTeam,
    setCurrentTeam,
  } = useContext(AppContext);

  const teamId = parseInt(location?.query?.team_id, 10) || 0;

  const {
    setLastEditedQueryName,
    setLastEditedQueryDescription,
    setLastEditedQueryBody,
    setLastEditedQueryResolution,
    setLastEditedQueryPlatform,
  } = useContext(PolicyContext);

  const canAddOrRemovePolicy =
    isGlobalAdmin || isGlobalMaintainer || isTeamMaintainer || isTeamAdmin;

  const policyUpdateInterval =
    secondsToHms(inMilliseconds(config?.osquery_policy || 0) / 1000) ||
    "osquery policy update interval";

  const [selectedPolicyIds, setSelectedPolicyIds] = useState<
    number[] | never[]
  >([]);
  const [showManageAutomationsModal, setShowManageAutomationsModal] = useState(
    false
  );
  const [showPreviewPayloadModal, setShowPreviewPayloadModal] = useState(false);
  const [showAddPolicyModal, setShowAddPolicyModal] = useState(false);
  const [showRemovePoliciesModal, setShowRemovePoliciesModal] = useState(false);
  const [showInheritedPolicies, setShowInheritedPolicies] = useState(false);

  const [
    isLoadingFailingPoliciesWebhook,
    setIsLoadingFailingPoliciesWebhook,
  ] = useState(true);
  const [
    isFailingPoliciesWebhookError,
    setIsFailingPoliciesWebhookError,
  ] = useState(false);
  const [failingPoliciesWebhook, setFailingPoliciesWebhook] = useState<
    IWebhookFailingPolicies | undefined
  >();
  const [currentAutomatedPolicies, setCurrentAutomatedPolicies] = useState<
    number[]
  >();

  const {
    data: globalPolicies,
    error: globalPoliciesError,
    isLoading: isLoadingGlobalPolicies,
    isStale: isStaleGlobalPolicies,
    refetch: refetchGlobalPolicies,
  } = useQuery<ILoadAllPoliciesResponse, Error, IPolicyStats[]>(
    ["globalPolicies"],
    () => {
      return globalPoliciesAPI.loadAll();
    },
    {
      enabled: !!availableTeams,
      select: (data) => data.policies,
      onSuccess: () => setLastEditedQueryPlatform(""),
      staleTime: 3000,
    }
  );

  const {
    data: teamPolicies,
    error: teamPoliciesError,
    isLoading: isLoadingTeamPolicies,
    refetch: refetchTeamPolicies,
  } = useQuery<ILoadAllPoliciesResponse, Error, IPolicyStats[]>(
    ["teamPolicies", currentTeam?.id],
    () => teamPoliciesAPI.loadAll(currentTeam?.id),
    {
      enabled: isPremiumTier && !!currentTeam?.id,
      select: (data) => data.policies,
    }
  );

  const refetchPolicies = (id?: number) => {
    if (id) {
      refetchTeamPolicies();
      refetchGlobalPolicies();
    } else {
      refetchGlobalPolicies();
    }
  };

  const findAvailableTeam = (id: number) => {
    return availableTeams?.find((t) => t.id === id);
  };

  const handleTeamSelect = (id: number) => {
    const { MANAGE_POLICIES } = PATHS;

    const selectedTeam = findAvailableTeam(id);
    const path = selectedTeam?.id
      ? `${MANAGE_POLICIES}?team_id=${selectedTeam.id}`
      : MANAGE_POLICIES;

    router.replace(path);
    setShowInheritedPolicies(false);
    setSelectedPolicyIds([]);
    setCurrentTeam(selectedTeam);
    isStaleGlobalPolicies && refetchGlobalPolicies();
  };

  const getFailingPoliciesWebhook = useCallback(async () => {
    setIsLoadingFailingPoliciesWebhook(true);
    setIsFailingPoliciesWebhookError(false);
    let result;
    try {
      result = await configAPI
        .loadAll()
        .then((response) => response.webhook_settings.failing_policies_webhook);
      setFailingPoliciesWebhook(result);
      setCurrentAutomatedPolicies(result.policy_ids);
    } catch (error) {
      console.log(error);
      setIsFailingPoliciesWebhookError(true);
    } finally {
      setIsLoadingFailingPoliciesWebhook(false);
    }
    return result;
  }, []);

  const toggleManageAutomationsModal = () =>
    setShowManageAutomationsModal(!showManageAutomationsModal);

  const togglePreviewPayloadModal = useCallback(() => {
    setShowPreviewPayloadModal(!showPreviewPayloadModal);
  }, [setShowPreviewPayloadModal, showPreviewPayloadModal]);

  const toggleAddPolicyModal = () => setShowAddPolicyModal(!showAddPolicyModal);

  const toggleRemovePoliciesModal = () =>
    setShowRemovePoliciesModal(!showRemovePoliciesModal);

  const toggleShowInheritedPolicies = () =>
    setShowInheritedPolicies(!showInheritedPolicies);

  const onManageAutomationsClick = () => {
    toggleManageAutomationsModal();
  };

  const onCreateWebhookSubmit = async ({
    destination_url,
    policy_ids,
    enable_failing_policies_webhook,
  }: IWebhookFailingPolicies) => {
    try {
      const request = configAPI.update({
        webhook_settings: {
          failing_policies_webhook: {
            destination_url,
            policy_ids,
            enable_failing_policies_webhook,
          },
        },
      });
      await request.then(() => {
        dispatch(
          renderFlash("success", "Successfully updated policy automations.")
        );
      });
    } catch {
      dispatch(
        renderFlash(
          "error",
          "Could not update policy automations. Please try again."
        )
      );
    } finally {
      toggleManageAutomationsModal();
      getFailingPoliciesWebhook();
    }
  };

  const onAddPolicyClick = () => {
    setLastEditedQueryName("");
    setLastEditedQueryDescription("");
    setLastEditedQueryBody(DEFAULT_POLICY.query);
    setLastEditedQueryResolution("");
    toggleAddPolicyModal();
  };

  const onRemovePoliciesClick = (selectedTableIds: number[]): void => {
    toggleRemovePoliciesModal();
    setSelectedPolicyIds(selectedTableIds);
  };

  const onRemovePoliciesSubmit = async () => {
    const id = currentTeam?.id;
    try {
      const request = id
        ? teamPoliciesAPI.destroy(id, selectedPolicyIds)
        : globalPoliciesAPI.destroy(selectedPolicyIds);

      await request.then(() => {
        dispatch(
          renderFlash(
            "success",
            `Successfully removed ${
              selectedPolicyIds?.length === 1 ? "policy" : "policies"
            }.`
          )
        );
      });
    } catch {
      dispatch(
        renderFlash(
          "error",
          `Unable to remove ${
            selectedPolicyIds?.length === 1 ? "policy" : "policies"
          }. Please try again.`
        )
      );
    } finally {
      toggleRemovePoliciesModal();
      refetchPolicies(id);
    }
  };

  const showDefaultDescription =
    isFreeTier || (isPremiumTier && !!availableTeams && !teamId);

  const showInfoBanner =
    (teamId && !teamPoliciesError && !!teamPolicies?.length) ||
    (!teamId && !globalPoliciesError && !!globalPolicies?.length);

  const showInheritedPoliciesButton =
    !!teamId &&
    !isLoadingTeamPolicies &&
    !teamPoliciesError &&
    !isLoadingGlobalPolicies &&
    !globalPoliciesError &&
    !!globalPolicies?.length;

  const inheritedPoliciesButtonText = (
    showPolicies: boolean,
    count: number
  ) => {
    return `${showPolicies ? "Hide" : "Show"} ${count} inherited ${
      count > 1 ? "policies" : "policy"
    }`;
  };

  // Validate team_id from URL query params and redirect invalid cases to default policy page
  useEffect(() => {
    if (availableTeams !== null && availableTeams !== undefined) {
      let validatedId: number;
      if (findAvailableTeam(teamId)) {
        validatedId = teamId;
      } else if (!teamId && currentTeam) {
        validatedId = currentTeam.id;
      } else if (!teamId && !currentTeam && !isOnGlobalTeam && availableTeams) {
        validatedId = availableTeams[0]?.id;
      } else {
        validatedId = 0;
      }

      if (validatedId !== currentTeam?.id || validatedId !== teamId) {
        handleTeamSelect(validatedId);
      }
    }
  }, [availableTeams]);

  return !availableTeams ? (
    <Spinner />
  ) : (
    <div className={baseClass}>
      <div className={`${baseClass}__wrapper body-wrap`}>
        <div className={`${baseClass}__header-wrap`}>
          <div className={`${baseClass}__header`}>
            <div className={`${baseClass}__text`}>
              <div className={`${baseClass}__title`}>
                {isFreeTier && <h1>Policies</h1>}
                {isPremiumTier &&
                  (availableTeams.length > 1 || isOnGlobalTeam) && (
                    <TeamsDropdown
                      currentUserTeams={availableTeams || []}
                      selectedTeamId={teamId}
                      onChange={(newSelectedValue: number) =>
                        handleTeamSelect(newSelectedValue)
                      }
                    />
                  )}
                {isPremiumTier &&
                  !isOnGlobalTeam &&
                  availableTeams.length === 1 && (
                    <h1>{availableTeams[0].name}</h1>
                  )}
              </div>
            </div>
          </div>
          <div className={`${baseClass} button-wrap`}>
            {canAddOrRemovePolicy && teamId === 0 && (
              <Button
                onClick={() => onManageAutomationsClick()}
                className={`${baseClass}__manage-automations button`}
                variant="inverse"
              >
                <span>Manage automations</span>
              </Button>
            )}
            {canAddOrRemovePolicy && (
              <div className={`${baseClass}__action-button-container`}>
                <Button
                  variant="brand"
                  className={`${baseClass}__select-policy-button`}
                  onClick={onAddPolicyClick}
                >
                  Add a policy
                </Button>
              </div>
            )}
          </div>
        </div>
        <div className={`${baseClass}__description`}>
          {isPremiumTier && !!teamId && (
            <p>
              Add additional policies for <b>all hosts assigned to this team</b>
              .
            </p>
          )}
          {showDefaultDescription && (
            <p>
              Add policies for <b>all of your hosts</b> to see which pass your
              organization’s standards.{" "}
            </p>
          )}
        </div>
        {!!policyUpdateInterval && showInfoBanner && (
          <InfoBanner className={`${baseClass}__sandbox-info`}>
            <p>
              Your policies are checked every{" "}
              <b>{policyUpdateInterval.trim()}</b>.{" "}
              {isGlobalAdmin && (
                <span>
                  Check out the Fleet documentation on{" "}
                  <a href={DOCS_LINK} target="_blank" rel="noreferrer">
                    <b>how to edit this frequency</b>
                  </a>
                  .
                </span>
              )}
            </p>
          </InfoBanner>
        )}
        <div>
          {!!teamId && teamPoliciesError && <TableDataError />}
          {!!teamId && !teamPoliciesError && teamPolicies === undefined && (
            <Spinner />
          )}
          {!!teamId && !teamPoliciesError && teamPolicies !== undefined && (
            <PoliciesListWrapper
              policiesList={teamPolicies || []}
              isLoading={
                isLoadingTeamPolicies && isLoadingFailingPoliciesWebhook
              }
              onRemovePoliciesClick={onRemovePoliciesClick}
              canAddOrRemovePolicy={canAddOrRemovePolicy}
              currentTeam={currentTeam}
              currentAutomatedPolicies={currentAutomatedPolicies}
            />
          )}
          {!teamId && globalPoliciesError && <TableDataError />}
          {!teamId &&
            !globalPoliciesError &&
            (globalPolicies === undefined ? (
              <Spinner />
            ) : (
              <PoliciesListWrapper
                policiesList={globalPolicies || []}
                isLoading={
                  isLoadingGlobalPolicies && isLoadingFailingPoliciesWebhook
                }
                onRemovePoliciesClick={onRemovePoliciesClick}
                canAddOrRemovePolicy={canAddOrRemovePolicy}
                currentTeam={currentTeam}
                currentAutomatedPolicies={currentAutomatedPolicies}
              />
            ))}
        </div>
        {showInheritedPoliciesButton && (
          <span>
            <Button
              variant="unstyled"
              className={`${showInheritedPolicies ? "upcarat" : "rightcarat"} 
                     ${baseClass}__inherited-policies-button`}
              onClick={toggleShowInheritedPolicies}
            >
              {inheritedPoliciesButtonText(
                showInheritedPolicies,
                globalPolicies.length
              )}
            </Button>
            <div className={`${baseClass}__details`}>
              <IconToolTip
                isHtml
                text={
                  "\
              <center><p>“All teams” policies are checked <br/> for this team’s hosts.</p></center>\
            "
                }
              />
            </div>
          </span>
        )}
        {showInheritedPoliciesButton && showInheritedPolicies && (
          <div className={`${baseClass}__inherited-policies-table`}>
            {globalPoliciesError && <TableDataError />}
            {!globalPoliciesError &&
              (globalPolicies === undefined ? (
                <Spinner />
              ) : (
                <PoliciesListWrapper
                  isLoading={
                    isLoadingGlobalPolicies && isLoadingFailingPoliciesWebhook
                  }
                  policiesList={globalPolicies || []}
                  onRemovePoliciesClick={noop}
                  resultsTitle="policies"
                  canAddOrRemovePolicy={canAddOrRemovePolicy}
                  tableType="inheritedPolicies"
                  currentTeam={currentTeam}
                  currentAutomatedPolicies={currentAutomatedPolicies}
                />
              ))}
          </div>
        )}
        {showManageAutomationsModal && (
          <ManageAutomationsModal
            onCancel={toggleManageAutomationsModal}
            onCreateWebhookSubmit={onCreateWebhookSubmit}
            togglePreviewPayloadModal={togglePreviewPayloadModal}
            showPreviewPayloadModal={showPreviewPayloadModal}
            availablePolicies={globalPolicies || []}
            currentAutomatedPolicies={currentAutomatedPolicies || []}
            currentDestinationUrl={
              (failingPoliciesWebhook &&
                failingPoliciesWebhook.destination_url) ||
              ""
            }
          />
        )}
        {showAddPolicyModal && (
          <AddPolicyModal
            onCancel={toggleAddPolicyModal}
            router={router}
            teamId={teamId}
            teamName={currentTeam?.name}
          />
        )}
        {showRemovePoliciesModal && (
          <RemovePoliciesModal
            onCancel={toggleRemovePoliciesModal}
            onSubmit={onRemovePoliciesSubmit}
          />
        )}
      </div>
    </div>
  );
};

export default ManagePolicyPage;

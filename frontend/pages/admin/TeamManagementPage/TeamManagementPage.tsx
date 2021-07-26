import React, { useState, useCallback } from "react";
import { useSelector, useDispatch } from "react-redux";

import { ITeam } from "interfaces/team";
import teamActions from "redux/nodes/entities/teams/actions";
// ignore TS error for now until these are rewritten in ts.
// @ts-ignore
import { renderFlash } from "redux/nodes/notifications/actions";
import TableContainer from "components/TableContainer";

import CreateTeamModal from "./components/CreateTeamModal";
import DeleteTeamModal from "./components/DeleteTeamModal";
import EditTeamModal from "./components/EditTeamModal";
import EmptyTeams from "./components/EmptyTeams";
import { ICreateTeamFormData } from "./components/CreateTeamModal/CreateTeamModal";
import { IEditTeamFormData } from "./components/EditTeamModal/EditTeamModal";
import { generateTableHeaders, generateDataSet } from "./TeamTableConfig";
import RemoveIcon from "../../../../assets/images/icon-close-vibrant-blue-16x16@2x.png";

const baseClass = "team-management";

// TODO: should probably live close to the store.js file and imported in.
interface RootState {
  entities: {
    teams: {
      isLoading: boolean;
      data: { [id: number]: ITeam };
    };
  };
}

const generateUpdateData = (
  currentTeamData: ITeam,
  formData: IEditTeamFormData
): IEditTeamFormData | null => {
  if (currentTeamData.name !== formData.name) {
    return {
      name: formData.name,
    };
  }
  return null;
};

const TeamManagementPage = (): JSX.Element => {
  const dispatch = useDispatch();
  const [showCreateTeamModal, setShowCreateTeamModal] = useState(false);
  const [showDeleteTeamModal, setShowDeleteTeamModal] = useState(false);
  const [showEditTeamModal, setShowEditTeamModal] = useState(false);
  const [teamEditing, setTeamEditing] = useState<ITeam>();

  const toggleCreateTeamModal = useCallback(() => {
    setShowCreateTeamModal(!showCreateTeamModal);
  }, [showCreateTeamModal, setShowCreateTeamModal]);

  const toggleDeleteTeamModal = useCallback(
    (team?: ITeam) => {
      setShowDeleteTeamModal(!showDeleteTeamModal);
      team ? setTeamEditing(team) : setTeamEditing(undefined);
    },
    [showDeleteTeamModal, setShowDeleteTeamModal, setTeamEditing]
  );

  const toggleEditTeamModal = useCallback(
    (team?: ITeam) => {
      setShowEditTeamModal(!showEditTeamModal);
      team ? setTeamEditing(team) : setTeamEditing(undefined);
    },
    [showEditTeamModal, setShowEditTeamModal, setTeamEditing]
  );

  // NOTE: called once on the initial render of this component.
  const onQueryChange = useCallback(
    (queryData) => {
      const { pageIndex, pageSize, searchQuery } = queryData;
      dispatch(
        teamActions.loadAll({
          page: pageIndex,
          perPage: pageSize,
          globalFilter: searchQuery,
        })
      );
    },
    [dispatch]
  );

  const onCreateSubmit = useCallback(
    (formData: ICreateTeamFormData) => {
      dispatch(teamActions.create(formData))
        .then(() => {
          dispatch(
            renderFlash("success", `Successfully created ${formData.name}.`)
          );
          dispatch(teamActions.loadAll({}));
        })
        .catch(() => {
          dispatch(
            renderFlash("error", "Could not create team. Please try again.")
          );
        });
      toggleCreateTeamModal();
    },
    [dispatch, toggleCreateTeamModal]
  );

  const onDeleteSubmit = useCallback(() => {
    dispatch(teamActions.destroy(teamEditing?.id))
      .then(() => {
        dispatch(
          renderFlash("success", `Successfully deleted ${teamEditing?.name}.`)
        );
        dispatch(teamActions.loadAll({}));
      })
      .catch(() => {
        dispatch(
          renderFlash(
            "error",
            `Could not delete ${teamEditing?.name}. Please try again.`
          )
        );
      });
    toggleDeleteTeamModal();
  }, [dispatch, teamEditing, toggleDeleteTeamModal]);

  const onEditSubmit = useCallback(
    (formData: IEditTeamFormData) => {
      const updatedAttrs = generateUpdateData(teamEditing as ITeam, formData);
      // no updates, so no need for a request.
      if (updatedAttrs === null) {
        toggleEditTeamModal();
        return;
      }
      dispatch(teamActions.update(teamEditing?.id, updatedAttrs))
        .then(() => {
          dispatch(
            renderFlash("success", `Successfully edited ${formData.name}.`)
          );
          dispatch(teamActions.loadAll({}));
        })
        .catch(() => {
          dispatch(
            renderFlash(
              "error",
              `Could not edit ${teamEditing?.name}. Please try again.`
            )
          );
        });
      toggleEditTeamModal();
    },
    [dispatch, teamEditing, toggleEditTeamModal]
  );

  const onSelectActionOneClick = useCallback((selectedTeamIds) => {
    console.log("Clicked Action One Button", selectedTeamIds);
  }, []);

  const onSelectActionTwoClick = useCallback((selectedTeamIds) => {
    console.log("Clicked Action Two Button", selectedTeamIds);
  }, []);
  
  const onActionSelection = (action: string, team: ITeam): void => {
    switch (action) {
      case "edit":
        toggleEditTeamModal(team);
        break;
      case "delete":
        toggleDeleteTeamModal(team);
        break;
      default:
    }
  };

  const tableHeaders = generateTableHeaders(onActionSelection);
  const loadingTableData = useSelector(
    (state: RootState) => state.entities.teams.isLoading
  );
  const teams = useSelector((state: RootState) =>
    generateDataSet(state.entities.teams.data)
  );

    // TODO replace stub actions; consider whether or not this ought to be in separate config
    const secondarySelectActions = [
      {
        name: "action-one",
        buttonText: "Action One",
        onActionButtonClick: onSelectActionOneClick,
        hideButton: (targetIds: number[] = []) => targetIds.length === 1,
        variant: "text-link",
      },
      {
        name: "action-two",
        buttonText: "Action Two",
        onActionButtonClick: onSelectActionTwoClick,
        hideButton: false,
        variant: "text-link",
        iconLink: RemoveIcon,
      },
    ];

  return (
    <div className={`${baseClass} body-wrap`}>
      <p className={`${baseClass}__page-description`}>
        Create, customize, and remove teams from Fleet.
      </p>
      <TableContainer
        columns={tableHeaders}
        data={teams}
        isLoading={loadingTableData}
        defaultSortHeader={"name"}
        defaultSortDirection={"asc"}
        inputPlaceHolder={"Search"}
        actionButtonText={"Create Team"}
        actionButtonVariant={"primary"}
        onActionButtonClick={toggleCreateTeamModal}
        onQueryChange={onQueryChange}
        resultsTitle={"teams"}
        emptyComponent={EmptyTeams}
        showMarkAllPages={false}
        isAllPagesSelected={false}
        searchable
        secondarySelectActions={secondarySelectActions}

      />
      {showCreateTeamModal ? (
        <CreateTeamModal
          onCancel={toggleCreateTeamModal}
          onSubmit={onCreateSubmit}
        />
      ) : null}
      {showDeleteTeamModal ? (
        <DeleteTeamModal
          onCancel={toggleDeleteTeamModal}
          onSubmit={onDeleteSubmit}
          name={teamEditing?.name || ""}
        />
      ) : null}
      {showEditTeamModal ? (
        <EditTeamModal
          onCancel={toggleEditTeamModal}
          onSubmit={onEditSubmit}
          defaultName={teamEditing?.name || ""}
        />
      ) : null}
    </div>
  );
};

export default TeamManagementPage;

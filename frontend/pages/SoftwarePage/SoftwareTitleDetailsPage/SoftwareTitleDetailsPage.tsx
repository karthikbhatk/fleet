import React, { useContext, useMemo } from "react";
import { RouteComponentProps } from "react-router";
import { useQuery } from "react-query";

import { AppContext } from "context/app";
import { ISoftwareTitle } from "interfaces/software";
import softwareAPI, {
  ISoftwareTitleResponse,
} from "services/entities/software";

import MainContent from "components/MainContent";

import SoftwareDetailsSummary from "../components/SoftwareDetailsSummary";
import SoftwareTitleDetailsTable from "./SoftwareTitleDetailsTable";

const baseClass = "SoftwareTitleDetailsPage";

interface ISoftwareTitleDetailsRouteParams {
  id: string;
}

type ISoftwareTitleDetailsPageProps = RouteComponentProps<
  undefined,
  ISoftwareTitleDetailsRouteParams
>;

const SoftwareTitleDetailsPage = ({
  router,
  routeParams,
}: ISoftwareTitleDetailsPageProps) => {
  // TODO: handle non integer values
  const softwareId = parseInt(routeParams.id, 10);
  const { isSandboxMode, filteredSoftwarePath } = useContext(AppContext);

  const {
    data: softwareTitle,
    isLoading: isSoftwareTitleLoading,
    isError: isSoftwareTitleError,
  } = useQuery<ISoftwareTitleResponse, Error, ISoftwareTitle>(
    ["softwareById", softwareId],
    () => softwareAPI.getSoftwareTitle(softwareId),
    {
      select: (data) => data.software_title,
    }
  );

  if (!softwareTitle) {
    return null;
  }

  return (
    <MainContent className={baseClass}>
      <>
        <SoftwareDetailsSummary
          softwareId={softwareId}
          title={softwareTitle.name}
          type={softwareTitle.source}
          hosts={softwareTitle.hosts_count}
        />
        <SoftwareTitleDetailsTable
          router={router}
          data={softwareTitle.versions}
          isLoading={isSoftwareTitleLoading}
        />
      </>
    </MainContent>
  );
  return <h1>script title details</h1>;
};

export default SoftwareTitleDetailsPage;

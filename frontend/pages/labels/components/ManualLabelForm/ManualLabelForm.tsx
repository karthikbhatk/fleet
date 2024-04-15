import React, { useDeferredValue, useState } from "react";
import { useQuery } from "react-query";
import { Row } from "react-table";

import { IHost } from "interfaces/host";
import targetsAPI, { ITargetsSearchResponse } from "services/entities/targets";

import TargetsInput from "components/LiveQuery/TargetsInput";

import LabelForm from "../LabelForm";
import { ILabelFormData } from "../LabelForm/LabelForm";
import { generateTableHeaders } from "./LabelHostTargetTableConfig";

const baseClass = "ManualLabelForm";

const LABEL_TARGET_HOSTS_INPUT_LABEL = "Select hosts";
const LABEL_TARGET_HOSTS_INPUT_PLACEHOLDER =
  "Search name, hostname, or serial number";

export interface IManualLabelFormData {
  name: string;
  description: string;
  targetedHosts: IHost[];
}

interface ITargetsQueryKey {
  scope: string;
  query?: string | null;
  excludedHostIds?: number[];
}

interface IManualLabelFormProps {
  defaultName?: string;
  defaultDescription?: string;
  defaultTargetedHosts?: IHost[];
  onSave: (formData: IManualLabelFormData) => void;
  onCancel: () => void;
}

const ManualLabelForm = ({
  defaultName = "",
  defaultDescription = "",
  defaultTargetedHosts = [],
  onSave,
  onCancel,
}: IManualLabelFormProps) => {
  const [searchQuery, setSearchQuery] = useState<string>("");
  const defferedQuery = useDeferredValue(searchQuery);

  const [targetedHosts, setTargetedHosts] = useState<IHost[]>(
    defaultTargetedHosts
  );

  const {
    data: hostTargets,
    isLoading: isLoadingSearchResults,
    isError: isErrorSearchResults,
  } = useQuery<ITargetsSearchResponse, Error, IHost[], ITargetsQueryKey[]>(
    [
      {
        scope: "labels-targets-search",
        query: defferedQuery,
        excludedHostIds: [], // TODO: add this
      },
    ],
    ({ queryKey }) => {
      const { query, excludedHostIds } = queryKey[0];
      return targetsAPI.search({
        query: query ?? "",
        excluded_host_ids: excludedHostIds ?? null,
      });
    },
    {
      select: (data) => data.hosts,
      enabled: defferedQuery !== "",
    }
  );

  const onHostSelect = (row: Row<IHost>) => {
    setTargetedHosts((prevHosts) => prevHosts.concat(row.original));
    setSearchQuery("");
  };

  const onHostRemove = (row: Row<IHost>) => {
    setTargetedHosts((prevHosts) =>
      prevHosts.filter((h) => h.id !== row.original.id)
    );
  };

  const onSaveNewLabel = (
    labelFormData: ILabelFormData,
    labelFormDataValid: boolean
  ) => {
    if (labelFormDataValid) {
      // values from LabelForm component must be valid too
      onSave({ ...labelFormData, targetedHosts });
    }
  };

  const onChangeSearchQuery = (value: string) => {
    setSearchQuery(value);
  };

  const resultsTableConfig = generateTableHeaders();
  const selectedHostsTableConfig = generateTableHeaders(onHostRemove);

  return (
    <div className={baseClass}>
      <LabelForm
        defaultName={defaultName}
        defaultDescription={defaultDescription}
        onCancel={onCancel}
        onSave={onSaveNewLabel}
        additionalFields={
          <TargetsInput
            label={LABEL_TARGET_HOSTS_INPUT_LABEL}
            placeholder={LABEL_TARGET_HOSTS_INPUT_PLACEHOLDER}
            searchText={searchQuery}
            searchResultsTableConfig={resultsTableConfig}
            selectedHostsTableConifg={selectedHostsTableConfig}
            isTargetsLoading={isLoadingSearchResults}
            hasFetchError={isErrorSearchResults}
            searchResults={hostTargets ?? []}
            targetedHosts={targetedHosts}
            setSearchText={onChangeSearchQuery}
            handleRowSelect={onHostSelect}
          />
        }
      />
    </div>
  );
};

export default ManualLabelForm;

/* eslint-disable  @typescript-eslint/explicit-module-boundary-types */
import sendRequest from "services";
import endpoints from "utilities/endpoints";
import { HostStatus } from "interfaces/host";
import {
  buildQueryStringFromParams,
  getLabelParam,
  reconcileMutuallyExclusiveHostParams,
  reconcileMutuallyInclusiveHostParams,
} from "utilities/url";
import { DiskEncryptionStatus } from "utilities/constants";
import { MacSettingsStatusQueryParam } from "./hosts";

export interface ISortOption {
  key: string;
  direction: string;
}

export interface IHostCountLoadOptions {
  page?: number;
  perPage?: number;
  selectedLabels?: string[];
  globalFilter?: string;
  status?: HostStatus;
  teamId?: number;
  policyId?: number;
  policyResponse?: string;
  macSettingsStatus?: MacSettingsStatusQueryParam;
  softwareId?: number;
  lowDiskSpaceHosts?: number;
  mdmId?: number;
  mdmEnrollmentStatus?: string;
  munkiIssueId?: number;
  osId?: number;
  osName?: string;
  osVersion?: string;
  diskEncryptionStatus?: DiskEncryptionStatus;
}

export default {
  load: (options: IHostCountLoadOptions | undefined) => {
    const selectedLabels = options?.selectedLabels || [];
    const policyId = options?.policyId;
    const policyResponse = options?.policyResponse;
    const globalFilter = options?.globalFilter || "";
    const teamId = options?.teamId;
    const softwareId = options?.softwareId;
    const macSettingsStatus = options?.macSettingsStatus;
    const status = options?.status;
    const mdmId = options?.mdmId;
    const mdmEnrollmentStatus = options?.mdmEnrollmentStatus;
    const munkiIssueId = options?.munkiIssueId;
    const lowDiskSpaceHosts = options?.lowDiskSpaceHosts;
    const label = getLabelParam(selectedLabels);
    const osId = options?.osId;
    const osName = options?.osName;
    const osVersion = options?.osVersion;
    const diskEncryptionStatus = options?.diskEncryptionStatus;

    const queryParams = {
      query: globalFilter,
      ...reconcileMutuallyInclusiveHostParams({ teamId, macSettingsStatus }),
      ...reconcileMutuallyExclusiveHostParams({
        label,
        policyId,
        policyResponse,
        mdmId,
        mdmEnrollmentStatus,
        munkiIssueId,
        softwareId,
        lowDiskSpaceHosts,
        osName,
        osId,
        osVersion,
        diskEncryptionStatus,
      }),
      label_id: label,
      status,
    };

    const queryString = buildQueryStringFromParams(queryParams);
    const endpoint = endpoints.HOSTS_COUNT;
    const path = `${endpoint}?${queryString}`;
    return sendRequest("GET", path);
  },
};

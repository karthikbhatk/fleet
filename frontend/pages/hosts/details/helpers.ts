/** Helpers used across the host details and my device pages and components. */

import { HostDeviceStatus, HostPendingAction } from "interfaces/host";
import {
  IHostMdmProfile,
  IWindowsDiskEncryptionStatus,
  MdmProfileStatus,
} from "interfaces/mdm";

const convertWinDiskEncryptionStatusToProfileStatus = (
  diskEncryptionStatus: IWindowsDiskEncryptionStatus
): MdmProfileStatus => {
  return diskEncryptionStatus === "enforcing"
    ? "pending"
    : diskEncryptionStatus;
};

/**
 * Manually generates a profile for the windows disk encryption status. We need
 * this as we don't have a windows disk encryption profile in the `profiles`
 * attribute coming back from the GET /hosts/:id API response.
 */
// eslint-disable-next-line import/prefer-default-export
export const generateWinDiskEncryptionProfile = (
  diskEncryptionStatus: IWindowsDiskEncryptionStatus,
  detail: string
): IHostMdmProfile => {
  return {
    profile_uuid: "0", // This s the only type of profile that can have this value
    platform: "windows",
    name: "Disk Encryption",
    status: convertWinDiskEncryptionStatusToProfileStatus(diskEncryptionStatus),
    detail,
    operation_type: null,
  };
};

/**
 * Gets the current UI state for the host device status. This helps us know what
 * to display in the UI depending host device status or pending device actions.
 */
export const getHostDeviceStatusUIState = (
  deviceStatus: HostDeviceStatus | null,
  pendingAction: HostPendingAction | null
) => {
  if (deviceStatus === null && pendingAction === null) {
    return null;
  } else if (pendingAction) {
    return pendingAction;
  }
  return deviceStatus;
};

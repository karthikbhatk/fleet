import React, { useState, useCallback } from "react";
import { useDispatch } from "react-redux";
import { Params, InjectedRouter } from "react-router/lib/Router";
import { useQuery } from "react-query";
import { useErrorHandler } from "react-error-boundary";
import { Tab, Tabs, TabList, TabPanel } from "react-tabs";

import classnames from "classnames";
import { isEmpty, pick, reduce } from "lodash";

import PATHS from "router/paths";
import hostAPI from "services/entities/hosts";
import {
  IHost,
  IDeviceMappingResponse,
  IMacadminsResponse,
} from "interfaces/host";
import { ISoftware } from "interfaces/software";
import { ILabel } from "interfaces/label";
// @ts-ignore
import { renderFlash } from "redux/nodes/notifications/actions";
import ReactTooltip from "react-tooltip";
import Spinner from "components/Spinner";
import Button from "components/buttons/Button";
import TableContainer from "components/TableContainer";
import TabsWrapper from "components/TabsWrapper";
import InfoBanner from "components/InfoBanner";
import TooltipWrapper from "components/TooltipWrapper";
import {
  humanHostUptime,
  humanHostLastSeen,
  humanHostEnrolled,
  humanHostMemory,
  humanHostDetailUpdated,
  secondsToHms,
} from "fleet/helpers";

import SoftwareTab from "./SoftwareTab/SoftwareTab";
import InfoModal from "./InfoModal";

import InfoIcon from "../../../assets/images/icon-info-purple-14x14@2x.png";
import IssueIcon from "../../../assets/images/icon-issue-fleet-black-50-16x16@2x.png";

const baseClass = "device-user";

interface IDeviceUserPageProps {
  router: InjectedRouter; // v3
  params: Params;
}

interface IHostResponse {
  host: IHost;
}

const DeviceUserPage = ({
  router,
  params: { host_id },
}: IDeviceUserPageProps): JSX.Element => {
  const hostIdFromURL = parseInt(host_id, 10);
  const dispatch = useDispatch();
  const handlePageError = useErrorHandler();

  const [showInfoModal, setShowInfoModal] = useState<boolean>(false);

  const [refetchStartTime, setRefetchStartTime] = useState<number | null>(null);
  const [showRefetchSpinner, setShowRefetchSpinner] = useState<boolean>(false);
  const [hostSoftware, setHostSoftware] = useState<ISoftware[]>([]);

  const { data: deviceMapping, refetch: refetchDeviceMapping } = useQuery(
    ["deviceMapping", hostIdFromURL],
    () => hostAPI.loadHostDetailsExtension(hostIdFromURL, "device_mapping"),
    {
      enabled: !!hostIdFromURL,
      refetchOnMount: false,
      refetchOnReconnect: false,
      refetchOnWindowFocus: false,
      retry: false,
      select: (data: IDeviceMappingResponse) => data.device_mapping,
    }
  );

  const { data: macadmins, refetch: refetchMacadmins } = useQuery(
    ["macadmins", hostIdFromURL],
    () => hostAPI.loadHostDetailsExtension(hostIdFromURL, "macadmins"),
    {
      enabled: !!hostIdFromURL,
      refetchOnMount: false,
      refetchOnReconnect: false,
      refetchOnWindowFocus: false,
      retry: false,
      select: (data: IMacadminsResponse) => data.macadmins,
    }
  );

  const refetchExtensions = () => {
    deviceMapping !== null && refetchDeviceMapping();
    macadmins !== null && refetchMacadmins();
  };

  const {
    isLoading: isLoadingHost,
    data: host,
    refetch: refetchHostDetails,
  } = useQuery<IHostResponse, Error, IHost>(
    ["host", hostIdFromURL],
    () => hostAPI.loadHostDetails(hostIdFromURL),
    {
      enabled: !!hostIdFromURL,
      refetchOnMount: false,
      refetchOnReconnect: false,
      refetchOnWindowFocus: false,
      retry: false,
      select: (data: IHostResponse) => data.host,
      onSuccess: (returnedHost) => {
        setShowRefetchSpinner(returnedHost.refetch_requested);
        if (returnedHost.refetch_requested) {
          // If the API reports that a Fleet refetch request is pending, we want to check back for fresh
          // host details. Here we set a one second timeout and poll the API again using
          // fullyReloadHost. We will repeat this process with each onSuccess cycle for a total of
          // 60 seconds or until the API reports that the Fleet refetch request has been resolved
          // or that the host has gone offline.
          if (!refetchStartTime) {
            // If our 60 second timer wasn't already started (e.g., if a refetch was pending when
            // the first page loads), we start it now if the host is online. If the host is offline,
            // we skip the refetch on page load.
            if (returnedHost.status === "online") {
              setRefetchStartTime(Date.now());
              setTimeout(() => {
                refetchHostDetails();
                refetchExtensions();
              }, 1000);
            } else {
              setShowRefetchSpinner(false);
            }
          } else {
            const totalElapsedTime = Date.now() - refetchStartTime;
            if (totalElapsedTime < 60000) {
              if (returnedHost.status === "online") {
                setTimeout(() => {
                  refetchHostDetails();
                  refetchExtensions();
                }, 1000);
              } else {
                dispatch(
                  renderFlash(
                    "error",
                    `This host is offline. Please try refetching host vitals later.`
                  )
                );
                setShowRefetchSpinner(false);
              }
            } else {
              dispatch(
                renderFlash(
                  "error",
                  `We're having trouble fetching fresh vitals for this host. Please try again later.`
                )
              );
              setShowRefetchSpinner(false);
            }
          }
          return; // exit early because refectch is pending so we can avoid unecessary steps below
        }
        setHostSoftware(returnedHost.software);
      },
      onError: (error) => handlePageError(error),
    }
  );

  const wrapFleetHelper = (
    helperFn: (value: any) => string,
    value: string
  ): any => {
    return value === "---" ? value : helperFn(value);
  };
  // returns a mixture of props from host
  const normalizeEmptyValues = (hostData: any): { [key: string]: any } => {
    return reduce(
      hostData,
      (result, value, key) => {
        if ((Number.isFinite(value) && value !== 0) || !isEmpty(value)) {
          Object.assign(result, { [key]: value });
        } else {
          Object.assign(result, { [key]: "---" });
        }
        return result;
      },
      {}
    );
  };

  const titleData = normalizeEmptyValues(
    pick(host, [
      "status",
      "issues",
      "memory",
      "cpu_type",
      "os_version",
      "osquery_version",
      "enroll_secret_name",
      "detail_updated_at",
      "percent_disk_space_available",
      "gigs_disk_space_available",
    ])
  );

  const operatingSystem = host?.os_version.slice(
    0,
    host?.os_version.lastIndexOf(" ")
  );
  const operatingSystemVersion = host?.os_version.slice(
    host?.os_version.lastIndexOf(" ") + 1
  );

  const aboutData = normalizeEmptyValues(
    pick(host, [
      "seen_time",
      "uptime",
      "last_enrolled_at",
      "hardware_model",
      "hardware_serial",
      "primary_ip",
    ])
  );

  const osqueryData = normalizeEmptyValues(
    pick(host, [
      "config_tls_refresh",
      "logger_tls_period",
      "distributed_interval",
    ])
  );

  const toggleInfoModal = useCallback(() => {
    setShowInfoModal(!showInfoModal);
  }, [showInfoModal, setShowInfoModal]);

  const onRefetchHost = async () => {
    if (host) {
      // Once the user clicks to refetch, the refetch loading spinner should continue spinning
      // unless there is an error. The spinner state is also controlled in the fullyReloadHost
      // method.
      setShowRefetchSpinner(true);
      try {
        await hostAPI.refetch(host).then(() => {
          setRefetchStartTime(Date.now());
          setTimeout(() => {
            refetchHostDetails();
            refetchExtensions();
          }, 1000);
        });
      } catch (error) {
        console.log(error);
        dispatch(renderFlash("error", `Host "${host.hostname}" refetch error`));
        setShowRefetchSpinner(false);
      }
    }
  };

  const onLabelClick = (label: ILabel) => {
    if (label.name === "All Hosts") {
      return router.push(PATHS.MANAGE_HOSTS);
    }

    return router.push(`${PATHS.MANAGE_HOSTS}/labels/${label.id}`);
  };

  const renderActionButtons = () => {
    return (
      <div className={`${baseClass}__action-button-container`}>
        <Button onClick={() => setShowInfoModal(true)} variant="text-icon">
          <>
            Info <img src={InfoIcon} alt="Host info icon" />
          </>
        </Button>
      </div>
    );
  };

  const renderLabels = () => {
    const { labels = [] } = host || {};

    const labelItems = labels.map((label) => {
      return (
        <li className="list__item" key={label.id}>
          <Button
            onClick={() => onLabelClick(label)}
            variant="label"
            className="list__button"
          >
            {label.name}
          </Button>
        </li>
      );
    });

    return (
      <div className="section labels col-50">
        <p className="section__header">Labels</p>
        {labels.length === 0 ? (
          <p className="info-flex__item">
            No labels are associated with this host.
          </p>
        ) : (
          <ul className="list">{labelItems}</ul>
        )}
      </div>
    );
  };

  const renderSoftware = () => {
    return <SoftwareTab isLoading={isLoadingHost} software={hostSoftware} />;
  };

  const renderRefetch = () => {
    const isOnline = host?.status === "online";

    return (
      <>
        <div
          className="refetch"
          data-tip
          data-for="refetch-tooltip"
          data-tip-disable={isOnline || showRefetchSpinner}
        >
          <Button
            className={`
              button
              button--unstyled
              ${!isOnline ? "refetch-offline" : ""} 
              ${showRefetchSpinner ? "refetch-spinner" : "refetch-btn"}
            `}
            disabled={!isOnline}
            onClick={onRefetchHost}
          >
            {showRefetchSpinner
              ? "Fetching fresh vitals...this may take a moment"
              : "Refetch"}
          </Button>
        </div>
        <ReactTooltip
          place="bottom"
          type="dark"
          effect="solid"
          id="refetch-tooltip"
          backgroundColor="#3e4771"
        >
          <span className={`${baseClass}__tooltip-text`}>
            You can’t fetch data from <br /> an offline host.
          </span>
        </ReactTooltip>
      </>
    );
  };

  const renderIssues = () => (
    <div className="info-flex__item info-flex__item--title">
      <span className="info-flex__header">Issues</span>
      <span className="info-flex__data">
        <span
          className="host-issue tooltip__tooltip-icon"
          data-tip
          data-for="host-issue-count"
          data-tip-disable={false}
        >
          <img alt="host issue" src={IssueIcon} />
        </span>
        <ReactTooltip
          place="bottom"
          type="dark"
          effect="solid"
          backgroundColor="#3e4771"
          id="host-issue-count"
          data-html
        >
          <span className={`tooltip__tooltip-text`}>
            Failing policies ({host?.issues.failing_policies_count})
          </span>
        </ReactTooltip>
        <span className={`total-issues-count`}>
          {host?.issues.total_issues_count}
        </span>
      </span>
    </div>
  );

  const renderHostTeam = () => (
    <div className="info-flex__item info-flex__item--title">
      <span className="info-flex__header">Team</span>
      <span className={`info-flex__data`}>
        {host?.team_name ? (
          `${host?.team_name}`
        ) : (
          <span className="info-flex__no-team">No team</span>
        )}
      </span>
    </div>
  );

  const renderDeviceUser = () => {
    const numUsers = deviceMapping?.length;
    if (numUsers) {
      return (
        <div className="info-grid__block">
          <span className="info-grid__header">Used by</span>
          <span className="info-grid__data">
            {numUsers === 1 && deviceMapping ? (
              deviceMapping[0].email || "---"
            ) : (
              <span className={`${baseClass}__device-mapping`}>
                <span
                  className="device-user"
                  data-tip
                  data-for="device-user-tooltip"
                >
                  {`${numUsers} users`}
                </span>
                <ReactTooltip
                  place="top"
                  type="dark"
                  effect="solid"
                  id="device-user-tooltip"
                  backgroundColor="#3e4771"
                >
                  <div
                    className={`${baseClass}__tooltip-text device-user-tooltip`}
                  >
                    {deviceMapping?.map((user, i, arr) => (
                      <span key={user.email}>{`${user.email}${
                        i < arr.length - 1 ? ", " : ""
                      }`}</span>
                    ))}
                  </div>
                </ReactTooltip>
              </span>
            )}
          </span>
        </div>
      );
    }
    return null;
  };

  const renderDiskSpace = () => {
    if (
      host &&
      (host.gigs_disk_space_available > 0 ||
        host.percent_disk_space_available > 0)
    ) {
      return (
        <span className="info-flex__data">
          <div className="info-flex__disk-space">
            <div
              className={
                titleData.percent_disk_space_available > 20
                  ? "info-flex__disk-space-used"
                  : "info-flex__disk-space-warning"
              }
              style={{
                width: `${100 - titleData.percent_disk_space_available}%`,
              }}
            />
          </div>
          {titleData.gigs_disk_space_available} GB available
        </span>
      );
    }
    return <span className="info-flex__data">No data available</span>;
  };

  const renderShowInfoModal = () => <InfoModal onCancel={toggleInfoModal} />;

  if (isLoadingHost) {
    return <Spinner />;
  }

  const statusClassName = classnames("status", `status--${host?.status}`);
  return (
    <div className={`${baseClass} body-wrap`}>
      <div className="header title">
        <div className="title__inner">
          <div className="hostname-container">
            <h1 className="hostname">My device</h1>
            <p className="last-fetched">
              {`Last reported vitals ${humanHostDetailUpdated(
                titleData.detail_updated_at
              )}`}
              &nbsp;
            </p>
            {renderRefetch()}
          </div>
        </div>
        {renderActionButtons()}
      </div>
      <div className="section title">
        <div className="title__inner">
          <div className="info-flex">
            <div className="info-flex__item info-flex__item--title">
              <span className="info-flex__header">Status</span>
              <span className={`${statusClassName} info-flex__data`}>
                {titleData.status}
              </span>
            </div>
            {titleData.issues?.total_issues_count > 0 && renderIssues()}
            <div className="info-flex__item info-flex__item--title">
              <span className="info-flex__header">Disk Space</span>
              {renderDiskSpace()}
            </div>
            <div className="info-flex__item info-flex__item--title">
              <span className="info-flex__header">Memory</span>
              <span className="info-flex__data">
                {wrapFleetHelper(humanHostMemory, titleData.memory)}
              </span>
            </div>
            <div className="info-flex__item info-flex__item--title">
              <span className="info-flex__header">Processor type</span>
              <span className="info-flex__data">{titleData.cpu_type}</span>
            </div>
            <div className="info-flex__item info-flex__item--title">
              <span className="info-flex__header">Operating system</span>
              <span className="info-flex__data">{titleData.os_version}</span>
            </div>
          </div>
        </div>
      </div>
      <TabsWrapper>
        <Tabs>
          <TabList>
            <Tab>Details</Tab>
            <Tab>Software</Tab>
          </TabList>
          <TabPanel>
            <div className="section about">
              <p className="section__header">About</p>
              <div className="info-grid">
                <div className="info-grid__block">
                  <span className="info-grid__header">Last restarted</span>
                  <span className="info-grid__data">
                    {wrapFleetHelper(humanHostUptime, aboutData.uptime)}
                  </span>
                </div>
                <div className="info-grid__block">
                  <span className="info-grid__header">Hardware model</span>
                  <span className="info-grid__data">
                    {aboutData.hardware_model}
                  </span>
                </div>
                <div className="info-grid__block">
                  <span className="info-grid__header">Added to Fleet</span>
                  <span className="info-grid__data">
                    {wrapFleetHelper(
                      humanHostEnrolled,
                      aboutData.last_enrolled_at
                    )}
                  </span>
                </div>
                <div className="info-grid__block">
                  <span className="info-grid__header">Serial number</span>
                  <span className="info-grid__data">
                    {aboutData.hardware_serial}
                  </span>
                </div>
                <div className="info-grid__block">
                  <span className="info-grid__header">IP address</span>
                  <span className="info-grid__data">
                    {aboutData.primary_ip}
                  </span>
                </div>
                {renderDeviceUser()}
              </div>
            </div>
          </TabPanel>
          <TabPanel>{renderSoftware()}</TabPanel>
        </Tabs>
      </TabsWrapper>

      {showInfoModal && renderShowInfoModal()}
    </div>
  );
};

export default DeviceUserPage;

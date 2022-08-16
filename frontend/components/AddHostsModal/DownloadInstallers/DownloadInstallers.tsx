import React, { FunctionComponent, useState } from "react";

import {
  IInstallerPlatform,
  IInstallerType,
  INSTALLER_PLATFORM_BY_TYPE,
  INSTALLER_TYPE_BY_PLATFORM,
} from "interfaces/installer";
import ENDPOINTS from "utilities/endpoints";
import local from "utilities/local";
import URL_PREFIX from "router/url_prefix";
import installerAPI from "services/entities/installers";

import Button from "components/buttons/Button";
import Checkbox from "components/forms/fields/Checkbox";
import DataError from "components/DataError";
import Spinner from "components/Spinner";
import TooltipWrapper from "components/TooltipWrapper";

import AppleIcon from "./../../../../assets/images/icon-apple-black-24x24@2x.png";
import AppleIconVibrant from "./../../../../assets/images/icon-apple-vibrant-blue-24x24@2x.png";
import LinuxIcon from "./../../../../assets/images/icon-linux-black-24x24@2x.png";
import LinuxIconVibrant from "./../../../../assets/images/icon-linux-vibrant-blue-24x24@2x.png";
import WindowsIcon from "./../../../../assets/images/icon-windows-black-24x24@2x.png";
import WindowsIconVibrant from "./../../../../assets/images/icon-windows-vibrant-blue-24x24@2x.png";
import SuccessIcon from "./../../../../assets/images/icon-circle-check-blue-48x48@2x.png";

interface IDownloadInstallersProps {
  enrollSecret: string;
  onCancel: () => void;
}

interface IDownloadFormProps {
  url: string;
  onSubmit: (event: React.FormEvent<HTMLFormElement>) => void;
  token: string | null;
  enrollSecret: string;
  includeDesktop: boolean;
  selectedInstaller: string | undefined;
  isDownloading: boolean;
  isDownloadSuccess: boolean;
}

const baseClass = "download-installers";

const displayOrder = [
  "macOS",
  "Windows",
  "Linux (RPM)",
  "Linux (deb)",
] as const;

const displayIcon = (platform: IInstallerPlatform, isSelected: boolean) => {
  switch (platform) {
    case "Linux (RPM)":
    case "Linux (deb)":
      return (
        <img src={isSelected ? LinuxIconVibrant : LinuxIcon} alt={platform} />
      );
    case "macOS":
      return (
        <img src={isSelected ? AppleIconVibrant : AppleIcon} alt={platform} />
      );
    case "Windows":
      return (
        <img
          src={isSelected ? WindowsIconVibrant : WindowsIcon}
          alt={platform}
        />
      );
    default:
      return null;
  }
};

const DownloadForm: FunctionComponent<IDownloadFormProps> = ({
  url,
  onSubmit,
  token,
  enrollSecret,
  includeDesktop,
  selectedInstaller,
  isDownloading,
  isDownloadSuccess,
}) => {
  return (
    <form
      key="form"
      method="POST"
      action={url}
      target="_self"
      onSubmit={onSubmit}
    >
      <input type="hidden" name="token" value={token || ""} />
      <input type="hidden" name="enroll_secret" value={enrollSecret} />
      <input type="hidden" name="desktop" value={String(includeDesktop)} />
      {!isDownloadSuccess && (
        <Button
          className={`${baseClass}__button--download`}
          disabled={!selectedInstaller}
          type="submit"
        >
          {isDownloading ? <Spinner /> : "Download installer"}
        </Button>
      )}
    </form>
  );
};

const DownloadInstallers = ({
  enrollSecret,
  onCancel,
}: IDownloadInstallersProps): JSX.Element => {
  const [includeDesktop, setIncludeDesktop] = useState(true);
  const [isDownloading, setIsDownloading] = useState(false);
  const [isDownloadError, setIsDownloadError] = useState(false);
  const [isDownloadSuccess, setIsDownloadSuccess] = useState(false);
  const [selectedInstaller, setSelectedInstaller] = useState<
    IInstallerType | undefined
  >();
  const path = `${ENDPOINTS.DOWNLOAD_INSTALLER}/${selectedInstaller}`;
  const { origin } = global.window.location;
  const url = `${origin}${URL_PREFIX}/api${path}`;
  const token = local.getItem("auth_token");

  const downloadInstaller = async (event: React.FormEvent<HTMLFormElement>) => {
    if (!selectedInstaller) {
      // do nothing
      return;
    }

    // Prevent the submit behavior, as we want to control when the POST is
    // actually performed.
    event.preventDefault();
    event.persist();

    setIsDownloading(true);
    try {
      // First check if the installer exists, no need to save the result of
      // this operation as any status other than 200 will throw an error
      await installerAPI.checkInstallerExistence({
        enrollSecret,
        includeDesktop,
        installerType: selectedInstaller,
      });

      (event.target as HTMLFormElement).submit();
      setIsDownloadSuccess(true);
    } catch (error) {
      setIsDownloadError(true);
    } finally {
      setIsDownloading(false);
    }
  };

  const onClickSelector = (type: IInstallerType) => {
    if (isDownloading) {
      // do nothing
      return;
    }
    if (type === selectedInstaller) {
      setSelectedInstaller(undefined);
      return;
    }
    setSelectedInstaller(type);
  };

  const form = (
    <DownloadForm
      key="downloadForm"
      url={url}
      onSubmit={downloadInstaller}
      token={token}
      enrollSecret={enrollSecret}
      includeDesktop={includeDesktop}
      selectedInstaller={selectedInstaller}
      isDownloading={isDownloading}
      isDownloadSuccess={isDownloadSuccess}
    />
  );

  if (isDownloadError) {
    return (
      <div className={`${baseClass}__error`}>
        <DataError />
      </div>
    );
  }

  if (isDownloadSuccess) {
    const installerPlatform =
      (selectedInstaller &&
        `${INSTALLER_PLATFORM_BY_TYPE[selectedInstaller]} `) ||
      "";
    return (
      <div className={`${baseClass}__success`}>
        <img src={SuccessIcon} alt="download successful" />
        <h2>You&rsquo;re almost there</h2>
        <p>{`Run the installer on a ${installerPlatform}laptop, workstation, or sever to add it to Fleet.`}</p>
        <Button onClick={onCancel}>Got it</Button>
        {form}
      </div>
    );
  }

  return (
    <div className={`${baseClass}`}>
      <p>Which platform is your host running?</p>
      <div className={`${baseClass}__select-installer`}>
        {displayOrder.map((platform) => {
          const installerType = INSTALLER_TYPE_BY_PLATFORM[platform];
          const isSelected = selectedInstaller === installerType;
          return (
            <div
              key={installerType}
              className={`${baseClass}__selector ${
                isSelected ? `${baseClass}__selector--selected` : ""
              }`}
              onClick={() => onClickSelector(installerType)}
            >
              <span>
                {displayIcon(platform, isSelected)}
                {platform}
              </span>
            </div>
          );
        })}
      </div>
      <Checkbox
        name="include-fleet-desktop"
        onChange={(value: boolean) => setIncludeDesktop(value)}
        value={includeDesktop}
      >
        <>
          Include&nbsp;
          <TooltipWrapper
            tipContent={
              "<p>Lightweight application that allows end users to see information about their device.</p>"
            }
          >
            Fleet Desktop
          </TooltipWrapper>
        </>
      </Checkbox>
      {form}
    </div>
  );
};

export default DownloadInstallers;

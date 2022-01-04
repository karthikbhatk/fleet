import React, { useState, useCallback } from "react";
import { syntaxHighlight } from "fleet/helpers";
import { size } from "lodash";

import yaml from "js-yaml";
// @ts-ignore
import validate from "components/forms/admin/AppConfigForm/validate";
// @ts-ignore
import constructErrorString from "utilities/yaml";

import { IConfigNested } from "interfaces/config";
import { IFormField } from "interfaces/form_field";
import { IEnrollSecret } from "interfaces/enroll_secret";

import Button from "components/buttons/Button";
import Checkbox from "components/forms/fields/Checkbox";
// @ts-ignore
import Dropdown from "components/forms/fields/Dropdown";
import EnrollSecretTable from "components/EnrollSecretTable";
// @ts-ignore
import InputField from "components/forms/fields/InputField";
// @ts-ignore
import OrgLogoIcon from "components/icons/OrgLogoIcon";
// @ts-ignore
import validate from "components/forms/admin/AppConfigForm/validate";
import IconToolTip from "components/IconToolTip";
import InfoBanner from "components/InfoBanner/InfoBanner";
// @ts-ignore
import YamlAce from "components/YamlAce";
import Modal from "components/Modal";
import OpenNewTabIcon from "../../../../../assets/images/open-new-tab-12x12@2x.png";

const authMethodOptions = [
  { label: "Plain", value: "authmethod_plain" },
  { label: "Cram MD5", value: "authmethod_cram_md5" },
  { label: "Login", value: "authmethod_login" },
];
const authTypeOptions = [
  { label: "Username and Password", value: "authtype_username_password" },
  { label: "None", value: "authtype_none" },
];
const percentageOfHosts = [
  { label: "1%", value: 1 },
  { label: "5%", value: 5 },
  { label: "10%", value: 10 },
  { label: "25%", value: 25 },
];
const numberOfDays = [
  { label: "1 day", value: 1 },
  { label: "3 days", value: 3 },
  { label: "7 days", value: 7 },
  { label: "14 days", value: 14 },
];

// consider breaking this up into components
const baseClass = "app-config-form";
// these will all be states...
const formFields = [
  // "authentication_method",
  // "authentication_type",
  // "domain",
  // "enable_ssl_tls",
  // "enable_start_tls",
  // "server_url",
  // "org_logo_url",
  // "org_name",
  // "osquery_enroll_secret", // not modified in UI
  // "password",
  // "port",
  // "sender_address",
  // "server",
  // "user_name",
  // "verify_ssl_certs",
  // "idp_name",
  // "entity_id",
  // "issuer_uri",
  // "idp_image_url",
  // "metadata",
  // "metadata_url",
  // "enable_sso",
  // "enable_sso_idp_login",
  // "enable_smtp",
  // "host_expiry_enabled",
  // "host_expiry_window",
  // "live_query_disabled",
  "agent_options",
  // "enable_host_status_webhook",
  // "destination_url",
  // "host_percentage",
  // "days_count",
  // "enable_analytics",
];

interface IAppConfigFormProps {
  // formData: {
  //   authentication_method: IFormField;
  //   authentication_type: IFormField;
  //   domain: IFormField;
  //   enable_ssl_tls: IFormField;
  //   enable_start_tls: IFormField;
  //   server_url: IFormField;
  //   org_logo_url: string;
  //   org_name: string;
  //   password: IFormField;
  //   port: IFormField;
  //   sender_address: IFormField;
  //   server: IFormField;
  //   user_name: IFormField;
  //   verify_ssl_certs: IFormField;
  //   entity_id: IFormField;
  //   issuer_uri: IFormField;
  //   idp_image_url: IFormField;
  //   metadata: IFormField;
  //   metadata_url: IFormField;
  //   idp_name: IFormField;
  //   enable_sso: IFormField;
  //   enable_sso_idp_login: IFormField;
  //   enable_smtp: IFormField;
  //   host_expiry_enabled: IFormField;
  //   host_expiry_window: IFormField;
  //   live_query_disabled: IFormField;
  //   agent_options: IFormField;
  //   enable_host_status_webhook: IFormField;
  //   destination_url?: IFormField;
  //   host_percentage?: IFormField;
  //   days_count?: IFormField;
  //   enable_analytics: IFormField;
  // };
  formData: IConfigNested;
  enrollSecret: IEnrollSecret[] | undefined;
  handleSubmit: any;
  smtpConfigured: boolean;
}

export interface IFormData {
  // email: string;
  // name: string;
  // newUserType?: NewUserType | null;
  // password?: string | null;
  // sso_enabled?: boolean;
  // global_role: string | null;
  // teams: ITeam[];
  // currentUserId?: number;
  // invited_by?: number;
}

interface IAppConfigFormState {
  showHostStatusWebhookPreviewModal: boolean;
  showUsageStatsPreviewModal: boolean;
}

interface IAppConfigFormErrors {
  // email: string | null;
  // name: string | null;
  // password: string | null;
  // sso_enabled: boolean | null;
}

const AppConfigFormFunctional = ({
  formData,
  enrollSecret,
  handleSubmit,
  smtpConfigured,
}: IAppConfigFormProps): JSX.Element => {
  // █▀ ▀█▀ ▄▀█ ▀█▀ █▀▀
  // ▄█ ░█░ █▀█ ░█░ ██▄
  const [
    showHostStatusWebhookPreviewModal,
    setShowHostStatusWebhookPreviewModal,
  ] = useState<boolean>(false);
  const [
    showUsageStatsPreviewModal,
    setShowUsageStatsPreviewModal,
  ] = useState<boolean>(false);

  // █▀▀ █▀█ █▀█ █▀▄▀█   █▀ ▀█▀ ▄▀█ ▀█▀ █▀▀
  // █▀░ █▄█ █▀▄ █░▀░█   ▄█ ░█░ █▀█ ░█░ ██▄
  // Organization info
  const [orgName, setOrgName] = useState<string>(
    formData.org_info.org_name || ""
  );
  const [orgLogoUrl, setOrgLogoUrl] = useState<string>(
    formData.org_info.org_logo_url || ""
  );
  // Fleet web address
  const [serverURL, setServerURL] = useState<string>(
    formData.server_settings.server_url || ""
  );
  // SAML single sign on options
  const [enableSSO, setEnableSSO] = useState<boolean>(
    formData.sso_settings.enable_sso || false
  );
  const [idpName, setIDPName] = useState<string>(
    formData.sso_settings.idp_name || ""
  );
  const [entityID, setEntityID] = useState<string>(
    formData.sso_settings.entity_id || ""
  );
  const [issuerURI, setIssuerURI] = useState<string>(
    formData.sso_settings.issuer_uri || ""
  );
  const [idpImageURL, setIDPImageURL] = useState<string>(
    formData.sso_settings.idp_image_url || ""
  );
  const [metadata, setMetadata] = useState<string>(
    formData.sso_settings.metadata || ""
  );
  const [metadataURL, setMetadataURL] = useState<string>(
    formData.sso_settings.metadata_url || ""
  );
  const [enableSSOIDLLogin, setEnableSSOIDPLogin] = useState<boolean>(
    formData.sso_settings.enable_sso_idp_login || false
  );
  // SMTP options
  const [enableSMTP, setEnableSMTP] = useState<boolean>(
    formData.smtp_settings.enable_smtp || false
  );
  const [smtpSenderAddress, setSMTPSenderAddress] = useState<string>(
    formData.smtp_settings.sender_address || ""
  );
  const [smtpServer, setSMTPServer] = useState<string>(
    formData.smtp_settings.server || ""
  );
  const [smtpPort, setSMTPPort] = useState<number | undefined>(
    formData.smtp_settings.port || undefined
  );
  const [smtpEnableSSLTLS, setSMTPEnableSSLTLS] = useState<boolean>(
    formData.smtp_settings.enable_ssl_tls || false
  );
  const [smtpAuthenticationType, setSMTPAuthenticationType] = useState<string>(
    formData.smtp_settings.authentication_type || ""
  );
  const [smtpUsername, setSMTPUsername] = useState<string>(
    formData.smtp_settings.user_name || ""
  );
  const [smtpPassword, setSMTPPassword] = useState<string>(
    formData.smtp_settings.password || ""
  );
  const [
    smtpAuthenticationMethod,
    setSMTPAuthenticationMethod,
  ] = useState<string>(formData.smtp_settings.authentication_method || "");
  // Global agent options
  const [agentOptions, setAgentOptions] = useState<any>(
    yaml.dump(formData.agent_options) || {}
  );
  // Host status webhook
  const [
    enableHostStatusWebhook,
    setEnableHostStatusWebhook,
  ] = useState<boolean>(
    formData.webhook_settings.host_status_webhook.enable_host_status_webhook ||
      false
  );
  const [
    hostStatusWebhookDestinationURL,
    setHostStatusWebhookDestinationURL,
  ] = useState<string>(
    formData.webhook_settings.host_status_webhook.destination_url || ""
  );
  const [
    hostStatusWebhookHostPercentage,
    setHostStatusWebhookHostPercentage,
  ] = useState<number | undefined>(
    formData.webhook_settings.host_status_webhook.host_percentage || undefined
  );
  const [hostStatusWebhookDaysCount, setHostStatusWebhookDaysCount] = useState<
    number | undefined
  >(formData.webhook_settings.host_status_webhook.days_count || undefined);
  // Usage statistics
  const [enableUsageStatistics, setEnableUsageStatistics] = useState<boolean>(
    formData.server_settings.enable_analytics || false
  );
  // Advanced options
  const [domain, setDomain] = useState<string>("");
  const [verifySSLCerts, setVerifySSLCerts] = useState<boolean>(
    formData.smtp_settings.verify_ssl_certs || false
  );
  const [enableStartTLS, setEnableStartTLS] = useState<boolean>(
    formData.smtp_settings.enable_start_tls || false
  );
  const [enableHostExpiry, setEnableHostExpiry] = useState<boolean>(
    formData.host_expiry_settings.host_expiry_enabled || false
  );
  const [hostExpiryWindow, setHostExpiryWindow] = useState<number | undefined>(
    formData.host_expiry_settings.host_expiry_window || undefined
  );
  const [disableLiveQuery, setDisableLiveQuery] = useState<boolean>(
    formData.server_settings.live_query_disabled || false
  );

  // █▀▀ █▀█ █▀█ █▀▄▀█   █▀▀ █░█ ▄▀█ █▄░█ █▀▀ █▀▀
  // █▀░ █▄█ █▀▄ █░▀░█   █▄▄ █▀█ █▀█ █░▀█ █▄█ ██▄
  // Organization info
  const onChangeOrgName = useCallback(
    (value: string) => {
      setOrgName(value);
    },
    [setOrgName]
  );
  const onChangeOrgLogoUrl = useCallback(
    (value: string) => {
      setOrgLogoUrl(value);
    },
    [setOrgLogoUrl]
  );
  // Fleet web address
  const onChangeServerURL = useCallback(
    (value: string) => {
      setServerURL(value);
    },
    [setServerURL]
  );
  // SAML single sign on options
  const onChangeEnableSSO = useCallback(
    (value: boolean) => {
      setEnableSSO(value);
    },
    [setEnableSSO]
  );
  const onChangeIDPName = useCallback(
    (value: string) => {
      setIDPName(value);
    },
    [setIDPName]
  );
  const onChangeEntityID = useCallback(
    (value: string) => {
      setEntityID(value);
    },
    [setEntityID]
  );
  const onChangeIssuerURI = useCallback(
    (value: string) => {
      setIssuerURI(value);
    },
    [setIssuerURI]
  );
  const onChangeIDPImageURL = useCallback(
    (value: string) => {
      setIDPImageURL(value);
    },
    [setIDPImageURL]
  );
  const onChangeMetadata = useCallback(
    (value: string) => {
      setMetadata(value);
    },
    [setMetadata]
  );
  const onChangeMetadataURL = useCallback(
    (value: string) => {
      setMetadataURL(value);
    },
    [setMetadataURL]
  );
  const onChangeEnableSSOIDPLogin = useCallback(
    (value: boolean) => {
      setEnableSSOIDPLogin(value);
    },
    [setEnableSSOIDPLogin]
  );
  // SMTP options
  const onChangeEnableSMTP = useCallback(
    (value: boolean) => {
      setEnableSMTP(value);
    },
    [setEnableSMTP]
  );
  const onChangeSMTPSenderAddress = useCallback(
    (value: string) => {
      setSMTPSenderAddress(value);
    },
    [setSMTPSenderAddress]
  );
  const onChangeSMTPServer = useCallback(
    (value: string) => {
      setSMTPServer(value);
    },
    [setSMTPServer]
  );
  const onChangeSMTPPort = useCallback(
    (value: number) => {
      setSMTPPort(value);
    },
    [setSMTPPort]
  );
  const onChangeSMTPEnableSSLTLS = useCallback(
    (value: boolean) => {
      setSMTPEnableSSLTLS(value);
    },
    [setSMTPEnableSSLTLS]
  );
  const onChangeSMTPAuthenticationType = useCallback(
    (value: string) => {
      setSMTPAuthenticationType(value);
    },
    [setSMTPAuthenticationType]
  );
  const onChangeSMTPUsername = useCallback(
    (value: string) => {
      setSMTPUsername(value);
    },
    [setSMTPUsername]
  );
  const onChangeSMTPPassword = useCallback(
    (value: string) => {
      setSMTPPassword(value);
    },
    [setSMTPPassword]
  );
  const onChangeSMTPAuthenticationMethod = useCallback(
    (value: string) => {
      setSMTPAuthenticationMethod(value);
    },
    [setSMTPAuthenticationMethod]
  );
  // Global agent options
  // const onChangeAgentOptions = useCallback(
  //   (value: string) => {
  //     setAgentOptions(value);
  //   },
  //   [setAgentOptions]
  // );

  // Host status webhook
  const onChangeEnableHostStatusWebhook = useCallback(
    (value: boolean) => {
      setEnableHostStatusWebhook(value);
    },
    [setEnableHostStatusWebhook]
  );
  const onChangeHostStatusWebhookDestinationURL = useCallback(
    (value: string) => {
      setHostStatusWebhookDestinationURL(value);
    },
    [setHostStatusWebhookDestinationURL]
  );
  const onChangeHostStatusWebhookHostPercentage = useCallback(
    (value: number) => {
      setHostStatusWebhookHostPercentage(value);
    },
    [setHostStatusWebhookHostPercentage]
  );
  const onChangeHostStatusWebhookDaysCount = useCallback(
    (value: number) => {
      setHostStatusWebhookDaysCount(value);
    },
    [setHostStatusWebhookDaysCount]
  );
  // Usage statistics
  const onChangeEnableUsageStatistics = useCallback(
    (value: boolean) => {
      setEnableUsageStatistics(value);
    },
    [setEnableUsageStatistics]
  );
  // Advanced options
  const onChangeDomain = useCallback(
    (value: string) => {
      setDomain(value);
    },
    [setDomain]
  );
  const onChangeVerifySSLCerts = useCallback(
    (value: boolean) => {
      setVerifySSLCerts(value);
    },
    [setVerifySSLCerts]
  );
  const onChangeEnableStartTLS = useCallback(
    (value: boolean) => {
      setEnableStartTLS(value);
    },
    [setEnableStartTLS]
  );
  const onChangeEnableHostExpiry = useCallback(
    (value: boolean) => {
      setEnableHostExpiry(value);
    },
    [setEnableHostExpiry]
  );
  const onChangeHostExpiryWindow = useCallback(
    (value: number) => {
      setHostExpiryWindow(value);
    },
    [setHostExpiryWindow]
  );
  const onChangeDisableLiveQuery = useCallback(
    (value: boolean) => {
      setDisableLiveQuery(value);
    },
    [setDisableLiveQuery]
  );

  // ▀█▀ █▀█ █▀▀ █▀▀ █░░ █▀▀   █▀▄▀█ █▀█ █▀▄ ▄▀█ █░░ █▀
  // ░█░ █▄█ █▄█ █▄█ █▄▄ ██▄   █░▀░█ █▄█ █▄▀ █▀█ █▄▄ ▄█

  const toggleHostStatusWebhookPreviewModal = () => {
    setShowHostStatusWebhookPreviewModal(!showHostStatusWebhookPreviewModal);
    return false;
  };

  const toggleUsageStatsPreviewModal = () => {
    setShowUsageStatsPreviewModal(!showUsageStatsPreviewModal);
    return false;
  };

  // █▀▀ █▀█ █▀█ █▀▄▀█   █▀ █░█ █▄▄ █▀▄▀█ █ ▀█▀
  // █▀░ █▄█ █▀▄ █░▀░█   ▄█ █▄█ █▄█ █░▀░█ █ ░█░
  const onFormSubmit = () => {
    // Validator

    handleSubmit({
      name: packName,
      description: packDescription,
      targets: [...packFormTargets],
    });
  };

  // █▀ █▀▀ █▀▀ ▀█▀ █ █▀█ █▄░█ █▀
  // ▄█ ██▄ █▄▄ ░█░ █ █▄█ █░▀█ ▄█
  const renderOrganizationInfoSection = () => {
    return (
      <div className={`${baseClass}__section`}>
        <h2>
          <a id="organization-info">Organization info</a>
        </h2>
        <div className={`${baseClass}__inputs`}>
          <InputField
            label="Organization name"
            onChange={onChangeOrgName}
            value={orgName}
          />
          <InputField
            label="Organization avatar URL"
            onChange={onChangeOrgLogoUrl}
            value={orgLogoUrl}
          />
        </div>
        <div className={`${baseClass}__details ${baseClass}__avatar-preview`}>
          <OrgLogoIcon src={orgLogoUrl} />
        </div>
      </div>
    );
  };

  const renderFleetWebAddressSection = () => {
    return (
      <div className={`${baseClass}__section`}>
        <h2>
          <a id="fleet-web-address">Fleet web address</a>
        </h2>
        <div className={`${baseClass}__inputs`}>
          <InputField
            label="Fleet app URL"
            hint={
              <span>
                Include base path only (eg. no <code>/v1</code>)
              </span>
            }
            onChange={onChangeServerURL}
            value={serverURL}
          />
        </div>
        <div className={`${baseClass}__details`}>
          <IconToolTip
            text={"The base URL of this instance for use in Fleet links."}
          />
        </div>
      </div>
    );
  };

  const renderSAMLSingleSignOnOptionsSection = () => {
    return (
      <div className={`${baseClass}__section`}>
        <h2>
          <a id="saml">SAML single sign on options</a>
        </h2>

        <div className={`${baseClass}__inputs`}>
          <Checkbox onChange={onChangeEnableSSO} value={enableSSO}>
            Enable single sign on
          </Checkbox>
        </div>

        <div className={`${baseClass}__inputs`}>
          <InputField
            label="Identity provider name"
            onChange={onChangeIDPName}
            value={idpName}
          />
        </div>
        <div className={`${baseClass}__details`}>
          <IconToolTip
            text={
              "A required human friendly name for the identity provider that will provide single sign on authentication."
            }
          />
        </div>

        <div className={`${baseClass}__inputs`}>
          <InputField
            label="Entity ID"
            hint={
              <span>
                The URI you provide here must exactly match the Entity ID field
                used in identity provider configuration.
              </span>
            }
            onChange={onChangeEntityID}
            value={entityID}
          />
        </div>
        <div className={`${baseClass}__details`}>
          <IconToolTip
            text={
              "The required entity ID is a URI that you use to identify Fleet when configuring the identity provider."
            }
          />
        </div>

        <div className={`${baseClass}__inputs`}>
          <InputField
            label="Issuer URI"
            onChange={onChangeIssuerURI}
            value={issuerURI}
          />
        </div>
        <div className={`${baseClass}__details`}>
          <IconToolTip
            text={"The issuer URI supplied by the identity provider."}
          />
        </div>

        <div className={`${baseClass}__inputs`}>
          <InputField
            label="IDP image URL"
            onChange={onChangeIDPImageURL}
            value={idpImageURL}
          />
        </div>
        <div className={`${baseClass}__details`}>
          <IconToolTip
            text={
              "An optional link to an image such as a logo for the identity provider."
            }
          />
        </div>

        <div className={`${baseClass}__inputs`}>
          <InputField
            label="Metadata"
            type="textarea"
            onChange={onChangeMetadata}
            value={metadata}
          />
        </div>
        <div className={`${baseClass}__details`}>
          <IconToolTip
            text={
              "Metadata provided by the identity provider. Either metadata or a metadata url must be provided."
            }
          />
        </div>

        <div className={`${baseClass}__inputs`}>
          <InputField
            label="Metadata URL"
            hint={
              <span>
                If available from the identity provider, this is the preferred
                means of providing metadata.
              </span>
            }
            onChange={onChangeMetadataURL}
            value={metadataURL}
          />
        </div>
        <div className={`${baseClass}__details`}>
          <IconToolTip
            text={"A URL that references the identity provider metadata."}
          />
        </div>

        <div className={`${baseClass}__inputs`}>
          <Checkbox
            onChange={onChangeEnableSSOIDPLogin}
            value={enableSSOIDLLogin}
          >
            Allow SSO login initiated by Identity Provider
          </Checkbox>
        </div>
      </div>
    );
  };

  const renderSMTPOptionsSection = () => {
    const renderSmtpSection = () => {
      if (smtpAuthenticationType === "authtype_none") {
        return false;
      }

      return (
        <div className={`${baseClass}__smtp-section`}>
          <InputField
            label="SMTP username"
            onChange={onChangeSMTPUsername}
            value={smtpUsername}
          />
          <InputField
            label="SMTP password"
            type="password"
            onChange={onChangeSMTPPassword}
            value={smtpPassword}
          />
          <Dropdown
            label="Auth method"
            options={authMethodOptions}
            placeholder=""
            onChange={onChangeSMTPAuthenticationMethod}
            value={smtpAuthenticationMethod}
          />
        </div>
      );
    };

    return (
      <div className={`${baseClass}__section`}>
        <h2>
          <a id="smtp">
            SMTP options{" "}
            <small
              className={`smtp-options smtp-options--${
                smtpConfigured ? "configured" : "notconfigured"
              }`}
            >
              STATUS:{" "}
              <em>{smtpConfigured ? "CONFIGURED" : "NOT CONFIGURED"}</em>
            </small>
          </a>
        </h2>
        <div className={`${baseClass}__inputs`}>
          <Checkbox onChange={onChangeEnableSMTP} value={enableSMTP}>
            Enable SMTP
          </Checkbox>
        </div>

        <div className={`${baseClass}__inputs`}>
          <InputField
            label="Sender address"
            onChange={onChangeSMTPSenderAddress}
            value={smtpSenderAddress}
          />
        </div>
        <div className={`${baseClass}__details`}>
          <IconToolTip text={"The sender address for emails from Fleet."} />
        </div>

        <div className={`${baseClass}__inputs ${baseClass}__inputs--smtp`}>
          <InputField
            label="SMTP server"
            onChange={onChangeSMTPServer}
            value={smtpServer}
          />
          <InputField
            label="&nbsp;"
            type="number"
            onChange={onChangeSMTPPort}
            value={smtpPort}
          />
          <Checkbox
            onChange={onChangeSMTPEnableSSLTLS}
            value={smtpEnableSSLTLS}
          >
            Use SSL/TLS to connect (recommended)
          </Checkbox>
        </div>
        <div className={`${baseClass}__details`}>
          <IconToolTip
            text={
              "The hostname / IP address and corresponding port of your organization's SMTP server."
            }
          />
        </div>

        <div className={`${baseClass}__inputs`}>
          <Dropdown label="Authentication type" options={authTypeOptions} />
          {renderSmtpSection()}
        </div>
        <div className={`${baseClass}__details`}>
          <IconToolTip
            isHtml
            text={
              "\
                  <p>If your mail server requires authentication, you need to specify the authentication type here.</p> \
                  <p><strong>No Authentication</strong> - Select this if your SMTP is open.</p> \
                  <p><strong>Username & Password</strong> - Select this if your SMTP server requires authentication with a username and password.</p>\
                "
            }
          />
        </div>
      </div>
    );
  };

  const renderOsqueryEnrollmentSecretsSection = () => {
    return (
      <div className={`${baseClass}__section`}>
        <h2>
          <a id="osquery-enrollment-secrets">Osquery enrollment secrets</a>
        </h2>
        <div className={`${baseClass}__inputs`}>
          <p className={`${baseClass}__enroll-secret-label`}>
            Manage secrets with <code>fleetctl</code>. Active secrets:
          </p>
          <EnrollSecretTable secrets={enrollSecret} />
        </div>
      </div>
    );
  };

  const renderGlobalAgentOptionsSection = () => {
    return (
      <div className={`${baseClass}__section`}>
        <h2>
          <a id="agent-options">Global agent options</a>
        </h2>
        <div className={`${baseClass}__yaml`}>
          <p className={`${baseClass}__section-description`}>
            This code will be used by osquery when it checks for configuration
            options.
            <br />
            <b>
              Changes to these configuration options will be applied to all
              hosts in your organization that do not belong to any team.
            </b>
          </p>
          <InfoBanner className={`${baseClass}__config-docs`}>
            How do global agent options interact with team-level agent
            options?&nbsp;
            <a
              href="https://github.com/fleetdm/fleet/blob/2f42c281f98e39a72ab4a5125ecd26d303a16a6b/docs/1-Using-Fleet/1-Fleet-UI.md#configuring-agent-options"
              className={`${baseClass}__learn-more ${baseClass}__learn-more--inline`}
              target="_blank"
              rel="noopener noreferrer"
            >
              Learn more about agent options&nbsp;
              <img className="icon" src={OpenNewTabIcon} alt="open new tab" />
            </a>
          </InfoBanner>
          <p className={`${baseClass}__component-label`}>
            <b>YAML</b>
          </p>
          <YamlAce
            // onChange={onChangeAgentOptions} TODO
            // value={agentOptions} TODO
            // error={fields.agent_options.error} TODO
            wrapperClassName={`${baseClass}__text-editor-wrapper`}
          />
          {/* this might be tricky */}
        </div>
      </div>
    );
  };

  const renderHostStatusWebhookSection = () => {
    return (
      <div className={`${baseClass}__section`}>
        <h2>
          <a id="host-status-webhook">Host status webhook</a>
        </h2>
        <div className={`${baseClass}__host-status-webhook`}>
          <p className={`${baseClass}__section-description`}>
            Send an alert if a portion of your hosts go offline.
          </p>
          <Checkbox
            onChange={onChangeEnableHostStatusWebhook}
            value={enableHostStatusWebhook}
          >
            Enable host status webhook
          </Checkbox>
          <p className={`${baseClass}__section-description`}>
            A request will be sent to your configured <b>Destination URL</b> if
            the configured <b>Percentage of hosts</b> have not checked into
            Fleet for the configured <b>Number of days</b>.
          </p>
        </div>
        <div className={`${baseClass}__inputs ${baseClass}__inputs--webhook`}>
          <Button
            type="button"
            variant="inverse"
            onClick={toggleHostStatusWebhookPreviewModal}
          >
            Preview request
          </Button>
        </div>
        <div className={`${baseClass}__inputs`}>
          <InputField
            placeholder="https://server.com/example"
            label="Destination URL"
            onChange={onChangeHostStatusWebhookDestinationURL}
            value={hostStatusWebhookDestinationURL}
          />
        </div>
        <div className={`${baseClass}__details`}>
          <IconToolTip
            isHtml
            text={
              "\
                  <center><p>Provide a URL to deliver <br/>the webhook request to.</p></center>\
                "
            }
          />
        </div>
        <div className={`${baseClass}__inputs ${baseClass}__host-percentage`}>
          <Dropdown
            label="Percentage of hosts"
            options={percentageOfHosts}
            onChange={onChangeHostStatusWebhookHostPercentage}
            value={hostStatusWebhookHostPercentage}
          />
        </div>
        <div className={`${baseClass}__details`}>
          <IconToolTip
            isHtml
            text={
              "\
                  <center><p>Select the minimum percentage of hosts that<br/>must fail to check into Fleet in order to trigger<br/>the webhook request.</p></center>\
                "
            }
          />
        </div>
        <div className={`${baseClass}__inputs ${baseClass}__days-count`}>
          <Dropdown
            label="Number of days"
            options={numberOfDays}
            onChange={onChangeHostStatusWebhookDaysCount}
            value={hostStatusWebhookDaysCount}
          />
        </div>
        <div className={`${baseClass}__details`}>
          <IconToolTip
            isHtml
            text={
              "\
                  <center><p>Select the minimum number of days that the<br/>configured <b>Percentage of hosts</b> must fail to<br/>check into Fleet in order to trigger the<br/>webhook request.</p></center>\
                "
            }
          />
        </div>
      </div>
    );
  };

  const renderUsageStatistics = () => {
    return (
      <div className={`${baseClass}__section`}>
        <h2>
          <a id="usage-stats">Usage statistics</a>
        </h2>
        <p className={`${baseClass}__section-description`}>
          Help improve Fleet by sending anonymous usage statistics.
          <br />
          <br />
          This information helps our team better understand feature adoption and
          usage, and allows us to see how Fleet is adding value, so that we can
          make better product decisions.
          <br />
          <br />
          <a
            href="https://github.com/fleetdm/fleet/blob/2f42c281f98e39a72ab4a5125ecd26d303a16a6b/docs/1-Using-Fleet/11-Usage-statistics.md"
            className={`${baseClass}__learn-more`}
            target="_blank"
            rel="noopener noreferrer"
          >
            Learn more about usage statistics&nbsp;
            <img className="icon" src={OpenNewTabIcon} alt="open new tab" />
          </a>
        </p>
        <div className={`${baseClass}__inputs ${baseClass}__inputs--usage`}>
          <Checkbox
            onChange={onChangeEnableUsageStatistics}
            value={enableUsageStatistics}
          >
            Enable usage statistics
          </Checkbox>
        </div>
        <div className={`${baseClass}__inputs ${baseClass}__inputs--usage`}>
          <Button
            type="button"
            variant="inverse"
            onClick={toggleUsageStatsPreviewModal}
          >
            Preview payload
          </Button>
        </div>
      </div>
    );
  };

  const renderAdvancedOptions = () => {
    return (
      <div className={`${baseClass}__section`}>
        <h2>
          <a id="advanced-options">Advanced options</a>
        </h2>
        <div className={`${baseClass}__advanced-options`}>
          <p className={`${baseClass}__section-description`}>
            Most users do not need to modify these options.
          </p>
          <div className={`${baseClass}__inputs`}>
            <div className={`${baseClass}__form-fields`}>
              <div className="tooltip-wrap tooltip-wrap--input">
                <InputField
                  label="Domain"
                  onChange={onChangeDomain}
                  value={domain}
                />
                <IconToolTip
                  isHtml
                  text={
                    '<p>If you need to specify a HELO domain, <br />you can do it here <em className="hint hint--brand">(Default: <strong>Blank</strong>)</em></p>'
                  }
                />
              </div>
              <div className="tooltip-wrap">
                <Checkbox
                  onChange={onChangeVerifySSLCerts}
                  value={verifySSLCerts}
                >
                  Verify SSL certs
                </Checkbox>
                <IconToolTip
                  isHtml
                  text={
                    '<p>Turn this off (not recommended) <br />if you use a self-signed certificate <em className="hint hint--brand"><br />(Default: <strong>On</strong>)</em></p>'
                  }
                />
              </div>
              <div className="tooltip-wrap">
                <Checkbox
                  onChange={onChangeEnableStartTLS}
                  value={enableStartTLS}
                >
                  Enable STARTTLS
                </Checkbox>
                <IconToolTip
                  isHtml
                  text={
                    '<p>Detects if STARTTLS is enabled <br />in your SMTP server and starts <br />to use it. <em className="hint hint--brand">(Default: <strong>On</strong>)</em></p>'
                  }
                />
              </div>
              <div className="tooltip-wrap">
                <Checkbox
                  onChange={onChangeEnableHostExpiry}
                  value={enableHostExpiry}
                >
                  Host expiry
                </Checkbox>
                <IconToolTip
                  isHtml
                  text={
                    '<p>When enabled, allows automatic cleanup <br />of hosts that have not communicated with Fleet <br />in some number of days. <em className="hint hint--brand">(Default: <strong>Off</strong>)</em></p>'
                  }
                />
              </div>
              <div className="tooltip-wrap tooltip-wrap--input">
                <InputField
                  onChange={onChangeHostExpiryWindow}
                  value={hostExpiryWindow}
                  // disabled={false} TODO!
                  label="Host Expiry Window"
                />
                <IconToolTip
                  isHtml
                  text={
                    "<p>If a host has not communicated with Fleet <br />in the specified number of days, it will be removed.</p>"
                  }
                />
              </div>
              <div className="tooltip-wrap">
                <Checkbox
                  onChange={onChangeDisableLiveQuery}
                  value={disableLiveQuery}
                >
                  Disable live queries
                </Checkbox>
                <IconToolTip
                  isHtml
                  text={
                    '<p>When enabled, disables the ability to run live queries <br />(ad hoc queries executed via the UI or fleetctl). <em className="hint hint--brand">(Default: <strong>Off</strong>)</em></p>'
                  }
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  };

  // █▀▄▀█ █▀█ █▀▄ ▄▀█ █░░ █▀
  // █░▀░█ █▄█ █▄▀ █▀█ █▄▄ ▄█
  const renderHostStatusWebhookPreviewModal = () => {
    if (!showHostStatusWebhookPreviewModal) {
      return null;
    }

    const json = {
      text:
        "More than X% of your hosts have not checked into Fleet for more than Y days. You’ve been sent this message because the Host status webhook is enabled in your Fleet instance.",
      data: {
        unseen_hosts: 1,
        total_hosts: 2,
        days_unseen: 3,
      },
    };

    return (
      <Modal
        title="Host status webhook"
        onExit={toggleHostStatusWebhookPreviewModal}
        className={`${baseClass}__host-status-webhook-preview-modal`}
      >
        <>
          <p>
            An example request sent to your configured <b>Destination URL</b>.
          </p>
          <div className={`${baseClass}__host-status-webhook-preview`}>
            <pre dangerouslySetInnerHTML={{ __html: syntaxHighlight(json) }} />
          </div>
          <div className="flex-end">
            <Button type="button" onClick={toggleHostStatusWebhookPreviewModal}>
              Done
            </Button>
          </div>
        </>
      </Modal>
    );
  };

  const renderUsageStatsPreviewModal = () => {
    if (!showUsageStatsPreviewModal) {
      return null;
    }

    const stats = {
      anonymousIdentifier: "9pnzNmrES3mQG66UQtd29cYTiX2+fZ4CYxDvh495720=",
      fleetVersion: "x.x.x",
      licenseTier: "free",
      numHostsEnrolled: 12345,
      numUsers: 12,
      numTeams: 3,
      numPolicies: 5,
      numLabels: 20,
      softwareInventoryEnabled: true,
      vulnDetectionEnabled: true,
      systemUsersEnabled: true,
      hostStatusWebhookEnabled: true,
    };

    return (
      <Modal
        title="Usage statistics"
        onExit={toggleUsageStatsPreviewModal}
        className={`${baseClass}__usage-stats-preview-modal`}
      >
        <>
          <p>An example JSON payload sent to Fleet Device Management Inc.</p>
          <pre dangerouslySetInnerHTML={{ __html: syntaxHighlight(stats) }} />
          <div className="flex-end">
            <Button type="button" onClick={toggleUsageStatsPreviewModal}>
              Done
            </Button>
          </div>
        </>
      </Modal>
    );
  };

  // █▀█ █▀▀ █▄░█ █▀▄ █▀▀ █▀█
  // █▀▄ ██▄ █░▀█ █▄▀ ██▄ █▀▄
  return (
    <>
      <form className={baseClass} onSubmit={handleSubmit} autoComplete="off">
        {renderOrganizationInfoSection()}
        {renderFleetWebAddressSection()}
        {renderSAMLSingleSignOnOptionsSection()}
        {renderSMTPOptionsSection()}
        {renderOsqueryEnrollmentSecretsSection()}
        {renderGlobalAgentOptionsSection()}
        {renderHostStatusWebhookSection()}
        {renderUsageStatistics()}
        {renderAdvancedOptions()}
        <Button type="submit" variant="brand">
          Update settings
        </Button>
        {/* this should rerender the page or scroll to top */}
      </form>
      {renderUsageStatsPreviewModal()}
      {renderHostStatusWebhookPreviewModal()}
    </>
  );
};

export default AppConfigFormFunctional;

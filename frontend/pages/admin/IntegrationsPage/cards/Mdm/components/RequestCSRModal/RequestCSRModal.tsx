import React, { FormEvent, useState, useContext, useEffect } from "react";

import { AppContext } from "context/app";

import MdmAPI from "services/entities/mdm";

import Button from "components/buttons/Button";
// @ts-ignore
import InputField from "components/forms/fields/InputField";
import DataError from "components/DataError";
import Icon from "components/Icon";
import Modal from "components/Modal";

export interface IRequestCSRFormData {
  email: string;
  orgName: string;
}

const baseClass = "modal request-csr-modal";
interface IRequestCSRModalProps {
  onCancel: () => void;
}

interface IFormField {
  name: string;
  value: string;
}

const FILES: CSRFile[] = [
  { name: "mdmcert.download.push.key", key: "apns_key" }, // APNS key
  { name: "fleet-mdm-apple-scep.key", key: "scep_key" }, // SCEP key
  { name: "fleet-mdm-apple-scep.crt", key: "scep_cert" }, // SCEP cert
];

const downloadFile = (tokens: string, fileName: string) => {
  const linkSource = `data:application/octet-stream;base64,${tokens}`;
  const downloadLink = document.createElement("a");

  downloadLink.href = linkSource;
  downloadLink.download = fileName;
  downloadLink.click();
};

type RequestCsrResponse = {
  apns_key: string;
  scep_key: string;
  scep_cert: string;
};

type ResponseKeys = keyof RequestCsrResponse;

type CSRFile = {
  name: string;
  key: ResponseKeys;
  value?: string;
};

const downloadCSRFiles = (data: RequestCsrResponse) => {
  FILES.forEach((file) => {
    downloadFile(data[file.key], file.name);
  });
};

const RequestCSRModal = ({ onCancel }: IRequestCSRModalProps): JSX.Element => {
  const { currentUser, config } = useContext(AppContext);

  const [formData, setFormData] = useState<IRequestCSRFormData>({
    email: currentUser?.email ?? "",
    orgName: config?.org_info?.org_name ?? "",
  });
  const [error, setError] = useState("");
  const [requestState, setRequestState] = useState<
    "loading" | "error" | "success" | undefined
  >(undefined);

  const { email, orgName } = formData;

  const onInputChange = ({ name, value }: IFormField) => {
    setFormData({ ...formData, [name]: value });
  };

  const onFormSubmit = async (evt: FormEvent) => {
    setRequestState("loading");
    evt.preventDefault();
    try {
      const data = await MdmAPI.requestCSR(email, orgName);
      downloadCSRFiles(data);
      setRequestState("success");
    } catch (e) {
      const err = e as any;
      console.log(e);
      if (err.status === 502) {
        setRequestState("error");
      } else {
        setError("This email is not valid.");
        setRequestState(undefined);
      }
    }
  };

  const RequestCSRSuccess = () => {
    return (
      <div className="success">
        <Icon name="success" size="extra-large" />
        <h2>You&apos;re almost there</h2>
        <p>
          Go to your <strong>{email}</strong> email to download your CSR.
        </p>
        <p>
          Your APNs key and SCEP certificate and key will be downloaded in the
          browser.
          <br />
          You&apos;ll need these later.
        </p>
        <Button
          onClick={() => {
            onCancel();
          }}
        >
          Got it
        </Button>
      </div>
    );
  };

  const renderRequestCSRForm = () => {
    if (requestState === "success") {
      return <RequestCSRSuccess />;
    }
    if (requestState === "error") {
      return <DataError />;
    }
    return (
      <>
        <p>
          A CSR and key for APNs and a certificate and key for SCEP are required
          to connect Fleet to Apple Developer. Apple Inc. requires the following
          information. <br />
          <br />
          fleetdm.com will send your CSR to the below email. Your APNs key and
          SCEP certificate and key will be downloaded in the browser.
        </p>
        <form
          className={`${baseClass}__form`}
          onSubmit={onFormSubmit}
          autoComplete="off"
        >
          <div className="bottom-label">
            <InputField
              name="email"
              onChange={onInputChange}
              label="Email"
              parseTarget
              value={email}
              error={error}
            />
            <p>
              Apple Inc. requires a work email (ex. name@your-organization.com).
            </p>
          </div>
          <InputField
            name="orgName"
            onChange={onInputChange}
            label="Organization name"
            parseTarget
            value={orgName}
          />
          <div className="modal-cta-wrap">
            <Button
              type="submit"
              variant="brand"
              isLoading={requestState === "loading"}
            >
              Request
            </Button>
            <Button onClick={onCancel} variant="inverse">
              Cancel
            </Button>
          </div>
        </form>
      </>
    );
  };

  return (
    <Modal title="Request" onExit={onCancel} className={baseClass}>
      {renderRequestCSRForm()}
    </Modal>
  );
};

export default RequestCSRModal;

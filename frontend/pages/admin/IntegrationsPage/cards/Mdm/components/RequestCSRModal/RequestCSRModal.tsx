import React from "react";

import Modal from "components/Modal";
import requestCSR from "services/entities/mdm_csr";
import RequestCSRForm from "./RequestCSRForm";

const baseClass = " modal request-csr-modal";

interface IRequestCSRModalProps {
  onCancel: () => void;
}

const RequestCSRModal = ({ onCancel }: IRequestCSRModalProps): JSX.Element => {
  return (
    <Modal title="Request" onExit={onCancel} className={baseClass}>
      <>
        <p>
          A CSR and key for APNs and a certificate and key for SCEP are required
          to connect Fleet to Apple Developer. Apple Inc. requires the following
          information. <br />
          <br />
          fleetdm.com will send your CSR to the below email. Your certificate
          and key for SCEP will be downloaded in the browser.
        </p>
        <RequestCSRForm onCancel={onCancel} />
      </>
    </Modal>
  );
};

export default RequestCSRModal;

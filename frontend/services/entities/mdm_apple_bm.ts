/* eslint-disable @typescript-eslint/explicit-module-boundary-types */
import { sendRequest } from "services/mock_service/service/service";
// import sendRequest from "services";
import endpoints from "utilities/endpoints";

export default {
  loadAll: () => {
    const { MDM_APPLE_BM } = endpoints;
    const path = MDM_APPLE_BM;
    return sendRequest("GET", path);
  },
  loadKeys: () => {
    const { MDM_APPLE } = endpoints;
    const path = `${MDM_APPLE}/keys`;

    // TODO: certificate_chain for certificate, refactor for keys
    return sendRequest("GET", path).then(({ certificate_chain }) => {
      let decodedKeys;
      try {
        decodedKeys = global.window.atob(certificate_chain);
      } catch (err) {
        return Promise.reject(`Unable to decode keys: ${err}`);
      }
      if (!decodedKeys) {
        return Promise.reject("Missing or undefined keys.");
      }

      return Promise.resolve(decodedKeys);
    });
  },
};

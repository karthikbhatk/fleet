import React from "react";

import { ISoftware } from "interfaces/software";
<<<<<<< HEAD:frontend/pages/hosts/SoftwareTab/SoftwareVulnCount/SoftwareVulnCount.tsx
import IssueIcon from "../../../../../assets/images/icon-issue-fleet-black-50-16x16@2x.png";
=======
import IssueIcon from "../../../../../../../assets/images/icon-issue-fleet-black-50-16x16@2x.png";
>>>>>>> b2894709e (Refactor host details page into components):frontend/pages/hosts/HostDetailsPage/cards/Software/SoftwareVulnCount/SoftwareVulnCount.tsx

const baseClass = "software-vuln-count";

interface ISoftwareVulnCountProps {
  softwareList: ISoftware[];
  deviceUser?: boolean;
}

const SoftwareVulnCount = ({
  softwareList,
  deviceUser,
}: ISoftwareVulnCountProps): JSX.Element => {
  const vulnCount = softwareList.reduce((sum, software) => {
    return software.vulnerabilities
      ? sum + software.vulnerabilities.length
      : sum;
  }, 0);

  return vulnCount ? (
    <div className={`${baseClass}`}>
      <div className={`${baseClass}__count`}>
        <img alt="Issue icon" src={IssueIcon} />
        {vulnCount === 1
          ? "1 vulnerability detected"
          : `${vulnCount} vulnerabilities detected`}
      </div>
      {!deviceUser && (
        <p>
          Click a vulnerable item below to see the associated Common
          Vulnerabilites and Exposures (CVEs).
        </p>
      )}
    </div>
  ) : (
    <></>
  );
};

export default SoftwareVulnCount;

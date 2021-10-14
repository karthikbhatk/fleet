import React, { useContext } from "react";
import { AppContext } from "context/app";

import paths from "router/paths";
import { Link } from "react-router";
import HostsSummary from "./HostsSummary";
import ActivityFeed from "./ActivityFeed";

import LinkArrow from "../../../assets/images/icon-arrow-right-vibrant-blue-10x18@2x.png";

const baseClass = "homepage";

const Homepage = (): JSX.Element => {
  const { MANAGE_HOSTS } = paths;
  const { config } = useContext(AppContext);

  return (
    <div className={baseClass}>
      <div className={`${baseClass}__wrapper body-wrap`}>
        <div className={`${baseClass}__header-wrap`}>
          <div className={`${baseClass}__header`}>
            <h1 className={`${baseClass}__title`}>
              <span>{config?.org_name}</span>
            </h1>
          </div>
        </div>
        <div className={`${baseClass}__section hosts-section`}>
          <div className={`${baseClass}__section-title`}>
            <div>
              <h2>Hosts</h2>
            </div>
            <Link to={MANAGE_HOSTS} className={`${baseClass}__host-link`}>
              <span>View all hosts</span>
              <img src={LinkArrow} alt="link arrow" id="link-arrow" />
            </Link>
          </div>
          <div className={`${baseClass}__section-details`}>
            <HostsSummary />
          </div>
        </div>
        <div className={`${baseClass}__section hosts-section`}>
          <div className={`${baseClass}__section-title`}>
            <div>
              <h2>Activity</h2>
            </div>
          </div>
          <div className={`${baseClass}__section-details`}>
            <ActivityFeed />
          </div>
        </div>
      </div>
    </div>
  );
};

export default Homepage;

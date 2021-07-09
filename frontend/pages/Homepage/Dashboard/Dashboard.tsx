import React, { useState, useCallback } from "react"; //, { useEffect }
import { useDispatch, useSelector } from "react-redux";

import { Link } from "react-router";
import { IUser } from "interfaces/user";
import HostSummary from "./HostSummary";

import paths from "router/paths";
import LinkArrow from "../../../../assets/images/icon-arrow-right-vibrant-blue-10x18@2x.png";

const baseClass = "dashboard";

interface RootState {
  auth: {
    user: IUser;
  };
  app: {
    config: {
      org_name: string;
    };
  };
}

const Dashboard = (): JSX.Element => {
  // Links to packs page
  const dispatch = useDispatch();
  const { MANAGE_HOSTS } = paths;

  const user = useSelector((state: RootState) => state.auth.user);
  const orgName = useSelector((state: RootState) => state.app.config.org_name);

  console.log("USER", user);
  return (
    <div className={baseClass}>
      <div className={`${baseClass}__wrapper body-wrap`}>
        <div className={`${baseClass}__header-wrap`}>
          <div className={`${baseClass}__header`}>
            <h1 className={`${baseClass}__title`}>
              <span>{orgName}</span>
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
            <HostSummary />
          </div>
        </div>
        <div className={`${baseClass}__section activity-section`}>
          <div className={`${baseClass}__section-title`}>
            <h2>Activity</h2>
          </div>
          <div className={`${baseClass}__section-details`}>Details</div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;

import React, { useState, useEffect, useContext } from "react";
import { InjectedRouter } from "react-router";
import { Params } from "react-router/lib/Router";

import paths from "router/paths";
import { AppContext } from "context/app";
import usersAPI from "services/entities/users";
import sessionsAPI from "services/entities/sessions"; // @ts-ignore
import { formatErrorResponse } from "redux/nodes/entities/base/helpers";

// @ts-ignore
import AuthenticationFormWrapper from "components/AuthenticationFormWrapper"; // @ts-ignore
import ConfirmSSOInviteForm from "components/forms/ConfirmSSOInviteForm";

interface IConfirmSSOInvitePageProps {
  location: any; // no type in react-router v3
  params: Params;
  router: InjectedRouter;
}

const baseClass = "confirm-invite-page";

const ConfirmSSOInvitePage = ({
  location,
  params,
  router,
}: IConfirmSSOInvitePageProps) => {
  const { email, name } = location.query;
  const { invite_token } = params;
  const inviteFormData = { email, invite_token, name };
  const { currentUser } = useContext(AppContext);
  const [errors, setErrors] = useState<{ [key: string]: string }>({});

  useEffect(() => {
    const { HOME } = paths;

    if (currentUser) {
      return router.push(HOME);
    }
  }, [currentUser]);
  
  const onSubmit = async (formData: any) => {
    const { HOME } = paths;

    formData.sso_invite = true;

    try {
      await usersAPI.create(formData);
      const { url } = await sessionsAPI.initializeSSO(HOME);
      window.location.href = url;
    } catch (response) {
      const errorObject = formatErrorResponse(response);
      setErrors(errorObject);
      return false;
    }
  };

  return (
    <AuthenticationFormWrapper>
      <div className={`${baseClass}`}>
        <div className={`${baseClass}__lead-wrapper`}>
          <p className={`${baseClass}__lead-text`}>Welcome to Fleet</p>
          <p className={`${baseClass}__sub-lead-text`}>
            Before you get started, please take a moment to complete the
            following information.
          </p>
        </div>
        <ConfirmSSOInviteForm
          className={`${baseClass}__form`}
          formData={inviteFormData}
          handleSubmit={onSubmit}
          serverErrors={errors}
        />
      </div>
    </AuthenticationFormWrapper>
  );
}

export default ConfirmSSOInvitePage;
import React, { Component } from "react";
import PropTypes from "prop-types";
import { connect } from "react-redux";
import { logoutUser } from "redux/nodes/auth/actions";
import { push } from "react-router-redux";
import { TableContext } from "context/table";

import { isEqual } from "lodash";

import configInterface from "interfaces/config";
import FlashMessage from "components/flash_messages/FlashMessage";
import PersistentFlash from "components/flash_messages/PersistentFlash";
import SiteTopNav from "components/side_panels/SiteTopNav";
import userInterface from "interfaces/user";
import notificationInterface from "interfaces/notification";
import { hideFlash } from "redux/nodes/notifications/actions";
import { licenseExpirationWarning } from "fleet/helpers";

const expirationMessage = (
  <>
    Your license for Fleet Premium is about to expire. If you’d like to renew or
    have questions about downgrading,{" "}
    <a
      href="https://github.com/fleetdm/fleet/blob/main/docs/1-Using-Fleet/10-Teams.md#expired_license"
      target="_blank"
      rel="noopener noreferrer"
    >
      please head to the Fleet documentation
    </a>
    .
  </>
);
export class CoreLayout extends Component {
  static propTypes = {
    children: PropTypes.node,
    config: configInterface,
    dispatch: PropTypes.func,
    user: userInterface,
    fullWidthFlash: PropTypes.bool,
    notifications: notificationInterface,
    persistentFlash: PropTypes.shape({
      showFlash: PropTypes.bool.isRequired,
      message: PropTypes.string.isRequired,
    }).isRequired,
  };

  constructor(props) {
    super(props);

    this.state = {
      showExpirationFlashMessage: true,
    };
  }

  componentWillReceiveProps(nextProps) {
    const { notifications } = nextProps;
    const table = this.context;

    // on success of an action, the table will reset its checkboxes.
    // setTimeout is to help with race conditions as table reloads
    // in some instances (i.e. Manage Hosts)
    if (!isEqual(this.props.notifications, notifications)) {
      if (notifications.alertType === "success") {
        setTimeout(() => {
          table.setResetSelectedRows(true);
          setTimeout(() => {
            table.setResetSelectedRows(false);
          }, 300);
        }, 0);
      }
    }
  }

  onLogoutUser = () => {
    const { dispatch } = this.props;

    dispatch(logoutUser());

    return false;
  };

  onNavItemClick = (path) => {
    return (evt) => {
      evt.preventDefault();

      const { dispatch } = this.props;

      if (path.indexOf("http") !== -1) {
        global.window.open(path, "_blank");

        return false;
      }

      dispatch(push(path));

      return false;
    };
  };

  onRemoveFlash = () => {
    const { dispatch } = this.props;

    dispatch(hideFlash);

    return false;
  };

  onRemoveExpirationWarning = () => {
    const { showExpirationFlashMessage } = this.state;

    this.setState({
      showExpirationFlashMessage: !showExpirationFlashMessage,
    });
  };

  onUndoActionClick = (undoAction) => {
    return (evt) => {
      evt.preventDefault();

      const { dispatch } = this.props;
      const { onRemoveFlash } = this;

      dispatch(undoAction);

      return onRemoveFlash();
    };
  };

  static contextType = TableContext;

  render() {
    const {
      fullWidthFlash,
      notifications,
      children,
      config,
      persistentFlash,
      user,
    } = this.props;
    const { showExpirationFlashMessage } = this.state;
    const {
      onRemoveFlash,
      onRemoveExpirationWarning,
      onUndoActionClick,
    } = this;

    const expirationNotification = {
      alertType: "warning-filled",
      isVisible: licenseExpirationWarning(config.expiration),
      message: expirationMessage,
    };

    if (!user) return false;

    const { onLogoutUser, onNavItemClick } = this;
    const { pathname } = global.window.location;

    return (
      <div className="app-wrap">
        <nav className="site-nav">
          <SiteTopNav
            config={config}
            onLogoutUser={onLogoutUser}
            onNavItemClick={onNavItemClick}
            pathname={pathname}
            user={user}
          />
        </nav>
        <div className="core-wrapper">
          {persistentFlash.showFlash && (
            <PersistentFlash message={persistentFlash.message} />
          )}
          {/* TODO: This will need a transition from redux to context API */}
          {expirationNotification && showExpirationFlashMessage && (
            <FlashMessage
              fullWidth={fullWidthFlash}
              notification={expirationNotification}
              onRemoveFlash={onRemoveExpirationWarning}
            />
          )}
          <FlashMessage
            fullWidth={fullWidthFlash}
            notification={notifications}
            onRemoveFlash={onRemoveFlash}
            onUndoActionClick={onUndoActionClick}
            className="app-flash"
          />
          {children}
        </div>
      </div>
    );
  }
}

const mapStateToProps = (state) => {
  const {
    app: { config },
    auth: { user },
    notifications,
    persistentFlash,
  } = state;

  const fullWidthFlash = !user;

  return {
    config,
    fullWidthFlash,
    notifications,
    persistentFlash,
    user,
  };
};

export default connect(mapStateToProps)(CoreLayout);

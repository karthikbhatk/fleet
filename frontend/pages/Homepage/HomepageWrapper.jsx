// Identical to PacksPageWrapper directory 6/28

import React from "react";
import PropTypes from "prop-types";

class HomepageWrapper extends React.Component {
  static propTypes = {
    children: PropTypes.node,
  };

  render() {
    const { children } = this.props;

    return children || null;
  }
}

export default HomepageWrapper;

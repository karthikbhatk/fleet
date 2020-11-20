import React, { Component } from 'react';

import fleetLogoText from '../../../assets/images/fleet-logo-text-white.svg';
import backgroundImg from '../../../assets/images/404.svg';

const baseClass = 'kolide-404';

class Kolide404 extends Component {
  render () {
    return (
      <div className={baseClass}>
        <header className="primary-header">
          <a href="/">
            <img className="primary-header__logo" src={fleetLogoText} alt="Kolide" />
          </a>
        </header>
        <img className="background-image" src={backgroundImg}></img>
        <main>
          <h1>404: Oops, sorry we can't find that page!</h1>
          <p>The page you are looking for has either moved, or doesn't exist.</p>
          <a href="https://fleetdm.com/support">Get help</a>
        </main>
      </div>
    );
  }
}

export default Kolide404;

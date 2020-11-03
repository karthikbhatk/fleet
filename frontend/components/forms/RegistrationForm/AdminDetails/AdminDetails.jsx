import React, { Component } from 'react';
import PropTypes from 'prop-types';

import Form from 'components/forms/Form';
import formFieldInterface from 'interfaces/form_field';
import Button from 'components/buttons/Button';
import InputFieldWithIcon from 'components/forms/fields/InputFieldWithIcon';
import helpers from './helpers';

const formFields = ['username', 'password', 'password_confirmation', 'email'];
const { validate } = helpers;

class AdminDetails extends Component {
  static propTypes = {
    className: PropTypes.string,
    currentPage: PropTypes.bool,
    fields: PropTypes.shape({
      email: formFieldInterface.isRequired,
      password: formFieldInterface.isRequired,
      password_confirmation: formFieldInterface.isRequired,
      username: formFieldInterface.isRequired,
    }).isRequired,
    handleSubmit: PropTypes.func.isRequired,
  };

  componentDidUpdate(prevProps) {
    if (this.props.currentPage && this.props.currentPage !== prevProps.currentPage) {
      // Component has a transition duration of 300ms set in
      // RegistrationForm/_styles.scss. We need to wait 300ms before
      // calling .focus() to preserve smooth transition.
      setTimeout(() => {
        this.firstInput.input.focus();
      }, 300);
    }
  }

  render () {
    const { className, currentPage, fields, handleSubmit } = this.props;
    const tabIndex = currentPage ? 1 : -1;

    return (
      <form onSubmit={handleSubmit} className={className}>
        <div className="registration-fields">
          <InputFieldWithIcon
            {...fields.username}
            iconName="username"
            placeholder="Username"
            tabIndex={tabIndex}
            autofocus={currentPage}
            ref={(input) => { this.firstInput = input; }}
          />
          <InputFieldWithIcon
            {...fields.password}
            iconName="password"
            placeholder="Password"
            type="password"
            tabIndex={tabIndex}
          />
          <InputFieldWithIcon
            {...fields.password_confirmation}
            iconName="password"
            placeholder="Confirm Password"
            type="password"
            tabIndex={tabIndex}
          />
          <InputFieldWithIcon
            {...fields.email}
            iconName="email"
            placeholder="Email"
            tabIndex={tabIndex}
          />
        </div>
        <Button type="submit" variant="gradient" tabIndex={tabIndex} disabled={!currentPage}>
          Submit
        </Button>
      </form>
    );
  }
}

export default Form(AdminDetails, {
  fields: formFields,
  validate,
});

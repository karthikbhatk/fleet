import React, { Component } from "react";
import PropTypes from "prop-types";

import Button from "components/buttons/Button";
// @ts-ignore
import Form from "components/forms/Form";
import { IFormField } from "interfaces/form_field";
import { IQuery } from "interfaces/query";
import { IPack } from "interfaces/pack";
import { ITarget } from "interfaces/target";
// @ts-ignore
import InputField from "components/forms/fields/InputField";
// @ts-ignore
import SelectTargetsDropdown from "components/forms/fields/SelectTargetsDropdown";

const fieldNames = ["description", "name", "targets"];
const baseClass = "edit-pack-form";

interface IEditPackForm {
  className?: string;
  handleSubmit?: (formData: any) => void;
  onCancelEditPack: () => void;
  onEditPack: () => void;
  onFetchTargets?: (query: IQuery, targetsResponse: any) => boolean;
  pack: IPack;
  packTargets?: ITarget[];
  targetsCount?: number;
  isPremiumTier?: boolean;
  fields: { description: IFormField; name: IFormField; targets: IFormField };
}
const EditPackForm = (props: IEditPackForm): JSX.Element => {
  const {
    className,
    handleSubmit,
    onCancelEditPack,
    onEditPack,
    onFetchTargets,
    pack,
    packTargets,
    targetsCount,
    isPremiumTier,
    fields,
  } = props;

  return (
    <form className={`${baseClass} ${className}`} onSubmit={handleSubmit}>
      <h1>Edit pack</h1>
      <InputField
        {...fields.name}
        placeholder="Name"
        label="Name"
        inputWrapperClass={`${baseClass}__pack-title`}
      />
      <InputField
        {...fields.description}
        inputWrapperClass={`${baseClass}__pack-description`}
        label="Description"
        placeholder="Add a description of your pack"
        type="textarea"
      />
      <SelectTargetsDropdown
        {...fields.targets}
        label="Select pack targets"
        name="selected-pack-targets"
        onFetchTargets={onFetchTargets}
        onSelect={fields.targets.onChange}
        selectedTargets={fields.targets.value}
        targetsCount={targetsCount}
        isPremiumTier={isPremiumTier}
      />
      <div className={`${baseClass}__pack-buttons`}>
        <Button onClick={onCancelEditPack} type="button" variant="inverse">
          Cancel
        </Button>
        <Button type="submit" variant="brand">
          Save
        </Button>
      </div>
    </form>
  );
};

export default Form(EditPackForm, {
  fields: fieldNames,
});

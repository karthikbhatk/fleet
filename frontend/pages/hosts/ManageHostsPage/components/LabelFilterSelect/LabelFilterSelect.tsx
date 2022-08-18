import React, { useMemo, useRef, useState } from "react";
import Select, { GroupBase, SelectInstance } from "react-select-5";
import classnames from "classnames";

import { ILabel } from "interfaces/label";
import { PLATFORM_LABEL_DISPLAY_NAMES } from "utilities/constants";

import CustomLabelGroupHeading from "../CustomLabelGroupHeading";
import { PLATFORM_TYPE_ICONS } from "./constants";
import { createDropdownOptions, IEmptyOption, IGroupOption } from "./helpers";

// Extending the react-select to add custom props we need for our custom group heading
// https://react-select.com/typescript#custom-select-props
declare module "react-select-5/dist/declarations/src/Select" {
  export interface Props<
    Option,
    IsMulti extends boolean,
    Group extends GroupBase<Option>
  > {
    labelQuery: string;
    canAddNewLabels: boolean;
    onAddLabel: () => void;
    onChangeLabelQuery: (event: React.ChangeEvent<HTMLInputElement>) => void;
    onClickLabelSeachInput: React.MouseEventHandler<HTMLInputElement>;
    onBlurLabelSearchInput: React.FocusEventHandler<HTMLInputElement>;
  }
}

/** A custom option label to show in the dropdown. Only used in this dropdown
 * component */
const OptionLabel = (data: ILabel | IEmptyOption) => {
  const isLabel = "display_text" in data;
  const isPlatform = isLabel && data.type === "platform";

  let labelText = isLabel ? data.display_text : data.label;

  // the display names for platform options are slightly different then the display_text
  // property, so we get the correct display name here
  if (isLabel && isPlatform) {
    labelText = PLATFORM_LABEL_DISPLAY_NAMES[data.display_text];
  }

  return (
    <div className={"option-label"}>
      {isPlatform && (
        <img src={PLATFORM_TYPE_ICONS[data.display_text]} alt="" />
      )}
      <span>{labelText}</span>
    </div>
  );
};

const baseClass = "label-filter-select";

interface ILabelFilterSelectProps {
  labels: ILabel[];
  selectedLabel: ILabel | null;
  canAddNewLabels: boolean;
  className?: string;
  onChange: (labelId: ILabel) => void;
  onAddLabel: () => void;
}

const LabelFilterSelect = ({
  labels,
  selectedLabel,
  canAddNewLabels,
  className,
  onChange,
  onAddLabel,
}: ILabelFilterSelectProps) => {
  const [labelQuery, setLabelQuery] = useState("");
  const [shouldOpenMenu, setShouldOpenMenu] = useState(false);
  const isLabelSearchInputFocusedRef = useRef(false);
  const selectRef = useRef<
    SelectInstance<ILabel | IEmptyOption, false, IGroupOption>
  >(null);

  const options = useMemo(() => createDropdownOptions(labels, labelQuery), [
    labels,
    labelQuery,
  ]);

  const handleChange = (option: ILabel | IEmptyOption | null) => {
    if (option === null) return;
    if ("type" in option) {
      setShouldOpenMenu(false);
      setLabelQuery("");
      selectRef.current?.blur();
      onChange(option);
    }
  };

  const handleLabelQueryChange = (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    event.stopPropagation();
    setLabelQuery(event.target.value);
  };

  const handleBlurSelect = () => {
    if (!isLabelSearchInputFocusedRef.current) {
      isLabelSearchInputFocusedRef.current = false;
      setShouldOpenMenu(false);
    }
  };

  const handleFocusSelect = () => {
    setShouldOpenMenu(true);
  };

  const handleClickLabelSearchInput = () => {
    isLabelSearchInputFocusedRef.current = true;
  };

  const handleBlurLabelSearchInput = () => {
    isLabelSearchInputFocusedRef.current = false;
    setShouldOpenMenu(false);
  };

  const getOptionLabel = (option: ILabel | IEmptyOption) => {
    if ("display_text" in option) {
      return option.display_text;
    }
    return option.label;
  };

  const getOptionValue = (option: ILabel | IEmptyOption) => {
    if ("id" in option) {
      return option.id.toString();
    }
    return option.label;
  };

  const classes = classnames(baseClass, className);

  return (
    <Select<ILabel | IEmptyOption, false, IGroupOption>
      ref={selectRef}
      options={options}
      className={classes}
      classNamePrefix={baseClass}
      defaultMenuIsOpen={false}
      placeholder={"Filter by operating System or Label..."}
      formatOptionLabel={OptionLabel}
      menuIsOpen={shouldOpenMenu}
      value={selectedLabel}
      isSearchable={false}
      getOptionLabel={getOptionLabel}
      getOptionValue={getOptionValue}
      components={{ GroupHeading: CustomLabelGroupHeading }}
      labelQuery={labelQuery}
      canAddNewLabels={canAddNewLabels}
      onChange={handleChange}
      onBlur={handleBlurSelect}
      onFocus={handleFocusSelect}
      onAddLabel={onAddLabel}
      onChangeLabelQuery={handleLabelQueryChange}
      onClickLabelSeachInput={handleClickLabelSearchInput}
      onBlurLabelSearchInput={handleBlurLabelSearchInput}
    />
  );
};

export default LabelFilterSelect;

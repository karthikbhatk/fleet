import Button from "components/buttons/Button";
import { ILabel } from "interfaces/label";
import React, { useRef } from "react";
import { components, GroupHeadingProps } from "react-select-5";

import { IEmptyOption, IGroupOption } from "../LabelFilterSelect/helpers";
import PlusIcon from "../../../../../../assets/images/icon-plus-16x16@2x.png";

const baseClass = "custom-label-group-heading";

const CustomLabelGroupHeading = (
  props: GroupHeadingProps<ILabel | IEmptyOption, false, IGroupOption>
) => {
  const { data, selectProps } = props;
  const {
    labelQuery,
    canAddNewLabels,
    onAddLabel,
    onChangeLabelQuery,
    onClickLabelSeachInput,
    onBlurLabelSearchInput,
  } = selectProps;
  const inputRef = useRef<HTMLInputElement | null>(null);

  const handleInputClick = (
    event: React.MouseEvent<HTMLInputElement, MouseEvent>
  ) => {
    onClickLabelSeachInput(event);
    inputRef.current?.focus();
    event.stopPropagation();
  };

  return data.type === "platform" ? (
    <components.GroupHeading {...props}>
      <div className={`${baseClass}__labels-header`}>
        <span className={`${baseClass}__label-title`}>{props.children}</span>
      </div>
    </components.GroupHeading>
  ) : (
    <components.GroupHeading {...props}>
      <div className={`${baseClass}__labels-header`}>
        <span className={`${baseClass}__label-title`}>{props.children}</span>
        <div className={`${baseClass}__add_new_label`}>
          {canAddNewLabels && (
            <Button
              variant="text-icon"
              onClick={onAddLabel}
              className={`${baseClass}__add-label-btn`}
            >
              <>
                <span>Add label</span>
                <img src={PlusIcon} alt="Add label icon" />
              </>
            </Button>
          )}
        </div>
      </div>
      <div className={`${baseClass}__field`}>
        <input
          className={`${baseClass}__input`}
          ref={inputRef}
          value={labelQuery}
          placeholder="Filter labels by name..."
          onKeyDown={(event) => {
            // Stops the parent dropdown from picking up on input keypresses
            event.stopPropagation();
          }}
          onChange={onChangeLabelQuery}
          onClick={handleInputClick}
          onBlur={onBlurLabelSearchInput}
        />
      </div>
    </components.GroupHeading>
  );
};

export default CustomLabelGroupHeading;

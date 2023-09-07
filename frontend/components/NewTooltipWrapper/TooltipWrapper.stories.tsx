import React from "react";
import { Meta, Story } from "@storybook/react";

import NewTooltipWrapper from ".";

import "../../index.scss";

interface INewTooltipWrapperProps {
  children: string;
  tipContent: string;
}

export default {
  component: NewTooltipWrapper,
  title: "Components/NewTooltipWrapper",
  args: {
    tipContent: "This is an example tooltip.",
  },
  argTypes: {
    position: {
      options: ["top", "bottom"],
      control: "radio",
    },
  },
} as Meta;

// using line breaks to create space for top position
const Template: Story<INewTooltipWrapperProps> = (props) => (
  <>
    <br />
    <br />
    <br />
    <br />
    <NewTooltipWrapper {...props}>Example text</NewTooltipWrapper>
  </>
);

export const Default = Template.bind({});

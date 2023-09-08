import React from "react";
import { fireEvent, render, screen } from "@testing-library/react";
import { renderWithSetup } from "test/test-utils";

import HumanTimeDiffWithDateTip from "./HumanTimeDiffWithDateTip";

const EMPTY_STRING = "Unavailable";
const INVALID_STRING = "Invalid date";

describe("HumanTimeDiffWithDateTip - component", () => {
  it("renders tooltip on hover", async () => {
    const { user } = renderWithSetup(
      <HumanTimeDiffWithDateTip timeString="2015-12-06T10:30:00Z" />
    );

    // Note: text varies as time passes
    await user.hover(screen.getByText(/years ago/i));

    // Note: text varies for timezones
    expect(screen.getByText(/12\/6\/2020/i)).toBeInTheDocument();
  });

  it("handles empty string error", async () => {
    render(<HumanTimeDiffWithDateTip timeString="" />);

    const emptyStringText = screen.getByText(EMPTY_STRING);
    expect(emptyStringText).toBeInTheDocument();
  });

  it("handles invalid string error", async () => {
    render(<HumanTimeDiffWithDateTip timeString="foobar" />);

    const invalidStringText = screen.getByText(INVALID_STRING);
    expect(invalidStringText).toBeInTheDocument();
  });
});

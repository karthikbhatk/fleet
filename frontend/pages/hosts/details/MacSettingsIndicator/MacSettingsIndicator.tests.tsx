import React from "react";
import { fireEvent, render, screen } from "@testing-library/react";
import MacSettingsIndicator from "./MacSettingsIndicator";

describe("MacSettingsIndicator", () => {
  it("Renders the text and icon", () => {
    const indicatorText = "test text";
    render(
      <MacSettingsIndicator indicatorText={indicatorText} iconName="success" />
    );
    const renderedIndicatorText = screen.getByText(indicatorText);
    const renderedIcon = screen.getByTestId("icon");

    expect(renderedIndicatorText).toBeInTheDocument();
    expect(renderedIcon).toBeInTheDocument();
  });

  it("Renders text, icon, and tooltip", () => {
    const indicatorText = "test text";
    const tooltipText = "test tooltip text";
    render(
      <MacSettingsIndicator
        indicatorText={indicatorText}
        iconName="success"
        tooltip={{ tooltipText }}
      />
    );
    const renderedIndicatorText = screen.getByText(indicatorText);
    const renderedIcon = screen.getByTestId("icon");
    const renderedTooltipText = screen.getByText(tooltipText);

    expect(renderedIndicatorText).toBeInTheDocument();
    expect(renderedIcon).toBeInTheDocument();
    expect(renderedTooltipText).toBeInTheDocument();
  });

  it("Renders text, icon, and onClick", () => {
    const indicatorText = "test text";
    const onClick = () => {
      const newDiv = document.createElement("div");
      newDiv.appendChild(document.createTextNode("onClick called"));
      document.body.appendChild(newDiv);
    };
    render(
      <MacSettingsIndicator
        indicatorText={indicatorText}
        iconName="success"
        onClick={() => {
          onClick();
        }}
      />
    );

    const renderedIndicatorText = screen.getByText(indicatorText);
    const renderedIcon = screen.getByTestId("icon");
    const renderedButton = screen.getByRole("button");

    expect(renderedIndicatorText).toBeInTheDocument();
    expect(renderedIcon).toBeInTheDocument();
    expect(renderedButton).toBeInTheDocument();

    fireEvent.click(renderedButton);
    expect(screen.getByText("onClick called")).toBeInTheDocument();
  });

  it("Renders text, icon, tooltip and onClick", () => {
    const indicatorText = "test text";
    const tooltipText = "test tooltip text";
    const onClick = () => {
      const newDiv = document.createElement("div");
      newDiv.appendChild(document.createTextNode("onClick called"));
      document.body.appendChild(newDiv);
    };
    render(
      <MacSettingsIndicator
        indicatorText={indicatorText}
        iconName="success"
        onClick={() => {
          onClick();
        }}
        tooltip={{ tooltipText }}
      />
    );

    const renderedIndicatorText = screen.getByText(indicatorText);
    const renderedIcon = screen.getByTestId("icon");
    const renderedButton = screen.getByRole("button");
    const renderedTooltipText = screen.getByText(tooltipText);

    expect(renderedIndicatorText).toBeInTheDocument();
    expect(renderedIcon).toBeInTheDocument();
    expect(renderedButton).toBeInTheDocument();
    expect(renderedTooltipText).toBeInTheDocument();

    fireEvent.click(renderedButton);
    expect(screen.getByText("onClick called")).toBeInTheDocument();
  });
});

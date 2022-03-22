import { render, screen } from "@testing-library/react";

import helpers from "test/helpers";
import { userStub } from "test/stubs";
import CoreLayout from "./CoreLayout";

const { connectedComponent, reduxMockStore } = helpers;

describe("CoreLayout - layouts", () => {
  const store = {
    app: { config: {} },
    auth: { user: userStub },
    notifications: {},
    persistentFlash: {
      showFlash: false,
      message: "",
    },
  };
  const mockStore = reduxMockStore(store);

  it("renders the FlashMessage component when notifications are present", () => {
    const storeWithNotifications = {
      app: { config: {} },
      auth: {
        user: userStub,
      },
      notifications: {
        alertType: "success",
        isVisible: true,
        message: "nice jerb!",
      },
    };
    const mockStoreWithNotifications = reduxMockStore(storeWithNotifications);
    const componentWithFlash = connectedComponent(CoreLayout, {
      mockStore: mockStoreWithNotifications,
    });
    const componentWithoutFlash = connectedComponent(CoreLayout, {
      mockStore,
    });

    const { container, rerender } = render(componentWithFlash);

    expect(container.querySelector(".flash-message")).toBeInTheDocument();

    rerender(componentWithoutFlash);

    expect(container.querySelector(".flash-message")).not.toBeInTheDocument();
  });
});

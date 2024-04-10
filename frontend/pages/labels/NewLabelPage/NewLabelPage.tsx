import React, { useCallback, useContext, useState } from "react";
import { Tab, TabList, Tabs } from "react-tabs";

import { QueryContext } from "context/query";
import useToggleSidePanel from "hooks/useToggleSidePanel";

import MainContent from "components/MainContent";
import SidePanelContent from "components/SidePanelContent";
import TabsWrapper from "components/TabsWrapper";
import QuerySidePanel from "components/side_panels/QuerySidePanel";

import PATHS from "router/paths";
import { InjectedRouter, RouteComponentProps } from "react-router";

interface ILabelSubNavItem {
  name: string;
  pathname: string;
}

const labelSubNav: ILabelSubNavItem[] = [
  {
    name: "Dynamic",
    pathname: PATHS.LABEL_NEW_DYNAMIC,
  },
  {
    name: "Manual",
    pathname: PATHS.LABEL_NEW_MANUAL,
  },
];

const getTabIndex = (path: string): number => {
  return labelSubNav.findIndex((navItem) => {
    // tab stays highlighted for paths that start with same pathname
    return path.startsWith(navItem.pathname);
  });
};

const baseClass = "new-label-page";

interface INewLabelPageProps extends RouteComponentProps<never, never> {
  children: JSX.Element;
}

const NewLabelPage = ({ router, children }: INewLabelPageProps) => {
  const { selectedOsqueryTable, setSelectedOsqueryTable } = useContext(
    QueryContext
  );

  const { isSidePanelOpen, setSidePanelOpen } = useToggleSidePanel(true);
  const [showOpenSchemaActionText, setShowOpenSchemaActionText] = useState(
    false
  );

  const navigateToNav = useCallback(
    (i: number): void => {
      router.replace(labelSubNav[i].pathname);
    },
    [router]
  );

  const onCloseSchemaSidebar = () => {
    setSidePanelOpen(false);
    setShowOpenSchemaActionText(true);
  };

  const onOsqueryTableSelect = (tableName: string) => {
    setSelectedOsqueryTable(tableName);
  };

  return (
    <>
      <MainContent className={baseClass}>
        <h1>Add label</h1>
        <TabsWrapper>
          <Tabs
            selectedIndex={getTabIndex(location?.pathname || "")}
            onSelect={navigateToNav}
          >
            <TabList>
              {labelSubNav.map((navItem) => {
                return (
                  <Tab key={navItem.name} data-text={navItem.name}>
                    {navItem.name}
                  </Tab>
                );
              })}
            </TabList>
          </Tabs>
        </TabsWrapper>
        {React.cloneElement(children, {
          test: "test",
        })}
      </MainContent>
      {isSidePanelOpen && (
        <SidePanelContent>
          <QuerySidePanel
            key="query-side-panel"
            onOsqueryTableSelect={onOsqueryTableSelect}
            selectedOsqueryTable={selectedOsqueryTable}
            onClose={onCloseSchemaSidebar}
          />
        </SidePanelContent>
      )}
    </>
  );
};

export default NewLabelPage;

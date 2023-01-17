import PATHS from "router/paths";
import URL_PREFIX from "router/url_prefix";
import { IUser } from "interfaces/user";

export interface INavItem {
  icon?: string;
  name: string;
  iconName?: string;
  location: {
    regex: RegExp;
    pathname: string;
  };
  withContext?: boolean;
  exclude?: boolean;
}

export default (
  user: IUser | null,
  isGlobalAdmin = false,
  isAnyTeamAdmin = false,
  isAnyTeamMaintainer = false,
  isGlobalMaintainer = false,
  isNoAccess = false
): INavItem[] => {
  if (!user) {
    return [];
  }

  const isMaintainerOrAdmin =
    isGlobalMaintainer ||
    isAnyTeamMaintainer ||
    isGlobalAdmin ||
    isAnyTeamAdmin;

  const logo = [
    {
      icon: "logo",
      name: "Home",
      iconName: "logo",
      location: {
        regex: new RegExp(`^${URL_PREFIX}/dashboard`),
        pathname: PATHS.DASHBOARD,
      },
    },
  ];

  const navItems = [
    {
      icon: "hosts",
      name: "Hosts",
      iconName: "hosts",
      location: {
        regex: new RegExp(`^${URL_PREFIX}/hosts/`),
        pathname: PATHS.MANAGE_HOSTS,
      },
      withContext: true,
    },
    {
      name: "Controls",
      location: {
        regex: new RegExp(`^${URL_PREFIX}/controls/`),
        pathname: PATHS.CONTROLS,
      },
      exclude: !isMaintainerOrAdmin,
    },
    {
      icon: "software",
      name: "Software",
      iconName: "software",
      location: {
        regex: new RegExp(`^${URL_PREFIX}/software/`),
        pathname: PATHS.MANAGE_SOFTWARE,
      },
      withContext: true,
    },
    {
      icon: "query",
      name: "Queries",
      iconName: "queries",
      location: {
        regex: new RegExp(`^${URL_PREFIX}/queries/`),
        pathname: PATHS.MANAGE_QUERIES,
      },
    },
    {
      icon: "packs",
      name: "Schedule",
      iconName: "packs",
      location: {
        regex: new RegExp(`^${URL_PREFIX}/(schedule|packs)/`),
        pathname: PATHS.MANAGE_SCHEDULE,
      },
      exclude: !isMaintainerOrAdmin,
    },
    {
      icon: "policies",
      name: "Policies",
      iconName: "policies",
      location: {
        regex: new RegExp(`^${URL_PREFIX}/(policies)/`),
        pathname: PATHS.MANAGE_POLICIES,
      },
    },
  ];

  if (isNoAccess) {
    return [...logo];
  }

  return [
    ...logo,
    ...navItems.filter((item) => {
      return !item.exclude;
    }),
  ];
};

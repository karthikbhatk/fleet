import React, { useCallback } from "react";
import { Async, OnChangeHandler, Option } from "react-select";
import classnames from "classnames";

// ignore TS error for now until these are rewritten in ts.
// @ts-ignore
import local from "utilities/local";
// @ts-ignore
import debounce from "utilities/debounce";
import permissionUtils from "utilities/permissions";
import { IDropdownOption } from "interfaces/dropdownOption";
import { ITeam } from "../../../../interfaces/team";
import { IUser } from "../../../../interfaces/user";

const baseClass = "autocomplete-dropdown";

interface IAutocompleteDropdownProps {
  id: string;
  team: ITeam;
  placeholder: string;
  onChange: OnChangeHandler;
  resourceUrl: string;
  value: Option[];
  disabledOptions: number[];
  disabled?: boolean;
  className?: string;
}

const debounceOptions = {
  timeout: 300,
  leading: false,
  trailing: true,
};

const createUrl = (baseUrl: string, input: string) => {
  return `/api${baseUrl}?query=${input}`;
};

const generateOptionLabel = (user: IUser, team: ITeam): string => {
  const userTeamIds = user.teams.map((currentTeam) => currentTeam.id);
  if (permissionUtils.isOnGlobalTeam(user)) {
    return `${user.name} - Global user`;
    // User is already in this team
  } else if (userTeamIds.includes(team.id)) {
    const teamName = user.teams.find(
      (currentTeam) => currentTeam.id === team.id
    )?.name;
    return `${user.name} - Already a member of ${teamName}`;
  }

  return user.name;
};

const AutocompleteDropdown = (
  props: IAutocompleteDropdownProps
): JSX.Element => {
  const {
    className,
    disabled,
    disabledOptions,
    placeholder,
    onChange,
    id,
    resourceUrl,
    value,
    team,
  } = props;

  const wrapperClass = classnames(baseClass, className);

  // We disable any filtering client side as the server filters the results
  // for us.
  const filterOptions = useCallback((options) => {
    return options;
  }, []);

  const createDropdownOptions = (users: IUser[]): IDropdownOption[] => {
    return users.map((user) => {
      return {
        value: user.id,
        label: generateOptionLabel(user, team),
        disabled:
          disabledOptions.includes(user.id) || user.global_role !== null,
      };
    });
  };

  // NOTE: It seems react-select v1 Async component does not work well with
  // returning results from promises in its loadOptions handler. That is why
  // we have decided to use callbacks as those seemed to make the component work
  // More info is here:
  // https://stackoverflow.com/questions/52984105/react-select-async-loadoptions-is-not-loading-options-properly
  const getOptions = debounce((input: string, callback: any) => {
    if (!input) {
      return callback([]);
    }

    fetch(createUrl(resourceUrl, input), {
      headers: {
        authorization: `Bearer ${local.getItem("auth_token")}`,
      },
    })
      .then((res) => {
        return res.json();
      })
      .then((json) => {
        const optionsData = createDropdownOptions(json.users);
        callback(null, { options: optionsData });
      })
      .catch((err) => {
        console.log("There was an error", err);
      });
  }, debounceOptions);

  return (
    <div className={wrapperClass}>
      <Async
        noResultsText={"Nothing found"}
        autoload={false}
        cache={false}
        id={id}
        loadOptions={getOptions}
        disabled={disabled}
        placeholder={placeholder}
        onChange={onChange}
        value={value}
        filterOptions={filterOptions}
        multi
        searchable
      />
    </div>
  );
};

export default AutocompleteDropdown;

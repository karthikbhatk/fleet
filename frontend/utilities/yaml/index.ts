import yaml from "js-yaml";

interface IYAMLError {
  name: string;
  reason: string;
  line: string;
}

export const constructErrorString = (yamlError: IYAMLError) => {
  return `${yamlError.name}: ${yamlError.reason} at line ${yamlError.line}`;
};

export const agentOptionsToYaml = (agentOpts) => {
  // hide the "overrides" key if it is empty
  if (!agentOpts.overrides || Object.keys(agentOpts.overrides).length === 0) {
    delete agentOpts.overrides;
  }

  // add a comment besides the "command_line_flags" if it is empty
  let addFlagsComment = false;
  if (
    !agentOpts.command_line_flags ||
    Object.keys(agentOpts.command_line_flags).length === 0
  ) {
    agentOpts.command_line_flags = {};
    addFlagsComment = true;
  }

  let yamlString = yaml.dump(agentOpts);
  if (addFlagsComment) {
    yamlString = yamlString.replace(
      "command_line_flags: {}\n",
      "command_line_flags: {} # requires Orbit\n"
    );
  }

  return yamlString;
};

export default constructErrorString;

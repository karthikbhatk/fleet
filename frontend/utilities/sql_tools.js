import sqliteParser from "sqlite-parser";
import { intersection } from "lodash";
import { osqueryTables } from "utilities/osquery_tables";

const platformsByTableDictionary = osqueryTables.reduce(
  (dictionary, osqueryTable) => {
    dictionary[osqueryTable.name] = osqueryTable.platforms;
    return dictionary;
  },
  {}
);

// The isNode and visit functionality is informed by https://lihautan.com/manipulating-ast-with-javascript/#traversing-an-ast
const _isNode = (node) => {
  // TODO: Improve type checking against shape of AST generated by sqliteParser
  return typeof node === "object";
};
const _visit = (abstractSyntaxTree, callback) => {
  if (abstractSyntaxTree) {
    callback(abstractSyntaxTree);

    Object.keys(abstractSyntaxTree).forEach((key) => {
      const childNode = abstractSyntaxTree[key];
      if (Array.isArray(childNode)) {
        childNode.forEach((grandchildNode) => _visit(grandchildNode, callback));
      } else if (_isNode(childNode)) {
        _visit(childNode, callback);
      }
    });
  }
};

export const listCompatiblePlatforms = (tablesList) => {
  if (
    tablesList[0] === "Invalid query" ||
    tablesList[0] === "No tables in query AST"
  ) {
    return tablesList;
  }
  const compatiblePlatforms = intersection(
    ...tablesList?.map((tableName) => platformsByTableDictionary[tableName])
  );

  return compatiblePlatforms.length ? compatiblePlatforms : ["None"];
};

export const parseSqlTables = (sqlString) => {
  let results = [];

  // Tables defined via common table expression will be excluded from results by default
  const cteTables = [];

  const _callback = (node) => {
    if (node) {
      if (
        (node.variant === "common" || node.variant === "recursive") &&
        node.format === "table" &&
        node.type === "expression"
      ) {
        cteTables.push(node.target?.name);
      } else if (node.variant === "table") {
        results.push(node.name);
      }
    }
  };

  try {
    const sqlTree = sqliteParser(sqlString);
    _visit(sqlTree, _callback);

    if (cteTables.length) {
      results = results.filter((t) => !cteTables.includes(t));
    }

    return results.length ? results : ["No tables in query AST"];
  } catch (err) {
    // console.log(`Invalid query syntax: ${err.message}\n\n${sqlString}`);

    return ["Invalid query"];
  }
};

export default { listCompatiblePlatforms, parseSqlTables };

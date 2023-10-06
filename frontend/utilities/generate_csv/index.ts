import convertToCSV from "utilities/convert_to_csv";
import { Row, Column } from "react-table";
import { ICampaignError } from "interfaces/campaign";
import { format } from "date-fns";

const reorderCSVFields = (tableHeaders: string[]) => {
  console.log("tableHeaders", tableHeaders);
  const result = tableHeaders.filter((field) => field !== "host_display_name");
  result.unshift("host_display_name");

  return result;
};

export const generateCSVFilename = (descriptor: string) => {
  return `${descriptor} (${format(new Date(), "MM-dd-yy hh-mm-ss")}).csv`;
};

// Query report
export const generateCSVQueryReport = (
  rows: Row[],
  filename: string,
  tableHeaders: Column[] | string[]
) => {
  return new global.window.File(
    [
      convertToCSV({
        objArray: rows,
        fieldSortFunc: reorderCSVFields,
        tableHeaders,
      }),
    ],
    filename,
    {
      type: "text/csv",
    }
  );
};

// Query results and query errors
export const generateCSVQueryResults = (
  rows: Row[],
  filename: string,
  tableHeaders: Column[] | string[]
) => {
  return new global.window.File(
    [
      convertToCSV({
        objArray: rows.map((r) => r.original),
        fieldSortFunc: reorderCSVFields,
        tableHeaders,
      }),
    ],
    filename,
    {
      type: "text/csv",
    }
  );
};

// Policy results only
export const generateCSVPolicyResults = (
  rows: { host: string; status: string }[],
  filename: string
) => {
  return new global.window.File([convertToCSV({ objArray: rows })], filename, {
    type: "text/csv",
  });
};

// Policy errors only
export const generateCSVPolicyErrors = (
  rows: ICampaignError[],
  filename: string
) => {
  return new global.window.File([convertToCSV({ objArray: rows })], filename, {
    type: "text/csv",
  });
};

export default {
  generateCSVFilename,
  generateCSVQueryResults,
  generateCSVPolicyResults,
  generateCSVPolicyErrors,
};

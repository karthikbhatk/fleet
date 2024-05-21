import { formatDistanceToNow } from "date-fns";

// eslint-disable-next-line import/prefer-default-export
export const uploadedFromNow = (date: string) => {
  // NOTE: Malformed dates will result in errors. This is expected "fail loudly" behavior.
  return `Uploaded ${formatDistanceToNow(new Date(date))} ago`;
};

export const dateAgo = (date: string) => {
  // NOTE: Malformed dates will result in errors. This is expected "fail loudly" behavior.
  return `${formatDistanceToNow(new Date(date))} ago`;
};

import React from "react";
import { useDispatch } from "react-redux";
import { push } from "react-router-redux";

import Button from "components/buttons/Button/Button";

interface ILinkCellProps<T> {
  value: string;
  path: string;
  title?: string;
}

const LinkCell = (props: ILinkCellProps<any>): JSX.Element => {
  const { value, path, title } = props;

  const dispatch = useDispatch();

  const onClick = (): void => {
    dispatch(push(path));
  };

  return (
    <Button onClick={onClick} variant="text-link" title={title}>
      {value}
    </Button>
  );
};

export default LinkCell;

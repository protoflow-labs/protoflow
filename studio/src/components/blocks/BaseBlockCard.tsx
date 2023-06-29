import {
  Card, CardProps,
  makeStyles,
  mergeClasses,
  shorthands,
  tokens,
} from "@fluentui/react-components";
import React, { ReactNode } from "react";

export interface BaseBlockCardProps extends CardProps {
  children: ReactNode;
  selected: boolean;
}

const useOverrides = makeStyles({
  card: {
    ...shorthands.border("1px", "solid", tokens.colorNeutralBackground1),
  },
  selected: {
    ...shorthands.border(
      "1px",
      "solid",
      tokens.colorNeutralForeground2BrandSelected
    ),
  },
});

export const BaseBlockCard: React.FC<BaseBlockCardProps> = ({ selected, ...props }) => {
  const overrides = useOverrides();
  const classes = mergeClasses(overrides.card, selected && overrides.selected);

  return <Card {...props} className={classes} />;
}

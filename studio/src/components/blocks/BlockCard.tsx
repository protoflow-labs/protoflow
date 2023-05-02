import {
  Card,
  makeStyles,
  mergeClasses,
  shorthands,
  tokens,
} from "@fluentui/react-components";
import { ReactNode } from "react";

type BlockCardProps = {
  children: ReactNode;
  selected: boolean;
};

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
export function BlockCard({ selected, ...props }: BlockCardProps) {
  const overrides = useOverrides();
  const classes = mergeClasses(overrides.card, selected && overrides.selected);

  return <Card {...props} className={classes} />;
}

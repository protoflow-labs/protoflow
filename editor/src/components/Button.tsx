import { MouseEvent, ReactNode } from "react";
import { tw } from "../lib/tw";

type ButtonProps = {
  className?: string;
  children?: ReactNode | ReactNode[];
  onClick?: (e: MouseEvent<HTMLButtonElement>) => void;
};
export function Button(props: ButtonProps) {
  const { children, className, onClick } = props;
  return (
    <button
      type="button"
      className={tw(
        "rounded bg-white/10 px-2 py-1 text-xs font-semibold text-white shadow-sm hover:bg-white/20",
        className
      )}
      onClick={onClick}
    >
      {children}
    </button>
  );
}

import { Caption1, Card, CardHeader, Text } from "@fluentui/react-components";
import { ChangeEvent, KeyboardEvent, useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { HiPencilSquare } from "react-icons/hi2";
import { Handle, NodeProps, Position } from "reactflow";

const handleStyle = {
  background: "#555",
};

type InputNodeProps = NodeProps<InputData>;

type Field = {
  name: string;
  type: "string" | "number" | "boolean";
  required?: boolean;
};

export type InputData = {
  name: string;
  fields: Field[];
  lastUpdate: number;
};

type InputProps = {
  defaultValue?: string;
  onKeyUp?: (e: KeyboardEvent<HTMLInputElement>) => void;
  onChange?: (e: ChangeEvent<HTMLInputElement>) => void;
};

export function InputNode(props: InputNodeProps) {
  const { data, selected } = props;

  return (
    <>
      <Card>
        <CardHeader
          image={<HiPencilSquare className="h-5 w-5 bg-gray-800" />}
          header={<Text weight="semibold">{data.name || "UntitledInput"}</Text>}
          description={<Caption1>Input</Caption1>}
        />
      </Card>
      <Handle type="source" position={Position.Bottom} />
      <Handle type="source" position={Position.Bottom} id="b" />
    </>
  );
}

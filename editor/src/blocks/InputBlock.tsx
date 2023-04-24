import { Caption1, Card, CardHeader, Text } from "@fluentui/react-components";
import { ChangeEvent, KeyboardEvent } from "react";
import { HiPencilSquare } from "react-icons/hi2";
import { Handle, NodeProps, Position } from "reactflow";
import { FieldDefinition } from "../../rpc/block_pb";

type InputBlockProps = NodeProps<InputData>;

export type InputData = {
  name: string;
  fields: FieldDefinition[];
  lastUpdate: number;
};

type InputProps = {
  defaultValue?: string;
  onKeyUp?: (e: KeyboardEvent<HTMLInputElement>) => void;
  onChange?: (e: ChangeEvent<HTMLInputElement>) => void;
};

export function InputBlock(props: InputBlockProps) {
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

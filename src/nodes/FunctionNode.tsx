import { useCallback } from "react";
import { Handle, NodeProps, Position } from "reactflow";
import { HiCircleStack, HiCodeBracket } from "react-icons/hi2";
import { Caption1, Card, CardHeader, Text } from "@fluentui/react-components";

const handleStyle = {
  background: "blue",
};
type FunctionNodeProps = NodeProps<FunctionData>;

export type FunctionData = {
  name: string;
  language: string;
};

export function FunctionNode(props: FunctionNodeProps) {
  const { data, selected } = props;

  return (
    <>
      <Card>
        <CardHeader
          image={<HiCodeBracket className="h-5 w-5 bg-gray-800" />}
          header={
            <Text weight="semibold">{data.name || "Untitled Function"}</Text>
          }
          description={<Caption1>Function</Caption1>}
        />
      </Card>
      <Handle type="source" position={Position.Bottom} className="z-10" />
      <Handle type="target" position={Position.Top} className="z-10" />
    </>
  );
}

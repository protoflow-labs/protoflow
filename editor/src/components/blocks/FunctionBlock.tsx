import { Caption1, CardHeader, Text } from "@fluentui/react-components";
import { HiCodeBracket } from "react-icons/hi2";
import { Handle, NodeProps, Position } from "reactflow";
import { BlockCard } from "./BlockCard";

type FunctionBlockProps = NodeProps<FunctionData>;

export type FunctionData = {
  name: string;
  config: {
    function?: {
      runtime: string;
    };
  };
};

export function FunctionBlock(props: FunctionBlockProps) {
  const { data, selected } = props;

  return (
    <>
      <BlockCard selected={selected}>
        <CardHeader
          image={<HiCodeBracket className="h-5 w-5 bg-gray-800" />}
          header={
            <Text weight="semibold">{data.name || "Untitled Function"}</Text>
          }
          description={<Caption1>Function</Caption1>}
        />
      </BlockCard>
      <Handle type="source" position={Position.Bottom} className="z-10" />
      <Handle type="target" position={Position.Top} className="z-10" />
    </>
  );
}
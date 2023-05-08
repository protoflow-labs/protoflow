import { Caption1, CardHeader, Text } from "@fluentui/react-components";
import {MdCode, MdHttp, MdOutbound} from "react-icons/md";
import { Handle, NodeProps, Position } from "reactflow";
import { BlockCard } from "./BlockCard";

export type GRPCBlockProps = NodeProps<GRPCData>;

export type GRPCData = {
  name: string;
  config: {
    grpc: { package: string; service: string; method: string };
  };
};

export function GRPCBlock(props: GRPCBlockProps) {
  const { data, selected } = props;

  return (
    <>
      <BlockCard selected={selected}>
        <CardHeader
          image={<MdOutbound className="h-5 w-5 bg-gray-800" />}
          header={<Text weight="semibold">{data.name || "Untitled GRPC"}</Text>}
          description={<Caption1>GRPC</Caption1>}
        />
      </BlockCard>
      <Handle type="source" position={Position.Bottom} className="z-10" />
      <Handle type="target" position={Position.Top} className="z-10" />
    </>
  );
}

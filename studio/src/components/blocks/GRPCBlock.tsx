import {Badge, Caption1, CardFooter, CardHeader, Text} from "@fluentui/react-components";
import {MdCode, MdHttp, MdOutbound} from "react-icons/md";
import { Handle, NodeProps, Position } from "reactflow";
import { BlockCard } from "./BlockCard";
import {GRPC} from "@/rpc/block_pb";
import {useProjectContext} from "@/providers/ProjectProvider";
import React from "react";

export type GRPCBlockProps = NodeProps<GRPCData>;

export type GRPCData = {
  name: string;
  config: {
    grpc: GRPC
  };
  resourceIds: string[];
};

export function GRPCBlock(props: GRPCBlockProps) {
  const { data, selected } = props;
  const {resources} = useProjectContext();

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
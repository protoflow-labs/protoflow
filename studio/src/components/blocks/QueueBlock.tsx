import { Caption1, CardHeader, Text } from "@fluentui/react-components";
import { AiOutlineMail } from "react-icons/ai";
import { Handle, NodeProps, Position } from "reactflow";
import { BaseBlockCard } from "./BaseBlockCard";

export type QueueBlockProps = NodeProps<QueueData>;

export type QueueData = {
  name: string;
  database: string;
  variables: Record<string, string>;
  sql: string;
};

export function QueueBlock(props: QueueBlockProps) {
  const { data, selected } = props;
  return (
    <>
      <BaseBlockCard selected={selected}>
        <CardHeader
          image={<AiOutlineMail className="h-5 w-5 bg-gray-800" />}
          header={
            <Text weight="semibold">{data.name || "Untitled Queue"}</Text>
          }
          description={<Caption1>Queue</Caption1>}
        />
      </BaseBlockCard>
      <Handle type="source" position={Position.Bottom} className="z-10" />
      <Handle type="target" position={Position.Top} className="z-10" />
    </>
  );
}

import { Caption1, CardHeader, Text } from "@fluentui/react-components";
import { HiOutlineMagnifyingGlass } from "react-icons/hi2";
import { Handle, NodeProps, Position } from "reactflow";
import { BaseBlockCard } from "./BaseBlockCard";

export type QueryBlockProps = NodeProps<QueryData>;

export type QueryData = {
  name: string;
  database: string;
  variables: Record<string, string>;
  sql: string;
};

export function QueryBlock(props: QueryBlockProps) {
  const { data, selected } = props;
  return (
    <>
      <BaseBlockCard selected={selected}>
        <CardHeader
          image={<HiOutlineMagnifyingGlass className="h-5 w-5 bg-gray-800" />}
          header={
            <Text weight="semibold">{data.name || "Untitled Query"}</Text>
          }
          description={<Caption1>Query</Caption1>}
        />
      </BaseBlockCard>
      <Handle type="source" position={Position.Bottom} className="z-10" />
      <Handle type="target" position={Position.Top} className="z-10" />
    </>
  );
}

import { Caption1, CardHeader, Text } from "@fluentui/react-components";
import { TbBucket } from "react-icons/tb";
import { Handle, NodeProps, Position } from "reactflow";
import { BaseBlockCard } from "./BaseBlockCard";

export type BucketBlockProps = NodeProps<BucketData>;

export type BucketData = {
  name: string;
};


export function BucketBlock(props: BucketBlockProps) {
  const { data, selected } = props;
  return (
    <>
      <BaseBlockCard selected={selected}>
        <CardHeader
          image={<TbBucket className="h-5 w-5 bg-gray-800" />}
          header={
            <Text weight="semibold">{data.name || "Untitled Bucket"}</Text>
          }
          description={<Caption1>Bucket</Caption1>}
        />
      </BaseBlockCard>
      <Handle type="source" position={Position.Bottom} className="z-10" />
      <Handle type="target" position={Position.Top} className="z-10" />
    </>
  );
}

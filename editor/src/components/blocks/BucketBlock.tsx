import { Caption1, Card, CardHeader, Text } from "@fluentui/react-components";
import { TbBucket } from "react-icons/tb";
import { Handle, NodeProps, Position } from "reactflow";

export type BucketBlockProps = NodeProps<BucketData>;

export type BucketData = {
  name: string;
};

export function BucketBlock(props: BucketBlockProps) {
  const { data, selected } = props;
  return (
    <>
      <Card>
        <CardHeader
          image={<TbBucket className="h-5 w-5 bg-gray-800" />}
          header={
            <Text weight="semibold">{data.name || "Untitled Bucket"}</Text>
          }
          description={<Caption1>Bucket</Caption1>}
        />
      </Card>
      <Handle type="source" position={Position.Bottom} className="z-10" />
      <Handle type="target" position={Position.Top} className="z-10" />
    </>
  );
}

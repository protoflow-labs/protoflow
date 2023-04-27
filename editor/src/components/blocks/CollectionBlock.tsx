import { Caption1, Card, CardHeader, Text } from "@fluentui/react-components";
import { HiCircleStack } from "react-icons/hi2";
import { Handle, NodeProps, Position } from "reactflow";

export type EntityBlockProps = NodeProps<EntityData>;

export type EntityData = {
  name: string;
  config: { entity?: { collection: string } };
};

export function CollectionBlock(props: EntityBlockProps) {
  const { data, selected } = props;

  return (
    <>
      <Card>
        <CardHeader
          image={<HiCircleStack className="h-5 w-5 bg-gray-800" />}
          header={
            <Text weight="semibold">{data.name || "Untitled Collection"}</Text>
          }
          description={<Caption1>Collection</Caption1>}
        />
      </Card>
      <Handle type="source" position={Position.Bottom} className="z-10" />
      <Handle type="target" position={Position.Top} className="z-10" />
    </>
  );
}

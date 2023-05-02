import { Caption1, CardHeader, Text } from "@fluentui/react-components";
import { HiCircleStack } from "react-icons/hi2";
import { Handle, NodeProps, Position } from "reactflow";
import { BlockCard } from "./BlockCard";

export type EntityBlockProps = NodeProps<EntityData>;

export type EntityData = {
  name: string;
  config: { entity?: { collection: string } };
};

export function CollectionBlock(props: EntityBlockProps) {
  const { data, selected } = props;

  return (
    <>
      <BlockCard selected={selected}>
        <CardHeader
          image={<HiCircleStack className="h-5 w-5 bg-gray-800" />}
          header={
            <Text weight="semibold">{data.name || "Untitled Collection"}</Text>
          }
          description={<Caption1>Collection</Caption1>}
        />
      </BlockCard>
      <Handle type="source" position={Position.Bottom} className="z-10" />
      <Handle type="target" position={Position.Top} className="z-10" />
    </>
  );
}

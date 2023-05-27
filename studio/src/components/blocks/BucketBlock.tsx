import { Caption1, CardHeader, Text } from "@fluentui/react-components";
import { TbBucket } from "react-icons/tb";
import { Handle, NodeProps, Position } from "reactflow";
import { BlockCard } from "./BlockCard";
import {useProjectContext} from "@/providers/ProjectProvider";

export function BucketBlock(props: NodeProps) {
  const { id, selected } = props;
  const {nodeLookup} = useProjectContext();
  const node = nodeLookup[id];
  return (
    <>
      <BlockCard selected={selected}>
        <CardHeader
          image={<TbBucket className="h-5 w-5 bg-gray-800" />}
          header={
            <Text weight="semibold">{node.name || "Untitled Bucket"}</Text>
          }
          description={<Caption1>Bucket</Caption1>}
        />
      </BlockCard>
      <Handle type="source" position={Position.Bottom} className="z-10" />
      <Handle type="target" position={Position.Top} className="z-10" />
    </>
  );
}

import { Caption1, CardHeader, Text } from "@fluentui/react-components";
import { AiOutlineMail } from "react-icons/ai";
import { Handle, NodeProps, Position } from "reactflow";
import { BlockCard } from "./BlockCard";
import {useProjectContext} from "@/providers/ProjectProvider";

export function QueueBlock(props: NodeProps) {
  const { id, selected } = props;
  const {nodeLookup} = useProjectContext();
  const node = nodeLookup[id];
  return (
    <>
      <BlockCard selected={selected}>
        <CardHeader
          image={<AiOutlineMail className="h-5 w-5 bg-gray-800" />}
          header={
            <Text weight="semibold">{node.name || "Untitled Queue"}</Text>
          }
          description={<Caption1>Queue</Caption1>}
        />
      </BlockCard>
      <Handle type="source" position={Position.Bottom} className="z-10" />
      <Handle type="target" position={Position.Top} className="z-10" />
    </>
  );
}

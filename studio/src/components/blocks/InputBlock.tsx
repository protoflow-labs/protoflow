import { Caption1, CardHeader, Text } from "@fluentui/react-components";
import { HiPencilSquare } from "react-icons/hi2";
import { Handle, NodeProps, Position } from "reactflow";
import { BlockCard } from "./BlockCard";
import {useProjectContext} from "@/providers/ProjectProvider";

export function InputBlock(props: NodeProps) {
  const {nodeLookup} = useProjectContext();
  const { id, selected } = props;

  const node = nodeLookup[id];

  return (
    <>
      <BlockCard selected={selected}>
        <CardHeader
          image={<HiPencilSquare className="h-5 w-5 bg-gray-800" />}
          header={<Text weight="semibold">{node.name || "UntitledInput"}</Text>}
          description={<Caption1>Input</Caption1>}
        />
      </BlockCard>
      <Handle type="source" position={Position.Bottom} />
      <Handle type="source" position={Position.Bottom} id="b" />
    </>
  );
}

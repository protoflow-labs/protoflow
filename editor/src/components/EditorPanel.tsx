import { Card } from "@fluentui/react-components";
import { useState } from "react";
import { Node, useOnSelectionChange } from "reactflow";
import { FieldType } from "../../rpc/block_pb";
import { EntityEditor } from "./blockEditors/EntityEditor";
import { FunctionEditor } from "./blockEditors/FunctionEditor";
import { InputEditor } from "./blockEditors/InputEditor";

export function EditorPanel() {
  const [activeNode, setActiveNode] = useState<Node | null>(null);
  useOnSelectionChange({
    onChange: ({ nodes }) => {
      if (nodes.length !== 1) {
        setActiveNode(null);
        return;
      }

      const [node] = nodes;
      setActiveNode(node);
    },
  });

  if (!activeNode) return null;

  return (
    <div className="absolute top-0 right-0 m-4 z-10">
      <Card>
        <NodeEditor node={activeNode} />
      </Card>
    </div>
  );
}

type NodeEditorProps = {
  node: Node | null;
};

function NodeEditor(props: NodeEditorProps) {
  switch (props.node?.type) {
    case "protoflow.input":
      return <InputEditor node={props.node} />;
    case "protoflow.entity":
      return <EntityEditor node={props.node} />;
    case "protoflow.function":
      return <FunctionEditor node={props.node} />;
    default:
      return null;
  }
}

export const fieldTypeMap = {
  [FieldType.STRING]: "String",
  [FieldType.INTEGER]: "Integer",
  [FieldType.BOOLEAN]: "Boolean",
};

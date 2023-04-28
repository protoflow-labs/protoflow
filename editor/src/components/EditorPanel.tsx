import { Card } from "@fluentui/react-components";
import { useState } from "react";
import { Node, useOnSelectionChange } from "reactflow";
import { FieldType } from "../../rpc/block_pb";
import { BucketEditor } from "./blockEditors/BucketEditor";
import { CollectionEditor } from "./blockEditors/CollectionEditor";
import { FunctionEditor } from "./blockEditors/FunctionEditor";
import { InputEditor } from "./blockEditors/InputEditor";
import { QueryEditor } from "./blockEditors/QueryEditor";
import { RESTEditor } from "./blockEditors/RESTEditor";

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
  if (!props.node || !props.node.type) {
    return null;
  }

  switch (props.node.type) {
    case "protoflow.input":
      return <InputEditor node={props.node} />;
    case "protoflow.collection":
      return <CollectionEditor node={props.node} />;
    case "protoflow.query":
      return <QueryEditor node={props.node} />;
    case "protoflow.function":
      return <FunctionEditor node={props.node} />;
    case "protoflow.bucket":
      return <BucketEditor node={props.node} />;
    case "protoflow.rest":
      return <RESTEditor node={props.node} />;
    default:
      return null;
  }
}

export const fieldTypeMap = {
  [FieldType.STRING]: "String",
  [FieldType.INTEGER]: "Integer",
  [FieldType.BOOLEAN]: "Boolean",
};

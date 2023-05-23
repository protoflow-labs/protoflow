import { Card } from "@fluentui/react-components";
import { useState } from "react";
import {Edge, EdgeChange, Node, useEdgesState, useOnSelectionChange} from "reactflow";
import { BucketEditor } from "./blockEditors/BucketEditor";
import { CollectionEditor } from "./blockEditors/CollectionEditor";
import { FunctionEditor } from "./blockEditors/FunctionEditor";
import { InputEditor } from "./blockEditors/InputEditor";
import { QueryEditor } from "./blockEditors/QueryEditor";
import { RESTEditor } from "./blockEditors/RESTEditor";
import {GRPCEditor} from "@/components/blockEditors/GRPCEditor";
import NodeProvider from "@/providers/NodeProvider";

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
    <NodeProvider nodeId={activeNode.id}>
      <div className="absolute top-0 right-0 m-4 z-10 overflow-auto">
        <Card>
          <NodeEditor node={activeNode} />
        </Card>
      </div>
    </NodeProvider>
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
    case "protoflow.grpc":
      return <GRPCEditor node={props.node} />;
    default:
      return null;
  }
}

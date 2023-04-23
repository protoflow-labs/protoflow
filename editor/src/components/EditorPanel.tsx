import { Button, Card, Input, Label, Select } from "@fluentui/react-components";
import { Node, useOnSelectionChange } from "reactflow";
import { EntityData } from "../nodes/EntityNode";
import { useState } from "react";
import { HiPlus } from "react-icons/hi2";
import { FunctionData } from "../nodes/FunctionNode";
import { projectService } from "../lib/api";
import { EndpointyData } from "../nodes/EndpointNode";

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
      <Card className="w-full max-w-sm">
        <NodeEditor node={activeNode} />
      </Card>
    </div>
  );
}

type NodeEditorProps = {
  node: Node | null;
};

interface NodeWithName {
  name: string;
}

function saveChanges(node: Node<NodeWithName>) {
  projectService.updateBlock({
    block: {
      id: node.id,
      x: node.position.x,
      y: node.position.y,
      name: node.data.name,
      type: node.type,
    },
  });
}

function NodeEditor(props: NodeEditorProps) {
  switch (props.node?.type) {
    case "message":
      return <InputEditor node={props.node} />;
    case "entity":
      return <EntityEditor node={props.node} />;
    case "function":
      return <FunctionEditor node={props.node} />;
    case "endpoint":
        return <EndpointEditor node={props.node} />;  
    default:
      return null;
  }
}

function EndpointEditor(props: { node: Node<EndpointyData> }) {
  return (
    <div className="flex flex-col gap-2 p-4">
      <div className="flex flex-col">
        <Label htmlFor="inputName">Name</Label>
        <Input
          id="inputName"
          defaultValue={props.node.data.name || ""}
          onChange={(e) => {
            props.node.data.name = e.currentTarget.value;
          }}
          onBlur={() => saveChanges(props.node)}
        />
      </div>
    </div>
  );
}

function InputEditor(props: { node: Node<EntityData> }) {
  return (
    <div className="flex flex-col gap-2 p-4">
      <div className="flex flex-col">
        <Label htmlFor="inputName">Name</Label>
        <Input
          id="inputName"
          defaultValue={props.node.data.name || ""}
          onChange={(e) => {
            props.node.data.name = e.currentTarget.value;
          }}
          onBlur={() => saveChanges(props.node)}
        />
      </div>
      <div className="flex flex-col">
        <Label htmlFor="entityName">Fields</Label>
        <Button icon={<HiPlus className="w-4 h-4" />}>Add Field</Button>
      </div>
    </div>
  );
}

function EntityEditor(props: { node: Node<EntityData> }) {
  return (
    <div className="flex flex-col gap-2 p-4">
      <div className="flex flex-col">
        <Label htmlFor="entityName">Name</Label>
        <Input
          id="entityName"
          defaultValue={props.node.data.name || ""}
          onChange={(e) => {
            props.node.data.name = e.target.value;
          }}
          onBlur={() => saveChanges(props.node)}
        />
      </div>
      <div className="flex flex-col">
        <Label htmlFor="entityTable">Storage Table</Label>
        <Input
          id="entityTable"
          defaultValue={props.node.data.table}
          onChange={(e) => {
            e.currentTarget.value = e.currentTarget.value
              .replace(/\s/g, "_")
              .toLowerCase();
            props.node.data.table = e.currentTarget.value;
          }}
        />
      </div>
    </div>
  );
}

function FunctionEditor(props: { node: Node<FunctionData> }) {
  return (
    <div className="flex flex-col gap-2 p-4">
      <div className="flex flex-col">
        <Label htmlFor="entityName">Name</Label>
        <Input
          id="entityName"
          defaultValue={props.node.data.name || ""}
          onChange={(e) => {
            props.node.data.name = e.target.value;
          }}
          onBlur={() => saveChanges(props.node)}
        />
      </div>
      <div className="flex flex-col">
        <Label htmlFor="entityLanguage">Language</Label>
        <Select
          onChange={(e) => {
            props.node.data.language = e.currentTarget.value;
          }}
        >
          <option>Go</option>
          <option>Node.js</option>
          <option>Python</option>
        </Select>
      </div>
    </div>
  );
}

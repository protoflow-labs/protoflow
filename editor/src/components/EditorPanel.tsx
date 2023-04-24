import {
  Button,
  Card,
  Dropdown,
  Input,
  Label,
  Select,
  Option,
} from "@fluentui/react-components";
import { Node, useOnSelectionChange } from "reactflow";
import { EntityData } from "../nodes/EntityNode";
import { useEffect, useState } from "react";
import { HiOutlineTrash, HiPlus } from "react-icons/hi2";
import { FunctionData } from "../nodes/FunctionNode";
import { projectService } from "../lib/api";
import { EndpointyData } from "../nodes/EndpointNode";
import { InputData } from "../nodes/InputNode";
import { useForm } from "react-hook-form";
import {FieldDefinition} from "../../rpc/block_pb";

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

interface BlockNode {
  name: string;
  inputFields?: FieldDefinition[];
  outputFields?: FieldDefinition[];
}

function saveChanges(node: Node<BlockNode>) {
  projectService.updateBlock({
    block: {
      id: node.id,
      x: node.position.x,
      y: node.position.y,
      name: node.data.name,
      type: node.type,
      inputFields: node.data.inputFields,
      outputFields: node.data.outputFields,
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

function InputEditor({ node }: { node: Node<InputData> }) {
  const { watch, setValue } = useForm({
    defaultValues: {
      name: node.data.name,
      inputFields: node.data.inputFields,
    },
  });

  const values = watch();

  useEffect(() => {
    if (!node) {
      return;
    }
    node.data.name = values.name;
    node.data.inputFields = values.inputFields;
  }, [values]);

  useEffect(() => {}, [node]);

  return (
    <div className="flex flex-col gap-2 p-4">
      <div className="flex flex-col">
        <Label htmlFor="inputName">Name</Label>
        <Input
          id="inputName"
          defaultValue={values.name || ""}
          onChange={(e) => {
            setValue("name", e.currentTarget.value);
          }}
          onBlur={() => saveChanges(node)}
        />
      </div>
      <div className="flex flex-col">
        <Label htmlFor="entityName">Fields</Label>
        {values.inputFields?.map((field, index) => (
          <div key={index} className="flex items-center gap-2 mb-2">
            <Input
              id={"fieldName" + index}
              defaultValue={field.name}
              onChange={(e) => {
                field.name = e.currentTarget.value;
              }}
              onBlur={() => saveChanges(node)}
            />
            <Dropdown
              id={"fieldType" + index}
              defaultValue={field.type === 1 ? "Integer" : "String"}
              onOptionSelect={(e, data) => {
                node.data.inputFields[index].type = Number(data.optionValue);
              }}
              onBlur={() => saveChanges(node)}
            >
              <Option value={String(FieldType.STRING)}>String</Option>
              <Option value={String(FieldType.INTEGER)}>Integer</Option>
            </Dropdown>
            <Button
              onClick={() => {
                node.data.inputFields = node.data.inputFields.filter(
                  (_, i) => i !== index
                );
                setValue("inputFields", node.data.inputFields);
                saveChanges(node);
              }}
              icon={<HiOutlineTrash className="h-4 w-4" />}
            />
          </div>
        ))}
        <Button
          onClick={() => {
            setValue("inputFields", [
              ...(node.data.inputFields || []),
              new FieldDefinition({
                name: "undefined",
                type: FieldType.STRING,
              }),
            ]);
          }}
          icon={<HiPlus className="w-4 h-4" />}
        >
          Add Field
        </Button>
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

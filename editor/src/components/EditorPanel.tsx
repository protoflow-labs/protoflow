import { EntityData } from "@/blocks/EntityBLock";
import { FunctionData } from "@/blocks/FunctionBlock";
import { InputData } from "@/blocks/InputBlock";
import {
  Button,
  Card,
  Divider,
  Dropdown,
  Input,
  Label,
  Option,
  Select,
} from "@fluentui/react-components";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { HiOutlineTrash, HiPlus } from "react-icons/hi2";
import { Node, useOnSelectionChange } from "reactflow";
import { FieldDefinition, FieldType } from "../../rpc/block_pb";

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
    case "message":
      return <InputEditor node={props.node} />;
    case "entity":
      return <EntityEditor node={props.node} />;
    case "function":
      return <FunctionEditor node={props.node} />;
    default:
      return null;
  }
}

function InputEditor({ node }: { node: Node<InputData> }) {
  const { watch, setValue, register } = useForm({
    defaultValues: {
      name: node.data.name,
      fields: node.data.fields,
    },
  });

  const values = watch();

  return (
    <div className="flex flex-col gap-2 p-3">
      <div className="flex flex-col">
        <Label htmlFor="inputName">Name</Label>
        <Input id="inputName" {...register("name")} />
      </div>
      <div className="flex flex-col">
        <Label htmlFor="entityName">Fields</Label>
        {values.fields?.map((field, index) => (
          <div key={index} className="flex items-center gap-2 mb-2">
            <Input
              id={"fieldName" + index}
              defaultValue={field.name}
              onChange={(e) => {
                field.name = e.currentTarget.value;
              }}
            />
            <Dropdown
              id={"fieldType" + index}
              defaultValue={field.type === 1 ? "Integer" : "String"}
              onOptionSelect={(e, data) => {
                node.data.fields[index].type = Number(data.optionValue);
              }}
            >
              <Option value={String(FieldType.STRING)}>String</Option>
              <Option value={String(FieldType.INTEGER)}>Integer</Option>
            </Dropdown>
            <Button
              onClick={() => {
                setValue(
                  "fields",
                  values.fields.filter((_, i) => i !== index)
                );
              }}
              icon={<HiOutlineTrash className="h-4 w-4" />}
            />
          </div>
        ))}
        <Button
          onClick={() => {
            setValue("fields", [
              ...(node.data.fields || []),
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
      <Divider />
      <div className="flex items-center gap-2">
        <Button onClick={() => {}}>Cancel</Button>
        <Button appearance="primary">Save</Button>
      </div>
    </div>
  );
}

function EntityEditor(props: { node: Node<EntityData> }) {
  return (
    <div className="flex flex-col gap-2 p-3">
      <div className="flex flex-col">
        <Label htmlFor="entityName">Name</Label>
        <Input
          id="entityName"
          defaultValue={props.node.data.name || ""}
          onChange={(e) => {
            // props.node.data.name = e.target.value;
          }}
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
            // props.node.data.table = e.currentTarget.value;
          }}
        />
      </div>
    </div>
  );
}

function FunctionEditor(props: { node: Node<FunctionData> }) {
  return (
    <div className="flex flex-col gap-2 p-3">
      <div className="flex flex-col">
        <Label htmlFor="entityName">Name</Label>
        <Input
          id="entityName"
          defaultValue={props.node.data.name || ""}
          onChange={(e) => {
            props.node.data.name = e.target.value;
          }}
        />
      </div>
      <div className="flex flex-col">
        <Label htmlFor="entityLanguage">Language</Label>
        <Select
          onChange={(e) => {
            // props.node.data.language = e.currentTarget.value;
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
